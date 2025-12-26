import { HomeHeader } from '@/feature/home/components/home-header'
import { PhotoSwiper } from '@/feature/home/components/photo-swiper'
import { useMediaItemsRange } from '@/feature/home/hooks/use-media-items-range'
import { FC } from 'react'

type Props = {
  date: string
}
export const HomeContainer: FC<Props> = ({ date }) => {
  const { range, isFetched } = useMediaItemsRange()
  return (
    <div className="flex h-full flex-col">
      <HomeHeader date={date ?? ''} range={range} rangeFetched={isFetched} />
      <PhotoSwiper date={date ?? ''} range={range} />
    </div>
  )
}
