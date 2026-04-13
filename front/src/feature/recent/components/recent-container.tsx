import { BatchCard } from '@/feature/recent/components/batch-card'
import { useGetUploadBatches } from '@/feature/recent/hooks/use-get-upload-batches'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

export const RecentContainer = () => {
  const { t } = useTranslation()
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } = useGetUploadBatches()
  const navigate = useNavigate()
  const batches = useMemo(() => data?.pages.flatMap((p) => p.items) ?? [], [data])

  return (
    <div className="flex h-full flex-col">
      <div className="border-b px-4 py-3">
        <p className="text-lg font-bold">{t('recent.title')}</p>
      </div>
      <div className="flex-1 overflow-y-auto">
        {batches.length === 0 && !isFetchingNextPage && (
          <p className="mt-10 text-center text-sm text-muted-foreground">{t('recent.noUploads')}</p>
        )}
        {batches.map((batch) => (
          <BatchCard
            key={batch.id}
            batch={batch}
            onClick={() => navigate(`/recent/${batch.id}`)}
          />
        ))}
        {hasNextPage && (
          <div className="flex justify-center py-4">
            <button
              type="button"
              onClick={() => fetchNextPage()}
              disabled={isFetchingNextPage}
              className="rounded-full border border-border px-5 py-2 text-sm text-muted-foreground transition-colors hover:bg-accent disabled:opacity-50"
            >
              {isFetchingNextPage ? t('recent.loading') : t('recent.loadMore')}
            </button>
          </div>
        )}
      </div>

    </div>
  )
}
