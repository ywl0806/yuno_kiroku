import { HomeFilter } from '@/feature/home/components/home-filter-panel'
import { HomeHeader } from '@/feature/home/components/home-header'
import { PhotoSwiper } from '@/feature/home/components/photo-swiper'
import { MediaItemRange } from '@/types'
import { FC, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

type Props = {
  range: MediaItemRange[]
  date: string
}
export const HomeContainer: FC<Props> = ({ date, range }) => {
  const nav = useNavigate()
  const [filter, setFilter] = useState<HomeFilter>({
    selectedAlbumId: null,
    selectedIdentityIds: [],
  })

  useEffect(() => {
    if (!range || date) return

    if (range.length > 0) {
      nav(`/${range[0].year}-${range[0].month}`)
    }
  }, [date, range])
  return (
    <div className="flex h-full flex-col">
      <HomeHeader date={date ?? ''} range={range} filter={filter} onFilterChange={setFilter} />
      <PhotoSwiper date={date ?? ''} range={range} filter={filter} />
    </div>
  )
}
