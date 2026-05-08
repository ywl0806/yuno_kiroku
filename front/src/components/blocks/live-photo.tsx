import { LivePhotoBadge } from '@/components/blocks/live-photo-badge'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { MediaItem } from '@/types'
import { FC, useEffect, useRef, useState } from 'react'

type Props = {
  photo: MediaItem
}

export const LivePhoto: FC<Props> = ({ photo }) => {
  const [livePhotoPlay, setLivePhotoPlay] = useState<boolean>(false)
  const videoRef = useRef<HTMLVideoElement>(null)

  const playVideo = () => {
    setLivePhotoPlay(true)
    if (videoRef.current) {
      videoRef.current.play()
    }
  }

  useEffect(() => {
    videoRef.current?.addEventListener('ended', () => {
      setLivePhotoPlay(false)
    })

    return () => {
      videoRef.current?.removeEventListener('ended', () => {
        setLivePhotoPlay(false)
        videoRef.current?.pause()
      })
    }
  }, [])

  return (
    <div className="relative">
      <div className="absolute">
        <Button variant="ghost" size="icon" className="rounded-full" onClick={playVideo}>
          <LivePhotoBadge />
        </Button>
      </div>
      <video
        ref={videoRef}
        className={cn('absolute w-full transition-opacity', livePhotoPlay ? 'opacity-100 duration-0' : 'opacity-0 duration-500')}
        src={photo.live_url}
        playsInline
        onClick={() => setLivePhotoPlay((prev) => !prev)}
      />
      <img
        src={photo.thumbnail_url}
        alt={photo.file_name}
        className={cn('transition-opacity', livePhotoPlay ? 'opacity-0 duration-500' : 'opacity-100 duration-0')}
      />
    </div>
  )
}
