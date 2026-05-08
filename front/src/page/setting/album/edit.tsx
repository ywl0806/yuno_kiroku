import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { AlbumGroupPermissionSelector } from '@/feature/settings/components/album-group-permission-selector'
import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { useDeleteAlbum } from '@/feature/settings/hooks/use-delete-album'
import { useGetAlbum } from '@/feature/settings/hooks/use-get-album'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { useUpdateAlbum } from '@/feature/settings/hooks/use-update-album'
import { AlbumGroupPermission } from '@/types'
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

  const isCommon = album?.is_common ?? false

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
    <SettingsSubPageLayout title={t('settings.album.editTitle')}>
      <div className="px-4 pt-3 space-y-3">
        {isCommon && (
          <p className="rounded-xl bg-amber-50 px-4 py-3 text-sm text-amber-700">
            {t('settings.album.commonNotEditable')}
          </p>
        )}

        <div className="rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04] divide-y divide-stone-100">
          <div className="px-4 py-4 space-y-1.5">
            <Label htmlFor="album-name">{t('settings.album.nameLabel')}</Label>
            <Input
              id="album-name"
              value={isCommon ? t('settings.album.commonName') : name}
              onChange={(e) => !isCommon && setName(e.target.value)}
              placeholder={t('settings.album.namePlaceholder')}
              className="border-stone-200 focus-visible:ring-stone-400"
              disabled={isCommon}
            />
          </div>

          {!isCommon && settingsData?.groups && settingsData.groups.length > 0 && (
            <div className="px-4 py-4 space-y-1.5">
              <Label>{t('settings.album.permissionsLabel')}</Label>
              <AlbumGroupPermissionSelector
                groups={settingsData.groups}
                value={permissions}
                onChange={setPermissions}
              />
            </div>
          )}
        </div>

        {!isCommon && (
          <Button
            className="w-full h-11 rounded-xl bg-stone-900 hover:bg-stone-800 font-medium mt-2"
            disabled={isUpdating || !name.trim()}
            onClick={handleSave}
          >
            {t('common.save')}
          </Button>
        )}
      </div>

      {!isCommon && (
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
      )}
    </SettingsSubPageLayout>
  )
}
