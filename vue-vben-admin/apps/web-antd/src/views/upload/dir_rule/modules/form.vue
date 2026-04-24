<script setup lang="ts">
import type { DirItem } from '#/api/upload/dir/types';
import type {
  DirRuleCreateParams,
  DirRuleItem,
  DirRuleStorageType,
  DirRuleUpdateParams,
} from '#/api/upload/dir_rule/types';

import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { getDirTree } from '#/api/upload/dir';
import {
  createDirRule,
  getDirRuleDetail,
  updateDirRule,
} from '#/api/upload/dir_rule';

const emit = defineEmits<{ success: [] }>();

/** 类别选项 */
const categoryOptions = [
  { label: '默认', value: 1 },
  { label: '类型', value: 2 },
  { label: '来源', value: 3 },
];

const storageTypeOptions = [
  { label: '本地', value: 1 },
  { label: 'OSS', value: 2 },
  { label: 'COS', value: 3 },
];

const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);

type DirRuleFormValues = DirRuleCreateParams & {
  storageTypes?: DirRuleStorageType[] | string;
};

const fileTypeDeps = {
  triggerFields: ['category'],
  show(values: Record<string, any>) {
    return values.category === 2 || values.category === 3;
  },
};

/** 目录下拉选项 */
const dirIDOptions = ref<{ label: string; value: string }[]>([]);

/** 将树形目录打平为选项 */
function flattenDirTree(
  items: DirItem[],
  prefix = '',
): { label: string; value: string }[] {
  const result: { label: string; value: string }[] = [];
  for (const item of items) {
    const label = prefix ? `${prefix} / ${item.name}` : item.name;
    result.push({ label, value: item.id });
    if (item.children?.length) {
      result.push(...flattenDirTree(item.children, label));
    }
  }
  return result;
}

/** 加载目录选项 */
async function loadDirOptions() {
  try {
    const list = await getDirTree();
    return flattenDirTree(list);
  } catch {
    return [];
  }
}

function applyDirOptions(options: { label: string; value: string }[]) {
  dirIDOptions.value = options;
  formApi.updateSchema([
    {
      fieldName: 'dirID',
      componentProps: { options: dirIDOptions.value },
    },
  ]);
}

function normalizeStorageTypes(value: unknown): DirRuleStorageType[] {
  const rawValues = Array.isArray(value)
    ? value
    : (typeof value === 'string'
      ? value.split(/[,\s，；;]+/g)
      : []);
  const normalized: DirRuleStorageType[] = [];
  const seen = new Set<DirRuleStorageType>();
  for (const item of rawValues) {
    const parsed = Number(String(item).trim());
    if (parsed !== 1 && parsed !== 2 && parsed !== 3) {
      continue;
    }
    const storageType = parsed as DirRuleStorageType;
    if (seen.has(storageType)) {
      continue;
    }
    seen.add(storageType);
    normalized.push(storageType);
  }
  return normalized.length > 0 ? normalized : [1, 2, 3];
}

function normalizeDetailValues(detail: DirRuleItem) {
  return {
    ...detail,
    storageTypes: normalizeStorageTypes(detail.storageTypes),
  };
}

function stringifyStorageTypes(value: unknown): string {
  return normalizeStorageTypes(value).join(',');
}

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'Select',
      fieldName: 'dirID',
      label: '所属目录',
      rules: 'selectRequired',
      componentProps: { options: dirIDOptions, placeholder: '请选择所属目录', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'category',
      label: '类别',
      rules: 'selectRequired',
      componentProps: { options: categoryOptions, placeholder: '请选择类别', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Textarea',
      fieldName: 'fileType',
      label: '匹配条件',
      dependencies: fileTypeDeps,
      componentProps: {
        placeholder: '类型规则可填扩展名、image、image/*、application/pdf；来源规则可填页面路由或上传场景，每行一个',
        autoSize: { minRows: 4, maxRows: 8 },
        maxlength: 1000,
      },
    },
    {
      component: 'Select',
      fieldName: 'storageTypes',
      label: '适用存储',
      rules: 'selectRequired',
      defaultValue: [1, 2, 3],
      componentProps: {
        options: storageTypeOptions,
        mode: 'multiple',
        placeholder: '请选择适用存储，可多选',
        allowClear: true,
        maxTagCount: 'responsive',
        class: 'w-full',
      },
    },
    {
      component: 'Input',
      fieldName: 'savePath',
      label: '保存目录',
      componentProps: { placeholder: '如 uploads/{systemUserId}/{Y-m-d}；本地存储可用 @up/cert/{ext}', maxlength: 500 },
    },
    {
      component: 'Switch',
      fieldName: 'keepName',
      label: '保留原文件名',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 0,
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
      | DirRuleFormValues
      | undefined;
    if (!values) return;
    const submitValues = {
      ...values,
      fileType: values.category === 2 || values.category === 3 ? values.fileType?.trim() : '',
      storageTypes: stringifyStorageTypes(values.storageTypes),
      keepName: values.keepName ?? 0,
    };
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateDirRule({ id: editId.value, ...submitValues } as DirRuleUpdateParams);
        message.success('更新成功');
      } else {
        await createDirRule(submitValues);
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
    const options = await loadDirOptions();
    if (currentOpenToken !== openToken.value) {
      return;
    }
    applyDirOptions(options);
    const data = modalApi.getData<{ id?: string }>();
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑文件目录规则' });
      try {
        const detail = await getDirRuleDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          formApi.setValues(normalizeDetailValues(detail));
        }
      } catch {
        if (currentOpenToken === openToken.value) {
          message.error('获取详情失败');
        }
      }
    } else {
      isEdit.value = false;
      editId.value = '';
      modalApi.setState({ title: '新建文件目录规则' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
