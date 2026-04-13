import { UploadBatchWithThumbnails } from '@/types'
import { FC } from 'react'
import { useTranslation } from 'react-i18next'

type Props = {
  batch: UploadBatchWithThumbnails
  onClick: () => void
}

export const BatchCard: FC<Props> = ({ batch, onClick }) => {
  const { t } = useTranslation()

  const date = new Date(batch.upload_at)
  const dateLabel = `${date.getFullYear()}.${String(date.getMonth() + 1).padStart(2, '0')}.${String(date.getDate()).padStart(2, '0')}`

  return (
    <div className="cursor-pointer border-b px-4 py-3 active:bg-accent" onClick={onClick}>
      <div className="mb-2 flex items-center justify-between">
        <span className="text-sm font-medium">{dateLabel}</span>
        <span className="text-xs text-muted-foreground">{t('recent.count', { count: batch.count })}</span>
      </div>
      <div className="flex gap-1">
        {batch.thumbnails.map((item) => (
          <div key={item.id} className="aspect-square flex-1 overflow-hidden rounded-sm bg-muted">
            <img src={item.thumbnail_url} alt="" className="h-full w-full object-cover" loading="lazy" />
          </div>
        ))}
        {Array.from({ length: Math.max(0, 5 - batch.thumbnails.length) }).map((_, i) => (
          <div key={`empty-${i}`} className="aspect-square flex-1 rounded-sm bg-muted" />
        ))}
      </div>
    </div>
  )
}
