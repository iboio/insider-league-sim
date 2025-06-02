package models

type CreateLeagueRequest struct {
	LeagueId   string `json:"leagueId"`
	LeagueName string `json:"leagueName"`
	TeamCount  string `json:"teamCount"`
}
type GetLeaguesIdsWithNameResponse struct {
	LeagueId   string `json:"leagueId"`
	LeagueName string `json:"leagueName"`
}

type GetActiveLeagueStandingsResponse struct {
	Standings []Standings `json:"standings"`
}
type GetFixturesResponse struct {
	UpcomingFixtures []Week `json:"upcomingFixtures"`
	PlayedFixtures   []Week `json:"playedFixtures"`
}

type EditMatchResult struct {
	LeagueId  string `json:"leagueId"`
	Home      string `json:"home"`
	Away      string `json:"away"`
	HomeScore int    `json:"homeScore"`
	AwayScore int    `json:"awayScore"`
	MatchWeek int    `json:"matchWeek"`
	Winner    string `json:"winner"`
}

type SimulateLeagueRequest struct {
	PlayAllFixture bool `json:"playAllFixture"`
}

type SimulationResponse struct {
	Matches          []Matches `json:"matches"`
	UpcomingFixtures []Week    `json:"upcomingFixtures"`
	PlayedFixtures   []Week    `json:"playedFixtures"`
}
