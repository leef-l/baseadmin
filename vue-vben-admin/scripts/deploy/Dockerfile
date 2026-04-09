FROM node:22-slim AS builder

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
ENV TZ=Asia/Shanghai
ENV CI=true
ENV PNPM_NETWORK_CONCURRENCY=3
ENV PNPM_CHILD_CONCURRENCY=1
ENV FRONTEND_INSTALL_NODE_MAX_OLD_SPACE_SIZE=768
ENV FRONTEND_NODE_MAX_OLD_SPACE_SIZE=1024

RUN npm i -g corepack

WORKDIR /app

# copy package.json and pnpm-lock.yaml to workspace
COPY . /app

# 安装依赖
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    npm_config_node_options="--max-old-space-size=${FRONTEND_INSTALL_NODE_MAX_OLD_SPACE_SIZE}" \
    pnpm install \
      --frozen-lockfile \
      --network-concurrency="${PNPM_NETWORK_CONCURRENCY}" \
      --child-concurrency="${PNPM_CHILD_CONCURRENCY}"
RUN pnpm run build --filter=\!./docs

RUN echo "Builder Success 🎉"

FROM nginx:stable-alpine AS production

# 配置 nginx
RUN echo "types { application/javascript js mjs; }" > /etc/nginx/conf.d/mjs.conf \
    && rm -rf /etc/nginx/conf.d/default.conf

# 复制构建产物
COPY --from=builder /app/playground/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY --from=builder /app/scripts/deploy/nginx.conf /etc/nginx/nginx.conf

EXPOSE 8080

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
