package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"whats-gpt/whatsgpt/twilio"
)

func SendCompletion(request CompletionRequest) (*http.Response, error) {
	reqJson, err := json.Marshal(request)
	if err != nil {
		log.Panic("fail to parse json")
		return nil, err
	}
	log.Print("chatgpt request:", string(reqJson))

	httpRequest, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(reqJson),
	)
	AddRequiredHeaders(httpRequest)
	httpResponse, err := http.DefaultClient.Do(httpRequest) //envia a httpRequest e recebe a httpResponse
	if err != nil {
		log.Println("fail to send request to chatgpt", err)
		return nil, err
	}
	return httpResponse, nil
}

func ParseCompletionResponse(response http.Response) (CompletionResponse, error) {
	defer response.Body.Close()
	var completionResponse CompletionResponse
	if err := json.NewDecoder(response.Body).Decode(&completionResponse); err != nil {
		log.Println("fail to parse chatgpt response", err)
		return CompletionResponse{}, err
	}
	return completionResponse, nil
}

func StreamMessageWithResponse(httpResponse http.Response, numberToSendMessage string) {
	fmt.Println("numberToSendMessage:", numberToSendMessage)
	defer httpResponse.Body.Close()
	fmt.Println("httpResponse.StatusCode:", httpResponse.StatusCode)
	scan := bufio.NewScanner(httpResponse.Body)
	var messageBody string
	for scan.Scan() {
		line := scan.Text()
		if line != "" && !strings.Contains(line, "[DONE]") {
			var response CompletionResponse
			err := json.Unmarshal([]byte(strings.Replace(line, "data: ", "", 1)), &response)
			if err != nil {
				log.Println("fail to parse chatgpt response", err)
				continue
			}
			choice := response.Choices[0]
			if choice.FinishReason == "stop" || choice.FinishReason == "length" {
				twilio.SendMessage(numberToSendMessage, messageBody)
				break
			}
			currentContent := choice.Delta.Content
			if strings.Contains(currentContent, "\n\n") && messageBody != "" {
				if strings.Contains(currentContent, "```") {
					messageBody += currentContent
					if strings.Count(messageBody, "```")%2 == 0 { //code block acabou

						twilio.SendMessage(numberToSendMessage, messageBody)
						messageBody = ""
						//send message to whatsapp
					}
				} else if !strings.Contains(messageBody, "```") {
					messageBody += strings.Replace(currentContent, "\n", "", -1)

					twilio.SendMessage(numberToSendMessage, messageBody)
					messageBody = ""
				} else {
					messageBody += currentContent
				}
			} else {
				messageBody += currentContent
			}
		}
	}
}
