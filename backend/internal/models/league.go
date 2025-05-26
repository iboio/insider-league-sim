package models

type Team struct {
	Name         string  `json:"name"`
	AttackPower  float64 `json:"attackPower"`
	DefensePower float64 `json:"defensePower"`
	Morale       float64 `json:"morale"`
	Stamina      float64 `json:"stamina"`
}

type Standings struct {
	Team    Team `json:"team"`
	Goals   int  `json:"goals"`
	Against int  `json:"against"`
	Played  int  `json:"played"`
	Wins    int  `json:"wins"`
	Losses  int  `json:"losses"`
	Points  int  `json:"points"`
}

type Match struct {
	Home *Team `json:"home"`
	Away *Team `json:"away"`
}

type Week struct {
	Number  int     `json:"number"`
	Matches []Match `json:"matches"`
}

type PredictedStanding struct {
	TeamName   string  `json:"team_name"`
	Points     int     `json:"points"`
	Strength   float64 `json:"strength"`
	Odds       float64 `json:"odds"`
	Eliminated bool    `json:"eliminated"`
}
