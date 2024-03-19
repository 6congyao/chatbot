package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func botHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if event.HTTPMethod == "GET" {
		return verifyHandler(ctx, event)
	}

	// init handler
	nullHandler := &NullHandler{}
	// responsibility chain
	nullHandler.SetNext(&ArgumentsHandler{}).
		SetNext(&TemplateHandler{}).
		//todo: call to ChatGPT/Claude analysing sentiment
		SetNext(&MessageSender{}).
		SetNext(&StorageHandler{})
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

func verifyHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var queryParameters = event.QueryStringParameters
	var response events.APIGatewayProxyResponse
	if queryParameters["hub.verify_token"] == "CONGYAO_VERIFIY_TOKEN" && queryParameters["hub.mode"] == "subscribe" {
		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       queryParameters["hub.challenge"],
		}
	} else {
		response = events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Missing/invalid token",
		}
	}

	return response, nil
}

// func main() {
// 	lambda.Start(botHandler)
// }
