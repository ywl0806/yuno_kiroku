import { PhotoGrid } from './PhotoGrid'
import { useGetPhotos } from '@/hooks/useGetPhotos'
import { FC } from 'react'

type Props = {
  year: number
  month: number
}

export const PhotoGridContainer: FC<Props> = ({ year, month }) => {
  const { photosGroups } = useGetPhotos({ year, month })

  return <PhotoGrid photos={photosGroups ? photosGroups[0].photos : []} />
}
