import { getPhotos } from '@/service/getPhotos'
import { useQuery } from '@tanstack/react-query'

export type UseGetPhotosProps = {
  year: number
  month: number
}
export const useGetPhotos = ({ year, month }: UseGetPhotosProps) => {
  const { data, refetch, isFetched } = useQuery({
    queryKey: ['photos', year, month],
    queryFn: () => {
      const from = new Date(year, month - 1, 1)
      const to = new Date(year, month)
      return getPhotos({ from, to })
    },
    // enabled: false,
  })

  return { photos: data, refetch, isFetched }
}
