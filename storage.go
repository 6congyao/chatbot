package main

type Storage interface {
	store(event FacebookEvent, replyMsg string) error
}
