import { requestClient } from '#/api/request';

import type {
  ConfigItem,
  ConfigListParams,
  ConfigCreateParams,
  ConfigUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/upload/config';

/** 获取上传配置列表 */
export function getConfigList(params: ConfigListParams) {
  return requestClient.get<{ list: ConfigItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取上传配置详情 */
export function getConfigDetail(id: string) {
  return requestClient.get<ConfigItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建上传配置 */
export function createConfig(data: ConfigCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新上传配置 */
export function updateConfig(data: ConfigUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除上传配置 */
export function deleteConfig(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}
