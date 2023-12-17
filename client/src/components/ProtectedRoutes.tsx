import { useContext } from "react"
import { Authentication } from "./DependencyInjector";

export default function ProtectedRoute({ children } : { children: JSX.Element}) {

    const auth = useContext(Authentication);

    if(auth.authCode == null){
        console.log(auth.authCode)
        auth.login();
        return (<></>);
    }

    return (
        <>
            {children}
        </>
    )
}