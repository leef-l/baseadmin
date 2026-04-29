<script setup lang="ts">
import { h, onMounted } from 'vue';
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
import { getUserTree, deleteUser, batchDeleteUser, exportUser, batchUpdateUser } from '#/api/member/user';
import { getLevelList } from '#/api/member/level';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
import type { UserItem } from '#/api/member/user/types';
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
  sort: 'sort',
  status: 'status',
  username: 'username',
  nickname: 'nickname',
  phone: 'phone',
  realName: 'real_name',
  remark: 'remark',
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}

/** 是否激活选项 */
const isActiveOptions = [
  { label: '未激活', value: 0 },
  { label: '已激活', value: 1 },
];

/** 是否激活映射 */
const isActiveMap: Record<EnumValue, string> = {
  0: '未激活',
  1: '已激活',
};

/** 是否激活颜色 */
function getIsActiveColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 仓库资格选项 */
const isQualifiedOptions = [
  { label: '已失效', value: 0 },
  { label: '有效', value: 1 },
];

/** 仓库资格映射 */
const isQualifiedMap: Record<EnumValue, string> = {
  0: '已失效',
  1: '有效',
};

/** 仓库资格颜色 */
function getIsQualifiedColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态选项 */
const statusOptions = [
  { label: '冻结', value: 0 },
  { label: '正常', value: 1 },
];

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '冻结',
  1: '正常',
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
const canBatchDelete = hasAccessByCodes(['member:user:batch-delete']);
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
      componentProps: { placeholder: '请输入用户名', allowClear: true },
      fieldName: 'username',
      label: '用户名',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入邀请码', allowClear: true },
      fieldName: 'inviteCode',
      label: '邀请码',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入昵称', allowClear: true },
      fieldName: 'nickname',
      label: '昵称',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入真实姓名', allowClear: true },
      fieldName: 'realName',
      label: '真实姓名',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '请输入手机号', allowClear: true },
      fieldName: 'phone',
      label: '手机号',
    },
    {
      component: 'TreeSelect',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'username', value: 'id', children: 'children' },
        placeholder: '请选择上级会员',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
      fieldName: 'parentID',
      label: '上级会员',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [],
        placeholder: '请选择当前等级',
        class: 'w-full',
      },
      fieldName: 'levelID',
      label: '当前等级',
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
        options: isActiveOptions,
        placeholder: '请选择是否激活',
        class: 'w-full',
      },
      fieldName: 'isActive',
      label: '是否激活',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: isQualifiedOptions,
        placeholder: '请选择仓库资格',
        class: 'w-full',
      },
      fieldName: 'isQualified',
      label: '仓库资格',
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
      fieldName: 'levelExpireAtRange',
      label: '等级到期时间',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        valueFormat: 'YYYY-MM-DD HH:mm:ss',
        class: 'w-full',
      },
    },
    {
      component: 'RangePicker',
      fieldName: 'lastLoginAtRange',
      label: '最后登录时间',
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
const gridOptions: VxeGridProps<UserItem> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
    { field: 'username', title: '用户名', slots: { header: tooltipHeader('用户名', '登录账号') }, treeNode: true },
    { field: 'nickname', title: '昵称' },
    { field: 'phone', title: '手机号' },
    { field: 'avatar', title: '头像', width: 100, slots: { default: 'avatar_cell' } },
    { field: 'realName', title: '真实姓名' },
    { field: 'levelName', title: '当前等级' },
    { field: 'teamCount', title: '团队总人数' },
    { field: 'directCount', title: '直推人数' },
    { field: 'activeCount', title: '有效用户数' },
    { field: 'teamTurnover', title: '团队总营业额', slots: { header: tooltipHeader('团队总营业额', '分') }, width: 140, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-', sortable: true },
    { field: 'isActive', title: '是否激活', width: 120, slots: { default: 'isActive_cell' } },
    { field: 'isQualified', title: '仓库资格', width: 120, slots: { default: 'isQualified_cell' } },
    { field: 'inviteCode', title: '邀请码' },
    { field: 'registerIP', title: '注册IP' },
    { field: 'remark', title: '备注' },
    { field: 'sort', title: '排序', slots: { header: tooltipHeader('排序', '升序') } },
    { field: 'status', title: '状态', width: 120, slots: { default: 'status_cell' } },
    ...(isPlatformSuperAdmin.value ? [
    { field: 'tenantName', title: '租户' },
    ] : []),
    ...(isPlatformSuperAdmin.value ? [
    { field: 'merchantName', title: '商户' },
    ] : []),
    { field: 'levelExpireAt', title: '等级到期时间', width: 180, formatter: 'formatDateTime' },
    { field: 'lastLoginAt', title: '最后登录时间', width: 180, formatter: 'formatDateTime' },
    { field: 'createdAt', title: '创建时间', width: 180, formatter: 'formatDateTime' },
    { title: '操作', width: 240, fixed: 'right', slots: { default: 'action' } },
  ],
  height: 'auto',
  pagerConfig: { enabled: false },
  treeConfig: {
    childrenField: 'children',
    expandAll: false,
  },
  proxyConfig: {
    ajax: {
      query: async (_params, formValues) => {
        const params: Record<string, any> = { ...formValues };
        if (params.timeRange && params.timeRange.length === 2) {
          params.startTime = params.timeRange[0];
          params.endTime = params.timeRange[1];
        }
        delete params.timeRange;
        if (!isPlatformSuperAdmin.value) {
          delete params.tenantID;
          delete params.merchantID;
        }
        if (params.levelExpireAtRange && params.levelExpireAtRange.length === 2) {
          params.levelExpireAtStart = params.levelExpireAtRange[0];
          params.levelExpireAtEnd = params.levelExpireAtRange[1];
        }
        delete params.levelExpireAtRange;
        if (params.lastLoginAtRange && params.lastLoginAtRange.length === 2) {
          params.lastLoginAtStart = params.lastLoginAtRange[0];
          params.lastLoginAtEnd = params.lastLoginAtRange[1];
        }
        delete params.lastLoginAtRange;
        return await getUserTree(params as any) ?? [];
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

async function initSearchOptions() {
  try {
    const treeRes = await getUserTree();
    gridApi.formApi.updateSchema([
      {
        fieldName: 'parentID',
        componentProps: {
          treeData: [
            { id: '0', username: '顶级节点', children: treeRes ?? [] },
          ],
        },
      },
    ]);
  } catch {
    // ignore
  }
  try {
    const levelIDRes = await getLevelList({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: 'levelID',
        componentProps: {
          options: (levelIDRes?.list ?? []).map((item: any) => ({
            label: item.name || item.id,
            value: item.id,
          })),
        },
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
function handleView(row: UserItem) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: UserItem) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: UserItem) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该会员用户吗？',
    okType: 'danger',
    async onOk() {
      await deleteUser(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

/** 批量删除 */
function handleBatchDelete() {
  const ids = getGridSelectedIds<UserItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要删除的数据');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${ids.length} 条会员用户吗？`,
    okType: 'danger',
    async onOk() {
      await batchDeleteUser(ids);
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
    if (params.timeRange && params.timeRange.length === 2) {
      params.startTime = params.timeRange[0];
      params.endTime = params.timeRange[1];
    }
    delete params.timeRange;
    if (!isPlatformSuperAdmin.value) {
      delete params.tenantID;
      delete params.merchantID;
    }
    if (params.levelExpireAtRange && params.levelExpireAtRange.length === 2) {
      params.levelExpireAtStart = params.levelExpireAtRange[0];
      params.levelExpireAtEnd = params.levelExpireAtRange[1];
    }
    delete params.levelExpireAtRange;
    if (params.lastLoginAtRange && params.lastLoginAtRange.length === 2) {
      params.lastLoginAtStart = params.lastLoginAtRange[0];
      params.lastLoginAtEnd = params.lastLoginAtRange[1];
    }
    delete params.lastLoginAtRange;
    const blob = await exportUser(params);
    downloadFileFromBlob({ fileName: '会员用户.csv', source: blob as Blob });
    message.success('导出成功');
  } catch {
    message.error('导出失败');
  }
}

/** 批量修改状态 */
function handleBatchUpdateStatus() {
  const ids = getGridSelectedIds<UserItem>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要修改的数据');
    return;
  }
  const rows = gridApi.grid.getCheckboxRecords() as UserItem[];
  Modal.confirm({
    title: '批量修改状态',
    content: `确定要将选中的 ${ids.length} 条数据的状态切换吗？`,
    async onOk() {
      const newStatus = rows[0]?.status === 1 ? 0 : 1;
      await batchUpdateUser({ ids, status: newStatus });
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
    <Grid>
      <template #toolbar-actions>
        <Button v-access:code="'member:user:create'" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'member:user:batch-delete'" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-access:code="'member:user:export'" class="ml-2" @click="handleExport">导出</Button>
        <Button v-access:code="'member:user:batch-update'" class="ml-2" @click="handleBatchUpdateStatus">批量修改状态</Button>
      </template>
      <template #avatar_cell="{ row }">
        <img v-if="row.avatar && /^https?:\/\//i.test(row.avatar)" :src="row.avatar" style="width: 48px; height: 48px; object-fit: cover; border-radius: 4px;" />
        <span v-else>-</span>
      </template>
      <template #isActive_cell="{ row }">
        <Tag :color="getIsActiveColor(row.isActive)">
          {{ getEnumLabel(isActiveMap, row.isActive) }}
        </Tag>
      </template>
      <template #isQualified_cell="{ row }">
        <Tag :color="getIsQualifiedColor(row.isQualified)">
          {{ getEnumLabel(isQualifiedMap, row.isQualified) }}
        </Tag>
      </template>
      <template #status_cell="{ row }">
        <Tag :color="getStatusColor(row.status)">
          {{ getEnumLabel(statusMap, row.status) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button v-access:code="'member:user:detail'" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-access:code="'member:user:update'" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-access:code="'member:user:delete'" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
