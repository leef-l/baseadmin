<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { Button, message, Modal, Tag } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  batchDeleteTenant,
  deleteTenant,
  getTenantList,
} from '#/api/system/tenant';
import type { TenantItem } from '#/api/system/tenant/types';
import { getGridSelectedIds } from '#/utils/grid-selection';

import FormModal from './modules/form.vue';

const TAG_COLORS = ['green', 'red', 'blue', 'orange'];

const statusOptions = [
  { label: '关闭', value: 0 },
  { label: '开启', value: 1 },
];

const statusMap: Record<number, string> = {
  0: '关闭',
  1: '开启',
};

function getStatusColor(val: number): string {
  const keys = [0, 1];
  const index = keys.indexOf(val);
  return TAG_COLORS[index >= 0 ? index % TAG_COLORS.length : 0] ?? 'default';
}

const { hasAccessByCodes } = useAccess();
const canBatchDelete = hasAccessByCodes(['system:tenant:batch-delete']);

const [FormModalComp, formModalApi] = useVbenModal({
  connectedComponent: FormModal,
  destroyOnClose: true,
});

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
        placeholder: '请输入租户名称/编码/联系人/电话',
      },
      fieldName: 'keyword',
      label: '关键词',
    },
    {
      component: 'Input',
      componentProps: {
        allowClear: true,
        placeholder: '请输入租户编码',
      },
      fieldName: 'code',
      label: '租户编码',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        class: 'w-full',
        options: statusOptions,
        placeholder: '请选择状态',
      },
      fieldName: 'status',
      label: '状态',
    },
  ],
};

const gridOptions: VxeGridProps<TenantItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { title: '序号', type: 'seq', width: 50 },
    { field: 'name', title: '租户名称', minWidth: 150 },
    { field: 'code', title: '租户编码', minWidth: 140 },
    { field: 'contactName', title: '联系人', minWidth: 120 },
    { field: 'contactPhone', title: '联系电话', minWidth: 140 },
    { field: 'domain', title: '域名', minWidth: 180 },
    {
      field: 'expireAt',
      formatter: 'formatDateTime',
      title: '到期时间',
      width: 180,
    },
    { field: 'status', slots: { default: 'status_cell' }, title: '状态', width: 100 },
    { field: 'remark', title: '备注', minWidth: 180 },
    {
      field: 'createdAt',
      formatter: 'formatDateTime',
      title: '创建时间',
      width: 180,
    },
    { fixed: 'right', slots: { default: 'action' }, title: '操作', width: 150 },
  ],
  height: 'auto',
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getTenantList({
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

function handleCreate() {
  formModalApi.setData(null).open();
}

function handleEdit(row: TenantItem) {
  formModalApi.setData({ id: row.id }).open();
}

function handleDelete(row: TenantItem) {
  Modal.confirm({
    async onOk() {
      await deleteTenant(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
    content: '确定要删除该租户吗？',
    okType: 'danger',
    title: '确认删除',
  });
}

function getSelectedIds() {
  return getGridSelectedIds<TenantItem>(gridApi.grid as any);
}

function handleBatchDelete() {
  const ids = getSelectedIds();
  if (ids.length === 0) {
    message.warning('请选择要删除的租户');
    return;
  }
  Modal.confirm({
    async onOk() {
      await batchDeleteTenant(ids);
      message.success('批量删除成功');
      gridApi.reload();
    },
    content: `确定要删除选中的 ${ids.length} 个租户吗？`,
    okType: 'danger',
    title: '确认批量删除',
  });
}
</script>

<template>
  <Page auto-content-height>
    <FormModalComp @success="() => gridApi.reload()" />
    <Grid>
      <template #toolbar-actions>
        <Button
          v-access:code="'system:tenant:create'"
          type="primary"
          @click="handleCreate"
        >
          新建
        </Button>
        <Button
          v-access:code="'system:tenant:batch-delete'"
          danger
          @click="handleBatchDelete"
        >
          批量删除
        </Button>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status ?? 0)">
          {{ statusMap[row.status ?? 0] || row.status }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button
          v-access:code="'system:tenant:update'"
          type="link"
          size="small"
          @click="handleEdit(row)"
        >
          编辑
        </Button>
        <Button
          v-access:code="'system:tenant:delete'"
          type="link"
          danger
          size="small"
          @click="handleDelete(row)"
        >
          删除
        </Button>
      </template>
    </Grid>
  </Page>
</template>
