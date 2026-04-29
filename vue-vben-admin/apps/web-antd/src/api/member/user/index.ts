import { requestClient } from '#/api/request';

import type {
  UserItem,
  UserListParams,
  UserCreateParams,
  UserUpdateParams,
  UserTreeParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/user';

/** 获取会员用户列表 */
export function getUserList(params: UserListParams) {
  return requestClient.get<{ list: UserItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取会员用户详情 */
export function getUserDetail(id: string) {
  return requestClient.get<UserItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建会员用户 */
export function createUser(data: UserCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新会员用户 */
export function updateUser(data: UserUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除会员用户 */
export function deleteUser(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除会员用户 */
export function batchDeleteUser(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出会员用户 */
export function exportUser(params?: Partial<UserListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 获取会员用户树形结构 */
export async function getUserTree(params?: UserTreeParams) {
  const res = await requestClient.get<{ list: UserItem[] }>(`${PREFIX}/tree`, { params });
  return res?.list ?? [];
}

/** 批量编辑会员用户 */
export function batchUpdateUser(data: { ids: string[]; isActive?: number; isQualified?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
