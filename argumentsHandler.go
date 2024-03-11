package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tidwall/gjson"
)

// ArgumentsHandler: The handler of prepare arguments
type ArgumentsHandler struct {
	// 合成复用Next
	Next
}

// Do prepare arguments
func (a *ArgumentsHandler) Do(ctx context.Context) (context.Context, error) {
	log.Println("prepare transferring arguments")
	event := new(FacebookEvent)
	fbReq := ctx.Value(ContextKey("fbReq")).(events.APIGatewayProxyRequest)

	log.Println(fbReq.Body)

	if gjson.Get(fbReq.Body, "entry.0.messaging.0.message.is_echo").Exists() {
		return ctx, errors.New("ignore echo messages")
	}

	event.Field = "messages"
	event.Id = fbReq.RequestContext.RequestID
	event.Message = gjson.Get(fbReq.Body, "entry.0.messaging.0.message.text").String()
	event.CustomerId = gjson.Get(fbReq.Body, "entry.0.messaging.0.sender.id").String()

	return context.WithValue(ctx, ContextKey("fbEvent"), *event), nil
}
