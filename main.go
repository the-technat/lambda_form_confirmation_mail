package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-mail/mail"

	log "github.com/sirupsen/logrus"
)

type FormData map[string]interface{}

type Config struct {
	MailHost string `json:"MAIL_HOST"`
	MailFrom string `json:"MAIL_FROM"`
	MailUser string `json:"MAIL_USER"`
	MailPW   string `json:"MAIL_PW"`
	MailPort int    `json:"MAIL_PORT"`
	MailMsg  string `json:"MAIL_MSG"`
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Debugf("event: %+v", req)

	// first get the config
	awsSecret, err := readSecret(os.Getenv("SECRET"))
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("could not retrive secret %s: %v", os.Getenv("SECRET"), err),
		}, nil
	}
	config := Config{}
	err = json.Unmarshal([]byte(awsSecret), &config)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("config secret has unknown structure: %v", err),
		}, nil

	}
	log.Printf("Using user %s on %s:%d", config.MailUser, config.MailHost, config.MailPort)

	// then take a look at the posted data
	data := FormData{}
	err = json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("error unmarshaling json: %v", err),
		}, nil
	}

	// and send a mail with this data
	m := mail.NewMessage()
	m.SetHeader("From", config.MailFrom)
	log.Printf("New form submitted %s for %s", data["Form Title"].(string), data["E-Mail"].(string))
	m.SetHeader("To", data["E-Mail"].(string), config.MailFrom) // always send a copy to the sender
	m.SetHeader("Subject", data["Form Title"].(string))
	m.SetBody("text/html", fmt.Sprintf("Hello <b>%s</b></br>%s</br>%s", data["Name"].(string), config.MailMsg, req.Body))
	d := mail.NewDialer(config.MailHost, config.MailPort, config.MailUser, config.MailPW)
	if err := d.DialAndSend(m); err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("couldn't send mail: %v", err),
		}, nil
	}

	log.Println("email sent successfully")
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "sent mail",
	}, nil

}

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

// func BuildMessage(mail Mail) string {
// 	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
// 	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
// 	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
// 	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
// 	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)
// 	return msg
// }
