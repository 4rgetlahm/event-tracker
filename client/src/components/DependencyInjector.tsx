import { ChakraProvider } from "@chakra-ui/react";
import { createContext } from "react";
import useAuth from "../hooks/useAuth";

export const Authentication = createContext({
    authCode: null as string | null,
    login: () => {},
    logout: () => {},
    setAuthCode: (code: string) => {},
});

export default function DependencyInjector({ children } : { children: JSX.Element}) {

    const auth = useAuth();
    
    return (
        <ChakraProvider>
            <Authentication.Provider value={auth}>
                {children}
            </Authentication.Provider>
        </ChakraProvider>
    )
}