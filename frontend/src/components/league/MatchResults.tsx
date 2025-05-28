import { useState } from 'react';
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from '@/components/ui/card';
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from '@/components/ui/dialog';
import { Edit, Trophy } from 'lucide-react';
import type {EditMatchData, MatchResult} from "@/interfaces/api.ts";


interface MatchResultsProps {
    matches: MatchResult[] | null | undefined;
    leagueId?: string;
    onMatchUpdate?: (data: EditMatchData) => void;
}

export const MatchResults = ({ matches, leagueId, onMatchUpdate }: MatchResultsProps) => {
    const [editDialogOpen, setEditDialogOpen] = useState(false);
    const [selectedMatch, setSelectedMatch] = useState<MatchResult | null>(null);
    const [editingTeam, setEditingTeam] = useState<'home' | 'away' | null>(null);
    const [newGoals, setNewGoals] = useState('');
    const [error, setError] = useState('');
    const [updating, setUpdating] = useState(false);

    const handleEditClick = (match: MatchResult, team: 'home' | 'away') => {
        // Only losing team can edit (in case of draw, both teams can)
        const isDraw = match.homeScore === match.awayScore;
        const isLoser = (team === 'home' && match.winner === match.away) ||
            (team === 'away' && match.winner === match.home);

        if (!isDraw && !isLoser) {
            setError('Only the losing team can edit the result!');
            return;
        }

        setSelectedMatch(match);
        setEditingTeam(team);
        setNewGoals(team === 'home' ? match.homeScore.toString() : match.awayScore.toString());
        setError('');
        setEditDialogOpen(true);
    };

    const handleSaveEdit = async () => {
        if (!selectedMatch || !editingTeam) return;

        const newGoalsNum = parseInt(newGoals);
        const opponentGoals = editingTeam === 'home' ? selectedMatch.awayScore : selectedMatch.homeScore;
        const oldGoals = editingTeam === 'home' ? selectedMatch.homeScore : selectedMatch.awayScore;
        const wasOriginallyDraw = selectedMatch.homeScore === selectedMatch.awayScore;

        // Validation: New goal count must be higher than opponent (except when originally draw)
        if (!wasOriginallyDraw && newGoalsNum <= opponentGoals) {
            setError('New goal count must be higher than opponent team!');
            return;
        }

        if (wasOriginallyDraw) {
            if (newGoalsNum <= opponentGoals) {
                setError('New goal count must be higher than current draw score to win!');
                return;
            }
        }

        const willBeDraw = selectedMatch.awayScore === selectedMatch.homeScore;
        const winnerTeam = editingTeam === 'home' ? selectedMatch.home : selectedMatch.away;

        setUpdating(true);
        setError('');

        try {
            const editData: EditMatchData = {
                weekNumber: selectedMatch.weekNumber,
                teamName: editingTeam === 'home' ? selectedMatch.home : selectedMatch.away,
                teamType: editingTeam,
                againstTeam: editingTeam === 'home' ? selectedMatch.away : selectedMatch.home,
                teamOldGoals: oldGoals,
                goals: newGoalsNum,
                isDraw: willBeDraw,
                winner: winnerTeam
            };

            if (onMatchUpdate) {
                await onMatchUpdate(editData);
            }

            setEditDialogOpen(false);
            setSelectedMatch(null);
            setEditingTeam(null);
            setNewGoals('');
        } catch (error) {
            setError('An error occurred during update!');
            console.error('Error updating match:', error);
        } finally {
            setUpdating(false);
        }
    };

    if (!matches || matches.length === 0) {
        return (
            <Card className="h-full min-h-[300px] shadow-md">
                <CardHeader className="bg-blue-600 py-3">
                    <CardTitle className="text-white text-lg font-bold flex items-center gap-2">
                        <Trophy className="w-5 h-5" />
                        Match Results
                    </CardTitle>
                </CardHeader>
                <CardContent className="p-8 text-center">
                    <div className="text-gray-400 text-lg">
                        No match results available yet
                    </div>
                </CardContent>
            </Card>
        );
    }

    // Maçları haftalar halinde grupla
    const matchesByWeek = matches.reduce((acc, match) => {
        if (!acc[match.weekNumber]) acc[match.weekNumber] = [];
        acc[match.weekNumber].push(match);
        return acc;
    }, {} as Record<number, MatchResult[]>);

    const sortedWeeks = Object.entries(matchesByWeek).sort(
        ([a], [b]) => Number(a) - Number(b)
    );

    const getTeamStyle = (team: string, match: MatchResult) => {
        const isWinner = match.winner === team;
        const isDraw = match.homeScore === match.awayScore;

        if (isDraw) return 'text-amber-600 font-medium';
        if (isWinner) return 'text-green-700 font-bold';
        return 'text-gray-700';
    };

    const getScoreStyle = (match: MatchResult, isHome: boolean) => {
        const isWinner = (isHome && match.winner === match.home) || (!isHome && match.winner === match.away);
        const isDraw = match.homeScore === match.awayScore;

        if (isDraw) return 'bg-amber-100 text-amber-800 border-amber-300';
        if (isWinner) return 'bg-green-100 text-green-800 border-green-300';
        return 'bg-gray-100 text-gray-700 border-gray-300';
    };

    const canEdit = (match: MatchResult, team: 'home' | 'away') => {
        const isDraw = match.homeScore === match.awayScore;
        const isLoser = (team === 'home' && match.winner === match.away) ||
            (team === 'away' && match.winner === match.home);
        return isDraw || isLoser;
    };

    return (
        <>
            <Card className="min-h-[300px] w-full shadow-md">
                <CardHeader className="bg-gradient-to-r from-green-600 to-green-700 py-4">
                    <CardTitle className="text-white text-xl font-bold flex items-center gap-2">
                        <Trophy className="w-6 h-6" />
                        Match Results
                    </CardTitle>
                </CardHeader>
                <CardContent className="p-0">
                    <Accordion type="single" collapsible className="w-full">
                        {sortedWeeks.map(([weekNumber, weekMatches]) => (
                            <AccordionItem key={weekNumber} value={`week-${weekNumber}`} className="border-b">
                                <AccordionTrigger className="px-6 py-4 text-base font-semibold hover:bg-gray-50 transition-colors">
                                    <div className="flex items-center gap-2">
                                        <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                                        Week {weekNumber}
                                        <span className="text-sm text-gray-500 ml-2">
                                            ({weekMatches.length} matches)
                                        </span>
                                    </div>
                                </AccordionTrigger>
                                <AccordionContent className="px-6 pb-4">
                                    <div className="space-y-3">
                                        {weekMatches.map((match, idx) => (
                                            <div
                                                key={idx}
                                                className="bg-white border border-gray-200 rounded-xl p-4 shadow-sm hover:shadow-md transition-shadow"
                                            >
                                                <div className="flex items-center justify-between">
                                                    {/* Home Team */}
                                                    <div className="flex items-center gap-3 flex-1">
                                                        <div className="text-right flex-1">
                                                            <span className={`text-lg ${getTeamStyle(match.home, match)}`}>
                                                                {match.home}
                                                            </span>
                                                        </div>
                                                        {leagueId && canEdit(match, 'home') && (
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                                className="h-8 w-8 p-0 hover:bg-blue-50 text-blue-600"
                                                                onClick={() => handleEditClick(match, 'home')}
                                                            >
                                                                <Edit className="w-4 h-4" />
                                                            </Button>
                                                        )}
                                                    </div>

                                                    {/* Score */}
                                                    <div className="flex items-center gap-3 px-4">
                                                        <div className={`w-12 h-12 rounded-xl border-2 font-bold text-lg flex items-center justify-center ${getScoreStyle(match, true)}`}>
                                                            {match.homeScore}
                                                        </div>
                                                        <div className="text-2xl font-light text-gray-400">:</div>
                                                        <div className={`w-12 h-12 rounded-xl border-2 font-bold text-lg flex items-center justify-center ${getScoreStyle(match, false)}`}>
                                                            {match.awayScore}
                                                        </div>
                                                    </div>

                                                    {/* Away Team */}
                                                    <div className="flex items-center gap-3 flex-1">
                                                        {leagueId && canEdit(match, 'away') && (
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                                className="h-8 w-8 p-0 hover:bg-blue-50 text-blue-600"
                                                                onClick={() => handleEditClick(match, 'away')}
                                                            >
                                                                <Edit className="w-4 h-4" />
                                                            </Button>
                                                        )}
                                                        <div className="text-left flex-1">
                                                            <span className={`text-lg ${getTeamStyle(match.away, match)}`}>
                                                                {match.away}
                                                            </span>
                                                        </div>
                                                    </div>
                                                </div>

                                                {/* Match Status */}
                                                {match.homeScore === match.awayScore ? (
                                                    <div className="mt-3 text-center">
                                                        <span className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-amber-100 text-amber-800">
                                                            Draw
                                                        </span>
                                                    </div>
                                                ) : (
                                                    <div className="mt-3 text-center">
                                                        <span className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-green-100 text-green-800">
                                                            Winner: {match.winner}
                                                        </span>
                                                    </div>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                </AccordionContent>
                            </AccordionItem>
                        ))}
                    </Accordion>
                </CardContent>
            </Card>

        {/* Edit Dialog */}
            <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2">
                            <Edit className="w-5 h-5" />
                            Edit Match Result
                        </DialogTitle>
                    </DialogHeader>

                    {selectedMatch && editingTeam && (
                        <div className="space-y-4">
                            <div className="bg-gray-50 p-4 rounded-lg">
                                <div className="text-center text-sm text-gray-600 mb-2">Editing Result For</div>
                                <div className="text-center font-bold text-lg text-blue-600">
                                    {editingTeam === 'home' ? selectedMatch.home : selectedMatch.away}
                                </div>
                                <div className="text-center text-sm text-gray-500 mt-1">
                                    vs {editingTeam === 'home' ? selectedMatch.away : selectedMatch.home}
                                </div>
                                {selectedMatch.homeScore === selectedMatch.awayScore && (
                                    <div className="text-center text-xs text-amber-600 mt-2 font-medium">
                                        If you set equal scores, this team will be marked as winner
                                    </div>
                                )}
                            </div>

                            <div>
                                <Label htmlFor="goals">New Goal Count</Label>
                                <Input
                                    id="goals"
                                    type="number"
                                    value={newGoals}
                                    onChange={(e) => setNewGoals(e.target.value)}
                                    min="0"
                                    className="mt-1"
                                />
                                <div className="text-xs text-gray-500 mt-1">
                                    Opponent team: {editingTeam === 'home' ? selectedMatch.awayScore : selectedMatch.homeScore} goals
                                </div>
                            </div>

                            {error && (
                                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
                                    {error}
                                </div>
                            )}
                        </div>
                    )}

                    <DialogFooter>
                        <Button
                            variant="outline"
                            onClick={() => setEditDialogOpen(false)}
                            disabled={updating}
                        >
                            Cancel
                        </Button>
                        <Button
                            onClick={handleSaveEdit}
                            disabled={updating || !newGoals}
                            className="bg-green-600 hover:bg-green-700 text-black"
                        >
                            {updating ? 'Saving...' : 'Save'}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
};

export default MatchResults;