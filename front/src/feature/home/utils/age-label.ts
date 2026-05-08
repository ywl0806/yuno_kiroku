import { TFunction } from 'i18next'

/**
 * 생년월일로부터 특정 시점까지 얼마나 지났는지 레이블을 반환합니다.
 * 예) 생년월일 2024-01 → "2살 1개월"
 */
export function getAgeLabel(when: Date | string, birthDate: Date | string, t: TFunction): string {
  const birth = new Date(birthDate)
  const whenDate = new Date(when)
  const diff = whenDate.getTime() - birth.getTime()
  const diffYears = diff / (1000 * 60 * 60 * 24 * 365)
  const diffMonths = Math.floor(diffYears * 12) % 12

  if (diffYears < 1) {
    return t('home.ageLabel_months', { months: diffMonths })
  }
  return (
    t('home.ageLabel_years', { years: Math.floor(diffYears) }) + ' ' + t('home.ageLabel_months', { months: diffMonths })
  )
}

/**
 *
 * 현재 시점으로부터 얼마나 지났는지 레이블을 반환합니다.
 * 예)  2026-04-10 현재 시점, 2025-03-10  → "1년 1개월전"
 * @param currentDate 현재 시점
 * @param whenDate 특정 시점
 * @param t 트랜스텔레레이터
 * @returns 레이블
 */
export function getAgoLabel(currentDate: Date | string, whenDate: Date | string, t: TFunction): string {
  const current = new Date(currentDate)
  const when = new Date(whenDate)
  const diff = current.getTime() - when.getTime()
  const diffYears = diff / (1000 * 60 * 60 * 24 * 365)
  const diffMonths = (Math.floor(diffYears * 12) % 12) + 1

  const beforeAfter = diff < 0 ? t('common.after') : t('common.before')
  if (diffYears < 1) {
    return diffMonths + ' ' + t('common.month') + beforeAfter
  }
  return (
    Math.floor(diffYears) +
    ' ' +
    t('common.year') +
    ' ' +
    beforeAfter +
    ' ' +
    diffMonths +
    ' ' +
    t('common.month') +
    beforeAfter
  )
}
