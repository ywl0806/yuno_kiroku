import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { useCreateInvite } from '@/feature/settings/hooks/use-create-invite'
import { useGetGroups } from '@/feature/settings/hooks/use-get-groups'
import { FAMILY_TITLE_OPTIONS, FamilyTitleType } from '@/types'
import { Check, Copy } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

export const MemberInviteContainer = () => {
  const { t } = useTranslation()
  const { data: groups } = useGetGroups()
  const { mutate: createInvite, isPending } = useCreateInvite()

  const [groupId, setGroupId] = useState<string>('')
  const [familyTitle, setFamilyTitle] = useState<FamilyTitleType | ''>('')
  const [customFamilyTitle, setCustomFamilyTitle] = useState('')
  const [inviteUrl, setInviteUrl] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

  const handleSubmit = () => {
    if (!groupId) return
    createInvite(
      {
        groupId: Number(groupId),
        familyTitle: familyTitle,
        customFamilyTitle: familyTitle === 'custom' ? customFamilyTitle : '',
      },
      {
        onSuccess: (data) => setInviteUrl(data.invite_url),
      },
    )
  }

  const handleCopy = () => {
    if (!inviteUrl) return
    navigator.clipboard.writeText(inviteUrl).then(() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    })
  }

  return (
    <SettingsSubPageLayout title={t('settings.member.inviteTitle')}>
      <div className="px-4 pt-3 space-y-3">
        <div className="rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04] divide-y divide-stone-100">
          <div className="px-4 py-4 space-y-1.5">
            <Label>{t('settings.member.inviteGroupLabel')}</Label>
            <Select value={groupId} onValueChange={setGroupId}>
              <SelectTrigger className="border-stone-200">
                <SelectValue placeholder={t('settings.member.inviteGroupPlaceholder')} />
              </SelectTrigger>
              <SelectContent>
                {groups?.map((group) => (
                  <SelectItem key={group.id} value={String(group.id)}>
                    {group.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="px-4 py-4 space-y-1.5">
            <Label>{t('settings.member.inviteTitleLabel')}</Label>
            <Select value={familyTitle} onValueChange={(v) => setFamilyTitle(v as FamilyTitleType)}>
              <SelectTrigger className="border-stone-200">
                <SelectValue placeholder={t('settings.member.inviteTitlePlaceholder')} />
              </SelectTrigger>
              <SelectContent>
                {FAMILY_TITLE_OPTIONS.map((option) => (
                  <SelectItem key={option} value={option}>
                    {t(`settings.member.familyTitle.${option}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {familyTitle === 'custom' && (
            <div className="px-4 py-4 space-y-1.5">
              <Label>{t('settings.member.customTitleLabel')}</Label>
              <Input
                value={customFamilyTitle}
                onChange={(e) => setCustomFamilyTitle(e.target.value)}
                placeholder={t('settings.member.customTitlePlaceholder')}
                className="border-stone-200 focus-visible:ring-stone-400"
              />
            </div>
          )}
        </div>

        {inviteUrl && (
          <div className="flex items-center gap-3 rounded-xl bg-stone-100 px-4 py-3">
            <span className="flex-1 truncate text-xs text-stone-500">{inviteUrl}</span>
            <button onClick={handleCopy} className="shrink-0 text-stone-400 transition-colors hover:text-stone-700">
              {copied ? <Check className="size-4 text-emerald-500" /> : <Copy className="size-4" />}
            </button>
          </div>
        )}

        <Button
          className="w-full h-11 rounded-xl bg-stone-900 hover:bg-stone-800 font-medium mt-2"
          disabled={!groupId || isPending}
          onClick={inviteUrl ? handleCopy : handleSubmit}
        >
          {inviteUrl
            ? copied
              ? t('settings.member.inviteCopied')
              : t('settings.member.inviteCopy')
            : t('settings.member.inviteGenerate')}
        </Button>
      </div>
    </SettingsSubPageLayout>
  )
}
