import { IdentityFaceOption } from '@/types'
import { useTranslation } from 'react-i18next'

interface Props {
  options: IdentityFaceOption[]
  selectedIdentityId: number | null
  onSelect: (identityId: number | null) => void
}

export const IdentityFaceSelector = ({ options, selectedIdentityId, onSelect }: Props) => {
  const { t } = useTranslation()

  return (
    <div className="grid grid-cols-5 gap-2">
      <button
        type="button"
        onClick={() => onSelect(null)}
        className={`flex aspect-square items-center justify-center rounded-md border-2 text-xs text-muted-foreground ${
          selectedIdentityId === null ? 'border-primary' : 'border-transparent bg-accent'
        }`}
      >
        {t('common.none')}
      </button>
      {options.map((option) => (
        <button
          key={option.identity_id}
          type="button"
          onClick={() => onSelect(option.identity_id)}
          className={`aspect-square overflow-hidden rounded-md border-2 ${
            selectedIdentityId === option.identity_id ? 'border-primary' : 'border-transparent'
          }`}
        >
          <img src={option.image_url} alt="" className="h-full w-full object-cover" />
        </button>
      ))}
    </div>
  )
}
