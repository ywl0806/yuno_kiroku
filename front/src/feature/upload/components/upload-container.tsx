import { PreviewImageInput } from '@/components/blocks/preview-image-input'
import { Button } from '@/components/ui/button'
import { UPLOAD_STATUS } from '@/enums'
import { useUploadPhoto } from '@/providers/upload-photo-provider'
import { UploadMediaItem } from '@/types'
import { Album, ArrowLeft, ChevronDown, ImagePlus, Upload } from 'lucide-react'
import { FC, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

export const UploadContainer: FC = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const {
    mediaItems,
    setMediaItems,
    handleUploadPhotos,
    clearPhotos,
    isUploaded,
    isUploading,
    progress,
    selectedAlbumId,
    selectedAlbumName,
    openAlbumSheet,
  } = useUploadPhoto()
  const imgInputRef = useRef<HTMLInputElement>(null)

  const handleBack = () => {
    clearPhotos()
    navigate(-1)
  }

  const MAX_FILE_SIZE = 50 * 1024 * 1024 // 50MB

  const handleAddImages = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (!files) return

    const validFiles: File[] = []
    const oversizedNames: string[] = []
    Array.from(files).forEach((file) => {
      if (file.size > MAX_FILE_SIZE) {
        oversizedNames.push(file.name)
      } else {
        validFiles.push(file)
      }
    })

    if (oversizedNames.length > 0) {
      alert(t('upload.fileTooLarge', { files: oversizedNames.join(', ') }))
    }

    setMediaItems((prev: UploadMediaItem[]) => [
      ...prev,
      ...validFiles.map((file) => ({ file, src: URL.createObjectURL(file), status: UPLOAD_STATUS.PENDING })),
    ])
    e.target.value = ''
  }

  const isEmpty = mediaItems.length === 0
  const inProgressCount = mediaItems.filter(
    (m) => m.status === UPLOAD_STATUS.PENDING || m.status === UPLOAD_STATUS.PROCESSING,
  ).length

  return (
    <div className="flex h-full w-full flex-col">
      {/* 헤더 */}
      <div className="flex w-full items-center justify-between gap-2 border-b border-border px-3 py-2 shrink-0">
        {/* 뒤로가기 */}
        <Button variant="ghost" size="sm" onClick={handleBack} className="gap-1 px-2">
          <ArrowLeft className="size-4" />
          <span className="text-sm">{t('upload.back')}</span>
        </Button>

        {/* 앨범 선택/변경 버튼 */}
        <button
          type="button"
          onClick={openAlbumSheet}
          disabled={isUploaded}
          className="flex items-center gap-1.5 rounded-full border border-border px-3 py-1.5 text-sm font-medium transition-colors hover:bg-muted disabled:cursor-not-allowed disabled:opacity-60"
        >
          <Album className="size-3.5 shrink-0" />
          <span className="max-w-[140px] truncate">
            {selectedAlbumName ?? t('upload.selectAlbum')}
          </span>
          {!isUploaded && <ChevronDown className="size-3.5 shrink-0 text-muted-foreground" />}
        </button>

        {/* 업로드 버튼 */}
        <div className="flex justify-end" style={{ minWidth: '80px' }}>
          {!isUploaded && mediaItems.length > 0 && selectedAlbumId !== null && (
            <Button size="sm" onClick={() => handleUploadPhotos(selectedAlbumId)} className="gap-1.5">
              <Upload className="size-3.5" />
              <span>{t('upload.upload')}</span>
            </Button>
          )}
        </div>
      </div>

      {/* 업로드 중 진행 요약 */}
      {isUploading && progress < 100 && (
        <div className="flex items-center gap-2 bg-primary/5 px-4 py-2 text-xs text-primary shrink-0">
          <span className="animate-pulse">●</span>
          <span>
            {t('upload.progress.uploading')} {inProgressCount} / {mediaItems.length}
          </span>
        </div>
      )}

      {/* 본문 */}
      {selectedAlbumId === null ? (
        /* 앨범 미선택 상태 */
        <div className="flex flex-1 flex-col items-center justify-center gap-4 p-6 text-center">
          <Album className="size-12 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">{t('upload.selectAlbumHint')}</p>
          <Button variant="outline" onClick={openAlbumSheet}>
            {t('upload.selectAlbum')}
          </Button>
        </div>
      ) : isEmpty ? (
        /* 이미지 미선택 상태 */
        <div className="flex flex-1 flex-col items-center justify-center gap-4 p-6 text-center">
          <ImagePlus className="size-12 text-muted-foreground/40" />
          <p className="text-sm text-muted-foreground">{t('upload.noImages')}</p>
          <Button
            variant="outline"
            className="rounded-full"
            onClick={() => imgInputRef.current?.click()}
          >
            <ImagePlus className="size-4" />
            <span>{t('upload.addImage')}</span>
          </Button>
          <input
            ref={imgInputRef}
            hidden
            type="file"
            accept="image/*, video/*"
            name="images"
            multiple
            onChange={handleAddImages}
          />
        </div>
      ) : (
        /* 이미지 선택 후 */
        <div className="flex-1 overflow-y-auto">
          <PreviewImageInput
            inputRef={imgInputRef}
            images={mediaItems}
            setImages={setMediaItems}
            albumId={selectedAlbumId ?? ''}
          />
        </div>
      )}
    </div>
  )
}
