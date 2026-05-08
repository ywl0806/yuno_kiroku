export const SESSION_STORAGE_KEY = {
  UPLOAD_ALBUM_ID: 'upload_album_id',
} as const

export const setSessionStorage = (
  key: (typeof SESSION_STORAGE_KEY)[keyof typeof SESSION_STORAGE_KEY],
  value: string,
) => {
  sessionStorage.setItem(key, value)
}

export const getSessionStorage = (key: (typeof SESSION_STORAGE_KEY)[keyof typeof SESSION_STORAGE_KEY]) => {
  return sessionStorage.getItem(key)
}

export const removeSessionStorage = (key: (typeof SESSION_STORAGE_KEY)[keyof typeof SESSION_STORAGE_KEY]) => {
  sessionStorage.removeItem(key)
}
