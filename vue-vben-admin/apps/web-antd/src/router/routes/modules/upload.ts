import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'CloudUploadOutlined',
      order: 30,
      title: '上传管理',
    },
    name: 'Upload',
    path: '/upload',
    children: [
      {
        path: '/upload/config',
        name: 'UploadConfig',
        component: () => import('#/views/upload/config/index.vue'),
        meta: { title: '上传配置' },
      },
      {
        path: '/upload/dir',
        name: 'UploadDir',
        component: () => import('#/views/upload/dir/index.vue'),
        meta: { title: '文件目录' },
      },
      {
        path: '/upload/dir-rule',
        name: 'UploadDirRule',
        component: () => import('#/views/upload/dir_rule/index.vue'),
        meta: { title: '目录规则' },
      },
      {
        path: '/upload/file',
        name: 'UploadFile',
        component: () => import('#/views/upload/file/index.vue'),
        meta: { title: '文件记录' },
      },
    ],
  },
];

export default routes;
