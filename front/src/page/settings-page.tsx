import { useTranslation } from 'react-i18next'

export const SettingsPage = () => {
  const { t } = useTranslation()
  return <div>{t('page.settings')}</div>
}
