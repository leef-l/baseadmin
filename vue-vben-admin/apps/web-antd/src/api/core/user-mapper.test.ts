import { describe, expect, it } from 'vitest';

import { mapToUserInfo } from './user-mapper';

describe('mapToUserInfo', () => {
  it('maps backend payload to Vben user info', () => {
    expect(
      mapToUserInfo(
        {
          userId: '1',
          username: 'admin',
          nickname: '管理员',
          email: 'admin@example.com',
          avatar: 'https://example.com/avatar.png',
          deptId: '10',
          status: 1,
          roles: ['admin'],
          perms: ['system:*'],
        },
        '/system/dept',
      ),
    ).toEqual({
      userId: '1',
      username: 'admin',
      realName: '管理员',
      avatar: 'https://example.com/avatar.png',
      roles: ['admin'],
      desc: '',
      homePath: '/system/dept',
      token: '',
    });
  });

  it('falls back to username, empty avatar and empty roles', () => {
    expect(
      mapToUserInfo(
        {
          userId: '2',
          username: 'editor',
          nickname: '',
          email: 'editor@example.com',
          avatar: '',
          deptId: '11',
          status: 1,
          roles: undefined as unknown as string[],
          perms: [],
        },
        '/dashboard',
      ),
    ).toEqual({
      userId: '2',
      username: 'editor',
      realName: 'editor',
      avatar: '',
      roles: [],
      desc: '',
      homePath: '/dashboard',
      token: '',
    });
  });
});
