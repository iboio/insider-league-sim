package models

type MatchOutcome struct {
	Winner      Team `json:"winner"`
	Loser       Team `json:"loser"`
	IsDraw      bool `json:"isDraw"`
	WinnerGoals int  `json:"winnerGoals"`
	LoserGoals  int  `json:"loserGoals"`
}

type MatchResult struct {
	MatchWeek int    `json:"matchWeek"`
	Home      string `json:"home"`
	HomeScore int    `json:"homeScore"`
	Away      string `json:"away"`
	AwayScore int    `json:"awayScore"`
	Winner    string `json:"winner"`
}
