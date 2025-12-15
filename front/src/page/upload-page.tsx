import { UploadContainer } from '@/feature/upload/components/upload-container'
import { useSearchParams } from 'react-router-dom'

export const UploadPage = () => {
  const [searchParams] = useSearchParams()
  const albumId = searchParams.get('album_id')
  return (
    <div className="flex h-screen flex-col items-center justify-center">
      <UploadContainer albumId={albumId ? Number(albumId) : null} />
    </div>
  )
}
