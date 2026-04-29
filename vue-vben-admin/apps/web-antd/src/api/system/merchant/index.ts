import { requestClient } from '#/api/request';

import type {
  MerchantCreateParams,
  MerchantItem,
  MerchantListParams,
  MerchantUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/system/merchant';

/** 获取商户列表 */
export function getMerchantList(params: MerchantListParams) {
  return requestClient.get<{ list: MerchantItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取商户详情 */
export function getMerchantDetail(id: string) {
  return requestClient.get<MerchantItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建商户 */
export function createMerchant(data: MerchantCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新商户 */
export function updateMerchant(data: MerchantUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除商户 */
export function deleteMerchant(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除商户 */
export function batchDeleteMerchant(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}
