import { LivePhoto } from '@/components/blocks/live-photo'
import { HeaderContainer } from '@/components/layouts/header-container'
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
        <HeaderContainer className="flex h-[4rem] items-center">
          <Button variant="ghost" size="icon" onClick={onClose}>
            <X className="size-6" />
          </Button>
        </HeaderContainer>

        <div className="my-[2rem] pb-5">
          <Swiper
            slidesPerView={1}
            onSlideChange={(swiper) => {
              setIndex(swiper.activeIndex)
            }}
            controller={{ control: swiper }}
            onSwiper={(swiper) => setSwiper(swiper)}
          >
            {photos.map((photo) => (
              <SwiperSlide key={photo.id} className="flex h-full w-full items-center justify-center px-1">
                <div className="flex h-full w-full items-center justify-center overflow-y-auto">
                  {photo.live_url ? (
                    <LivePhoto photo={photo} />
                  ) : (
                    <LazyLoadImage
                      className="mx-auto max-h-[90vh] object-cover"
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
