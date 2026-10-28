package core_http_types

import (
	"encoding/json"
	"security/internal/core/domains"
)

type Put[T any] struct {
	domains.Put[T]
}

func (n *Put[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = value

	return nil
}

func (n *Put[T]) ToDomain() domains.Put[T] {
	return domains.Put[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
