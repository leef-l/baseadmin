# EasyMVP对接OpenHands与Aider引擎设计实现文档

本文档定义仓库在引入 OpenHands / Aider 类 AI 执行引擎时的基础约定，当前仓库也沿用同一套目录和入口规范。

## 目标

- 统一 AI 代理读取项目的入口文档
- 统一 Docker 开发目录布局
- 统一环境变量来源，避免多份 `.env` 漂移
- 统一文档目录，保证引用路径可点击、可追踪

## 目录原则

```text
docker/
├── dev/
│   ├── .env
│   ├── compose.ps1
│   ├── compose.sh
│   ├── docker-compose.yml
│   └── docker-compose.cn.yml
└── prod/
    └── docker-compose.yml

docs/

docker/build/
docker/mysql/
docker/nginx/
```

规则如下：

1. `docker/` 负责 compose 入口与启动脚本
2. `docker/build/`、`docker/mysql/`、`docker/nginx/` 负责镜像构建资源与初始化资源
3. `docs/` 负责长期文档沉淀
4. 根目录只保留简短入口说明，不堆放大段文档

## 环境变量原则

- 开发环境唯一源文件是 `docker/dev/.env`
- 执行 `docker/dev/compose.ps1` 或 `docker/dev/compose.sh` 时，必须先同步到 `admin-go/.env`
- `admin-go/.env` 视为运行时副本，不作为长期维护入口

## 文档引用原则

- 文档之间尽量使用 Markdown 相对路径链接
- `CLAUDE.md` 必须先引导阅读本文档
- 所有关键入口都应能在 IDE 中直接点击打开

## 对当前仓库的落地要求

- Docker Compose 入口放在根目录 `docker/`
- 详细说明放在根目录 `docs/`
- `CLAUDE.md`、`README.MD` 只保留索引和关键引导
