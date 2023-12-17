import { useContext, useEffect } from "react";
import { Navigate, useNavigate, useSearchParams } from "react-router-dom";
import { Authentication } from "../../components/DependencyInjector";


export function LoginCallback(){
    const auth = useContext(Authentication);
    const navigate = useNavigate();
    const [urlParams] = useSearchParams();

    useEffect(() => {
        if(urlParams.has("code")){
            auth.setAuthCode(urlParams.get("code")!);
            navigate("/");
        }
    }, [urlParams, auth]);


    return (<></>);
}