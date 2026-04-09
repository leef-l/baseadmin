import type { UserInfo } from '@vben/types';

import type { BackendAuthInfo } from './auth-info';

export type BackendUserInfo = BackendAuthInfo;

export function mapToUserInfo(
  payload: BackendUserInfo,
  defaultHomePath: string,
): UserInfo {
  return {
    userId: payload.userId,
    username: payload.username,
    realName: payload.nickname || payload.username,
    avatar: payload.avatar || '',
    roles: payload.roles || [],
    desc: '',
    homePath: defaultHomePath,
    token: '',
  };
}
