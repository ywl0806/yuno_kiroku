import { PhotoDetailSwipeDialog } from './PhotoDetailSwipeDialog'
import { PhotoGrid } from './PhotoGrid'
import { useGetPhotos } from '@/hooks/useGetPhotos'
import { FC, useState } from 'react'

type Props = {
  year: number
  month: number
}

export const PhotoGridContainer: FC<Props> = ({ year, month }) => {
  const { photosGroups } = useGetPhotos({ year, month })
  const [detailViewIndex, setDetailViewIndex] = useState<number>(0)
  const [openDetailView, setOpenDetailView] = useState<boolean>(false)
  return (
    <>
      <PhotoGrid
        photos={photosGroups ? photosGroups[0].photos : []}
        onClick={(index) => {
          setDetailViewIndex(index)
          setOpenDetailView(true)
        }}
      />
      <PhotoDetailSwipeDialog
        photos={photosGroups ? photosGroups[0].photos : []}
        index={detailViewIndex}
        setIndex={setDetailViewIndex}
        open={openDetailView}
        onClose={() => setOpenDetailView(false)}
      />
    </>
  )
}
