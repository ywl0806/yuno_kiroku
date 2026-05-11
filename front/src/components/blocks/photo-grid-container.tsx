import { PhotoDetailSwipeDialog } from '@/components/blocks/photo-detail-swipe-dialog'
import { PhotoGrid } from '@/components/blocks/photo-grid'
import { MonthHeroSection } from '@/feature/home/components/month-hero-section'
import { useGetMediaItems } from '@/feature/home/hooks/use-get-media-items'
import { FC, useEffect, useMemo, useState } from 'react'
import { Photo as AlbumPhoto } from 'react-photo-album'
import { useLocation, useNavigate } from 'react-router-dom'

type Props = {
  year: number
  month: number
  date: string
  onScroll: (scrollTop: number) => void
}

type OpenMediaItemState = {
  mediaItemId: string
  year: number
  month: number
}

export const PhotoGridContainer: FC<Props> = ({ year, month, date, onScroll }) => {
  const { data: photos, refetch, isFetched } = useGetMediaItems({
    year,
    month,
    enabled: false,
  })
  const [detailViewIndex, setDetailViewIndex] = useState<number>(0)
  const [openDetailView, setOpenDetailView] = useState<boolean>(false)

  const location = useLocation()
  const nav = useNavigate()

  useEffect(() => {
    if (!date || isFetched) return
    const [y, m] = date.split('-')
    if (parseInt(y) === year && parseInt(m) === month) {
      refetch()
    }
  }, [date, year, month, refetch])

  // URL state로 특정 사진 상세뷰 자동 열기
  useEffect(() => {
    const req = location.state?.openMediaItem as OpenMediaItemState | undefined
    if (!req || !photos || req.year !== year || req.month !== month) return
    const index = photos.findIndex((p) => p.id === req.mediaItemId)
    if (index !== -1) {
      setDetailViewIndex(index)
      setOpenDetailView(true)
      // state 소비 후 제거
      nav(location.pathname, { replace: true, state: {} })
    }
  }, [photos, location.state, year, month])


  const photoAlbum: AlbumPhoto[] = useMemo(() => {
    return photos?.map((mediaItem) => {
      return {
        key: mediaItem.id,
        src: mediaItem.thumbnail_url,
        width: mediaItem.thumbnail_width,
        height: mediaItem.thumbnail_height,
        alt: mediaItem.file_name,
      }
    }) ?? []
  }, [photos])

  return (
    <div className="scroll-container h-full w-full overflow-y-auto" onScroll={(e) => onScroll(e.currentTarget.scrollTop)}>
      {photos && photos.length > 0 && (
        <MonthHeroSection year={year} month={month} allPhotos={photos} />
      )}
      <PhotoGrid
        photos={photoAlbum}
        onClick={(index) => {
          setDetailViewIndex(index)
          setOpenDetailView(true)
        }}
      />
      <PhotoDetailSwipeDialog
        photos={photos ?? []}
        index={detailViewIndex}
        setIndex={setDetailViewIndex}
        open={openDetailView}
        onClose={() => setOpenDetailView(false)}
      />
    </div>
  )
}
