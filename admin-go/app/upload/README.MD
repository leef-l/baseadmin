# GBaseAdmin - Upload 应用

文件管理应用，提供文件上传、目录管理、上传配置等功能。

## 数据库表

| 表名 | 说明 |
|------|------|
| `upload_dir` | 文件目录（树形结构） |
| `upload_file` | 文件记录 |
| `upload_config` | 上传配置（本地/OSS/COS） |
| `upload_dir_rule` | 文件目录规则 |

## API 路由

所有接口前缀：`/api/upload/`，需要 Auth 认证。

| 模块 | 路径 | 说明 |
|------|------|------|
| 文件上传 | `POST /api/upload/uploader/upload` | 上传文件（multipart/form-data） |
| 文件目录 | `/api/upload/dir/list,create,update,delete,detail,tree` | 目录管理（树形） |
| 文件记录 | `/api/upload/file/list,create,update,delete,detail` | 文件管理 |
| 上传配置 | `/api/upload/config/list,create,update,delete,detail` | 配置管理 |
| 目录规则 | `/api/upload/dir_rule/list,create,update,delete,detail` | 目录规则管理 |

静态文件访问：`/upload/` 映射到 `resource/upload/` 目录。

默认初始化数据中的本地存储目录也是 `resource/upload`。

## 上传流程

1. 前端调用 `POST /api/upload/uploader/upload`，发送 `file`（文件）+ `dirId`（可选目录ID）
2. 后端读取 `upload_config` 中 `is_default=1` 的默认配置，获取存储类型和大小限制
3. 验证文件大小，解析扩展名和 MIME 类型，判断是否为图片
4. 本地存储：保存到 `resource/upload/{日期}/{唯一文件名}`
5. 写入 `upload_file` 记录，返回文件信息（id, url, name, size, ext, mime, isImage）
6. 如果数据库写入失败，已上传成功的 OSS/COS 对象会执行回滚删除，避免残留孤儿文件

## 存储类型

| 值 | 类型 | 配置字段 |
|----|------|---------|
| 1 | 本地存储 | `local_path` |
| 2 | 阿里云 OSS | `oss_endpoint`, `oss_bucket`, `oss_access_key`, `oss_secret_key` |
| 3 | 腾讯云 COS | `cos_region`, `cos_bucket`, `cos_secret_id`, `cos_secret_key` |

## 文件删除

删除文件时同时执行物理删除，根据 `storage` 字段判断存储类型：
- 本地：直接删除磁盘文件
- OSS/COS：按文件 URL 匹配真实配置后删除对应云端对象

## 前端组件

本应用提供以下可复用组件：

| 组件 | 路径 | 说明 |
|------|------|------|
| FileManager | `components/file-manager/` | 通用文件管理器（目录树 + 文件网格 + 上传 + 勾选） |
| ImageUpload | `components/upload/image-upload.vue` | 图片上传（调用 FileManager，mode=image） |
| FileUpload | `components/upload/file-upload.vue` | 文件上传（调用 FileManager，mode=all） |
| RichText | `components/tinymce/index.vue` | TinyMCE v8 富文本编辑器（图片上传走 upload API） |
| JsonEditor | `components/json-editor/index.vue` | JSON 编辑器（tree + code 双模式） |

### FileManager 文件管理器

核心组件，支持三种模式：

| mode | 说明 | 使用场景 |
|------|------|---------|
| `image` | 只显示图片文件 | ImageUpload 选择图片 |
| `file` | 只显示非图片文件 | 选择文档/附件 |
| `all` | 显示所有文件 | FileUpload 选择任意文件 |

Props: `mode`, `multiple`, `maxCount`, `accept`, `maxSize`

功能：左侧目录树筛选、右侧文件网格、勾选多选、搜索、分页、上传文件

调用方式：

```typescript
// 方式1: 通过 composable
import { useFileManager } from '#/components/file-manager/use-file-manager';
const { ModalComp, open } = useFileManager();
open({ mode: 'image', multiple: true, maxCount: 5 });

// 方式2: 通过 useVbenModal
import FileManagerModal from '#/components/file-manager/modal.vue';
```

### 表单组件注册

这些组件已注册到 `adapter/component/index.ts`，可在 vben-form 中直接使用：

```typescript
{ component: 'ImageUpload', fieldName: 'avatar', label: '头像', componentProps: { maxCount: 1 } }
{ component: 'FileUpload', fieldName: 'attachment', label: '附件', componentProps: { maxCount: 3 } }
{ component: 'RichText', fieldName: 'content', label: '内容', formItemClass: 'col-span-full' }
{ component: 'JsonEditor', fieldName: 'config', label: '配置', formItemClass: 'col-span-full' }
```

### 前端 API

所有 API 通过 `requestClient`（baseURL = `/api`）调用：

| API 文件 | 前缀 | 说明 |
|---------|------|------|
| `api/upload/uploader/index.ts` | `/upload/uploader/upload` | 文件上传 |
| `api/upload/file/index.ts` | `/upload/file/` | 文件记录 CRUD |
| `api/upload/dir/index.ts` | `/upload/dir/` | 目录 CRUD + tree |
| `api/upload/config/index.ts` | `/upload/config/` | 配置 CRUD |
| `api/upload/dir_rule/index.ts` | `/upload/dir_rule/` | 目录规则 CRUD |

## 前端页面

| 页面 | 说明 |
|------|------|
| 文件列表 | 左侧目录树 + 列表/区块双视图切换 + 图片预览 + 上传 |
| 上传配置 | storage RadioGroup 联动显示本地/OSS/COS 配置字段 |
| 文件目录 | 树形 CRUD |
| 目录规则 | 标准 CRUD |
