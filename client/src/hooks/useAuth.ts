import { useEffect, useState } from "react";
import { getAuthUrl, getClientId, getRedirectUri } from "../environment/environment";

export default function useAuth(){

    const [authCode, setAuthCode] = useState<string | null>(localStorage.getItem('auth_code') || null)

    useEffect(() => {
        if(authCode){
            localStorage.setItem('auth_code', authCode)
            return;
        }
        localStorage.removeItem('auth_code')
    }, [authCode])

    function login(){
        window.location.href = `${getAuthUrl()}?client_id=${getClientId()}&redirect_uri=${getRedirectUri()}&response_type=code&scope=openid%20email%20profile`
    }

    function logout(){
        setAuthCode(null)
    }

    return { authCode, login, logout, setAuthCode };

}