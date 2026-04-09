import type { UserInfo } from '@vben/types';

import { preferences } from '@vben/preferences';

import { getAuthInfoApi } from './auth-info';
import { mapToUserInfo } from './user-mapper';

/**
 * 获取用户信息
 */
export async function getUserInfoApi() {
  const res = await getAuthInfoApi();
  return mapToUserInfo(res, preferences.app.defaultHomePath) as UserInfo;
}
