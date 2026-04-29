<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getWarehouseGoodsDetail,
  createWarehouseGoods,
  updateWarehouseGoods,
} from '#/api/member/warehouse_goods';
import type {
  WarehouseGoodsCreateParams,
  WarehouseGoodsUpdateParams
} from '#/api/member/warehouse_goods/types';
import { getUserTree } from '#/api/member/user';
import type { UserItem } from '#/api/member/user/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 商品状态选项 */
const goodsStatusOptions = [
  { label: '持有中', value: 1 },
  { label: '挂卖中', value: 2 },
  { label: '交易中', value: 3 },
];
const ownerIDOptions = ref<UserItem[]>([]);
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
      fieldName: 'goodsNo',
      label: '商品编号',
      rules: 'required',
      componentProps: { placeholder: '请输入商品编号', maxlength: 64 },
    },
    {
      component: 'Input',
      fieldName: 'title',
      label: '商品名称',
      rules: 'required',
      componentProps: { placeholder: '请输入商品名称', maxlength: 200 },
    },
    {
      component: 'ImageUpload',
      fieldName: 'cover',
      label: '商品封面',
      componentProps: { maxCount: 1 },
    },
    {
      component: 'InputNumber',
      fieldName: 'initPrice',
      label: tooltipLabel('初始价格', '分'),
      componentProps: { placeholder: '请输入初始价格（分）', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'currentPrice',
      label: tooltipLabel('当前价格', '分'),
      componentProps: { placeholder: '请输入当前价格（分）', class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'priceRiseRate',
      label: tooltipLabel('每次加价比例', '百分比，如10=10%'),
      componentProps: { placeholder: '请输入每次加价比例（百分比，如10=10%）' },
    },
    {
      component: 'Input',
      fieldName: 'platformFeeRate',
      label: tooltipLabel('平台扣除比例', '百分比，如5=5%'),
      componentProps: { placeholder: '请输入平台扣除比例（百分比，如5=5%）' },
    },
    {
      component: 'TreeSelect',
      fieldName: 'ownerID',
      label: '当前持有人',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'nickname', value: 'id', children: 'children' },
        placeholder: '请选择当前持有人',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'Input',
      fieldName: 'tradeCount',
      label: '流转次数',
      componentProps: { placeholder: '请输入流转次数' },
    },
    {
      component: 'Select',
      fieldName: 'goodsStatus',
      label: '商品状态',
      componentProps: { options: goodsStatusOptions, placeholder: '请选择商品状态', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注', maxlength: 500 },
    },
    {
      component: 'InputNumber',
      fieldName: 'sort',
      label: tooltipLabel('排序', '升序'),
      componentProps: { placeholder: '请输入排序（升序）', class: 'w-full' },
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
      | WarehouseGoodsCreateParams
      | undefined;
    if (!values) return;
    if (values.initPrice != null) {
      (values as any).initPrice = Math.round(Number(values.initPrice) * 100);
    }
    if (values.currentPrice != null) {
      (values as any).currentPrice = Math.round(Number(values.currentPrice) * 100);
    }
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateWarehouseGoods({ id: editId.value, ...values } as WarehouseGoodsUpdateParams);
        message.success('更新成功');
      } else {
        await createWarehouseGoods(values);
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
    // 加载当前持有人树形数据
    try {
      const userRes = await getUserTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      ownerIDOptions.value = userRes ?? [];
      formApi.updateSchema([
        {
          fieldName: 'ownerID',
          componentProps: { treeData: ownerIDOptions.value },
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
      modalApi.setState({ title: '编辑仓库商品' });
      try {
        const detail = await getWarehouseGoodsDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          const formData = { ...detail };
          if (formData.initPrice != null) {
            formData.initPrice = formData.initPrice / 100;
          }
          if (formData.currentPrice != null) {
            formData.currentPrice = formData.currentPrice / 100;
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
      modalApi.setState({ title: '新建仓库商品' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
