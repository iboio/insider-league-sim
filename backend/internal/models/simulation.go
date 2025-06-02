package models

type MatchOutcome struct {
	Winner      Team `json:"winner"`
	Loser       Team `json:"loser"`
	IsDraw      bool `json:"isDraw"`
	WinnerGoals int  `json:"winnerGoals"`
	LoserGoals  int  `json:"loserGoals"`
}

type Matches struct {
	LeagueId  string `json:"leagueId"`
	MatchWeek int    `json:"matchWeek"`
	Home      string `json:"home"`
	HomeScore int    `json:"homeScore"`
	Away      string `json:"away"`
	AwayScore int    `json:"awayScore"`
	IsPlayed  bool   `json:"isPlayed"`
	Winner    string `json:"winner"`
}
