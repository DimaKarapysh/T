package service

import (
	"T/internal/config"
	"T/internal/repo/ports"
	"T/internal/service/adapters"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service struct {
	auth adapters.SubService
}

type Params struct {
	fx.In

	Logger     *zap.Logger
	Repository ports.Repository
	//SessionRepo ports.Session
	//Token       token.Token
	Config *config.Config
}

func NewService(params Params) *Service {

	serviceFactory := newServiceFactory(params)

	return &Service{
		auth: (serviceFactory).Build(),
	}
}

func (s *Service) Auth() adapters.SubService {
	return s.auth
}
