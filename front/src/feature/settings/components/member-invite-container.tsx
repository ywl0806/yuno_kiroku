import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useCreateInvite } from '@/feature/settings/hooks/use-create-invite'
import { useGetGroups } from '@/feature/settings/hooks/use-get-groups'
import { FAMILY_TITLE_OPTIONS, FamilyTitleType } from '@/types'
import { Check, ChevronLeft, Copy } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

export const MemberInviteContainer = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
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
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.member.inviteTitle')}</span>
      </div>

      <div className="space-y-5 px-4 pt-6">
        <div className="space-y-1.5">
          <Label>{t('settings.member.inviteGroupLabel')}</Label>
          <Select value={groupId} onValueChange={setGroupId}>
            <SelectTrigger>
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

        <div className="space-y-1.5">
          <Label>{t('settings.member.inviteTitleLabel')}</Label>
          <Select value={familyTitle} onValueChange={(v) => setFamilyTitle(v as FamilyTitleType)}>
            <SelectTrigger>
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
          <div className="space-y-1.5">
            <Label>{t('settings.member.customTitleLabel')}</Label>
            <Input
              value={customFamilyTitle}
              onChange={(e) => setCustomFamilyTitle(e.target.value)}
              placeholder={t('settings.member.customTitlePlaceholder')}
            />
          </div>
        )}

        {inviteUrl && (
          <div className="flex items-center gap-2 rounded-md bg-muted px-3 py-2">
            <span className="flex-1 truncate text-xs text-muted-foreground">{inviteUrl}</span>
            <button onClick={handleCopy} className="shrink-0 text-muted-foreground hover:text-foreground">
              {copied ? <Check className="size-4 text-green-500" /> : <Copy className="size-4" />}
            </button>
          </div>
        )}

        <Button className="w-full" disabled={!groupId || isPending} onClick={inviteUrl ? handleCopy : handleSubmit}>
          {inviteUrl ? (copied ? t('settings.member.inviteCopied') : t('settings.member.inviteCopy')) : t('settings.member.inviteGenerate')}
        </Button>
      </div>
    </div>
  )
}
