import { requestClient } from '#/api/request';

import type {
  FileItem,
  FileListParams,
  FileCreateParams,
  FileUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/upload/file';

export interface UploadResult {
  id: string;
  url: string;
  name: string;
  size: number;
  ext: string;
  mime: string;
  isImage: number;
}

function normalizeUploadSource(value: string) {
  let next = value.trim();
  if (!next) {
    return '';
  }
  if (next.startsWith('#')) {
    next = next.slice(1);
  }
  next = next.replace(/\\/g, '/');
  next = next.split('#')[0] || '';
  next = next.split('?')[0] || '';
  const segments = next.split('/').filter(Boolean);
  return segments.length > 0 ? `/${segments.join('/')}` : '/';
}

function resolveCurrentUploadSource() {
  if (typeof window === 'undefined') {
    return '';
  }
  const hashSource = normalizeUploadSource(window.location.hash || '');
  if (hashSource && hashSource !== '/') {
    return hashSource;
  }
  return normalizeUploadSource(window.location.pathname || '');
}

/** 获取文件记录列表 */
export function getFileList(params: FileListParams) {
  return requestClient.get<{ list: FileItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取文件记录详情 */
export function getFileDetail(id: string) {
  return requestClient.get<FileItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建文件记录 */
export function createFile(data: FileCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 上传文件 */
export function uploadFile(file: File, dirId?: string, source = resolveCurrentUploadSource()) {
  const formData = new FormData();
  formData.append('file', file);
  if (dirId) {
    formData.append('dirId', dirId);
  }
  if (source) {
    formData.append('source', source);
  }
  return requestClient.post<UploadResult>(`${PREFIX}/upload`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
}

/** 更新文件记录 */
export function updateFile(data: FileUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除文件记录 */
export function deleteFile(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除文件记录 */
export function batchDeleteFile(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}
