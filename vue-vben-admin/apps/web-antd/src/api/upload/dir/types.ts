/** 文件目录类型定义 */

/** 文件目录项 */
export interface DirItem {
  id: string;
  parentID?: string;
  dirName?: string;
  name: string;
  path: string;
  keepName?: number;
  sort?: number;
  status?: number;
  createdAt?: string;
  updatedAt?: string;
  children?: DirItem[];
}

/** 文件目录列表查询参数 */
export interface DirListParams {
  pageNum: number;
  pageSize: number;
  keyword?: string;
  status?: number;
}

/** 文件目录树形查询参数 */
export interface DirTreeParams {
  keyword?: string;
  status?: number;
}

/** 文件目录创建参数 */
export interface DirCreateParams {
  parentID?: string;
  name: string;
  path: string;
  keepName?: number;
  sort?: number;
  status?: number;
}

/** 文件目录更新参数 */
export interface DirUpdateParams {
  id: string;
  parentID?: string;
  name: string;
  path: string;
  keepName?: number;
  sort?: number;
  status?: number;
}
