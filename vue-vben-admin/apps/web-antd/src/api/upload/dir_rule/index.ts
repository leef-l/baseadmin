import { requestClient } from '#/api/request';

import type {
  DirRuleItem,
  DirRuleListParams,
  DirRuleCreateParams,
  DirRuleUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/upload/dir_rule';

/** 获取文件目录规则列表 */
export function getDirRuleList(params: DirRuleListParams) {
  return requestClient.get<{ list: DirRuleItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取文件目录规则详情 */
export function getDirRuleDetail(id: string) {
  return requestClient.get<DirRuleItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建文件目录规则 */
export function createDirRule(data: DirRuleCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新文件目录规则 */
export function updateDirRule(data: DirRuleUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除文件目录规则 */
export function deleteDirRule(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}
