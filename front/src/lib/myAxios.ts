import axios from 'axios'

const baseURL = import.meta.env.VITE_API_URL as string

export const MyAxios = axios.create({
  baseURL: baseURL,
})

const MyAxiosWithAuth = axios.create({
  baseURL: baseURL,
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
