import axios from 'axios'

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
  headers: {
    Authorization: localStorage.getItem('token') ? `Bearer ${localStorage.getItem('token')}` : undefined,
  },
})

MyAxiosWithAuth.interceptors.response.use(
  (res) => {
    return res
  },
  (error) => {
    if (error.response.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  },
)

export { MyAxiosWithAuth }
