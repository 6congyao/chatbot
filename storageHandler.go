package main

import (
	"context"
	"log"
)

// StorageHandler: Handler for storing contents
type StorageHandler struct {
	Next
}

func (s *StorageHandler) Do(ctx context.Context) (context.Context, error) {
	log.Println("prepare storing contents")
	fbReply := ctx.Value(ContextKey("fbReply")).(string)
	fbEvent := ctx.Value(ContextKey("fbEvent")).(FacebookEvent)
	var stor Storage = new(Ddb)
	if err := stor.store(fbEvent, fbReply); err != nil {
		return ctx, err
	}

	return ctx, nil
}
