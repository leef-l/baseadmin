import { Suspense, lazy } from 'react';
import { Navigate, Route, Routes, useLocation } from 'react-router-dom';
import { DotLoading } from 'antd-mobile';
import { useAuth } from '@/stores/auth';
import TabLayout from '@/components/layout/TabLayout';

const Login = lazy(() => import('@/pages/auth/Login'));
const Register = lazy(() => import('@/pages/auth/Register'));
const ForgetPassword = lazy(() => import('@/pages/auth/ForgetPassword'));

const Home = lazy(() => import('@/pages/home/Home'));
const Mall = lazy(() => import('@/pages/mall/MallList'));
const MallDetail = lazy(() => import('@/pages/mall/MallDetail'));
const MallOrders = lazy(() => import('@/pages/mall/MallOrders'));

const Warehouse = lazy(() => import('@/pages/warehouse/Warehouse'));

const Me = lazy(() => import('@/pages/me/Me'));
const Profile = lazy(() => import('@/pages/me/Profile'));
const ChangePassword = lazy(() => import('@/pages/me/ChangePassword'));
const ChangePhone = lazy(() => import('@/pages/me/ChangePhone'));
const InviteShare = lazy(() => import('@/pages/me/InviteShare'));
const Wallet = lazy(() => import('@/pages/wallet/Wallet'));
const Team = lazy(() => import('@/pages/team/Team'));

function RequireAuth({ children }: { children: JSX.Element }) {
  const isAuthed = useAuth((s) => s.isAuthed());
  const loc = useLocation();
  if (!isAuthed) {
    const redirect = encodeURIComponent(loc.pathname + loc.search);
    return <Navigate to={`/login?redirect=${redirect}`} replace />;
  }
  return children;
}

function PageFallback() {
  return (
    <div className="flex items-center justify-center h-[80vh]">
      <DotLoading color="primary" />
    </div>
  );
}

export default function AppRouter() {
  return (
    <Suspense fallback={<PageFallback />}>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/forget-password" element={<ForgetPassword />} />

        <Route
          path="/"
          element={
            <RequireAuth>
              <TabLayout />
            </RequireAuth>
          }
        >
          <Route index element={<Home />} />
          <Route path="mall" element={<Mall />} />
          <Route path="warehouse" element={<Warehouse />} />
          <Route path="me" element={<Me />} />
        </Route>

        <Route
          path="/mall/detail/:id"
          element={
            <RequireAuth>
              <MallDetail />
            </RequireAuth>
          }
        />
        <Route
          path="/mall/orders"
          element={
            <RequireAuth>
              <MallOrders />
            </RequireAuth>
          }
        />
        <Route
          path="/me/profile"
          element={
            <RequireAuth>
              <Profile />
            </RequireAuth>
          }
        />
        <Route
          path="/me/change-password"
          element={
            <RequireAuth>
              <ChangePassword />
            </RequireAuth>
          }
        />
        <Route
          path="/me/change-phone"
          element={
            <RequireAuth>
              <ChangePhone />
            </RequireAuth>
          }
        />
        <Route
          path="/me/invite"
          element={
            <RequireAuth>
              <InviteShare />
            </RequireAuth>
          }
        />
        <Route
          path="/wallet"
          element={
            <RequireAuth>
              <Wallet />
            </RequireAuth>
          }
        />
        <Route
          path="/team"
          element={
            <RequireAuth>
              <Team />
            </RequireAuth>
          }
        />

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Suspense>
  );
}
