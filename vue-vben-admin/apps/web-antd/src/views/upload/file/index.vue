<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { Button, message, Modal, Tag } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { batchDeleteFile, deleteFile, getFileList } from '#/api/upload/file';
import type { FileItem } from '#/api/upload/file/types';
import { getGridSelectedIds } from '#/utils/grid-selection';
import FormModal from './modules/form.vue';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

/** 存储类型选项 */
const storageOptions = [
  { label: '本地', value: 1 },
  { label: '阿里云OSS', value: 2 },
  { label: '腾讯云COS', value: 3 },
];

/** 存储类型映射 */
const storageMap: Record<number, string> = {
  1: '本地',
  2: '阿里云OSS',
  3: '腾讯云COS',
};

/** 存储类型颜色 */
function getStorageColor(val: number): string {
  const keys = [1, 2, 3];
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 是否图片选项 */
const isImageOptions = [
  { label: '否', value: 0 },
  { label: '是', value: 1 },
];

/** 是否图片映射 */
const isImageMap: Record<number, string> = {
  0: '否',
  1: '是',
};

/** 是否图片颜色 */
function getIsImageColor(val: number): string {
  const keys = [0, 1];
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 表单弹窗 */
const [FormModalComp, formModalApi] = useVbenModal({
  connectedComponent: FormModal,
  destroyOnClose: true,
});
const { hasAccessByCodes } = useAccess();
const canBatchDelete = hasAccessByCodes(['upload:file:batch-delete']);
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
        placeholder: '请输入文件名/地址/类型',
      },
      fieldName: 'keyword',
      label: '关键词',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: storageOptions,
        placeholder: '请选择存储类型',
        class: 'w-full',
      },
      fieldName: 'storage',
      label: '存储类型',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: isImageOptions,
        placeholder: '请选择是否图片',
        class: 'w-full',
      },
      fieldName: 'isImage',
      label: '是否图片',
    },
  ],
};

/** 表格列配置 */
const gridOptions: VxeGridProps<FileItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { field: 'dirName', title: '所属目录' },
    { field: 'name', title: '文件名称' },
    { field: 'url', title: '文件地址' },
    { field: 'ext', title: '文件扩展名' },
    { field: 'size', title: '文件大小' },
    { field: 'mime', title: 'MIME类型' },
    { field: 'storage', title: '存储类型', width: 120, slots: { default: 'storage_cell' } },
    { field: 'isImage', title: '是否图片', width: 120, slots: { default: 'isImage_cell' } },
    { field: 'createdAt', title: '创建时间', width: 180, formatter: 'formatDateTime' },
    { title: '操作', width: 200, fixed: 'right', slots: { default: 'action' } },
  ],
  height: 'auto',
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getFileList({
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
function handleEdit(row: FileItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: FileItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该文件记录吗？',
    okType: 'danger',
    async onOk() {
      await deleteFile(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

function getSelectedIds() {
  return getGridSelectedIds<FileItem>(gridApi.grid as any);
}

function handleBatchDelete() {
  const ids = getSelectedIds();
  if (ids.length === 0) {
    message.warning('请选择要删除的文件');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${ids.length} 个文件记录吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteFile(ids);
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
        <Button v-access:code="'upload:file:create'" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'upload:file:batch-delete'" danger @click="handleBatchDelete">批量删除</Button>
      </template>
      <template #storage_cell="{ row }">
        <Tag :color="getStorageColor(row.storage ?? 1)">
          {{ storageMap[row.storage ?? 1] || row.storage }}
        </Tag>
      </template>
      <template #isImage_cell="{ row }">
        <Tag :color="getIsImageColor(row.isImage ?? 0)">
          {{ isImageMap[row.isImage ?? 0] || row.isImage }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-access:code="'upload:file:update'" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-access:code="'upload:file:delete'" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
