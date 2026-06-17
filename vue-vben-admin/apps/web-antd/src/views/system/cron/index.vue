<script setup lang="ts">
import { h, onMounted, ref } from 'vue';
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { ActionMoreItem } from '#/components/action-more/index.vue';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { downloadFileFromBlob } from '@vben/utils';
import { Button, message, Modal, Tag, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import ActionMore from '#/components/action-more/index.vue';
import { getGridSelectedIds } from '#/utils/grid-selection';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getCronList, deleteCron, batchDeleteCron, exportCron, importCron, downloadImportTemplateCron, batchUpdateCron, triggerCron, getCronLogList } from '#/api/system/cron';
import type { CronItem, CronLogItem } from '#/api/system/cron/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
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
  name: 'name',
  remark: 'remark',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
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
const canBatchDelete = hasAccessByCodes(['system:cron:batch-delete']);
const canDelete = hasAccessByCodes(['system:cron:delete']);
const canDetail = hasAccessByCodes(['system:cron:detail']);
const canUpdate = hasAccessByCodes(['system:cron:update']);
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
      componentProps: { placeholder: '请输入任务名称', allowClear: true },
      fieldName: 'name',
      label: '任务名称',
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
      component: 'Input',
      componentProps: { placeholder: '请输入备注', allowClear: true },
      fieldName: 'remark',
      label: '备注',
    },
    {
      component: 'RangePicker',
      fieldName: 'lastRunAtRange',
      label: '上次执行时间',
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
const gridOptions: VxeGridProps<CronItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { field: 'name', title: '任务名称', sortable: true },
    { field: 'expression', title: 'Cron表达式' },
    { field: 'handler', title: '处理器名称' },
    { field: 'params', title: '参数', slots: { header: tooltipHeader('参数', 'JSON') } },
    { field: 'remark', title: '备注', sortable: true },
    { field: 'status', title: '状态', width: 120, slots: { default: 'status_cell' }, sortable: true },
    { field: 'lastResult', title: '上次执行结果' },
    ...(isPlatformSuperAdmin.value ? [
    { field: 'tenantName', title: '租户' },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    { field: 'merchantName', title: '商户' },
    ] : []),
    { field: 'lastRunAt', title: '上次执行时间', width: 180, formatter: 'formatDateTime' },
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
        if (params.lastRunAtRange && params.lastRunAtRange.length === 2) {
          params.lastRunAtStart = params.lastRunAtRange[0];
          params.lastRunAtEnd = params.lastRunAtRange[1];
        }
        delete params.lastRunAtRange;
        if (sorts && sorts.length > 0) {
          const sort = sorts[0];
          if (sort && sort.field && sort.order) {
            params.orderBy = resolveSortField(String(sort.field));
            params.orderDir = sort.order;
          }
        }
        const res = await getCronList(params as any);
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
  showSearchForm: false,
});

const importInputRef = ref<HTMLInputElement | null>(null);

async function initSearchOptions() {
  if (isPlatformSuperAdmin.value) {
  try {
    const tenantIDRes = await getTenantList({ pageNum: 1, pageSize: 500 });
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
    const merchantIDRes = await getMerchantList({ pageNum: 1, pageSize: 500 });
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
function handleView(row: CronItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: CronItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: CronItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该定时任务表吗？',
    okType: 'danger',
    async onOk() {
      await deleteCron(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

/** 批量删除 */
function handleBatchDelete() {
  const ids = getGridSelectedIds<CronItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要删除的数据');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${ids.length} 条定时任务表吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteCron(ids);
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
    if (params.lastRunAtRange && params.lastRunAtRange.length === 2) {
      params.lastRunAtStart = params.lastRunAtRange[0];
      params.lastRunAtEnd = params.lastRunAtRange[1];
    }
    delete params.lastRunAtRange;
    if (sorts.length > 0) {
      const sort = sorts[0];
      if (sort?.field && sort?.order) {
        params.orderBy = resolveSortField(String(sort.field));
        params.orderDir = sort.order;
      }
    }
    const blob = await exportCron(params);
    downloadFileFromBlob({ fileName: '定时任务表.csv', source: blob as Blob });
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
    const res = await importCron(formData);
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
    const blob = await downloadImportTemplateCron();
    downloadFileFromBlob({ fileName: '定时任务表导入模板.csv', source: blob as Blob });
  } catch {
    message.error('下载模板失败');
  }
}

/** 批量修改状态 */
function handleBatchUpdateStatus() {
  const ids = getGridSelectedIds<CronItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要修改的数据');
    return;
  }
  const rows = gridApi.grid.getCheckboxRecords() as CronItem[];
  Modal.confirm({
    title: '批量修改状态',
    content: `确定要将选中的 ${ids.length} 条数据的状态切换吗？`,
    async onOk() {
      const newStatus = rows[0]?.status === 1 ? 0 : 1;
      await batchUpdateCron({ ids, status: newStatus });
      message.success('批量修改成功');
      gridApi.reload();
    },
  });
}

const logData = ref<CronLogItem[]>([]);
const logTotal = ref(0);
const logVisible = ref(false);
const currentLogCronName = ref('');

async function handleTrigger(row: CronItem) {
  try {
    const res = await triggerCron(row.name);
    if (res?.success) {
      message.success(`执行成功: ${res.message}`);
    } else {
      message.error(`执行失败: ${res?.message || '未知错误'}`);
    }
    gridApi.reload();
  } catch {
    message.error('执行失败');
  }
}

async function handleViewLog(row: CronItem) {
  currentLogCronName.value = row.name;
  try {
    const res = await getCronLogList({ cronId: row.id, pageNum: 1, pageSize: 50 });
    logData.value = res?.list ?? [];
    logTotal.value = res?.total ?? 0;
  } catch {
    logData.value = [];
    logTotal.value = 0;
  }
  logVisible.value = true;
}

function getRowActions(row: CronItem): ActionMoreItem[] {
  return [
    {
      key: 'detail',
      label: '查看',
      onClick: () => handleView(row),
      visible: canDetail,
    },
    {
      key: 'edit',
      label: '编辑',
      onClick: () => handleEdit(row),
      visible: canUpdate,
    },
    {
      key: 'trigger',
      label: '手动执行',
      onClick: () => handleTrigger(row),
      visible: canUpdate,
    },
    {
      key: 'log',
      label: '执行日志',
      onClick: () => handleViewLog(row),
      visible: true,
    },
    {
      danger: true,
      key: 'delete',
      label: '删除',
      onClick: () => handleDelete(row),
      visible: canDelete,
    },
  ];
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
        <Button v-access:code="'system:cron:create'" class="mr-2" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'system:cron:batch-delete'" danger class="mr-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-access:code="'system:cron:export'" class="mr-2" @click="handleExport">导出</Button>
        <Button v-access:code="'system:cron:import'" class="mr-2" @click="handleImportTrigger">导入</Button>
        <Button v-access:code="'system:cron:import'" class="mr-2" @click="handleDownloadTemplate">模板下载</Button>
        <Button v-access:code="'system:cron:batch-update'" class="mr-2" @click="handleBatchUpdateStatus">批量修改状态</Button>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status)">
          {{ getEnumLabel(statusMap, row.status) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <ActionMore :actions="getRowActions(row)" />
      </template>
    </Grid>
	  </Page>
	  <Modal v-model:open="logVisible" title="执行日志 - {{ currentLogCronName }}" width="800px" :footer="null">
		    <div v-if="logData.length === 0" style="text-align:center;padding:40px;color:#999">暂无执行日志</div>
			    <VxeTable v-else :data="logData" size="small">
			      <VxeColumn field="startAt" title="开始时间" width="160">
			        <template #default="{ row }"> {{ row.startAt?.substring(0, 19) }} </template>
			      </VxeColumn>
			      <VxeColumn field="durationMs" title="耗时" width="80">
			        <template #default="{ row }"> {{ row.durationMs }}ms </template>
			      </VxeColumn>
			      <VxeColumn field="result" title="结果" width="80">
			        <template #default="{ row }">
			          <Tag :color="row.result === 'success' ? 'green' : row.result === 'fail' ? 'red' : 'orange'">
			            {{ row.result === 'success' ? '成功' : row.result === 'fail' ? '失败' : '执行中' }}
			          </Tag>
			        </template>
			      </VxeColumn>
			      <VxeColumn field="message" title="消息" :show-overflow="true" />
			    </VxeTable>
	  </Modal>
	</template>
