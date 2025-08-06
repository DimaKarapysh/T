package service

import "T/internal/service/adapters"

type serviceFactory struct {
	deps Params
}

func newServiceFactory(deps Params) *serviceFactory {
	return &serviceFactory{deps: deps}
}

func (receiver *serviceFactory) Build() adapters.SubService {
	return NewSubService(SubParams{

		receiver.deps.Logger,
		receiver.deps.Repository,
		//	receiver.deps.SessionRepo,
		//	receiver.deps.Token,
		receiver.deps.Config})
}
