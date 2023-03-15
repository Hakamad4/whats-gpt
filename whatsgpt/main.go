package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	. "github.com/aws/aws-lambda-go/events"
)

func createChatGptRequest(query string) ChatGptRequest {
	req := ChatGptRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{"user", query},
		},
		MaxTokens: 150,
	}
	return req
}

func GenerateGptText(query string) (string, error) {
	request := createChatGptRequest(query)

	reqJson, err := json.Marshal(request)
	if err != nil {
		log.Panic("fail to parse json")
		return "", err
	}
	log.Print("chatgpt request:", string(reqJson))

	httpRequest := createHttpRequest(err, reqJson)
	httpResponse, err := http.DefaultClient.Do(httpRequest) //envia a httpRequest e recebe a httpResponse
	if err != nil {
		log.Println("fail to send request to chatgpt", err)
		return "", err
	}
	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	log.Print("chatgpt response:", string(responseBody))
	var response Response
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}

func createHttpRequest(err error, reqJson []byte) *http.Request {
	httpRequest, err := http.NewRequest( //isso builda uma httpRequest
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(reqJson),
	)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer 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")
	return httpRequest
}

func parseBase64(str string) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(decodeString), nil
}

func getBodyFromUrlParams(urlQuery string) (string, error) {
	data, err := url.ParseQuery(urlQuery)
	if err != nil {
		return "", err
	}
	if data.Has("Body") {
		return data.Get("Body"), nil
	}
	return "", errors.New("body not found")
}

func process(request APIGatewayProxyRequest) (APIGatewayProxyResponse, error) {
	parsed, err := parseBase64(request.Body)
	result, err := getBodyFromUrlParams(parsed)
	log.Println("body from twilio: ", result)
	if err != nil {
		return APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	text, err := GenerateGptText(result)
	fmt.Println("chatgpt generated text ", text)
	if err != nil {
		return APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	return APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       text,
	}, nil
}

func main() {
	GenerateGptText("Whats is a GPT-3?")
}
