import axios from 'axios'

const baseURL = import.meta.env.VITE_API_URL as string

export const MyAxios = axios.create({
  baseURL: baseURL,
})

export const MyAxiosWithAuth = axios.create({
  baseURL: baseURL,
  headers: {
    Authorization: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjY4NmNiZmY2NTk4NjFkMzI1NDU5ZTI4MSIsImVtYWlsIjoidGVzdDFAZXhhbXBsZS5jb20iLCJncm91cF9pZCI6IjY4NmNiZmY2NTk4NjFkMzI1NDU5ZTI4MCIsImNsYW5fZ3JvdXBfaWQiOiIiLCJyb2xlIjoiIiwiaXNzIjoieW91ci1hcHAtbmFtZSIsImV4cCI6MTc1NTgzMjkwNywiaWF0IjoxNzUzMjQwOTA3fQ.IaEcjDe4JYnlKQl0ACS76M3bHqFYfVhd5sWW3fvQ3AE`,
    // Authorization: `Bearer ${localStorage.getItem('token')}`,
  },
})
