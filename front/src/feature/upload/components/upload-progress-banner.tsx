import { Progress } from '@/components/ui/progress'
import { CheckCircle2, ChevronDown, ChevronUp, Loader2 } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@radix-ui/react-collapsible'

type StatusCounts = {
  uploading: number
  processing: number
  completed: number
  failed: number
  duplicate: number
}

type Props = {
  isUploading: boolean
  uploadPhase: 'uploading' | 'processing' | 'done'
  uploadCompleted: number
  totalCount: number
  progress: number
  statusCounts: StatusCounts
}

export function UploadProgressBanner({ isUploading, uploadPhase, uploadCompleted, totalCount, progress, statusCounts }: Props) {
  const { t } = useTranslation()
  const [bannerOpen, setBannerOpen] = useState(true)

  if (!isUploading) return null

  return (
    <Collapsible open={bannerOpen} onOpenChange={setBannerOpen}>
      <div className="fixed bottom-[calc(4rem+env(safe-area-inset-bottom)+0.5rem)] left-1/2 z-50 w-[320px] -translate-x-1/2 rounded-2xl border bg-white px-5 py-4 shadow-xl transition-all duration-300">
        <div className="flex items-center gap-2">
          {uploadPhase === 'done' ? (
            <CheckCircle2 className="size-5 shrink-0 text-emerald" />
          ) : (
            <Loader2 className="size-5 shrink-0 animate-spin text-peter-river" />
          )}
          <span className="text-sm font-medium text-gray-800">
            {uploadPhase === 'uploading' && t('upload.progress.uploading')}
            {uploadPhase === 'processing' && t('upload.progress.processing')}
            {uploadPhase === 'done' && t('upload.progress.done')}
          </span>
          <span className="ml-auto text-xs text-gray-400">
            {uploadCompleted} / {totalCount}
          </span>
          <CollapsibleTrigger asChild>
            <button type="button" className="ml-1 rounded p-0.5 text-gray-400 hover:text-gray-600">
              {bannerOpen ? <ChevronDown className="size-4" /> : <ChevronUp className="size-4" />}
            </button>
          </CollapsibleTrigger>
        </div>
        <CollapsibleContent>
          <Progress value={progress} className="mt-3 h-1.5 w-full" />
          <div className="mt-2 flex gap-3 text-xs text-gray-400">
            {statusCounts.uploading > 0 && <span>{t('upload.progress.statusUploading')} {statusCounts.uploading}</span>}
            {statusCounts.processing > 0 && <span>{t('upload.progress.statusProcessing')} {statusCounts.processing}</span>}
            {statusCounts.completed > 0 && <span className="text-emerald">{t('upload.progress.statusCompleted')} {statusCounts.completed}</span>}
            {statusCounts.failed > 0 && <span className="text-alizarin">{t('upload.progress.statusFailed')} {statusCounts.failed}</span>}
            {statusCounts.duplicate > 0 && <span>{t('upload.progress.statusDuplicate')} {statusCounts.duplicate}</span>}
          </div>
        </CollapsibleContent>
      </div>
    </Collapsible>
  )
}
