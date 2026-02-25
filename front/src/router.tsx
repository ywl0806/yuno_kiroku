import { DefaultLayout } from '@/components/layouts/default-layout'
import { HomePage } from '@/page/home-page'
import { LoginCallbackPage } from '@/page/login-callback-page'
import { LoginPage } from '@/page/login-page'
import { LogoutPage } from '@/page/logout-page'
import { NotFoundPage } from '@/page/not-found-page'
import { SettingsPage } from '@/page/settings-page'
import { UploadPage } from '@/page/upload-page'
import { Route, createBrowserRouter, createRoutesFromElements } from 'react-router-dom'

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route>
      <Route>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/login/callback" element={<LoginCallbackPage />} />
      </Route>
      <Route path="/" element={<DefaultLayout />}>
        <Route index element={<HomePage />} />
        <Route path={'/:date'} element={<HomePage />} />
        <Route path="/upload" element={<UploadPage />} />
        <Route path="/settings" element={<SettingsPage />} />
        <Route path="/logout" element={<LogoutPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Route>,
  ),
)
