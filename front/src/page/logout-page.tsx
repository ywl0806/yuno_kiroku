import { useTranslation } from 'react-i18next'

export const LogoutPage = () => {
  const { t } = useTranslation()
  return <div>{t('page.logout')}</div>
}
