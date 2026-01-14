import { LivePhoto } from '@/components/blocks/live-photo'
import { HeaderContainer } from '@/components/layouts/header-container'
import { MediaItem } from '@/types'
import CloseIcon from '@mui/icons-material/Close'
import { Dialog, IconButton, Slide } from '@mui/material'
import { TransitionProps } from '@mui/material/transitions'
import { FC, forwardRef, useEffect, useState } from 'react'
import { LazyLoadImage } from 'react-lazy-load-image-component'
import { Swiper, SwiperClass, SwiperSlide } from 'swiper/react'

type Props = {
  photos: MediaItem[]
  index: number
  setIndex: (index: number) => void
  open: boolean
  onClose: () => void
}

const Transition = forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement<any, any>
  },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="up" ref={ref} {...props} />
})

export const PhotoDetailSwipeDialog: FC<Props> = ({ photos, index, setIndex, open, onClose }) => {
  const [swiper, setSwiper] = useState<SwiperClass | null>(null)

  useEffect(() => {
    if (swiper) {
      swiper.slideTo(index)
    }
  }, [index, swiper])

  return (
    <Dialog open={open} onClose={onClose} fullScreen TransitionComponent={Transition}>
      <div className="h-screen overflow-hidden">
        <HeaderContainer className=" flex h-[4rem] items-center">
          <IconButton onClick={onClose}>
            <CloseIcon fontSize="large" />
          </IconButton>
        </HeaderContainer>

        <div className="my-[2rem] h-[90vh] overflow-y-hidden">
          <Swiper
            slidesPerView={1}
            onSlideChange={(swiper) => {
              setIndex(swiper.activeIndex)
            }}
            controller={{ control: swiper }}
            onSwiper={(swiper) => setSwiper(swiper)}
            className="h-full w-full"
            wrapperClass="h-full w-full"
          >
            {photos.map((photo) => (
              <SwiperSlide key={photo.id} className="flex h-full w-full items-center justify-center px-1 pb-5">
                <div className="flex h-full w-full items-center justify-center">
                  {photo.live_url ? (
                    <LivePhoto photo={photo} />
                  ) : (
                    <LazyLoadImage
                      className="mx-auto max-h-[90vh] object-cover"
                      src={photo.view_url}
                      alt={photo.file_name}
                      effect="blur"
                      onClick={() => {
                        window.open(photo.view_url, '_blank')
                      }}
                    />
                  )}
                </div>
              </SwiperSlide>
            ))}
          </Swiper>
        </div>
      </div>
    </Dialog>
  )
}
