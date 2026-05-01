import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { Toast } from 'antd-mobile';
import { tokenStorage, userStorage } from '@/utils/storage';

const baseURL = (import.meta.env.VITE_API_BASE as string) || '/api';

const instance: AxiosInstance = axios.create({
  baseURL,
  timeout: 15000,
});

instance.interceptors.request.use((cfg) => {
  const token = tokenStorage.get();
  if (token) {
    cfg.headers = cfg.headers || {};
    (cfg.headers as any).Authorization = `Bearer ${token}`;
  }
  return cfg;
});

let redirecting = false;

instance.interceptors.response.use(
  (resp: AxiosResponse) => {
    const data = resp.data;
    if (data && typeof data === 'object' && 'code' in data) {
      const code = (data as any).code;
      const msg = (data as any).message || '请求失败';
      if (code === 0) return (data as any).data;
      if (code === 401 || code === 403) {
        if (!redirecting) {
          redirecting = true;
          tokenStorage.clear();
          userStorage.clear();
          Toast.show({ icon: 'fail', content: '登录已失效，请重新登录' });
          setTimeout(() => {
            redirecting = false;
            const path = window.location.pathname;
            if (!path.startsWith('/login') && !path.startsWith('/register')) {
              window.location.href = '/login?redirect=' + encodeURIComponent(path);
            }
          }, 600);
        }
        return Promise.reject(new Error(msg));
      }
      Toast.show({ icon: 'fail', content: msg });
      return Promise.reject(new Error(msg));
    }
    return data;
  },
  (err) => {
    const msg = err?.response?.data?.message || err?.message || '网络异常';
    if (err?.response?.status === 401) {
      if (!redirecting) {
        redirecting = true;
        tokenStorage.clear();
        userStorage.clear();
        setTimeout(() => {
          redirecting = false;
          const path = window.location.pathname;
          if (!path.startsWith('/login')) {
            window.location.href = '/login?redirect=' + encodeURIComponent(path);
          }
        }, 600);
      }
    } else {
      Toast.show({ icon: 'fail', content: msg });
    }
    return Promise.reject(err);
  },
);

export default instance;
