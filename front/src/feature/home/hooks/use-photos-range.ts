import { API_ROUTES } from '@/constants/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { PhotoRange } from '@/types'
import { useQuery } from '@tanstack/react-query'
import { useCallback, useMemo } from 'react'

export const usePhotosRange = () => {
  const fetchFunc = useCallback(async () => {
    const response = await MyAxiosWithAuth.get<PhotoRange[]>(API_ROUTES.PHOTO.RANGE)
    return response.data
  }, [])

  const query = useQuery({
    queryKey: ['photosRange'],
    queryFn: () => {
      return fetchFunc()
    },
  })

  const range = useMemo(() => {
    if (!query.data) return []
    const sotedData = query.data.sort((a, b) => {
      if (a.year === b.year) {
        return b.month - a.month
      }
      return b.year - a.year
    })

    return sotedData
  }, [query.data])

  return { ...query, range }
}
