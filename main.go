package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type FormData struct {
	SubmissionDate string `json:"Submission Date"`
	FormTitle      string `json:"Form Title"`
	Mail           string `json:"E-Mail"`
}

func HandleLambdaEvent(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("Event: %+v", req)

	data := FormData{}
	err := json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("error unmarshaling json: %v", err),
		}, nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusTeapot,
		Body:       fmt.Sprintf("Sent event: %+v", data),
	}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
