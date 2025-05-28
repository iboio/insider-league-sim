import { useEffect, useState } from 'react'
import { createLeague, getLeagues, deleteLeague } from '../services/api'
import { Button } from '../components/ui/button'
import { useNavigate } from 'react-router-dom'
import type { LeagueIdWithName } from '../interfaces/dto'

import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from '../components/ui/dialog'
import { Input } from '../components/ui/input'
import { Label } from '@/components/ui/label'

// Separate component for the empty state
const EmptyState = ({ onCreateClick }: { onCreateClick: () => void }) => (
    <div className="flex flex-col items-center justify-center rounded-2xl gap-4 p-6 bg-gray-50">
        <h1 className="text-4xl font-bold text-gray-900">League Sim</h1>
        <h2 className="text-2xl text-gray-700">Welcome to League Sim!</h2>
        <p className="max-w-md text-center text-gray-600">
            It seems you don't have any leagues yet. Create one to get started.
        </p>
        <Button className="bg-blue-600 hover:bg-blue-700 text-black w-full sm:w-auto" size="lg" onClick={onCreateClick}>
            Create League
        </Button>
    </div>
)

// Separate component for the league card
const LeagueCard = ({ 
    league, 
    onDelete, 
    isDeleting 
}: { 
    league: LeagueIdWithName, 
    onDelete: (id: string) => void,
    isDeleting: boolean 
}) => {
    const navigate = useNavigate()

    const handleCardClick = () => {
        navigate(`/league/${league.leagueId}`, {
            state: { leagueName: league.leagueName }
        })
    }

    const handleDeleteClick = (e: React.MouseEvent) => {
        e.stopPropagation()
        onDelete(league.leagueId)
    }

    return (
        <div
            className="bg-white rounded-lg shadow-md p-5 hover:shadow-lg transition-shadow cursor-pointer relative"
            onClick={handleCardClick}
        >
            <div className="absolute top-2 right-2">
                <Button 
                    variant="outline" 
                    size="sm" 
                    className="text-red-600 hover:text-red-800 hover:bg-red-50"
                    onClick={handleDeleteClick}
                    disabled={isDeleting}
                >
                    {isDeleting ? (
                        <>
                            <span className="mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-current border-r-transparent"></span>
                            Deleting...
                        </>
                    ) : 'Sil'}
                </Button>
            </div>
            <h3 className="text-lg font-semibold text-blue-600 pr-16">
                {league.leagueName || `League ${league.leagueId}`}
            </h3>
            <p className="text-sm text-gray-500 mt-2">
                ID: {league.leagueId.substring(0, 8)}...
            </p>
        </div>
    )
}

// Separate component for the create league dialog
const CreateLeagueDialog = ({
                                open,
                                setOpen,
                                leagueName,
                                setLeagueName,
                                teamCount,
                                setTeamCount,
                                onCreateLeague,
                            }: {
    open: boolean
    setOpen: (open: boolean) => void
    leagueName: string
    setLeagueName: (name: string) => void
    teamCount: string
    setTeamCount: (count: string) => void
    onCreateLeague: () => void
}) => (
    <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-[25vw] sm:max-w-[425px] mx-auto p-0">
            <div className="p-6">
                <DialogHeader className="pb-4">
                    <DialogTitle className="text-xl">Create a New League</DialogTitle>
                    <DialogDescription className="mt-2 text-base">
                        Please enter the league name and number of teams.
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4 py-4">
                    <div className="space-y-2">
                        <Label htmlFor="leagueName">League Name</Label>
                        <Input
                            id="leagueName"
                            value={leagueName}
                            onChange={(e) => setLeagueName(e.target.value)}
                            placeholder="Enter league name"
                            autoFocus
                        />
                    </div>
                    
                    <div className="space-y-2">
                        <Label htmlFor="teamCount">Team Count</Label>
                        <Input
                            id="teamCount"
                            type="number"
                            min={1}
                            value={teamCount}
                            onChange={(e) => setTeamCount(e.target.value)}
                            placeholder="Enter number of teams"
                        />
                    </div>
                </div>

                <DialogFooter className="pt-4">
                    <Button
                        onClick={() => {
                            onCreateLeague()
                            setOpen(false)
                        }}
                        className="bg-blue-600 hover:bg-blue-700 text-black"
                        disabled={!leagueName.trim() || !teamCount.trim()}
                    >
                        Create League
                    </Button>
                </DialogFooter>
            </div>
        </DialogContent>
    </Dialog>
)

function Home() {
    const [leagues, setLeagues] = useState<LeagueIdWithName[]>([])
    const [notFound, setNotFound] = useState(false)
    const [loading, setLoading] = useState(true)
    const [deleting, setDeleting] = useState(false)
    const [open, setOpen] = useState(false)
    const [leagueName, setLeagueName] = useState("")
    const [teamCount, setTeamCount] = useState("")

    useEffect(() => {
        fetchLeagues()
    }, [])

    const fetchLeagues = async () => {
        try {
            setLoading(true)
            const leaguesList = await getLeagues()
            
            if (!leaguesList || leaguesList.length === 0) {
                setLeagues([])
                setNotFound(true)
            } else {
                setLeagues(leaguesList)
                setNotFound(false)
            }
        } catch (error) {
            console.error('Error fetching leagues:', error)
            setNotFound(true)
        } finally {
            setLoading(false)
        }
    }

    const handleCreateLeague = async () => {
        try {
            if (!leagueName.trim() || isNaN(Number(teamCount))) {
                return
            }
            await createLeague(leagueName, teamCount)
            setLeagueName("")
            setTeamCount("")
            await fetchLeagues()
        } catch (error) {
            console.error('Error creating league:', error)
        }
    }
    
    const handleDeleteLeague = async (leagueId: string) => {
        try {
            setDeleting(true)
            await deleteLeague(leagueId)
            // Update state by removing the deleted league
            setLeagues(prevLeagues => prevLeagues.filter(league => league.leagueId !== leagueId))
            // If no leagues left, show empty state
            if (leagues.length === 1) {
                setNotFound(true)
            }
        } catch (error) {
            console.error('Error deleting league:', error)
        } finally {
            setDeleting(false)
        }
    }

    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
            </div>
        )
    }

    return (
        <>
            <CreateLeagueDialog
                open={open}
                setOpen={setOpen}
                leagueName={leagueName}
                setLeagueName={setLeagueName}
                teamCount={teamCount}
                setTeamCount={setTeamCount}
                onCreateLeague={handleCreateLeague}
            />

            {notFound ? (
                <EmptyState onCreateClick={() => setOpen(true)} />
            ) : (
                <div className="container flex flex-col mx-auto p-4">
                    <div className="flex flex-col justify-between items-center mb-6 p-4">
                        <h1 className="text-3xl font-bold mb-4 p-6 sm:mb-0">Your Leagues</h1>
                        <Button
                            onClick={() => setOpen(true)}
                            className="bg-blue-600 hover:bg-blue-700 text-black sm:w-auto"
                        >
                            Create New League
                        </Button>
                    </div>

                    <div className="grid justify-center align-middle grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 p-4">
                        {leagues.map((league) => (
                            <LeagueCard 
                                key={league.leagueId} 
                                league={league} 
                                onDelete={handleDeleteLeague}
                                isDeleting={deleting}
                            />
                        ))}
                    </div>
                </div>
            )}
        </>
    )
}

export default Home