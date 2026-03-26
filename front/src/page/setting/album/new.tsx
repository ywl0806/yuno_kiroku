import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { AlbumGroupPermissionSelector } from '@/feature/settings/components/album-group-permission-selector'
import { useCreateAlbum } from '@/feature/settings/hooks/use-create-album'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { AlbumGroupPermission } from '@/types'
import { ChevronLeft } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

export const SettingsAlbumNewPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { data: settingsData } = useGetSettingsData()
  const { mutate: createAlbum, isPending } = useCreateAlbum()

  const [name, setName] = useState('')
  const [permissions, setPermissions] = useState<AlbumGroupPermission[]>([])

  const handleSubmit = () => {
    if (!name.trim()) return
    createAlbum({ name, permissions }, { onSuccess: () => navigate(-1) })
  }

  return (
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.album.newTitle')}</span>
      </div>

      <div className="space-y-5 px-4 pt-6">
        <div className="space-y-1.5">
          <Label htmlFor="album-name">{t('settings.album.nameLabel')}</Label>
          <Input
            id="album-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder={t('settings.album.namePlaceholder')}
          />
        </div>

        {settingsData?.groups && settingsData.groups.length > 0 && (
          <div className="space-y-1.5">
            <Label>{t('settings.album.permissionsLabel')}</Label>
            <AlbumGroupPermissionSelector
              groups={settingsData.groups}
              value={permissions}
              onChange={setPermissions}
            />
          </div>
        )}

        <Button className="w-full" disabled={isPending || !name.trim()} onClick={handleSubmit}>
          {t('common.save')}
        </Button>
      </div>
    </div>
  )
}
