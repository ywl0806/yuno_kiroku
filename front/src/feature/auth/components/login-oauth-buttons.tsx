import { Button } from '@/components/ui/button'
import { AUTH_OAUTH_URLS } from '@/feature/auth/consts'
import { useTranslation } from 'react-i18next'
import { useSearchParams } from 'react-router-dom'

export function LoginOAuthButtons() {
  const { t } = useTranslation()
  const [searchParams] = useSearchParams()
  const inviteToken = searchParams.get('invite_token')
  return (
    <>
      <Button
        type="button"
        variant="outline"
        className="bg-line text-white hover:bg-line-hover hover:text-white"
        onClick={() => (window.location.href = `${AUTH_OAUTH_URLS.line}?invite_token=${inviteToken}`)}
      >
        {t('auth.loginWithLine')}
      </Button>
      <Button
        type="button"
        variant="outline"
        className="bg-kakao text-kakao-text hover:bg-kakao-hover hover:text-kakao-text"
        onClick={() => (window.location.href = `${AUTH_OAUTH_URLS.kakao}?invite_token=${inviteToken}`)}
      >
        {t('auth.loginWithKakao')}
      </Button>
    </>
  )
}
