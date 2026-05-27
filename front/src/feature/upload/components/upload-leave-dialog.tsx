import { AlertTriangle } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { unstable_useBlocker as useBlocker } from 'react-router-dom'

type Blocker = ReturnType<typeof useBlocker>

type Props = {
  blocker: Blocker
}

export function UploadLeaveDialog({ blocker }: Props) {
  const { t } = useTranslation()

  if (blocker.state !== 'blocked') return null

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50">
      <div className="mx-4 w-full max-w-sm rounded-2xl bg-white px-6 py-5 shadow-2xl">
        <div className="mb-3 flex items-center gap-2">
          <AlertTriangle className="size-5 shrink-0 text-amber-500" />
          <span className="font-semibold text-gray-900">{t('upload.progress.leaveWarningTitle')}</span>
        </div>
        <p className="mb-5 text-sm text-gray-500">{t('upload.progress.leaveWarningBody')}</p>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => blocker.reset?.()}
            className="flex-1 rounded-lg border border-gray-200 py-2.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50"
          >
            {t('upload.progress.leaveWarningStay')}
          </button>
          <button
            type="button"
            onClick={() => blocker.proceed?.()}
            className="flex-1 rounded-lg bg-alizarin py-2.5 text-sm font-medium text-white transition-colors hover:bg-alizarin/90"
          >
            {t('upload.progress.leaveWarningLeave')}
          </button>
        </div>
      </div>
    </div>
  )
}
