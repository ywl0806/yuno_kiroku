import { API_ROUTES } from '@/constants/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { MediaItemRange } from '@/types'
import { useQuery } from '@tanstack/react-query'
import { useCallback, useMemo } from 'react'

export const useMediaItemsRange = () => {
  const fetchFunc = useCallback(async () => {
    const response = await MyAxiosWithAuth.get<MediaItemRange[]>(API_ROUTES.MEDIA_ITEM.RANGE)
    return response.data
  }, [])

  const query = useQuery({
    queryKey: ['mediaItemsRange'],
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
