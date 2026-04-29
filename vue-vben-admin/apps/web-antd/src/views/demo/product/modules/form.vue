<script setup lang="ts">
import { h, ref } from 'vue';
import { isPlatformSuperAdminUser } from '@/utils/auth-scope';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getProductDetail,
  createProduct,
  updateProduct,
} from '#/api/demo/product';
import type {
  ProductCreateParams,
  ProductUpdateParams
} from '#/api/demo/product/types';
import { getCategoryTree } from '#/api/demo/category';
import type { CategoryItem } from '#/api/demo/category/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 类型选项 */
const typeOptions = [
  { label: '普通', value: 1 },
  { label: '置顶', value: 2 },
  { label: '推荐', value: 3 },
  { label: '热门', value: 4 },
];

/** 状态选项 */
const statusOptions = [
  { label: '草稿', value: 0 },
  { label: '上架', value: 1 },
  { label: '下架', value: 2 },
];
const categoryIDOptions = ref<CategoryItem[]>([]);
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
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'TreeSelect',
      fieldName: 'categoryID',
      label: '商品分类',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'name', value: 'id', children: 'children' },
        placeholder: '请选择商品分类',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'Input',
      fieldName: 'skuNo',
      label: 'SKU编号',
      rules: 'required',
      componentProps: { placeholder: '请输入SKU编号', maxlength: 50 },
    },
    {
      component: 'Input',
      fieldName: 'name',
      label: '商品名称',
      rules: 'required',
      componentProps: { placeholder: '请输入商品名称', maxlength: 120 },
    },
    {
      component: 'ImageUpload',
      fieldName: 'cover',
      label: '封面',
      componentProps: { maxCount: 1 },
    },
    {
      component: 'FileUpload',
      fieldName: 'manualFile',
      label: '说明书文件',
      componentProps: { maxCount: 1 },
    },
    {
      component: 'RichText',
      fieldName: 'detailContent',
      label: '详情内容',
      formItemClass: 'col-span-full',
    },
    {
      component: 'JsonEditor',
      fieldName: 'specJSON',
      label: '规格JSON',
      formItemClass: 'col-span-full',
    },
    {
      component: 'Input',
      fieldName: 'websiteURL',
      label: '官网URL',
      componentProps: { placeholder: '请输入URL地址', maxlength: 500, addonBefore: 'https://' },
    },
    {
      component: 'Select',
      fieldName: 'type',
      label: '类型',
      componentProps: { options: typeOptions, placeholder: '请选择类型', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Switch',
      fieldName: 'isRecommend',
      label: '是否推荐',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 0,
    },
    {
      component: 'InputNumber',
      fieldName: 'salePrice',
      label: tooltipLabel('销售价', '分'),
      componentProps: { placeholder: '请输入销售价（分）', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'stockNum',
      label: '库存数量',
      componentProps: { placeholder: '请输入库存数量', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'weightNum',
      label: tooltipLabel('重量', '克'),
      componentProps: { placeholder: '请输入重量（克）', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'sort',
      label: tooltipLabel('排序', '升序'),
      componentProps: { placeholder: '请输入排序（升序）', class: 'w-full' },
    },
    {
      component: 'IconPicker',
      fieldName: 'icon',
      label: '图标',
      componentProps: { placeholder: '请选择图标' },
    },
    {
      component: 'RadioGroup',
      fieldName: 'status',
      label: '状态',
      componentProps: { options: statusOptions },
    },
    {
      component: 'Select',
      ifShow: () => isPlatformSuperAdminUser(),
      fieldName: 'tenantID',
      label: '租户',
      componentProps: { options: [], placeholder: '请选择租户', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      ifShow: () => isPlatformSuperAdminUser(),
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
      | ProductCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateProduct({ id: editId.value, ...values } as ProductUpdateParams);
        message.success('更新成功');
      } else {
        await createProduct(values);
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
    // 加载商品分类树形数据
    try {
      const categoryRes = await getCategoryTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      categoryIDOptions.value = categoryRes ?? [];
      formApi.updateSchema([
        {
          fieldName: 'categoryID',
          componentProps: { treeData: categoryIDOptions.value },
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
      modalApi.setState({ title: '编辑体验商品' });
      try {
        const detail = await getProductDetail(data.id);
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
      modalApi.setState({ title: '新建体验商品' });
    }
  },
});
</script>

<template>
  <Modal class="w-[800px]">
    <Form />
  </Modal>
</template>
