import { getTags } from '@/service/tag-service'
import { useQuery } from '@tanstack/react-query'

export const useGetTags = () => {
  return useQuery({
    queryKey: ['tags'],
    queryFn: getTags,
    staleTime: 1000 * 60 * 5,
  })
}
