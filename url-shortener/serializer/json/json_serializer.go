package json

import (
	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
	"url-shortener/domain"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*domain.Redirect, error) {
	redirect := &domain.Redirect{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}

func (r *Redirect) Encode(input *domain.Redirect) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}
