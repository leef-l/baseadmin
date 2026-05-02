import { requestClient } from '#/api/request';

import type { BizConfigData } from './types';

const PREFIX = '/member/biz-config';

/** 获取业务配置（单例） */
export function getBizConfig() {
  return requestClient.get<BizConfigData>(PREFIX);
}

/** 保存业务配置（整体覆盖） */
export function saveBizConfig(data: BizConfigData) {
  return requestClient.put(PREFIX, data);
}
