<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page, useVbenModal } from '@vben/common-ui';
import { Button, message, Modal, Tag } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getConfigList, deleteConfig } from '#/api/upload/config';
import type { ConfigItem } from '#/api/upload/config/types';
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

/** 是否默认选项 */
const isDefaultOptions = [
  { label: '否', value: 0 },
  { label: '是', value: 1 },
];

/** 是否默认映射 */
const isDefaultMap: Record<number, string> = {
  0: '否',
  1: '是',
};

/** 是否默认颜色 */
function getIsDefaultColor(val: number): string {
  const keys = [0, 1];
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
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

/** 状态颜色 */
function getStatusColor(val: number): string {
  const keys = [0, 1];
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 表单弹窗 */
const [FormModalComp, formModalApi] = useVbenModal({
  connectedComponent: FormModal,
  destroyOnClose: true,
});

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
        placeholder: '请输入配置名/路径/桶名',
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
        options: isDefaultOptions,
        placeholder: '请选择是否默认',
        class: 'w-full',
      },
      fieldName: 'isDefault',
      label: '是否默认',
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
const gridOptions: VxeGridProps<ConfigItem> = {
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    { field: 'name', title: '配置名称' },
    { field: 'storage', title: '存储类型', width: 120, slots: { default: 'storage_cell' } },
    { field: 'isDefault', title: '是否默认', width: 100, slots: { default: 'isDefault_cell' } },
    { field: 'maxSize', title: '大小限制(MB)', width: 120 },
    { field: 'status', title: '状态', width: 100, slots: { default: 'status_cell' } },
    { field: 'createdAt', title: '创建时间', width: 180, formatter: 'formatDateTime' },
    { title: '操作', width: 200, fixed: 'right', slots: { default: 'action' } },
  ],
  height: 'auto',
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getConfigList({
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
function handleEdit(row: ConfigItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: ConfigItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该上传配置吗？',
    okType: 'danger',
    async onOk() {
      await deleteConfig(row.id);
      message.success('删除成功');
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
        <Button type="primary" @click="handleCreate">新建</Button>
      </template>
      <template #storage_cell="{ row }">
        <Tag :color="getStorageColor(row.storage)">
          {{ storageMap[row.storage] || row.storage }}
        </Tag>
      </template>
      <template #isDefault_cell="{ row }">
        <Tag :color="getIsDefaultColor(row.isDefault)">
          {{ isDefaultMap[row.isDefault] || row.isDefault }}
        </Tag>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status)">
          {{ statusMap[row.status] || row.status }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
