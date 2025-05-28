import type {PredictedStanding, Standings, Team, Week} from "@/interfaces/league.ts";
import type {MatchResult} from "@/interfaces/simulation";

export interface LeagueData {
    leagueId: string;
    leagueName: string;
    teams: Team[];
    standings: Standings[];
    totalWeeks: number;
    currentWeek: number;
    upcomingFixtures: Week[];
    playedFixtures: Week[];
    predict: PredictedStanding[];
    matches: MatchResult[];
}