import { API_ROUTES } from '@/constants/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Photos } from '@/service/get-photos'
import { useQuery } from '@tanstack/react-query'
import { useCallback } from 'react'

export type UseGetPhotosProps = {
  year: number
  month: number
}
export const useGetPhotos = ({ year, month }: UseGetPhotosProps) => {
  const fetchFunc = useCallback(async () => {
    const from = new Date(year, month - 1, 1)
    const to = new Date(year, month)
    const response = await MyAxiosWithAuth.get<Photos>(API_ROUTES.PHOTO.LIST, {
      params: {
        from: from.toISOString(),
        to: to.toISOString(),
      },
    })
    return response.data
  }, [year, month])

  return useQuery({
    queryKey: ['photos', year, month],
    queryFn: fetchFunc,
  })
}
