<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import FileManager from './index.vue';
import type { FileManagerItem } from './index.vue';

const pickerRef = ref<InstanceType<typeof FileManager>>();
const mode = ref<'all' | 'file' | 'image'>('all');
const multiple = ref(false);
const maxCount = ref(1);
const accept = ref('');
const maxSize = ref(10);

const emit = defineEmits<{
  confirm: [files: FileManagerItem[]];
}>();

const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  onCancel() {
    modalApi.close();
  },
  onOpenChange(isOpen: boolean) {
    if (isOpen) {
      const data = modalApi.getData<{
        mode?: 'all' | 'file' | 'image';
        multiple?: boolean;
        maxCount?: number;
        accept?: string;
        maxSize?: number;
      } | null>();
      mode.value = data?.mode ?? 'all';
      multiple.value = data?.multiple ?? false;
      maxCount.value = data?.maxCount ?? 1;
      accept.value = data?.accept ?? '';
      maxSize.value = data?.maxSize ?? 10;
      pickerRef.value?.reset();
      pickerRef.value?.loadDirTree?.();
      pickerRef.value?.loadFileList();
    }
  },
});

function handleConfirm(files: FileManagerItem[]) {
  emit('confirm', files);
  modalApi.close();
}
</script>

<template>
  <Modal title="文件管理器" class="w-[960px]">
    <FileManager
      ref="pickerRef"
      :mode="mode"
      :multiple="multiple"
      :max-count="maxCount"
      :accept="accept"
      :max-size="maxSize"
      @confirm="handleConfirm"
    />
  </Modal>
</template>
