<script setup lang="ts">
import { onMounted, ref } from 'vue';
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { downloadFileFromBlob } from '@vben/utils';
import { Button, message, Modal, Tag } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getGridSelectedIds } from '#/utils/grid-selection';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getAuditLogList, deleteAuditLog, batchDeleteAuditLog, exportAuditLog, importAuditLog, downloadImportTemplateAuditLog } from '#/api/demo/audit_log';
import { getUsersList } from '#/api/system/users';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
import type { AuditLogItem } from '#/api/demo/audit_log/types';
import FormModal from './modules/form.vue';
import DetailDrawer from './modules/detail-drawer.vue';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

const sortableFieldMap: Record<string, string> = {
  createdAt: 'created_at',
  logNo: 'log_no',
  remark: 'remark',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}

/** 动作选项 */
const actionOptions = [
  { label: '创建', value: 1 },
  { label: '修改', value: 2 },
  { label: '删除', value: 3 },
  { label: '导出', value: 4 },
  { label: '导入', value: 5 },
];

/** 动作映射 */
const actionMap: Record<EnumValue, string> = {
  1: '创建',
  2: '修改',
  3: '删除',
  4: '导出',
  5: '导入',
};

/** 动作颜色 */
function getActionColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4, 5];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 对象类型选项 */
const targetTypeOptions = [
  { label: '客户', value: 1 },
  { label: '商品', value: 2 },
  { label: '订单', value: 3 },
  { label: '工单', value: 4 },
];

/** 对象类型映射 */
const targetTypeMap: Record<EnumValue, string> = {
  1: '客户',
  2: '商品',
  3: '订单',
  4: '工单',
};

/** 对象类型颜色 */
function getTargetTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 结果选项 */
const resultOptions = [
  { label: '失败', value: 0 },
  { label: '成功', value: 1 },
];

/** 结果映射 */
const resultMap: Record<EnumValue, string> = {
  0: '失败',
  1: '成功',
};

/** 结果颜色 */
function getResultColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 表单弹窗 */
const [FormModalComp, formModalApi] = useVbenModal({
  connectedComponent: FormModal,
  destroyOnClose: true,
});

/** 详情抽屉 */
const [DetailDrawerComp, detailDrawerApi] = useVbenModal({
  connectedComponent: DetailDrawer,
  destroyOnClose: true,
});
const { hasAccessByCodes } = useAccess();
const canBatchDelete = hasAccessByCodes(['demo:audit_log:batch-delete']);
const isPlatformSuperAdmin = usePlatformSuperAdmin();

/** 搜索表单配置 */
const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: true,
  submitOnChange: false,
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      componentProps: { placeholder: '请输入关键词', allowClear: true },
      fieldName: 'keyword',
      label: '关键词',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入日志编号', allowClear: true },
      fieldName: 'logNo',
      label: '日志编号',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入对象编号', allowClear: true },
      fieldName: 'targetCode',
      label: '对象编号',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择操作人',
        class: 'w-full',
      },
      fieldName: 'operatorID',
      label: '操作人',
    },
    ...(isPlatformSuperAdmin.value ? [
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择租户',
        class: 'w-full',
      },
      fieldName: 'tenantID',
      label: '租户',
    },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择商户',
        class: 'w-full',
      },
      fieldName: 'merchantID',
      label: '商户',
    },
    ] : []),
    {
      component: 'Input',
      componentProps: { placeholder: '请输入客户端IP', allowClear: true },
      fieldName: 'clientIP',
      label: '客户端IP',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: actionOptions,
        placeholder: '请选择动作',
        class: 'w-full',
      },
      fieldName: 'action',
      label: '动作',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: targetTypeOptions,
        placeholder: '请选择对象类型',
        class: 'w-full',
      },
      fieldName: 'targetType',
      label: '对象类型',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: resultOptions,
        placeholder: '请选择结果',
        class: 'w-full',
      },
      fieldName: 'result',
      label: '结果',
    },
    {
      component: 'RangePicker',
      fieldName: 'occurredAtRange',
      label: '发生时间',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        valueFormat: 'YYYY-MM-DD HH:mm:ss',
        class: 'w-full',
      },
    },
    {
      component: 'RangePicker',
      fieldName: 'timeRange',
      label: '创建时间',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        valueFormat: 'YYYY-MM-DD HH:mm:ss',
        class: 'w-full',
      },
    },
  ],
};

/** 表格列配置 */
const gridOptions: VxeGridProps<AuditLogItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { field: 'logNo', title: '日志编号', sortable: true },
    { field: 'usersUsername', title: '操作人' },
    { field: 'action', title: '动作', width: 120, slots: { default: 'action_cell' } },
    { field: 'targetType', title: '对象类型', width: 120, slots: { default: 'targetType_cell' } },
    { field: 'targetCode', title: '对象编号' },
    { field: 'result', title: '结果', width: 120, slots: { default: 'result_cell' } },
    { field: 'clientIP', title: '客户端IP' },
    { field: 'remark', title: '备注', sortable: true },
    ...(isPlatformSuperAdmin.value ? [
    { field: 'tenantName', title: '租户' },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    { field: 'merchantName', title: '商户' },
    ] : []),
    { field: 'occurredAt', title: '发生时间', width: 180, formatter: 'formatDateTime' },
    { field: 'createdAt', title: '创建时间', width: 180, formatter: 'formatDateTime', sortable: true },
    { title: '操作', width: 240, fixed: 'right', slots: { default: 'action' } },
  ],
  height: 'auto',
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page, sorts }, formValues) => {
        const params: Record<string, any> = {
          pageNum: page.currentPage,
          pageSize: page.pageSize,
          ...formValues,
        };
        if (params.timeRange && params.timeRange.length === 2) {
          params.startTime = params.timeRange[0];
          params.endTime = params.timeRange[1];
        }
        delete params.timeRange;
        if (!isPlatformSuperAdmin.value) {
          delete params.tenantID;
          delete params.merchantID;
        }
        if (params.occurredAtRange && params.occurredAtRange.length === 2) {
          params.occurredAtStart = params.occurredAtRange[0];
          params.occurredAtEnd = params.occurredAtRange[1];
        }
        delete params.occurredAtRange;
        if (sorts && sorts.length > 0) {
          const sort = sorts[0];
          if (sort && sort.field && sort.order) {
            params.orderBy = resolveSortField(String(sort.field));
            params.orderDir = sort.order;
          }
        }
        const res = await getAuditLogList(params as any);
        return { items: res?.list ?? [], total: res?.total ?? 0 };
      },
    },
  },
  sortConfig: {
    remote: true,
    trigger: 'cell',
    defaultSort: { field: 'createdAt', order: 'desc' },
  },
  toolbarConfig: {
    custom: true,
    refresh: true,
    search: true,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
});

const importInputRef = ref<HTMLInputElement | null>(null);

async function initSearchOptions() {
  try {
    const operatorIDRes = await getUsersList({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: 'operatorID',
        componentProps: {
          options: (operatorIDRes?.list ?? []).map((item: any) => ({
            label: item.username || item.id,
            value: item.id,
          })),
        },
      },
    ]);
  } catch {
    // ignore
  }
  if (isPlatformSuperAdmin.value) {
  try {
    const tenantIDRes = await getTenantList({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: 'tenantID',
        componentProps: {
          options: (tenantIDRes?.list ?? []).map((item: any) => ({
            label: item.name || item.id,
            value: item.id,
          })),
        },
      },
    ]);
  } catch {
    // ignore
  }
  }
  if (isPlatformSuperAdmin.value) {
  try {
    const merchantIDRes = await getMerchantList({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: 'merchantID',
        componentProps: {
          options: (merchantIDRes?.list ?? []).map((item: any) => ({
            label: item.name || item.id,
            value: item.id,
          })),
        },
      },
    ]);
  } catch {
    // ignore
  }
  }
}

onMounted(() => {
  void initSearchOptions();
});

/** 新建 */
function handleCreate() {
  formModalApi.setData(null).open();
}

/** 查看 */
function handleView(row: AuditLogItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: AuditLogItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: AuditLogItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该体验审计日志吗？',
    okType: 'danger',
    async onOk() {
      await deleteAuditLog(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

/** 批量删除 */
function handleBatchDelete() {
  const ids = getGridSelectedIds<AuditLogItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要删除的数据');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${ids.length} 条体验审计日志吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteAuditLog(ids);
      message.success('批量删除成功');
      gridApi.reload();
    },
  });
}

/** 导出 */
async function handleExport() {
  try {
    const formValues = await gridApi.formApi.getValues();
    const params: Record<string, any> = { ...formValues };
    const sorts = gridApi.grid?.getSortColumns?.() ?? [];
    if (params.timeRange && params.timeRange.length === 2) {
      params.startTime = params.timeRange[0];
      params.endTime = params.timeRange[1];
    }
    delete params.timeRange;
    if (!isPlatformSuperAdmin.value) {
      delete params.tenantID;
      delete params.merchantID;
    }
    if (params.occurredAtRange && params.occurredAtRange.length === 2) {
      params.occurredAtStart = params.occurredAtRange[0];
      params.occurredAtEnd = params.occurredAtRange[1];
    }
    delete params.occurredAtRange;
    if (sorts.length > 0) {
      const sort = sorts[0];
      if (sort?.field && sort?.order) {
        params.orderBy = resolveSortField(String(sort.field));
        params.orderDir = sort.order;
      }
    }
    const blob = await exportAuditLog(params);
    downloadFileFromBlob({ fileName: '体验审计日志.csv', source: blob as Blob });
    message.success('导出成功');
  } catch {
    message.error('导出失败');
  }
}

function handleImportTrigger() {
  const input = importInputRef.value;
  if (!input) {
    return;
  }
  input.value = '';
  input.click();
}

/** 导入 */
async function handleImportChange(event: Event) {
  const input = event.target as HTMLInputElement | null;
  const file = input?.files?.[0];
  if (!file) {
    return;
  }
  const formData = new FormData();
  formData.append('file', file);
  try {
    const res = await importAuditLog(formData);
    message.success(`导入完成：成功 ${res?.success ?? 0} 条，失败 ${res?.fail ?? 0} 条`);
    gridApi.reload();
  } catch {
    message.error('导入失败');
  } finally {
    if (input) {
      input.value = '';
    }
  }
}

/** 下载导入模板 */
async function handleDownloadTemplate() {
  try {
    const blob = await downloadImportTemplateAuditLog();
    downloadFileFromBlob({ fileName: '体验审计日志导入模板.csv', source: blob as Blob });
  } catch {
    message.error('下载模板失败');
  }
}
</script>

<template>
  <Page auto-content-height>
    <FormModalComp @success="() => gridApi.reload()" />
    <DetailDrawerComp />
    <input
      ref="importInputRef"
      type="file"
      accept=".csv"
      class="hidden"
      @change="handleImportChange"
    />
    <Grid>
      <template #toolbar-actions>
        <Button v-access:code="'demo:audit_log:create'" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'demo:audit_log:batch-delete'" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-access:code="'demo:audit_log:export'" class="ml-2" @click="handleExport">导出</Button>
        <Button v-access:code="'demo:audit_log:import'" class="ml-2" @click="handleImportTrigger">导入</Button>
        <Button v-access:code="'demo:audit_log:import'" class="ml-2" @click="handleDownloadTemplate">模板下载</Button>
      </template>
      <template #action_cell="{ row }">
        <Tag :color="getActionColor(row.action)">
          {{ getEnumLabel(actionMap, row.action) }}
        </Tag>
      </template>
      <template #targetType_cell="{ row }">
        <Tag :color="getTargetTypeColor(row.targetType)">
          {{ getEnumLabel(targetTypeMap, row.targetType) }}
        </Tag>
      </template>
      <template #result_cell="{ row }">
        <Tag :color="getResultColor(row.result)">
          {{ getEnumLabel(resultMap, row.result) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-access:code="'demo:audit_log:detail'" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-access:code="'demo:audit_log:update'" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-access:code="'demo:audit_log:delete'" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
