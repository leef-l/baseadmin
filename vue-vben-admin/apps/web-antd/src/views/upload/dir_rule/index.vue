<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { DirRuleItem, DirRuleStorageTypesValue } from '#/api/upload/dir_rule/types';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';

import { Button, message, Modal, Tag } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { batchDeleteDirRule, deleteDirRule, getDirRuleList } from '#/api/upload/dir_rule';
import { getGridSelectedIds } from '#/utils/grid-selection';

import FormModal from './modules/form.vue';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

/** 类别选项 */
const categoryOptions = [
  { label: '默认', value: 1 },
  { label: '类型', value: 2 },
  { label: '来源', value: 3 },
];

/** 类别映射 */
const categoryMap: Record<number, string> = {
  1: '默认',
  2: '类型',
  3: '来源',
};

const storageTypeMap: Record<number, string> = {
  1: '本地',
  2: 'OSS',
  3: 'COS',
};

const keepNameMap: Record<number, string> = {
  0: '否',
  1: '是',
};

/** 类别颜色 */
function getCategoryColor(val: number): string {
  const keys = [1, 2, 3];
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx === -1 ? 0 : idx % TAG_COLORS.length] ?? 'default';
}

function getStorageTypeColor(val: number): string {
  const keys = [1, 2, 3];
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx === -1 ? 0 : idx % TAG_COLORS.length] ?? 'default';
}

/** 状态选项 */
const statusOptions = [
  { label: '禁用', value: 0 },
  { label: '启用', value: 1 },
];

/** 状态映射 */
const statusMap: Record<number, string> = {
  0: '禁用',
  1: '启用',
};

function getStorageTypes(value?: DirRuleStorageTypesValue): number[] {
  if (!value) {
    return [];
  }
  const rawItems = Array.isArray(value) ? value : String(value).split(/[,\s，；;]+/g);
  return rawItems
    .map((item) => Number(String(item).trim()))
    .filter((item) => item === 1 || item === 2 || item === 3);
}

/** 状态颜色 */
function getStatusColor(val: number): string {
  const keys = [0, 1];
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx === -1 ? 0 : idx % TAG_COLORS.length] ?? 'default';
}

/** 表单弹窗 */
const [FormModalComp, formModalApi] = useVbenModal({
  connectedComponent: FormModal,
  destroyOnClose: true,
});
const { hasAccessByCodes } = useAccess();
const canBatchDelete = hasAccessByCodes(['upload:dir_rule:batch-delete']);
/** 搜索表单配置 */
const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: true,
  submitOnChange: false,
  submitOnEnter: true,
  schema: [
    {
      component: 'Input',
      componentProps: {
        allowClear: true,
        placeholder: '请输入保存目录、匹配条件或适用存储',
      },
      fieldName: 'keyword',
      label: '关键词',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: categoryOptions,
        placeholder: '请选择类别',
        class: 'w-full',
      },
      fieldName: 'category',
      label: '类别',
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
  ],
};

/** 表格列配置 */
const gridOptions: VxeGridProps<DirRuleItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { field: 'dirName', title: '所属目录' },
    { field: 'category', title: '类别', width: 120, slots: { default: 'category_cell' } },
    { field: 'fileType', title: '匹配条件', width: 220 },
    { field: 'storageTypes', title: '使用存储', width: 180, slots: { default: 'storage_types_cell' } },
    { field: 'savePath', title: '保存目录' },
    { field: 'keepName', title: '保留原名', width: 110, slots: { default: 'keep_name_cell' } },
    { field: 'status', title: '状态', width: 120, slots: { default: 'status_cell' } },
    { field: 'createdAt', title: '创建时间', width: 180, formatter: 'formatDateTime' },
    { title: '操作', width: 200, fixed: 'right', slots: { default: 'action' } },
  ],
  height: 'auto',
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getDirRuleList({
          pageNum: page.currentPage,
          pageSize: page.pageSize,
          ...formValues,
        });
        return { items: res?.list ?? [], total: res?.total ?? 0 };
      },
    },
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

/** 新建 */
function handleCreate() {
  formModalApi.setData(null).open();
}

/** 编辑 */
function handleEdit(row: DirRuleItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: DirRuleItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该文件目录规则吗？',
    okType: 'danger',
    async onOk() {
      await deleteDirRule(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

function getSelectedIds() {
  return getGridSelectedIds<DirRuleItem>(gridApi.grid as any);
}

function handleBatchDelete() {
  const ids = getSelectedIds();
  if (ids.length === 0) {
    message.warning('请选择要删除的目录规则');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${ids.length} 条目录规则吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteDirRule(ids);
      message.success('批量删除成功');
      gridApi.reload();
    },
  });
}
</script>

<template>
  <Page auto-content-height>
    <FormModalComp @success="() => gridApi.reload()" />
    <Grid>
      <template #toolbar-actions>
        <Button v-access:code="'upload:dir_rule:create'" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'upload:dir_rule:batch-delete'" danger @click="handleBatchDelete">批量删除</Button>
      </template>
      <template #category_cell="{ row }">
        <Tag :color="getCategoryColor(row.category ?? 1)">
          {{ categoryMap[row.category ?? 1] || row.category }}
        </Tag>
      </template>
      <template #storage_types_cell="{ row }">
        <template v-if="getStorageTypes(row.storageTypes).length > 0">
          <Tag
            v-for="storageType in getStorageTypes(row.storageTypes)"
            :key="storageType"
            :color="getStorageTypeColor(storageType)"
          >
            {{ storageTypeMap[storageType] || storageType }}
          </Tag>
        </template>
        <span v-else>-</span>
      </template>
      <template #keep_name_cell="{ row }">
        <Tag :color="getStatusColor(row.keepName ?? 0)">
          {{ keepNameMap[row.keepName ?? 0] || row.keepName }}
        </Tag>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status ?? 0)">
          {{ statusMap[row.status ?? 0] || row.status }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-access:code="'upload:dir_rule:update'" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-access:code="'upload:dir_rule:delete'" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
