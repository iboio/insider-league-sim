package appContext

import (
	"league-sim/internal/layers/adapt"
	"league-sim/internal/layers/infra"
	"league-sim/internal/layers/service"
)

type AppContext struct {
	Infra   *infra.Infra
	Adapt   *adapt.Adapt
	Service *service.Service
}

func BuildAppContext() (AppContext, error) {
	infraLayer := infra.BuildInfraLayer()
	adaptLayer := adapt.BuildAdaptLayer(infraLayer)
	serviceLayer := service.BuildServiceLayer(adaptLayer)
	return AppContext{
		Infra:   infraLayer,
		Adapt:   adaptLayer,
		Service: serviceLayer,
	}, nil
}
