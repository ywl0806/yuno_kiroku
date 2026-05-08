import { KidCircleRow } from '@/feature/home/components/kid-circle-row'
import { MediaItem } from '@/types'
import { FC, useMemo } from 'react'

type Props = {
  year: number
  month: number
  allPhotos: MediaItem[]
}

export const MonthHeroSection: FC<Props> = ({ year, month, allPhotos }) => {
  const heroPhoto = useMemo(() => {
    if (allPhotos.length === 0) return null
    return allPhotos[Math.floor(Math.random() * allPhotos.length)]
  }, [allPhotos])

  if (!heroPhoto) return null

  return (
    <div className="relative w-full">
      <img
        src={heroPhoto.view_url}
        alt="이달의 사진"
        className="w-full max-h-[90vh] object-contain"
      />
      <div className='absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/40 to-black/10 h-full pb-3 pt-8 flex items-end'>
        <KidCircleRow year={year} month={month} />
      </div>
    </div>
  )
}
