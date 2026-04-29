import { requestClient } from '#/api/request';

import type {
  DaemonBatchOperationResult,
  DaemonCreateParams,
  DaemonItem,
  DaemonListParams,
  DaemonLogResult,
  DaemonOperationResult,
  DaemonUpdateParams,
} from './types';

const PREFIX = '/system/daemon';

export function getDaemonList(params: DaemonListParams) {
  return requestClient.get<{ list: DaemonItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

export function getDaemonDetail(id: string) {
  return requestClient.get<DaemonItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

export function createDaemon(data: DaemonCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

export function updateDaemon(data: DaemonUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

export function deleteDaemon(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

export function batchDeleteDaemon(ids: string[]) {
  return requestClient.delete<DaemonBatchOperationResult>(
    `${PREFIX}/batch-delete`,
    { data: { ids } },
  );
}

export function restartDaemon(id: string) {
  return requestClient.post<DaemonOperationResult>(`${PREFIX}/restart`, { id });
}

export function batchRestartDaemon(ids: string[]) {
  return requestClient.post<DaemonBatchOperationResult>(
    `${PREFIX}/batch-restart`,
    { ids },
  );
}

export function stopDaemon(id: string) {
  return requestClient.post<DaemonOperationResult>(`${PREFIX}/stop`, { id });
}

export function batchStopDaemon(ids: string[]) {
  return requestClient.post<DaemonBatchOperationResult>(`${PREFIX}/batch-stop`, {
    ids,
  });
}

export function getDaemonLog(id: string, logType: 'error' | 'normal') {
  return requestClient.get<DaemonLogResult>(`${PREFIX}/log`, {
    params: { id, lines: 500, logType },
  });
}
