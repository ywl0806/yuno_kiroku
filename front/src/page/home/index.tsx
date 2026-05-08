import { HomeContainer } from '@/feature/home/components/home-container'
import { useMediaItemsRange } from '@/feature/home/hooks/use-media-items-range'
import { Loader } from 'lucide-react'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useParams } from 'react-router-dom'

export const HomePage = () => {
  const { t } = useTranslation()
  const { date } = useParams<{ date: string }>()
  const { range } = useMediaItemsRange()
  if (!range)
    return (
      <div className="flex h-full flex-col items-center justify-center gap-2">
        <Loader className="animate-spin" aria-hidden />
        <span className="text-sm text-muted-foreground">{t('common.loading')}</span>
      </div>
    )

  const defaultDate = useMemo(() => {
    if (!range || range.length === 0) return `${new Date().getFullYear()}-${new Date().getMonth() + 1}`
    return `${range[0].year}-${range[0].month}`
  }, [range])

  return <HomeContainer date={date ?? defaultDate} range={range} />
}
