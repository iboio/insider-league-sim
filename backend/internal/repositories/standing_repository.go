package repositories

import (
	"database/sql"
	"fmt"
	"league-sim/internal/layers/infra"
	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"
	"strings"
)

type standingRepository struct {
	db *sql.DB
}

func NewStandingRepository(db *infra.Infra) interfaces.StandingRepository {
	return &standingRepository{db: db.MysqlConn}
}

func (sr *standingRepository) GetStandings(leagueId string) ([]models.Standings, error) {
	query := `SELECT leagueId,teamName, goals, against, diff, played, wins,draws, losses, points FROM standings WHERE leagueId = ? ORDER BY points DESC`
	rows, err := sr.db.Query(query, leagueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var standings []models.Standings
	for rows.Next() {
		var standing models.Standings
		if err := rows.Scan(
			&standing.LeagueId,
			&standing.TeamName,
			&standing.Goals,
			&standing.Against,
			&standing.Diff,
			&standing.Played,
			&standing.Wins,
			&standing.Draws,
			&standing.Losses,
			&standing.Points,
		); err != nil {
			return nil, err
		}
		standings = append(standings, standing)
	}

	return standings, nil
}

func (sr *standingRepository) GetStandingsByTeam(leagueId, teamName string) (models.Standings, error) {
	query := `SELECT * FROM standings WHERE leagueId = ? AND teamName = ?`
	row := sr.db.QueryRow(query, leagueId, teamName)

	var standing models.Standings
	if err := row.Scan(
		&standing.TeamName,
		&standing.Goals,
		&standing.Against,
		&standing.Diff,
		&standing.Played,
		&standing.Wins,
		&standing.Draws,
		&standing.Losses,
		&standing.Points,
	); err != nil {
		return models.Standings{}, err
	}

	return standing, nil
}

func (sr *standingRepository) SetStandings(standings []models.Standings) error {
	if len(standings) == 0 {
		return nil
	}

	placeholder := "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	placeholders := strings.Repeat(placeholder+", ", len(standings))
	placeholders = placeholders[:len(placeholders)-2]

	query := fmt.Sprintf(
		`
        INSERT INTO standings 
        (leagueId, teamName, goals, against, played, wins, draws, losses, points) 
        VALUES %s
    `, placeholders)

	var vals []interface{}
	for _, standing := range standings {
		vals = append(
			vals,
			standing.LeagueId,
			standing.TeamName,
			standing.Goals,
			standing.Against,
			standing.Played,
			standing.Wins,
			standing.Draws,
			standing.Losses,
			standing.Points,
		)
	}

	_, err := sr.db.Exec(query, vals...)
	return err
}

func (sr *standingRepository) UpdateStanding(standing models.Standings) error {

	checkQuery := `SELECT COUNT(*) FROM standings WHERE leagueId = ? AND teamName = ?`
	var count int
	err := sr.db.QueryRow(checkQuery, standing.LeagueId, standing.TeamName).Scan(&count)
	if err != nil {
		return fmt.Errorf("check query error: %w", err)
	}

	if count == 0 {
		return fmt.Errorf(
			"no standing found with leagueId: '%s' and teamName: '%s'",
			standing.LeagueId,
			standing.TeamName)
	}

	// Update yap - rowsAffected kontrolü yapma
	updateQuery := `
       UPDATE standings 
       SET goals = ?, against = ?, played = ?, wins = ?, draws = ?, losses = ?, points = ? 
       WHERE leagueId = ? AND teamName = ?`

	_, err = sr.db.Exec(
		updateQuery,
		standing.Goals,
		standing.Against,
		standing.Played,
		standing.Wins,
		standing.Draws,
		standing.Losses,
		standing.Points,
		standing.LeagueId,
		standing.TeamName,
	)

	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	return nil
}
