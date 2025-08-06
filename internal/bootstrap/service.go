package bootstrap

import (
	v01 "T/internal/api/http/v01"
	"T/internal/repo/ports"
	"T/internal/repo/postgres"
	"T/internal/service"
	"T/internal/service/adapters"
	"go.uber.org/fx"
)

func ServiceModules() fx.Option {
	return fx.Options(
		fx.Provide(

			//tokens
			//token.NewAccessToken,
			//token.NewRefreshToken,
			//NewToken,
			//middleware
			//middlewarepackage.NewJWTMiddleware,

			//repository
			fx.Annotate(postgres.NewRepoSub, fx.As(new(ports.Repository))),
			//fx.Annotate(redis.NewSessionRepository, fx.As(new(ports.Session))),
			//service
			service.NewService,
			fx.Annotate(
				func(s *service.Service) adapters.SubService {
					return s.Auth()
				},
				fx.As(new(adapters.SubService)),
			),
			//handler
			v01.NewHandler,
		),
	)
}
