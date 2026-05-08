import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { IdentityFaceSelector } from '@/feature/settings/components/identity-face-selector'
import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { useDeleteKid } from '@/feature/settings/hooks/use-delete-kid'
import { useGetIdentityFaceOptions } from '@/feature/settings/hooks/use-get-identity-face-options'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { useUpdateKid } from '@/feature/settings/hooks/use-update-kid'
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
    <SettingsSubPageLayout title={t('settings.kid.editTitle')}>
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
          disabled={isUpdating || isDeleting}
          onClick={handleUpdate}
        >
          {t('common.save')}
        </Button>
      </div>

      <div className="px-4 pt-10">
        <div className="border-t border-stone-200 pt-4">
          <button
            className="w-full py-2.5 text-sm text-red-500 transition-colors hover:text-red-600 disabled:opacity-40"
            disabled={isUpdating || isDeleting}
            onClick={handleDelete}
          >
            {t('common.delete')}
          </button>
        </div>
      </div>
    </SettingsSubPageLayout>
  )
}
