import { TabBar } from 'antd-mobile';
import {
  AppOutline,
  ShopbagOutline,
  UnorderedListOutline,
  UserOutline,
} from 'antd-mobile-icons';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

const tabs = [
  { key: '/', title: '首页', icon: <AppOutline /> },
  { key: '/mall', title: '商城', icon: <ShopbagOutline /> },
  { key: '/warehouse', title: '仓库', icon: <UnorderedListOutline /> },
  { key: '/me', title: '我的', icon: <UserOutline /> },
];

export default function TabLayout() {
  const nav = useNavigate();
  const loc = useLocation();
  const active = tabs.find((t) => t.key === loc.pathname)?.key || '/';

  return (
    <div className="app-page tab-content">
      <Outlet />
      <div
        className="fixed bottom-0 left-0 right-0 bg-white"
        style={{
          borderTop: '1px solid #f0f0f0',
          paddingBottom: 'env(safe-area-inset-bottom)',
          zIndex: 100,
        }}
      >
        <TabBar activeKey={active} onChange={(key) => nav(key)}>
          {tabs.map((item) => (
            <TabBar.Item key={item.key} icon={item.icon} title={item.title} />
          ))}
        </TabBar>
      </div>
    </div>
  );
}
