<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getShopOrderDetail,
  createShopOrder,
  updateShopOrder,
} from '#/api/member/shop_order';
import type {
  ShopOrderCreateParams,
  ShopOrderUpdateParams
} from '#/api/member/shop_order/types';
import { getUserTree } from '#/api/member/user';
import type { UserItem } from '#/api/member/user/types';
import { getShopGoodsList } from '#/api/member/shop_goods';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 支付钱包选项 */
const payWalletOptions = [
  { label: '优惠券余额', value: 1 },
];

/** 订单状态选项 */
const orderStatusOptions = [
  { label: '已完成', value: 1 },
  { label: '已取消', value: 2 },
];
const userIDOptions = ref<UserItem[]>([]);
const goodsIDOptions = ref<{ label: string; value: string | number }[]>([]);
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
      fieldName: 'orderNo',
      label: '订单号',
      rules: 'required',
      componentProps: { placeholder: '请输入订单号', maxlength: 64 },
    },
    {
      component: 'TreeSelect',
      fieldName: 'userID',
      label: '购买会员',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'nickname', value: 'id', children: 'children' },
        placeholder: '请选择购买会员',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'Select',
      fieldName: 'goodsID',
      label: '商品',
      componentProps: { options: [], placeholder: '请选择商品', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'goodsTitle',
      label: tooltipLabel('商品名称', '快照'),
      componentProps: { placeholder: '请输入商品名称（快照）', maxlength: 200 },
    },
    {
      component: 'ImageUpload',
      fieldName: 'goodsCover',
      label: '商品封面（快照）',
      componentProps: { maxCount: 1 },
    },
    {
      component: 'InputNumber',
      fieldName: 'quantity',
      label: '购买数量',
      componentProps: { placeholder: '请输入购买数量', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'totalPrice',
      label: tooltipLabel('订单总价', '分'),
      componentProps: { placeholder: '请输入订单总价（分）', class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'payWallet',
      label: '支付钱包',
      componentProps: { options: payWalletOptions, placeholder: '请选择支付钱包', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'orderStatus',
      label: '订单状态',
      componentProps: { options: orderStatusOptions, placeholder: '请选择订单状态', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'remark',
      label: '订单备注',
      componentProps: { placeholder: '请输入订单备注', maxlength: 500 },
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
      | ShopOrderCreateParams
      | undefined;
    if (!values) return;
    if (values.totalPrice != null) {
      (values as any).totalPrice = Math.round(Number(values.totalPrice) * 100);
    }
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateShopOrder({ id: editId.value, ...values } as ShopOrderUpdateParams);
        message.success('更新成功');
      } else {
        await createShopOrder(values);
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
    // 加载购买会员树形数据
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
    // 加载商品选项
    try {
      const shopGoodsRes = await getShopGoodsList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      goodsIDOptions.value = (shopGoodsRes?.list ?? []).map((item: any) => ({
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
      modalApi.setState({ title: '编辑商城订单' });
      try {
        const detail = await getShopOrderDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          const formData = { ...detail };
          if (formData.totalPrice != null) {
            formData.totalPrice = formData.totalPrice / 100;
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
      modalApi.setState({ title: '新建商城订单' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
