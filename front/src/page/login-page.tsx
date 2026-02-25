import { LoginContainer } from '@/feature/auth/components/login-container'
import { useSearchParams } from 'react-router-dom'

export const LoginPage = () => {
  const [searchParams] = useSearchParams()
  const errorParam = searchParams.get('error')
  return <LoginContainer errorParam={errorParam} />
}
