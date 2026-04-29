<script setup lang="ts">
import { h, onMounted, ref } from 'vue';
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { downloadFileFromBlob } from '@vben/utils';
import { Button, message, Modal, Tag, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getContractList, deleteContract, batchDeleteContract, exportContract, importContract, downloadImportTemplateContract, batchUpdateContract } from '#/api/demo/contract';
import { getCustomerList } from '#/api/demo/customer';
import { getOrderList } from '#/api/demo/order';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
import type { ContractItem } from '#/api/demo/contract/types';
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
  contractNo: 'contract_no',
  title: 'title',
  contractFile: 'contract_file',
  signImage: 'sign_image',
  contractAmount: 'contract_amount',
  signPassword: 'sign_password',
  signedAt: 'signed_at',
  expiresAt: 'expires_at',
  status: 'status',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}

/** 状态选项 */
const statusOptions = [
  { label: '待审核', value: 0 },
  { label: '已通过', value: 1 },
  { label: '已拒绝', value: 2 },
  { label: '已取消', value: 3 },
];

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '待审核',
  1: '已通过',
  2: '已拒绝',
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
/** 渲染带 Tooltip 的列标题 */
function tooltipHeader(label: string, tip: string) {
  return () => h('span', {}, [
    label + ' ',
    h(Tooltip, { title: tip }, {
      default: () => h(QuestionCircleOutlined, { style: { color: '#999', marginLeft: '4px' } }),
    }),
  ]);
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
const canBatchDelete = hasAccessByCodes(['demo:contract:batch-delete']);
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
      componentProps: { placeholder: '请输入合同编号', allowClear: true },
      fieldName: 'contractNo',
      label: '合同编号',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入合同标题', allowClear: true },
      fieldName: 'title',
      label: '合同标题',
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
        options: statusOptions,
        placeholder: '请选择状态',
        class: 'w-full',
      },
      fieldName: 'status',
      label: '状态',
    },
    {
      component: 'RangePicker',
      fieldName: 'signedAtRange',
      label: '签署时间',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        valueFormat: 'YYYY-MM-DD HH:mm:ss',
        class: 'w-full',
      },
    },
    {
      component: 'RangePicker',
      fieldName: 'expiresAtRange',
      label: '到期时间',
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
const gridOptions: VxeGridProps<ContractItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { title: '序号', type: 'seq', width: 50 },
    { field: 'contractNo', title: '合同编号' },
    { field: 'customerName', title: '客户' },
    { field: 'orderOrderNo', title: '订单' },
    { field: 'title', title: '合同标题' },
    { field: 'contractFile', title: '合同文件', slots: { default: 'contractFile_cell' } },
    { field: 'signImage', title: '签章图片', width: 100, slots: { default: 'signImage_cell' } },
    { field: 'contractAmount', title: '合同金额', slots: { header: tooltipHeader('合同金额', '分') }, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-' },
    { field: 'status', title: '状态', width: 120, slots: { default: 'status_cell' } },
    ...(isPlatformSuperAdmin.value ? [
    { field: 'tenantName', title: '租户' },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    { field: 'merchantName', title: '商户' },
    ] : []),
    { field: 'signedAt', title: '签署时间', width: 180, formatter: 'formatDateTime' },
    { field: 'expiresAt', title: '到期时间', width: 180, formatter: 'formatDateTime' },
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
        if (params.signedAtRange && params.signedAtRange.length === 2) {
          params.signedAtStart = params.signedAtRange[0];
          params.signedAtEnd = params.signedAtRange[1];
        }
        delete params.signedAtRange;
        if (params.expiresAtRange && params.expiresAtRange.length === 2) {
          params.expiresAtStart = params.expiresAtRange[0];
          params.expiresAtEnd = params.expiresAtRange[1];
        }
        delete params.expiresAtRange;
        if (sorts && sorts.length > 0) {
          const sort = sorts[0];
          if (sort && sort.field && sort.order) {
            params.orderBy = resolveSortField(String(sort.field));
            params.orderDir = sort.order;
          }
        }
        const res = await getContractList(params as any);
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
function handleView(row: ContractItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: ContractItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: ContractItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该体验合同吗？',
    okType: 'danger',
    async onOk() {
      await deleteContract(row.id);
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
    content: `确定要删除选中的 ${rows.length} 条体验合同吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteContract(rows.map((r: ContractItem) => r.id));
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
    if (params.signedAtRange && params.signedAtRange.length === 2) {
      params.signedAtStart = params.signedAtRange[0];
      params.signedAtEnd = params.signedAtRange[1];
    }
    delete params.signedAtRange;
    if (params.expiresAtRange && params.expiresAtRange.length === 2) {
      params.expiresAtStart = params.expiresAtRange[0];
      params.expiresAtEnd = params.expiresAtRange[1];
    }
    delete params.expiresAtRange;
    if (sorts.length > 0) {
      const sort = sorts[0];
      if (sort?.field && sort?.order) {
        params.orderBy = resolveSortField(String(sort.field));
        params.orderDir = sort.order;
      }
    }
    const blob = await exportContract(params);
    downloadFileFromBlob({ fileName: '体验合同.csv', source: blob as Blob });
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
    const res = await importContract(formData);
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
    const blob = await downloadImportTemplateContract();
    downloadFileFromBlob({ fileName: '体验合同导入模板.csv', source: blob as Blob });
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
      await batchUpdateContract({ ids: rows.map((r: ContractItem) => r.id), status: newStatus });
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
        <Button v-auth="['demo:contract:create']" type="primary" @click="handleCreate">新建</Button>
        <Button v-auth="['demo:contract:batch-delete']" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-auth="['demo:contract:export']" class="ml-2" @click="handleExport">导出</Button>
        <Button v-auth="['demo:contract:import']" class="ml-2" @click="handleImportTrigger">导入</Button>
        <Button class="ml-2" @click="handleDownloadTemplate">模板下载</Button>
        <Button v-auth="['demo:contract:batch-update']" class="ml-2" @click="handleBatchUpdateStatus">批量修改状态</Button>
      </template>
      <template #contractFile_cell="{ row }">
        <a v-if="row.contractFile" :href="row.contractFile" target="_blank" rel="noreferrer noopener" style="color: #1890ff;">下载</a>
        <span v-else>-</span>
      </template>
      <template #signImage_cell="{ row }">
        <img v-if="row.signImage" :src="row.signImage" style="width: 48px; height: 48px; object-fit: cover; border-radius: 4px;" />
        <span v-else>-</span>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status)">
          {{ getEnumLabel(statusMap, row.status) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-auth="['demo:contract:detail']" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-auth="['demo:contract:update']" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-auth="['demo:contract:delete']" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
