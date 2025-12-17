import { UploadPhotoProvider } from './providers/upload-photo-provider'
import { queryClient } from '@/lib/query-client'
import { router } from '@/router'
import { QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from 'react-router-dom'

const App = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <UploadPhotoProvider>
        <div className="min-h-screen w-screen bg-clouds">
          <RouterProvider router={router} />
        </div>
      </UploadPhotoProvider>
    </QueryClientProvider>
  )
}

export default App
