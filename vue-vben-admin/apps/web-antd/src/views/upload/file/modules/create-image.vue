<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { message } from 'ant-design-vue';

import ImageUpload from '#/components/upload/image-upload.vue';

const emit = defineEmits<{ success: [] }>();
const imageValue = ref('');

const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  onCancel() {
    modalApi.close();
  },
  onConfirm() {
    if (!imageValue.value) {
      message.warning('请先选择图片');
      return;
    }
    emit('success');
    message.success('图片已创建');
    modalApi.close();
  },
  onOpenChange(isOpen: boolean) {
    if (isOpen) {
      imageValue.value = '';
      modalApi.setState({ title: '新建图片' });
    }
  },
});
</script>

<template>
  <Modal class="w-[640px]">
    <ImageUpload v-model:value="imageValue" :max-count="10" />
  </Modal>
</template>
