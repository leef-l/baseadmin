import { requestClient } from '#/api/request';

import type {
  DeptItem,
  DeptListParams,
  DeptTreeParams,
  DeptCreateParams,
  DeptUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/system/dept';

/** 获取部门表列表 */
export function getDeptList(params: DeptListParams) {
  return requestClient.get<{ list: DeptItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取部门表详情 */
export function getDeptDetail(id: string) {
  return requestClient.get<DeptItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建部门表 */
export function createDept(data: DeptCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新部门表 */
export function updateDept(data: DeptUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除部门表 */
export function deleteDept(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除部门表 */
export function batchDeleteDept(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 获取部门表树形结构 */
export async function getDeptTree(params?: DeptTreeParams) {
  const res = await requestClient.get<{ list: DeptItem[] }>(
    `${PREFIX}/tree`,
    { params },
  );
  return res?.list ?? [];
}
