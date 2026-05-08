import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { useDeleteGroup } from '@/feature/settings/hooks/use-delete-group'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { useUpdateGroup } from '@/feature/settings/hooks/use-update-group'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate, useParams } from 'react-router-dom'

export const SettingsFamilyEditPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { familyId } = useParams<{ familyId: string }>()
  const groupId = Number(familyId)

  const { data: settingsData } = useGetSettingsData()
  const { mutate: updateGroup, isPending: isUpdating } = useUpdateGroup()
  const { mutate: deleteGroup, isPending: isDeleting } = useDeleteGroup()

  const group = settingsData?.groups.find((g) => g.id === groupId)
  const [name, setName] = useState('')

  useEffect(() => {
    if (group?.name) setName(group.name)
  }, [group])

  const handleSave = () => {
    if (!name.trim()) return
    updateGroup({ id: groupId, name }, { onSuccess: () => navigate(-1) })
  }

  const handleDelete = () => {
    if (!window.confirm(t('settings.group.deleteConfirm'))) return
    deleteGroup(groupId, { onSuccess: () => navigate('/settings') })
  }

  return (
    <SettingsSubPageLayout title={t('settings.group.editTitle')}>
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
          disabled={isUpdating || !name.trim()}
          onClick={handleSave}
        >
          {t('common.save')}
        </Button>
      </div>

      <div className="px-4 pt-10">
        <div className="border-t border-stone-200 pt-4">
          <button
            className="w-full py-2.5 text-sm text-red-500 transition-colors hover:text-red-600 disabled:opacity-40"
            disabled={isDeleting}
            onClick={handleDelete}
          >
            {t('common.delete')}
          </button>
        </div>
      </div>
    </SettingsSubPageLayout>
  )
}
