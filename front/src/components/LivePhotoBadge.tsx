import MotionPhotosAutoIcon from '@mui/icons-material/MotionPhotosAuto'
import { FC } from 'react'

export const LivePhotoBadge: FC = () => {
  return (
    <span className="z-20 flex items-center justify-center gap-1 rounded-full border-[1px] bg-white/90 p-2 text-[0.8rem]">
      <MotionPhotosAutoIcon fontSize="small" />
    </span>
  )
}
