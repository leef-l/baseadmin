<script setup lang="ts">
import { onMounted, ref } from 'vue';
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { downloadFileFromBlob } from '@vben/utils';
import { Button, message, Modal, Tag } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getWorkOrderList, deleteWorkOrder, batchDeleteWorkOrder, exportWorkOrder, importWorkOrder, downloadImportTemplateWorkOrder, batchUpdateWorkOrder } from '#/api/demo/work_order';
import { getCustomerList } from '#/api/demo/customer';
import { getProductList } from '#/api/demo/product';
import { getOrderList } from '#/api/demo/order';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
import type { WorkOrderItem } from '#/api/demo/work_order/types';
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
  ticketNo: 'ticket_no',
  title: 'title',
  priority: 'priority',
  sourceType: 'source_type',
  description: 'description',
  attachmentFile: 'attachment_file',
  dueAt: 'due_at',
  status: 'status',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}

/** 优先级选项 */
const priorityOptions = [
  { label: '低', value: 1 },
  { label: '普通', value: 2 },
  { label: '高', value: 3 },
  { label: '紧急', value: 4 },
];

/** 优先级映射 */
const priorityMap: Record<EnumValue, string> = {
  1: '低',
  2: '普通',
  3: '高',
  4: '紧急',
};

/** 优先级颜色 */
function getPriorityColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 来源选项 */
const sourceTypeOptions = [
  { label: '官网', value: 1 },
  { label: '电话', value: 2 },
  { label: '微信', value: 3 },
  { label: '后台', value: 4 },
];

/** 来源映射 */
const sourceTypeMap: Record<EnumValue, string> = {
  1: '官网',
  2: '电话',
  3: '微信',
  4: '后台',
};

/** 来源颜色 */
function getSourceTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态选项 */
const statusOptions = [
  { label: '待处理', value: 0 },
  { label: '进行中', value: 1 },
  { label: '已完成', value: 2 },
  { label: '已取消', value: 3 },
];

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '待处理',
  1: '进行中',
  2: '已完成',
  3: '已取消',
};

/** 状态颜色 */
function getStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1, 2, 3];
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
const canBatchDelete = hasAccessByCodes(['demo:work_order:batch-delete']);
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
      componentProps: { placeholder: '请输入工单号', allowClear: true },
      fieldName: 'ticketNo',
      label: '工单号',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入工单标题', allowClear: true },
      fieldName: 'title',
      label: '工单标题',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择客户',
        class: 'w-full',
      },
      fieldName: 'customerID',
      label: '客户',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择商品',
        class: 'w-full',
      },
      fieldName: 'productID',
      label: '商品',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择订单',
        class: 'w-full',
      },
      fieldName: 'orderID',
      label: '订单',
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
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: priorityOptions,
        placeholder: '请选择优先级',
        class: 'w-full',
      },
      fieldName: 'priority',
      label: '优先级',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: sourceTypeOptions,
        placeholder: '请选择来源',
        class: 'w-full',
      },
      fieldName: 'sourceType',
      label: '来源',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: statusOptions,
        placeholder: '请选择状态',
        class: 'w-full',
      },
      fieldName: 'status',
      label: '状态',
    },
    {
      component: 'RangePicker',
      fieldName: 'dueAtRange',
      label: '截止时间',
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
const gridOptions: VxeGridProps<WorkOrderItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { title: '序号', type: 'seq', width: 50 },
    { field: 'ticketNo', title: '工单号' },
    { field: 'customerName', title: '客户' },
    { field: 'productSkuNo', title: '商品' },
    { field: 'orderOrderNo', title: '订单' },
    { field: 'title', title: '工单标题' },
    { field: 'priority', title: '优先级', width: 120, slots: { default: 'priority_cell' } },
    { field: 'sourceType', title: '来源', width: 120, slots: { default: 'sourceType_cell' } },
    { field: 'description', title: '问题描述' },
    { field: 'attachmentFile', title: '附件', slots: { default: 'attachmentFile_cell' } },
    { field: 'status', title: '状态', width: 120, slots: { default: 'status_cell' } },
    ...(isPlatformSuperAdmin.value ? [
    { field: 'tenantName', title: '租户' },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    { field: 'merchantName', title: '商户' },
    ] : []),
    { field: 'dueAt', title: '截止时间', width: 180, formatter: 'formatDateTime' },
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
        if (params.dueAtRange && params.dueAtRange.length === 2) {
          params.dueAtStart = params.dueAtRange[0];
          params.dueAtEnd = params.dueAtRange[1];
        }
        delete params.dueAtRange;
        if (sorts && sorts.length > 0) {
          const sort = sorts[0];
          if (sort && sort.field && sort.order) {
            params.orderBy = resolveSortField(String(sort.field));
            params.orderDir = sort.order;
          }
        }
        const res = await getWorkOrderList(params as any);
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
    const customerIDRes = await getCustomerList({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: 'customerID',
        componentProps: {
          options: (customerIDRes?.list ?? []).map((item: any) => ({
            label: item.name || item.id,
            value: item.id,
          })),
        },
      },
    ]);
  } catch {
    // ignore
  }
  try {
    const productIDRes = await getProductList({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: 'productID',
        componentProps: {
          options: (productIDRes?.list ?? []).map((item: any) => ({
            label: item.skuNo || item.id,
            value: item.id,
          })),
        },
      },
    ]);
  } catch {
    // ignore
  }
  try {
    const orderIDRes = await getOrderList({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: 'orderID',
        componentProps: {
          options: (orderIDRes?.list ?? []).map((item: any) => ({
            label: item.orderNo || item.id,
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
function handleView(row: WorkOrderItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: WorkOrderItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: WorkOrderItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该体验工单吗？',
    okType: 'danger',
    async onOk() {
      await deleteWorkOrder(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

/** 批量删除 */
function handleBatchDelete() {
  const rows = gridApi.grid.getCheckboxRecords();
  if (rows.length === 0) {
    message.warning('请先选择要删除的数据');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${rows.length} 条体验工单吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteWorkOrder(rows.map((r: WorkOrderItem) => r.id));
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
    if (params.dueAtRange && params.dueAtRange.length === 2) {
      params.dueAtStart = params.dueAtRange[0];
      params.dueAtEnd = params.dueAtRange[1];
    }
    delete params.dueAtRange;
    if (sorts.length > 0) {
      const sort = sorts[0];
      if (sort?.field && sort?.order) {
        params.orderBy = resolveSortField(String(sort.field));
        params.orderDir = sort.order;
      }
    }
    const blob = await exportWorkOrder(params);
    downloadFileFromBlob({ fileName: '体验工单.csv', source: blob as Blob });
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
    const res = await importWorkOrder(formData);
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
    const blob = await downloadImportTemplateWorkOrder();
    downloadFileFromBlob({ fileName: '体验工单导入模板.csv', source: blob as Blob });
  } catch {
    message.error('下载模板失败');
  }
}

/** 批量修改状态 */
function handleBatchUpdateStatus() {
  const rows = gridApi.grid.getCheckboxRecords();
  if (rows.length === 0) {
    message.warning('请先选择要修改的数据');
    return;
  }
  Modal.confirm({
    title: '批量修改状态',
    content: `确定要将选中的 ${rows.length} 条数据的状态切换吗？`,
    async onOk() {
      const newStatus = rows[0]?.status === 1 ? 0 : 1;
      await batchUpdateWorkOrder({ ids: rows.map((r: WorkOrderItem) => r.id), status: newStatus });
      message.success('批量修改成功');
      gridApi.reload();
    },
  });
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
        <Button v-auth="['demo:work_order:create']" type="primary" @click="handleCreate">新建</Button>
        <Button v-auth="['demo:work_order:batch-delete']" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-auth="['demo:work_order:export']" class="ml-2" @click="handleExport">导出</Button>
        <Button v-auth="['demo:work_order:import']" class="ml-2" @click="handleImportTrigger">导入</Button>
        <Button class="ml-2" @click="handleDownloadTemplate">模板下载</Button>
        <Button v-auth="['demo:work_order:batch-update']" class="ml-2" @click="handleBatchUpdateStatus">批量修改状态</Button>
      </template>
      <template #priority_cell="{ row }">
        <Tag :color="getPriorityColor(row.priority)">
          {{ getEnumLabel(priorityMap, row.priority) }}
        </Tag>
      </template>
      <template #sourceType_cell="{ row }">
        <Tag :color="getSourceTypeColor(row.sourceType)">
          {{ getEnumLabel(sourceTypeMap, row.sourceType) }}
        </Tag>
      </template>
      <template #attachmentFile_cell="{ row }">
        <a v-if="row.attachmentFile" :href="row.attachmentFile" target="_blank" rel="noreferrer noopener" style="color: #1890ff;">下载</a>
        <span v-else>-</span>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status)">
          {{ getEnumLabel(statusMap, row.status) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-auth="['demo:work_order:detail']" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-auth="['demo:work_order:update']" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-auth="['demo:work_order:delete']" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
