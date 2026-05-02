<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getWarehouseTradeDetail,
  createWarehouseTrade,
  updateWarehouseTrade,
} from '#/api/member/warehouse_trade';
import type {
  WarehouseTradeCreateParams,
  WarehouseTradeUpdateParams
} from '#/api/member/warehouse_trade/types';
import { getWarehouseGoodsList } from '#/api/member/warehouse_goods';
import { getWarehouseListingList } from '#/api/member/warehouse_listing';
import { getUserTree } from '#/api/member/user';
import type { UserItem } from '#/api/member/user/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 交易状态选项 */
const tradeStatusOptions = [
  { label: '待卖家确认', value: 1 },
  { label: '已确认完成', value: 2 },
  { label: '已取消', value: 3 },
];
const goodsIDOptions = ref<{ label: string; value: string | number }[]>([]);
const listingIDOptions = ref<{ label: string; value: string | number }[]>([]);
const sellerIDOptions = ref<UserItem[]>([]);
const buyerIDOptions = ref<UserItem[]>([]);
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
      component: 'Input',
      fieldName: 'tradeNo',
      label: '交易编号',
      rules: 'required',
      componentProps: { placeholder: '请输入交易编号', maxlength: 64 },
    },
    {
      component: 'Select',
      fieldName: 'goodsID',
      label: '仓库商品',
      componentProps: { options: [], placeholder: '请选择仓库商品', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'listingID',
      label: '挂卖记录',
      componentProps: { options: [], placeholder: '请选择挂卖记录', allowClear: true, class: 'w-full' },
    },
    {
      component: 'TreeSelect',
      fieldName: 'sellerID',
      label: '卖家',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'nickname', value: 'id', children: 'children' },
        placeholder: '请选择卖家',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'TreeSelect',
      fieldName: 'buyerID',
      label: '买家',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'nickname', value: 'id', children: 'children' },
        placeholder: '请选择买家',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'InputNumber',
      fieldName: 'tradePrice',
      label: tooltipLabel('成交价格', '分'),
      componentProps: { placeholder: '请输入成交价格（分）', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'platformFee',
      label: tooltipLabel('平台扣除费用', '分'),
      componentProps: { placeholder: '请输入平台扣除费用（分）', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'sellerIncome',
      label: tooltipLabel('卖家实收', '分'),
      componentProps: { placeholder: '请输入卖家实收（分）', class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'tradeStatus',
      label: '交易状态',
      componentProps: { options: tradeStatusOptions, placeholder: '请选择交易状态', allowClear: true, class: 'w-full' },
    },
    {
      component: 'DatePicker',
      fieldName: 'confirmedAt',
      label: '确认时间',
      componentProps: { showTime: true, placeholder: '请选择确认时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Input',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注', maxlength: 500 },
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
      | WarehouseTradeCreateParams
      | undefined;
    if (!values) return;
    if (values.tradePrice != null) {
      (values as any).tradePrice = Math.round(Number(values.tradePrice) * 100);
    }
    if (values.platformFee != null) {
      (values as any).platformFee = Math.round(Number(values.platformFee) * 100);
    }
    if (values.sellerIncome != null) {
      (values as any).sellerIncome = Math.round(Number(values.sellerIncome) * 100);
    }
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateWarehouseTrade({ id: editId.value, ...values } as WarehouseTradeUpdateParams);
        message.success('更新成功');
      } else {
        await createWarehouseTrade(values);
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
    // 加载仓库商品选项
    try {
      const warehouseGoodsRes = await getWarehouseGoodsList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      goodsIDOptions.value = (warehouseGoodsRes?.list ?? []).map((item: any) => ({
        label: item.title || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'goodsID',
          componentProps: { options: goodsIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    // 加载挂卖记录选项
    try {
      const warehouseListingRes = await getWarehouseListingList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      listingIDOptions.value = (warehouseListingRes?.list ?? []).map((item: any) => ({
        label: item.id || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'listingID',
          componentProps: { options: listingIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    // 加载卖家树形数据
    try {
      const userRes = await getUserTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      sellerIDOptions.value = userRes ?? [];
      formApi.updateSchema([
        {
          fieldName: 'sellerID',
          componentProps: { treeData: sellerIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    // 加载买家树形数据
    try {
      const userRes = await getUserTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      buyerIDOptions.value = userRes ?? [];
      formApi.updateSchema([
        {
          fieldName: 'buyerID',
          componentProps: { treeData: buyerIDOptions.value },
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
      modalApi.setState({ title: '编辑仓库交易记录' });
      try {
        const detail = await getWarehouseTradeDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          const formData = { ...detail };
          if (formData.tradePrice != null) {
            formData.tradePrice = formData.tradePrice / 100;
          }
          if (formData.platformFee != null) {
            formData.platformFee = formData.platformFee / 100;
          }
          if (formData.sellerIncome != null) {
            formData.sellerIncome = formData.sellerIncome / 100;
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
      modalApi.setState({ title: '新建仓库交易记录' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
