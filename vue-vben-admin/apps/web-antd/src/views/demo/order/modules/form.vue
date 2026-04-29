<script setup lang="ts">
import { h, ref } from 'vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getOrderDetail,
  createOrder,
  updateOrder,
} from '#/api/demo/order';
import type {
  OrderCreateParams,
  OrderUpdateParams
} from '#/api/demo/order/types';
import { getCustomerList } from '#/api/demo/customer';
import { getProductList } from '#/api/demo/product';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 支付状态选项 */
const payStatusOptions = [
  { label: '待支付', value: 0 },
  { label: '已支付', value: 1 },
  { label: '已退款', value: 2 },
];

/** 发货状态选项 */
const deliverStatusOptions = [
  { label: '待发货', value: 0 },
  { label: '已发货', value: 1 },
  { label: '已签收', value: 2 },
];

/** 状态选项 */
const statusOptions = [
  { label: '待确认', value: 0 },
  { label: '已确认', value: 1 },
  { label: '已取消', value: 2 },
];
const customerIDOptions = ref<{ label: string; value: string | number }[]>([]);
const productIDOptions = ref<{ label: string; value: string | number }[]>([]);
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
      componentProps: { placeholder: '请输入订单号', maxlength: 50 },
    },
    {
      component: 'Select',
      fieldName: 'customerID',
      label: '客户',
      componentProps: { options: [], placeholder: '请选择客户', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'productID',
      label: '商品',
      componentProps: { options: [], placeholder: '请选择商品', allowClear: true, class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'quantity',
      label: '购买数量',
      componentProps: { placeholder: '请输入购买数量', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'amount',
      label: tooltipLabel('订单金额', '分'),
      componentProps: { placeholder: '请输入订单金额（分）', class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'payStatus',
      label: '支付状态',
      componentProps: { options: payStatusOptions, placeholder: '请选择支付状态', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'deliverStatus',
      label: '发货状态',
      componentProps: { options: deliverStatusOptions, placeholder: '请选择发货状态', allowClear: true, class: 'w-full' },
    },
    {
      component: 'DatePicker',
      fieldName: 'paidAt',
      label: '支付时间',
      componentProps: { showTime: true, placeholder: '请选择支付时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'DatePicker',
      fieldName: 'deliverAt',
      label: '发货时间',
      componentProps: { showTime: true, placeholder: '请选择发货时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Input',
      fieldName: 'receiverPhone',
      label: '收货电话',
      componentProps: { placeholder: '请输入收货电话', maxlength: 30 },
    },
    {
      component: 'Input',
      fieldName: 'address',
      label: '收货地址',
      componentProps: { placeholder: '请输入收货地址', maxlength: 255 },
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注', rows: 4, maxlength: 65535 },
    },
    {
      component: 'RadioGroup',
      fieldName: 'status',
      label: '状态',
      componentProps: { options: statusOptions },
    },
    {
      component: 'Select',
      ifShow: () => isPlatformSuperAdmin.value,
      fieldName: 'tenantID',
      label: '租户',
      componentProps: { options: [], placeholder: '请选择租户', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      ifShow: () => isPlatformSuperAdmin.value,
      fieldName: 'merchantID',
      label: '商户',
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
      | OrderCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateOrder({ id: editId.value, ...values } as OrderUpdateParams);
        message.success('更新成功');
      } else {
        await createOrder(values);
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
    // 加载客户选项
    try {
      const customerRes = await getCustomerList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      customerIDOptions.value = (customerRes?.list ?? []).map((item: any) => ({
        label: item.name || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'customerID',
          componentProps: { options: customerIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    // 加载商品选项
    try {
      const productRes = await getProductList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      productIDOptions.value = (productRes?.list ?? []).map((item: any) => ({
        label: item.skuNo || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'productID',
          componentProps: { options: productIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
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
    if (currentOpenToken !== openToken.value) {
      return;
    }
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑体验订单' });
      try {
        const detail = await getOrderDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          formApi.setValues(detail);
        }
      } catch {
        if (currentOpenToken === openToken.value) {
          message.error('获取详情失败');
        }
      }
    } else {
      isEdit.value = false;
      editId.value = '';
      modalApi.setState({ title: '新建体验订单' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
