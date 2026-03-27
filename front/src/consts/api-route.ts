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
    UPDATE_GROUP: (id: number) => `/user/${id}/group`,
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
    INVITE: (familyId: number) => `/family/${familyId}/invite`,
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
    BY_ID: (id: number) => `/album/${id}`,
    UPDATE: (id: number) => `/album/${id}`,
    DELETE: (id: number) => `/album/${id}`,
  },

  // MediaItem
  MEDIA_ITEM: {
    LIST: '/media-item',
    SEARCH: '/media-item/search',
    UPLOAD: '/media-item/image/upload',
    RANGE: '/media-item/range',
    UPLOAD_BATCH: '/media-item/upload-batch',
    UPLOAD_BATCH_STATUS: '/media-item/upload-batch/status',
  },
}
