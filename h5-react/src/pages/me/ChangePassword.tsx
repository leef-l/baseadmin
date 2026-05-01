import { Button, Form, Input, Toast } from 'antd-mobile';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import PageHeader from '@/components/layout/PageHeader';
import { meApi } from '@/api/me';

export default function ChangePassword() {
  const [oldPwd, setOldPwd] = useState('');
  const [newPwd, setNewPwd] = useState('');
  const [confirm, setConfirm] = useState('');
  const [loading, setLoading] = useState(false);
  const nav = useNavigate();

  const submit = async () => {
    if (!oldPwd || !newPwd) {
      Toast.show({ icon: 'fail', content: '请填写完整' });
      return;
    }
    if (newPwd.length < 6) {
      Toast.show({ icon: 'fail', content: '新密码至少 6 位' });
      return;
    }
    if (newPwd !== confirm) {
      Toast.show({ icon: 'fail', content: '两次输入不一致' });
      return;
    }
    setLoading(true);
    try {
      await meApi.changePassword(oldPwd, newPwd);
      Toast.show({ icon: 'success', content: '已修改' });
      nav(-1);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-page bg-white min-h-screen">
      <PageHeader title="修改密码" />
      <div className="px-5 pt-4">
        <Form layout="vertical">
          <Form.Item label="旧密码">
            <Input type="password" value={oldPwd} onChange={setOldPwd} clearable />
          </Form.Item>
          <Form.Item label="新密码">
            <Input type="password" value={newPwd} onChange={setNewPwd} clearable />
          </Form.Item>
          <Form.Item label="确认新密码">
            <Input type="password" value={confirm} onChange={setConfirm} clearable />
          </Form.Item>
          <Button block color="primary" loading={loading} onClick={submit} style={{ marginTop: 16, borderRadius: 12 }}>
            确认修改
          </Button>
        </Form>
      </div>
    </div>
  );
}
