package main

import (
	"context"
	"log"
)

// MessageSender: Handler for sending messages
type MessageSender struct {
	Next
}

func (m *MessageSender) Do(ctx context.Context) (context.Context, error) {
	log.Println("prepare sending reply messages")
	fbReply := ctx.Value(ContextKey("fbReply")).(string)
	fbEvent := ctx.Value(ContextKey("fbEvent")).(FacebookEvent)
	if err := sendMessage(fbReply, fbEvent.CustomerId, "RESPONSE"); err != nil {
		return ctx, err
	}

	return ctx, nil
}
