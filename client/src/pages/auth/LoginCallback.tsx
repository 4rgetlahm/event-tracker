import { useContext, useEffect } from 'react'
import {
    Navigate,
    useLocation,
    useNavigate,
    useSearchParams,
} from 'react-router-dom'
import { Authentication } from '../../components/DependencyInjector'

export function LoginCallback() {
    const auth = useContext(Authentication)
    const navigate = useNavigate()
    const { hash } = useLocation()

    useEffect(() => {
        console.log(hash)

        const params = new Map()

        hash.slice(1)
            .split('&')
            .forEach((param) => {
                const [key, value] = param.split('=')
                params.set(key, value)
            })

        if (params.get('access_token') != null) {
            auth.setAccessToken(params.get('access_token')!)
            navigate('/')
        }
    }, [hash, auth])

    return <></>
}
