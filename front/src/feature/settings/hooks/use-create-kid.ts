import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Kid } from '@/types'
import { useMutation, useQueryClient } from '@tanstack/react-query'

type CreateKidPayload = {
  name?: string | null
  birth_date?: string | null
  identity_id?: number | null
}

export const useCreateKid = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: CreateKidPayload) => {
      const response = await MyAxiosWithAuth.post<Kid>(API_ROUTES.KID.CREATE, payload)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
    },
  })
}
