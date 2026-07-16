# Shadcn 管理后台瘦身版

这是基于 Shadcn、Vite、React 和 TanStack Router 的后台基础模板，当前项目已经裁剪为登录 + 仪表盘的轻量起点

![项目截图](public/images/shadcn-admin.png)

## 功能

- 登录、注册、忘记密码、一次性密码等基础认证页面
- 登录后进入仪表盘骨架页
- 侧边栏仅保留仪表盘菜单
- 用户菜单仅保留退出登录
- 保留主题切换、布局设置、方向切换和基础样式能力
- 保留 shadcn/ui 基础组件与 Tailwind CSS v4 样式变量

## 技术栈

- React
- Vite
- TypeScript
- TanStack Router
- Tailwind CSS v4
- shadcn/ui
- Zustand
- Vitest + Playwright

## 本地运行

安装依赖：

```bash
pnpm install
```

启动开发服务：

```bash
pnpm run dev
```

执行校验：

```bash
pnpm run format:check
pnpm run lint
pnpm run build
pnpm test
```

## 瘦身说明

详细裁剪范围和后续扩展建议见 `SLIM.md`

## 许可

本项目沿用 MIT 许可证
