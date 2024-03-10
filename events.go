package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tidwall/gjson"
)

type FacebookEvent struct {
	Field      string
	Message    string
	CustomerId string
	Id         string
}

func makeFacebookEvent(request events.APIGatewayProxyRequest) (FacebookEvent, error) {
	event := new(FacebookEvent)
	requestBody := request.Body
	if gjson.Get(requestBody, "entry.0.changes").Exists() {
		event.Field = gjson.Get(requestBody, "entry.0.changes.0.field").String()
		event.Message = gjson.Get(requestBody, "entry.0.changes.0.field.value.message").String()
		event.CustomerId = gjson.Get(requestBody, "entry.0.changes.0.field.value.from.id").String()
	}
	if gjson.Get(requestBody, "entry.0.messaging.0.message.is_echo").Exists() {
		return FacebookEvent{}, errors.New("ignore echo messages")
	} else {
		event.Field = "messages"
		event.Id = request.RequestContext.RequestID
		event.Message = gjson.Get(requestBody, "entry.0.messaging.0.message.text").String()
		event.CustomerId = gjson.Get(requestBody, "entry.0.messaging.0.sender.id").String()
	}

	if len(event.Field) == 0 || len(event.Message) == 0 || len(event.Message) == 0 {
		return FacebookEvent{}, errors.New("received invalid webhook response")
	}

	return *event, nil
}
