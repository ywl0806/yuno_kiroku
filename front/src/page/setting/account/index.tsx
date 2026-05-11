import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { clearAuthToken } from '@/feature/auth/lib/set-auth-token'
import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { useGetMe } from '@/feature/settings/hooks/use-get-me'
import { useUpdateMe } from '@/feature/settings/hooks/use-update-me'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

export const SettingsAccountPage = () => {
  const { t } = useTranslation()
  const { data: me } = useGetMe()
  const { mutate: updateMe, isPending } = useUpdateMe()

  const [name, setName] = useState('')

  useEffect(() => {
    if (me?.name) setName(me.name)
  }, [me])

  const handleSave = () => {
    updateMe({ name })
  }

  const handleLogout = async () => {
    await clearAuthToken()
    window.location.href = '/login'
  }

  return (
    <SettingsSubPageLayout title={t('settings.account.title')}>
      <div className="px-4 pt-3 space-y-3">
        <div className="rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04] divide-y divide-stone-100">
          <div className="px-4 py-4 space-y-1.5">
            <Label htmlFor="account-name">{t('settings.account.nameLabel')}</Label>
            <Input
              id="account-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t('settings.account.namePlaceholder')}
              className="border-stone-200 focus-visible:ring-stone-400"
            />
          </div>

          {me?.provider && (
            <div className="px-4 py-4 space-y-1">
              <Label>{t('settings.account.social')}</Label>
              <p className="text-sm capitalize text-stone-500">{me.provider}</p>
            </div>
          )}
        </div>

        <Button
          className="w-full h-11 rounded-xl bg-stone-900 hover:bg-stone-800 font-medium mt-2"
          disabled={isPending}
          onClick={handleSave}
        >
          {t('common.save')}
        </Button>
      </div>

      <div className="px-4 pt-10">
        <div className="border-t border-stone-200 pt-4">
          <button
            className="w-full py-2.5 text-sm text-red-500 transition-colors hover:text-red-600"
            onClick={handleLogout}
          >
            {t('settings.account.logout')}
          </button>
        </div>
      </div>
    </SettingsSubPageLayout>
  )
}
