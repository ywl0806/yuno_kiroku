import { AlbumSelectSheet } from '@/feature/upload/components/album-select-sheet'
import { UploadLeaveDialog } from '@/feature/upload/components/upload-leave-dialog'
import { UploadProgressBanner } from '@/feature/upload/components/upload-progress-banner'
import { UploadRecoveryBanner } from '@/feature/upload/components/upload-recovery-banner'
import { useAlbumSelection } from '@/feature/upload/hooks/use-album-selection'
import { useBatchRecovery } from '@/feature/upload/hooks/use-batch-recovery'
import { useUpload } from '@/feature/upload/hooks/use-upload'
import { UploadMediaItem } from '@/types'
import { createContext, Dispatch, FC, SetStateAction, useContext } from 'react'

type UploadPhotoContextType = {
  mediaItems: UploadMediaItem[]
  setMediaItems: Dispatch<SetStateAction<UploadMediaItem[]>>
  handleUploadPhotos: (albumId: string) => Promise<void>
  clearPhotos: () => void
  progress: number
  isUploading: boolean
  reUploadPhoto: (index: number, albumId: string) => void
  isUploaded: boolean
  albumSheetOpen: boolean
  openAlbumSheet: () => void
  closeAlbumSheet: () => void
  selectedAlbumId: string | null
  selectedAlbumName: string | null
  setSelectedAlbumId: (id: string) => void
}

export const UploadPhotoContext = createContext<UploadPhotoContextType>({
  mediaItems: [],
  setMediaItems: () => {},
  handleUploadPhotos: () => Promise.resolve(),
  clearPhotos: () => {},
  progress: 0,
  isUploading: false,
  reUploadPhoto: () => {},
  isUploaded: false,
  albumSheetOpen: false,
  openAlbumSheet: () => {},
  closeAlbumSheet: () => {},
  selectedAlbumId: null,
  selectedAlbumName: null,
  setSelectedAlbumId: () => {},
})

export const UploadPhotoProvider: FC<{ children: React.ReactNode }> = ({ children }) => {
  const upload = useUpload()
  const recovery = useBatchRecovery()
  const albumSelection = useAlbumSelection()

  return (
    <UploadPhotoContext.Provider
      value={{
        mediaItems: upload.mediaItems,
        setMediaItems: upload.setMediaItems,
        handleUploadPhotos: upload.handleUploadPhotos,
        clearPhotos: upload.clearPhotos,
        progress: upload.progress,
        isUploading: upload.isUploading,
        reUploadPhoto: upload.reUploadPhoto,
        isUploaded: upload.isUploaded,
        albumSheetOpen: albumSelection.albumSheetOpen,
        openAlbumSheet: albumSelection.openAlbumSheet,
        closeAlbumSheet: albumSelection.closeAlbumSheet,
        selectedAlbumId: albumSelection.selectedAlbumId,
        selectedAlbumName: albumSelection.selectedAlbumName,
        setSelectedAlbumId: albumSelection.setSelectedAlbumId,
      }}
    >
      <UploadLeaveDialog blocker={upload.blocker} />
      <UploadProgressBanner
        isUploading={upload.isUploading}
        uploadPhase={upload.uploadPhase}
        uploadCompleted={upload.uploadCompleted}
        totalCount={upload.mediaItems.length}
        progress={upload.progress}
        statusCounts={upload.statusCounts}
      />
      <UploadRecoveryBanner
        isRecovering={recovery.isRecovering}
        recoveredBatch={recovery.recoveredBatch}
        recoveryDone={recovery.recoveryDone}
        onCancel={recovery.cancelRecovery}
      />
      <AlbumSelectSheet
        albums={albumSelection.albums}
        open={albumSelection.albumSheetOpen}
        onClose={albumSelection.closeAlbumSheet}
        selectedAlbum={albumSelection.selectedAlbum}
        onSelect={albumSelection.setSelectedAlbumId}
      />
      {children}
    </UploadPhotoContext.Provider>
  )
}

export const useUploadPhoto = () => {
  return useContext(UploadPhotoContext)
}
