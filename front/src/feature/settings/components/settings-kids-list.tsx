import { Button } from '@/components/ui/button'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { Kid } from '@/types'
import { ChevronDown, Pencil, Plus } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'

interface Props {
  kids?: Kid[]
}

export const SettingsKidsList = ({ kids }: Props) => {
  const { t } = useTranslation()
  const [open, setOpen] = useState(true)

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <CollapsibleTrigger asChild>
        <button className="flex w-full items-center justify-between px-4 py-2 hover:bg-accent">
          <span className="text-xs font-semibold text-muted-foreground">{t('settings.kid.title')}</span>
          <ChevronDown className={`size-4 text-muted-foreground transition-transform duration-200 ${open ? 'rotate-180' : ''}`} />
        </button>
      </CollapsibleTrigger>
      <CollapsibleContent>
        {kids?.map((kid) => (
          <Link key={kid.id} to={`/settings/kid/${kid.id}/edit`}>
            <div className="flex items-center justify-between px-4 py-3 hover:bg-accent">
              <span className="text-sm">{kid.name ?? kid.birth_date ?? `Kid #${kid.id}`}</span>
              <Pencil className="size-4 text-muted-foreground" />
            </div>
          </Link>
        ))}
        <Button asChild variant="ghost" className="w-full justify-start gap-2 rounded-none px-4 py-3 text-sm text-muted-foreground hover:text-foreground">
          <Link to="/settings/kid/new">
            <Plus className="size-4" />
            {t('settings.kid.add')}
          </Link>
        </Button>
      </CollapsibleContent>
    </Collapsible>
  )
}
