package events

import (
	"github.com/adm87/finch-core/errors"
	"github.com/adm87/finch-core/types"
)

var (
	ErrMessageHandleAlreadyExists = errors.NewDuplicateError("message handler already exists")
)

type MessageHandler[T Message] interface {
	HandleMessage(msg T) error
}

type MessageBus[T Message] struct {
	handlers types.HashSet[MessageHandler[T]]
}

func NewMessageBus[T Message]() *MessageBus[T] {
	return &MessageBus[T]{
		handlers: make(types.HashSet[MessageHandler[T]]),
	}
}

func (bus *MessageBus[T]) Subscribe(handler ...MessageHandler[T]) error {
	for _, h := range handler {
		if bus.handlers.Contains(h) {
			return ErrMessageHandleAlreadyExists
		}
		bus.handlers.Add(h)
	}
	return nil
}

func (bus *MessageBus[T]) Unsubscribe(handler ...MessageHandler[T]) error {
	for _, h := range handler {
		bus.handlers.Remove(h)
	}
	return nil
}

func (bus *MessageBus[T]) Send(msg T) error {
	for handler := range bus.handlers {
		if err := handler.HandleMessage(msg); err != nil {
			return err
		}
	}
	return nil
}
