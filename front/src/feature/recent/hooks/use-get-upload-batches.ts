import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { UploadBatchesResponse } from '@/types'
import { useInfiniteQuery } from '@tanstack/react-query'

export const useGetUploadBatches = () => {
  return useInfiniteQuery({
    queryKey: ['upload-batches'],
    queryFn: async ({ pageParam = 1 }) => {
      const response = await MyAxiosWithAuth.get<UploadBatchesResponse>(API_ROUTES.MEDIA_ITEM.UPLOAD_BATCHES, {
        params: { page: pageParam },
      })
      return response.data
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage) => (lastPage.has_next ? lastPage.page + 1 : undefined),
    staleTime: 1000 * 60,
  })
}
