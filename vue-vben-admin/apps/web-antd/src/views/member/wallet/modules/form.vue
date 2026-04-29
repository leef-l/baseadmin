<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getWalletDetail,
  createWallet,
  updateWallet,
} from '#/api/member/wallet';
import type {
  WalletCreateParams,
  WalletUpdateParams
} from '#/api/member/wallet/types';
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
        fieldNames: { label: 'username', value: 'id', children: 'children' },
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
      component: 'InputNumber',
      fieldName: 'balance',
      label: tooltipLabel('当前余额', '分'),
      componentProps: { placeholder: '请输入当前余额（分）', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'totalIncome',
      label: tooltipLabel('累计收入', '分'),
      componentProps: { placeholder: '请输入累计收入（分）', class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'totalExpense',
      label: tooltipLabel('累计支出', '分'),
      componentProps: { placeholder: '请输入累计支出（分）' },
    },
    {
      component: 'InputNumber',
      fieldName: 'frozenAmount',
      label: tooltipLabel('冻结金额', '分'),
      componentProps: { placeholder: '请输入冻结金额（分）', class: 'w-full' },
    },
    {
      component: 'Switch',
      fieldName: 'status',
      label: '状态',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
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
      | WalletCreateParams
      | undefined;
    if (!values) return;
    if (values.balance != null) {
      (values as any).balance = Math.round(Number(values.balance) * 100);
    }
    if (values.totalIncome != null) {
      (values as any).totalIncome = Math.round(Number(values.totalIncome) * 100);
    }
    if (values.frozenAmount != null) {
      (values as any).frozenAmount = Math.round(Number(values.frozenAmount) * 100);
    }
    if (values.totalExpense != null) {
      (values as any).totalExpense = Math.round(Number(values.totalExpense) * 100);
    }
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateWallet({ id: editId.value, ...values } as WalletUpdateParams);
        message.success('更新成功');
      } else {
        await createWallet(values);
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
      modalApi.setState({ title: '编辑会员钱包' });
      try {
        const detail = await getWalletDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          const formData = { ...detail };
          if (formData.balance != null) {
            formData.balance = formData.balance / 100;
          }
          if (formData.totalIncome != null) {
            formData.totalIncome = formData.totalIncome / 100;
          }
          if (formData.frozenAmount != null) {
            formData.frozenAmount = formData.frozenAmount / 100;
          }
          if (formData.totalExpense != null) {
            formData.totalExpense = formData.totalExpense / 100;
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
      modalApi.setState({ title: '新建会员钱包' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
