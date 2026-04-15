import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { getTagsForMediaItem } from '@/service/tag-service'
import { useGetTags } from '@/feature/tag/hooks/use-get-tags'
import { useTagMutations } from '@/feature/tag/hooks/use-tag-mutations'
import { Tag } from '@/types'
import { Check, Plus, Tag as TagIcon, X } from 'lucide-react'
import { FC, useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { TagBadge } from './tag-badge'

type Props = {
  mediaItemId: number
}

export const TagSelector: FC<Props> = ({ mediaItemId }) => {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  const [openAddTag, setOpenAddTag] = useState(false)
  const [loaded, setLoaded] = useState(false)
  const [currentTags, setCurrentTags] = useState<Tag[]>([])
  const [showCreate, setShowCreate] = useState(false)
  const [newName, setNewName] = useState('')
  const containerRef = useRef<HTMLDivElement>(null)
  const queryClient = useQueryClient()

  const { data: allTags = [] } = useGetTags()
  const { addTag, removeTag, createTag, isPending } = useTagMutations(mediaItemId)

  const currentTagIds = new Set(currentTags.map((t) => t.id))

  // 바깥 클릭 시 닫기
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false)
        setOpenAddTag(false)
      }
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])


  // 버튼 클릭 시 최초 한 번만 현재 태그 fetch
  const handleOpen = async () => {
    if (!open && !loaded) {
      try {
        const tags = await getTagsForMediaItem(mediaItemId)
        setCurrentTags(tags)
        setLoaded(true)
      } catch {
        setLoaded(true)
      }
    }
    if (open) {
      setOpen(false)
      setOpenAddTag(false)
    } else {
      setOpen(true)
    }

  }

  const handleOpenAddTag = () => {
    setOpenAddTag(v => !v)
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
    <div ref={containerRef} className="relative h-full w-full">
      <div className="flex gap-1 items-end">

        <button
          onClick={handleOpen}
          className="size-9 flex items-center gap-1 rounded-full bg-black/50 p-3 text-xs text-white backdrop-blur-sm hover:bg-black/70"
        >
          <TagIcon className="size-4" />
        </button>
        {open && (
          <div className="relative flex items-center h-full gap-1 flex-wrap animate-in fade-in duration-100">
            {currentTags.map((tag) => (
              <TagBadge key={tag.id} tag={tag} onRemove={() => handleToggleTag(tag)} />
            ))}
            <button onClick={handleOpenAddTag} className="flex rounded-lg items-center gap-1 bg-black/20 px-2 py-1 text-xs text-white backdrop-blur-sm hover:bg-black/70">
              <Plus className="size-4" />
              {t('tag.add')}
            </button>
          </div>
        )}
      </div>
      {openAddTag && (
        <div className="absolute bottom-full left-0 mb-2 w-56 rounded-lg border bg-background p-2 shadow-lg">
          {/* 헤더 */}
          <div className="mb-1 flex items-center justify-between">
            <span className="text-xs font-semibold text-muted-foreground">{t('tag.title')}</span>
            <button onClick={() => setOpenAddTag(false)} className="rounded p-0.5 hover:bg-muted">
              <X className="size-3 text-muted-foreground" />
            </button>
          </div>

          {/* 태그 목록 */}
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
        </div>
      )}
    </div>
  )
}
