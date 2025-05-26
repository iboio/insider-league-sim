import {useEffect, useState} from 'react'
import {useParams, useNavigate} from 'react-router-dom'
import {getLeagueById, simulateMatches, resetLeague, updateMatchResult} from '../services/api'
import {Button} from '../components/ui/button'
import type {LeagueData} from '../interfaces/full'

// Import our components
import StandingsTable from '../components/league/StandingsTable'
import Fixtures from '../components/league/Fixtures'
import MatchResults from '../components/league/MatchResults'
import Predictions from '../components/league/Predictions'
import type {EditMatchData} from "@/interfaces/api.ts";

// Main League Component
function League() {
    const {leagueId} = useParams<{ leagueId: string }>();
    const navigate = useNavigate();
    const [allData, setAllData] = useState<LeagueData>({
        leagueId: '',
        leagueName: '',
        teams: [],
        standings: [],
        totalWeeks: 0,
        currentWeek: 0,
        upcomingFixtures: [],
        playedFixtures: [],
        predict: [],
        matches: [],
    });
    const [loading, setLoading] = useState(true);
    const [simulating, setSimulating] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Fetch league data function
    const fetchLeagueData = async () => {
        if (!leagueId) return;

        try {
            setLoading(true);
            const data = await getLeagueById(leagueId);
            setAllData(prevState => ({...prevState, ...data}));
            setError(null);
        } catch (err) {
            console.error('Error fetching league data:', err);
            setError('Failed to load league data. Please try again.');
        } finally {
            setLoading(false);
        }
    };

    const handleMatchUpdate = async (editData: EditMatchData) => {
        if (!leagueId) return;
        try {
            await updateMatchResult(leagueId, editData);

            const updatedLeagueData = await getLeagueById(leagueId);
            setAllData(prevState => ({...prevState, ...updatedLeagueData}));
            await fetchLeagueData();
        } catch (error) {
            console.error('Match update failed:', error);
        }
    };

    useEffect(() => {
        fetchLeagueData();
    }, [leagueId]);

    // Reset league function
    const handleResetLeague = async () => {
        if (!leagueId) return;

        try {
            setSimulating(true);
            await resetLeague(leagueId);

            // Get updated league data after reset
            await fetchLeagueData();

            setError(null);
        } catch (err) {
            console.error('Error resetting league:', err);
            setError('Failed to reset league. Please try again.');
        } finally {
            setSimulating(false);
        }
    };

    // Simulate matches function
    const handleSimulateMatches = async (playAll = false) => {
        if (!leagueId) return;

        try {
            setSimulating(true);
            const data = await simulateMatches(leagueId, playAll);

            // Get updated league data after simulation without causing a refresh
            const updatedLeagueData = await getLeagueById(leagueId);

            // Update all data with the combined results
            setAllData(prevState => ({...prevState, ...updatedLeagueData, ...data}));

            setError(null);
        } catch (err) {
            console.error('Error simulating matches:', err);
            setError('Failed to simulate matches. Please try again.');
        } finally {
            setSimulating(false);
        }
    };

    // Calculate current week and total week
    const currentWeek = allData.playedFixtures.length;
    const totalWeek = allData.upcomingFixtures.length + allData.playedFixtures.length;
    const allWeeksPlayed = currentWeek === totalWeek || allData.upcomingFixtures.length === 0;

    if (loading) {
        return (
            <div className="container mx-auto px-4 py-8 min-h-screen">
                <div
                    className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="text-red-500 text-center">
                    <p className="text-xl font-semibold">{error}</p>
                    <Button
                        onClick={() => window.location.reload()}
                        className="mt-4"
                    >
                        Try Again
                    </Button>
                </div>
            </div>
        );
    }

    if (!allData.leagueId) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="text-center">
                    <p className="text-xl font-semibold">League not found</p>
                    <Button
                        onClick={() => window.history.back()}
                        className="mt-4"
                    >
                        Go Back
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div className="w-full max-w-[98%] lg:max-w-[95%] mx-auto py-2" style={{minHeight: '100vh'}}>
            {/* No overlay, keep content stable */}
            {/* Header with league name and buttons */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-3 gap-2">
                <div>
                    <h1 className="text-xl sm:text-2xl font-bold text-gray-900">League {leagueId}</h1>
                    <div className="text-sm text-gray-600 mt-1">
                        <span className="font-medium">Week:</span> {currentWeek} / {totalWeek}
                    </div>
                </div>
                <div className="flex items-center space-x-2 w-full sm:w-auto">
                    <Button
                        onClick={() => handleSimulateMatches(false)}
                        disabled={simulating || allWeeksPlayed}
                        className="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 px-4 py-2 text-sm text-black flex-1 sm:flex-none min-w-[120px] flex items-center justify-center"
                    >
                        {simulating ? (
                            <>
                                <span
                                    className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-current border-r-transparent mr-2"></span>
                                Simulating
                            </>
                        ) : 'Play One Week'}
                    </Button>
                    <Button
                        onClick={() => handleSimulateMatches(true)}
                        disabled={simulating || allWeeksPlayed}
                        className="bg-green-600 hover:bg-green-700 disabled:opacity-50 px-4 py-2 text-sm text-black flex-1 sm:flex-none min-w-[120px] flex items-center justify-center"
                    >
                        {simulating ? (
                            <>
                                <span
                                    className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-current border-r-transparent mr-2"></span>
                                Simulating
                            </>
                        ) : 'Play All Weeks'}
                    </Button>
                    <Button
                        onClick={handleResetLeague}
                        disabled={simulating}
                        className="bg-red-600 hover:bg-red-700 disabled:opacity-50 px-4 py-2 text-sm text-black flex-1 sm:flex-none min-w-[120px] flex items-center justify-center"
                    >
                        {simulating ? (
                            <>
                                <span
                                    className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-current border-r-transparent mr-2"></span>
                                Resetting
                            </>
                        ) : 'Reset League'}
                    </Button>
                    <Button
                        onClick={() => navigate('/')}
                        className="bg-gray-600 hover:bg-gray-700 px-4 py-2 text-sm text-black flex-1 sm:flex-none min-w-[120px] flex items-center justify-center"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-2" viewBox="0 0 20 20"
                             fill="currentColor">
                            <path
                                d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
                        </svg>
                        Home
                    </Button>
                </div>
            </div>

            {/* Show champion banner if no upcoming fixtures */}
            {allData.upcomingFixtures.length === 0 && allData.standings.length > 0 && (
                <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-6 rounded-lg shadow-lg mb-8">
                    <h2 className="text-2xl font-bold mb-2">🏆 Champion Declared!</h2>
                    <p className="text-xl">
                        {(() => {
                            // Sort standings by points, then by goal difference
                            const sortedStandings = [...allData.standings].sort((a, b) => {
                                if (a.points !== b.points) {
                                    return b.points - a.points; // Higher points first
                                }

                                const aGoalDiff = a.goals - a.against;
                                const bGoalDiff = b.goals - b.against;
                                return bGoalDiff - aGoalDiff;
                            });

                            // Get the champion (first team after sorting)
                            const champion = sortedStandings[0];
                            return champion ? `${champion.team.name} wins the league with ${champion.points} points and a goal difference of ${champion.goals - champion.against}!` : 'No champion determined';
                        })()}
                        <p className="text-white/60">If two teams have the same number of points, the one with the lower goal difference will be the champion.
                        </p>
                    </p>
                </div>
            )}

            {/* Main content layout - Split into separate sections */}
            {/* Top row with standings, predictions, and match results */}
            <div className="mb-8">
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-3">
                    <div className="w-full">
                        <div className="bg-white rounded-lg shadow-md p-2 min-h-[300px]">
                            <h2 className="text-lg font-semibold mb-2">Standings</h2>
                            <StandingsTable standings={allData.standings}/>
                        </div>
                    </div>
                    <div className="w-full">
                        <div className="bg-white rounded-lg shadow-md p-2 min-h-[300px]">
                            <h2 className="text-lg font-semibold mb-2">Predictions</h2>
                            <Predictions predictions={allData.predict}/>
                        </div>
                    </div>
                    <div className="w-full">
                        <div className="bg-white rounded-lg shadow-md p-2 min-h-[300px]">
                            <h2 className="text-lg font-semibold mb-2">Match Results</h2>
                            <MatchResults
                                matches={allData.matches}
                                leagueId={leagueId}
                                onMatchUpdate={handleMatchUpdate}
                            />
                        </div>
                    </div>
                </div>
            </div>

            {/* Bottom row with fixtures */}
            <div className="mb-8">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-3">
                    <div className="w-full">
                        <div className="bg-white rounded-lg shadow-md p-2 min-h-[300px]">
                            <h2 className="text-lg font-semibold mb-2">Played Fixtures</h2>
                            <Fixtures weeks={allData.playedFixtures} title=""/>
                        </div>
                    </div>
                    <div className="w-full">
                        <div className="bg-white rounded-lg shadow-md p-2 min-h-[300px]">
                            <h2 className="text-lg font-semibold mb-2">Upcoming Fixtures</h2>
                            <Fixtures weeks={allData.upcomingFixtures} title=""/>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default League