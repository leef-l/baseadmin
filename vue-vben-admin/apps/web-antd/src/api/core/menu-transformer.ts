import type { RouteRecordStringComponent } from '@vben/types';

export interface BackendMenu {
  id: string;
  parentId?: string;
  parentID?: string;
  title: string;
  type: number;
  path: string;
  component?: null | string;
  permission: string;
  icon: string;
  sort: number;
  isShow: number;
  isCache: number;
  linkUrl?: null | string;
  linkURL?: null | string;
  status: number;
  children?: BackendMenu[];
}

const routeMenuTypes = new Set([1, 2, 4, 5]);

function getMenuLink(menu: BackendMenu) {
  return menu.linkUrl?.trim() || menu.linkURL?.trim() || '';
}

function getMenuComponent(menu: BackendMenu, link: string) {
  const component = menu.component?.trim() || '';

  // Link-style routes still need a valid component fallback so the generated
  // route record remains mountable when the backend leaves component empty.
  if ((menu.type === 4 || menu.type === 5) && link) {
    return component || 'IFrameView';
  }

  return component;
}

export function transformMenus(
  menus: BackendMenu[],
): RouteRecordStringComponent[] {
  return menus
    .filter((menu) => menu.isShow === 1 && routeMenuTypes.has(menu.type))
    .map((menu) => {
      const link = getMenuLink(menu);
      const baseMeta = {
        title: menu.title,
        icon: menu.icon || undefined,
        order: menu.sort,
        hideInMenu: menu.isShow !== 1,
        keepAlive: menu.isCache === 1,
        authority: menu.permission ? [menu.permission] : undefined,
      };
      const route: RouteRecordStringComponent = {
        name: menu.path?.replace(/\//g, '-').replace(/^-/, '') || `menu-${menu.id}`,
        path: menu.path || '',
        component: getMenuComponent(menu, link),
        meta: baseMeta,
      };

      if (menu.children?.length) {
        route.children = transformMenus(menu.children);
      }

      if (menu.type === 4 && link) {
        route.meta = { ...baseMeta, link };
      }

      if (menu.type === 5 && link) {
        route.meta = { ...baseMeta, iframeSrc: link };
      }

      return route;
    });
}
