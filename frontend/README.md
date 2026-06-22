# Lykn Frontend

Lykn 授权平台前端，基于 Soybean Admin、Vue 3、Vite、TypeScript、Antdv Next 和 UnoCSS

## Scripts

```bash
pnpm install
pnpm dev
pnpm typecheck
pnpm build
```

## Backend

默认后端接口为 `/api/v1`，本地开发可在 `.env.test` 或本机 `.env` 中调整：

```bash
VITE_API_BASE=http://127.0.0.1:8080/api/v1
VITE_SERVICE_BASE_URL=http://127.0.0.1:8080/api/v1
```
