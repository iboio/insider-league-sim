package repositories

import (
	"database/sql"
	"fmt"
	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"
	"strings"
)

type matchResultRepository struct {
	db *sql.DB
}

func NewMatchResultRepository(db *sql.DB) interfaces.MatchResultRepository {
	return &matchResultRepository{db: db}
}

func (mrr *matchResultRepository) GetMatchResults(leagueId string) ([]models.MatchResult, error) {
	query := `SELECT homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek FROM match_results WHERE leagueId = ? ORDER BY matchWeek`

	rows, err := mrr.db.Query(query, leagueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.MatchResult

	for rows.Next() {
		var mr models.MatchResult
		err := rows.Scan(
			&mr.Home,
			&mr.HomeScore,
			&mr.Away,
			&mr.AwayScore,
			&mr.Winner,
			&mr.WeekNumber,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, mr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (mrr *matchResultRepository) SetMatchResults(leagueId string, matchResults []models.MatchResult) error {
	if len(matchResults) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(matchResults))
	args := make([]interface{}, 0, len(matchResults)*7)

	for _, mr := range matchResults {
		// Her satır için "(?, ?, ?, ?, ?, ?, ?)" şeklinde placeholder
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?)")
		args = append(
			args,
			leagueId,
			mr.Home,
			mr.HomeScore,
			mr.Away,
			mr.AwayScore,
			mr.Winner,
			mr.WeekNumber,
		)
	}

	query := fmt.Sprintf(
		"INSERT INTO match_results (leagueId, homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek) VALUES %s",
		strings.Join(placeholders, ","),
	)

	_, err := mrr.db.Exec(query, args...)
	if err != nil {

		return err
	}
	return nil
}

func (mrr *matchResultRepository) EditMatchScore(data models.EditMatchResult) error {
	changedTeam := fmt.Sprintf("%sTeam = '%s'", data.TeamType, data.TeamName)
	changedTeamGoal := fmt.Sprintf("%sGoals = %d", data.TeamType, data.Goals)

	query := `UPDATE match_results SET winnerName = ?, ` + changedTeamGoal + ` WHERE leagueId = ? AND matchWeek = ? AND ` + changedTeam

	_, err := mrr.db.Exec(query, data.TeamName, data.LeagueId, data.WeekNumber)

	if err != nil {
		panic(err)
		return err
	}

	return nil
}

func (mrr *matchResultRepository) DeleteMatchResults(leagueId string) error {
	query := `DELETE FROM match_results WHERE leagueId = ?`

	_, err := mrr.db.Exec(query, leagueId)

	if err != nil {

		return err
	}

	return nil
}
