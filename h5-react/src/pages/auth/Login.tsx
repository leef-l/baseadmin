import { Button, Form, Input, Toast } from 'antd-mobile';
import { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { authApi } from '@/api/auth';
import { useAuth } from '@/stores/auth';

export default function Login() {
  const [account, setAccount] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const setSession = useAuth((s) => s.setSession);
  const nav = useNavigate();
  const [search] = useSearchParams();
  const redirect = search.get('redirect') || '/';

  const submit = async () => {
    if (!account || !password) {
      Toast.show({ icon: 'fail', content: '请填写账号与密码' });
      return;
    }
    setLoading(true);
    try {
      const r = await authApi.login(account, password);
      setSession(r.token, {
        userID: r.memberId,
        phone: r.phone,
        nickname: r.nickname,
        avatar: r.avatar,
        inviteCode: r.inviteCode,
        isQualified: r.isQualified,
      });
      Toast.show({ icon: 'success', content: '登录成功' });
      nav(redirect, { replace: true });
    } catch {
      // 全局已 toast
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-page min-h-screen gradient-primary">
      <div className="px-6 pt-16 pb-8 text-white">
        <div className="text-3xl font-bold">欢迎回来</div>
        <div className="text-sm opacity-80 mt-2">用账号 / 手机号 / 邀请码登录</div>
      </div>
      <div className="bg-white rounded-t-3xl px-5 pt-6 pb-10 min-h-[60vh]">
        <Form layout="vertical">
          <Form.Item label="账号">
            <Input
              placeholder="手机号 / 邀请码 / 用户名"
              value={account}
              onChange={setAccount}
              clearable
            />
          </Form.Item>
          <Form.Item label="密码">
            <Input
              placeholder="请输入密码"
              type="password"
              value={password}
              onChange={setPassword}
              clearable
            />
          </Form.Item>
          <Button
            block
            color="primary"
            size="large"
            loading={loading}
            onClick={submit}
            style={{ marginTop: 24, borderRadius: 12 }}
          >
            登录
          </Button>
        </Form>
        <div className="flex justify-between mt-5 text-sm">
          <Link to="/forget-password" className="text-gray-500">
            忘记密码？
          </Link>
          <Link to="/register" className="text-primary font-medium">
            没有账号？立即注册
          </Link>
        </div>
      </div>
    </div>
  );
}
