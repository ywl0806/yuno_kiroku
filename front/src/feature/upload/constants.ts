import { UploadStatus } from '@/enums'

export const UPLOAD_BATCH_STORAGE_KEY = 'yuno_upload_batch'

export type StoredBatch = {
  batchId: number
  totalCount: number
}

export type UploadBatchStatusResponse = {
  statuses: {
    id: string
    upload_status: UploadStatus
    failure_reason?: string
  }[]
  is_completed: boolean
}
