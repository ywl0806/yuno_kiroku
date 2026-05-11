import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

export function LoginCallbackContainer() {
  const { t } = useTranslation()
  const navigate = useNavigate()

  useEffect(() => {
    navigate('/', { replace: true })
  }, [navigate])

  return (
    <div className="flex h-screen w-screen items-center justify-center">
      <p className="text-muted-foreground">{t('auth.loggingIn')}</p>
    </div>
  )
}
