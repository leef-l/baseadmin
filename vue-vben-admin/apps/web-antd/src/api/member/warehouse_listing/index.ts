import { requestClient } from '#/api/request';

import type {
  WarehouseListingItem,
  WarehouseListingListParams,
  WarehouseListingCreateParams,
  WarehouseListingUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/warehouse_listing';

/** 获取仓库挂卖记录列表 */
export function getWarehouseListingList(params: WarehouseListingListParams) {
  return requestClient.get<{ list: WarehouseListingItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取仓库挂卖记录详情 */
export function getWarehouseListingDetail(id: string) {
  return requestClient.get<WarehouseListingItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建仓库挂卖记录 */
export function createWarehouseListing(data: WarehouseListingCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新仓库挂卖记录 */
export function updateWarehouseListing(data: WarehouseListingUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除仓库挂卖记录 */
export function deleteWarehouseListing(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除仓库挂卖记录 */
export function batchDeleteWarehouseListing(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出仓库挂卖记录 */
export function exportWarehouseListing(params?: Partial<WarehouseListingListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入仓库挂卖记录 */
export function importWarehouseListing(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载仓库挂卖记录导入模板 */
export function downloadImportTemplateWarehouseListing() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑仓库挂卖记录 */
export function batchUpdateWarehouseListing(data: { ids: string[]; listingStatus?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
