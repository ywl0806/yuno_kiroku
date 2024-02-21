import { HomePage } from './page/Home'
import { Route, createBrowserRouter, createRoutesFromElements } from 'react-router-dom'

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<HomePage />}>
      <Route path="/login" element={<div>Login</div>} />
      <Route path="/register" element={<div>Register</div>} />
      <Route path="/dashboard" element={<div>Dashboard</div>} />
      <Route path="/profile" element={<div>Profile</div>} />
      <Route path="/settings" element={<div>Settings</div>} />
      <Route path="/logout" element={<div>Logout</div>} />
      <Route path="*" element={<div>404</div>} />
    </Route>,
  ),
)
