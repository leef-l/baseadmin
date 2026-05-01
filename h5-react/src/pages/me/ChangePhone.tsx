import { Button, Form, Input, Toast } from 'antd-mobile';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import SmsCodeButton from '@/components/SmsCodeButton';
import PageHeader from '@/components/layout/PageHeader';
import { meApi } from '@/api/me';

export default function ChangePhone() {
  const [newPhone, setNewPhone] = useState('');
  const [smsCode, setSmsCode] = useState('');
  const [loading, setLoading] = useState(false);
  const nav = useNavigate();

  const submit = async () => {
    if (!newPhone || !smsCode) {
      Toast.show({ icon: 'fail', content: '请填写完整' });
      return;
    }
    setLoading(true);
    try {
      await meApi.changePhone(newPhone, smsCode);
      Toast.show({ icon: 'success', content: '已修改' });
      nav(-1);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-page bg-white min-h-screen">
      <PageHeader title="修改手机号" />
      <div className="px-5 pt-4">
        <Form layout="vertical">
          <Form.Item label="新手机号">
            <Input value={newPhone} onChange={setNewPhone} clearable placeholder="11 位手机号" />
          </Form.Item>
          <Form.Item label="验证码" extra={<SmsCodeButton phone={newPhone} scene="change_phone" />}>
            <Input value={smsCode} onChange={setSmsCode} clearable />
          </Form.Item>
          <Button block color="primary" loading={loading} onClick={submit} style={{ marginTop: 16, borderRadius: 12 }}>
            确认修改
          </Button>
        </Form>
      </div>
    </div>
  );
}
