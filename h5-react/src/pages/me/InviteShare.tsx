import { Button, Card, Toast } from 'antd-mobile';
import { QRCodeCanvas } from 'qrcode.react';
import { useEffect, useState } from 'react';
import PageHeader from '@/components/layout/PageHeader';
import { meApi } from '@/api/me';
import { MeProfile } from '@/api/types';

export default function InviteShare() {
  const [profile, setProfile] = useState<MeProfile | null>(null);

  useEffect(() => {
    meApi.profile().then(setProfile);
  }, []);

  const inviteUrl =
    profile?.inviteUrl ||
    (profile?.inviteCode ? `${window.location.origin}/register?invite=${profile.inviteCode}` : '');

  const copy = async (text: string) => {
    if (!text) return;
    try {
      await navigator.clipboard.writeText(text);
      Toast.show({ icon: 'success', content: '已复制' });
    } catch {
      const ta = document.createElement('textarea');
      ta.value = text;
      document.body.appendChild(ta);
      ta.select();
      document.execCommand('copy');
      document.body.removeChild(ta);
      Toast.show({ icon: 'success', content: '已复制' });
    }
  };

  return (
    <div className="app-page bg-[#f5f5f7] min-h-screen">
      <PageHeader title="邀请好友" />
      <div className="px-3 pt-4 space-y-3">
        <Card className="rounded-xl text-center">
          <div className="text-sm text-gray-500 mb-2">我的邀请码</div>
          <div className="text-3xl font-bold text-primary tracking-widest font-mono">
            {profile?.inviteCode || '-'}
          </div>
          <Button
            size="small"
            color="primary"
            fill="outline"
            className="mt-3"
            onClick={() => copy(profile?.inviteCode || '')}
          >
            复制邀请码
          </Button>
        </Card>

        <Card className="rounded-xl text-center">
          <div className="text-sm text-gray-500 mb-2">邀请二维码</div>
          {inviteUrl ? (
            <div className="flex justify-center my-3">
              <QRCodeCanvas value={inviteUrl} size={200} />
            </div>
          ) : (
            <div className="text-gray-400 my-3">加载中…</div>
          )}
          <div className="text-xs text-gray-500 break-all px-4">{inviteUrl}</div>
          <Button color="primary" className="mt-3" onClick={() => copy(inviteUrl)}>
            复制邀请链接
          </Button>
        </Card>

        <Card className="rounded-xl">
          <div className="text-sm font-medium mb-2">邀请奖励</div>
          <div className="text-xs text-gray-600 leading-relaxed">
            1. 好友通过你的邀请码注册成为团队成员；
            <br />
            2. 好友购买商城商品 / 仓库交易，按等级返推广奖；
            <br />
            3. 团队达成业绩自动晋升等级，奖金更多。
          </div>
        </Card>
      </div>
    </div>
  );
}
