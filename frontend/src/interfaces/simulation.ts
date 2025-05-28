export interface MatchResult {
    matchWeek: number;
    home: string;
    homeScore: number;
    away: string;
    awayScore: number;
    winner: string;
}

export interface EditMatchData {
    home: string;
    homeScore: number;
    away: string;
    awayScore: number;
    matchWeek: number;
}