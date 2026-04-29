<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { Button, message, Modal, Tag } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  applyDomainNginx,
  applyDomainSSL,
  batchDeleteDomain,
  deleteDomain,
  getDomainList,
} from '#/api/system/domain';
import type { DomainItem } from '#/api/system/domain/types';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getGridSelectedIds } from '#/utils/grid-selection';

import FormModal from './modules/form.vue';

const statusOptions = [
  { label: '关闭', value: 0 },
  { label: '开启', value: 1 },
];

const ownerOptions = [
  { label: '租户', value: 1 },
  { label: '商户', value: 2 },
];

const binaryMap: Record<number, string> = {
  0: '否',
  1: '是',
};

const ownerMap: Record<number, string> = {
  1: '租户',
  2: '商户',
};

const statusMap: Record<number, string> = {
  0: '关闭',
  1: '开启',
};

function getBinaryColor(value: number) {
  return value === 1 ? 'green' : 'orange';
}

function getStatusColor(value: number) {
  return value === 1 ? 'green' : 'red';
}

const { hasAccessByCodes } = useAccess();
const canBatchDelete = hasAccessByCodes(['system:domain:batch-delete']);
const isPlatformSuperAdmin = usePlatformSuperAdmin();

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
        placeholder: '请输入域名/应用/备注',
      },
      fieldName: 'keyword',
      label: '关键词',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        class: 'w-full',
        options: ownerOptions,
        placeholder: '请选择主体类型',
      },
      fieldName: 'ownerType',
      label: '主体类型',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        class: 'w-full',
        options: [{ label: '后台', value: 'admin' }],
        placeholder: '请选择应用',
      },
      fieldName: 'appCode',
      label: '应用',
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

const gridOptions: VxeGridProps<DomainItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { title: '序号', type: 'seq', width: 50 },
    { field: 'domain', title: '域名', minWidth: 220 },
    { field: 'ownerType', slots: { default: 'owner_cell' }, title: '主体', width: 90 },
    ...(isPlatformSuperAdmin.value
      ? [
          { field: 'tenantName', title: '租户', minWidth: 150 },
          { field: 'merchantName', title: '商户', minWidth: 150 },
        ]
      : []),
    { field: 'appCode', title: '应用', width: 90 },
    {
      field: 'verifyStatus',
      slots: { default: 'verify_cell' },
      title: '已校验',
      width: 90,
    },
    {
      field: 'sslStatus',
      slots: { default: 'ssl_cell' },
      title: 'SSL',
      width: 80,
    },
    {
      field: 'nginxStatus',
      slots: { default: 'nginx_cell' },
      title: 'Nginx',
      width: 90,
    },
    { field: 'status', slots: { default: 'status_cell' }, title: '状态', width: 90 },
    {
      field: 'createdAt',
      formatter: 'formatDateTime',
      title: '创建时间',
      width: 180,
    },
    { fixed: 'right', slots: { default: 'action' }, title: '操作', width: 320 },
  ],
  height: 'auto',
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getDomainList({
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

function handleEdit(row: DomainItem) {
  formModalApi.setData({ id: row.id }).open();
}

function handleDelete(row: DomainItem) {
  Modal.confirm({
    async onOk() {
      await deleteDomain(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
    content: '确定要删除该域名吗？',
    okType: 'danger',
    title: '确认删除',
  });
}

function handleApplyNginx(row: DomainItem) {
  Modal.confirm({
    async onOk() {
      const res = await applyDomainNginx(row.id);
      message.success(
        res?.sslStatus === 1 ? 'Nginx和SSL配置已应用' : 'Nginx配置已应用',
      );
      gridApi.reload();
    },
    content: '将为该域名生成宝塔 Nginx vhost 并重载 Nginx，确认继续？',
    title: '应用Nginx配置',
  });
}

function handleApplySSL(row: DomainItem) {
  Modal.confirm({
    async onOk() {
      const res = await applyDomainSSL(row.id);
      message.success(
        res?.certPath ? `SSL证书已申请并启用：${res.certPath}` : 'SSL证书已申请并启用',
      );
      gridApi.reload();
    },
    content:
      '将先应用80端口Nginx配置，再调用宝塔ACME申请SSL证书，成功后自动启用443，确认继续？',
    title: '申请SSL证书',
  });
}

function getSelectedIds() {
  return getGridSelectedIds<DomainItem>(gridApi.grid as any);
}

function handleBatchDelete() {
  const ids = getSelectedIds();
  if (ids.length === 0) {
    message.warning('请选择要删除的域名');
    return;
  }
  Modal.confirm({
    async onOk() {
      await batchDeleteDomain(ids);
      message.success('批量删除成功');
      gridApi.reload();
    },
    content: `确定要删除选中的 ${ids.length} 个域名吗？`,
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
          v-access:code="'system:domain:create'"
          type="primary"
          @click="handleCreate"
        >
          新建
        </Button>
        <Button
          v-access:code="'system:domain:batch-delete'"
          danger
          @click="handleBatchDelete"
        >
          批量删除
        </Button>
      </template>
      <template #owner_cell="{ row }">
        <Tag color="blue">{{ ownerMap[row.ownerType] || row.ownerType }}</Tag>
      </template>
      <template #verify_cell="{ row }">
        <Tag :color="getBinaryColor(row.verifyStatus ?? 0)">
          {{ binaryMap[row.verifyStatus ?? 0] }}
        </Tag>
      </template>
      <template #ssl_cell="{ row }">
        <Tag :color="getBinaryColor(row.sslStatus ?? 0)">
          {{ binaryMap[row.sslStatus ?? 0] }}
        </Tag>
      </template>
      <template #nginx_cell="{ row }">
        <Tag :color="getBinaryColor(row.nginxStatus ?? 0)">
          {{ binaryMap[row.nginxStatus ?? 0] }}
        </Tag>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status ?? 0)">
          {{ statusMap[row.status ?? 0] || row.status }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button
          v-access:code="'system:domain:ssl'"
          type="link"
          size="small"
          @click="handleApplySSL(row)"
        >
          申请SSL
        </Button>
        <Button
          v-access:code="'system:domain:apply'"
          type="link"
          size="small"
          @click="handleApplyNginx(row)"
        >
          应用Nginx
        </Button>
        <Button
          v-access:code="'system:domain:update'"
          type="link"
          size="small"
          @click="handleEdit(row)"
        >
          编辑
        </Button>
        <Button
          v-access:code="'system:domain:delete'"
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
