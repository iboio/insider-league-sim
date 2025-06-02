package models

type Team struct {
	LeagueId     string  `json:"leagueId"`
	TeamName     string  `json:"teamName"`
	AttackPower  float64 `json:"attackPower"`
	DefensePower float64 `json:"defensePower"`
	Morale       float64 `json:"morale"`
	Stamina      float64 `json:"stamina"`
}

type Standings struct {
	LeagueId string `json:"leagueId"`
	TeamName string `json:"teamName"`
	Goals    int    `json:"goals"`
	Against  int    `json:"against"`
	Diff     int    `json:"diff"`
	Played   int    `json:"played"`
	Wins     int    `json:"wins"`
	Draws    int    `json:"draws"`
	Losses   int    `json:"losses"`
	Points   int    `json:"points"`
}

type Match struct {
	Home string `json:"home"`
	Away string `json:"away"`
}

type Week struct {
	Number  int     `json:"number"`
	Matches []Match `json:"matches"`
}

type PredictedStanding struct {
	TeamName   string  `json:"teamName"`
	Points     int     `json:"points"`
	Strength   float64 `json:"strength"`
	Odds       float64 `json:"odds"`
	Eliminated bool    `json:"eliminated"`
}
