package services

import (
	appContext "league-sim/internal/contexts/appContexts"
	"league-sim/internal/league"
	interfaces1 "league-sim/internal/league/interfaces"
	"league-sim/internal/predict"
	interfaces2 "league-sim/internal/predict/interfaces"
	"league-sim/internal/simulation"
	interfaces3 "league-sim/internal/simulation/interfaces"
)

type Service interface {
	LeagueService() interfaces1.LeagueServiceInterface
	PredictService() interfaces2.PredictServiceInterface
	SimulationService() interfaces3.SimulationServiceInterface
}

type ServiceImpl struct {
	leagueService     interfaces1.LeagueServiceInterface
	predictService    interfaces2.PredictServiceInterface
	simulationService interfaces3.SimulationServiceInterface
}

func (s *ServiceImpl) LeagueService() interfaces1.LeagueServiceInterface {

	return s.leagueService
}

func (s *ServiceImpl) PredictService() interfaces2.PredictServiceInterface {

	return s.predictService
}

func (s *ServiceImpl) SimulationService() interfaces3.SimulationServiceInterface {

	return s.simulationService
}

func BuildService(ctx appContext.AppContext) (*ServiceImpl, error) {
	newLeagueService := league.NewLeagueService(ctx)
	newPredictService := predict.NewPredictService(ctx)
	newSimulationService := simulation.NewSimulationService(ctx)
	return &ServiceImpl{
		leagueService:     newLeagueService,
		predictService:    newPredictService,
		simulationService: newSimulationService,
	}, nil
}
