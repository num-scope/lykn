# Lykn

Python 项目授权平台。管理 license 的签发与验证。

## 项目结构

- `backend/` — Go 管理平台后端
- `sdk/python/` — Python 验证 SDK
- `frontend/` — 前端（Antdv Next）

## Backend 开发

1. 复制 `backend/config/config.yaml.example` 为 `backend/config/config.yaml`
2. 启动 PostgreSQL，或使用 `docker-compose up postgres`
3. 进入 `backend/` 目录运行 `make run`
4. 默认用户：`admin / admin123`（仅开发环境，启动后请立即修改）

常用命令：

```bash
cd backend
make help
make test
make build
make demo
```

`make demo` 会重新生成跨语言测试 fixture：`tests/fixtures/public.pem`、`tests/fixtures/license.lic` 和完整的 `tests/fixtures/license.json`。

### 已提供接口

- `POST /api/v1/auth/login`
- `GET/POST/PUT/DELETE /api/v1/features`
- `GET/POST/PUT/DELETE /api/v1/plans`
- `GET/POST/PUT/DELETE /api/v1/projects`
- `GET /api/v1/projects/:id/public-key`
- `GET/POST /api/v1/projects/:id/licenses`
- `GET /api/v1/licenses/:id`
- `GET /api/v1/licenses/:id/download`

## Frontend 开发

1. 复制 `frontend/.env.example` 为 `frontend/.env`
2. 进入 `frontend/` 目录运行 `pnpm dev`
3. `frontend/.env` 必须提供 `VITE_API_BASE`，代码不会为该变量设置默认值；跨域由后端 CORS 处理

## Docker Compose

当前 compose 只覆盖后端链路，故意不包含 `frontend/`。

```bash
docker-compose up -d --build
docker-compose logs -f server
docker-compose down -v
```

默认服务：

- `postgres`：`postgres:16`
- `server`：Go 管理后端，监听 `http://127.0.0.1:8080`

本地 `make run` 使用 `backend/config/config.yaml`。Compose 使用 `backend/config/config.compose.yaml`，后端所有配置都从 YAML 读取。
