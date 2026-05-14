import { BottomSheet } from '@/components/ui/bottom-sheet'
import { Button } from '@/components/ui/button'
import { Album as AlbumType } from '@/types'
import { Album } from 'lucide-react'
import { useTranslation } from 'react-i18next'

type Props = {
  albums: AlbumType[] | undefined
  open: boolean
  onClose: () => void
  selectedAlbum: AlbumType | null
  onSelect: (id: string) => void
}

export function AlbumSelectSheet({ albums, open, onClose, selectedAlbum, onSelect }: Props) {
  const { t } = useTranslation()

  return (
    <BottomSheet open={open} onClose={onClose} title={t('upload.selectAlbum')}>
      <div className="flex flex-col gap-2 px-4 pb-6 pt-2">
        {albums?.map((album) => (
          <Button
            key={album.id}
            variant={selectedAlbum?.id === album.id ? 'default' : 'outline'}
            className="h-12 justify-start text-[1rem]"
            onClick={() => onSelect(album.id)}
          >
            <Album className="size-4 shrink-0" />
            {album.is_common ? t('settings.album.commonName') : album.name}
          </Button>
        ))}
      </div>
    </BottomSheet>
  )
}
