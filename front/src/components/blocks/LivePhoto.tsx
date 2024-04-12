import { LivePhotoBadge } from '../LivePhotoBadge'
import { Photo } from '@/types/photo'
import { Fade, IconButton } from '@mui/material'
import { FC, useEffect, useRef, useState } from 'react'

type Props = {
  photo: Photo
}

export const LivePhoto: FC<Props> = ({ photo }) => {
  const [livePhotoPlay, setLivePhotoPlay] = useState<boolean>(false)
  const videoRef = useRef<HTMLVideoElement>(null)

  const playVideo = () => {
    setLivePhotoPlay(true)
    if (videoRef.current) {
      videoRef.current.play().then(() => {
        console.log('hoge')
      })
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
        <IconButton onClick={playVideo}>
          <LivePhotoBadge />
        </IconButton>
      </div>
      <Fade
        in={livePhotoPlay}
        timeout={{
          enter: 0,
          exit: 500,
        }}
        className="absolute"
      >
        <video
          ref={videoRef}
          className="w-full"
          src={photo.live_url}
          playsInline
          onClick={() => setLivePhotoPlay((prev) => !prev)}
        />
      </Fade>
      <Fade
        in={!livePhotoPlay}
        timeout={{
          enter: 0,
          exit: 500,
        }}
      >
        <img src={photo.thumbnail_url} alt={photo.file_name} />
      </Fade>
    </div>
  )
}
