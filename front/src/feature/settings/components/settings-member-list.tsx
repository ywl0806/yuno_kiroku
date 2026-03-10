import { Button } from '@/components/ui/button'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { Member } from '@/types'
import { ChevronDown, User, UserPlus } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'

interface Props {
  members?: Member[]
}

export const SettingsMemberList = ({ members }: Props) => {
  const { t } = useTranslation()
  const [open, setOpen] = useState(true)

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <CollapsibleTrigger asChild>
        <button className="flex w-full items-center justify-between px-4 py-2 hover:bg-accent">
          <span className="text-xs font-semibold text-muted-foreground">{t('settings.member.title')}</span>
          <ChevronDown className={`size-4 text-muted-foreground transition-transform duration-200 ${open ? 'rotate-180' : ''}`} />
        </button>
      </CollapsibleTrigger>
      <CollapsibleContent>
        {members?.map((member) => (
          <div key={member.id} className="flex items-center gap-3 px-4 py-3">
            <User className="size-4 shrink-0 text-muted-foreground" />
            <span className="text-sm">{member.name || member.username}</span>
          </div>
        ))}
        <Button asChild variant="ghost" className="w-full justify-start gap-2 rounded-none px-4 py-3 text-sm text-muted-foreground hover:text-foreground">
          <Link to="/settings/member/invite">
            <UserPlus className="size-4" />
            {t('settings.member.invite')}
          </Link>
        </Button>
      </CollapsibleContent>
    </Collapsible>
  )
}
