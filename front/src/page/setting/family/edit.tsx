import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useDeleteGroup } from '@/feature/settings/hooks/use-delete-group'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { useUpdateGroup } from '@/feature/settings/hooks/use-update-group'
import { ChevronLeft } from 'lucide-react'
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
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.group.editTitle')}</span>
      </div>

      <div className="space-y-5 px-4 pt-6">
        <div className="space-y-1.5">
          <Label htmlFor="group-name">{t('settings.group.nameLabel')}</Label>
          <Input
            id="group-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder={t('settings.group.namePlaceholder')}
          />
        </div>

        <Button className="w-full" disabled={isUpdating || !name.trim()} onClick={handleSave}>
          {t('common.save')}
        </Button>

        <div className="pt-4">
          <Button variant="destructive" className="w-full" disabled={isDeleting} onClick={handleDelete}>
            {t('common.delete')}
          </Button>
        </div>
      </div>
    </div>
  )
}
