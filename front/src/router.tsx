import { DefaultLayout } from '@/components/layouts/default-layout'
import { HomePage } from '@/page/home-page'
import { LoginPage } from '@/page/login-page'
import { UploadPage } from '@/page/upload-page'
import { Route, createBrowserRouter, createRoutesFromElements } from 'react-router-dom'

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route>
      <Route>
        <Route path="/login" element={<LoginPage />} />
      </Route>
      <Route path="/" element={<DefaultLayout />}>
        <Route index element={<HomePage />} />
        <Route path={'/:date'} element={<HomePage />} />
        <Route path="/upload" element={<UploadPage />} />
        <Route path="/settings" element={<div>Settings</div>} />
        <Route path="/logout" element={<div>Logout</div>} />
        <Route path="*" element={<div>404</div>} />
      </Route>
    </Route>,
  ),
)
