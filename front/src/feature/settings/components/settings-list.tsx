import { ChevronRight, Plus } from 'lucide-react'
import { ReactNode } from 'react'
import { Link } from 'react-router-dom'

export const SettingsList = ({ children }: { children: ReactNode }) => (
  <div className="overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04]">{children}</div>
)

export const SettingsListItem = ({ to, children }: { to: string; children: ReactNode }) => (
  <Link to={to} className="group relative flex items-center justify-between px-4 py-3.5 transition-colors hover:bg-stone-50 active:bg-stone-100">
    <span className="text-sm font-medium text-stone-800">{children}</span>
    <ChevronRight className="size-4 flex-shrink-0 text-stone-300 transition-transform group-hover:translate-x-0.5" />
    <span className="absolute inset-x-4 bottom-0 h-px bg-stone-100 group-last:hidden" />
  </Link>
)

export const SettingsListAddButton = ({ to, label }: { to: string; label: string }) => (
  <Link
    to={to}
    className="flex w-full items-center gap-2.5 px-4 py-3.5 text-sm font-medium text-amber-600 transition-colors hover:bg-amber-50 active:bg-amber-100"
  >
    <span className="flex size-5 items-center justify-center rounded-full bg-amber-100">
      <Plus className="size-3 text-amber-600" strokeWidth={2.5} />
    </span>
    {label}
  </Link>
)
