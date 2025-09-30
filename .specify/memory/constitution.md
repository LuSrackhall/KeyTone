<!--
KeyTone Constitution Update Report:
Version: 1.0.0 (Initial version)
Created: 2025-09-30
New Principles Added:
- I. 代码质量与架构分离 (Code Quality & Architecture Separation)
- II. 测试优先与覆盖率标准 (Test-First & Coverage Standards)  
- III. 用户体验一致性 (User Experience Consistency)
- IV. 性能与响应性要求 (Performance & Responsiveness Requirements)
- V. 跨平台兼容性 (Cross-Platform Compatibility)
Templates Status:
✅ constitution.md (created)
✅ plan-template.md (constitution check section updated with specific principle verification)
⚠ spec-template.md (no changes needed - already principle-agnostic)
✅ tasks-template.md (updated with constitution compliance annotations and quality gates)
Follow-up TODOs:
- No placeholders intentionally deferred
- All principle-based verification points implemented
- Ready for use in development workflow
-->

# KeyTone 项目宪章

## 核心原则

### I. 代码质量与架构分离

所有功能必须基于清晰的架构分层实现：Go 后端负责系统级操作（键盘监听、音频播放），前端负责用户界面和配置管理。

- 代码必须遵循 Go 和 TypeScript 的最佳实践和惯用法
- 每个模块必须具有单一职责，接口清晰定义
- 跨层通信必须通过定义良好的 API 契约（JSON-RPC 或 REST）
- 禁止在前端直接调用系统 API，所有系统交互必须通过 Go 后端
- 代码审查必须验证架构边界未被违反

### II. 测试优先与覆盖率标准（不可妥协）

所有新功能和修改必须采用测试驱动开发（TDD）：

- 必须先编写测试 → 用户批准 → 测试失败 → 然后实现功能
- Go 后端代码覆盖率必须达到 85% 以上
- 前端组件测试覆盖率必须达到 80% 以上
- 音频处理和键盘监听等关键功能必须包含集成测试
- 严格执行红-绿-重构循环
- 每个 PR 必须包含相应的测试代码

### III. 用户体验一致性

确保跨平台和多语言环境下的一致用户体验：

- UI 组件必须遵循统一的设计系统（基于 Quasar 组件）
- 所有用户界面文本必须支持国际化（i18n），至少支持中文和英文
- 错误消息必须用户友好，提供明确的解决方案指引
- 键音配置和用户设置必须在所有平台上保持一致的行为
- 用户反馈（音效、视觉提示）必须即时且直观
- 无障碍访问（ARIA）支持必须在所有交互元素中实现

### IV. 性能与响应性要求

保证实时音频响应和系统资源效率：

- 键盘事件响应延迟必须低于 10ms
- 音频播放延迟必须低于 20ms
- 应用启动时间必须控制在 3 秒以内
- 内存使用必须保持在 100MB 以下（空闲状态）
- CPU 使用率在正常使用下必须低于 5%
- 音频文件加载和缓存策略必须优化，支持预加载常用音效
- 性能关键路径必须包含性能测试和基准测试

### V. 跨平台兼容性

确保在所有目标平台上的功能完整性和稳定性：

- 必须支持 Windows、macOS 和 Linux 主流版本
- 平台特定功能（如系统托盘、开机启动）必须有统一的抽象接口
- 音频格式支持必须在所有平台上保持一致
- 文件路径和权限处理必须跨平台兼容
- 键盘和鼠标事件监听必须适配不同操作系统的差异
- 发布构建必须为每个平台生成优化的二进制文件

## 技术栈约束

### 后端技术要求

- Go 1.21+ 用于核心逻辑实现
- 结构化日志记录（使用标准 log 包或 slog）
- 配置管理必须支持 JSON 格式
- 错误处理必须遵循 Go 的错误处理惯例
- 并发控制必须使用 context 包进行适当的取消和超时处理

### 前端技术要求

- Vue 3 + TypeScript + Quasar Framework
- 状态管理使用 Pinia
- 路由管理使用 Vue Router
- 国际化使用 Vue i18n
- 构建工具使用 Vite/Quasar CLI
- 代码风格必须通过 ESLint 和 Prettier 检查

### 电子应用要求

- Electron 主进程和渲染进程之间的通信必须安全（contextBridge）
- 必须实现适当的安全策略（CSP）
- 自动更新机制必须安全可靠
- 系统集成功能必须通过主进程暴露

## 开发工作流程

### 代码审查要求

- 所有代码必须通过至少一次代码审查
- 审查必须验证架构原则、测试覆盖率和性能要求
- 安全相关代码必须由具有安全背景的审查员审查
- UI/UX 变更必须包含设计审查

### 质量门禁

- 所有自动化测试必须通过
- 代码覆盖率必须达到最低要求
- 静态代码分析不得有高危或中危问题
- 性能基准测试不得出现回归
- 多语言支持必须完整，无遗漏的翻译键

### 发布标准

- 遵循语义化版本控制（SemVer）
- 每个版本必须包含完整的变更日志
- 必须在所有目标平台上完成回归测试
- 用户文档必须与功能变更同步更新

## 治理规则

本宪章优于所有其他开发实践和约定。所有架构决策、代码审查和功能实现都必须符合这些原则。

复杂性必须有充分的理由支持。如果无法简化复杂的方法，必须在设计文档中提供详细的论证。

所有 PR 和代码审查必须验证对这些原则的遵守情况。违反原则的代码不得合并到主分支。

运行时开发指导请参考项目文档和相关的代理特定指导文件。

**版本**: 1.0.0 | **批准日期**: 2025-09-30 | **最后修订**: 2025-09-30