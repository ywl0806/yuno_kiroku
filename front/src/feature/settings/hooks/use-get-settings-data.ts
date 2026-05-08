import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { SettingsData } from '@/types'
import { useQuery } from '@tanstack/react-query'

export const useGetSettingsData = () => {
  return useQuery({
    queryKey: ['settings-data'],
    queryFn: async () => {
      const response = await MyAxiosWithAuth.get<SettingsData>(API_ROUTES.SETTINGS.DATA)
      return response.data
    },
    staleTime: 1000 * 60 * 1,
  })
}
