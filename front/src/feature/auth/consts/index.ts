import { API_ROUTES } from '@/consts/api-route'
import { z } from 'zod'

export const LANGUAGE_OPTIONS = [
  { value: 'jp', label: '日本語' },
  { value: 'kr', label: '한국어' },
] as const

export const LOGIN_SCHEMA = z.object({
  username: z.string().min(1),
  password: z.string().min(8),
})

export const OAUTH_ERROR_KEYS = ['line_token', 'kakao_token', 'token', 'missing_code', 'missing_token'] as const

const apiBaseUrl = (import.meta.env.VITE_API_URL as string)?.replace(/\/$/, '') ?? ''

export const AUTH_OAUTH_URLS = {
  line: `${apiBaseUrl}${API_ROUTES.AUTH.LINE_REDIRECT}`,
  kakao: `${apiBaseUrl}${API_ROUTES.AUTH.KAKAO_REDIRECT}`,
} as const
