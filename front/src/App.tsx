import { router } from './router'
import { RouterProvider } from 'react-router-dom'

const App = () => {
  return (
    <>
      <div className="h-screen w-screen bg-[url('/uploads/window.jpg')] bg-cover bg-center">
        <RouterProvider router={router} />
      </div>
    </>
  )
}

export default App
