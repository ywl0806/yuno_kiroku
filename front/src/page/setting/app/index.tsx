import { ChevronLeft } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

const LANGUAGES = [
  { code: 'kr', label: '한국어', key: 'settings.app.korean' },
  { code: 'jp', label: '日本語', key: 'settings.app.japanese' },
] as const

export const SettingsAppPage = () => {
  const { t, i18n } = useTranslation()
  const navigate = useNavigate()

  const currentLang = i18n.language

  const handleChange = (code: string) => {
    i18n.changeLanguage(code)
  }

  return (
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.app.title')}</span>
      </div>

      <div className="space-y-1 px-4 pt-6">
        <p className="mb-2 text-sm font-medium">{t('settings.app.language')}</p>
        {LANGUAGES.map((lang) => (
          <button
            key={lang.code}
            onClick={() => handleChange(lang.code)}
            className={`flex w-full items-center justify-between rounded-md px-3 py-3 text-sm transition-colors ${
              currentLang === lang.code ? 'bg-accent font-medium' : 'hover:bg-accent/50'
            }`}
          >
            <span>{t(lang.key)}</span>
            {currentLang === lang.code && <span className="text-primary">✓</span>}
          </button>
        ))}
      </div>
    </div>
  )
}
