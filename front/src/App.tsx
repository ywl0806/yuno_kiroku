import { router } from './router'
import { RouterProvider } from 'react-router-dom'

const App = () => {
  return (
    <>
      <div className="mx-auto w-fit rounded-lg bg-red-500 p-[2rem] text-[5rem]">
        <h1>hogehoge</h1>
        <RouterProvider router={router} />
      </div>
    </>
  )
}

export default App
