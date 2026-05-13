import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useGetTags } from '@/feature/tag/hooks/use-get-tags'
import { useTagMutations } from '@/feature/tag/hooks/use-tag-mutations'
import { Tag } from '@/types'
import { Check, Plus, Tag as TagIcon } from 'lucide-react'
import { FC, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { TagBadge } from './tag-badge'
import { useGetMediaItemTags } from '../hooks/use-get-media-item-tags'

type Props = {
  mediaItemId: string
}

export const TagSelector: FC<Props> = ({ mediaItemId }) => {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  const [popoverOpen, setPopoverOpen] = useState(false)
  const [loaded, setLoaded] = useState(false)
  const [showCreate, setShowCreate] = useState(false)
  const [newName, setNewName] = useState('')
  const queryClient = useQueryClient()

  const { data: allTags = [] } = useGetTags()
  const { addTag, removeTag, createTag, isPending } = useTagMutations(mediaItemId)
  const { data: mediaItemTags = [] } = useGetMediaItemTags(mediaItemId, open)

  const handleOpen = () => {
    if (!open && !loaded) setLoaded(true)
    if (open) {
      setOpen(false)
      setPopoverOpen(false)
    } else {
      setOpen(true)
    }
  }

  const handleToggleTag = (tag: Tag) => {
    if (mediaItemTags.some((t) => t.id === tag.id)) {
      removeTag(tag.id, {
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['media-item-tags', mediaItemId] }),
      })
    } else {
      addTag(tag.id, {
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['media-item-tags', mediaItemId] }),
      })
    }
  }

  const handleCreateTag = () => {
    if (!newName.trim()) return
    createTag({ name: newName.trim() }, {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: ['media-item-tags', mediaItemId] })
        queryClient.invalidateQueries({ queryKey: ['tags'] })
      },
    })
    setNewName('')
    setShowCreate(false)
  }

  return (
    <div className="relative h-full w-full">
      <div className="flex items-end gap-1">
        <button
          onClick={handleOpen}
          className="size-10 flex items-center gap-1 rounded-full bg-black/50 p-3 text-xs text-white backdrop-blur-sm hover:bg-black/70"
        >
          <TagIcon className="size-5" />
        </button>

        {open && (
          <div className="relative flex h-full flex-wrap items-center gap-1 animate-in fade-in duration-100">
            {mediaItemTags.map((tag) => (
              <TagBadge key={tag.id} tag={tag} onRemove={() => handleToggleTag(tag)} />
            ))}

            <Popover open={popoverOpen} onOpenChange={setPopoverOpen}>
              <PopoverTrigger asChild>
                <button className="flex items-center gap-1 rounded-lg bg-black/20 px-2 py-1 text-xs text-white backdrop-blur-sm hover:bg-black/70">
                  <Plus className="size-4" />
                  {t('tag.add')}
                </button>
              </PopoverTrigger>
              <PopoverContent className="w-56 p-2" side="top" align="start" sideOffset={8}>
                {/* 헤더 */}
                <div className="mb-1 flex items-center justify-between">
                  <span className="text-xs font-semibold text-muted-foreground">{t('tag.title')}</span>
                </div>

                {/* 태그 목록 */}
                <div className="max-h-48 space-y-0.5 overflow-y-auto">
                  {!loaded && (
                    <p className="px-2 py-1.5 text-xs text-muted-foreground">{t('common.loading')}</p>
                  )}
                  {loaded && allTags.map((tag) => {
                    const isSelected = mediaItemTags.some((t) => t.id === tag.id)
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
          </div>
        )}
      </div>
    </div>
  )
}
