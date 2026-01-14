import { HomeContainer } from '@/feature/home/components/home-container'
import { useMediaItemsRange } from '@/feature/home/hooks/use-media-items-range'
import { Loader } from 'lucide-react'
import { useMemo } from 'react'
import { useParams } from 'react-router-dom'

export const HomePage = () => {
  const { date } = useParams<{ date: string }>()
  const { range } = useMediaItemsRange()
  if (!range)
    return (
      <div className="flex h-full items-center justify-center">
        <Loader className="animate-spin" />
      </div>
    )

  const defaultDate = useMemo(() => {
    if (!range || range.length === 0) return `${new Date().getFullYear()}-${new Date().getMonth() + 1}`
    return `${range[0].year}-${range[0].month}`
  }, [range])

  return <HomeContainer date={date ?? defaultDate} range={range} />
}
