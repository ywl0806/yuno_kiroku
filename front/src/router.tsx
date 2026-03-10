import { DefaultLayout } from '@/components/layouts/default-layout'
import { HomePage } from '@/page/home'
import { LoginCallbackPage } from '@/page/auth/login-callback'
import { LoginPage } from '@/page/auth/login'
import { LogoutPage } from '@/page/auth/logout'
import { NotFoundPage } from '@/page/not-found-page'
import { SettingsPage } from '@/page/setting'
import { SettingsAppPage } from '@/page/setting/app'
import { SettingsAccountPage } from '@/page/setting/account'
import { SettingsFamilyNewPage } from '@/page/setting/family/new'
import { SettingsFamilyEditPage } from '@/page/setting/family/edit'
import { SettingsFamilyInvitePage } from '@/page/setting/family/invite'
import { SettingsAlbumNewPage } from '@/page/setting/album/new'
import { SettingsAlbumEditPage } from '@/page/setting/album/edit'
import { SettingsMemberInvitePage } from '@/page/setting/member/invite'
import { UploadPage } from '@/page/upload/upload'
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
        <Route path="/logout" element={<LogoutPage />} />
        <Route path="*" element={<NotFoundPage />} />

        <Route path="/settings">
          <Route index element={<SettingsPage />} />
          <Route path="app" element={<SettingsAppPage />} />
          <Route path="account" element={<SettingsAccountPage />} />

          <Route path="family">
            <Route path="new" element={<SettingsFamilyNewPage />} />
            <Route path=":familyId/edit" element={<SettingsFamilyEditPage />} />
            <Route path=":familyId/invite" element={<SettingsFamilyInvitePage />} />
          </Route>

          <Route path="album">
            <Route path="new" element={<SettingsAlbumNewPage />} />
            <Route path=":albumId/edit" element={<SettingsAlbumEditPage />} />
          </Route>

          <Route path="member">
            <Route path="invite" element={<SettingsMemberInvitePage />} />
          </Route>
        </Route>

      </Route>
    </Route>,
  ),
)
