import { getSessionStorage, SESSION_STORAGE_KEY, setSessionStorage } from '@/lib/session-storage'
import { useCallback, useMemo, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useGetAlbums } from './use-get-albums'

export function useAlbumSelection() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const location = useLocation()
  const { data: albums } = useGetAlbums()

  const [albumSheetOpen, setAlbumSheetOpen] = useState(false)
  const [selectedAlbumId, setSelectedAlbumIdState] = useState<string | null>(
    getSessionStorage(SESSION_STORAGE_KEY.UPLOAD_ALBUM_ID) || null,
  )

  const openAlbumSheet = useCallback(() => setAlbumSheetOpen(true), [])
  const closeAlbumSheet = useCallback(() => setAlbumSheetOpen(false), [])

  const setSelectedAlbumId = useCallback(
    (id: string) => {
      setSelectedAlbumIdState(id)
      setSessionStorage(SESSION_STORAGE_KEY.UPLOAD_ALBUM_ID, id)
      closeAlbumSheet()
      if (location.pathname !== '/upload') {
        navigate('/upload')
      }
    },
    [closeAlbumSheet, navigate, location.pathname],
  )

  const selectedAlbum = useMemo(() => albums?.find((a) => a.id === selectedAlbumId) ?? null, [albums, selectedAlbumId])

  const selectedAlbumName = useMemo(() => {
    if (!selectedAlbum) return null
    return selectedAlbum.is_common ? t('settings.album.commonName') : selectedAlbum.name
  }, [selectedAlbum, t])

  return {
    albums,
    albumSheetOpen,
    selectedAlbumId,
    selectedAlbum,
    selectedAlbumName,
    openAlbumSheet,
    closeAlbumSheet,
    setSelectedAlbumId,
  }
}
