export const API_ROUTES = {
  // Auth
  AUTH: {
    LOGIN: '/auth/login',
  },

  // User
  USER: {
    CREATE: '/user',
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
    LIST: '/identity',
    BY_ID: (id: number) => `/identity/${id}`,
    UPDATE: '/identity',
  },

  // Album
  ALBUM: {
    LIST_FOR_WRITE: '/album/write',
  },

  // MediaItem
  MEDIA_ITEM: {
    LIST: '/media-item',
    UPLOAD: '/media-item/upload',
    RANGE: '/media-item/range',
  },
}
