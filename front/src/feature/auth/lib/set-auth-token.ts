import { MyAxios } from '@/lib/my-axios'

export async function clearAuthToken(): Promise<void> {
  try {
    await MyAxios.post('/auth/logout', null, { withCredentials: true })
  } catch {}
}
