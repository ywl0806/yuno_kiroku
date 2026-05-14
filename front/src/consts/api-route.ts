export const API_ROUTES = {
  // Auth
  AUTH: {
    LOGIN: '/auth/login',
    LINE_REDIRECT: '/auth/line',
    KAKAO_REDIRECT: '/auth/kakao',
  },

  // User
  USER: {
    CREATE: '/user',
    ME: '/user/me',
    MEMBERS: '/user/members',
    MEMBER: (id: string) => `/user/${id}`,
    UPDATE: (id: string) => `/user/${id}`,
  },

  // Photo
  PHOTO: {
    LIST: '/photo',
    UPLOAD: '/photo/upload',
    UPLOAD_LIVE: '/photo/upload-live',
    RANGE: '/photo/range',
    IDENTITY_RANDOM: (identityId: number) => `/photo/identity/${identityId}/random`,
  },

  // Identity
  IDENTITY: {
    OPTIONS: '/identity/options',
    LIST: '/identity',
    BY_ID: (id: number) => `/identity/${id}`,
    UPDATE: '/identity',
  },

  // Family
  FAMILY: {
    LIST: '/family',
    INVITE: (familyId: string) => `/family/${familyId}/invite`,
  },

  // Invite
  INVITE: {
    CREATE: '/invite',
  },

  // Group (앨범 그룹)
  GROUP: {
    LIST: '/group',
    CREATE: '/group',
    UPDATE: (id: number) => `/group/${id}`,
    DELETE: (id: number) => `/group/${id}`,
  },

  // Settings
  SETTINGS: {
    DATA: '/settings/data',
    IDENTITY_FACE_OPTIONS: '/settings/identity-face-options',
  },

  // Kid
  KID: {
    CREATE: '/kid',
    UPDATE: (kidId: number) => `/kid/${kidId}`,
    DELETE: (kidId: number) => `/kid/${kidId}`,
    FACE_IMGS: '/kid/face-imgs',
  },

  // Album
  ALBUM: {
    LIST_FOR_WRITE: '/album/write',
    OPTIONS: '/album/options',
    ALL: '/album',
    CREATE: '/album',
    BY_ID: (id: string) => `/album/${id}`,
    UPDATE: (id: string) => `/album/${id}`,
    DELETE: (id: string) => `/album/${id}`,
  },

  // MediaItem
  MEDIA_ITEM: {
    LIST: '/media-item',
    SEARCH: '/media-item/search',
    UPLOAD: '/media-item/image/upload',
    PRESIGNED_URL: '/media-item/presigned-url',
    RANGE: '/media-item/range',
    UPLOAD_BATCH: '/media-item/upload-batch',
    UPLOAD_BATCH_STATUS: '/media-item/upload-batch/status',
    UPLOAD_BATCHES: '/media-item/upload-batch',
    UPLOAD_BATCH_ITEMS: (id: number) => `/media-item/upload-batch/${id}/items`,
    LIKE: (id: string) => `/media-item/${id}/like`,
    TAGS: (id: string) => `/media-item/${id}/tag`,
    TAG_REMOVE: (id: string, tagId: number) => `/media-item/${id}/tag/${tagId}`,
    DELETE: (id: string) => `/media-item/${id}`,
    UPDATE_ALBUM: (id: string) => `/media-item/${id}/album`,
  },

  // Tag
  TAG: {
    LIST: '/tag',
    CREATE: '/tag',
    DELETE: (id: number) => `/tag/${id}`,
  },
}
