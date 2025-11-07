import { HeaderContainer } from '../layouts/HeaderContainer'
import { LivePhoto } from './LivePhoto'
import { Photo } from '@/types/photo'
import CloseIcon from '@mui/icons-material/Close'
import { Dialog, IconButton, Slide } from '@mui/material'
import { TransitionProps } from '@mui/material/transitions'
import { FC, forwardRef, useEffect, useState } from 'react'
import { Swiper, SwiperClass, SwiperSlide } from 'swiper/react'

type Props = {
  photos: Photo[]
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

        <div className="my-[2rem] h-[90vh] overflow-y-auto">
          <Swiper
            slidesPerView={1}
            onSlideChange={(swiper) => {
              setIndex(swiper.activeIndex)
            }}
            controller={{ control: swiper }}
            onSwiper={(swiper) => setSwiper(swiper)}
          >
            {photos.map((photo) => (
              <SwiperSlide key={photo._id} className="flex items-center justify-center px-1">
                {photo.live_url ? <LivePhoto photo={photo} /> : <img src={photo.thumbnail_url} alt={photo.file_name} />}
              </SwiperSlide>
            ))}
          </Swiper>
        </div>
      </div>
    </Dialog>
  )
}
