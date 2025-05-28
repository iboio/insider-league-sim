package repositories

import (
	"database/sql"
	"errors"
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
			&mr.MatchWeek,
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
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?)")
		args = append(
			args,
			leagueId,
			mr.Home,
			mr.HomeScore,
			mr.Away,
			mr.AwayScore,
			mr.Winner,
			mr.MatchWeek,
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
	query := `UPDATE match_results SET homeGoals = ?, awayGoals = ?, winnerName = ? WHERE leagueId = ? AND matchWeek = ? AND homeTeam = ? AND awayTeam = ?`
	_, err := mrr.db.Exec(
		query,
		data.HomeScore,
		data.AwayScore,
		data.Winner,
		data.LeagueId,
		data.MatchWeek,
		data.Home,
		data.Away)
	if err != nil {
		fmt.Println("Error updating match score:", err)
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

func (mmr *matchResultRepository) GetMatchResultByWeekAndTeam(data models.EditMatchResult) (
	models.MatchResult, error) {

	query := `
		SELECT homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek
		FROM match_results
		WHERE leagueId = ? AND matchWeek = ? AND homeTeam = ? AND awayTeam = ?
	`

	row := mmr.db.QueryRow(query, data.LeagueId, data.MatchWeek, data.Home, data.Away)

	var queryData models.MatchResult
	err := row.Scan(
		&queryData.Home,
		&queryData.HomeScore,
		&queryData.Away,
		&queryData.AwayScore,
		&queryData.Winner,
		&queryData.MatchWeek)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No match found for the given week and teams")
			return models.MatchResult{}, err
		}
		fmt.Println("Error scanning match result:", err)
		return models.MatchResult{}, err
	}

	return queryData, nil
}
