import { Button, Form, Input, Toast } from 'antd-mobile';
import { useEffect, useState } from 'react';
import PageHeader from '@/components/layout/PageHeader';
import { meApi } from '@/api/me';
import { MeProfile } from '@/api/types';
import { useAuth } from '@/stores/auth';

export default function Profile() {
  const [profile, setProfile] = useState<MeProfile | null>(null);
  const [nickname, setNickname] = useState('');
  const [realName, setRealName] = useState('');
  const [loading, setLoading] = useState(false);
  const setUser = useAuth((s) => s.setUser);

  useEffect(() => {
    meApi.profile().then((p) => {
      setProfile(p);
      setNickname(p.nickname);
      setRealName(p.realName);
    });
  }, []);

  const submit = async () => {
    setLoading(true);
    try {
      await meApi.update({ nickname, realName });
      Toast.show({ icon: 'success', content: '已保存' });
      setUser({
        userID: profile?.memberId || '',
        phone: profile?.phone,
        nickname,
        avatar: profile?.avatar,
        inviteCode: profile?.inviteCode,
        levelName: profile?.levelName,
        isQualified: profile?.isQualified,
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-page bg-white min-h-screen">
      <PageHeader title="个人资料" />
      <div className="px-5 pt-4 pb-10">
        <Form layout="vertical">
          <Form.Item label="昵称">
            <Input value={nickname} onChange={setNickname} clearable />
          </Form.Item>
          <Form.Item label="真实姓名">
            <Input value={realName} onChange={setRealName} clearable />
          </Form.Item>
          <Button block color="primary" loading={loading} onClick={submit} style={{ marginTop: 16, borderRadius: 12 }}>
            保存
          </Button>
        </Form>
      </div>
    </div>
  );
}
