package appContext

import (
	"database/sql"

	"league-sim/internal/builder"
	"league-sim/internal/repositories"
	"league-sim/internal/repositories/interfaces"
)

type AppContext interface {
	LeagueRepository() interfaces.LeagueRepository
	ActiveLeagueRepository() interfaces.ActiveLeagueRepository
	MatchResultRepository() interfaces.MatchResultRepository
	DB() *DB
}
type DB struct {
	Sql *sql.DB
}
type AppContextImpl struct {
	db                     *DB
	leagueRepository       interfaces.LeagueRepository
	activeLeagueRepository interfaces.ActiveLeagueRepository
	matchResultRepository  interfaces.MatchResultRepository
}

func (a *AppContextImpl) DB() *DB {

	return a.db
}

func (a *AppContextImpl) LeagueRepository() interfaces.LeagueRepository {

	return a.leagueRepository
}

func (a *AppContextImpl) ActiveLeagueRepository() interfaces.ActiveLeagueRepository {

	return a.activeLeagueRepository
}
func (a *AppContextImpl) MatchResultRepository() interfaces.MatchResultRepository {

	return a.matchResultRepository
}

func AppContextInit() (*AppContextImpl, error) {
	db, err := AppContextDBInit()
	if err != nil {
		return nil, err
	}

	leagueRepository := repositories.NewLeagueRepository(db.Sql)
	activeLeagueRepository := repositories.NewActiveLeagueRepository(db.Sql)
	matchResultRepository := repositories.NewMatchResultRepository(db.Sql)

	return &AppContextImpl{
		db:                     db,
		activeLeagueRepository: activeLeagueRepository,
		leagueRepository:       leagueRepository,
		matchResultRepository:  matchResultRepository,
	}, nil
}

func AppContextDBInit() (*DB, error) {
	sql, err := builder.SqlConnectionInit()
	if err != nil {

		return nil, err
	}

	return &DB{
		Sql: sql,
	}, nil
}
