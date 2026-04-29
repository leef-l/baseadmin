<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getUserDetail } from '#/api/member/user';
import type { UserItem } from '#/api/member/user/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

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

const isPlatformSuperAdmin = usePlatformSuperAdmin();
const detail = ref<UserItem | null>(null);
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
      modalApi.setState({ title: '会员用户详情' });
      try {
        const res = await getUserDetail(data.id);
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
      <DescriptionsItem label="上级会员">{{ detail.userUsername || '-' }}</DescriptionsItem>
      <DescriptionsItem label="用户名">{{ displayValue(detail.username) }}</DescriptionsItem>
      <DescriptionsItem label="昵称">{{ displayValue(detail.nickname) }}</DescriptionsItem>
      <DescriptionsItem label="手机号">{{ displayValue(detail.phone) }}</DescriptionsItem>
      <DescriptionsItem label="头像">
        <img v-if="detail.avatar && /^https?:\/\//i.test(detail.avatar)" :src="detail.avatar" style="max-width: 200px; max-height: 200px; object-fit: contain;" />
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="真实姓名">{{ displayValue(detail.realName) }}</DescriptionsItem>
      <DescriptionsItem label="当前等级">{{ detail.levelName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="团队总人数">{{ displayValue(detail.teamCount) }}</DescriptionsItem>
      <DescriptionsItem label="直推人数">{{ displayValue(detail.directCount) }}</DescriptionsItem>
      <DescriptionsItem label="有效用户数">{{ displayValue(detail.activeCount) }}</DescriptionsItem>
      <DescriptionsItem label="团队总营业额">{{ detail.teamTurnover != null ? (detail.teamTurnover / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="是否激活">
        <Tag :color="getIsActiveColor(detail.isActive)">{{ getEnumLabel(isActiveMap, detail.isActive) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="仓库资格">
        <Tag :color="getIsQualifiedColor(detail.isQualified)">{{ getEnumLabel(isQualifiedMap, detail.isQualified) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="邀请码">{{ displayValue(detail.inviteCode) }}</DescriptionsItem>
      <DescriptionsItem label="注册IP">{{ displayValue(detail.registerIP) }}</DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="排序">{{ displayValue(detail.sort) }}</DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="等级到期时间">{{ displayValue(detail.levelExpireAt) }}</DescriptionsItem>
      <DescriptionsItem label="最后登录时间">{{ displayValue(detail.lastLoginAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
