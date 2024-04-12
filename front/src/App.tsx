import { queryClient } from './lib/queryClient'
import { router } from './router'
import { QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from 'react-router-dom'

const App = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="min-h-screen w-screen bg-clouds">
        <RouterProvider router={router} />
      </div>
    </QueryClientProvider>
  )
}

export default App
