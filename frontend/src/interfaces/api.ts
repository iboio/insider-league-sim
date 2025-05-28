// API Response and Request interfaces


export interface CreateLeagueRequest {
    leagueName: string;
    teamCount: string;
}

export interface UpdateStandingTableRequest {
    leagueId: string;
    teamName: string;
    points: number;
}

export interface SimulateLeagueRequest {
    playAllFixture: boolean;
}

export interface GetActiveLeagueFixturesResponse {
    upcomingFixtures: import('./league').Week[];
    playedFixtures: import('./league').Week[];
}

export interface SimulationResponse {
    matches: import('./simulation').MatchResult[];
    upcomingFixtures: import('./league').Week[];
    playedFixtures: import('./league').Week[];
}

export interface EditMatchData {
    weekNumber: number;
    teamName: string;
    teamType: 'home' | 'away';
    againstTeam: string;
    teamOldGoals: number;
    goals: number;
    isDraw: boolean;
    winner: string;
}

export interface MatchResult {
    weekNumber: number;
    home: string;
    homeScore: number;
    away: string;
    awayScore: number;
    winner: string;
}