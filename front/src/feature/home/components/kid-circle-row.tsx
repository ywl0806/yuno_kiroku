import { FC } from 'react'
import { useGetKidFaceImgs } from '../hooks/use-get-kid-face-imgs'
import { useNavigate } from 'react-router-dom'
import { getAgeLabel } from '../utils/age-label'
import { useTranslation } from 'react-i18next'

type Props = {
  year: number
  month: number
}

export const KidCircleRow: FC<Props> = ({ year, month }) => {

  const { data: kidFaceImgs } = useGetKidFaceImgs(year, month)
  const nav = useNavigate()
  const { t } = useTranslation()
  if (!kidFaceImgs) return null
  return (
    <div className="flex gap-4 overflow-x-auto px-4 pb-1">
      {kidFaceImgs.map((kid) => (
        <button
          key={kid.kid_id}
          type="button"
          onClick={() => {
            nav(`/${kid.taken_at_year}-${kid.taken_at_month}`, {
              state: { openMediaItem: { mediaItemId: kid.media_item_id, year: kid.taken_at_year, month: kid.taken_at_month } },
            })
          }}
          className="flex shrink-0 flex-col items-center gap-1"
        >
          <div className="h-12 w-12 overflow-hidden rounded-full border-2 border-white/80 shadow-md">
            <img src={kid.face_img_url} alt={kid.name} className="h-full w-full object-cover" />
          </div>
          <span className="text-xs font-semibold text-white drop-shadow">{kid.name}</span>
          <span className="text-xs text-white/80 drop-shadow">{getAgeLabel(new Date(year, month), kid.birth_date, t)}</span>
        </button>
      ))}
    </div>
  )
}
