<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { getAuditLogDetail } from '#/api/demo/audit_log';
import type { AuditLogItem } from '#/api/demo/audit_log/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 动作映射 */
const actionMap: Record<EnumValue, string> = {
  1: '创建',
  2: '修改',
  3: '删除',
  4: '导出',
  5: '导入',
};

/** 动作颜色 */
function getActionColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4, 5];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 对象类型映射 */
const targetTypeMap: Record<EnumValue, string> = {
  1: '客户',
  2: '商品',
  3: '订单',
  4: '工单',
};

/** 对象类型颜色 */
function getTargetTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 结果映射 */
const resultMap: Record<EnumValue, string> = {
  0: '失败',
  1: '成功',
};

/** 结果颜色 */
function getResultColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

const detail = ref<AuditLogItem | null>(null);
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
      modalApi.setState({ title: '体验审计日志详情' });
      try {
        const res = await getAuditLogDetail(data.id);
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
  <Modal class="w-[800px]">
    <Descriptions v-if="detail" bordered :column="1" size="small">
      <DescriptionsItem label="ID">{{ detail.id }}</DescriptionsItem>
      <DescriptionsItem label="日志编号">{{ displayValue(detail.logNo) }}</DescriptionsItem>
      <DescriptionsItem label="操作人">{{ detail.usersUsername || '-' }}</DescriptionsItem>
      <DescriptionsItem label="动作">
        <Tag :color="getActionColor(detail.action)">{{ getEnumLabel(actionMap, detail.action) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="对象类型">
        <Tag :color="getTargetTypeColor(detail.targetType)">{{ getEnumLabel(targetTypeMap, detail.targetType) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="对象编号">{{ displayValue(detail.targetCode) }}</DescriptionsItem>
      <DescriptionsItem label="请求JSON">
        <pre style="max-height: 300px; overflow: auto; white-space: pre-wrap; word-break: break-all; margin: 0; font-size: 12px;">{{ (() => { const value = detail.requestJSON; if (!value) return '-'; try { return JSON.stringify(JSON.parse(value), null, 2) } catch { return value } })() }}</pre>
      </DescriptionsItem>
      <DescriptionsItem label="结果">
        <Tag :color="getResultColor(detail.result)">{{ getEnumLabel(resultMap, detail.result) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="客户端IP">{{ displayValue(detail.clientIP) }}</DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="发生时间">{{ displayValue(detail.occurredAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
