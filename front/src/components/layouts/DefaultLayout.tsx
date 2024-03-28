import { Header } from './partials/Header'
import { Outlet } from 'react-router-dom'

export const DefaultLayout = () => {
  return (
    <div>
      <Header />

      <div className="container mx-auto">
        <Outlet />
      </div>
    </div>
  )
}
