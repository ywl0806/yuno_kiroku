import { ChevronLeft } from 'lucide-react'
import { ReactNode } from 'react'
import { useNavigate } from 'react-router-dom'

interface Props {
  title: string
  children: ReactNode
}

export const SettingsSubPageLayout = ({ title, children }: Props) => {
  const navigate = useNavigate()
  return (
    <div className="h-full bg-stone-50">
      <div className="mx-auto h-full max-w-[50rem] overflow-y-auto pb-10">
        <div className="sticky top-0 z-10 flex items-center gap-1 bg-stone-50/80 px-3 py-2 backdrop-blur-sm">
          <button
            onClick={() => navigate(-1)}
            className="flex size-9 items-center justify-center rounded-full transition-colors hover:bg-stone-100 active:bg-stone-200"
          >
            <ChevronLeft className="size-5 text-stone-500" />
          </button>
          <span className="text-base font-semibold text-stone-900">{title}</span>
        </div>
        {children}
      </div>
    </div>
  )
}
