# Shadcn 管理后台瘦身版

本文件记录当前项目的瘦身范围、保留能力和后续扩展建议

## 当前目标

项目被裁剪为一个后台基础起点：

- 访问 `/` 时，如果未登录，会自动跳转到 `/sign-in`
- 登录成功后进入仪表盘
- 侧边栏只保留仪表盘菜单
- 主内容区不再限制为 `7xl` 最大宽度，默认铺满可用空间
- 顶部用户菜单只保留退出登录
- 退出登录后回到 `/sign-in`

## 保留模块

- 登录、注册、忘记密码、一次性密码等基础认证页面
- 登录态存储与路由守卫
- 后台布局壳：侧边栏、顶部栏、主内容区
- 主内容区基础间距，移除大屏居中和最大宽度限制
- 仪表盘示例页
- 主题切换、方向切换、布局设置和搜索命令面板
- 基础 UI 组件、图标和全局样式
- Tailwind CSS v4 与 shadcn 主题变量
- Vitest 与 Playwright 浏览器测试配置

## 已删除模块

- 任务管理模块
- 用户管理模块
- 应用集成模块
- 聊天模块
- 设置页面模块
- 错误页路由示例
- 表格业务封装
- Clerk 集成示例
- 帮助中心示例页
- 即将上线、了解更多等占位组件
- TanStack Table 相关业务封装

## 目录说明

- `src/features/auth`：认证相关页面和表单
- `src/features/dashboard`：登录后的后台首页模板
- `src/components/layout`：后台布局骨架
- `src/components/ui`：基础 UI 组件
- `src/styles/index.css`：全局样式入口
- `src/routes`：当前保留的路由定义
- `src/stores`：登录态存储
- `src/lib`：通用工具函数

## 后续扩展建议

1. 先按业务域新增 `src/features/<module>`
2. 再为新模块补充 `src/routes/_authenticated/<module>/index.tsx`
3. 需要表格时再引入 TanStack Table 或恢复表格封装
4. 需要组织体系时再单独新增组织切换、组织创建和成员权限模块
