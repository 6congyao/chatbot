package main

import (
	"context"
	"log"
)

// ArgumentsHandler: The handler of prepare arguments
type AggregationHandler struct {
	// 合成复用Next
	Next
}

// Do prepare arguments
func (a *AggregationHandler) Do(ctx context.Context) (context.Context, error) {
	log.Println("aggregation starting")
	chRes := ctx.Value(ContextKey("chLlmRes")).(chan string)
	llmRes := <-chRes
	log.Println("got llm reply:" + llmRes)
	return ctx, nil
}
