package kiwi

import "errors"

var (
	ErrUnauthorized = errors.New("kiwi: unauthorized")
	ErrNotFound     = errors.New("kiwi: not found")
)

type RPCError struct {
	Code    int
	Message string
	Data    any
}

func (e *RPCError) Error() string {
	return e.Message
}
