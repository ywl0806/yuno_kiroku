import { LivePhoto } from '@/components/blocks/live-photo'
import { Button } from '@/components/ui/button'
import { FullScreenModal } from '@/components/ui/full-screen-modal'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useDeleteMediaItem } from '@/feature/home/hooks/use-delete-media-item'
import { useGetAlbumOptions } from '@/feature/home/hooks/use-get-album-options'
import { useUpdateMediaItemAlbum } from '@/feature/home/hooks/use-update-media-item-album'
import { useLikeMutation } from '@/feature/like/hooks/use-like-mutation'
import { useGetMe } from '@/feature/settings/hooks/use-get-me'
import { TagSelector } from '@/feature/tag/components/tag-selector'
import { MediaItem } from '@/types'
import { Download, Heart, MoreVertical, Trash2, X } from 'lucide-react'
import { FC, useEffect, useMemo, useState } from 'react'
import { LazyLoadImage } from 'react-lazy-load-image-component'
import { Zoom } from 'swiper/modules'
import { Swiper, SwiperClass, SwiperSlide } from 'swiper/react'

type Props = {
  photos: MediaItem[]
  index: number
  setIndex: (index: number) => void
  open: boolean
  onClose: () => void
}

export const PhotoDetailSwipeDialog: FC<Props> = ({ photos, index, setIndex, open, onClose }) => {
  const [swiper, setSwiper] = useState<SwiperClass | null>(null)
  const [isLiked, setIsLiked] = useState(false)
  const [likeMaps, setLikeMaps] = useState<Record<string, boolean>>({})
  const [menuOpen, setMenuOpen] = useState(false)
  const [showAlbumSelect, setShowAlbumSelect] = useState(false)

  const currentPhoto = useMemo(() => photos[index], [photos, index])
  const currentMediaItemId = currentPhoto ? currentPhoto.id : null

  const { toggle: toggleLike, isPending: isLikePending } = useLikeMutation()
  const { data: me } = useGetMe()
  const { data: albumOptions } = useGetAlbumOptions()
  const deleteMutation = useDeleteMediaItem()
  const updateAlbumMutation = useUpdateMediaItemAlbum()

  const hasWritePermission = useMemo(() => {
    if (!me || !currentPhoto) return false
    return me.writable_album_ids.includes(currentPhoto.album_id)
  }, [me, currentPhoto])

  useEffect(() => {
    if (currentPhoto) {
      const like = likeMaps[currentPhoto.id] ? likeMaps[currentPhoto.id] : currentPhoto.is_liked ?? false
      setIsLiked(like)
    }
  }, [currentPhoto?.id])

  useEffect(() => {
    if (swiper) {
      swiper.slideTo(index)
    }
  }, [index, swiper])

  const handleLike = () => {
    if (!currentMediaItemId || isLikePending) return
    const next = !isLiked
    setIsLiked(next)
    setLikeMaps((prev) => ({ ...prev, [currentMediaItemId]: next }))
    toggleLike(currentMediaItemId, isLiked)
  }

  const handleDelete = () => {
    if (!currentMediaItemId) return
    if (!confirm('이 사진을 삭제하시겠습니까?')) return
    setMenuOpen(false)
    deleteMutation.mutate(currentMediaItemId, {
      onSuccess: () => onClose(),
    })
  }

  const handleAlbumChange = (albumId: string) => {
    if (!currentMediaItemId) return
    setMenuOpen(false)
    setShowAlbumSelect(false)
    updateAlbumMutation.mutate({ mediaItemId: currentMediaItemId, albumId })
  }

  const handleDownload = () => {
    if (!currentPhoto?.original_url) return
    setMenuOpen(false)
    const a = document.createElement('a')
    a.href = currentPhoto.original_url
    a.download = currentPhoto.file_name || 'photo'
    a.target = '_blank'
    a.click()
  }

  return (
    <FullScreenModal open={open}>
      <div className="relative h-screen">
        <div className='flex items-center gap-3 px-4 py-3 h-[4rem]'>
          {/* 닫기 버튼 */}
          <Button className="rounded-full" variant="default" size="icon" onClick={onClose}>
            <X className="size-6" />
          </Button>

          {/* 좋아요 버튼 */}
          <button
            className="flex size-10 items-center justify-center rounded-full bg-black/30 backdrop-blur-sm transition-colors hover:bg-black/50 disabled:opacity-50"
            onClick={handleLike}
            disabled={isLikePending}
          >
            <Heart className={`size-5 transition-colors ${isLiked ? 'fill-red-500 text-red-500' : 'text-white'}`} />
          </button>

          <div className="flex-1" />

          {/* 더보기 메뉴 */}
          <Popover open={menuOpen} onOpenChange={(o) => { setMenuOpen(o); if (!o) setShowAlbumSelect(false) }}>
            <PopoverTrigger asChild>
              <button className="flex size-10 items-center justify-center rounded-full bg-black/30 backdrop-blur-sm transition-colors hover:bg-black/50">
                <MoreVertical className="size-5 text-white" />
              </button>
            </PopoverTrigger>
            <PopoverContent className="w-52 p-1" align="end">
              {!showAlbumSelect ? (
                <div className="flex flex-col">
                  {/* 원본 다운로드 */}
                  <button
                    className="flex items-center gap-3 rounded-md px-3 py-2.5 text-sm hover:bg-muted transition-colors text-left"
                    onClick={handleDownload}
                  >
                    <Download className="size-4 shrink-0" />
                    원본 다운로드
                  </button>

                  {/* W 권한이 있을 때만 표시 */}
                  {hasWritePermission && (
                    <>
                      <button
                        className="flex items-center gap-3 rounded-md px-3 py-2.5 text-sm hover:bg-muted transition-colors text-left"
                        onClick={() => setShowAlbumSelect(true)}
                      >
                        <span className="size-4 shrink-0 flex items-center justify-center text-base">📁</span>
                        앨범 변경
                      </button>

                      <button
                        className="flex items-center gap-3 rounded-md px-3 py-2.5 text-sm hover:bg-muted transition-colors text-left text-destructive"
                        onClick={handleDelete}
                        disabled={deleteMutation.isPending}
                      >
                        <Trash2 className="size-4 shrink-0" />
                        삭제
                      </button>
                    </>
                  )}
                </div>
              ) : (
                <div className="flex flex-col">
                  <button
                    className="flex items-center gap-2 px-3 py-2 text-sm text-muted-foreground hover:bg-muted rounded-md transition-colors"
                    onClick={() => setShowAlbumSelect(false)}
                  >
                    ← 뒤로
                  </button>
                  <div className="my-1 border-t" />
                  <p className="px-3 py-1.5 text-xs text-muted-foreground font-medium">앨범 선택</p>
                  {albumOptions
                    ?.filter((a) => String(a.id) !== currentPhoto?.album_id)
                    .map((album) => (
                      <button
                        key={album.id}
                        className="flex items-center gap-3 rounded-md px-3 py-2.5 text-sm hover:bg-muted transition-colors text-left"
                        onClick={() => handleAlbumChange(String(album.id))}
                        disabled={updateAlbumMutation.isPending}
                      >
                        {album.name}
                      </button>
                    ))}
                </div>
              )}
            </PopoverContent>
          </Popover>
        </div>

        <div className="h-[calc(100vh-4rem)] w-full overflow-y-auto">
          <Swiper
            modules={[Zoom]}
            zoom={{ maxRatio: 4 }}
            slidesPerView={1}
            className="h-full w-full"
            wrapperClass="h-full w-full"
            onSlideChange={(swiper) => {
              setIndex(swiper.activeIndex)
            }}
            controller={{ control: swiper }}
            onSwiper={(swiper) => setSwiper(swiper)}
          >
            {photos.map((photo) => (
              <SwiperSlide key={photo.id}>
                <div className="swiper-zoom-container relative flex h-full w-full items-center justify-center">
                  {photo.video_url ? (
                    <video
                      className="max-h-[calc(100vh-5rem)] w-fit object-contain"
                      src={photo.video_url}
                      poster={photo.thumbnail_url}
                      controls
                      playsInline
                    />
                  ) : photo.live_url ? (
                    <LivePhoto photo={photo} />
                  ) : (
                    <LazyLoadImage
                      className="max-h-[calc(100vh-5rem)] object-contain"
                      src={photo.view_url}
                      alt={photo.file_name}
                      effect="blur"
                    />
                  )}
                </div>
              </SwiperSlide>
            ))}
          </Swiper>
        </div>

        {/* 태그 버튼 (하단 좌측) */}
        {currentMediaItemId && (
          <div className="absolute bottom-[calc(4rem+env(safe-area-inset-top)+0.5rem)] left-4 z-50">
            <TagSelector mediaItemId={currentMediaItemId} />
          </div>
        )}
      </div>
    </FullScreenModal>
  )
}
