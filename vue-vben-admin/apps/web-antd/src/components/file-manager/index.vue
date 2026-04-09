<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { PlusOutlined } from '@ant-design/icons-vue';
import {
  Button,
  Card,
  Checkbox,
  Col,
  Empty,
  Input,
  message,
  Pagination,
  Row,
  Spin,
  Tree,
  Upload,
} from 'ant-design-vue';
import type { UploadProps } from 'ant-design-vue';

import { getDirTree } from '#/api/upload/dir';
import type { DirItem } from '#/api/upload/dir/types';
import { getFileList } from '#/api/upload/file';
import type { FileItem } from '#/api/upload/file/types';
import { uploadFile } from '#/api/upload/uploader';

export interface FileManagerItem {
  id: string;
  url: string;
  name: string;
  size: number;
  ext: string;
  mime: string;
  isImage: number;
}

interface Props {
  /** 模式: image=只看图片, file=只看非图片, all=全部 */
  mode?: 'all' | 'file' | 'image';
  /** 是否多选 */
  multiple?: boolean;
  /** 最多选几个 */
  maxCount?: number;
  /** 上传文件类型限制 (如 image/*,.pdf) */
  accept?: string;
  /** 最大文件大小 MB */
  maxSize?: number;
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'all',
  multiple: false,
  maxCount: 1,
  accept: '',
  maxSize: 10,
});

const emit = defineEmits<{
  confirm: [files: FileManagerItem[]];
}>();

const effectiveMax = computed(() => (props.multiple ? props.maxCount : 1));
const uploadMultiple = computed(() => props.multiple && effectiveMax.value > 1);

/** State */
const loading = ref(false);
const treeLoading = ref(false);
const uploadingCount = ref(0);
const uploading = computed(() => uploadingCount.value > 0);
const keyword = ref('');
const selectedDirId = ref<string | undefined>(undefined);
const fileListData = ref<FileItem[]>([]);
const selectedIds = ref<Set<string>>(new Set());
const dirTreeData = ref<DirItem[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0 });

/** 图片扩展名集合 */
const IMAGE_EXTS = new Set(['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp', 'ico']);

function isImageFile(file: FileItem): boolean {
  if (file.isImage === 1) return true;
  return IMAGE_EXTS.has((file.ext || '').toLowerCase());
}

/** 文件图标 */
function getFileIcon(ext: string): string {
  const e = (ext || '').toLowerCase();
  if (['pdf'].includes(e)) return '📄';
  if (['doc', 'docx'].includes(e)) return '📝';
  if (['xls', 'xlsx'].includes(e)) return '📊';
  if (['ppt', 'pptx'].includes(e)) return '📎';
  if (['zip', 'rar', '7z', 'gz', 'tar'].includes(e)) return '📦';
  if (['mp4', 'avi', 'mov', 'mkv'].includes(e)) return '🎬';
  if (['mp3', 'wav', 'flac', 'aac'].includes(e)) return '🎵';
  return '📁';
}

/** 格式化文件大小 */
function formatSize(bytes: number | string | undefined): string {
  const b = Number(bytes) || 0;
  if (b < 1024) return `${b}B`;
  if (b < 1024 * 1024) return `${(b / 1024).toFixed(1)}KB`;
  return `${(b / (1024 * 1024)).toFixed(1)}MB`;
}

/** 目录树节点 */
const treeNodes = computed(() => {
  return [{ id: undefined, name: '全部文件', children: dirTreeData.value }] as any[];
});

/** 加载目录树 */
async function loadDirTree() {
  treeLoading.value = true;
  try {
    dirTreeData.value = await getDirTree();
  } finally {
    treeLoading.value = false;
  }
}

/** 加载文件列表 */
async function loadFileList() {
  loading.value = true;
  try {
    const isImageParam = props.mode === 'image' ? 1 : props.mode === 'file' ? 0 : undefined;
    const res = await getFileList({
      pageNum: pagination.current,
      pageSize: pagination.pageSize,
      isImage: isImageParam,
      dirID: selectedDirId.value,
      name: keyword.value || undefined,
    });
    fileListData.value = res?.list ?? [];
    pagination.total = res?.total ?? 0;
    syncSelectedIdsWithVisibleFiles();
  } finally {
    loading.value = false;
  }
}

function syncSelectedIdsWithVisibleFiles() {
  if (selectedIds.value.size === 0) {
    return;
  }
  const visible = new Set(fileListData.value.map((item) => item.id));
  selectedIds.value = new Set(
    [...selectedIds.value].filter((id) => visible.has(id)),
  );
}

/** 切换选中 */
function toggleSelect(file: FileItem) {
  const nextIds = new Set(selectedIds.value);
  if (nextIds.has(file.id)) {
    nextIds.delete(file.id);
  } else {
    if (!props.multiple) {
      nextIds.clear();
    }
    if (nextIds.size < effectiveMax.value) {
      nextIds.add(file.id);
    }
  }
  selectedIds.value = nextIds;
}

/** 目录选择 */
function onDirSelect(keys: (string | undefined)[]) {
  selectedDirId.value = keys[0] ?? undefined;
  pagination.current = 1;
  loadFileList();
}

/** 搜索 */
function handleSearch() {
  pagination.current = 1;
  loadFileList();
}

/** 分页 */
function onPageChange(page: number) {
  pagination.current = page;
  loadFileList();
}

function onPageSizeChange(_current: number, size: number) {
  pagination.pageSize = size;
  pagination.current = 1;
  loadFileList();
}

/** 上传 */
const customUpload: UploadProps['customRequest'] = async (options) => {
  const { file, onSuccess, onError } = options;
  const f = file as File;
  if (props.maxSize && f.size > props.maxSize * 1024 * 1024) {
    message.error(`文件大小不能超过 ${props.maxSize}MB`);
    onError?.(new Error('文件过大') as any);
    return;
  }
  uploadingCount.value += 1;
  try {
    await uploadFile(f, selectedDirId.value);
    message.success('上传成功');
    onSuccess?.({});
    loadFileList();
  } catch (err: any) {
    message.error('上传失败');
    onError?.(err);
  } finally {
    uploadingCount.value = Math.max(0, uploadingCount.value - 1);
  }
};

/** 确认选择 */
function handleConfirm() {
  const selected = fileListData.value
    .filter((f) => selectedIds.value.has(f.id))
    .map((f) => ({
      id: f.id,
      url: f.url,
      name: f.name,
      size: Number(f.size) || 0,
      ext: f.ext ?? '',
      mime: f.mime ?? '',
      isImage: f.isImage ?? 0,
    }));
  emit('confirm', selected);
}

/** 重置 */
function reset() {
  selectedIds.value = new Set();
  keyword.value = '';
  selectedDirId.value = undefined;
  fileListData.value = [];
  pagination.current = 1;
  pagination.total = 0;
}

defineExpose({ reset, loadDirTree, loadFileList });

onMounted(() => {
  loadDirTree();
  loadFileList();
});
// PLACEHOLDER_TEMPLATE_MARKER
</script>

<template>
  <div class="file-manager">
    <!-- 顶部工具栏 -->
    <div class="fm-header">
      <Upload
        :accept="accept"
        :custom-request="customUpload"
        :show-upload-list="false"
        :multiple="uploadMultiple"
        :disabled="uploading"
      >
        <Button type="primary" :loading="uploading">
          <PlusOutlined /> 上传文件
        </Button>
      </Upload>
      <Input.Search
        v-model:value="keyword"
        placeholder="搜索文件名"
        style="width: 220px; margin-left: 12px"
        allow-clear
        @search="handleSearch"
      />
      <span class="fm-count">
        已选 {{ selectedIds.size }} / {{ effectiveMax }}
      </span>
      <Button type="primary" :disabled="selectedIds.size === 0" @click="handleConfirm">
        确认选择
      </Button>
    </div>

    <div class="fm-body">
      <!-- 左侧目录树 -->
      <div class="fm-sidebar">
        <Spin :spinning="treeLoading">
          <Tree
            :tree-data="treeNodes"
            :field-names="{ title: 'name', key: 'id', children: 'children' }"
            default-expand-all
            block-node
            @select="onDirSelect"
          />
        </Spin>
      </div>

      <!-- 右侧文件网格 -->
      <div class="fm-content">
        <Spin :spinning="loading">
          <template v-if="fileListData.length > 0">
            <Row :gutter="[12, 12]">
              <Col v-for="file in fileListData" :key="file.id" :span="4">
                <Card
                  size="small"
                  hoverable
                  :class="['fm-card', { 'fm-card--selected': selectedIds.has(file.id) }]"
                  @click="toggleSelect(file)"
                >
                  <template #cover>
                    <div class="fm-thumb">
                      <Checkbox
                        class="fm-checkbox"
                        :checked="selectedIds.has(file.id)"
                        @click.stop
                        @change="toggleSelect(file)"
                      />
                      <img v-if="isImageFile(file)" :src="file.url" :alt="file.name" />
                      <span v-else class="fm-icon">{{ getFileIcon(file.ext || '') }}</span>
                    </div>
                  </template>
                  <Card.Meta>
                    <template #title>
                      <span class="fm-name" :title="file.name">{{ file.name }}</span>
                    </template>
                    <template #description>
                      <span class="fm-size">{{ formatSize(file.size) }}</span>
                    </template>
                  </Card.Meta>
                </Card>
              </Col>
            </Row>
          </template>
          <Empty v-else description="暂无文件" />
        </Spin>

        <div v-if="pagination.total > 0" class="fm-pagination">
          <Pagination
            v-model:current="pagination.current"
            :total="pagination.total"
            :page-size="pagination.pageSize"
            size="small"
            show-size-changer
            :page-size-options="['20', '40', '60']"
            @change="onPageChange"
            @show-size-change="onPageSizeChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.file-manager {
  display: flex;
  flex-direction: column;
  height: 520px;
}

.fm-header {
  display: flex;
  align-items: center;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.fm-count {
  flex: 1;
  text-align: right;
  color: #666;
  font-size: 13px;
  margin-right: 12px;
}

.fm-body {
  display: flex;
  flex: 1;
  gap: 12px;
  padding-top: 12px;
  overflow: hidden;
}

.fm-sidebar {
  width: 200px;
  min-width: 200px;
  overflow-y: auto;
  border-right: 1px solid #f0f0f0;
  padding-right: 12px;
}

.fm-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.fm-card {
  cursor: pointer;
}

.fm-card--selected {
  border-color: #1677ff;
  box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.2);
}

.fm-thumb {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
}

.fm-thumb img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.fm-icon {
  font-size: 36px;
}

.fm-checkbox {
  position: absolute;
  top: 4px;
  left: 4px;
  z-index: 1;
}

.fm-name {
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
}

.fm-size {
  font-size: 11px;
  color: #999;
}

.fm-pagination {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
  margin-top: auto;
}
</style>
