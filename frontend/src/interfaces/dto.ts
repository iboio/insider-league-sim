import type {Standings, Week} from "@/interfaces/league.ts";

export interface CreateLeagueRequest {
    leagueName: string;
    teamCount: string;
}

export interface GetLeaguesIdsWithNameResponse {
    leagueId: string;
    leagueName: string;
}

export interface GetActiveLeagueStandingsResponse {
    standings: Standings[];
}

export interface GetActiveLeagueFixturesResponse {
    upcomingFixtures: Week[];
    playedFixtures: Week[];
}

export interface UpdateStandingTableRequest {
    leagueId: string;
    teamName: string;
    points: number;
}

export interface SimulateLeagueRequest {
    playAllFixture: boolean;
}

export interface SimulationResponse {
    upcomingFixtures: Week[];
    playedFixtures: Week[];
}

export interface LeagueIdWithName {
    leagueId: string;
    leagueName: string;
}