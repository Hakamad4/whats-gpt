package twilio

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Message struct {
	From string `json:"From"`
	To   string `json:"To"`
	Body string `json:"Body"`
}

type Request struct {
	Body string `json:"Body"`
	From string `json:"From"`
}

var (
	AccountSid = "ACf30341ecb84b114d6aa99bd6b0c8424d"
	AuthToken  = "50350d8a94e126ed1ee91facca72a7bc"
	MyNumber   = "whatsapp:+14155238886"
)

func SendMessage(from string, body string) {
	twilioMessage := url.Values{
		"From": {MyNumber},
		"To":   {from},
		"Body": {body},
	}

	httpRequest, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/ACf30341ecb84b114d6aa99bd6b0c8424d/Messages.json", strings.NewReader(twilioMessage.Encode()))
	if err != nil {
		log.Panic("fail to create request")
		return
	}
	httpRequest.SetBasicAuth(AccountSid, AuthToken)
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResponse, err := http.DefaultClient.Do(httpRequest)

	if err != nil {
		log.Panic("fail to send request")
		return
	}

	defer httpResponse.Body.Close()
}

func BuildTwilioRequestFromUrlParams(urlQuery string) (Request, error) {
	fmt.Println("urlQuery from twilio:", urlQuery)
	data, err := url.ParseQuery(urlQuery)
	if err != nil {
		return Request{}, err
	}
	var twilioRequest Request
	if data.Has("Body") && data.Has("From") {
		twilioRequest.Body = data.Get("Body")
		twilioRequest.From = data.Get("From")
	} else {
		return Request{}, errors.New("invalid request")
	}
	return twilioRequest, nil
}
