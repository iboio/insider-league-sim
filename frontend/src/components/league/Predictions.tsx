import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from '../../components/ui/card';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '../../components/ui/table';
import type {PredictedStanding} from "@/interfaces/league.ts";

interface PredictionsProps {
    predictions: PredictedStanding[] | null | undefined;
}

export const Predictions = ({predictions}: PredictionsProps) => {
    const sortedPredictions = (predictions as PredictedStanding[]).sort((a, b) => b.odds - a.odds);
    if (!predictions || predictions.length === 0) {
        return (
            <Card className="h-full min-h-[300px]">
                <CardHeader className="bg-purple-600 py-2">
                    <CardTitle className="text-white text-sm">Championship Predictions</CardTitle>
                </CardHeader>
                <CardContent className="p-4 text-center text-gray-500 text-sm">
                    No prediction data available
                </CardContent>
            </Card>
        );
    }

    return (
        <Card className="min-h-[300px] w-full shadow-md">
            <CardHeader className="bg-purple-600 py-3">
                <CardTitle className="text-white text-base font-bold">Championship Predictions</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
                <div className="overflow-x-auto">
                    <Table className="w-full">
                        <TableHeader>
                            <TableRow className="bg-gray-100 border-b border-gray-200">
                                <TableHead className="text-sm py-3 text-center font-semibold">#</TableHead>
                                <TableHead className="text-sm py-3 text-center font-semibold">Team</TableHead>
                                <TableHead className="text-sm py-3 text-center font-semibold">Odds</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {sortedPredictions.map((team, index) => (
                                <TableRow
                                    key={`${team.teamName}-${index}`}
                                    className={`hover:bg-gray-50 ${index === 0 ? 'bg-purple-50' : ''}`}
                                >
                                    <TableCell className="text-sm text-center font-medium text-gray-700 py-3">
                                        {index + 1}
                                    </TableCell>
                                    <TableCell className="text-sm font-medium py-3 text-center">
                                        {team.teamName}
                                    </TableCell>
                                    <TableCell className="text-sm text-center py-3">
                        <span className={`font-semibold px-2 py-1 rounded-full ${
                            team.odds > 0.5 ? 'bg-green-100 text-green-800'
                                : team.odds > 0.2 ? 'bg-yellow-100 text-yellow-800'
                                    : 'bg-red-100 text-red-800'
                        }`}>
                          {team.odds.toFixed(1)}%
                        </span>
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

export default Predictions;