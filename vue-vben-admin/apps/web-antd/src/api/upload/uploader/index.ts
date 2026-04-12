import { requestClient } from '#/api/request';

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

export function uploadFile(file: File, dirId?: string, source = resolveCurrentUploadSource()) {
  const formData = new FormData();
  formData.append('file', file);
  if (dirId) {
    formData.append('dirId', dirId);
  }
  if (source) {
    formData.append('source', source);
  }
  return requestClient.post<UploadResult>('/upload/uploader/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
}
