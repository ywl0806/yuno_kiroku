import { FC } from 'react'

const LivePhotoIcon = () => (
  <svg className="size-4" viewBox="0 0 24 24" fill="currentColor">
    <circle cx="12" cy="12" r="3.5" />
    <path
      d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8z"
      opacity="0.35"
    />
  </svg>
)

export const LivePhotoBadge: FC = () => {
  return (
    <span className="z-20 flex items-center justify-center gap-1 rounded-full border-[1px] bg-white/90 p-2 text-[0.8rem]">
      <LivePhotoIcon />
    </span>
  )
}
