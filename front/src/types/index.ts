import { UploadStatus } from '@/enums'

export interface Tag {
  id: number
  family_id: string | null
  name: string
  is_preset: boolean
}

export interface MediaItem {
  id: string
  family_id: string
  album_id: string
  thumbnail_url: string
  original_url: string
  view_url: string
  live_url: string
  video_url: string
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
  is_liked: boolean
  tags: Tag[]
}

export type MediaItemRange = {
  year: number
  month: number
}

export interface Family {
  id: string
  created_at: string
  updated_at: string
}

export type FamilyTitleType = 'dad' | 'mom' | 'grandfather' | 'grandmother' | 'uncle' | 'aunt' | 'other' | 'custom'

export const FAMILY_TITLE_OPTIONS: FamilyTitleType[] = ['dad', 'mom', 'grandfather', 'grandmother', 'uncle', 'aunt', 'other', 'custom']

export interface Member {
  id: string
  name: string
  username: string
  family_id: string
  group_id: number
  family_title: FamilyTitleType | null
  custom_family_title: string | null
}

export interface Group {
  id: number
  family_id: string
  is_admin: boolean
  name: string
  created_at: string
  updated_at: string
}

export interface Album {
  id: string
  name: string
  is_common: boolean
  created_at: string
  updated_at: string
}

export interface Kid {
  id: number
  family_id: string
  name: string | null
  birth_date: string | null
  identity_id: number | null
  created_at: string
  updated_at: string
}

export interface IdentityFaceOption {
  identity_id: number
  image_url: string
}

export interface SettingsData {
  groups: Group[]
  members: Member[]
  albums: Album[]
  kids: Kid[]
}

export interface Me {
  id: string
  name: string | null
  username: string
  provider: string | null
  is_admin: boolean
}

export interface AlbumGroupPermission {
  group_id: number
  permission: 'R' | 'W'
}

export interface AlbumWithPermissions extends Album {
  permissions: AlbumGroupPermission[]
}

export interface BatchThumbnail {
  id: string
  thumbnail_url: string
  thumbnail_width: number
  thumbnail_height: number
}

export interface UploadBatchWithThumbnails {
  id: number
  album_id: string
  upload_at: string
  count: number
  thumbnails: BatchThumbnail[]
}

export interface UploadBatchesResponse {
  items: UploadBatchWithThumbnails[]
  has_next: boolean
  page: number
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
  media_item_id?: string
  file: File
  src: string
  status: UploadStatus
  error?: UploadMediaItemError
}
