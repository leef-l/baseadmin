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
import { getSurveyList, deleteSurvey, batchDeleteSurvey, exportSurvey, importSurvey, downloadImportTemplateSurvey, batchUpdateSurvey } from '#/api/demo/survey';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
import type { SurveyItem } from '#/api/demo/survey/types';
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
  surveyNo: 'survey_no',
  title: 'title',
  poster: 'poster',
  questionJSON: 'question_json',
  introContent: 'intro_content',
  publishAt: 'publish_at',
  expireAt: 'expire_at',
  isAnonymous: 'is_anonymous',
  status: 'status',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}

/** 是否匿名选项 */
const isAnonymousOptions = [
  { label: '否', value: 0 },
  { label: '是', value: 1 },
];

/** 是否匿名映射 */
const isAnonymousMap: Record<EnumValue, string> = {
  0: '否',
  1: '是',
};

/** 是否匿名颜色 */
function getIsAnonymousColor(val: EnumValue | null | undefined): string {
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
const canBatchDelete = hasAccessByCodes(['demo:survey:batch-delete']);
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
      componentProps: { placeholder: '请输入问卷编号', allowClear: true },
      fieldName: 'surveyNo',
      label: '问卷编号',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入问卷标题', allowClear: true },
      fieldName: 'title',
      label: '问卷标题',
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
        options: isAnonymousOptions,
        placeholder: '请选择是否匿名',
        class: 'w-full',
      },
      fieldName: 'isAnonymous',
      label: '是否匿名',
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
      fieldName: 'publishAtRange',
      label: '发布时间',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        valueFormat: 'YYYY-MM-DD HH:mm:ss',
        class: 'w-full',
      },
    },
    {
      component: 'RangePicker',
      fieldName: 'expireAtRange',
      label: '过期时间',
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
const gridOptions: VxeGridProps<SurveyItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { title: '序号', type: 'seq', width: 50 },
    { field: 'surveyNo', title: '问卷编号' },
    { field: 'title', title: '问卷标题' },
    { field: 'poster', title: '海报', width: 100, slots: { default: 'poster_cell' } },
    { field: 'isAnonymous', title: '是否匿名', width: 120, slots: { default: 'isAnonymous_cell' } },
    { field: 'status', title: '状态', width: 120, slots: { default: 'status_cell' } },
    ...(isPlatformSuperAdmin.value ? [
    { field: 'tenantName', title: '租户' },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    { field: 'merchantName', title: '商户' },
    ] : []),
    { field: 'publishAt', title: '发布时间', width: 180, formatter: 'formatDateTime' },
    { field: 'expireAt', title: '过期时间', width: 180, formatter: 'formatDateTime' },
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
        if (params.publishAtRange && params.publishAtRange.length === 2) {
          params.publishAtStart = params.publishAtRange[0];
          params.publishAtEnd = params.publishAtRange[1];
        }
        delete params.publishAtRange;
        if (params.expireAtRange && params.expireAtRange.length === 2) {
          params.expireAtStart = params.expireAtRange[0];
          params.expireAtEnd = params.expireAtRange[1];
        }
        delete params.expireAtRange;
        if (sorts && sorts.length > 0) {
          const sort = sorts[0];
          if (sort && sort.field && sort.order) {
            params.orderBy = resolveSortField(String(sort.field));
            params.orderDir = sort.order;
          }
        }
        const res = await getSurveyList(params as any);
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
function handleView(row: SurveyItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: SurveyItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: SurveyItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该体验问卷吗？',
    okType: 'danger',
    async onOk() {
      await deleteSurvey(row.id);
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
    content: `确定要删除选中的 ${rows.length} 条体验问卷吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteSurvey(rows.map((r: SurveyItem) => r.id));
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
    if (params.publishAtRange && params.publishAtRange.length === 2) {
      params.publishAtStart = params.publishAtRange[0];
      params.publishAtEnd = params.publishAtRange[1];
    }
    delete params.publishAtRange;
    if (params.expireAtRange && params.expireAtRange.length === 2) {
      params.expireAtStart = params.expireAtRange[0];
      params.expireAtEnd = params.expireAtRange[1];
    }
    delete params.expireAtRange;
    if (sorts.length > 0) {
      const sort = sorts[0];
      if (sort?.field && sort?.order) {
        params.orderBy = resolveSortField(String(sort.field));
        params.orderDir = sort.order;
      }
    }
    const blob = await exportSurvey(params);
    downloadFileFromBlob({ fileName: '体验问卷.csv', source: blob as Blob });
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
    const res = await importSurvey(formData);
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
    const blob = await downloadImportTemplateSurvey();
    downloadFileFromBlob({ fileName: '体验问卷导入模板.csv', source: blob as Blob });
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
      await batchUpdateSurvey({ ids: rows.map((r: SurveyItem) => r.id), status: newStatus });
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
        <Button v-auth="['demo:survey:create']" type="primary" @click="handleCreate">新建</Button>
        <Button v-auth="['demo:survey:batch-delete']" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-auth="['demo:survey:export']" class="ml-2" @click="handleExport">导出</Button>
        <Button v-auth="['demo:survey:import']" class="ml-2" @click="handleImportTrigger">导入</Button>
        <Button class="ml-2" @click="handleDownloadTemplate">模板下载</Button>
        <Button v-auth="['demo:survey:batch-update']" class="ml-2" @click="handleBatchUpdateStatus">批量修改状态</Button>
      </template>
      <template #poster_cell="{ row }">
        <img v-if="row.poster" :src="row.poster" style="width: 48px; height: 48px; object-fit: cover; border-radius: 4px;" />
        <span v-else>-</span>
      </template>
      <template #isAnonymous_cell="{ row }">
        <Tag :color="getIsAnonymousColor(row.isAnonymous)">
          {{ getEnumLabel(isAnonymousMap, row.isAnonymous) }}
        </Tag>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status)">
          {{ getEnumLabel(statusMap, row.status) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-auth="['demo:survey:detail']" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-auth="['demo:survey:update']" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-auth="['demo:survey:delete']" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
