export interface Team {
    name: string;
    attackPower: number;
    defensePower: number;
    morale: number;
    stamina: number;
}

export interface Standings {
    teamName: string;
    goals: number;
    against: number;
    played: number;
    wins: number;
    losses: number;
    points: number;
    draws: number;
    diff: number;
}

export interface Match {
    home: string;
    away: string;
}

export interface Week {
    number: number;
    matches: Match[];
}

export interface PredictedStanding {
    teamName: string;
    points: number;
    strength: number;
    odds: number;
    eliminated: boolean;
}