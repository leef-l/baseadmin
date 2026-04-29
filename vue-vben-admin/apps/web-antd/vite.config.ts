import { defineConfig } from '@vben/vite-config';
import { loadEnv } from 'vite';

export default defineConfig(async (configEnv) => {
  const env = loadEnv(configEnv?.mode ?? 'development', process.cwd(), '');
  const systemProxyTarget =
    process.env.VITE_PROXY_SYSTEM_TARGET ||
    env.VITE_PROXY_SYSTEM_TARGET ||
    'http://127.0.0.1:10022';
  const uploadProxyTarget =
    process.env.VITE_PROXY_UPLOAD_TARGET ||
    env.VITE_PROXY_UPLOAD_TARGET ||
    'http://127.0.0.1:10023';
  const demoProxyTarget =
    process.env.VITE_PROXY_DEMO_TARGET ||
    env.VITE_PROXY_DEMO_TARGET ||
    'http://127.0.0.1:10026';

  return {
    application: {},
    vite: {
      define: {
        tinymce: 'window.tinymce',
      },
      server: {
        proxy: {
          '/api/system': {
            changeOrigin: true,
            target: systemProxyTarget,
            ws: true,
          },
          '/api/upload': {
            changeOrigin: true,
            target: uploadProxyTarget,
            ws: true,
          },
          '/api/demo': {
            changeOrigin: true,
            target: demoProxyTarget,
            ws: true,
          },
        },
      },
    },
  };
});
