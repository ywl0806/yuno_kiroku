import { LivePhoto } from '@/components/blocks/live-photo'
import { Button } from '@/components/ui/button'
import { FullScreenModal } from '@/components/ui/full-screen-modal'
import { MediaItem } from '@/types'
import { X } from 'lucide-react'
import { FC, useEffect, useState } from 'react'
import { LazyLoadImage } from 'react-lazy-load-image-component'
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

  useEffect(() => {
    if (swiper) {
      swiper.slideTo(index)
    }
  }, [index, swiper])

  return (
    <FullScreenModal open={open}>
      <div className="h-screen">
        <Button className='absolute top-5 right-5 z-10' variant="default" size="icon" onClick={onClose}>
          <X className="size-6" />
        </Button>


        <div className="h-full w-full">
          <Swiper
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
                <div className='h-full w-full flex items-center justify-center'>
                  {photo.live_url ? (
                    <LivePhoto photo={photo} />
                  ) : (
                    <LazyLoadImage
                      className='max-h-[calc(100vh-5rem)] object-contain'
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
      </div>
    </FullScreenModal>
  )
}
