package msgpack

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/vmihailenco/msgpack/v5"
	"url-shortener/domain"
)

type RedirectMsgPackSerializer struct{}

func (m *RedirectMsgPackSerializer) Decode(input []byte) (*domain.Redirect, error) {
	redirect := &domain.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}

func (m *RedirectMsgPackSerializer) Encode(input *domain.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}

func (m *RedirectMsgPackSerializer) Info() string {
	return fmt.Sprintf("From RedirectMsgPackSerializer")
}
