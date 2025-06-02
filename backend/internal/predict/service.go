package predict

import (
	"fmt"
	"league-sim/config"
	adaptInterface "league-sim/internal/layers/adapt/interfaces"
	"league-sim/internal/league"
	"league-sim/internal/models"
	"league-sim/internal/predict/interfaces"
	repoInterfaces "league-sim/internal/repositories/interfaces"
)

type Predict struct {
	matchResultRepo repoInterfaces.MatchesRepository
	standingsRepo   repoInterfaces.StandingRepository
	teamsRepo       repoInterfaces.TeamsRepository
}

func NewPredictService(adapt adaptInterface.AdaptInterface) interfaces.PredictServiceInterface {
	return &Predict{
		matchResultRepo: adapt.MatchesRepository(),
		standingsRepo:   adapt.StandingsRepository(),
		teamsRepo:       adapt.TeamRepository(),
	}
}

func (a *Predict) PredictChampionShipSession(leagueId string) ([]models.PredictedStanding, error) {
	standings, err := a.standingsRepo.GetStandings(leagueId)
	if err != nil {
		fmt.Println("Error getting league standings:", err)
		return nil, err
	}
	teams, err := a.teamsRepo.GetTeams(leagueId)
	if err != nil {
		fmt.Println("Error getting teams:", err)
		return nil, err
	}

	var teamMap = make(map[string]models.Team)

	for _, team := range teams {
		teamMap[team.TeamName] = team
	}

	if len(standings) == 0 {
		return []models.PredictedStanding{}, nil
	}

	leaderPoints := findLeaderPoints(standings)
	totalTeams := len(standings)
	totalScore := 0.0

	type scoredTeam struct {
		TeamName string
		Score    float64
		Adjusted float64
		Points   int
		Strength float64
		Played   int
		Goals    int
		Against  int
	}

	var scored []scoredTeam

	for _, s := range standings {
		str := league.CalculateStrength(teamMap[s.TeamName])
		score := float64(s.Points)*config.WeightPoints + str*config.WeightsStrength
		scored = append(
			scored, scoredTeam{
				TeamName: s.TeamName,
				Score:    score,
				Points:   s.Points,
				Strength: str,
				Played:   s.Played,
				Goals:    s.Goals,
				Against:  s.Against,
			})
		totalScore += score
	}

	avgScore := totalScore / float64(len(scored))
	totalAdjusted := 0.0

	allRemainingMatchesZero := true
	for _, s := range scored {
		remainingMatches := (totalTeams - 1) - s.Played
		if remainingMatches > 0 {
			allRemainingMatchesZero = false
			break
		}
	}

	if allRemainingMatchesZero {

		leader := scored[0]
		for _, s := range scored[1:] {
			if s.Points > leader.Points {
				leader = s
			} else if s.Points == leader.Points {
				if (s.Goals - s.Against) > (leader.Goals - leader.Against) {
					leader = s
				}
			}
		}

		var result []models.PredictedStanding
		for _, s := range scored {
			eliminated := s.TeamName != leader.TeamName
			odds := 0.0
			if !eliminated {
				odds = 100.0
			}
			result = append(
				result, models.PredictedStanding{
					TeamName:   s.TeamName,
					Points:     s.Points,
					Strength:   s.Strength,
					Odds:       odds,
					Eliminated: eliminated,
				})
		}

		return result, nil
	}

	for i, s := range scored {
		remainingMatches := (totalTeams - 1) - s.Played
		maxPossiblePoints := s.Points + remainingMatches*3

		adjusted := scored[i].Score / avgScore
		if maxPossiblePoints < leaderPoints {
			adjusted = 0
		}
		scored[i].Adjusted = adjusted
		totalAdjusted += adjusted
	}

	var result []models.PredictedStanding

	for _, s := range scored {
		eliminated := s.Adjusted == 0
		odds := 0.0
		if !eliminated && totalAdjusted > 0 {
			odds = s.Adjusted / totalAdjusted * 100
		}
		result = append(
			result, models.PredictedStanding{
				TeamName:   s.TeamName,
				Points:     s.Points,
				Strength:   s.Strength,
				Odds:       odds,
				Eliminated: eliminated,
			})
	}
	return result, nil
}

func findLeaderPoints(standings []models.Standings) int {
	maxPoint := 0

	for _, s := range standings {
		if s.Points > maxPoint {
			maxPoint = s.Points
		}
	}

	return maxPoint
}
