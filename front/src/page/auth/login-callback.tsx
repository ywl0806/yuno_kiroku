import { LoginCallbackContainer } from '@/feature/auth/components/login-callback-container'
import { useSearchParams } from 'react-router-dom'

export const LoginCallbackPage = () => {
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token')
  return <LoginCallbackContainer token={token} />
}
