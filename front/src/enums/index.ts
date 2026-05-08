export const UPLOAD_STATUS = {
    PENDING: '01',
    PROCESSING: '02',
    COMPLETED: '03',
    FAILED: '04',
    DUPLICATE: '05',
} as const

export type UploadStatus = (typeof UPLOAD_STATUS)[keyof typeof UPLOAD_STATUS]

export const MEDIA_ITEM_ROLE = {
    ORIGINAL: '01',
    THUMBNAIL: '02',
    VIEW: '03',
    LIVE: '04',
    VIDEO: '05',
} as const

export type MediaItemRole = (typeof MEDIA_ITEM_ROLE)[keyof typeof MEDIA_ITEM_ROLE]