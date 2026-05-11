import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { getTagsForMediaItem } from '@/service/tag-service'
import { useGetTags } from '@/feature/tag/hooks/use-get-tags'
import { useTagMutations } from '@/feature/tag/hooks/use-tag-mutations'
import { Tag } from '@/types'
import { Check, Plus, Tag as TagIcon } from 'lucide-react'
import { FC, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'

type Props = {
  mediaItemId: string
}

export const TagSelector: FC<Props> = ({ mediaItemId }) => {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  const [loaded, setLoaded] = useState(false)
  const [currentTags, setCurrentTags] = useState<Tag[]>([])
  const [showCreate, setShowCreate] = useState(false)
  const [newName, setNewName] = useState('')
  const queryClient = useQueryClient()

  const { data: allTags = [] } = useGetTags()
  const { addTag, removeTag, createTag, isPending } = useTagMutations(mediaItemId)

  const currentTagIds = new Set(currentTags.map((t) => t.id))

  const handleOpenChange = async (next: boolean) => {
    if (next && !loaded) {
      try {
        const tags = await getTagsForMediaItem(mediaItemId)
        setCurrentTags(tags)
        setLoaded(true)
      } catch {
        setLoaded(true)
      }
    }
    if (!next) {
      setShowCreate(false)
      setNewName('')
    }
    setOpen(next)
  }

  const handleToggleTag = (tag: Tag) => {
    if (currentTagIds.has(tag.id)) {
      removeTag(tag.id, {
        onSuccess: () => setCurrentTags((prev) => prev.filter((t) => t.id !== tag.id)),
      })
    } else {
      addTag(tag.id, {
        onSuccess: () => setCurrentTags((prev) => [...prev, tag]),
      })
    }
  }

  const handleCreateTag = () => {
    if (!newName.trim()) return
    createTag({ name: newName.trim() }, {
      onSuccess: (newTag) => {
        setCurrentTags((prev) => [...prev, newTag])
        queryClient.invalidateQueries({ queryKey: ['tags'] })
      },
    })
    setNewName('')
    setShowCreate(false)
  }

  return (
    <Popover open={open} onOpenChange={handleOpenChange}>
      <PopoverTrigger asChild>
        <button className="flex size-9 items-center justify-center rounded-full bg-black/50 text-white backdrop-blur-sm hover:bg-black/70">
          <TagIcon className="size-5" />
        </button>
      </PopoverTrigger>
      <PopoverContent className="w-64 p-2" side="top" align="start">
        <p className="mb-1 text-xs font-semibold text-muted-foreground">{t('tag.title')}</p>

        {/* 현재 선택된 태그 */}
        {loaded && currentTags.length > 0 && (
          <div className="mb-1.5 flex flex-wrap gap-1 border-b pb-1.5">
            {currentTags.map((tag) => (
              <span
                key={tag.id}
                className="flex items-center gap-1 rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary"
              >
                {tag.name}
                <button
                  onClick={() => handleToggleTag(tag)}
                  disabled={isPending}
                  className="leading-none hover:text-destructive"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}

        {/* 전체 태그 목록 */}
        <div className="max-h-48 space-y-0.5 overflow-y-auto">
          {!loaded && (
            <p className="px-2 py-1.5 text-xs text-muted-foreground">{t('common.loading')}</p>
          )}
          {loaded && allTags.map((tag) => {
            const isSelected = currentTagIds.has(tag.id)
            return (
              <button
                key={tag.id}
                className="flex w-full items-center justify-between rounded px-2 py-1.5 text-sm hover:bg-muted"
                onClick={() => handleToggleTag(tag)}
                disabled={isPending}
              >
                <span>{tag.name}</span>
                {isSelected && <Check className="size-4 text-primary" />}
              </button>
            )
          })}
          {loaded && allTags.length === 0 && (
            <p className="px-2 py-1.5 text-xs text-muted-foreground">{t('tag.no_tags')}</p>
          )}
        </div>

        {/* 새 태그 생성 */}
        <div className="mt-1.5 border-t pt-1.5">
          {showCreate ? (
            <div className="space-y-1.5">
              <Input
                placeholder={t('tag.name')}
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                className="h-7 text-xs"
                autoFocus
              />
              <div className="flex gap-1">
                <Button
                  size="sm"
                  className="h-7 flex-1 text-xs"
                  onClick={handleCreateTag}
                  disabled={isPending || !newName.trim()}
                >
                  {t('common.save')}
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  className="h-7 text-xs"
                  onClick={() => { setShowCreate(false); setNewName('') }}
                >
                  {t('common.cancel')}
                </Button>
              </div>
            </div>
          ) : (
            <button
              className="flex w-full items-center gap-1 rounded px-2 py-1.5 text-xs text-muted-foreground hover:bg-muted"
              onClick={() => setShowCreate(true)}
            >
              <Plus className="size-3" />
              {t('tag.create_custom')}
            </button>
          )}
        </div>
      </PopoverContent>
    </Popover>
  )
}
