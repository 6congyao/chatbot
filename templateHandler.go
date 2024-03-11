package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

var templates = map[string]string{
	"thank":                 "Dear %s, Your unwavering support and trust in our products/services mean the world to us. We are truly grateful for the opportunity to serve you and for the strong partnership we have built.",
	"transaction completed": "We sincerely appreciate your business and the trust you have placed in us. If you have any further questions or need assistance in the future, please don't hesitate to reach out. We value your satisfaction and look forward to serving you again.",
}

// TemplateHandler: Handler for templates message
type TemplateHandler struct {
	Next
}

func (t *TemplateHandler) Do(ctx context.Context) (context.Context, error) {
	log.Println("prepare template reply messages")
	fbEvent := ctx.Value(ContextKey("fbEvent")).(FacebookEvent)
	message := ""
	if fbEvent.Field == "messages" {
		if strings.Contains(fbEvent.Message, "thank") {
			message = fmt.Sprintf(templates["thank"], fbEvent.CustomerId)
		}
		if strings.Contains(fbEvent.Message, "transaction completed") {
			message = templates["transaction completed"]
		}
	}
	if message == "" || len(message) == 0 {
		return ctx, errors.New("no need to reply")
	}

	return context.WithValue(ctx, ContextKey("fbReply"), message), nil
}
