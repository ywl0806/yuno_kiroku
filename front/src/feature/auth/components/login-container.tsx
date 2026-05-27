import { LoginForm } from '@/feature/auth/components/login-form'
import { LoginLanguageSelect } from '@/feature/auth/components/login-language-select'
import { getOAuthErrorMessage } from '@/feature/auth/hooks/use-oauth-error'
import { useTranslation } from 'react-i18next'

type Props = {
  errorParam: string | null
}

export function LoginContainer({ errorParam }: Props) {
  const { t } = useTranslation()
  const oauthError = getOAuthErrorMessage(errorParam, t)
  return (
    <div className="flex h-screen w-screen items-center justify-center">
      <LoginForm oauthError={oauthError} />
      <div className="absolute top-4 right-4">
        <LoginLanguageSelect />
      </div>
    </div>
  )
}
