package msgpack

import (
	"github.com/cockroachdb/errors"
	"github.com/vmihailenco/msgpack/v5"
	"url-shortener/domain"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*domain.Redirect, error) {
	redirect := &domain.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}

func (r *Redirect) Encode(input *domain.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}
