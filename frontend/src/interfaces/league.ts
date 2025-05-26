export interface Team {
    name: string;
    attackPower: number;
    defensePower: number;
    morale: number;
    stamina: number;
}

export interface Standings {
    team: Team;
    goals: number;
    against: number;
    played: number;
    wins: number;
    losses: number;
    points: number;
}

export interface Match {
    home: Team;
    away: Team;
}

export interface Week {
    number: number;
    matches: Match[];
}

export interface PredictedStanding {
    team_name: string;
    points: number;
    strength: number;
    odds: number;
    eliminated: boolean;
}