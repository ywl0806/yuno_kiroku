import { CheckCircle2, Loader2, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { StoredBatch } from '../constants'

type Props = {
  isRecovering: boolean
  recoveredBatch: StoredBatch | null
  recoveryDone: boolean
  onCancel: () => void
}

const BANNER_CLASS =
  'fixed bottom-[calc(4rem+env(safe-area-inset-bottom)+0.5rem)] left-1/2 z-50 w-[320px] -translate-x-1/2 rounded-2xl border bg-white px-5 py-4 shadow-xl transition-all duration-300'

export function UploadRecoveryBanner({ isRecovering, recoveredBatch, recoveryDone, onCancel }: Props) {
  const { t } = useTranslation()

  if (recoveryDone) {
    return (
      <div className={BANNER_CLASS}>
        <div className="flex items-center gap-2">
          <CheckCircle2 className="size-5 shrink-0 text-emerald" />
          <span className="text-sm font-medium text-gray-800">{t('upload.progress.recoveryDone')}</span>
        </div>
      </div>
    )
  }

  if (!isRecovering || !recoveredBatch) return null

  return (
    <div className={BANNER_CLASS}>
      <div className="flex items-center gap-2">
        <Loader2 className="size-5 shrink-0 animate-spin text-peter-river" />
        <span className="text-sm font-medium text-gray-800">{t('upload.progress.recoveryInProgress')}</span>
        <span className="ml-auto text-xs text-gray-400">{recoveredBatch.totalCount}장</span>
        <button type="button" onClick={onCancel} className="ml-1 rounded p-0.5 text-gray-400 hover:text-gray-600">
          <X className="size-3.5" />
        </button>
      </div>
    </div>
  )
}
