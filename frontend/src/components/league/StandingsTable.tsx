import type {Standings} from '../../interfaces/league';
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from '../../components/ui/table';
import {Card, CardContent, CardHeader, CardTitle,} from '../../components/ui/card';

interface StandingsTableProps {
    standings: Standings[] | null | undefined;
}

export const StandingsTable = ({standings}: StandingsTableProps) => {
    const sortedStandings = (standings as Standings[]).sort((a, b) => b.points - a.points);

    if (!standings || !Array.isArray(standings) || standings.length === 0) {
        return (
            <Card className="min-h-[300px] w-full shadow-md">
                <CardHeader className="bg-blue-600 py-3">
                    <CardTitle className="text-white text-base font-bold">Standings</CardTitle>
                </CardHeader>
                <CardContent className="p-4 text-center text-gray-500 text-sm">
                    No standings data available
                </CardContent>
            </Card>
        );
    }

    return (
        <Card className="min-h-[300px] w-full shadow-md">
            <CardHeader className="bg-blue-600 py-3">
                <CardTitle className="text-white text-base font-bold">Standings</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
                <div className="overflow-x-auto">
                    <Table className="w-full">
                        <TableHeader>
                            <TableRow className="bg-gray-100 border-b border-gray-200">
                                <TableHead className="text-sm py-3 text-center font-semibold">#</TableHead>
                                <TableHead className="text-sm py-3 font-semibold">Team</TableHead>
                                <TableHead className="text-sm py-3 text-center font-semibold">Played</TableHead>
                                <TableHead className="text-sm py-3 text-center font-semibold">Win</TableHead>
                                <TableHead className="text-sm py-3 text-center font-semibold">Lose</TableHead>
                                <TableHead className="text-sm py-3 text-center font-semibold">Goal</TableHead>
                                <TableHead className="text-sm py-3 text-center font-semibold">Against</TableHead>
                                <TableHead className="text-sm py-3 text-center font-semibold">Diff</TableHead>
                                <TableHead className="text-sm py-3 text-center font-semibold">Points</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {sortedStandings.map((team, index) => (
                                <TableRow
                                    key={`${team.team.name}-${index}`}
                                    className={`hover:bg-gray-50 ${index === 0 ? 'bg-green-50' : ''}`}
                                >
                                    <TableCell className="text-sm text-center font-medium text-gray-700 py-3">
                                        {index + 1}
                                    </TableCell>
                                    <TableCell className="text-sm font-medium py-3">
                                        {team.team.name}
                                    </TableCell>
                                    <TableCell className="text-sm text-center py-3">
                                        {team.played}
                                    </TableCell>
                                    <TableCell className="text-sm text-center py-3">
                                        {team.wins}
                                    </TableCell>
                                    <TableCell className="text-sm text-center py-3">
                                        {team.losses}
                                    </TableCell>
                                    <TableCell className="text-sm text-center py-3">
                                        {team.goals}
                                    </TableCell>
                                    <TableCell className="text-sm text-center py-3">
                                        {team.against}
                                    </TableCell>
                                    <TableCell className="text-sm text-center py-3">
                                        {team.goals - team.against} </TableCell>
                                    <TableCell className="text-sm text-center font-bold py-3 text-blue-600">
                                        {team.points}
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </div>
            </CardContent>
        </Card>
    );
};

export default StandingsTable;