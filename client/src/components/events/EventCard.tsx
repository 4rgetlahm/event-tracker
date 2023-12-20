import {
    Button,
    ButtonGroup,
    Card,
    CardBody,
    CardFooter,
    CardHeader,
    Flex,
    Heading,
    Stack,
    StackDivider,
    Text,
} from '@chakra-ui/react'
import { getApiUrl } from '../../environment/environment'
import { Authentication } from '../DependencyInjector'
import { useContext } from 'react'

export interface EventCardProps {
    id: string
    title: string
    description: string
    date: Date
    isRegistered: boolean
    slots: number
    slotsLeft: number
    onRegister?: (eventId: string) => void
    onCancelRegistration?: (eventId: string) => void
}

export function EventCard({
    id,
    title,
    description,
    date,
    isRegistered,
    slots,
    slotsLeft,
    onRegister,
    onCancelRegistration,
}: EventCardProps) {
    const auth = useContext(Authentication)

    async function registerToEvent() {
        const response = await fetch(`${getApiUrl()}/v1/event/${id}/register`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${auth.accessToken}`,
            },
        })

        if (response.status === 401) {
            auth.logout()
            return
        }

        if (!response.ok) {
            return
        }

        if (onRegister) {
            onRegister(id)
        }
    }

    async function cancelRegistration() {
        const response = await fetch(
            `${getApiUrl()}/v1/event/${id}/cancel-registration`,
            {
                method: 'POST',
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
            return
        }

        if (onCancelRegistration) {
            onCancelRegistration(id)
        }
    }

    const eventInPast = new Date() > date

    return (
        <Card variant={'outline'}>
            <CardBody>
                <Heading size="md">{title}</Heading>
                <Text pt={2} fontSize="sm">
                    {description}
                </Text>
            </CardBody>
            <CardFooter justify={'right'} flexWrap={'wrap'}>
                <Stack
                    direction="row"
                    divider={<StackDivider borderColor="gray.200" />}
                    spacing={4}
                    pt={4}
                >
                    <Flex direction="column">
                        <Heading size="xs" textTransform="uppercase">
                            Renginio data
                        </Heading>
                        <Text fontSize="sm">{date.toLocaleString()}</Text>
                    </Flex>
                    <Flex direction="column">
                        <Heading size="xs" textTransform="uppercase">
                            Vietos
                        </Heading>
                        <Text fontSize="sm">
                            {slotsLeft}/{slots}
                        </Text>
                    </Flex>
                    <ButtonGroup>
                        <>
                            {!isRegistered && !eventInPast && (
                                <Button
                                    variant="solid"
                                    colorScheme={'blue'}
                                    onClick={() => {
                                        registerToEvent()
                                    }}
                                >
                                    Registruotis
                                </Button>
                            )}
                            {isRegistered && !eventInPast && (
                                <Button
                                    variant="outline"
                                    colorScheme={'red'}
                                    onClick={() => {
                                        cancelRegistration()
                                    }}
                                >
                                    Atšaukti registraciją
                                </Button>
                            )}
                            {eventInPast && (
                                <Button variant="outline" disabled>
                                    Įvykęs renginys
                                </Button>
                            )}
                        </>
                    </ButtonGroup>
                </Stack>
            </CardFooter>
        </Card>
    )
}
