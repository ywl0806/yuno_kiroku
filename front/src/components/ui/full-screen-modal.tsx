import { cn } from '@/lib/utils'
import { FC, ReactNode } from 'react'
import { createPortal } from 'react-dom'

type Props = {
  open: boolean
  children: ReactNode
  className?: string
}

export const FullScreenModal: FC<Props> = ({ open, children, className }) => {
  if (!open) return null
  return createPortal(
    <div
      className={cn('fixed inset-0 z-50 bg-background animate-in slide-in-from-bottom duration-300', className)}
      aria-modal="true"
      role="dialog"
    >
      {children}
    </div>,
    document.body,
  )
}
