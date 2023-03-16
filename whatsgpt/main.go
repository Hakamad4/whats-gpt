package main

import (
	"encoding/base64"
	"errors"
	. "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"whats-gpt/whatsgpt/chatgpt/client"
	"whats-gpt/whatsgpt/twilio"
)

func streamMessage(twilioRequest twilio.Request) error {
	request := client.CompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []client.Message{
			{"user", twilioRequest.Body},
		},
		Stream:    true,
		MaxTokens: 500,
	}
	completion, err := client.SendCompletion(request)
	if err != nil {
		return err
	}
	if completion.StatusCode != http.StatusOK {
		return errors.New("ChatGPT returned a non 200 status code")
	}

	if completion.Header.Values("Content-Type")[0] == "text/event-stream" {
		client.StreamMessageWithResponse(*completion, twilioRequest.From)
	} else {
		response, err := client.ParseCompletionResponse(*completion)
		if err != nil {
			return err
		}
		twilio.SendMessage(twilioRequest.From, response.Choices[0].Message.Content)
	}

	return nil
}

func parseBase64(str string) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(decodeString), nil
}

func process(request APIGatewayProxyRequest) (APIGatewayProxyResponse, error) {
	parsed, err := parseBase64(request.Body)
	result, err := twilio.BuildTwilioRequestFromUrlParams(parsed)
	log.Println("body from twilio: ", result)
	if err != nil {
		return APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	err = streamMessage(result)
	if err != nil {
		return APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	return APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(process)
}
