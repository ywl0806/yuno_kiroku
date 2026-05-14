import { API_ROUTES } from '@/consts/api-route'
import { MyAxiosWithAuth } from '@/lib/my-axios'
import { useCallback, useEffect, useRef, useState } from 'react'
import { UPLOAD_BATCH_STORAGE_KEY, StoredBatch, UploadBatchStatusResponse } from '../constants'

export function useBatchRecovery() {
  const [isRecovering, setIsRecovering] = useState(false)
  const [recoveredBatch, setRecoveredBatch] = useState<StoredBatch | null>(null)
  const [recoveryDone, setRecoveryDone] = useState(false)
  const recoveryPollRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const cancelRecovery = useCallback(() => {
    if (recoveryPollRef.current) clearTimeout(recoveryPollRef.current)
    localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
    setIsRecovering(false)
    setRecoveredBatch(null)
  }, [])

  const pollRecoveredBatch = useCallback(
    async (batch: StoredBatch) => {
      const maxAttempts = 60
      let attempts = 0

      const poll = async () => {
        try {
          const response = await MyAxiosWithAuth.get<UploadBatchStatusResponse>(
            `${API_ROUTES.MEDIA_ITEM.UPLOAD_BATCH_STATUS}?upload_batch_id=${batch.batchId}`,
          )

          if (response.data.is_completed) {
            localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
            setIsRecovering(false)
            setRecoveredBatch(null)
            setRecoveryDone(true)
            setTimeout(() => setRecoveryDone(false), 5000)
            return
          }

          attempts++
          if (attempts < maxAttempts) {
            recoveryPollRef.current = setTimeout(poll, 5000)
          } else {
            cancelRecovery()
          }
        } catch {
          console.error('복구 상태 조회 실패')
          cancelRecovery()
        }
      }

      poll()
    },
    [cancelRecovery],
  )

  useEffect(() => {
    const stored = localStorage.getItem(UPLOAD_BATCH_STORAGE_KEY)
    if (!stored) return
    try {
      const batch: StoredBatch = JSON.parse(stored)
      setRecoveredBatch(batch)
      setIsRecovering(true)
      pollRecoveredBatch(batch)
    } catch {
      localStorage.removeItem(UPLOAD_BATCH_STORAGE_KEY)
    }

    return () => {
      if (recoveryPollRef.current) clearTimeout(recoveryPollRef.current)
    }
  }, [])

  return { isRecovering, recoveredBatch, recoveryDone, cancelRecovery }
}
