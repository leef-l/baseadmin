# GBaseAdmin - Upload 应用

文件管理应用，提供文件上传、目录管理、上传配置等功能。

## 代码入口

- 启动入口：`main.go`
- 路由绑定：`internal/cmd/cmd.go`
- 当前挂载控制器：`config`、`dir`、`dir_rule`、`file`、`uploader`、`health`
- 默认开发端口：`10023`

## 数据库表

| 表名 | 说明 |
|------|------|
| `upload_dir` | 文件目录（树形结构） |
| `upload_file` | 文件记录 |
| `upload_config` | 上传配置（本地/OSS/COS） |
| `upload_dir_rule` | 文件目录规则 |

## API 路由

对外接口分两层：

- 公开路由：`GET /healthz`、`GET /readyz`
- 业务接口前缀：`/api/upload/`，全部需要 Auth 认证

| 模块 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| 健康检查 | `GET /healthz` | 否 | 返回 `{ status: "ok", service: "upload" }` |
| 就绪检查 | `GET /readyz` | 否 | 返回数据库与上传存储检查结果 |
| 文件上传 | `POST /api/upload/file/upload` | 是 | 上传文件（multipart/form-data） |
| 文件目录 | `/api/upload/dir/list,create,update,delete,batch-delete,detail,tree` | 是 | 目录管理（树形） |
| 文件记录 | `/api/upload/file/list,create,update,delete,batch-delete,detail` | 是 | 文件管理 |
| 上传配置 | `/api/upload/config/list,create,update,delete,batch-delete,detail` | 是 | 配置管理 |
| 目录规则 | `/api/upload/dir_rule/list,create,update,delete,batch-delete,detail` | 是 | 目录规则管理 |

静态文件访问：`/upload/` 映射到 `resource/upload/` 目录。

默认初始化数据中的本地存储目录也是 `resource/upload`。

## 后端模块

| 模块 | 目录 | 说明 |
|------|------|------|
| 健康检查 | `internal/controller/health` | 无鉴权探活接口 |
| 上传配置 | `internal/controller/config` | 本地 / OSS / COS 配置 CRUD |
| 文件目录 | `internal/controller/dir` | 树形目录 CRUD |
| 目录规则 | `internal/controller/dir_rule` | 目录命名规则 CRUD |
| 文件记录 | `internal/controller/file` | 文件记录查询与删除 |
| 文件上传 | `internal/controller/uploader` | 上传入口和落库 |

目录规则类别约定：
- `1=默认规则`
- `2=类型规则`
- `3=来源规则`，按上传来源标识匹配，通常来自当前页面路由

## 上传流程

1. 前端调用 `POST /api/upload/file/upload`，发送 `file`（文件）+ `dirId`（可选目录ID）
2. 后端读取 `upload_config` 中 `is_default=1` 的默认配置，获取存储类型和大小限制
3. 验证文件大小，解析扩展名和 MIME 类型，判断是否为图片
4. 本地存储：保存到 `resource/upload/{日期}/{唯一文件名}`
5. 写入 `upload_file` 记录，返回文件信息（id, url, name, size, ext, mime, isImage）
6. 如果数据库写入失败，已上传成功的 OSS/COS 对象会执行回滚删除，避免残留孤儿文件

目录规则命中优先级：`来源规则 > 类型规则 > 默认规则`。

线上 `Deploy GBaseAdmin` 会用专用 smoke 账号做一次真实链路冒烟：创建临时来源规则 -> 上传文件 -> 校验 `upload_file.dir_id` -> 校验公网 `/upload/...` 可访问。

## 存储类型

| 值 | 类型 | 配置字段 |
|----|------|---------|
| 1 | 本地存储 | `local_path` |
| 2 | 阿里云 OSS | `oss_endpoint`, `oss_bucket`, `oss_access_key`, `oss_secret_key` |
| 3 | 腾讯云 COS | `cos_region`, `cos_bucket`, `cos_secret_id`, `cos_secret_key` |

## 文件删除

删除文件时同时执行物理删除，根据 `storage` 字段判断存储类型：
- 本地：先把文件原子移动到 `resource/upload/.trash/` 暂存，再删除数据库记录，最后清理暂存文件
- OSS/COS：先软删除数据库记录，再按文件 URL 匹配真实配置删除云端对象；如果云端删除失败，会把对应记录恢复回来
- 批量删除会先做一轮配置匹配预检；本地文件按批次原子暂存，云文件失败时只回滚失败记录，避免出现“记录还在但文件已没”的坏状态

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

当 `multiple=true` 且 `maxCount > 1` 时，文件管理器顶部上传入口也支持一次选择多个文件并并发上传。

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
| `api/upload/file/index.ts` | `/upload/file/upload` + `/upload/file/` | 文件上传 + 文件记录 CRUD |
| `api/upload/dir/index.ts` | `/upload/dir/` | 目录 CRUD + tree |
| `api/upload/config/index.ts` | `/upload/config/` | 配置 CRUD |
| `api/upload/dir_rule/index.ts` | `/upload/dir_rule/` | 目录规则 CRUD |

## 前端页面

| 页面 | 说明 |
|------|------|
| `views/upload/config/index.vue` | 上传配置 CRUD，按存储类型联动显示本地 / OSS / COS 字段 |
| `views/upload/dir/index.vue` | 文件目录树形 CRUD |
| `views/upload/dir_rule/index.vue` | 文件目录规则 CRUD |
| `views/upload/file/index.vue` | 文件记录 CRUD |
| `components/file-manager/index.vue` | 复用文件选择器，左侧目录树 + 右侧文件卡片 + 上传 + 多选 |

## 本地直跑

```bash
cd admin-go
go run ./app/upload/main.go
```

历史记录回填目录归属可用：

```bash
cd admin-go
go run ./app/upload/cmd/backfill_upload_dir --config ./app/upload/manifest/config/config.yaml --dry-run
```
