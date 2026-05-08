import { setAuthToken } from '@/feature/auth/lib/set-auth-token'
import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

type Props = {
  token: string | null
}

/**
 * OAuth 콜백. 백엔드가 /login/callback?token=... 으로 리다이렉트한 뒤
 * 토큰을 저장하고 홈으로 이동한다.
 */
export function LoginCallbackContainer({ token }: Props) {
  const { t } = useTranslation()
  const navigate = useNavigate()

  useEffect(() => {
    if (!token) {
      navigate('/login?error=missing_token', { replace: true })
      return
    }
    setAuthToken(token)
    navigate('/', { replace: true })
  }, [token, navigate])

  return (
    <div className="flex h-screen w-screen items-center justify-center">
      <p className="text-muted-foreground">{t('auth.loggingIn')}</p>
    </div>
  )
}
