import axios, { AxiosRequestConfig } from 'axios'

const baseURL = import.meta.env.VITE_API_URL as string

/** Echo `query` 바인딩: 배열은 `key=1&key=2` 반복. Axios 기본 `key[]=1`은 매칭되지 않음. */
function serializeEchoQueryParams(params: Record<string, unknown>): string {
  const parts: string[] = []
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null) continue
    const encKey = encodeURIComponent(key)
    if (Array.isArray(value)) {
      for (const v of value) {
        parts.push(`${encKey}=${encodeURIComponent(String(v))}`)
      }
    } else {
      parts.push(`${encKey}=${encodeURIComponent(String(value))}`)
    }
  }
  return parts.join('&')
}

const axiosDefaults = {
  baseURL,
  paramsSerializer: serializeEchoQueryParams,
} as const

export const MyAxios = axios.create(axiosDefaults)

const MyAxiosWithAuth = axios.create({
  ...axiosDefaults,
  withCredentials: true,
})

let isRefreshing = false
let pendingQueue: Array<{ resolve: () => void; reject: (err: unknown) => void }> = []

function flushQueue(error: unknown) {
  for (const { resolve, reject } of pendingQueue) {
    if (error) reject(error)
    else resolve()
  }
  pendingQueue = []
}

async function tryRefresh(): Promise<void> {
  await MyAxios.post('/auth/refresh', null, { withCredentials: true })
}

MyAxiosWithAuth.interceptors.response.use(
  (res) => res,
  async (error) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }

    if (error.response?.status !== 401 || originalRequest._retry) {
      return Promise.reject(error)
    }

    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        pendingQueue.push({
          resolve: () => resolve(MyAxiosWithAuth(originalRequest)),
          reject,
        })
      })
    }

    originalRequest._retry = true
    isRefreshing = true

    try {
      await tryRefresh()
      flushQueue(null)
      return MyAxiosWithAuth(originalRequest)
    } catch (refreshError) {
      flushQueue(refreshError)
      window.location.href = '/login'
      return Promise.reject(refreshError)
    } finally {
      isRefreshing = false
    }
  },
)

export { MyAxiosWithAuth }
