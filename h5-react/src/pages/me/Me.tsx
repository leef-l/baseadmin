import { Avatar, Dialog, List, Toast } from 'antd-mobile';
import { useEffect, useState } from 'react';
import {
  GiftOutline,
  KeyOutline,
  PayCircleOutline,
  PhoneFill,
  RightOutline,
  TeamOutline,
  UserContactOutline,
  UserSetOutline,
} from 'antd-mobile-icons';
import { useNavigate } from 'react-router-dom';
import { meApi } from '@/api/me';
import { MeProfile } from '@/api/types';
import { useAuth } from '@/stores/auth';

export default function Me() {
  const [profile, setProfile] = useState<MeProfile | null>(null);
  const user = useAuth((s) => s.user);
  const setUser = useAuth((s) => s.setUser);
  const clear = useAuth((s) => s.clear);
  const nav = useNavigate();

  useEffect(() => {
    meApi.profile().then((p) => {
      setProfile(p);
      setUser({
        userID: p.memberId,
        phone: p.phone,
        nickname: p.nickname,
        avatar: p.avatar,
        inviteCode: p.inviteCode,
        levelName: p.levelName,
        isQualified: p.isQualified,
      });
    });
  }, []);

  const logout = async () => {
    const ok = await Dialog.confirm({ content: '确认退出登录？' });
    if (!ok) return;
    clear();
    Toast.show({ icon: 'success', content: '已退出' });
    nav('/login', { replace: true });
  };

  return (
    <div className="bg-[#f5f5f7] min-h-screen">
      <div className="gradient-primary text-white px-5 pt-4 pb-12">
        <div className="flex items-center gap-3">
          <Avatar
            src={user?.avatar || ''}
            style={{ '--size': '64px', '--border-radius': '32px' } as any}
          />
          <div className="flex-1">
            <div className="text-lg font-bold">{user?.nickname || '会员'}</div>
            <div className="text-xs opacity-80 mt-1">手机号：{profile?.phone || '-'}</div>
            <div className="text-xs opacity-80 mt-1">
              邀请码 <span className="font-mono">{user?.inviteCode || '-'}</span>
            </div>
          </div>
        </div>
      </div>

      <div className="-mt-8 mx-3">
        <div className="card grid grid-cols-3 gap-3 text-center">
          <Stat label="团队总数" value={profile?.teamCount} />
          <Stat label="活跃用户" value={profile?.activeCount} />
          <Stat label="直推人数" value={profile?.directCount} />
        </div>
      </div>

      <div className="m-3">
        <List style={{ '--border-top': 'none', '--border-bottom': 'none', borderRadius: 12 } as any}>
          <List.Item prefix={<UserSetOutline />} arrow={<RightOutline />} onClick={() => nav('/me/profile')}>
            个人资料
          </List.Item>
          <List.Item prefix={<KeyOutline />} arrow={<RightOutline />} onClick={() => nav('/me/change-password')}>
            修改密码
          </List.Item>
          <List.Item prefix={<PhoneFill />} arrow={<RightOutline />} onClick={() => nav('/me/change-phone')}>
            修改手机号
          </List.Item>
          <List.Item prefix={<GiftOutline />} arrow={<RightOutline />} onClick={() => nav('/me/invite')}>
            邀请好友
          </List.Item>
        </List>
      </div>

      <div className="m-3">
        <List style={{ borderRadius: 12 } as any}>
          <List.Item prefix={<PayCircleOutline />} arrow={<RightOutline />} onClick={() => nav('/wallet')}>
            我的钱包
          </List.Item>
          <List.Item prefix={<TeamOutline />} arrow={<RightOutline />} onClick={() => nav('/team')}>
            我的团队
          </List.Item>
          <List.Item prefix={<UserContactOutline />} arrow={<RightOutline />} onClick={() => nav('/mall/orders')}>
            我的订单
          </List.Item>
        </List>
      </div>

      <div className="m-3 mt-6 text-center">
        <button
          onClick={logout}
          className="bg-white text-red-500 w-full py-3 rounded-xl text-sm font-medium"
          style={{ border: '1px solid #f0f0f0' }}
        >
          退出登录
        </button>
      </div>
    </div>
  );
}

function Stat({ label, value }: { label: string; value: number | undefined }) {
  return (
    <div>
      <div className="text-xl font-bold text-primary">{value ?? 0}</div>
      <div className="text-xs text-gray-500 mt-1">{label}</div>
    </div>
  );
}
