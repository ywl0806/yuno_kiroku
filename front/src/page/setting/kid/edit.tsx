import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { IdentityFaceSelector } from '@/feature/settings/components/identity-face-selector'
import { useDeleteKid } from '@/feature/settings/hooks/use-delete-kid'
import { useGetIdentityFaceOptions } from '@/feature/settings/hooks/use-get-identity-face-options'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { useUpdateKid } from '@/feature/settings/hooks/use-update-kid'
import { ChevronLeft } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate, useParams } from 'react-router-dom'

export const SettingsKidEditPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { kidId } = useParams<{ kidId: string }>()
  const kidIdNum = Number(kidId)

  const { data: settingsData } = useGetSettingsData()
  const kid = settingsData?.kids.find((k) => k.id === kidIdNum)

  const { mutate: updateKid, isPending: isUpdating } = useUpdateKid()
  const { mutate: deleteKid, isPending: isDeleting } = useDeleteKid()
  const { data: faceOptions = [] } = useGetIdentityFaceOptions(true, [kidIdNum])

  const [name, setName] = useState('')
  const [birthDate, setBirthDate] = useState('')
  const [selectedIdentityId, setSelectedIdentityId] = useState<number | null>(null)

  useEffect(() => {
    if (kid) {
      setName(kid.name ?? '')
      setBirthDate(kid.birth_date ? kid.birth_date.slice(0, 10) : '')
      setSelectedIdentityId(kid.identity_id ?? null)
    }
  }, [kid])

  const handleUpdate = () => {
    updateKid(
      {
        kidId: kidIdNum,
        name: name || null,
        birth_date: birthDate || null,
        identity_id: selectedIdentityId,
      },
      { onSuccess: () => navigate(-1) },
    )
  }

  const handleDelete = () => {
    if (!window.confirm(t('settings.kid.deleteConfirm'))) return
    deleteKid(kidIdNum, { onSuccess: () => navigate(-1) })
  }

  return (
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.kid.editTitle')}</span>
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

        <Button className="w-full" disabled={isUpdating || isDeleting} onClick={handleUpdate}>
          {t('common.save')}
        </Button>

        <Button className="w-full" variant="destructive" disabled={isUpdating || isDeleting} onClick={handleDelete}>
          {t('common.delete')}
        </Button>
      </div>
    </div>
  )
}
