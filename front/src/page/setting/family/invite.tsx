import { useGetFamilyInvite } from '@/feature/settings/hooks/use-get-family-invite'
import { useTranslation } from 'react-i18next'
import { useParams } from 'react-router-dom'

export const SettingsFamilyInvitePage = () => {
  const { familyId } = useParams<{ familyId: string }>()
  const { t } = useTranslation()
  const { data, isLoading, isError } = useGetFamilyInvite(familyId ?? '')

  const handleCopy = () => {
    if (!data?.inviteUrl) return
    navigator.clipboard.writeText(data.inviteUrl).then(() => {
      alert(t('settings.family.invite.copied'))
    })
  }

  return (
    <div>
      <h1>{t('settings.family.invite.title')}</h1>
      <p>{t('settings.family.invite.description')}</p>

      {isLoading && <p>{t('common.loading')}</p>}
      {isError && <p>{t('settings.family.invite.error')}</p>}

      {data?.inviteUrl && (
        <div>
          <input type="text" readOnly value={data.inviteUrl} />
          <button onClick={handleCopy}>{t('settings.family.invite.copy')}</button>
        </div>
      )}
    </div>
  )
}
