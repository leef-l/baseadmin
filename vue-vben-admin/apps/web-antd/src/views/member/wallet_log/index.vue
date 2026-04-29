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
import { getWalletLogList, deleteWalletLog, batchDeleteWalletLog, exportWalletLog, importWalletLog, downloadImportTemplateWalletLog } from '#/api/member/wallet_log';
import { getUserTree } from '#/api/member/user';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
import type { WalletLogItem } from '#/api/member/wallet_log/types';
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
  changeAmount: 'change_amount',
  beforeBalance: 'before_balance',
  afterBalance: 'after_balance',
  relatedOrderNo: 'related_order_no',
  remark: 'remark',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}

/** 钱包类型选项 */
const walletTypeOptions = [
  { label: '优惠券余额', value: 1 },
  { label: '奖金余额', value: 2 },
  { label: '推广奖余额', value: 3 },
];

/** 钱包类型映射 */
const walletTypeMap: Record<EnumValue, string> = {
  1: '优惠券余额',
  2: '奖金余额',
  3: '推广奖余额',
};

/** 钱包类型颜色 */
function getWalletTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 变动类型选项 */
const changeTypeOptions = [
  { label: '充值', value: 1 },
  { label: '消费', value: 2 },
  { label: '推广奖', value: 3 },
  { label: '仓库卖出收入', value: 4 },
  { label: '平台扣除', value: 5 },
  { label: '后台调整', value: 6 },
];

/** 变动类型映射 */
const changeTypeMap: Record<EnumValue, string> = {
  1: '充值',
  2: '消费',
  3: '推广奖',
  4: '仓库卖出收入',
  5: '平台扣除',
  6: '后台调整',
};

/** 变动类型颜色 */
function getChangeTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4, 5, 6];
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
const canBatchDelete = hasAccessByCodes(['member:wallet_log:batch-delete']);
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
      componentProps: { placeholder: '请输入关联单号', allowClear: true },
      fieldName: 'relatedOrderNo',
      label: '关联单号',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择会员',
        class: 'w-full',
      },
      fieldName: 'userID',
      label: '会员',
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
        options: walletTypeOptions,
        placeholder: '请选择钱包类型',
        class: 'w-full',
      },
      fieldName: 'walletType',
      label: '钱包类型',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: changeTypeOptions,
        placeholder: '请选择变动类型',
        class: 'w-full',
      },
      fieldName: 'changeType',
      label: '变动类型',
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
const gridOptions: VxeGridProps<WalletLogItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { field: 'userNickname', title: '会员' },
    { field: 'walletType', title: '钱包类型', width: 120, slots: { default: 'walletType_cell' } },
    { field: 'changeType', title: '变动类型', width: 120, slots: { default: 'changeType_cell' } },
    { field: 'changeAmount', title: '变动金额', slots: { header: tooltipHeader('变动金额', '分，正增负减') }, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-', sortable: true },
    { field: 'beforeBalance', title: '变动前余额', slots: { header: tooltipHeader('变动前余额', '分') }, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-', sortable: true },
    { field: 'afterBalance', title: '变动后余额', slots: { header: tooltipHeader('变动后余额', '分') }, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-', sortable: true },
    { field: 'relatedOrderNo', title: '关联单号', sortable: true },
    { field: 'remark', title: '备注说明', sortable: true },
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
        const res = await getWalletLogList(params as any);
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
function handleView(row: WalletLogItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: WalletLogItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: WalletLogItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该钱包流水记录吗？',
    okType: 'danger',
    async onOk() {
      await deleteWalletLog(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

/** 批量删除 */
function handleBatchDelete() {
  const ids = getGridSelectedIds<WalletLogItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要删除的数据');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${ids.length} 条钱包流水记录吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteWalletLog(ids);
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
    const blob = await exportWalletLog(params);
    downloadFileFromBlob({ fileName: '钱包流水记录.csv', source: blob as Blob });
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
    const res = await importWalletLog(formData);
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
    const blob = await downloadImportTemplateWalletLog();
    downloadFileFromBlob({ fileName: '钱包流水记录导入模板.csv', source: blob as Blob });
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
        <Button v-access:code="'member:wallet_log:create'" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'member:wallet_log:batch-delete'" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-access:code="'member:wallet_log:export'" class="ml-2" @click="handleExport">导出</Button>
        <Button v-access:code="'member:wallet_log:import'" class="ml-2" @click="handleImportTrigger">导入</Button>
        <Button v-access:code="'member:wallet_log:import'" class="ml-2" @click="handleDownloadTemplate">模板下载</Button>
      </template>
      <template #walletType_cell="{ row }">
        <Tag :color="getWalletTypeColor(row.walletType)">
          {{ getEnumLabel(walletTypeMap, row.walletType) }}
        </Tag>
      </template>
      <template #changeType_cell="{ row }">
        <Tag :color="getChangeTypeColor(row.changeType)">
          {{ getEnumLabel(changeTypeMap, row.changeType) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-access:code="'member:wallet_log:detail'" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-access:code="'member:wallet_log:update'" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-access:code="'member:wallet_log:delete'" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
