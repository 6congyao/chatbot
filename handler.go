package main

import "context"

type ContextKey string

// Handler: Base handler
type Handler interface {
	Do(ctx context.Context) (context.Context, error)
	SetNext(h Handler) Handler
	Run(ctx context.Context) error
}

// Next: An abstract structure that can be synthesized and reused
type Next struct {
	nextHandler Handler
}

func (n *Next) SetNext(h Handler) Handler {
	n.nextHandler = h
	return h
}

func (n *Next) Run(ctx context.Context) (err error) {
	if n.nextHandler != nil {
		var nctx context.Context
		if nctx, err = (n.nextHandler).Do(ctx); err != nil {
			return
		}
		return (n.nextHandler).Run(nctx)
	}
	return
}

// NullHandler: A null handler for the first implementaion
type NullHandler struct {
	Next
}

func (n *NullHandler) Do(ctx context.Context) (err error) {
	// do nothing...
	return
}
