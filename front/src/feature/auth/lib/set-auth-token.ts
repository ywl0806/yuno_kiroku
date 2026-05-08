import { MyAxiosWithAuth } from '@/lib/my-axios'

export function setAuthToken(token: string): void {
  MyAxiosWithAuth.interceptors.request.use((config) => {
    config.headers.Authorization = `Bearer ${token}`
    return config
  })
  localStorage.setItem('token', token)
}
