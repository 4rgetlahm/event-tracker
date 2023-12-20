import { ChakraProvider } from '@chakra-ui/react'
import { createContext } from 'react'
import useAuth, { AuthenticationContextParams } from '../hooks/useAuth'

export const Authentication = createContext<AuthenticationContextParams>({
    accessToken: null,
    tokenStore: null,
    login: () => {},
    logout: () => {},
    setAccessToken: () => {},
})

export default function DependencyInjector({
    children,
}: {
    children: JSX.Element
}) {
    const auth = useAuth()

    return (
        <ChakraProvider>
            <Authentication.Provider value={auth}>
                {children}
            </Authentication.Provider>
        </ChakraProvider>
    )
}
