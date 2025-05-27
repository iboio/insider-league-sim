import axios from 'axios';
import type {
    LeagueIdWithName,
    SimulationResponse,
    GetActiveLeagueStandingsResponse,
    GetActiveLeagueFixturesResponse,
} from '../interfaces/dto';
import type {PredictedStanding} from '../interfaces/league';
import type {LeagueData} from '../interfaces/full';
import type {EditMatchData} from "@/interfaces/api.ts";
import type {MatchResult} from "@/interfaces/simulation.ts";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    headers: {
        'Content-Type': 'application/json',
    },

});

// API functions
export async function getLeagues(): Promise<LeagueIdWithName[]> {
    try {
        const response = await api.get('/league');
        return response.data;
    } catch (error) {
        console.error('Error fetching leagues:', error);
        throw error;
    }
}

export async function createLeague(leagueName: string, teamCount: string) {
    try {
        const response = await api.post('/league', {
            leagueName,
            teamCount
        });
        return response.data;
    } catch (error) {
        console.error('Error creating league:', error);
        throw error;
    }
}

export async function deleteLeague(leagueId: string) {
    try {
        const response = await api.delete(`/league/${leagueId}`);
        return response.data;
    } catch (error) {
        console.error('Error deleting league:', error);
        throw error;
    }
}

export async function updateMatchResult(leagueId: string, data: EditMatchData) {
    try {
        const response = await api.put(`/league/${leagueId}`, data);
        return response.data;
    } catch (error) {
        console.error('Error updating match result:', error);
        throw error;
    }
}

export async function resetLeague(leagueId: string) {
    try {
        const response = await api.post(`/league/${leagueId}/reset`);
        return response.data;
    } catch (error) {
        console.error('Error resetting league:', error);
        throw error;
    }
}

export async function getLeagueById(leagueId: string): Promise<LeagueData> {
    try {
        // Get league details

        // Get standings
        const standingsResponse = await api.get<GetActiveLeagueStandingsResponse>(`/league/${leagueId}/standing`);

        // Get fixtures
        const fixturesResponse = await api.get<GetActiveLeagueFixturesResponse>(`/league/${leagueId}/fixtures`);

        // Get predictions
        const predictionsResponse = await api.get<PredictedStanding[]>(`/league/${leagueId}/predict`);

        // Get matches
        const matchResponse = await api.get<MatchResult[]>(`/league/${leagueId}/matchResults`);

        return {
            currentWeek: 0,
            leagueName: "",
            teams: [],
            totalWeeks: 0,
            leagueId: leagueId,
            standings: standingsResponse.data?.standings || standingsResponse.data || [],
            upcomingFixtures: fixturesResponse.data?.upcomingFixtures || [],
            playedFixtures: fixturesResponse.data?.playedFixtures || [],
            predict: Array.isArray(predictionsResponse.data) ? predictionsResponse.data : [],
            matches: matchResponse.data
        };
    } catch (error) {
        console.error('Error fetching league:', error);
        throw error;
    }
}

export async function simulateMatches(leagueId: string, playAllFixture: boolean = false): Promise<SimulationResponse> {
    try {
        // Run the simulation
        const simulationResponse = await api.post(`/league/${leagueId}/simulation`, {
            playAllFixture
        });

        return {
            upcomingFixtures: simulationResponse.data.upcomingFixtures || [],
            playedFixtures: simulationResponse.data.playedFixtures || [],
        };
    } catch (error) {
        console.error('Error simulating matches:', error);
        throw error;
    }
}
