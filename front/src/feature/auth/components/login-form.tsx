import { Button } from '@/components/ui/button'
import { Form, FormField, FormItem, FormLabel } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useLoginMutation } from '@/feature/auth/hooks/use-login-mutation'
import { LOGIN_SCHEMA } from '@/feature/auth/consts'
import type { LoginForm } from '@/feature/auth/types'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

type Props = {
  oauthError: string | null
  oauthButtons: React.ReactNode
}

export function LoginForm({ oauthError, oauthButtons }: Props) {
  const { t } = useTranslation()
  const form = useForm<LoginForm>({
    defaultValues: { username: '', password: '' },
    resolver: zodResolver(LOGIN_SCHEMA),
  })
  const { mutate: login } = useLoginMutation(form)

  return (
    <Form {...form}>
      <form
        className="flex min-w-[400px] flex-col gap-4 rounded-2xl border-2 p-20"
        onSubmit={form.handleSubmit((data) => login(data))}
      >
        <div className="flex items-center justify-between gap-4">
          <h1 className="text-2xl font-bold">{t('auth.login')}</h1>
        </div>
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
        {oauthError && (
          <p className="text-sm text-destructive" role="alert">
            {oauthError}
          </p>
        )}
        <Button type="submit">{t('auth.login')}</Button>
        <div className="flex flex-col gap-2 pt-2">{oauthButtons}</div>
      </form>
    </Form>
  )
}
