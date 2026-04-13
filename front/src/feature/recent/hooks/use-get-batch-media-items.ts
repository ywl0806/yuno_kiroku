import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { MediaItem } from '@/types'
import { useInfiniteQuery } from '@tanstack/react-query'

type BatchItemsResponse = {
  items: MediaItem[]
  has_next: boolean
  page: number
}

export const useGetBatchMediaItems = (batchId: number | null) => {
  return useInfiniteQuery({
    queryKey: ['batch-media-items', batchId],
    queryFn: async ({ pageParam = 1 }) => {
      const response = await MyAxiosWithAuth.get<BatchItemsResponse>(API_ROUTES.MEDIA_ITEM.UPLOAD_BATCH_ITEMS(batchId!), {
        params: { page: pageParam },
      })
      return response.data
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage) => (lastPage.has_next ? lastPage.page + 1 : undefined),
    enabled: !!batchId,
    staleTime: 1000 * 60 * 1,
  })
}
