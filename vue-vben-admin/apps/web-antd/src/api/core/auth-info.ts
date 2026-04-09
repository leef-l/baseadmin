import { requestClient } from '#/api/request';

export interface BackendAuthInfo {
  userId: string;
  username: string;
  nickname: string;
  email: string;
  avatar: string;
  deptId: string;
  status: number;
  roles: string[];
  perms: string[];
}

export async function getAuthInfoApi() {
  return requestClient.get<BackendAuthInfo>('/system/auth/info');
}

export function extractAccessCodes(
  payload: null | Partial<BackendAuthInfo> | undefined,
) {
  return payload?.perms ?? [];
}
