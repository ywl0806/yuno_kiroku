import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useCreateGroup } from '@/feature/settings/hooks/use-create-group'
import { ChevronLeft } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'

export const SettingsFamilyNewPage = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { mutate: createGroup, isPending } = useCreateGroup()

  const [name, setName] = useState('')

  const handleSubmit = () => {
    if (!name.trim()) return
    createGroup({ name }, { onSuccess: () => navigate(-1) })
  }

  return (
    <div className="mx-auto h-full max-w-[50rem] pt-4">
      <div className="flex items-center gap-2 px-4 py-2">
        <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
          <ChevronLeft className="size-5" />
        </button>
        <span className="text-sm font-medium">{t('settings.group.newTitle')}</span>
      </div>

      <div className="space-y-5 px-4 pt-6">
        <div className="space-y-1.5">
          <Label htmlFor="group-name">{t('settings.group.nameLabel')}</Label>
          <Input
            id="group-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder={t('settings.group.namePlaceholder')}
          />
        </div>

        <Button className="w-full" disabled={isPending || !name.trim()} onClick={handleSubmit}>
          {t('common.save')}
        </Button>
      </div>
    </div>
  )
}
