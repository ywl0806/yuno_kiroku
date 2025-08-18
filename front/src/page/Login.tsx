import { MyAxios, MyAxiosWithAuth } from '../lib/myAxios'
import { Button } from '@/components/ui/button'
import { Form, FormField, FormItem, FormLabel } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { z } from 'zod'

const loginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
})

type LoginForm = z.infer<typeof loginSchema>

export const LoginPage = () => {
  const nav = useNavigate()
  const form = useForm<LoginForm>({
    defaultValues: {
      email: '',
      password: '',
    },
    resolver: zodResolver(loginSchema),
  })

  const { mutate: login } = useMutation({
    mutationFn: (data: LoginForm) => MyAxios.post('/auth/login', data),
    onSuccess: (data) => {
      MyAxiosWithAuth.interceptors.request.use((config) => {
        config.headers.Authorization = `Bearer ${data.data.token}`
        return config
      })
      localStorage.setItem('token', data.data.token)
      nav('/')
    },
    onError: (error) => {
      form.setError('root.serverError', { message: error.message })
    },
  })
  return (
    <div className="flex h-screen w-screen items-center justify-center">
      <Form {...form}>
        <form
          className="flex min-w-[400px] flex-col gap-4 rounded-2xl border-2 p-20"
          onSubmit={form.handleSubmit((data) => login(data))}
        >
          <h1 className="text-2xl font-bold">Login</h1>
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <Input type="email" {...field} />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password</FormLabel>
                <Input type="password" {...field} />
              </FormItem>
            )}
          />
          <Button type="submit">Login</Button>
        </form>
      </Form>
    </div>
  )
}
