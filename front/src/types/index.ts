import { UploadStatus } from "@/enums"

export interface MediaItem {
  id: string
  group_id: string
  album_id: string
  thumbnail_url: string
  original_url: string
  view_url: string
  live_url: string
  original_live_url: string
  original_width: number
  original_height: number
  thumbnail_width: number
  thumbnail_height: number
  view_width: number
  view_height: number
  file_name: string
  taken_at: string
  created_at: string
  updated_at: string
  created_by: string
  updated_by: string
}

export type MediaItemRange = {
  year: number
  month: number
}

export interface Album {
  id: number
  name: string
  created_at: string
  updated_at: string
}


export const UPLOAD_MEDIA_ITEM_ERROR_CODE = {
  DUPLICATE: 'duplicate',
  INTERNAL_SERVER_ERROR: 'internal_server_error',
} as const

export type UploadMediaItemError = {
  code: string
  message: string
}

export type UploadMediaItem = {
  id?: string
    media_item_id?: number
  file: File
  src: string
  status: UploadStatus
  error?: UploadMediaItemError
}
