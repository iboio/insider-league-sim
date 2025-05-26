import type {Team} from "@/interfaces/league.ts";


export interface MatchOutcome {
    winner: Team;
    loser: Team;
    isDraw: boolean;
    winnerGoals: number;
    loserGoals: number;
}

export interface MatchResult {
    weekNumber: number;
    home: string;
    homeScore: number;
    away: string;
    awayScore: number;
    winner: string;
}
