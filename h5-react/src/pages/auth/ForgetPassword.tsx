import { Button, Form, Input, Toast } from 'antd-mobile';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import SmsCodeButton from '@/components/SmsCodeButton';
import PageHeader from '@/components/layout/PageHeader';
import { authApi } from '@/api/auth';

export default function ForgetPassword() {
  const [phone, setPhone] = useState('');
  const [smsCode, setSmsCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const nav = useNavigate();

  const submit = async () => {
    if (!phone || !smsCode || !newPassword) {
      Toast.show({ icon: 'fail', content: '请完整填写表单' });
      return;
    }
    if (newPassword.length < 6) {
      Toast.show({ icon: 'fail', content: '密码至少 6 位' });
      return;
    }
    setLoading(true);
    try {
      await authApi.forgetPassword({ phone, smsCode, newPassword });
      Toast.show({ icon: 'success', content: '密码已重置，请重新登录' });
      nav('/login', { replace: true });
    } catch {
      // 全局 toast
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-page min-h-screen bg-white">
      <PageHeader title="找回密码" />
      <div className="px-5 pt-4 pb-10">
        <Form layout="vertical">
          <Form.Item label="注册手机号">
            <Input placeholder="请输入手机号" value={phone} onChange={setPhone} clearable />
          </Form.Item>
          <Form.Item
            label="短信验证码"
            extra={<SmsCodeButton phone={phone} scene="forget_password" />}
          >
            <Input placeholder="请输入验证码" value={smsCode} onChange={setSmsCode} clearable />
          </Form.Item>
          <Form.Item label="新密码">
            <Input
              placeholder="6-32 位新密码"
              type="password"
              value={newPassword}
              onChange={setNewPassword}
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
            确认重置
          </Button>
        </Form>
      </div>
    </div>
  );
}
