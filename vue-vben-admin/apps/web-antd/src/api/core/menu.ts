import { requestClient } from '#/api/request';
import { transformMenus } from './menu-transformer';
import type { BackendMenu } from './menu-transformer';

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
