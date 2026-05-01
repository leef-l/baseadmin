# fund-disk H5 用户端

基于 React 18 + TypeScript + Vite + antd-mobile 的会员 H5 站点，对接 `admin-go/app/member` 的 `/api/member-portal/*` 接口。

## 目录结构

```
h5-react/
  src/
    api/                     # 后端接口封装
    components/
      SmsCodeButton.tsx      # 验证码倒计时按钮
      layout/
        TabLayout.tsx        # 底部 Tab
        PageHeader.tsx       # 通用顶栏
    pages/
      auth/                  # 登录 / 注册 / 找回密码
      home/                  # 首页（轮播+等级+三钱包+推荐+热门寄售）
      mall/                  # 商城（列表+详情+订单）
      warehouse/             # 仓库（市场+我的库存+我的交易）
      me/                    # 我的（资料/改密码/换手机/邀请码）
      wallet/                # 三钱包 + 流水
      team/                  # 团队
    router/                  # 路由 + 鉴权守卫
    stores/                  # zustand 全局状态
    styles/                  # 全局 CSS + Tailwind
    utils/                   # 通用工具
```

## 快速开始

```bash
cd h5-react

# 1. 安装依赖（仓库根目录铁律：禁止本机直接 npm/pnpm，请走 scripts/run-node-task-with-limits.sh）
../scripts/run-node-task-with-limits.sh npm install

# 2. 配置环境变量
cp .env.example .env
# 修改 VITE_API_PROXY_TARGET 指向后端 member 服务（默认 http://127.0.0.1:10027）

# 3. 启动开发服务器
../scripts/run-node-task-with-limits.sh npm run dev
# 访问 http://本机IP:5173
```

## 生产构建

```bash
../scripts/run-node-task-with-limits.sh npm run build
# 产物在 dist/，由 Nginx 托管，/api 反代到后端 member 服务
```

## 部署：Nginx 配置示例

```nginx
server {
    listen 80;
    server_name h5.your-domain.com;

    root /www/wwwroot/project/fund-disk/h5-react/dist;
    index index.html;

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|svg|woff2?|ttf)$ {
        expires 7d;
        add_header Cache-Control "public, max-age=604800, immutable";
    }

    # SPA 回退
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 反向代理后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:10027;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 技术栈

- **React 18** + **TypeScript 5** + **Vite 5**
- **antd-mobile 5**：移动端组件库
- **Tailwind CSS 3**：原子化样式（`preflight: false` 避免与 antd-mobile 冲突）
- **react-router-dom 6**：路由
- **zustand 5**：全局状态（轻量、无 boilerplate）
- **axios**：HTTP（统一拦截器，自动带 JWT，401 跳登录）
- **swiper / qrcode.react**：轮播 + 二维码

## 关键约定

- **接口基址**：所有请求走 `/api/member-portal/*`，开发时由 vite proxy 转发到后端，生产由 Nginx 反代。
- **JWT**：登录后存 `localStorage.fd_member_token`，axios 拦截器自动加 `Authorization: Bearer <token>`。
- **邀请码**：注册必填，URL `?invite=XXX` 自动填充，输入后实时调 `/auth/invite-preview` 显示邀请人头像昵称。
- **金额**：后端接口返回的均是元字符串（保留两位小数），前端直接展示，不做换算。
- **单位元换算**：仅 `teamTurnover` 字段是分（int64），前端除以 100 再展示。

## 与后端联调

- 后端 member 服务默认监听 `:10027`
- 启动后端：`cd admin-go && go run app/member/main.go`
- 启动 H5：`cd h5-react && npm run dev`，访问 `http://本机IP:5173`
- 推荐手机访问真机调试；用 Chrome DevTools 设备模拟器也可。
