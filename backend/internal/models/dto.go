package models

type CreateLeagueRequest struct {
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
type GetActiveLeagueFixturesResponse struct {
	UpcomingFixtures []Week `json:"upcomingFixtures"`
	PlayedFixtures   []Week `json:"playedFixtures"`
}

type EditMatchResult struct {
	WeekNumber   int    `json:"weekNumber"`
	LeagueId     string `json:"leagueId"`
	TeamName     string `json:"teamName"`
	TeamType     string `json:"teamType"`
	AgainstTeam  string `json:"againstTeam"`
	TeamOldGoals int    `json:"teamOldPoints"`
	Goals        int    `json:"goals"`
	IsDraw       bool   `json:"isDraw"`
}

type SimulateLeagueRequest struct {
	PlayAllFixture bool `json:"playAllFixture"`
}

type SimulationResponse struct {
	Matches          []MatchResult `json:"matches"`
	UpcomingFixtures []Week        `json:"upcomingFixtures"`
	PlayedFixtures   []Week        `json:"playedFixtures"`
}
