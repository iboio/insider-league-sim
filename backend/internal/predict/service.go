package predict

import (
	"fmt"

	"league-sim/config"
	appContext "league-sim/internal/contexts/appContexts"
	"league-sim/internal/league"
	"league-sim/internal/models"
)

type Predict struct {
	appCtx appContext.AppContext
}

func NewPredictService(ctx appContext.AppContext) *Predict {
	return &Predict{
		appCtx: ctx,
	}
}

func (a *Predict) PredictChampionshipPreSeason(id string) ([]models.PredictedStanding, error) {
	teams, err := a.appCtx.ActiveLeagueRepository().GetActiveLeagueTeams(id)

	if err != nil {
		fmt.Println("Error getting league data:", err)
		return nil, err
	}

	totalStrength := 0.0
	var predictTable []models.PredictedStanding

	for _, t := range teams {
		str := league.CalculateStrength(t)
		totalStrength += str
	}

	for _, t := range teams {
		str := league.CalculateStrength(t)
		predictTable = append(
			predictTable, models.PredictedStanding{
				TeamName:   t.Name,
				Strength:   str,
				Odds:       str / totalStrength * 100,
				Eliminated: false,
			})
	}

	return predictTable, nil
}

func (a *Predict) PredictChampionShipSession(id string) ([]models.PredictedStanding, error) {
	standings, err := a.appCtx.ActiveLeagueRepository().GetActiveLeaguesStandings(id)

	if err != nil {
		fmt.Println("Error getting league standings:", err)
		return nil, err
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
		Team     models.Team
	}

	var scored []scoredTeam

	for _, s := range standings {
		str := league.CalculateStrength(s.Team)
		score := float64(s.Points)*config.WeightPoints + str*config.WeightsStrength
		scored = append(
			scored, scoredTeam{
				TeamName: s.Team.Name,
				Score:    score,
				Points:   s.Points,
				Strength: str,
				Played:   s.Played,
				Goals:    s.Goals,
				Against:  s.Against,
				Team:     s.Team,
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

func FindLeader(standings []models.Standings) models.Team {
	var leader models.Team
	if len(standings) == 0 {
		return leader
	}

	leaderPoint := findLeaderPoints(standings)
	var candidates []models.Standings

	for _, s := range standings {
		if s.Points == leaderPoint {
			candidates = append(candidates, s)
		}
	}

	if len(candidates) == 1 {
		return candidates[0].Team
	}

	leader = candidates[0].Team
	maxGoalDiff := candidates[0].Goals - candidates[0].Against

	for _, s := range candidates[1:] {
		currentDiff := s.Goals - s.Against
		if currentDiff > maxGoalDiff {
			leader = s.Team
			maxGoalDiff = currentDiff
		}
	}
	return leader
}
