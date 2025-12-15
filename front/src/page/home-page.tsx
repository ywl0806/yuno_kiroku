import { HomeContainer } from '@/feature/home/components/home-container'
import { useParams } from 'react-router-dom'

export const HomePage = () => {
  const { date } = useParams<{ date: string }>()

  return <HomeContainer date={date ?? `${new Date().getFullYear()}-${new Date().getMonth() + 1}`} />
}
