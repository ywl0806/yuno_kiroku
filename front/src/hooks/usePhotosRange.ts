import { getPhotosRange } from '@/service/getPhotosRange'
import { useQuery } from '@tanstack/react-query'
import { useMemo } from 'react'

export const usePhotosRange = () => {
  const { data, isFetched, refetch } = useQuery({
    queryKey: ['photosRange'],
    queryFn: () => {
      return getPhotosRange()
    },
  })

  const range = useMemo(() => {
    if (!data) return []
    const sotedData = data.sort((a, b) => {
      if (a.year === b.year) {
        return b.month - a.month
      }
      return b.year - a.year
    })

    return sotedData
  }, [data])

  return { range, isFetched, refetch }
}
