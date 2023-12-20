import { useEffect, useState } from 'react'
import {
    getApiUrl,
    getAuthUrl,
    getClientId,
    getRedirectUri,
} from '../environment/environment'
import { useLocation, useNavigate } from 'react-router-dom'
import { TokenStore } from '../models/tokenStore'

export interface AuthenticationContextParams {
    accessToken: string | null
    tokenStore: TokenStore | null
    login: () => void
    logout: () => void
    setAccessToken: (accessToken: string | null) => void
}

export default function useAuth() {
    const [accessToken, setAccessToken] = useState<string | null>(
        localStorage.getItem('access_token') || null
    )
    const [tokenStore, setTokenStore] = useState<TokenStore | null>(
        readTokenStore()
    )

    const navigate = useNavigate()

    function readTokenStore() {
        const tokenStore = localStorage.getItem('token_store')
        if (tokenStore) {
            const tokenStoreObj = JSON.parse(tokenStore)
            return tokenStoreObj
        }
        return null
    }

    useEffect(() => {
        if (!tokenStore) {
            getIdentity()
        }
    }, [])

    useEffect(() => {
        if (accessToken) {
            localStorage.setItem('access_token', accessToken)
            return
        }
        localStorage.removeItem('access_token')
    }, [accessToken])

    useEffect(() => {
        if (tokenStore) {
            localStorage.setItem('token_store', JSON.stringify(tokenStore))
            return
        }

        localStorage.removeItem('token_store')
    }, [tokenStore])

    useEffect(() => {
        if (accessToken && !tokenStore) {
            getIdentity()
        }

        if (tokenStore) {
            const expiration = new Date(tokenStore.expiration)
            const now = new Date()
            if (now > expiration) {
                logout()
            }
        }
    }, [accessToken, tokenStore])

    function getIdentity() {
        if (!accessToken) {
            return
        }
        fetchIdentity(accessToken).then((identity) => {
            if (!identity) {
                logout()
            }
            setTokenStore(identity)
        })
    }

    function login() {
        window.location.href = `${getAuthUrl()}?client_id=${getClientId()}&redirect_uri=${getRedirectUri()}&response_type=token&scope=openid%20email%20profile`
    }

    function logout() {
        setAccessToken(null)
        setTokenStore(null)
        navigate('/')
    }

    return { accessToken, tokenStore, login, logout, setAccessToken }
}

async function fetchIdentity(accessToken: string) {
    const response = await fetch(`${getApiUrl()}/v1/identity`, {
        method: 'GET',
        headers: {
            Authorization: `Bearer ${accessToken}`,
        },
    })

    if (response.status !== 200) {
        return null
    }

    const data = await response.json()
    console.log(data)
    return data
}
