/** 会员合同模板类型定义 */

/** 会员合同模板项 */
export interface ContractTemplateItem {
  id: string;
  templateName: string;
  templateType?: string;
  content: string;
  isDefault?: number;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 会员合同模板列表查询参数 */
export interface ContractTemplateListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  templateName?: string;
  tenantID?: string;
  merchantID?: string;
  isDefault?: number;
  status?: number;
  templateType?: string;
}

/** 会员合同模板创建参数 */
export interface ContractTemplateCreateParams {
  templateName: string;
  templateType?: string;
  content: string;
  isDefault?: number;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 会员合同模板更新参数 */
export interface ContractTemplateUpdateParams {
  id: string;
  templateName: string;
  templateType?: string;
  content: string;
  isDefault?: number;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
