import i18n from 'i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import { initReactI18next } from 'react-i18next'

import jp from './locales/jp.json'
import kr from './locales/kr.json'

// 브라우저 언어 코드(ja, ko)도 매핑해 감지되도록 함
const resources = {
  ko: { translation: kr },
  ja: { translation: jp },
  jp: { translation: jp },
  kr: { translation: kr },
  'ja-JP': { translation: jp },
  'ko-KR': { translation: kr },
} as const

export const supportedLanguages = ['jp', 'kr', 'ja-JP', 'ko-KR'] as const
export type SupportedLanguage = (typeof supportedLanguages)[number]

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    interpolation: {
      escapeValue: false,
    },
  })

export default i18n
