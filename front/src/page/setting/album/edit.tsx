import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { AlbumGroupPermissionSelector } from '@/feature/settings/components/album-group-permission-selector'
import { useDeleteAlbum } from '@/feature/settings/hooks/use-delete-album'
import { useGetAlbum } from '@/feature/settings/hooks/use-get-album'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { useUpdateAlbum } from '@/feature/settings/hooks/use-update-album'
import { AlbumGroupPermission } from '@/types'
import { ChevronLeft } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate, useParams } from 'react-router-dom'

export const SettingsAlbumEditPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { albumId } = useParams<{ albumId: string }>()
  const id = Number(albumId)

  const { data: album } = useGetAlbum(id)
  const { data: settingsData } = useGetSettingsData()
  const { mutate: updateAlbum, isPending: isUpdating } = useUpdateAlbum()
  const { mutate: deleteAlbum, isPending: isDeleting } = useDeleteAlbum()

  const [name, setName] = useState('')
  const [permissions, setPermissions] = useState<AlbumGroupPermission[]>([])

  useEffect(() => {
    if (album) {
      setName(album.name)
      setPermissions(album.permissions ?? [])
    }
  }, [album])

  const handleSave = () => {
    if (!name.trim()) return
    updateAlbum({ id, name, permissions }, { onSuccess: () => navigate(-1) })
  }

  const handleDelete = () => {
    if (!window.confirm(t('settings.album.deleteConfirm'))) return
    deleteAlbum(id, { onSuccess: () => navigate('/settings') })
  }

  return (
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.album.editTitle')}</span>
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
