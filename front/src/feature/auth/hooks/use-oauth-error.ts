import { OAUTH_ERROR_KEYS } from '@/feature/auth/consts'
import type { TFunction } from 'i18next'

export function getOAuthErrorMessage(
  errorParam: string | null,
  t: TFunction
): string | null {
  if (!errorParam || !OAUTH_ERROR_KEYS.includes(errorParam as (typeof OAUTH_ERROR_KEYS)[number])) {
    return null
  }
  return t(`auth.error.${errorParam}`)
}
