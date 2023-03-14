package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGptRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type Response struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

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

	httpRequest := createHttpRequest(err, reqJson)
	httpResponse, err := http.DefaultClient.Do(httpRequest) //envia a httpRequest e recebe a httpResponse
	if err != nil {
		return "", err
	}
	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	log.Print("response:", string(responseBody))
	var response Response
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}
	log.Print("response:", response)
	return response.Choices[0].Message.Content, nil
}

func createHttpRequest(err error, reqJson []byte) *http.Request {
	httpRequest, err := http.NewRequest( //isso builda uma httpRequest
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(reqJson),
	)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer {your token key here}")
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

func process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	parsed, err := parseBase64(request.Body)
	result, err := getBodyFromUrlParams(parsed)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	text, err := GenerateGptText(result)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       text,
	}, nil
}

func main() {
	GenerateGptText("o que é uma mosca?")
}
