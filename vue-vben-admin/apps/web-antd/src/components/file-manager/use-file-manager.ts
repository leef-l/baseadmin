import { useVbenModal } from '@vben/common-ui';

import FileManagerModal from './modal.vue';

export type { FileManagerItem } from './index.vue';

export function useFileManager() {
  const [ModalComp, modalApi] = useVbenModal({
    connectedComponent: FileManagerModal,
  });

  function open(options?: {
    mode?: 'all' | 'file' | 'image';
    multiple?: boolean;
    maxCount?: number;
    accept?: string;
    maxSize?: number;
  }) {
    modalApi.setData(options || {});
    modalApi.open();
  }

  return { ModalComp, modalApi, open };
}
