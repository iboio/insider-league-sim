package models

type League struct {
	LeagueID         string              `json:"leagueId"`
	LeagueName       string              `json:"leagueName"`
	Teams            []Team              `json:"teams"`
	Standings        []Standings         `json:"standings"`
	TotalWeeks       int                 `json:"totalWeeks"`
	CurrentWeek      int                 `json:"currentWeek"`
	UpcomingFixtures []Week              `json:"upcomingFixtures"`
	PlayedFixtures   []Week              `json:"playedFixtures"`
	Predict          []PredictedStanding `json:"predict"`
}
