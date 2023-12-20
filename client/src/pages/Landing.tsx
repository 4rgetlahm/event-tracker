import { useContext, useEffect, useState } from 'react'
import { Authentication } from '../components/DependencyInjector'
import { getApiUrl } from '../environment/environment'
import { EventWithUserState } from '../models/event'
import { EventCard } from '../components/events/EventCard'
import {
    Box,
    Button,
    Center,
    Container,
    Heading,
    Spinner,
} from '@chakra-ui/react'

const itemsPerPage = 5

export default function Landing() {
    const auth = useContext(Authentication)
    const [fetching, setFetching] = useState(false)
    const [page, setPage] = useState(0)
    const [events, setEvents] = useState<EventWithUserState[]>([])

    useEffect(() => {
        fetchEvents()
    }, [page, auth])

    async function fetchEvents() {
        setFetching(true)
        setEvents([])
        const response = await fetch(
            `${getApiUrl()}/v1/events/${page * itemsPerPage}/${
                (page + 1) * itemsPerPage
            }`,
            {
                headers: {
                    Authorization: `Bearer ${auth.accessToken}`,
                },
            }
        )

        if (response.status === 401) {
            auth.logout()
            return
        }

        if (!response.ok) {
            setFetching(false)
            return
        }

        const data = await response.json()
        setFetching(false)
        if (!data.events) {
            setEvents([])
            return
        }
        setEvents(data.events)
    }

    return (
        <>
            <Container maxW="container.lg" sx={{ mt: 20 }}>
                {fetching && (
                    <Center sx={{ mt: 16, mb: 20 }}>
                        <Spinner size={'xl'} />
                    </Center>
                )}
                {events.length === 0 && !fetching && (
                    <Center sx={{ mb: 20 }}>
                        <Box sx={{ mt: 16 }}>
                            <Heading size={'lg'}>Renginių nerasta</Heading>
                        </Box>
                    </Center>
                )}
                {events.map((event) => (
                    <Box key={`box-${event.id}`} sx={{ mb: 8 }}>
                        <EventCard
                            key={event.id}
                            id={event.id}
                            title={event.title}
                            description={event.description}
                            date={new Date(Date.parse(event.eventDate))}
                            isRegistered={event.isRegistered}
                            slots={event.slots}
                            slotsLeft={event.slotsLeft}
                            onRegister={() => {
                                fetchEvents()
                            }}
                            onCancelRegistration={() => {
                                fetchEvents()
                            }}
                        />
                    </Box>
                ))}
                <Center>
                    {page > 0 && (
                        <Button
                            onClick={() => setPage(page - 1)}
                            sx={{ mr: 8 }}
                        >
                            Buvęs puslapis
                        </Button>
                    )}
                    <Button onClick={() => setPage(page + 1)}>
                        Kitas puslapis
                    </Button>
                </Center>
            </Container>
        </>
    )
}
