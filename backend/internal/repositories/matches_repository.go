package repositories

import (
	"database/sql"
	"fmt"
	"league-sim/internal/layers/infra"
	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"
	"strings"
)

type matchRepository struct {
	db *sql.DB
}

func NewMatchResultRepository(db *infra.Infra) interfaces.MatchesRepository {
	return &matchRepository{db: db.MysqlConn}
}

func (mr *matchRepository) SetMatches(matches []models.Matches) error {
	if len(matches) == 0 {
		return nil
	}

	const baseQuery = `
		INSERT INTO matches 
		(leagueId, home, homeGoals, away, awayGoals, winner, matchWeek, isPlayed) 
		VALUES %s
	`

	var vals []interface{}
	placeholders := make([]string, 0, len(matches))

	for _, match := range matches {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?)")
		vals = append(
			vals,
			match.LeagueId,
			match.Home,
			match.HomeScore,
			match.Away,
			match.AwayScore,
			match.Winner,
			match.MatchWeek,
			match.IsPlayed,
		)
	}

	query := fmt.Sprintf(baseQuery, strings.Join(placeholders, ", "))

	_, err := mr.db.Exec(query, vals...)
	if err != nil {
		return err
	}

	return nil
}

func (mr *matchRepository) GetMatches(leagueId string) ([]models.Matches, error) {
	query := `SELECT home, homeGoals, away, awayGoals, isPlayed , winner, matchWeek FROM matches WHERE leagueId = ? ORDER BY matchWeek`

	rows, err := mr.db.Query(query, leagueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Matches

	for rows.Next() {
		var mr models.Matches
		err := rows.Scan(
			&mr.Home,
			&mr.HomeScore,
			&mr.Away,
			&mr.AwayScore,
			&mr.IsPlayed,
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

func (mr *matchRepository) GetPlayedMatches(leagueId string) ([]models.Matches, error) {
	query := `SELECT home, homeGoals, away, awayGoals, winner, matchWeek FROM matches WHERE leagueId = ? AND isPlayed = true ORDER BY matchWeek`

	rows, err := mr.db.Query(query, leagueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Matches

	for rows.Next() {
		var mr models.Matches
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

func (mr *matchRepository) EditMatch(data models.Matches) error {
	query := `UPDATE matches SET homeGoals = ?, awayGoals = ?, winner = ?, isPlayed = ? WHERE leagueId = ? AND matchWeek = ? AND home = ? AND away = ?`
	_, err := mr.db.Exec(
		query,
		data.HomeScore,
		data.AwayScore,
		data.Winner,
		data.IsPlayed,
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

func (mr *matchRepository) GetMatchByTeams(data models.EditMatchResult) (models.Matches, error) {
	query := `SELECT home, homeGoals, away, awayGoals, winner, matchWeek FROM matches WHERE leagueId = ? AND home = ? AND away = ?`
	row := mr.db.QueryRow(query, data.LeagueId, data.Home, data.Away)

	var match models.Matches
	err := row.Scan(
		&match.Home,
		&match.HomeScore,
		&match.Away,
		&match.AwayScore,
		&match.Winner,
		&match.MatchWeek,
	)
	if err != nil {
		fmt.Println("Error getting match by teams:", err)
		return models.Matches{}, err
	}

	return match, nil
}

func (mr *matchRepository) GetFixtures(leagueId string) (models.GetFixturesResponse, error) {
	allMatches, err := mr.GetMatches(leagueId)
	if err != nil {
		fmt.Println("Error getting fixtures:", err)
		return models.GetFixturesResponse{}, err
	}
	var upcomingMatches []models.Week
	var playedMatches []models.Week
	for _, match := range allMatches {
		if match.IsPlayed {

			var matchWeek models.Week
			var matchedTeams models.Match

			matchedTeams.Home = match.Home
			matchedTeams.Away = match.Away
			matchWeek.Number = match.MatchWeek

			matchWeek.Matches = append(matchWeek.Matches, matchedTeams)

			playedMatches = append(playedMatches, matchWeek)

		} else {
			var matchWeek models.Week
			var matchedTeams models.Match

			matchedTeams.Home = match.Home
			matchedTeams.Away = match.Away
			matchWeek.Number = match.MatchWeek

			matchWeek.Matches = append(matchWeek.Matches, matchedTeams)

			upcomingMatches = append(upcomingMatches, matchWeek)
		}
	}
	return models.GetFixturesResponse{
		UpcomingFixtures: upcomingMatches,
		PlayedFixtures:   playedMatches,
	}, nil
}
