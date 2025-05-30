package service

import (
	"league-sim/internal/layers/adapt"
	"league-sim/internal/league"
	leagueInterfaces "league-sim/internal/league/interfaces"
	"league-sim/internal/predict"
	predictInterfaces "league-sim/internal/predict/interfaces"
	"league-sim/internal/simulation"
	simInterfaces "league-sim/internal/simulation/interfaces"
)

type ServiceInterface interface {
	SimulationService() simInterfaces.SimulationServiceInterface
	LeagueService() leagueInterfaces.LeagueServiceInterface
	PredictionService() predictInterfaces.PredictServiceInterface
}

type Service struct {
	simulationService simInterfaces.SimulationServiceInterface
	leagueService     leagueInterfaces.LeagueServiceInterface
	predictionService predictInterfaces.PredictServiceInterface
}

func BuildServiceLayer(adapt *adapt.Adapt) *Service {
	newSimulationService := simulation.NewSimulationService(adapt)
	newLeagueService := league.NewLeagueService(adapt)
	newPredictionService := predict.NewPredictService(adapt)

	return &Service{
		simulationService: newSimulationService,
		leagueService:     newLeagueService,
		predictionService: newPredictionService,
	}
}

func (s *Service) SimulationService() simInterfaces.SimulationServiceInterface {
	return s.simulationService
}

func (s *Service) LeagueService() leagueInterfaces.LeagueServiceInterface {
	return s.leagueService
}

func (s *Service) PredictionService() predictInterfaces.PredictServiceInterface {
	return s.predictionService
}
