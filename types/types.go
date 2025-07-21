package types

import "errors"

var ErrKeyNotFound = errors.New("Key not found")

type KeyValue struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}
