import { API_ROUTES } from '@/consts/api-route'
import type { LoginForm } from '@/feature/auth/types'
import { MyAxios } from '@/lib/my-axios'
import { useMutation } from '@tanstack/react-query'
import type { UseFormReturn } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'

export function useLoginMutation(form: UseFormReturn<LoginForm>) {
  const navigate = useNavigate()

  return useMutation({
    mutationFn: (data: LoginForm) => MyAxios.post(API_ROUTES.AUTH.LOGIN, data, { withCredentials: true }),
    onSuccess: () => {
      navigate('/')
    },
    onError: (error: Error) => {
      form.setError('root.serverError', { message: error.message })
    },
  })
}
