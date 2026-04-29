import { requestClient } from '#/api/request';

import type {
  DomainApplyNginxResult,
  DomainApplySSLResult,
  DomainCreateParams,
  DomainItem,
  DomainListParams,
  DomainUpdateParams,
} from './types';

const PREFIX = '/system/domain';

export function getDomainList(params: DomainListParams) {
  return requestClient.get<{ list: DomainItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

export function getDomainDetail(id: string) {
  return requestClient.get<DomainItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

export function createDomain(data: DomainCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

export function updateDomain(data: DomainUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

export function deleteDomain(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

export function batchDeleteDomain(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

export function applyDomainNginx(id: string) {
  return requestClient.post<DomainApplyNginxResult>(`${PREFIX}/apply-nginx`, {
    id,
  });
}

export function applyDomainSSL(id: string) {
  return requestClient.post<DomainApplySSLResult>(`${PREFIX}/apply-ssl`, {
    id,
  });
}
