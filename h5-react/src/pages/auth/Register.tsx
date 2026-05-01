import { Avatar, Button, Form, Input, Toast } from 'antd-mobile';
import { useEffect, useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import SmsCodeButton from '@/components/SmsCodeButton';
import PageHeader from '@/components/layout/PageHeader';
import { authApi } from '@/api/auth';
import { InvitePreview } from '@/api';
import { useAuth } from '@/stores/auth';

export default function Register() {
  const [search] = useSearchParams();
  const initInvite = search.get('invite') || search.get('inviteCode') || '';
  const [phone, setPhone] = useState('');
  const [smsCode, setSmsCode] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPwd, setConfirmPwd] = useState('');
  const [inviteCode, setInviteCode] = useState(initInvite);
  const [nickname, setNickname] = useState('');
  const [preview, setPreview] = useState<InvitePreview | null>(null);
  const [loading, setLoading] = useState(false);
  const setSession = useAuth((s) => s.setSession);
  const nav = useNavigate();

  useEffect(() => {
    let cancelled = false;
    if (!inviteCode || inviteCode.length < 4) {
      setPreview(null);
      return;
    }
    authApi
      .invitePreview(inviteCode)
      .then((p) => {
        if (!cancelled) setPreview(p);
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, [inviteCode]);

  const submit = async () => {
    if (!phone || !smsCode || !password || !inviteCode) {
      Toast.show({ icon: 'fail', content: '请完整填写表单' });
      return;
    }
    if (password.length < 6) {
      Toast.show({ icon: 'fail', content: '密码至少 6 位' });
      return;
    }
    if (password !== confirmPwd) {
      Toast.show({ icon: 'fail', content: '两次输入的密码不一致' });
      return;
    }
    if (!preview?.found) {
      Toast.show({ icon: 'fail', content: '邀请码无效，请确认' });
      return;
    }
    setLoading(true);
    try {
      const r = await authApi.register({ phone, smsCode, password, inviteCode, nickname });
      setSession(r.token, {
        userID: r.memberId,
        phone: r.phone,
        nickname: r.nickname,
        avatar: r.avatar,
        inviteCode: r.inviteCode,
        isQualified: r.isQualified,
      });
      Toast.show({ icon: 'success', content: '注册成功' });
      // 注册成功后引导签署会员协议；不签也可继续，去首页
      nav('/sign-contract?type=register', { replace: true });
    } catch {
      // 全局 toast
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-page min-h-screen bg-white">
      <PageHeader title="会员注册" />
      <div className="px-5 pt-2 pb-10">
        {preview?.found && (
          <div className="card flex items-center gap-3 mb-4">
            <Avatar src={preview.avatar || ''} style={{ '--size': '40px' } as any} />
            <div>
              <div className="text-xs text-gray-500">邀请人</div>
              <div className="text-base font-medium">{preview.nickname}</div>
            </div>
            <div className="ml-auto text-xs text-primary">已绑定</div>
          </div>
        )}

        <Form layout="vertical">
          <Form.Item label="手机号">
            <Input placeholder="请输入手机号" value={phone} onChange={setPhone} clearable />
          </Form.Item>
          <Form.Item label="短信验证码" extra={<SmsCodeButton phone={phone} scene="register" />}>
            <Input placeholder="请输入验证码" value={smsCode} onChange={setSmsCode} clearable />
          </Form.Item>
          <Form.Item label="登录密码">
            <Input
              placeholder="6-32 位"
              type="password"
              value={password}
              onChange={setPassword}
              clearable
            />
          </Form.Item>
          <Form.Item label="确认密码">
            <Input
              placeholder="再次输入密码"
              type="password"
              value={confirmPwd}
              onChange={setConfirmPwd}
              clearable
            />
          </Form.Item>
          <Form.Item
            label="邀请码"
            extra={
              preview && !preview.found ? (
                <span className="text-xs text-red-500">无效</span>
              ) : null
            }
          >
            <Input placeholder="必填" value={inviteCode} onChange={setInviteCode} clearable />
          </Form.Item>
          <Form.Item label="昵称（选填）">
            <Input placeholder="不填将使用脱敏手机号" value={nickname} onChange={setNickname} clearable />
          </Form.Item>
          <Button
            block
            color="primary"
            size="large"
            loading={loading}
            onClick={submit}
            style={{ marginTop: 12, borderRadius: 12 }}
          >
            立即注册
          </Button>
        </Form>
        <div className="text-center text-sm text-gray-500 mt-4">
          已有账号？
          <Link to="/login" className="text-primary">
            立即登录
          </Link>
        </div>
      </div>
    </div>
  );
}
