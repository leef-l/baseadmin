import type { RouteRecordStringComponent } from '@vben/types';

import { requestClient } from '#/api/request';

interface BackendMenu {
  id: string;
  parentId: string;
  title: string;
  type: number;
  path: string;
  component: string;
  permission: string;
  icon: string;
  sort: number;
  isShow: number;
  isCache: number;
  linkUrl: string;
  status: number;
  children?: BackendMenu[];
}

/** 将后端菜单转换为 Vben 路由格式 */
function transformMenus(menus: BackendMenu[]): RouteRecordStringComponent[] {
  return menus
    .filter((m) => m.isShow === 1)
    .map((menu) => {
      const route: RouteRecordStringComponent = {
        name: menu.path?.replace(/\//g, '-').replace(/^-/, '') || `menu-${menu.id}`,
        path: menu.path || '',
        component: menu.component || '',
        meta: {
          title: menu.title,
          icon: menu.icon || undefined,
          order: menu.sort,
          hideInMenu: menu.isShow !== 1,
          keepAlive: menu.isCache === 1,
          authority: menu.permission ? [menu.permission] : undefined,
        },
      };

      if (menu.children?.length) {
        route.children = transformMenus(menu.children);
      }

      // 外链
      if (menu.type === 4 && menu.linkUrl) {
        route.meta = { ...route.meta, link: menu.linkUrl };
      }

      return route;
    });
}

/**
 * 获取用户所有菜单
 */
export async function getAllMenusApi() {
  const res = await requestClient.get<{ menus: BackendMenu[] }>(
    '/system/auth/menus',
  );
  const menus = res?.menus ?? [];
  return transformMenus(menus);
}
