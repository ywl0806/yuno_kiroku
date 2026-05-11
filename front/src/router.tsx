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
import { SettingsMemberEditPage } from '@/page/setting/member/edit'
import { SettingsMemberInvitePage } from '@/page/setting/member/invite'
import { SettingsKidNewPage } from '@/page/setting/kid/new'
import { SettingsKidEditPage } from '@/page/setting/kid/edit'
import { RecentPage } from '@/page/recent'
import { SearchPage } from '@/page/search'
import { UploadPage } from '@/page/upload/upload'
import { useGetMe } from '@/feature/settings/hooks/use-get-me'
import { Route, Navigate, createBrowserRouter, createRoutesFromElements } from 'react-router-dom'
import { BatchDetailPage } from '@/page/recent/detail'

const AdminRoute = ({ children }: { children: React.ReactNode }) => {
  const { data: me, isLoading } = useGetMe()
  if (isLoading) return null
  if (!me?.is_admin) return <Navigate to="/settings" replace />
  return <>{children}</>

}

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
        <Route path="/recent" element={<RecentPage />} />
        <Route path="/recent/:batchId" element={<BatchDetailPage />} />
        <Route path="/search" element={<SearchPage />} />
        <Route path="/upload" element={<UploadPage />} />
        <Route path="/logout" element={<LogoutPage />} />
        <Route path="*" element={<NotFoundPage />} />

        <Route path="/settings">
          <Route index element={<SettingsPage />} />
          <Route path="app" element={<SettingsAppPage />} />
          <Route path="account" element={<SettingsAccountPage />} />

          <Route path="family">
            <Route path="new" element={<AdminRoute><SettingsFamilyNewPage /></AdminRoute>} />
            <Route path=":familyId/edit" element={<AdminRoute><SettingsFamilyEditPage /></AdminRoute>} />
            <Route path=":familyId/invite" element={<AdminRoute><SettingsFamilyInvitePage /></AdminRoute>} />
          </Route>

          <Route path="album">
            <Route path="new" element={<AdminRoute><SettingsAlbumNewPage /></AdminRoute>} />
            <Route path=":albumId/edit" element={<AdminRoute><SettingsAlbumEditPage /></AdminRoute>} />
          </Route>

          <Route path="member">
            <Route path="invite" element={<AdminRoute><SettingsMemberInvitePage /></AdminRoute>} />
            <Route path=":memberId/edit" element={<AdminRoute><SettingsMemberEditPage /></AdminRoute>} />
          </Route>

          <Route path="kid">
            <Route path="new" element={<AdminRoute><SettingsKidNewPage /></AdminRoute>} />
            <Route path=":kidId/edit" element={<AdminRoute><SettingsKidEditPage /></AdminRoute>} />
          </Route>
        </Route>

      </Route>
    </Route>,
  ),
)
