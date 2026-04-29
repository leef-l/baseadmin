<script setup lang="ts">
import type { UploadFile, UploadProps } from 'ant-design-vue';

import { computed, ref } from 'vue';

import { UploadOutlined } from '@ant-design/icons-vue';
import { Button, message, Upload } from 'ant-design-vue';

import { uploadFile } from '#/api/upload/file';
import { parseUploadValue } from './shared';

interface Props {
  value?: string;
  maxCount?: number;
  accept?: string;
  maxSize?: number;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  value: '',
  maxCount: 1,
  accept: '',
  maxSize: 10,
  disabled: false,
});

const emit = defineEmits<{
  'update:value': [val: string];
}>();

type ManagedUploadFile = UploadFile & {
  savedUrl?: string;
};

const effectiveMaxCount = computed(() => Math.max(1, Number(props.maxCount) || 1));
const uploadingFiles = ref<ManagedUploadFile[]>([]);
const uploading = computed(() => uploadingFiles.value.length > 0);

const selectedFiles = computed<ManagedUploadFile[]>(() =>
  parseUploadValue(props.value).map((url, i) => ({
    uid: `saved-${i}-${url}`,
    name: url.split('/').pop() || url,
    status: 'done',
    url,
    savedUrl: url,
  })),
);

const displayFileList = computed<UploadProps['fileList']>(() => [
  ...selectedFiles.value,
  ...uploadingFiles.value,
]);

const remainingCount = computed(() =>
  Math.max(0, effectiveMaxCount.value - (displayFileList.value?.length ?? 0)),
);

function emitUrls(urls: string[]) {
  emit('update:value', parseUploadValue(urls.join(',')).join(','));
}

const beforeUpload: UploadProps['beforeUpload'] = (file, fileList) => {
  const batchIndex = fileList.findIndex((item) => item.uid === file.uid);
  if (remainingCount.value <= 0 || batchIndex >= remainingCount.value) {
    if (batchIndex <= 0) {
      message.warning(`最多只能选择 ${effectiveMaxCount.value} 个文件`);
    }
    return false;
  }
  if (props.maxSize && file.size > props.maxSize * 1024 * 1024) {
    message.error(`文件大小不能超过 ${props.maxSize}MB`);
    return false;
  }
  return true;
};

const customUpload: UploadProps['customRequest'] = async (options) => {
  const { file, onError, onSuccess } = options;
  const f = file as File & { uid?: string };
  const uid = f.uid || `${Date.now()}-${f.name}`;
  const uploadingFile: ManagedUploadFile = {
    uid,
    name: f.name,
    status: 'uploading',
    size: f.size,
    type: f.type,
  };
  uploadingFiles.value = [...uploadingFiles.value, uploadingFile];

  try {
    const result = await uploadFile(f);
    emitUrls([...parseUploadValue(props.value), result.url]);
    message.success('上传成功');
    onSuccess?.(result as any);
  } catch (error: any) {
    message.error('上传失败');
    onError?.(error);
  } finally {
    uploadingFiles.value = uploadingFiles.value.filter((item) => item.uid !== uid);
  }
};

const handleRemove: UploadProps['onRemove'] = (file) => {
  if (props.disabled) {
    return false;
  }
  const managedFile = file as ManagedUploadFile;
  const savedUrl = managedFile.savedUrl || managedFile.url;
  if (savedUrl) {
    emitUrls(parseUploadValue(props.value).filter((url) => url !== savedUrl));
    return true;
  }
  uploadingFiles.value = uploadingFiles.value.filter((item) => item.uid !== managedFile.uid);
  return true;
};

const handlePreview: UploadProps['onPreview'] = (file) => {
  const url = (file as ManagedUploadFile).savedUrl || file.url;
  if (url) {
    window.open(url, '_blank', 'noopener,noreferrer');
  }
};
</script>

<template>
  <Upload
    :accept="accept || undefined"
    :before-upload="beforeUpload"
    :custom-request="customUpload"
    :disabled="disabled"
    :file-list="displayFileList"
    :multiple="effectiveMaxCount > 1"
    :show-upload-list="{ showRemoveIcon: !disabled }"
    @preview="handlePreview"
    @remove="handleRemove"
  >
    <Button
      v-if="(displayFileList?.length ?? 0) < effectiveMaxCount && !disabled"
      :loading="uploading"
    >
      <UploadOutlined />
      选择文件
    </Button>
  </Upload>
</template>
