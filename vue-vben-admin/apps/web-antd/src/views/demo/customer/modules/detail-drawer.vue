<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getCustomerDetail } from '#/api/demo/customer';
import type { CustomerItem } from '#/api/demo/customer/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 性别映射 */
const genderMap: Record<EnumValue, string> = {
  0: '未知',
  1: '男',
  2: '女',
};

/** 性别颜色 */
function getGenderColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1, 2];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 等级映射 */
const levelMap: Record<EnumValue, string> = {
  1: '普通',
  2: 'VIP',
  3: '付费',
  4: '冻结',
};

/** 等级颜色 */
function getLevelColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 来源映射 */
const sourceTypeMap: Record<EnumValue, string> = {
  1: '官网',
  2: '小程序',
  3: '线下',
  4: '导入',
};

/** 来源颜色 */
function getSourceTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 是否VIP映射 */
const isVipMap: Record<EnumValue, string> = {
  0: '否',
  1: '是',
};

/** 是否VIP颜色 */
function getIsVipColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '禁用',
  1: '启用',
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
const detail = ref<CustomerItem | null>(null);
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
      modalApi.setState({ title: '体验客户详情' });
      try {
        const res = await getCustomerDetail(data.id);
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
      <DescriptionsItem label="头像">
        <img v-if="detail.avatar && /^https?:\/\//i.test(detail.avatar)" :src="detail.avatar" style="max-width: 200px; max-height: 200px; object-fit: contain;" />
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="客户名称">{{ displayValue(detail.name) }}</DescriptionsItem>
      <DescriptionsItem label="客户编号">{{ displayValue(detail.customerNo) }}</DescriptionsItem>
      <DescriptionsItem label="联系电话">{{ displayValue(detail.phone) }}</DescriptionsItem>
      <DescriptionsItem label="邮箱">{{ displayValue(detail.email) }}</DescriptionsItem>
      <DescriptionsItem label="性别">
        <Tag :color="getGenderColor(detail.gender)">{{ getEnumLabel(genderMap, detail.gender) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="等级">
        <Tag :color="getLevelColor(detail.level)">{{ getEnumLabel(levelMap, detail.level) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="来源">
        <Tag :color="getSourceTypeColor(detail.sourceType)">{{ getEnumLabel(sourceTypeMap, detail.sourceType) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="是否VIP">
        <Tag :color="getIsVipColor(detail.isVip)">{{ getEnumLabel(isVipMap, detail.isVip) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="注册时间">{{ displayValue(detail.registeredAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
