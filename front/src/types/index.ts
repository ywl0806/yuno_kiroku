export interface Photo {
  id: string
  group_id: string
  album_id: string
  thumbnail_url: string
  original_url: string
  live_url: string
  original_live_url: string
  original_width: number
  original_height: number
  thumbnail_width: number
  thumbnail_height: number
  file_name: string
  photo_created_at: string
  created_at: string
  updated_at: string
  created_by: string
  updated_by: string
}

export type PhotoRange = {
  year: number
  month: number
}

export interface Album {
  id: number
  name: string
  created_at: string
  updated_at: string
}

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

export const UPLOAD_PHOTO_ERROR_CODE = {
  DUPLICATE: 'duplicate',
  INTERNAL_SERVER_ERROR: 'internal_server_error',
} as const

export type UploadPhotoError = {
  code: string
  message: string
}

export type UploadPhoto = {
  file: File
  src: string
  status: UploadStatus
  error?: UploadPhotoError
}
