package xsetup

import (
	"gopkg.in/dealancer/validate.v2"
	"url-shortener/domain"
)

type Environment struct {
	DaprClient         *DaprClient                 `validate:"nil=false"`
	InitConfig         *InitConfig                 `validate:"nil=false"`
	RedirectRepository *domain.IRedirectRepository `validate:"nil=false"`
	RedirectSerializer *domain.IRedirectSerializer // `validate:"nil=false"`
}

func NewEnvironment(daprClient DaprClient, initConfig InitConfig, repository domain.IRedirectRepository) *Environment {
	return &Environment{
		DaprClient:         &daprClient,
		InitConfig:         &initConfig,
		RedirectRepository: &repository,
		// RedirectSerializer: &serializer,
	}
}

func (e *Environment) Info() string {
	if err := validate.Validate(e); err != nil {
		return "Environment not initialized."
	}
	return "Environment successfully initialized."
}
