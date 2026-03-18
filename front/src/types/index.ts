import { UploadStatus } from '@/enums'

export interface MediaItem {
  id: string
  family_id: string
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

export interface Family {
  id: number
  name: string
  created_at: string
  updated_at: string
}

export interface Member {
  id: number
  name: string
  username: string
  family_id: number
  group_id: number
}

export interface Group {
  id: number
  family_id: number
  is_admin: boolean
  name: string
  created_at: string
  updated_at: string
}

export interface Album {
  id: number
  name: string
  created_at: string
  updated_at: string
}

export interface Kid {
  id: number
  family_id: number
  name: string | null
  birth_date: string | null
  identity_id: number | null
  created_at: string
  updated_at: string
}

export interface SettingsData {
  groups: Group[]
  members: Member[]
  albums: Album[]
  kids: Kid[]
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
