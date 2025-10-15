# 项目上下文

## 目的

KeyTone 是一款实用的键盘音效模拟软件，旨在为用户在需要保持安静的环境下提供实时的按键声音反馈，提升使用电脑的舒适度和个性化体验。

**核心目标：**

- 提供高性能的实时按键音效触发系统
- 支持高度可定制的键音专辑系统
- 释放用户创造力，打造独一无二的专属键音体验
- 遵循 GPL 开源协议，鼓励社区参与和分享

## 技术栈

### 核心后端 (SDK)

- **Go 1.24.5** - 高性能核心引擎
  - `github.com/gin-gonic/gin` - HTTP 服务器框架
  - `github.com/gopxl/beep/v2` - 音频播放引擎
  - `github.com/robotn/gohook` - 键盘/鼠标事件监听
  - `github.com/fsnotify/fsnotify` - 文件系统监控

### 前端应用 (Frontend)

- **Vue 3.5.14** - 渐进式前端框架
- **Quasar 2.17.1** - Vue 组件框架
- **TypeScript 4.9.5** - 类型安全的 JavaScript 超集
- **Electron 31.7.7** - 跨平台桌面应用框架
- **UnoCSS 0.58.9** - 原子化 CSS 引擎
- **Pinia 2.0.11** - Vue 状态管理
- **Vue Router 4.0** - 前端路由
- **Vue I18n 9.14.4** - 国际化支持
  - **Axios 1.2.1** - HTTP 客户端

### 文档系统 (Docs)

- **VitePress 1.2.3** - 基于 Vite 的静态站点生成器
  - **Vue 3.2** - 文档组件框架

### 构建与工具

- **Vite** - 前端构建工具 (通过 @quasar/app-vite)
- **Electron Builder 24.13.3** - 应用打包工具
- **ESLint 8.10.0** - 代码质量检查
- **Prettier 2.5.1** - 代码格式化

## 项目约定

### 代码风格

- **命名规范：**
  - Go: 驼峰命名，导出符号首字母大写
  - TypeScript/Vue: 驼峰命名，组件使用 PascalCase
  - 文件名: kebab-case (如 `key-event.go`, `audio-package.ts`)
  
- **格式化：**
  - 使用 Prettier 统一格式化，配置见根目录
  - 所有代码文件应包含 GPL 许可证头部声明
  
- **ESLint 规则：**
  - 启用 TypeScript ESLint
  - Vue 推荐规则集
  - 与 Prettier 集成，避免冲突

### 架构模式

**混合架构 (Hybrid Architecture):**

1. **Go 后端服务 (SDK)**
   - 独立的 HTTP 服务器 (Gin)
   - 负责键盘/鼠标事件监听
   - 音频播放引擎
   - 配置文件管理
   - 音频包管理

2. **Electron + Vue 前端**
   - Electron 主进程与渲染进程分离
   - Vue 3 Composition API
   - Quasar 组件库提供 UI 组件
   - 通过 Axios/EventSource 与 Go 后端通信

3. **模块化设计：**

   ```text
   sdk/
   ├── audioPackage/    # 音频包配置与列表管理
   ├── config/          # 应用配置管理
   ├── keyEvent/        # 键盘事件处理
   ├── keySound/        # 音效播放逻辑
   ├── logger/          # 日志系统
   └── server/          # HTTP 服务器
   ```

4. **通信模式：**
   - 前端通过 HTTP API 与 Go 后端交互
   - 使用 EventSource (SSE) 实现实时事件推送

### 测试策略

- **当前状态：** 测试脚本返回 "No test specified"
- **推荐策略：**
  - Go 后端：使用 Go 标准测试框架
  - 前端：Vue Test Utils + Vitest (待实施)
  - E2E 测试：Playwright/Spectron (待实施)

### Git 工作流

**提交规范 (Commitlint):**
- 遵循 Conventional Commits 规范
- 使用 `@commitlint/config-conventional`
- Header 最大长度：150 字符 (扩展配置)
- 支持 Commitizen 交互式提交

**分支策略：**
- `main` - 主分支，稳定版本
- `openSpec` - 当前开发分支 (OpenSpec 规范)
- 功能分支：使用描述性名称

**Husky 钩子：**
- Pre-commit: ESLint 检查
- Commit-msg: Commitlint 验证

## 领域上下文

### 键音专辑系统

**核心概念：**
- **音频源文件 (Audio Source):** 用户自行提供的音频文件
- **裁剪定义声音 (Sound Clipping & Definition):** 用户可对导入的音频源文件进行裁剪, 裁剪后的声音是可以独立绑定到按键且不影响音频源文件的使用的，指定起止时间，提取所需片段作为按键音效。支持精确到毫秒的裁剪、音量增减的设置，便于个性化定制每个声音反馈(毕竟音频源文件有时在片段长度和音量上无法达到需求, 这个步骤可以在一定程度上解决这些问题)
- **高级声音 (Advanced Sound):** 多个音频组合，支持随机或顺序播放, 以及任意层级的嵌套
- **键音绑定 (Key Binding):** 将声音绑定到特定按键
- **键音专辑 (Album):** 完整的键音配置集合，可导入导出分享

**设计理念：**
- **不提供音频文件** - 鼓励用户创作和收集
- **高度组合性** - 高级声音支持嵌套、继承、组合
- **社区驱动** - 通过 itch.io 社区分享键音专辑

### 按键触发逻辑

- 基于真实按键触发状态设计
- 仅在按下和抬起瞬间触发音效
- 长按不重复播放

### 支持的输入设备

- 键盘 (所有按键)
- 鼠标 (按键和滚轮)

## 重要约束

### 技术约束

1. **性能要求：** 必须实现低延迟的实时音效触发 (< 50ms)
2. **Go 版本：** 最低 Go 1.24.5
3. **Node 版本：** 最低 Node ^18
4. **浏览器目标：** ES2019+, Edge 88+, Firefox 78+, Chrome 87+, Safari 13.1+

### 许可证约束

- **GPL-3.0** - 所有源代码必须遵循 GPL v3 协议
- 所有代码文件头部必须包含 GPL 声明
- 不得包含版权受限的音频资源

### 平台约束

- **主要支持：** Windows (已上架 Microsoft Store)
- **潜在支持：** macOS, Linux
- **不支持：** 移动平台

### 业务约束

- **免费开源** - 不接受付费功能
- **无音频提供** - 不在项目中包含或分发音频文件
- **社区驱动** - 鼓励用户在 itch.io 社区分享

## 外部依赖

### 官方资源

- **官方网站：** <https://keytone.xuanhall.com>
- **文档站点：** <https://keytone.xuanhall.com> (VitePress)
- **GitHub 仓库：** <https://github.com/LuSrackhall/KeyTone>
- **itch.io 页面：** <https://lusrackhall.itch.io/keytone>
- **Microsoft Store：** <https://apps.microsoft.com/store/detail/9NGKDXHPGJXD>

### 关键外部库

- **beep** - 音频播放的核心依赖，历史上存在内存泄漏问题 (已协助其作者进行解决修复)
- **robotn/gohook** - 系统级键盘/鼠标事件监听
- **Electron** - 跨平台桌面框架

### 开发工具

- **OpenSpec** - 本项目使用的规范驱动开发流程, 且规范文档尽量使用中文编写
- **Commitizen** - 交互式提交工具
- **Husky** - Git hooks 管理
