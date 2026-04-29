import { requestClient } from '#/api/request';

import type {
  WarehouseTradeItem,
  WarehouseTradeListParams,
  WarehouseTradeCreateParams,
  WarehouseTradeUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/warehouse_trade';

/** 获取仓库交易记录列表 */
export function getWarehouseTradeList(params: WarehouseTradeListParams) {
  return requestClient.get<{ list: WarehouseTradeItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取仓库交易记录详情 */
export function getWarehouseTradeDetail(id: string) {
  return requestClient.get<WarehouseTradeItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建仓库交易记录 */
export function createWarehouseTrade(data: WarehouseTradeCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新仓库交易记录 */
export function updateWarehouseTrade(data: WarehouseTradeUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除仓库交易记录 */
export function deleteWarehouseTrade(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除仓库交易记录 */
export function batchDeleteWarehouseTrade(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出仓库交易记录 */
export function exportWarehouseTrade(params?: Partial<WarehouseTradeListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入仓库交易记录 */
export function importWarehouseTrade(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载仓库交易记录导入模板 */
export function downloadImportTemplateWarehouseTrade() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑仓库交易记录 */
export function batchUpdateWarehouseTrade(data: { ids: string[]; tradeStatus?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
