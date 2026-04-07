<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { message } from 'ant-design-vue';
import {
  getFileDetail,
  createFile,
  updateFile,
} from '#/api/upload/file';
import { getDirTree } from '#/api/upload/dir';
import type { DirItem } from '#/api/upload/dir/types';

/** 存储类型选项 */
const storageOptions = [
  { label: '本地', value: 1 },
  { label: '阿里云OSS', value: 2 },
  { label: '腾讯云COS', value: 3 },
];

/** 是否图片选项 */
const isImageOptions = [
  { label: '否', value: 0 },
  { label: '是', value: 1 },
];

const emit = defineEmits<{ success: [] }>();
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);

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
  let options: { label: string; value: string }[] = [];
  try {
    const list = await getDirTree();
    options = flattenDirTree(list);
  } catch {
    options = [];
  }
  return options;
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

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'Select',
      fieldName: 'dirID',
      label: '所属目录',
      componentProps: { options: dirIDOptions, placeholder: '请选择所属目录', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'name',
      label: '文件名称',
      rules: 'required',
      componentProps: { placeholder: '请输入文件名称', maxlength: 255 },
    },
    {
      component: 'Input',
      fieldName: 'url',
      label: '文件地址',
      rules: 'required',
      componentProps: { placeholder: '请输入文件地址', maxlength: 500 },
    },
    {
      component: 'Input',
      fieldName: 'ext',
      label: '文件扩展名',
      componentProps: { placeholder: '请输入文件扩展名', maxlength: 20 },
    },
    {
      component: 'Input',
      fieldName: 'size',
      label: '文件大小',
      componentProps: { placeholder: '请输入文件大小' },
    },
    {
      component: 'Input',
      fieldName: 'mime',
      label: 'MIME类型',
      componentProps: { placeholder: '请输入MIME类型', maxlength: 100 },
    },
    {
      component: 'Select',
      fieldName: 'storage',
      label: '存储类型',
      componentProps: { options: storageOptions, placeholder: '请选择存储类型', allowClear: true, class: 'w-full' },
    },
    {
      component: 'RadioGroup',
      fieldName: 'isImage',
      label: '是否图片',
      componentProps: { options: isImageOptions },
      defaultValue: 0,
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
    const values = await formApi.validateAndSubmitForm();
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateFile({ id: editId.value, ...values });
        message.success('更新成功');
      } else {
        await createFile(values);
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
    const data = modalApi.getData<{ id?: string } | null>();
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑文件记录' });
      try {
        const detail = await getFileDetail(data.id);
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
      modalApi.setState({ title: '新建文件记录' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
