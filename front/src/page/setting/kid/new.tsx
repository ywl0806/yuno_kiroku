import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { IdentityFaceSelector } from '@/feature/settings/components/identity-face-selector'
import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { useCreateKid } from '@/feature/settings/hooks/use-create-kid'
import { useGetIdentityFaceOptions } from '@/feature/settings/hooks/use-get-identity-face-options'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

export const SettingsKidNewPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { mutate: createKid, isPending } = useCreateKid()
  const { data: faceOptions = [] } = useGetIdentityFaceOptions(true)

  const [name, setName] = useState('')
  const [birthDate, setBirthDate] = useState('')
  const [selectedIdentityId, setSelectedIdentityId] = useState<number | null>(null)

  const handleSubmit = () => {
    createKid(
      {
        name: name || null,
        birth_date: birthDate || null,
        identity_id: selectedIdentityId,
      },
      { onSuccess: () => navigate(-1) },
    )
  }

  return (
    <SettingsSubPageLayout title={t('settings.kid.newTitle')}>
      <div className="px-4 pt-3 space-y-3">
        <div className="rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04] divide-y divide-stone-100">
          <div className="px-4 py-4 space-y-1.5">
            <Label htmlFor="kid-name">{t('settings.kid.nameLabel')}</Label>
            <Input
              id="kid-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t('settings.kid.namePlaceholder')}
              className="border-stone-200 focus-visible:ring-stone-400"
            />
          </div>

          <div className="px-4 py-4 space-y-1.5">
            <Label htmlFor="kid-birth-date">{t('settings.kid.birthDateLabel')}</Label>
            <Input
              id="kid-birth-date"
              type="date"
              value={birthDate}
              onChange={(e) => setBirthDate(e.target.value)}
              className="border-stone-200 focus-visible:ring-stone-400"
            />
          </div>

          {faceOptions.length > 0 && (
            <div className="px-4 py-4 space-y-1.5">
              <Label>{t('settings.kid.faceLabel')}</Label>
              <IdentityFaceSelector
                options={faceOptions}
                selectedIdentityId={selectedIdentityId}
                onSelect={setSelectedIdentityId}
              />
            </div>
          )}
        </div>

        <Button
          className="w-full h-11 rounded-xl bg-stone-900 hover:bg-stone-800 font-medium mt-2"
          disabled={isPending}
          onClick={handleSubmit}
        >
          {t('common.save')}
        </Button>
      </div>
    </SettingsSubPageLayout>
  )
}
