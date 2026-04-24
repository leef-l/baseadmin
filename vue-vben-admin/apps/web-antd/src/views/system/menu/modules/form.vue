<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { message } from 'ant-design-vue';
import {
  getMenuDetail,
  createMenu,
  updateMenu,
  getMenuTree,
} from '#/api/system/menu';
import type {
  MenuCreateParams,
  MenuItem,
  MenuUpdateParams,
} from '#/api/system/menu/types';

const treeData = ref<MenuItem[]>([]);

/** 类型选项 */
const typeOptions = [
  { label: '目录', value: 1 },
  { label: '菜单', value: 2 },
  { label: '按钮', value: 3 },
  { label: '外链', value: 4 },
  { label: '内链', value: 5 },
];

const cacheDeps = {
  triggerFields: ['type'],
  show(values: Record<string, any>) {
    return Number(values.type ?? 0) === 2;
  },
};

const emit = defineEmits<{ success: [] }>();
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);
const routeMenuTypes = new Set([1, 2, 4, 5]);

function normalizeText(value?: string) {
  return value?.trim() ?? '';
}

function normalizeMenuValues(values: MenuCreateParams): MenuCreateParams {
  const type = Number(values.type ?? 0);
  const normalized: MenuCreateParams = {
    ...values,
    title: normalizeText(values.title),
    path: normalizeText(values.path),
    component: normalizeText(values.component),
    permission: normalizeText(values.permission),
    icon: normalizeText(values.icon),
    linkURL: normalizeText(values.linkURL),
    isCache: type === 2 ? Number(values.isCache ?? 0) : 0,
  };

  switch (type) {
    case 1: {
      normalized.component = '';
      normalized.linkURL = '';
      break;
    }
    case 2: {
      normalized.linkURL = '';
      break;
    }
    case 3: {
      normalized.path = '';
      normalized.component = '';
      normalized.linkURL = '';
      break;
    }
  }

  return normalized;
}

function validateMenuValues(values: MenuCreateParams) {
  const type = Number(values.type ?? 0);
  const path = normalizeText(values.path);
  const component = normalizeText(values.component);
  const permission = normalizeText(values.permission);
  const linkURL = normalizeText(values.linkURL);

  if (routeMenuTypes.has(type) && !path) {
    return '当前菜单类型必须填写前端路由路径';
  }
  if (routeMenuTypes.has(type) && !path.startsWith('/')) {
    return '前端路由路径必须以 / 开头';
  }
  if (type === 2 && !component) {
    return '菜单类型必须填写前端组件路径';
  }
  if (type === 3 && !permission) {
    return '按钮类型必须填写权限标识';
  }
  if ((type === 4 || type === 5) && !linkURL) {
    return '外链/内链类型必须填写地址';
  }
  return '';
}

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'TreeSelect',
      fieldName: 'parentID',
      label: '上级菜单',
      componentProps: {
        treeData: treeData.value,
        fieldNames: { label: 'title', value: 'id', children: 'children' },
        placeholder: '请选择上级菜单',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'Input',
      fieldName: 'title',
      label: '菜单名称',
      rules: 'required',
      componentProps: { placeholder: '请输入菜单名称', maxlength: 50 },
    },
    {
      component: 'Select',
      fieldName: 'type',
      label: '类型',
      rules: 'selectRequired',
      componentProps: { options: typeOptions, placeholder: '请选择类型', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'path',
      label: '前端路由路径',
      componentProps: { placeholder: '请输入前端路由路径', maxlength: 200 },
    },
    {
      component: 'Input',
      fieldName: 'component',
      label: '前端组件路径',
      componentProps: { placeholder: '请输入前端组件路径', maxlength: 200 },
    },
    {
      component: 'Input',
      fieldName: 'permission',
      label: '权限标识',
      componentProps: { placeholder: '请输入权限标识，如 system:dept:list', maxlength: 100 },
    },
    {
      component: 'Input',
      fieldName: 'icon',
      label: '菜单图标（图标名称）',
      componentProps: { placeholder: '请输入菜单图标（图标名称）', maxlength: 100 },
    },
    {
      component: 'InputNumber',
      fieldName: 'sort',
      label: '排序（升序）',
      componentProps: { placeholder: '请输入排序（升序）', min: 0, class: 'w-full' },
    },
    {
      component: 'Switch',
      fieldName: 'isShow',
      label: '是否显示',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
    },
    {
      component: 'Switch',
      fieldName: 'isCache',
      label: '是否缓存',
      dependencies: cacheDeps,
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 0,
    },
    {
      component: 'Input',
      fieldName: 'linkURL',
      label: '外链/内链地址（type=4或5时有效）',
      componentProps: { placeholder: '请输入外链/内链地址（type=4或5时有效）', maxlength: 500 },
    },
    {
      component: 'Switch',
      fieldName: 'status',
      label: '状态',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
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
      | MenuCreateParams
      | undefined;
    if (!values) return;
    const normalizedValues = normalizeMenuValues(values);
    const validationError = validateMenuValues(normalizedValues);
    if (validationError) {
      message.error(validationError);
      return;
    }
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateMenu({ id: editId.value, ...normalizedValues } as MenuUpdateParams);
        message.success('更新成功');
      } else {
        await createMenu(normalizedValues);
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
    // 加载树形数据
    try {
      const res = await getMenuTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      treeData.value = [
        { id: '0', title: '顶级节点', children: res ?? [] } as MenuItem,
      ];
      formApi.updateSchema([
        {
          fieldName: 'parentID',
          componentProps: { treeData: treeData.value },
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
      modalApi.setState({ title: '编辑菜单表' });
      try {
        const detail = await getMenuDetail(data.id);
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
      modalApi.setState({ title: '新建菜单表' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
