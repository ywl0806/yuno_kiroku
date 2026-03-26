import { Label } from '@/components/ui/label'
import { useGetSettingsData } from '@/feature/settings/hooks/use-get-settings-data'
import { useUpdateMemberGroup } from '@/feature/settings/hooks/use-update-member-group'
import { ChevronLeft } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useNavigate, useParams } from 'react-router-dom'

export const SettingsMemberEditPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { memberId } = useParams<{ memberId: string }>()
  const id = Number(memberId)

  const { data: settingsData } = useGetSettingsData()
  const { mutate: updateGroup, isPending } = useUpdateMemberGroup()

  const member = settingsData?.members.find((m) => m.id === id)

  const handleGroupChange = (groupId: number) => {
    updateGroup({ id, groupId }, { onSuccess: () => navigate(-1) })
  }

  return (
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.member.editTitle')}</span>
      </div>

      <div className="space-y-5 px-4 pt-6">
        <div className="space-y-1.5">
          <p className="text-sm font-medium">{member?.name || member?.username}</p>
        </div>

        <div className="space-y-1.5">
          <Label>{t('settings.member.groupLabel')}</Label>
          <div className="space-y-2">
            {settingsData?.groups.map((group) => (
              <button
                key={group.id}
                disabled={isPending}
                onClick={() => handleGroupChange(group.id)}
                className={`flex w-full items-center justify-between rounded-md border px-3 py-2 text-sm hover:bg-accent ${member?.group_id === group.id ? 'border-primary bg-primary/5 font-medium' : ''}`}
              >
                <span>{group.name}</span>
                {member?.group_id === group.id && <span className="text-xs text-primary">{'\u2713'}</span>}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
