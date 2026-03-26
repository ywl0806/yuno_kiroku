import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useGetMe } from '@/feature/settings/hooks/use-get-me'
import { useUpdateMe } from '@/feature/settings/hooks/use-update-me'
import { ChevronLeft } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

export const SettingsAccountPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { data: me } = useGetMe()
  const { mutate: updateMe, isPending } = useUpdateMe()

  const [name, setName] = useState('')

  useEffect(() => {
    if (me?.name) setName(me.name)
  }, [me])

  const handleSave = () => {
    updateMe({ name })
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    navigate('/login')
  }

  return (
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.account.title')}</span>
      </div>

      <div className="space-y-5 px-4 pt-6">
        <div className="space-y-1.5">
          <Label htmlFor="account-name">{t('settings.account.nameLabel')}</Label>
          <div className="flex gap-2">
            <Input
              id="account-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t('settings.account.namePlaceholder')}
            />
            <Button disabled={isPending} onClick={handleSave}>
              {t('common.save')}
            </Button>
          </div>
        </div>

        {me?.provider && (
          <div className="space-y-1.5">
            <Label>{t('settings.account.social')}</Label>
            <p className="text-sm text-muted-foreground capitalize">{me.provider}</p>
          </div>
        )}

        <div className="pt-4">
          <Button variant="destructive" className="w-full" onClick={handleLogout}>
            {t('settings.account.logout')}
          </Button>
        </div>
      </div>
    </div>
  )
}
