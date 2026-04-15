import { Tag } from '@/types'
import { X } from 'lucide-react'
import { FC } from 'react'

type Props = {
  tag: Tag
  onRemove?: () => void
}

export const TagBadge: FC<Props> = ({ tag, onRemove }) => {
  return (
    <span className="inline-flex items-center gap-1 rounded-full bg-black/50 px-2 py-1 text-xs text-white backdrop-blur-sm">
      {tag.name}
      {onRemove && (
        <button
          onClick={(e) => {
            e.stopPropagation()
            onRemove()
          }}
          className="ml-0.5 rounded-full hover:bg-white/30"
        >
          <X className="size-3" />
        </button>
      )}
    </span>
  )
}
