import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { Kid } from '@/types'
import { useMutation, useQueryClient } from '@tanstack/react-query'

type UpdateKidPayload = {
  kidId: number
  name?: string | null
  birth_date?: string | null
  identity_id?: number | null
}

export const useUpdateKid = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({ kidId, ...payload }: UpdateKidPayload) => {
      const response = await MyAxiosWithAuth.put<Kid>(API_ROUTES.KID.UPDATE(kidId), payload)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings-data'] })
    },
  })
}
