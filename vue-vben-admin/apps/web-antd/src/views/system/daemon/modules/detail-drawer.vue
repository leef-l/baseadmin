<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getDaemonDetail } from '#/api/system/daemon';
import type { DaemonItem } from '#/api/system/daemon/types';

const isPlatformSuperAdmin = usePlatformSuperAdmin();
const detail = ref<DaemonItem | null>(null);
const openToken = ref(0);

function displayValue(value: null | number | string | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return value;
}

const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  footer: false,
  async onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      openToken.value += 1;
      detail.value = null;
      return;
    }

    const currentOpenToken = ++openToken.value;
    const data = modalApi.getData<{ id: string }>();
    if (data?.id) {
      modalApi.setState({ title: '守护进程配置表详情' });
      try {
        const res = await getDaemonDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        detail.value = res;
      } catch {
        if (currentOpenToken === openToken.value) {
          detail.value = null;
        }
      }
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Descriptions v-if="detail" bordered :column="1" size="small">
      <DescriptionsItem label="ID">{{ detail.id }}</DescriptionsItem>
      <DescriptionsItem label="显示名称">{{ displayValue(detail.name) }}</DescriptionsItem>
      <DescriptionsItem label="Supervisor进程名">{{ displayValue(detail.program) }}</DescriptionsItem>
      <DescriptionsItem label="启动命令">{{ displayValue(detail.command) }}</DescriptionsItem>
      <DescriptionsItem label="运行目录">{{ displayValue(detail.directory) }}</DescriptionsItem>
      <DescriptionsItem label="运行用户">{{ displayValue(detail.runUser) }}</DescriptionsItem>
      <DescriptionsItem label="进程数量">{{ displayValue(detail.numprocs) }}</DescriptionsItem>
      <DescriptionsItem label="启动优先级">{{ displayValue(detail.priority) }}</DescriptionsItem>
      <DescriptionsItem label="是否随Supervisor启动">{{ displayValue(detail.autostart) }}</DescriptionsItem>
      <DescriptionsItem label="异常退出是否自动重启">{{ displayValue(detail.autorestart) }}</DescriptionsItem>
      <DescriptionsItem label="启动稳定秒数">{{ displayValue(detail.startsecs) }}</DescriptionsItem>
      <DescriptionsItem label="启动重试次数">{{ displayValue(detail.startretries) }}</DescriptionsItem>
      <DescriptionsItem label="停止信号">{{ displayValue(detail.stopSignal) }}</DescriptionsItem>
      <DescriptionsItem label="环境变量，Supervisor environment格式">{{ displayValue(detail.environment) }}</DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
