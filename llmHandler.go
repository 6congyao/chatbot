package main

import (
	"context"
	"log"
	"time"
)

// MessageSender: Handler for sending messages
type LLMHandler struct {
	Next
}

func (l *LLMHandler) Do(ctx context.Context) (context.Context, error) {
	log.Println("prepare sending to LLM service")
	fbEvent := ctx.Value(ContextKey("fbEvent")).(FacebookEvent)

	chResult := make(chan string, 1)
	ctx = context.WithValue(ctx, ContextKey("chLlmRes"), chResult)

	go func() {
		res, _ := callLLM(ctx, fbEvent.Message)
		chResult <- res
	}()

	return ctx, nil
}

func callLLM(ctx context.Context, prompt string) (string, error) {
	log.Println("got prompt:" + prompt)
	time.Sleep(5 * time.Second)

	return "llm gives replay here", nil
}
