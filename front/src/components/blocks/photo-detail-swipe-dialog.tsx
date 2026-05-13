import { LivePhoto } from '@/components/blocks/live-photo'
import { Button } from '@/components/ui/button'
import { FullScreenModal } from '@/components/ui/full-screen-modal'
import { useLikeMutation } from '@/feature/like/hooks/use-like-mutation'
import { TagSelector } from '@/feature/tag/components/tag-selector'
import { MediaItem } from '@/types'
import { Heart, X } from 'lucide-react'
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

  const currentPhoto = useMemo(() => photos[index], [photos, index])
  const currentMediaItemId = currentPhoto ? currentPhoto.id : null

  const { toggle: toggleLike, isPending: isLikePending } = useLikeMutation()

  // 현재 사진이 바뀌면 is_liked 초기값 설정
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

  return (
    <FullScreenModal open={open}>
      <div className="relative h-screen">
        <div className='flex items-center gap-3 px-4 py-3 h-[4rem]'>
          {/* 닫기 버튼 */}
          <Button className="" variant="default" size="icon" onClick={onClose}>
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

        {/* 태그 버튼 (상단 좌측) */}
        {currentMediaItemId && (
          <div className="absolute bottom-[calc(4rem+env(safe-area-inset-top))] left-4 z-50">
            <TagSelector mediaItemId={currentMediaItemId} />
          </div>
        )}
      </div>
    </FullScreenModal>
  )
}
