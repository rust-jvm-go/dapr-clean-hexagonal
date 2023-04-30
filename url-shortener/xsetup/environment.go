package xsetup

import (
	"fmt"
	"url-shortener/domain"
)

type Environment struct {
	Dapr               *Dapr
	InitConfig         InitConfig
	RedirectRepository domain.IRedirectRepository
	RedirectSerializer domain.IRedirectSerializer
}

func NewEnvironment(dapr *Dapr, initConfig InitConfig, repository domain.IRedirectRepository) Environment {
	return Environment{
		Dapr:               dapr,
		InitConfig:         initConfig,
		RedirectRepository: repository,
	}
}

func (e Environment) Info() string {
	return fmt.Sprintf("\n%s\n%s\n", e.RedirectRepository.Info(), "")
}
