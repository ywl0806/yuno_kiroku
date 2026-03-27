import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { MediaItem, MediaItemRange } from '@/types'
import { useInfiniteQuery } from '@tanstack/react-query'

export type SearchFilter = {
  selectedAlbumId: number | null
  selectedIdentityIds: number[]
  selectedRange: MediaItemRange | null
}

type SearchResponse = {
  items: MediaItem[]
  has_next: boolean
  page: number
}

export const useSearchMediaItems = (filter: SearchFilter) => {
  const hasFilter =
    filter.selectedAlbumId !== null || filter.selectedIdentityIds.length > 0 || filter.selectedRange !== null

  return useInfiniteQuery({
    queryKey: ['search-media-items', filter],
    queryFn: async ({ pageParam = 1 }) => {
      const from = filter.selectedRange
        ? new Date(filter.selectedRange.year, filter.selectedRange.month - 1, 1)
        : undefined
      const to = filter.selectedRange ? new Date(filter.selectedRange.year, filter.selectedRange.month) : undefined

      const response = await MyAxiosWithAuth.get<SearchResponse>(API_ROUTES.MEDIA_ITEM.SEARCH, {
        params: {
          page: pageParam,
          ...(from ? { from: from.toISOString() } : {}),
          ...(to ? { to: to.toISOString() } : {}),
          ...(filter.selectedIdentityIds.length > 0 ? { identity_ids: filter.selectedIdentityIds } : {}),
          ...(filter.selectedAlbumId ? { album_id: filter.selectedAlbumId } : {}),
        },
      })
      return response.data
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage) => (lastPage.has_next ? lastPage.page + 1 : undefined),
    enabled: hasFilter,
    staleTime: 1000 * 60 * 1,
  })
}
