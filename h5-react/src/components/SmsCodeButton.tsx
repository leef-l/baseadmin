import { Button, Toast } from 'antd-mobile';
import { useEffect, useRef, useState } from 'react';
import { authApi, SmsScene } from '@/api/auth';

interface Props {
  phone: string;
  scene: SmsScene;
  beforeSend?: () => string | undefined;
}

export default function SmsCodeButton({ phone, scene, beforeSend }: Props) {
  const [count, setCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const timer = useRef<number>();

  useEffect(
    () => () => {
      if (timer.current) window.clearInterval(timer.current);
    },
    [],
  );

  const startCountdown = (seconds: number) => {
    setCount(seconds);
    timer.current = window.setInterval(() => {
      setCount((s) => {
        if (s <= 1) {
          window.clearInterval(timer.current);
          return 0;
        }
        return s - 1;
      });
    }, 1000);
  };

  const handleSend = async () => {
    if (count > 0) return;
    const errMsg = beforeSend?.();
    if (errMsg) {
      Toast.show({ icon: 'fail', content: errMsg });
      return;
    }
    if (!/^1[3-9]\d{9}$/.test(phone)) {
      Toast.show({ icon: 'fail', content: '请输入正确的手机号' });
      return;
    }
    setLoading(true);
    try {
      const r = await authApi.sendSmsCode(phone, scene);
      Toast.show({ icon: 'success', content: '验证码已发送' });
      startCountdown(r?.expiresIn || 60);
    } catch {
      // 全局拦截器已 toast
    } finally {
      setLoading(false);
    }
  };

  return (
    <Button
      size="small"
      color="primary"
      fill="none"
      loading={loading}
      onClick={handleSend}
      disabled={count > 0}
      style={{ minWidth: 88 }}
    >
      {count > 0 ? `${count}s 后重发` : '获取验证码'}
    </Button>
  );
}
