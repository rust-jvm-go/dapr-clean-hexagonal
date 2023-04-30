package json

import (
	"fmt"
	"github.com/cockroachdb/errors"
	json "github.com/json-iterator/go"
	"url-shortener/domain"
)

type RedirectJsonSerializer struct{}

func (j *RedirectJsonSerializer) Decode(input []byte) (*domain.Redirect, error) {
	redirect := &domain.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "jsonSerializer.Redirect.Decode")
	}

	return redirect, nil
}

func (j *RedirectJsonSerializer) Encode(input *domain.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "jsonSerializer.Redirect.Encode")
	}

	return rawMsg, nil
}

func (j *RedirectJsonSerializer) Info() string {
	return fmt.Sprintf("From RedirectJsonSerializer")
}
