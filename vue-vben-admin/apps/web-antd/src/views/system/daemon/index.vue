<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { Button, message, Modal, Tag } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  batchDeleteDaemon,
  batchRestartDaemon,
  batchStopDaemon,
  deleteDaemon,
  getDaemonDetail,
  getDaemonList,
  getDaemonLog,
  restartDaemon,
  stopDaemon,
} from '#/api/system/daemon';
import type { DaemonItem } from '#/api/system/daemon/types';
import { getGridSelectedIds } from '#/utils/grid-selection';

import FormModal from './modules/form.vue';

const { hasAccessByCodes } = useAccess();
const canBatchDelete = hasAccessByCodes(['system:daemon:batch-delete']);
const canBatchRestart = hasAccessByCodes(['system:daemon:restart']);
const canBatchStop = hasAccessByCodes(['system:daemon:stop']);

const detailOpen = ref(false);
const detailLoading = ref(false);
const detail = ref<DaemonItem>();
const normalLog = ref('');
const errorLog = ref('');

const runtimeColorMap: Record<string, string> = {
  BACKOFF: 'orange',
  EXITED: 'default',
  FATAL: 'red',
  MISSING: 'default',
  RUNNING: 'green',
  STARTING: 'blue',
  STOPPED: 'orange',
  STOPPING: 'blue',
  UNKNOWN: 'default',
};

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
        placeholder: '请输入名称/进程名/命令/目录/备注',
      },
      fieldName: 'keyword',
      label: '关键词',
    },
    {
      component: 'Input',
      componentProps: {
        allowClear: true,
        placeholder: '请输入进程名',
      },
      fieldName: 'program',
      label: '进程名',
    },
  ],
};

const gridOptions: VxeGridProps<DaemonItem> = {
  checkboxConfig:
    canBatchDelete || canBatchRestart || canBatchStop ? { highlight: true } : undefined,
  columns: [
    ...(canBatchDelete || canBatchRestart || canBatchStop
      ? [{ type: 'checkbox', width: 50 }]
      : []),
    { title: '序号', type: 'seq', width: 50 },
    { field: 'name', title: '显示名称', minWidth: 140 },
    { field: 'program', title: '进程名', minWidth: 160 },
    {
      field: 'runStatus',
      slots: { default: 'runtime_cell' },
      title: '运行状态',
      width: 120,
    },
    { field: 'pid', title: 'PID', width: 90 },
    { field: 'uptime', title: '运行时长', minWidth: 140 },
    { field: 'directory', title: '运行目录', minWidth: 260 },
    { field: 'command', title: '启动命令', minWidth: 320 },
    { field: 'runUser', title: '用户', width: 90 },
    { field: 'numprocs', title: '数量', width: 80 },
    {
      field: 'createdAt',
      formatter: 'formatDateTime',
      title: '创建时间',
      width: 180,
    },
    { fixed: 'right', slots: { default: 'action' }, title: '操作', width: 260 },
  ],
  height: 'auto',
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getDaemonList({
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

function handleEdit(row: DaemonItem) {
  formModalApi.setData({ id: row.id }).open();
}

async function handleView(row: DaemonItem) {
  detailOpen.value = true;
  detailLoading.value = true;
  detail.value = undefined;
  normalLog.value = '';
  errorLog.value = '';
  try {
    const [info, normal, error] = await Promise.all([
      getDaemonDetail(row.id),
      getDaemonLog(row.id, 'normal'),
      getDaemonLog(row.id, 'error'),
    ]);
    detail.value = info;
    normalLog.value = normal?.content ?? '';
    errorLog.value = error?.content ?? '';
  } finally {
    detailLoading.value = false;
  }
}

function handleRestart(row: DaemonItem) {
  Modal.confirm({
    async onOk() {
      await restartDaemon(row.id);
      message.success('重启指令已执行');
      gridApi.reload();
    },
    content: `确定要重启 ${row.program} 吗？`,
    title: '确认重启',
  });
}

function handleStop(row: DaemonItem) {
  Modal.confirm({
    async onOk() {
      await stopDaemon(row.id);
      message.success('暂停指令已执行');
      gridApi.reload();
    },
    content: `确定要暂停 ${row.program} 吗？`,
    title: '确认暂停',
  });
}

function handleDelete(row: DaemonItem) {
  Modal.confirm({
    async onOk() {
      await deleteDaemon(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
    content: `删除后会移除宝塔 Supervisor 配置并清理日志，确定要删除 ${row.program} 吗？`,
    okType: 'danger',
    title: '确认删除',
  });
}

function getSelectedIds() {
  return getGridSelectedIds<DaemonItem>(gridApi.grid as any);
}

function handleBatchRestart() {
  const ids = getSelectedIds();
  if (ids.length === 0) {
    message.warning('请选择要重启的守护进程');
    return;
  }
  Modal.confirm({
    async onOk() {
      await batchRestartDaemon(ids);
      message.success('批量重启指令已执行');
      gridApi.reload();
    },
    content: `确定要重启选中的 ${ids.length} 个守护进程吗？`,
    title: '确认批量重启',
  });
}

function handleBatchStop() {
  const ids = getSelectedIds();
  if (ids.length === 0) {
    message.warning('请选择要暂停的守护进程');
    return;
  }
  Modal.confirm({
    async onOk() {
      await batchStopDaemon(ids);
      message.success('批量暂停指令已执行');
      gridApi.reload();
    },
    content: `确定要暂停选中的 ${ids.length} 个守护进程吗？`,
    title: '确认批量暂停',
  });
}

function handleBatchDelete() {
  const ids = getSelectedIds();
  if (ids.length === 0) {
    message.warning('请选择要删除的守护进程');
    return;
  }
  Modal.confirm({
    async onOk() {
      await batchDeleteDaemon(ids);
      message.success('批量删除成功');
      gridApi.reload();
    },
    content: `确定要删除选中的 ${ids.length} 个守护进程吗？`,
    okType: 'danger',
    title: '确认批量删除',
  });
}

function runtimeColor(status?: string) {
  return runtimeColorMap[status || 'UNKNOWN'] || 'default';
}
</script>

<template>
  <Page auto-content-height>
    <FormModalComp @success="() => gridApi.reload()" />
    <Grid>
      <template #toolbar-actions>
        <Button
          v-access:code="'system:daemon:create'"
          type="primary"
          @click="handleCreate"
        >
          新建
        </Button>
        <Button
          v-access:code="'system:daemon:restart'"
          @click="handleBatchRestart"
        >
          批量重启
        </Button>
        <Button v-access:code="'system:daemon:stop'" @click="handleBatchStop">
          批量暂停
        </Button>
        <Button
          v-access:code="'system:daemon:batch-delete'"
          danger
          @click="handleBatchDelete"
        >
          批量删除
        </Button>
      </template>
      <template #runtime_cell="{ row }">
        <Tag :color="runtimeColor(row.runStatus)">
          {{ row.statusText || row.runStatus || '未知' }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Button
          v-access:code="'system:daemon:view'"
          type="link"
          size="small"
          @click="handleView(row)"
        >
          查看
        </Button>
        <Button
          v-access:code="'system:daemon:update'"
          type="link"
          size="small"
          @click="handleEdit(row)"
        >
          编辑
        </Button>
        <Button
          v-access:code="'system:daemon:restart'"
          type="link"
          size="small"
          @click="handleRestart(row)"
        >
          重启
        </Button>
        <Button
          v-access:code="'system:daemon:stop'"
          type="link"
          size="small"
          @click="handleStop(row)"
        >
          暂停
        </Button>
        <Button
          v-access:code="'system:daemon:delete'"
          type="link"
          danger
          size="small"
          @click="handleDelete(row)"
        >
          删除
        </Button>
      </template>
    </Grid>
    <Modal
      v-model:open="detailOpen"
      :confirm-loading="detailLoading"
      :footer="null"
      title="守护进程详情"
      width="900px"
    >
      <div v-if="detail" class="daemon-detail">
        <div class="daemon-detail-grid">
          <span>显示名称</span><strong>{{ detail.name }}</strong>
          <span>进程名</span><strong>{{ detail.program }}</strong>
          <span>运行状态</span><strong>{{ detail.statusText || detail.runStatus }}</strong>
          <span>PID</span><strong>{{ detail.pid || '-' }}</strong>
          <span>运行目录</span><strong>{{ detail.directory }}</strong>
          <span>启动命令</span><strong>{{ detail.command }}</strong>
          <span>配置文件</span><strong>{{ detail.configPath }}</strong>
          <span>标准日志</span><strong>{{ detail.outLogPath }}</strong>
          <span>错误日志</span><strong>{{ detail.errLogPath }}</strong>
          <span>环境变量</span><strong>{{ detail.environment || '-' }}</strong>
          <span>备注</span><strong>{{ detail.remark || '-' }}</strong>
        </div>
        <div class="daemon-log-title">标准输出</div>
        <pre class="daemon-log">{{ normalLog || '暂无日志' }}</pre>
        <div class="daemon-log-title">错误输出</div>
        <pre class="daemon-log">{{ errorLog || '暂无日志' }}</pre>
      </div>
    </Modal>
  </Page>
</template>

<style scoped>
.daemon-detail {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.daemon-detail-grid {
  display: grid;
  grid-template-columns: 96px minmax(0, 1fr);
  gap: 8px 12px;
  font-size: 13px;
}

.daemon-detail-grid span {
  color: hsl(var(--muted-foreground));
}

.daemon-detail-grid strong {
  min-width: 0;
  overflow-wrap: anywhere;
  font-weight: 500;
}

.daemon-log-title {
  font-size: 13px;
  font-weight: 600;
}

.daemon-log {
  max-height: 220px;
  padding: 10px;
  overflow: auto;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  background: hsl(var(--muted));
  border-radius: 6px;
}
</style>
