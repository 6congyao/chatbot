package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/carlmjohnson/requests"
	"github.com/ybbus/httpretry"
)

type Recipient struct {
	Id string `json:"id,omitempty"`
}

type Message struct {
	Text string `json:"text,omitempty"`
}

type Request struct {
	Recipient     Recipient `json:"recipient,omitempty"`
	Message       Message   `json:"message,omitempty"`
	MessagingType string    `json:"messaging_type,omitempty"`
}

func main() {
	localRun()
}

func botHandlerTest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if event.HTTPMethod == "GET" {
		return verifyHandler(ctx, event)
	}

	// init handler
	nullHandler := &NullHandler{}
	// responsibility chain
	nullHandler.SetNext(&ArgumentsHandler{}).
		SetNext(&TemplateHandler{}).
		SetNext(&LLMHandler{}).
		SetNext(&AggregationHandler{}).
		//todo: call to ChatGPT/Claude analysing sentiment
		SetNext(&MessageSender{})
		// SetNext(&StorageHandler{})
	// launch
	rootCtx := context.Background()
	if err := nullHandler.Run(context.WithValue(rootCtx, ContextKey("fbReq"), event)); err != nil {
		log.Println("Fail | Error:" + err.Error())
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       event.Body,
	}
	return response, nil
}

// func launch() {
// 	// 初始化空handler
// 	nullHandler := new(NullHandler)
// 	// nullHandler.SetNext(&ArgumentsHandler{})
// 	// 开始执行业务
// 	if err := nullHandler.Run(context.Background()); err != nil {
// 		// 异常
// 		fmt.Println("Fail | Error:" + err.Error())
// 		return
// 	}
// 	// 成功
// 	fmt.Println("Success")
// }

func localRun() {
	req := events.APIGatewayProxyRequest{
		Body: `{
			"object": "page",
			"entry": [
				{
					"id": "257662950760575",
					"time": 1710820811552,
					"messaging": [
						{
							"sender": {
								"id": "7436026376434696"
							},
							"recipient": {
								"id": "257662950760575"
							},
							"timestamp": 1710820811286,
							"message": {
								"mid": "m_OhMe1ddb67e5E_BubHX6QSuUu2NqmjKv-0bILXgkXTQMbxB9WbbqzZJA7SbJQKJCNyeRB3HuQeLrXRr3fc_snA",
								"text": "thank"
							}
						}
					]
				}
			]
		}`,
	}
	botHandlerTest(context.TODO(), req)
}

func SendMessage(message string, customerId string, messageType string) error {
	rcp := Recipient{
		Id: customerId,
	}
	msg := Message{
		Text: message,
	}
	req := Request{
		Recipient:     rcp,
		Message:       msg,
		MessagingType: messageType,
	}
	var res string

	cl := httpretry.NewDefaultClient() //Used for retries
	var err = requests.
		URL("https://graph.facebook.com/v19.0/257662950760575/messages").
		Param("access_token", "EAAzg2pLO8LoBOzZAXxXCEQLWGDxFlM1D2vmv3wyEcAQtTOxOLGRodkFdoggJiFritaOUKErKeVbcPDvq6AcNihwaXnI7NHAFcUi0rYTxCc1jO9AKAtSC4a76ywP7c6yB13VyDPpjOiWYw4Wrb4sHHSrptxHH2wUf9U7CXKWphifd0GjZAiSIKreMZCQS2pAWtwHZAumE").
		BodyJSON(&req).ToString(&res).
		Client(cl).
		Fetch(context.Background())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	return err
}
