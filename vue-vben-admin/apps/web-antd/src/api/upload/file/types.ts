/** 文件记录类型定义 */

/** 文件记录项 */
export interface FileItem {
  id: string;
  dirID?: string;
  dirName?: string;
  name: string;
  url: string;
  ext?: string;
  size?: string;
  mime?: string;
  storage?: number;
  isImage?: number;
  createdAt?: string;
  updatedAt?: string;
}

/** 文件记录列表查询参数 */
export interface FileListParams {
  pageNum: number;
  pageSize: number;
  keyword?: string;
  dirID?: string;
  name?: string;
  storage?: number;
  isImage?: number;
}

/** 文件记录创建参数 */
export interface FileCreateParams {
  dirID?: string;
  name: string;
  url: string;
  ext?: string;
  size?: string;
  mime?: string;
  storage?: number;
  isImage?: number;
}

/** 文件记录更新参数 */
export interface FileUpdateParams {
  id: string;
  dirID?: string;
  name: string;
  url: string;
  ext?: string;
  size?: string;
  mime?: string;
  storage?: number;
  isImage?: number;
}
