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
import { getShopOrderList, deleteShopOrder, batchDeleteShopOrder, exportShopOrder, importShopOrder, downloadImportTemplateShopOrder, batchUpdateShopOrder } from '#/api/member/shop_order';
import { getUserTree } from '#/api/member/user';
import { getShopGoodsList } from '#/api/member/shop_goods';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
import type { ShopOrderItem } from '#/api/member/shop_order/types';
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
  orderNo: 'order_no',
  goodsTitle: 'goods_title',
  totalPrice: 'total_price',
  remark: 'remark',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}

/** 支付钱包选项 */
const payWalletOptions = [
  { label: '优惠券余额', value: 1 },
];

/** 支付钱包映射 */
const payWalletMap: Record<EnumValue, string> = {
  1: '优惠券余额',
};

/** 支付钱包颜色 */
function getPayWalletColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 订单状态选项 */
const orderStatusOptions = [
  { label: '已完成', value: 1 },
  { label: '已取消', value: 2 },
];

/** 订单状态映射 */
const orderStatusMap: Record<EnumValue, string> = {
  1: '已完成',
  2: '已取消',
};

/** 订单状态颜色 */
function getOrderStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2];
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
const canBatchDelete = hasAccessByCodes(['member:shop_order:batch-delete']);
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
      componentProps: { placeholder: '请输入订单号', allowClear: true },
      fieldName: 'orderNo',
      label: '订单号',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择购买会员',
        class: 'w-full',
      },
      fieldName: 'userID',
      label: '购买会员',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择商品',
        class: 'w-full',
      },
      fieldName: 'goodsID',
      label: '商品',
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
        options: payWalletOptions,
        placeholder: '请选择支付钱包',
        class: 'w-full',
      },
      fieldName: 'payWallet',
      label: '支付钱包',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: orderStatusOptions,
        placeholder: '请选择订单状态',
        class: 'w-full',
      },
      fieldName: 'orderStatus',
      label: '订单状态',
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
      component: 'Input',
      componentProps: { placeholder: '请输入商品名称', allowClear: true },
      fieldName: 'goodsTitle',
      label: '商品名称',
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
const gridOptions: VxeGridProps<ShopOrderItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { field: 'orderNo', title: '订单号', sortable: true },
    { field: 'userNickname', title: '购买会员' },
    { field: 'shopGoodsTitle', title: '商品' },
    { field: 'goodsTitle', title: '商品名称', slots: { header: tooltipHeader('商品名称', '快照') }, sortable: true },
    { field: 'goodsCover', title: '商品封面', width: 100, slots: { default: 'goodsCover_cell' } },
    { field: 'quantity', title: '购买数量' },
    { field: 'totalPrice', title: '订单总价', slots: { header: tooltipHeader('订单总价', '分') }, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-', sortable: true },
    { field: 'payWallet', title: '支付钱包', width: 120, slots: { default: 'payWallet_cell' } },
    { field: 'orderStatus', title: '订单状态', width: 120, slots: { default: 'orderStatus_cell' } },
    { field: 'remark', title: '订单备注', sortable: true },
    { field: 'status', title: '状态', width: 120, slots: { default: 'status_cell' }, sortable: true },
    ...(isPlatformSuperAdmin.value ? [
    { field: 'tenantName', title: '租户' },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    { field: 'merchantName', title: '商户' },
    ] : []),
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
        if (sorts && sorts.length > 0) {
          const sort = sorts[0];
          if (sort && sort.field && sort.order) {
            params.orderBy = resolveSortField(String(sort.field));
            params.orderDir = sort.order;
          }
        }
        const res = await getShopOrderList(params as any);
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
    const userIDTree = await getUserTree();
    gridApi.formApi.updateSchema([
      {
        fieldName: 'userID',
        componentProps: { treeData: userIDTree ?? [] },
      },
    ]);
  } catch {
    // ignore
  }
  try {
    const goodsIDRes = await getShopGoodsList({ pageNum: 1, pageSize: 1000 });
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
function handleView(row: ShopOrderItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: ShopOrderItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: ShopOrderItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该商城订单吗？',
    okType: 'danger',
    async onOk() {
      await deleteShopOrder(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

/** 批量删除 */
function handleBatchDelete() {
  const ids = getGridSelectedIds<ShopOrderItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要删除的数据');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${ids.length} 条商城订单吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteShopOrder(ids);
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
    if (sorts.length > 0) {
      const sort = sorts[0];
      if (sort?.field && sort?.order) {
        params.orderBy = resolveSortField(String(sort.field));
        params.orderDir = sort.order;
      }
    }
    const blob = await exportShopOrder(params);
    downloadFileFromBlob({ fileName: '商城订单.csv', source: blob as Blob });
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
    const res = await importShopOrder(formData);
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
    const blob = await downloadImportTemplateShopOrder();
    downloadFileFromBlob({ fileName: '商城订单导入模板.csv', source: blob as Blob });
  } catch {
    message.error('下载模板失败');
  }
}

/** 批量修改状态 */
function handleBatchUpdateStatus() {
  const ids = getGridSelectedIds<ShopOrderItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要修改的数据');
    return;
  }
  const rows = gridApi.grid.getCheckboxRecords() as ShopOrderItem[];
  Modal.confirm({
    title: '批量修改状态',
    content: `确定要将选中的 ${ids.length} 条数据的状态切换吗？`,
    async onOk() {
      const newStatus = rows[0]?.status === 1 ? 0 : 1;
      await batchUpdateShopOrder({ ids, status: newStatus });
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
        <Button v-access:code="'member:shop_order:create'" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'member:shop_order:batch-delete'" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-access:code="'member:shop_order:export'" class="ml-2" @click="handleExport">导出</Button>
        <Button v-access:code="'member:shop_order:import'" class="ml-2" @click="handleImportTrigger">导入</Button>
        <Button v-access:code="'member:shop_order:import'" class="ml-2" @click="handleDownloadTemplate">模板下载</Button>
        <Button v-access:code="'member:shop_order:batch-update'" class="ml-2" @click="handleBatchUpdateStatus">批量修改状态</Button>
      </template>
      <template #goodsCover_cell="{ row }">
        <img v-if="row.goodsCover && /^https?:\/\//i.test(row.goodsCover)" :src="row.goodsCover" style="width: 48px; height: 48px; object-fit: cover; border-radius: 4px;" />
        <span v-else>-</span>
      </template>
      <template #payWallet_cell="{ row }">
        <Tag :color="getPayWalletColor(row.payWallet)">
          {{ getEnumLabel(payWalletMap, row.payWallet) }}
        </Tag>
      </template>
      <template #orderStatus_cell="{ row }">
        <Tag :color="getOrderStatusColor(row.orderStatus)">
          {{ getEnumLabel(orderStatusMap, row.orderStatus) }}
        </Tag>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status)">
          {{ getEnumLabel(statusMap, row.status) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-access:code="'member:shop_order:detail'" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-access:code="'member:shop_order:update'" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-access:code="'member:shop_order:delete'" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
