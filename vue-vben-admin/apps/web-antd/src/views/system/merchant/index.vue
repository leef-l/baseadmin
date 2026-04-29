<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { Button, message, Modal, Tag } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  batchDeleteMerchant,
  deleteMerchant,
  getMerchantList,
} from '#/api/system/merchant';
import type { MerchantItem } from '#/api/system/merchant/types';
import { getTenantList } from '#/api/system/tenant';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getGridSelectedIds } from '#/utils/grid-selection';

import FormModal from './modules/form.vue';

interface SelectOption {
  label: string;
  value: string;
}

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

const tenantOptions = ref<SelectOption[]>([]);
const isPlatformSuperAdmin = usePlatformSuperAdmin();

async function loadTenantOptions() {
  try {
    const res = await getTenantList({ pageNum: 1, pageSize: 500 });
    tenantOptions.value = (res?.list ?? []).map((item) => ({
      label: `${item.name}（${item.code}）`,
      value: item.id,
    }));
  } catch {
    tenantOptions.value = [];
  }
}

onMounted(() => {
  if (isPlatformSuperAdmin.value) {
    loadTenantOptions();
  }
});

const { hasAccessByCodes } = useAccess();
const canBatchDelete = hasAccessByCodes(['system:merchant:batch-delete']);

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
        placeholder: '请输入商户名称/编码/联系人/电话',
      },
      fieldName: 'keyword',
      label: '关键词',
    },
    ...(isPlatformSuperAdmin.value
      ? [
          {
            component: 'Select',
            componentProps: () => ({
              allowClear: true,
              class: 'w-full',
              options: tenantOptions.value,
              placeholder: '请选择租户',
            }),
            fieldName: 'tenantId',
            label: '租户',
          },
        ]
      : []),
    {
      component: 'Input',
      componentProps: {
        allowClear: true,
        placeholder: '请输入商户编码',
      },
      fieldName: 'code',
      label: '商户编码',
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

const gridOptions: VxeGridProps<MerchantItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { title: '序号', type: 'seq', width: 50 },
    { field: 'name', title: '商户名称', minWidth: 150 },
    { field: 'code', title: '商户编码', minWidth: 140 },
    ...(isPlatformSuperAdmin.value
      ? [{ field: 'tenantName', title: '所属租户', minWidth: 150 }]
      : []),
    { field: 'contactName', title: '联系人', minWidth: 120 },
    { field: 'contactPhone', title: '联系电话', minWidth: 140 },
    { field: 'address', title: '地址', minWidth: 200 },
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
        const params: Record<string, any> = { ...formValues };
        if (!isPlatformSuperAdmin.value) {
          delete params.tenantId;
        }
        const res = await getMerchantList({
          pageNum: page.currentPage,
          pageSize: page.pageSize,
          ...params,
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

function handleEdit(row: MerchantItem) {
  formModalApi.setData({ id: row.id }).open();
}

function handleDelete(row: MerchantItem) {
  Modal.confirm({
    async onOk() {
      await deleteMerchant(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
    content: '确定要删除该商户吗？',
    okType: 'danger',
    title: '确认删除',
  });
}

function getSelectedIds() {
  return getGridSelectedIds<MerchantItem>(gridApi.grid as any);
}

function handleBatchDelete() {
  const ids = getSelectedIds();
  if (ids.length === 0) {
    message.warning('请选择要删除的商户');
    return;
  }
  Modal.confirm({
    async onOk() {
      await batchDeleteMerchant(ids);
      message.success('批量删除成功');
      gridApi.reload();
    },
    content: `确定要删除选中的 ${ids.length} 个商户吗？`,
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
          v-access:code="'system:merchant:create'"
          type="primary"
          @click="handleCreate"
        >
          新建
        </Button>
        <Button
          v-access:code="'system:merchant:batch-delete'"
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
          v-access:code="'system:merchant:update'"
          type="link"
          size="small"
          @click="handleEdit(row)"
        >
          编辑
        </Button>
        <Button
          v-access:code="'system:merchant:delete'"
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
