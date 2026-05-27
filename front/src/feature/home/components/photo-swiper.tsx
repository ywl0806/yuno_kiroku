import { PhotoGridContainer } from '@/components/blocks/photo-grid-container'
import { MediaItemRange } from '@/types'
import { FC, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Swiper as SwiperCore } from 'swiper'
import 'swiper/css'
import { Controller } from 'swiper/modules'
import { Swiper, SwiperClass, SwiperSlide } from 'swiper/react'

SwiperCore.use([Controller])

type Props = {
  date: string
  range: MediaItemRange[]
  onScroll: (scrollTop: number) => void
}
export const PhotoSwiper: FC<Props> = ({ date, range, onScroll }) => {
  const [swiper, setSwiper] = useState<SwiperClass | null>(null)

  const nav = useNavigate()

  useEffect(() => {
    if (!date) return
    const [y, m] = date.split('-')
    if (swiper) {
      const findIndex = range.findIndex((r) => r.year === parseInt(y) && r.month === parseInt(m))
      swiper.slideTo(findIndex)
    }
  }, [date, range])

  return (
    <Swiper
      spaceBetween={10}
      slidesPerView={1}
      onSlideChange={(swiper) => {
        nav(`/${range[swiper.activeIndex].year}-${range[swiper.activeIndex].month}`)
      }}
      controller={{ control: swiper }}
      onSwiper={(swiper) => setSwiper(swiper)}
      className="h-full w-full"
      wrapperClass="h-full w-full"
    >
      {range.map((ran) => {
        return (
          <SwiperSlide key={`${ran.year}-${ran.month}`}>
            <PhotoGridContainer year={ran.year} month={ran.month} isActive={date === `${ran.year}-${ran.month}`} onScroll={onScroll} />
          </SwiperSlide>
        )
      })}
    </Swiper>
  )
}
