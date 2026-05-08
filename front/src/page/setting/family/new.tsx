import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { useCreateGroup } from '@/feature/settings/hooks/use-create-group'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

export const SettingsFamilyNewPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { mutate: createGroup, isPending } = useCreateGroup()

  const [name, setName] = useState('')

  const handleSubmit = () => {
    if (!name.trim()) return
    createGroup({ name }, { onSuccess: () => navigate(-1) })
  }

  return (
    <SettingsSubPageLayout title={t('settings.group.newTitle')}>
      <div className="px-4 pt-3 space-y-3">
        <div className="rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04]">
          <div className="px-4 py-4 space-y-1.5">
            <Label htmlFor="group-name">{t('settings.group.nameLabel')}</Label>
            <Input
              id="group-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t('settings.group.namePlaceholder')}
              className="border-stone-200 focus-visible:ring-stone-400"
            />
          </div>
        </div>

        <Button
          className="w-full h-11 rounded-xl bg-stone-900 hover:bg-stone-800 font-medium mt-2"
          disabled={isPending || !name.trim()}
          onClick={handleSubmit}
        >
          {t('common.save')}
        </Button>
      </div>
    </SettingsSubPageLayout>
  )
}
