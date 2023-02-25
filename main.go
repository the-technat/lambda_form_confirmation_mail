package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-mail/mail"

	log "github.com/sirupsen/logrus"
)

// FormData represents the JSON sent in the request body
type FormData map[string]interface{}

// Config represent the keys that should be present in the AWS SecretsManager secret
type Config struct {
	MailHost string `json:"MAIL_HOST"`
	MailPort int    `json:"MAIL_PORT"`
	MailFrom string `json:"MAIL_FROM"`
	MailUser string `json:"MAIL_USER"`
	MailPW   string `json:"MAIL_PW"`
	MailMsg  string `json:"MAIL_MSG"`
  MailCopy bool   `json:"MAIL_COPY"`
}

func main() {
  // main is the handler that was configured to be called when the lambda should start
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
  log.Infof("new incoming event: %+v", req)

	// read config from secret
  config := Config{}
	awsSecret, err := readSecret(os.Getenv("SECRET"))
	if err != nil {
    return nil, fmt.Errorf("couldn't read secret: %v", err) // by returning an error, the lambda will return 500 to the client and mark the lambda as failed internally
	}
	err = json.Unmarshal([]byte(awsSecret), &config)
	if err != nil {
    return nil, fmt.Errorf("couldn't unmarshal secret into known structure: %v", err)
	}
  log.Info("successfully read config from secret")

  // parse incoming json data
	data := FormData{}
	err = json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
    // in case the unmarshal failed in can only be because the client send a bad json
    // so we inform the client about the bad request and don't just error out
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("error unmarshaling json: %v", err),
		}, nil
	}
  log.Infof("successfully read posted data: %+v", data)

  // parse msg template from config
  msgT := template.Must(template.New("msg").Parse(config.MailMsg))
  msg := &bytes.Buffer{}
  err = msgT.Execute(msg, data)
  if err != nil {
    return &events.APIGatewayProxyResponse{
      StatusCode: http.StatusInternalServerError,
      Body: fmt.Sprintf("failed parsing your message template: %v", err),
    }, nil
  }

	// and send a mail with the parsed message
	m := mail.NewMessage()
	m.SetHeader("From", config.MailFrom)
  if config.MailCopy {
    m.SetHeader("To", data["E-Mail"].(string), config.MailFrom)
  } else {
    m.SetHeader("To", data["E-Mail"].(string))
  }
	m.SetHeader("Subject", data["Form Title"].(string))
	m.SetBody("text/html", msg.String())

	d := mail.NewDialer(config.MailHost, config.MailPort, config.MailUser, config.MailPW)
	if err := d.DialAndSend(m); err != nil {
    return nil, fmt.Errorf("coulnd't send mail: %v", err)
  }

  log.Infof("new mail for %s sent to %s", data["Form Title"].(string), data["E-Mail"].(string))
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "sent mail",
	}, nil

}

// readSecrets returns the content of the specified secret
func readSecret(secretName string) (string, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return "", err
	}

	svc := secretsmanager.NewFromConfig(config)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return "", err
	}

	// Decrypts secret using the associated KMS key.
	return *result.SecretString, nil
}