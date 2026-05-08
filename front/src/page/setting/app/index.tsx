import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { Check } from 'lucide-react'
import { useTranslation } from 'react-i18next'

const LANGUAGES = [
  { code: 'kr', label: '한국어', key: 'settings.app.language.korean' },
  { code: 'jp', label: '日本語', key: 'settings.app.language.japanese' },
] as const

export const SettingsAppPage = () => {
  const { t, i18n } = useTranslation()
  const currentLang = i18n.language

  const handleChange = (code: string) => {
    i18n.changeLanguage(code)
  }

  return (
    <SettingsSubPageLayout title={t('settings.app.title')}>
      <div className="px-4 pt-3">
        <p className="mb-2 px-1 text-[0.7rem] font-semibold uppercase tracking-widest text-stone-400">
          {t('settings.app.language.title')}
        </p>
        <div className="overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04] divide-y divide-stone-100">
          {LANGUAGES.map((lang) => (
            <button
              key={lang.code}
              onClick={() => handleChange(lang.code)}
              className="flex w-full items-center justify-between px-4 py-3.5 text-sm transition-colors hover:bg-stone-50 active:bg-stone-100"
            >
              <span className={currentLang === lang.code ? 'font-medium text-stone-900' : 'text-stone-700'}>
                {t(lang.key)}
              </span>
              {currentLang === lang.code && <Check className="size-4 text-amber-500" strokeWidth={2.5} />}
            </button>
          ))}
        </div>
      </div>
    </SettingsSubPageLayout>
  )
}
