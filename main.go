package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"

	log "github.com/sirupsen/logrus"
)

type FormData map[string]interface{}

type Secret struct {
	MailHost string `json:"MAIL_HOST"`
	MailFrom string `json:"MAIL_FROM"`
	MailUser string `json:"MAIL_USER"`
	MailPW   string `json:"MAIL_PW"`
}

type Mail struct {
	Sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Debugf("event: %+v", req)

	// first get the secrets
	secret, err := readSecret(os.Getenv("SECRET"))
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("could not retrive secret %s: %v", os.Getenv("SECRET"), err),
		}, nil
	}
	secretData := Secret{}
	err = json.Unmarshal([]byte(secret), &secretData)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("secret has unknown structure: %v", err),
		}, nil

	}

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
	to := []string{data["E-Mail"].(string)}
	log.Printf("New mail for %s", to[0])

	msg := "Hello geeks!!!"
	log.Printf("Using user %s on %s", secretData.MailUser, secretData.MailHost)
	auth := smtp.PlainAuth("", secretData.MailUser, secretData.MailPW, secretData.MailHost)
	err = smtp.SendMail(secretData.MailHost, auth, secretData.MailFrom, to, []byte(msg))
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("couldn't send mail: %v", err),
		}, nil
	}

	fmt.Println("Email sent successfully")
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
