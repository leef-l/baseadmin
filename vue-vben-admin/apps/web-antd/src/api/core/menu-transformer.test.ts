import { describe, expect, it } from 'vitest';

import { transformMenus } from './menu-transformer';

describe('transformMenus', () => {
  it('filters hidden and non-route menu types', () => {
    expect(
      transformMenus([
        {
          id: '1',
          parentId: '0',
          title: 'Dept',
          type: 2,
          path: '/system/dept',
          component: '/system/dept/index',
          permission: 'system:dept:list',
          icon: 'TeamOutlined',
          sort: 1,
          isShow: 1,
          isCache: 1,
          linkUrl: '',
          status: 1,
        },
        {
          id: '2',
          parentId: '0',
          title: 'Hidden button',
          type: 3,
          path: '/system/button',
          component: '',
          permission: 'system:button:list',
          icon: '',
          sort: 2,
          isShow: 1,
          isCache: 0,
          linkUrl: '',
          status: 1,
        },
        {
          id: '3',
          parentId: '0',
          title: 'Not shown',
          type: 2,
          path: '/system/hidden',
          component: '',
          permission: '',
          icon: '',
          sort: 3,
          isShow: 0,
          isCache: 0,
          linkUrl: '',
          status: 1,
        },
      ]),
    ).toEqual([
      {
        name: 'system-dept',
        path: '/system/dept',
        component: '/system/dept/index',
        meta: {
          title: 'Dept',
          icon: 'TeamOutlined',
          order: 1,
          hideInMenu: false,
          keepAlive: true,
          authority: ['system:dept:list'],
        },
      },
    ]);
  });

  it('keeps link metadata for external menus and transforms children', () => {
    expect(
      transformMenus([
        {
          id: '10',
          parentId: '0',
          title: 'Docs',
          type: 4,
          path: '/docs',
          component: '',
          permission: '',
          icon: '',
          sort: 10,
          isShow: 1,
          isCache: 0,
          linkUrl: 'https://example.com',
          status: 1,
          children: [
            {
              id: '11',
              parentId: '10',
              title: 'Child',
              type: 2,
              path: '/docs/child',
              component: '/docs/child/index',
              permission: '',
              icon: '',
              sort: 1,
              isShow: 1,
              isCache: 0,
              linkUrl: '',
              status: 1,
            },
          ],
        },
      ]),
    ).toEqual([
      {
        name: 'docs',
        path: '/docs',
        component: 'IFrameView',
        meta: {
          title: 'Docs',
          icon: undefined,
          order: 10,
          hideInMenu: false,
          keepAlive: false,
          authority: undefined,
          link: 'https://example.com',
        },
        children: [
          {
            name: 'docs-child',
            path: '/docs/child',
            component: '/docs/child/index',
            meta: {
              title: 'Child',
              icon: undefined,
              order: 1,
              hideInMenu: false,
              keepAlive: false,
              authority: undefined,
            },
          },
        ],
      },
    ]);
  });

  it('maps embedded menus to iframe routes', () => {
    expect(
      transformMenus([
        {
          id: '20',
          parentId: '0',
          title: 'Embedded Docs',
          type: 5,
          path: '/embedded/docs',
          component: '',
          permission: '',
          icon: '',
          sort: 20,
          isShow: 1,
          isCache: 0,
          linkURL: 'https://example.com/embed',
          status: 1,
        },
      ]),
    ).toEqual([
      {
        name: 'embedded-docs',
        path: '/embedded/docs',
        component: 'IFrameView',
        meta: {
          title: 'Embedded Docs',
          icon: undefined,
          order: 20,
          hideInMenu: false,
          keepAlive: false,
          authority: undefined,
          iframeSrc: 'https://example.com/embed',
        },
      },
    ]);
  });

  it('skips invalid and duplicate route records', () => {
    expect(
      transformMenus([
        {
          id: '30',
          parentId: '0',
          title: 'Broken Menu',
          type: 2,
          path: '/broken',
          component: '',
          permission: '',
          icon: '',
          sort: 1,
          isShow: 1,
          isCache: 0,
          linkUrl: '',
          status: 1,
        },
        {
          id: '31',
          parentId: '0',
          title: 'Valid Menu',
          type: 2,
          path: '/system/users',
          component: '/system/users/index',
          permission: '',
          icon: '',
          sort: 2,
          isShow: 1,
          isCache: 0,
          linkUrl: '',
          status: 1,
        },
        {
          id: '32',
          parentId: '0',
          title: 'Duplicate Path',
          type: 4,
          path: '/system/users',
          component: '',
          permission: '',
          icon: '',
          sort: 3,
          isShow: 1,
          isCache: 0,
          linkUrl: 'https://example.com/users',
          status: 1,
        },
      ]),
    ).toEqual([
      {
        name: 'system-users',
        path: '/system/users',
        component: '/system/users/index',
        meta: {
          title: 'Valid Menu',
          icon: undefined,
          order: 2,
          hideInMenu: false,
          keepAlive: false,
          authority: undefined,
        },
      },
    ]);
  });
});
