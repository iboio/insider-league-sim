package repositories

import (
	"database/sql"
	"fmt"
	"league-sim/internal/layers/infra"
	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"
	"strings"
)

type teamRepository struct {
	db *sql.DB
}

func NewTeamsRepository(infra *infra.Infra) interfaces.TeamsRepository {
	return &teamRepository{db: infra.MysqlConn}
}

func (tr *teamRepository) GetTeamByName(leagueId, teamName string) (models.Team, error) {
	query := `
       SELECT leagueId, teamName, attackPower, defensePower, morale, stamina 
       FROM teams 
       WHERE leagueId = ? AND teamName = ?
   `

	var team models.Team
	err := tr.db.QueryRow(query, leagueId, teamName).Scan(
		&team.LeagueId,
		&team.TeamName,
		&team.AttackPower,
		&team.DefensePower,
		&team.Morale,
		&team.Stamina,
	)

	if err != nil {
		return models.Team{}, err
	}

	return team, nil
}

func (tr *teamRepository) GetTeams(leagueId string) ([]models.Team, error) {
	query := `
       SELECT leagueId, teamName, attackPower, defensePower, morale, stamina 
       FROM teams 
       WHERE leagueId = ?`

	rows, err := tr.db.Query(query, leagueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(
			&team.LeagueId,
			&team.TeamName,
			&team.AttackPower,
			&team.DefensePower,
			&team.Morale,
			&team.Stamina,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func (tr *teamRepository) SetTeams(leagueId string, teams []models.Team) error {
	if len(teams) == 0 {
		return nil
	}

	placeholder := "(?, ?, ?, ?, ?, ?)"
	placeholders := strings.Repeat(placeholder+", ", len(teams))
	placeholders = placeholders[:len(placeholders)-2]

	query := fmt.Sprintf(
		`
       INSERT INTO teams 
       (leagueId, teamName, attackPower, defensePower, morale, stamina) 
       VALUES %s
   `, placeholders)

	var vals []interface{}
	for _, team := range teams {
		vals = append(
			vals,
			leagueId,
			team.TeamName,
			team.AttackPower,
			team.DefensePower,
			team.Morale,
			team.Stamina,
		)
	}

	_, err := tr.db.Exec(query, vals...)
	return err
}

func (tr *teamRepository) UpdateTeam(team models.Team) error {

	attackPower := int(team.AttackPower)
	defensePower := int(team.DefensePower)
	morale := int(team.Morale)
	stamina := int(team.Stamina)

	query := `UPDATE teams SET attackPower = ?, defensePower = ?, morale = ?, stamina = ? WHERE leagueId = ? AND teamName = ?`

	_, err := tr.db.Exec(query, attackPower, defensePower, morale, stamina, team.LeagueId, team.TeamName)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	return nil
}
