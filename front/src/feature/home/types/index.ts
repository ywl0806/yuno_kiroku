export const UPLOAD_STATUS = {
  IDLE: 'idle',
  PENDING: 'pending',
  SUCCESS: 'success',
  ERROR: 'error',
} as const

export type UploadStatus = (typeof UPLOAD_STATUS)[keyof typeof UPLOAD_STATUS]

export type UploadPhotoStatus = {
  [key: string]: UploadStatus
}

export type UploadPhoto = {
  file: File
  src: string
}
