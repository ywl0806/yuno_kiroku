import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'
import { ReactNode } from 'react'
import { Link } from 'react-router-dom'

export const SettingsList = ({ children }: { children: ReactNode }) => (
  <div className="bg-white pl-2">{children}</div>
)

export const SettingsListItem = ({ to, children }: { to: string; children: ReactNode }) => (
  <Link to={to}>
    <div className="flex items-center justify-between border-b bg-white px-4 py-3 hover:bg-accent">{children}</div>
  </Link>
)

export const SettingsListAddButton = ({ to, label }: { to: string; label: string }) => (
  <Button asChild variant="ghost" className="w-full justify-start gap-2 rounded-none bg-white px-4 py-3 text-sm text-blue-400 hover:text-foreground">
    <Link to={to}>
      <Plus className="size-4" />
      {label}
    </Link>
  </Button>
)
