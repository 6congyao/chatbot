package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func botHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if event.HTTPMethod == "GET" {
		return verifyHandler(ctx, event)
	}

	log.Println(event.Body)
	var facebookEvent, _err = makeFacebookEvent(event)
	var response events.APIGatewayProxyResponse
	if _err != nil {
		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       event.Body,
		}
		return response, nil
	}

	replyMsg, _err := parseEvent(facebookEvent)
	if _err != nil {
		log.Println(_err)
		response = events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       _err.Error(),
		}
	} else {
		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       event.Body,
		}
		var stor Storage = new(Ddb)
		_err = stor.store(facebookEvent, replyMsg)
		if _err != nil {
			log.Println(_err)
		}

	}

	return response, nil
}

func verifyHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var queryParameters = event.QueryStringParameters
	if queryParameters["hub.verify_token"] == "CONGYAO_VERIFIY_TOKEN" && queryParameters["hub.mode"] == "subscribe" {
		response := events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       queryParameters["hub.challenge"],
		}
		return response, nil
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 403,
		Body:       "Missing/invalid token",
	}

	return response, nil
}

func main() {
	lambda.Start(botHandler)
}
