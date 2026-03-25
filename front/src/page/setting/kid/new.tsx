import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { IdentityFaceSelector } from '@/feature/settings/components/identity-face-selector'
import { useCreateKid } from '@/feature/settings/hooks/use-create-kid'
import { useGetIdentityFaceOptions } from '@/feature/settings/hooks/use-get-identity-face-options'
import { ChevronLeft } from 'lucide-react'
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
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.kid.newTitle')}</span>
      </div>

      <div className="space-y-5 px-4 pt-6">
        <div className="space-y-1.5">
          <Label htmlFor="kid-name">{t('settings.kid.nameLabel')}</Label>
          <Input
            id="kid-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder={t('settings.kid.namePlaceholder')}
          />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="kid-birth-date">{t('settings.kid.birthDateLabel')}</Label>
          <Input id="kid-birth-date" type="date" value={birthDate} onChange={(e) => setBirthDate(e.target.value)} />
        </div>

        {faceOptions.length > 0 && (
          <div className="space-y-1.5">
            <Label>{t('settings.kid.faceLabel')}</Label>
            <IdentityFaceSelector
              options={faceOptions}
              selectedIdentityId={selectedIdentityId}
              onSelect={setSelectedIdentityId}
            />
          </div>
        )}

        <Button className="w-full" disabled={isPending} onClick={handleSubmit}>
          {t('common.save')}
        </Button>
      </div>
    </div>
  )
}
