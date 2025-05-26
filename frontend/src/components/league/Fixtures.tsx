import type { Week } from '../../interfaces/league';
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from '../../components/ui/card';
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from '../../components/ui/accordion';

interface FixturesProps {
    weeks: Week[] | null | undefined;
    title: string;
}

export const Fixtures = ({ weeks, title }: FixturesProps) => {
    if (!weeks || weeks.length === 0) {
        return (
            <Card className="h-full min-h-[300px]">
                <CardHeader className="bg-blue-600 py-2">
                    <CardTitle className="text-white text-sm">{title}</CardTitle>
                </CardHeader>
                <CardContent className="p-4 text-center text-gray-500 text-sm">
                    No fixture data available
                </CardContent>
            </Card>
        );
    }

    return (
        <Card className="min-h-[300px]">
            <CardHeader className="bg-blue-600 py-2">
                <CardTitle className="text-white text-sm">{title}</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
                <Accordion type="single" collapsible className="w-full">
                    {weeks.map((week) => (
                        <AccordionItem key={week.number} value={`week-${week.number}`}>
                            <AccordionTrigger className="px-3 py-2 text-xs">
                                <span className="font-semibold">Week {week.number}</span>
                                <span className="text-xs text-gray-500 ml-2">
                  ({week.matches?.length || 0} matches)
                </span>
                            </AccordionTrigger>
                            <AccordionContent className="px-3 pb-2">
                                <div className="space-y-2">
                                    {week.matches?.map((match, idx) => (
                                        <div
                                            key={idx}
                                            className="bg-gray-50 rounded p-2 border text-xs"
                                        >
                                            <div className="flex items-center justify-between">
                                                <div className="flex-1 text-right pr-2 text-gray-700">
                          <span className="block font-medium">
                            {match.home?.name || 'TBD'}
                          </span>
                                                </div>
                                                <div className="px-2 text-gray-400 font-semibold">vs</div>
                                                <div className="flex-1 text-left pl-2 text-gray-700">
                          <span className="block font-medium">
                            {match.away?.name || 'TBD'}
                          </span>
                                                </div>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </AccordionContent>
                        </AccordionItem>
                    ))}
                </Accordion>
            </CardContent>
        </Card>
    );
};

export default Fixtures;