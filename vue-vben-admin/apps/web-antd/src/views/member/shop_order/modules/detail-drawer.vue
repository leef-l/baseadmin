<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getShopOrderDetail } from '#/api/member/shop_order';
import type { ShopOrderItem } from '#/api/member/shop_order/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 支付钱包映射 */
const payWalletMap: Record<EnumValue, string> = {
  1: '优惠券余额',
};

/** 支付钱包颜色 */
function getPayWalletColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 订单状态映射 */
const orderStatusMap: Record<EnumValue, string> = {
  1: '已完成',
  2: '已取消',
};

/** 订单状态颜色 */
function getOrderStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '关闭',
  1: '开启',
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

const isPlatformSuperAdmin = usePlatformSuperAdmin();
const detail = ref<ShopOrderItem | null>(null);
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
      modalApi.setState({ title: '商城订单详情' });
      try {
        const res = await getShopOrderDetail(data.id);
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
      <DescriptionsItem label="订单号">{{ displayValue(detail.orderNo) }}</DescriptionsItem>
      <DescriptionsItem label="购买会员">{{ detail.userNickname || '-' }}</DescriptionsItem>
      <DescriptionsItem label="商品">{{ detail.shopGoodsTitle || '-' }}</DescriptionsItem>
      <DescriptionsItem label="商品名称">{{ displayValue(detail.goodsTitle) }}</DescriptionsItem>
      <DescriptionsItem label="商品封面">
        <img v-if="detail.goodsCover && /^https?:\/\//i.test(detail.goodsCover)" :src="detail.goodsCover" style="max-width: 200px; max-height: 200px; object-fit: contain;" />
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="购买数量">{{ displayValue(detail.quantity) }}</DescriptionsItem>
      <DescriptionsItem label="订单总价">{{ detail.totalPrice != null ? (detail.totalPrice / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="支付钱包">
        <Tag :color="getPayWalletColor(detail.payWallet)">{{ getEnumLabel(payWalletMap, detail.payWallet) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="订单状态">
        <Tag :color="getOrderStatusColor(detail.orderStatus)">{{ getEnumLabel(orderStatusMap, detail.orderStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="订单备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
