package domain

import (
	"errors"
	errs "github.com/cockroachdb/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
	"time"
)

var (
	ErrRedirectNotFound = errors.New("redirect not found")
	ErrRedirectInvalid  = errors.New("redirect invalid")
)

type redirectService struct {
	redirectRepository IRedirectRepository
}

func NewRedirectService(redirectRepository IRedirectRepository) IRedirectService {
	return &redirectService{
		redirectRepository,
	}
}

func (s *redirectService) Find(code string) (*Redirect, error) {
	return s.redirectRepository.Find(code)
}

func (s *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}

	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()

	return s.redirectRepository.Store(redirect)
}
