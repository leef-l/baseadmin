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
import { getCampaignList, deleteCampaign, batchDeleteCampaign, exportCampaign, importCampaign, downloadImportTemplateCampaign, batchUpdateCampaign } from '#/api/demo/campaign';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
import type { CampaignItem } from '#/api/demo/campaign/types';
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
  campaignNo: 'campaign_no',
  title: 'title',
  banner: 'banner',
  type: 'type',
  channel: 'channel',
  budgetAmount: 'budget_amount',
  landingURL: 'landing_url',
  ruleJSON: 'rule_json',
  introContent: 'intro_content',
  startAt: 'start_at',
  endAt: 'end_at',
  isPublic: 'is_public',
  status: 'status',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}

/** 活动类型选项 */
const typeOptions = [
  { label: '免费', value: 1 },
  { label: '付费', value: 2 },
  { label: '公开', value: 3 },
  { label: '私密', value: 4 },
];

/** 活动类型映射 */
const typeMap: Record<EnumValue, string> = {
  1: '免费',
  2: '付费',
  3: '公开',
  4: '私密',
};

/** 活动类型颜色 */
function getTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 投放渠道选项 */
const channelOptions = [
  { label: '官网', value: 1 },
  { label: '小程序', value: 2 },
  { label: '短信', value: 3 },
  { label: '线下', value: 4 },
];

/** 投放渠道映射 */
const channelMap: Record<EnumValue, string> = {
  1: '官网',
  2: '小程序',
  3: '短信',
  4: '线下',
};

/** 投放渠道颜色 */
function getChannelColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 是否公开选项 */
const isPublicOptions = [
  { label: '否', value: 0 },
  { label: '是', value: 1 },
];

/** 是否公开映射 */
const isPublicMap: Record<EnumValue, string> = {
  0: '否',
  1: '是',
};

/** 是否公开颜色 */
function getIsPublicColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态选项 */
const statusOptions = [
  { label: '草稿', value: 0 },
  { label: '已发布', value: 1 },
  { label: '已下架', value: 2 },
];

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '草稿',
  1: '已发布',
  2: '已下架',
};

/** 状态颜色 */
function getStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1, 2];
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
const canBatchDelete = hasAccessByCodes(['demo:campaign:batch-delete']);
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
      componentProps: { placeholder: '请输入活动编号', allowClear: true },
      fieldName: 'campaignNo',
      label: '活动编号',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入活动标题', allowClear: true },
      fieldName: 'title',
      label: '活动标题',
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
        options: typeOptions,
        placeholder: '请选择活动类型',
        class: 'w-full',
      },
      fieldName: 'type',
      label: '活动类型',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: channelOptions,
        placeholder: '请选择投放渠道',
        class: 'w-full',
      },
      fieldName: 'channel',
      label: '投放渠道',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: isPublicOptions,
        placeholder: '请选择是否公开',
        class: 'w-full',
      },
      fieldName: 'isPublic',
      label: '是否公开',
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
      fieldName: 'startAtRange',
      label: '开始时间',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        valueFormat: 'YYYY-MM-DD HH:mm:ss',
        class: 'w-full',
      },
    },
    {
      component: 'RangePicker',
      fieldName: 'endAtRange',
      label: '结束时间',
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
const gridOptions: VxeGridProps<CampaignItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { title: '序号', type: 'seq', width: 50 },
    { field: 'campaignNo', title: '活动编号' },
    { field: 'title', title: '活动标题' },
    { field: 'banner', title: '横幅图', width: 100, slots: { default: 'banner_cell' } },
    { field: 'type', title: '活动类型', width: 120, slots: { default: 'type_cell' } },
    { field: 'channel', title: '投放渠道', width: 120, slots: { default: 'channel_cell' } },
    { field: 'budgetAmount', title: '预算金额', slots: { header: tooltipHeader('预算金额', '分') }, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-' },
    { field: 'landingURL', title: '落地页URL', slots: { default: 'landingURL_cell' } },
    { field: 'isPublic', title: '是否公开', width: 120, slots: { default: 'isPublic_cell' } },
    { field: 'status', title: '状态', width: 120, slots: { default: 'status_cell' } },
    ...(isPlatformSuperAdmin.value ? [
    { field: 'tenantName', title: '租户' },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    { field: 'merchantName', title: '商户' },
    ] : []),
    { field: 'startAt', title: '开始时间', width: 180, formatter: 'formatDateTime' },
    { field: 'endAt', title: '结束时间', width: 180, formatter: 'formatDateTime' },
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
        if (params.startAtRange && params.startAtRange.length === 2) {
          params.startAtStart = params.startAtRange[0];
          params.startAtEnd = params.startAtRange[1];
        }
        delete params.startAtRange;
        if (params.endAtRange && params.endAtRange.length === 2) {
          params.endAtStart = params.endAtRange[0];
          params.endAtEnd = params.endAtRange[1];
        }
        delete params.endAtRange;
        if (sorts && sorts.length > 0) {
          const sort = sorts[0];
          if (sort && sort.field && sort.order) {
            params.orderBy = resolveSortField(String(sort.field));
            params.orderDir = sort.order;
          }
        }
        const res = await getCampaignList(params as any);
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
function handleView(row: CampaignItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: CampaignItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: CampaignItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该体验活动吗？',
    okType: 'danger',
    async onOk() {
      await deleteCampaign(row.id);
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
    content: `确定要删除选中的 ${rows.length} 条体验活动吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteCampaign(rows.map((r: CampaignItem) => r.id));
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
    if (params.startAtRange && params.startAtRange.length === 2) {
      params.startAtStart = params.startAtRange[0];
      params.startAtEnd = params.startAtRange[1];
    }
    delete params.startAtRange;
    if (params.endAtRange && params.endAtRange.length === 2) {
      params.endAtStart = params.endAtRange[0];
      params.endAtEnd = params.endAtRange[1];
    }
    delete params.endAtRange;
    if (sorts.length > 0) {
      const sort = sorts[0];
      if (sort?.field && sort?.order) {
        params.orderBy = resolveSortField(String(sort.field));
        params.orderDir = sort.order;
      }
    }
    const blob = await exportCampaign(params);
    downloadFileFromBlob({ fileName: '体验活动.csv', source: blob as Blob });
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
    const res = await importCampaign(formData);
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
    const blob = await downloadImportTemplateCampaign();
    downloadFileFromBlob({ fileName: '体验活动导入模板.csv', source: blob as Blob });
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
      await batchUpdateCampaign({ ids: rows.map((r: CampaignItem) => r.id), status: newStatus });
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
        <Button v-access:code="'demo:campaign:create'" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'demo:campaign:batch-delete'" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-access:code="'demo:campaign:export'" class="ml-2" @click="handleExport">导出</Button>
        <Button v-access:code="'demo:campaign:import'" class="ml-2" @click="handleImportTrigger">导入</Button>
        <Button class="ml-2" @click="handleDownloadTemplate">模板下载</Button>
        <Button v-access:code="'demo:campaign:batch-update'" class="ml-2" @click="handleBatchUpdateStatus">批量修改状态</Button>
      </template>
      <template #banner_cell="{ row }">
        <img v-if="row.banner" :src="row.banner" style="width: 48px; height: 48px; object-fit: cover; border-radius: 4px;" />
        <span v-else>-</span>
      </template>
      <template #type_cell="{ row }">
        <Tag :color="getTypeColor(row.type)">
          {{ getEnumLabel(typeMap, row.type) }}
        </Tag>
      </template>
      <template #channel_cell="{ row }">
        <Tag :color="getChannelColor(row.channel)">
          {{ getEnumLabel(channelMap, row.channel) }}
        </Tag>
      </template>
      <template #landingURL_cell="{ row }">
        <a v-if="row.landingURL" :href="row.landingURL" target="_blank" rel="noreferrer noopener" style="color: #1890ff;">{{ row.landingURL }}</a>
        <span v-else>-</span>
      </template>
      <template #isPublic_cell="{ row }">
        <Tag :color="getIsPublicColor(row.isPublic)">
          {{ getEnumLabel(isPublicMap, row.isPublic) }}
        </Tag>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status)">
          {{ getEnumLabel(statusMap, row.status) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-access:code="'demo:campaign:detail'" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-access:code="'demo:campaign:update'" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-access:code="'demo:campaign:delete'" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
