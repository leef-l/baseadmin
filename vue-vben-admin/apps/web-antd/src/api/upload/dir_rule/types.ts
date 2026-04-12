/** 文件目录规则类型定义 */

export type DirRuleStorageType = 1 | 2 | 3;
export type DirRuleStorageTypesValue = DirRuleStorageType[] | string;

/** 文件目录规则项 */
export interface DirRuleItem {
  id: string;
  dirID: string;
  dirName?: string;
  category?: number;
  fileType?: string;
  storageTypes?: string;
  savePath?: string;
  status?: number;
  createdAt?: string;
  updatedAt?: string;
}

/** 文件目录规则列表查询参数 */
export interface DirRuleListParams {
  pageNum: number;
  pageSize: number;
  keyword?: string;
  category?: number;
  status?: number;
}

/** 文件目录规则创建参数 */
export interface DirRuleCreateParams {
  dirID: string;
  category?: number;
  fileType?: string;
  storageTypes?: DirRuleStorageTypesValue;
  savePath?: string;
  status?: number;
}

/** 文件目录规则更新参数 */
export interface DirRuleUpdateParams {
  id: string;
  dirID: string;
  category?: number;
  fileType?: string;
  storageTypes?: DirRuleStorageTypesValue;
  savePath?: string;
  status?: number;
}
