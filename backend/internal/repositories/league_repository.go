package repositories

import (
	"database/sql"

	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"
)

type leagueRepository struct {
	db *sql.DB
}

func NewLeagueRepository(db *sql.DB) interfaces.LeagueRepository {
	return &leagueRepository{
		db: db,
	}
}

func (lr *leagueRepository) SetLeague(id string, data models.CreateLeagueRequest) error {
	query := `INSERT INTO league (leagueId, name) VALUES (?,?)`

	_, err := lr.db.Exec(query, id, data.LeagueName)

	if err != nil {

		return err
	}

	return nil
}

func (lr *leagueRepository) GetLeague() ([]models.GetLeaguesIdsWithNameResponse, error) {
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
