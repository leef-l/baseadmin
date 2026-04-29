<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getWalletLogDetail } from '#/api/member/wallet_log';
import type { WalletLogItem } from '#/api/member/wallet_log/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 钱包类型映射 */
const walletTypeMap: Record<EnumValue, string> = {
  1: '优惠券余额',
  2: '奖金余额',
  3: '推广奖余额',
};

/** 钱包类型颜色 */
function getWalletTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 变动类型映射 */
const changeTypeMap: Record<EnumValue, string> = {
  1: '充值',
  2: '消费',
  3: '推广奖',
  4: '仓库卖出收入',
  5: '平台扣除',
  6: '后台调整',
};

/** 变动类型颜色 */
function getChangeTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4, 5, 6];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

const isPlatformSuperAdmin = usePlatformSuperAdmin();
const detail = ref<WalletLogItem | null>(null);
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
      modalApi.setState({ title: '钱包流水记录详情' });
      try {
        const res = await getWalletLogDetail(data.id);
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
      <DescriptionsItem label="会员">{{ detail.userNickname || '-' }}</DescriptionsItem>
      <DescriptionsItem label="钱包类型">
        <Tag :color="getWalletTypeColor(detail.walletType)">{{ getEnumLabel(walletTypeMap, detail.walletType) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="变动类型">
        <Tag :color="getChangeTypeColor(detail.changeType)">{{ getEnumLabel(changeTypeMap, detail.changeType) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="变动金额">{{ detail.changeAmount != null ? (detail.changeAmount / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="变动前余额">{{ detail.beforeBalance != null ? (detail.beforeBalance / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="变动后余额">{{ detail.afterBalance != null ? (detail.afterBalance / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="关联单号">{{ displayValue(detail.relatedOrderNo) }}</DescriptionsItem>
      <DescriptionsItem label="备注说明">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
