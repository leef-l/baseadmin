<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { message } from 'ant-design-vue';

import FileUpload from '#/components/upload/file-upload.vue';

const emit = defineEmits<{ success: [] }>();
const fileValue = ref('');

const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  onCancel() {
    modalApi.close();
  },
  onConfirm() {
    if (!fileValue.value) {
      message.warning('请先选择文件');
      return;
    }
    emit('success');
    message.success('文件已创建');
    modalApi.close();
  },
  onOpenChange(isOpen: boolean) {
    if (isOpen) {
      fileValue.value = '';
      modalApi.setState({ title: '新建文件' });
    }
  },
});
</script>

<template>
  <Modal class="w-[640px]">
    <FileUpload v-model:value="fileValue" :max-count="10" />
  </Modal>
</template>
