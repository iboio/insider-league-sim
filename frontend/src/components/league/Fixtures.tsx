import { useMemo } from 'react';

// Mock interfaces for the component
interface Match {
    home: string;
    away: string;
}

interface Week {
    number: number;
    matches: Match[];
}

interface FixturesProps {
    weeks: Week[] | null | undefined;
    title: string;
}

// Group matches by week number
const groupMatchesByWeek = (weeks: Week[]): Week[] => {
    const grouped: { [key: number]: Match[] } = {};

    // Group all matches by week number
    weeks.forEach(week => {
        if (!grouped[week.number]) {
            grouped[week.number] = [];
        }
        grouped[week.number].push(...week.matches);
    });

    // Convert back to Week array and sort by week number
    return Object.entries(grouped)
        .map(([weekNumber, matches]) => ({
            number: parseInt(weekNumber),
            matches: matches
        }))
        .sort((a, b) => a.number - b.number);
};

export const Fixtures = ({ weeks, title }: FixturesProps) => {
    // Group matches by week number
    const groupedWeeks = useMemo(() => {
        if (!weeks || weeks.length === 0) return [];
        return groupMatchesByWeek(weeks);
    }, [weeks]);

    if (!weeks || weeks.length === 0) {
        return (
            <div className="bg-white rounded-lg shadow-sm border min-h-[300px]">
                <div className="bg-blue-600 py-3 px-4 rounded-t-lg">
                    <h3 className="text-white text-sm font-medium">{title}</h3>
                </div>
                <div className="p-4 text-center text-gray-500 text-sm">
                    No fixture data available
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white rounded-lg shadow-sm border min-h-[300px]">
            <div className="bg-blue-600 py-3 px-4 rounded-t-lg">
                <h3 className="text-white text-sm font-medium">{title}</h3>
            </div>
            <div className="divide-y divide-gray-200">
                {groupedWeeks.map((week) => (
                    <details key={week.number} className="group">
                        <summary className="px-4 py-3 cursor-pointer hover:bg-gray-50 flex justify-between items-center">
                            <div className="flex items-center space-x-2">
                                <span className="font-semibold text-sm">Week {week.number}</span>
                                <span className="text-xs text-gray-500">
                  ({week.matches?.length || 0} matches)
                </span>
                            </div>
                            <svg
                                className="w-4 h-4 text-gray-400 transform transition-transform group-open:rotate-180"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                            </svg>
                        </summary>
                        <div className="px-4 pb-3">
                            <div className="space-y-2">
                                {week.matches?.map((match, idx) => (
                                    <div
                                        key={idx}
                                        className="bg-gray-50 rounded-md p-3 border text-sm"
                                    >
                                        <div className="flex items-center justify-between">
                                            <div className="flex-1 text-right pr-3 text-gray-700">
                        <span className="font-medium">
                          {match.home || 'TBD'}
                        </span>
                                            </div>
                                            <div className="px-3 text-gray-400 font-semibold text-xs">
                                                VS
                                            </div>
                                            <div className="flex-1 text-left pl-3 text-gray-700">
                        <span className="font-medium">
                          {match.away || 'TBD'}
                        </span>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </details>
                ))}
            </div>
        </div>
    );
};

export default Fixtures