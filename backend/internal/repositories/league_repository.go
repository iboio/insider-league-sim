package repositories

import (
	"database/sql"
	"errors"
	"league-sim/internal/layers/infra"

	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"
)

type leagueRepository struct {
	db *sql.DB
}

func NewLeagueRepository(db *infra.Infra) interfaces.LeagueRepository {
	return &leagueRepository{
		db: db.MysqlConn,
	}
}

func (lr *leagueRepository) SetLeague(data models.CreateLeagueRequest) error {
	query := `INSERT INTO league (leagueId, name, teamCount) VALUES (?,?,?)`

	_, err := lr.db.Exec(query, data.LeagueId, data.LeagueName, data.TeamCount)

	if err != nil {

		return err
	}

	return nil
}

func (lr *leagueRepository) GetLeagueById(id string) (models.CreateLeagueRequest, error) {
	query := `SELECT leagueId, name, teamCount FROM league WHERE leagueId = ?`

	row := lr.db.QueryRow(query, id)

	var league models.CreateLeagueRequest
	err := row.Scan(&league.LeagueId, &league.LeagueName, &league.TeamCount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.CreateLeagueRequest{}, nil
		}
		return models.CreateLeagueRequest{}, err
	}

	return league, nil

}

func (lr *leagueRepository) GetLeagues() ([]models.GetLeaguesIdsWithNameResponse, error) {
	query := `SELECT leagueId, name FROM league`

	rows, err := lr.db.Query(query)

	if err != nil {
		panic(err)
		return nil, err
	}

	var leagues []models.GetLeaguesIdsWithNameResponse

	for rows.Next() {
		var league models.GetLeaguesIdsWithNameResponse
		err := rows.Scan(&league.LeagueId, &league.LeagueName)
		if err != nil {
			panic(err)
			return nil, err
		}
		leagues = append(leagues, league)
	}

	return leagues, nil
}

func (lr *leagueRepository) DeleteLeague(id string) error {
	query := `DELETE FROM league WHERE leagueId = ?`

	_, err := lr.db.Exec(query, id)

	if err != nil {
		panic(err)

		return err
	}

	return nil
}
