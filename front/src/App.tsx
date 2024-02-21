import { router } from './router'
import { RouterProvider } from 'react-router-dom'

const App = () => {
  return (
    <>
      <div className="h-screen w-screen bg-slate-200">
        <RouterProvider router={router} />
      </div>
    </>
  )
}

export default App
