/** 上传配置类型定义 */

/** 上传配置项 */
export interface ConfigItem {
  id: string;
  name: string;
  storage?: number;
  isDefault?: number;
  localPath?: string;
  ossEndpoint?: string;
  ossBucket?: string;
  ossAccessKey?: string;
  ossSecretKey?: string;
  cosRegion?: string;
  cosBucket?: string;
  cosSecretID?: string;
  cosSecretKey?: string;
  maxSize?: number;
  status?: number;
  createdAt?: string;
  updatedAt?: string;
}

/** 上传配置列表查询参数 */
export interface ConfigListParams {
  pageNum: number;
  pageSize: number;
  storage?: number;
  isDefault?: number;
  status?: number;
}

/** 上传配置创建参数 */
export interface ConfigCreateParams {
  name: string;
  storage?: number;
  isDefault?: number;
  localPath?: string;
  ossEndpoint?: string;
  ossBucket?: string;
  ossAccessKey?: string;
  ossSecretKey?: string;
  cosRegion?: string;
  cosBucket?: string;
  cosSecretID?: string;
  cosSecretKey?: string;
  maxSize?: number;
  status?: number;
}

/** 上传配置更新参数 */
export interface ConfigUpdateParams {
  id: string;
  name: string;
  storage?: number;
  isDefault?: number;
  localPath?: string;
  ossEndpoint?: string;
  ossBucket?: string;
  ossAccessKey?: string;
  ossSecretKey?: string;
  cosRegion?: string;
  cosBucket?: string;
  cosSecretID?: string;
  cosSecretKey?: string;
  maxSize?: number;
  status?: number;
}
