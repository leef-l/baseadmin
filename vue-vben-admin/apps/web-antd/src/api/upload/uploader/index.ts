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

export function uploadFile(file: File, dirId?: string) {
  const formData = new FormData();
  formData.append('file', file);
  if (dirId) {
    formData.append('dirId', dirId);
  }
  return requestClient.post<UploadResult>('/upload/uploader/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
}
