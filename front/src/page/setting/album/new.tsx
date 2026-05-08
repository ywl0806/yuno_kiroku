import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { AlbumGroupPermissionSelector } from '@/feature/settings/components/album-group-permission-selector'
import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { useCreateAlbum } from '@/feature/settings/hooks/use-create-album'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { AlbumGroupPermission } from '@/types'
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
    <SettingsSubPageLayout title={t('settings.album.newTitle')}>
      <div className="px-4 pt-3 space-y-3">
        <div className="rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04] divide-y divide-stone-100">
          <div className="px-4 py-4 space-y-1.5">
            <Label htmlFor="album-name">{t('settings.album.nameLabel')}</Label>
            <Input
              id="album-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t('settings.album.namePlaceholder')}
              className="border-stone-200 focus-visible:ring-stone-400"
            />
          </div>

          {settingsData?.groups && settingsData.groups.length > 0 && (
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
