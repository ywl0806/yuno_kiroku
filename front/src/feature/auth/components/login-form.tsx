import { useState } from 'react'
import { ChevronDown, ChevronUp } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Form, FormField, FormItem, FormLabel } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useLoginMutation } from '@/feature/auth/hooks/use-login-mutation'
import { LOGIN_SCHEMA, AUTH_OAUTH_URLS } from '@/feature/auth/consts'
import type { LoginForm } from '@/feature/auth/types'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { useSearchParams } from 'react-router-dom'

const LAST_LOGIN_KEY = 'lastLoginMethod'
type LoginMethod = 'email' | 'line' | 'kakao'

function getLastLoginMethod(): LoginMethod | null {
  return localStorage.getItem(LAST_LOGIN_KEY) as LoginMethod | null
}

function saveLastLoginMethod(method: LoginMethod) {
  localStorage.setItem(LAST_LOGIN_KEY, method)
}

type Props = {
  oauthError: string | null
}

export function LoginForm({ oauthError }: Props) {
  const { t } = useTranslation()
  const [searchParams] = useSearchParams()
  const inviteToken = searchParams.get('invite_token')
  const form = useForm<LoginForm>({
    defaultValues: { username: '', password: '' },
    resolver: zodResolver(LOGIN_SCHEMA),
  })
  const { mutate: login } = useLoginMutation(form)
  const [lastMethod] = useState<LoginMethod | null>(getLastLoginMethod)
  const [showSecondary, setShowSecondary] = useState(false)

  const handleEmailSubmit = (data: LoginForm) => {
    saveLastLoginMethod('email')
    login(data)
  }

  const handleLineClick = () => {
    saveLastLoginMethod('line')
    window.location.href = `${AUTH_OAUTH_URLS.line}?invite_token=${inviteToken}`
  }

  const handleKakaoClick = () => {
    saveLastLoginMethod('kakao')
    window.location.href = `${AUTH_OAUTH_URLS.kakao}?invite_token=${inviteToken}`
  }

  const lineButton = (
    <Button
      type="button"
      variant="outline"
      className="bg-line text-white hover:bg-line-hover hover:text-white"
      onClick={handleLineClick}
    >
      {t('auth.loginWithLine')}
    </Button>
  )

  const kakaoButton = (
    <Button
      type="button"
      variant="outline"
      className="bg-kakao text-kakao-text hover:bg-kakao-hover hover:text-kakao-text"
      onClick={handleKakaoClick}
    >
      {t('auth.loginWithKakao')}
    </Button>
  )

  const emailSection = (
    <>
      <FormField
        control={form.control}
        name="username"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('auth.username')}</FormLabel>
            <Input type="text" {...field} />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name="password"
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('auth.password')}</FormLabel>
            <Input type="password" {...field} />
          </FormItem>
        )}
      />
      <Button type="submit">{t('auth.login')}</Button>
    </>
  )

  let primaryContent: React.ReactNode
  let secondaryContent: React.ReactNode
  let secondaryLabel: string

  if (lastMethod === 'email') {
    primaryContent = emailSection
    secondaryContent = (
      <>
        {lineButton}
        {kakaoButton}
      </>
    )
    secondaryLabel = t('auth.loginWithOther')
  } else if (lastMethod === 'line') {
    primaryContent = lineButton
    secondaryContent = (
      <>
        {kakaoButton}
        {emailSection}
      </>
    )
    secondaryLabel = t('auth.loginWithOther')
  } else if (lastMethod === 'kakao') {
    primaryContent = kakaoButton
    secondaryContent = (
      <>
        {lineButton}
        {emailSection}
      </>
    )
    secondaryLabel = t('auth.loginWithOther')
  } else {
    primaryContent = (
      <>
        {lineButton}
        {kakaoButton}
      </>
    )
    secondaryContent = emailSection
    secondaryLabel = t('auth.loginWithEmail')
  }

  return (
    <Form {...form}>
      <form
        className="flex min-w-[400px] flex-col gap-4 rounded-2xl border-2 p-20"
        onSubmit={form.handleSubmit(handleEmailSubmit)}
      >
        <h1 className="text-2xl font-bold">{t('auth.login')}</h1>

        {oauthError && (
          <p className="text-sm text-destructive" role="alert">
            {oauthError}
          </p>
        )}

        <div className="flex flex-col gap-2">{primaryContent}</div>

        <button
          type="button"
          className="flex items-center justify-center gap-1 text-sm text-muted-foreground hover:text-foreground"
          onClick={() => setShowSecondary((v) => !v)}
        >
          {secondaryLabel}
          {showSecondary ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
        </button>

        <div className={`overflow-hidden transition-all duration-300 ${showSecondary ? 'max-h-96' : 'max-h-0'}`}>
          <div className="flex flex-col gap-2">{secondaryContent}</div>
        </div>
      </form>
    </Form>
  )
}
