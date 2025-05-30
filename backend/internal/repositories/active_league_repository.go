package repositories

import (
	"database/sql"
	"league-sim/internal/layers/infra"

	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"
	"league-sim/utils"
)

type activeLeagueRepository struct {
	db *sql.DB
}

func NewActiveLeagueRepository(db *infra.Infra) interfaces.ActiveLeagueRepository {
	return &activeLeagueRepository{
		db: db.MysqlConn,
	}
}

func (alr *activeLeagueRepository) GetActiveLeague(id string) (models.League, error) {
	query := `SELECT upcomingFixtures, playedFixtures, currentWeek,teams, standings FROM active_league WHERE leagueId = ? order by createdAt desc limit 1`
	row := alr.db.QueryRow(query, id)

	var upcomingFixturesJson, playedFixturesJson, standingsJson, teamsJson string
	var CurrentWeek int
	err := row.Scan(&upcomingFixturesJson, &playedFixturesJson, &CurrentWeek, &teamsJson, &standingsJson)
	if err != nil {
		return models.League{}, err
	}

	var league models.League
	league.CurrentWeek = CurrentWeek

	upcomingFixtures, err := utils.StringToStruct[[]models.Week](upcomingFixturesJson)

	if err != nil {

		return models.League{}, err
	}

	playedFixtures, err := utils.StringToStruct[[]models.Week](playedFixturesJson)

	if err != nil {

		return models.League{}, err
	}

	standings, err := utils.StringToStruct[[]models.Standings](standingsJson)

	if err != nil {

		return models.League{}, err
	}

	teams, err := utils.StringToStruct[[]models.Team](teamsJson)

	if err != nil {

		return models.League{}, err
	}

	league.UpcomingFixtures = upcomingFixtures
	league.PlayedFixtures = playedFixtures
	league.Standings = standings
	league.LeagueID = id
	league.Teams = teams

	return league, nil
}

func (alr *activeLeagueRepository) GetActiveLeagueTeams(id string) ([]models.Team, error) {
	query := `SELECT teams FROM active_league WHERE leagueId = ? order by createdAt desc limit 1`
	row := alr.db.QueryRow(query, id)
	var teamsJson string

	err := row.Scan(&teamsJson)
	if err != nil {
		panic(err)

		return nil, err
	}

	teams, err := utils.StringToStruct[[]models.Team](teamsJson)

	if err != nil {
		panic(err)

		return nil, err
	}

	return teams, nil
}

func (alr *activeLeagueRepository) SetActiveLeague(data models.League) error {
	upcomingFixtures := utils.StructToString[[]models.Week](data.UpcomingFixtures)
	playedFixtures := utils.StructToString[[]models.Week](data.PlayedFixtures)
	standings := utils.StructToString[[]models.Standings](data.Standings)
	teams := utils.StructToString[[]models.Team](data.Teams)

	query := `INSERT INTO active_league (leagueId, upcomingFixtures ,playedFixtures,teams, currentWeek, standings) VALUES (?, ?, ?, ?, ?,?)`

	_, err := alr.db.Exec(
		query,
		data.LeagueID,
		upcomingFixtures,
		playedFixtures,
		teams,
		data.CurrentWeek,
		standings)

	if err != nil {
		panic(err)

		return err
	}

	return nil
}

func (alr *activeLeagueRepository) GetActiveLeaguesFixtures(id string) (models.GetActiveLeagueFixturesResponse, error) {
	query := `SELECT upcomingFixtures, playedFixtures FROM active_league WHERE leagueId = ? order by createdAt desc limit 1`

	row := alr.db.QueryRow(query, id)
	var upcomingFixturesJson, playedFixturesJson string

	err := row.Scan(&upcomingFixturesJson, &playedFixturesJson)

	if err != nil {

		return models.GetActiveLeagueFixturesResponse{}, err
	}

	upcomingFixtures, err := utils.StringToStruct[[]models.Week](upcomingFixturesJson)

	if err != nil {

		return models.GetActiveLeagueFixturesResponse{}, err
	}

	playedFixtures, err := utils.StringToStruct[[]models.Week](playedFixturesJson)

	if err != nil {

		return models.GetActiveLeagueFixturesResponse{}, err
	}

	return models.GetActiveLeagueFixturesResponse{
		UpcomingFixtures: upcomingFixtures,
		PlayedFixtures:   playedFixtures,
	}, nil
}

func (alr *activeLeagueRepository) GetActiveLeaguesStandings(id string) (
	[]models.Standings, error,
) {
	query := `SELECT standings FROM active_league WHERE leagueId = ? order by createdAt desc limit 1`

	row := alr.db.QueryRow(query, id)
	var standingsJson string

	err := row.Scan(&standingsJson)

	if err != nil {

		return nil, err
	}

	standings, err := utils.StringToStruct[[]models.Standings](standingsJson)

	if err != nil {

		return nil, err
	}

	return standings, nil
}
