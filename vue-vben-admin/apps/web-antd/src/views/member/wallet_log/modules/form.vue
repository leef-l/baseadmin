<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getWalletLogDetail,
  createWalletLog,
  updateWalletLog,
} from '#/api/member/wallet_log';
import type {
  WalletLogCreateParams,
  WalletLogUpdateParams
} from '#/api/member/wallet_log/types';
import { getUserTree } from '#/api/member/user';
import type { UserItem } from '#/api/member/user/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 钱包类型选项 */
const walletTypeOptions = [
  { label: '优惠券余额', value: 1 },
  { label: '奖金余额', value: 2 },
  { label: '推广奖余额', value: 3 },
];

/** 变动类型选项 */
const changeTypeOptions = [
  { label: '充值', value: 1 },
  { label: '消费', value: 2 },
  { label: '推广奖', value: 3 },
  { label: '仓库卖出收入', value: 4 },
  { label: '平台扣除', value: 5 },
  { label: '后台调整', value: 6 },
];
const userIDOptions = ref<UserItem[]>([]);
const tenantIDOptions = ref<{ label: string; value: string | number }[]>([]);
const merchantIDOptions = ref<{ label: string; value: string | number }[]>([]);
/** 渲染带 Tooltip 的表单 label */
function tooltipLabel(label: string, tip: string) {
  return () => h('span', {}, [
    label + ' ',
    h(Tooltip, { title: tip }, {
      default: () => h(QuestionCircleOutlined, { style: { color: '#999', marginLeft: '4px' } }),
    }),
  ]);
}

const emit = defineEmits<{ success: [] }>();
const isPlatformSuperAdmin = usePlatformSuperAdmin();
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'TreeSelect',
      fieldName: 'userID',
      label: '会员',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'nickname', value: 'id', children: 'children' },
        placeholder: '请选择会员',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'Select',
      fieldName: 'walletType',
      label: '钱包类型',
      componentProps: { options: walletTypeOptions, placeholder: '请选择钱包类型', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'changeType',
      label: '变动类型',
      componentProps: { options: changeTypeOptions, placeholder: '请选择变动类型', allowClear: true, class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'changeAmount',
      label: tooltipLabel('变动金额', '分，正增负减'),
      componentProps: { placeholder: '请输入变动金额（分，正增负减）', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'beforeBalance',
      label: tooltipLabel('变动前余额', '分'),
      componentProps: { placeholder: '请输入变动前余额（分）', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'afterBalance',
      label: tooltipLabel('变动后余额', '分'),
      componentProps: { placeholder: '请输入变动后余额（分）', class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'relatedOrderNo',
      label: '关联单号',
      componentProps: { placeholder: '请输入关联单号', maxlength: 64 },
    },
    {
      component: 'Input',
      fieldName: 'remark',
      label: '备注说明',
      componentProps: { placeholder: '请输入备注说明', maxlength: 500 },
    },
    {
      component: 'Select',
      fieldName: 'tenantID',
      label: '租户',
      ifShow: () => isPlatformSuperAdmin.value,
      componentProps: { options: [], placeholder: '请选择租户', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'merchantID',
      label: '商户',
      ifShow: () => isPlatformSuperAdmin.value,
      componentProps: { options: [], placeholder: '请选择商户', allowClear: true, class: 'w-full' },
    },
  ],
});

/** Modal 配置 */
const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  onCancel() {
    modalApi.close();
  },
  onConfirm: async () => {
    const values = await formApi.validateAndSubmitForm() as
      | WalletLogCreateParams
      | undefined;
    if (!values) return;
    if (values.changeAmount != null) {
      (values as any).changeAmount = Math.round(Number(values.changeAmount) * 100);
    }
    if (values.beforeBalance != null) {
      (values as any).beforeBalance = Math.round(Number(values.beforeBalance) * 100);
    }
    if (values.afterBalance != null) {
      (values as any).afterBalance = Math.round(Number(values.afterBalance) * 100);
    }
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateWalletLog({ id: editId.value, ...values } as WalletLogUpdateParams);
        message.success('更新成功');
      } else {
        await createWalletLog(values);
        message.success('创建成功');
      }
      emit('success');
      modalApi.close();
    } finally {
      modalApi.lock(false);
    }
  },
  async onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      openToken.value += 1;
      return;
    }

    const currentOpenToken = ++openToken.value;
    formApi.resetForm();
    const data = modalApi.getData<{ id?: string }>();
    // 加载会员树形数据
    try {
      const userRes = await getUserTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      userIDOptions.value = userRes ?? [];
      formApi.updateSchema([
        {
          fieldName: 'userID',
          componentProps: { treeData: userIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    if (isPlatformSuperAdmin.value) {
    // 加载租户选项
    try {
      const tenantRes = await getTenantList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      tenantIDOptions.value = (tenantRes?.list ?? []).map((item: any) => ({
        label: item.name || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'tenantID',
          componentProps: { options: tenantIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    }
    if (isPlatformSuperAdmin.value) {
    // 加载商户选项
    try {
      const merchantRes = await getMerchantList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      merchantIDOptions.value = (merchantRes?.list ?? []).map((item: any) => ({
        label: item.name || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'merchantID',
          componentProps: { options: merchantIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    }
    if (currentOpenToken !== openToken.value) {
      return;
    }
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑钱包流水记录' });
      try {
        const detail = await getWalletLogDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          const formData = { ...detail };
          if (formData.changeAmount != null) {
            formData.changeAmount = formData.changeAmount / 100;
          }
          if (formData.beforeBalance != null) {
            formData.beforeBalance = formData.beforeBalance / 100;
          }
          if (formData.afterBalance != null) {
            formData.afterBalance = formData.afterBalance / 100;
          }
          formApi.setValues(formData);
        }
      } catch {
        if (currentOpenToken === openToken.value) {
          message.error('获取详情失败');
        }
      }
    } else {
      isEdit.value = false;
      editId.value = '';
      modalApi.setState({ title: '新建钱包流水记录' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
