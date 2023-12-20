import { useContext } from 'react'
import { Authentication } from './DependencyInjector'
import {
    Box,
    Button,
    ButtonGroup,
    Center,
    Flex,
    Heading,
    Spacer,
    Text,
} from '@chakra-ui/react'

export function Navbar() {
    const auth = useContext(Authentication)

    if (auth.tokenStore == null) return <></>

    return (
        <Flex minWidth="max-content" justify={'center'} gap="2" sx={{ mt: 10 }}>
            <Text size="md">Prisijungta kaip {auth.tokenStore?.email}</Text>
            <Button
                size="sm"
                colorScheme="blue"
                onClick={() => {
                    auth.logout()
                }}
            >
                Atsijungti
            </Button>
        </Flex>
    )
}
