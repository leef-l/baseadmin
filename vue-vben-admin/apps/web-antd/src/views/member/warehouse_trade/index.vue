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
import { getGridSelectedIds } from '#/utils/grid-selection';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getWarehouseTradeList, deleteWarehouseTrade, batchDeleteWarehouseTrade, exportWarehouseTrade, importWarehouseTrade, downloadImportTemplateWarehouseTrade, batchUpdateWarehouseTrade } from '#/api/member/warehouse_trade';
import { getWarehouseGoodsList } from '#/api/member/warehouse_goods';
import { getWarehouseListingList } from '#/api/member/warehouse_listing';
import { getUserTree } from '#/api/member/user';
import { getUserTree } from '#/api/member/user';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
import type { WarehouseTradeItem } from '#/api/member/warehouse_trade/types';
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
  status: 'status',
  tradeNo: 'trade_no',
  tradePrice: 'trade_price',
  platformFee: 'platform_fee',
  sellerIncome: 'seller_income',
  remark: 'remark',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}

/** 交易状态选项 */
const tradeStatusOptions = [
  { label: '待卖家确认', value: 1 },
  { label: '已确认完成', value: 2 },
  { label: '已取消', value: 3 },
];

/** 交易状态映射 */
const tradeStatusMap: Record<EnumValue, string> = {
  1: '待卖家确认',
  2: '已确认完成',
  3: '已取消',
};

/** 交易状态颜色 */
function getTradeStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态选项 */
const statusOptions = [
  { label: '关闭', value: 0 },
  { label: '开启', value: 1 },
];

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '关闭',
  1: '开启',
};

/** 状态颜色 */
function getStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
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
const canBatchDelete = hasAccessByCodes(['member:warehouse_trade:batch-delete']);
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
      componentProps: { placeholder: '请输入交易编号', allowClear: true },
      fieldName: 'tradeNo',
      label: '交易编号',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择仓库商品',
        class: 'w-full',
      },
      fieldName: 'goodsID',
      label: '仓库商品',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择挂卖记录',
        class: 'w-full',
      },
      fieldName: 'listingID',
      label: '挂卖记录',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择卖家',
        class: 'w-full',
      },
      fieldName: 'sellerID',
      label: '卖家',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择买家',
        class: 'w-full',
      },
      fieldName: 'buyerID',
      label: '买家',
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
        options: tradeStatusOptions,
        placeholder: '请选择交易状态',
        class: 'w-full',
      },
      fieldName: 'tradeStatus',
      label: '交易状态',
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
      fieldName: 'confirmedAtRange',
      label: '确认时间',
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
const gridOptions: VxeGridProps<WarehouseTradeItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { field: 'tradeNo', title: '交易编号', sortable: true },
    { field: 'warehouseGoodsTitle', title: '仓库商品' },
    { field: 'warehouseListingID', title: '挂卖记录' },
    { field: 'userNickname', title: '卖家' },
    { field: 'buyerNickname', title: '买家' },
    { field: 'tradePrice', title: '成交价格', slots: { header: tooltipHeader('成交价格', '分') }, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-', sortable: true },
    { field: 'platformFee', title: '平台扣除费用', slots: { header: tooltipHeader('平台扣除费用', '分') }, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-', sortable: true },
    { field: 'sellerIncome', title: '卖家实收', slots: { header: tooltipHeader('卖家实收', '分') }, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-', sortable: true },
    { field: 'tradeStatus', title: '交易状态', width: 120, slots: { default: 'tradeStatus_cell' } },
    { field: 'remark', title: '备注', sortable: true },
    { field: 'status', title: '状态', width: 120, slots: { default: 'status_cell' }, sortable: true },
    ...(isPlatformSuperAdmin.value ? [
    { field: 'tenantName', title: '租户' },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    { field: 'merchantName', title: '商户' },
    ] : []),
    { field: 'confirmedAt', title: '确认时间', width: 180, formatter: 'formatDateTime' },
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
        if (params.confirmedAtRange && params.confirmedAtRange.length === 2) {
          params.confirmedAtStart = params.confirmedAtRange[0];
          params.confirmedAtEnd = params.confirmedAtRange[1];
        }
        delete params.confirmedAtRange;
        if (sorts && sorts.length > 0) {
          const sort = sorts[0];
          if (sort && sort.field && sort.order) {
            params.orderBy = resolveSortField(String(sort.field));
            params.orderDir = sort.order;
          }
        }
        const res = await getWarehouseTradeList(params as any);
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
    const goodsIDRes = await getWarehouseGoodsList({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: 'goodsID',
        componentProps: {
          options: (goodsIDRes?.list ?? []).map((item: any) => ({
            label: item.title || item.id,
            value: item.id,
          })),
        },
      },
    ]);
  } catch {
    // ignore
  }
  try {
    const listingIDRes = await getWarehouseListingList({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: 'listingID',
        componentProps: {
          options: (listingIDRes?.list ?? []).map((item: any) => ({
            label: item.id || item.id,
            value: item.id,
          })),
        },
      },
    ]);
  } catch {
    // ignore
  }
  try {
    const sellerIDTree = await getUserTree();
    gridApi.formApi.updateSchema([
      {
        fieldName: 'sellerID',
        componentProps: { treeData: sellerIDTree ?? [] },
      },
    ]);
  } catch {
    // ignore
  }
  try {
    const buyerIDTree = await getUserTree();
    gridApi.formApi.updateSchema([
      {
        fieldName: 'buyerID',
        componentProps: { treeData: buyerIDTree ?? [] },
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
function handleView(row: WarehouseTradeItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: WarehouseTradeItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: WarehouseTradeItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该仓库交易记录吗？',
    okType: 'danger',
    async onOk() {
      await deleteWarehouseTrade(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

/** 批量删除 */
function handleBatchDelete() {
  const ids = getGridSelectedIds<WarehouseTradeItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要删除的数据');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${ids.length} 条仓库交易记录吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteWarehouseTrade(ids);
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
    if (params.confirmedAtRange && params.confirmedAtRange.length === 2) {
      params.confirmedAtStart = params.confirmedAtRange[0];
      params.confirmedAtEnd = params.confirmedAtRange[1];
    }
    delete params.confirmedAtRange;
    if (sorts.length > 0) {
      const sort = sorts[0];
      if (sort?.field && sort?.order) {
        params.orderBy = resolveSortField(String(sort.field));
        params.orderDir = sort.order;
      }
    }
    const blob = await exportWarehouseTrade(params);
    downloadFileFromBlob({ fileName: '仓库交易记录.csv', source: blob as Blob });
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
    const res = await importWarehouseTrade(formData);
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
    const blob = await downloadImportTemplateWarehouseTrade();
    downloadFileFromBlob({ fileName: '仓库交易记录导入模板.csv', source: blob as Blob });
  } catch {
    message.error('下载模板失败');
  }
}

/** 批量修改状态 */
function handleBatchUpdateStatus() {
  const ids = getGridSelectedIds<WarehouseTradeItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要修改的数据');
    return;
  }
  const rows = gridApi.grid.getCheckboxRecords() as WarehouseTradeItem[];
  Modal.confirm({
    title: '批量修改状态',
    content: `确定要将选中的 ${ids.length} 条数据的状态切换吗？`,
    async onOk() {
      const newStatus = rows[0]?.status === 1 ? 0 : 1;
      await batchUpdateWarehouseTrade({ ids, status: newStatus });
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
        <Button v-access:code="'member:warehouse_trade:create'" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'member:warehouse_trade:batch-delete'" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-access:code="'member:warehouse_trade:export'" class="ml-2" @click="handleExport">导出</Button>
        <Button v-access:code="'member:warehouse_trade:import'" class="ml-2" @click="handleImportTrigger">导入</Button>
        <Button v-access:code="'member:warehouse_trade:import'" class="ml-2" @click="handleDownloadTemplate">模板下载</Button>
        <Button v-access:code="'member:warehouse_trade:batch-update'" class="ml-2" @click="handleBatchUpdateStatus">批量修改状态</Button>
      </template>
      <template #tradeStatus_cell="{ row }">
        <Tag :color="getTradeStatusColor(row.tradeStatus)">
          {{ getEnumLabel(tradeStatusMap, row.tradeStatus) }}
        </Tag>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status)">
          {{ getEnumLabel(statusMap, row.status) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-access:code="'member:warehouse_trade:detail'" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-access:code="'member:warehouse_trade:update'" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-access:code="'member:warehouse_trade:delete'" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
