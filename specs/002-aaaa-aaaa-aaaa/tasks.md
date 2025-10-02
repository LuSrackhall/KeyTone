# Tasks: 键音专辑签名系统

**Input**: 设计文档来自 `/specs/002-aaaa-aaaa-aaaa/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. 从特性目录加载 plan.md
   → 如未找到: ERROR "未找到实现计划"
   → 提取: 技术栈、库、结构
2. 加载可选设计文档:
   → data-model.md: 提取实体 → 模型任务
   → contracts/: 每个文件 → 契约测试任务
   → research.md: 提取决策 → 设置任务
3. 按类别生成任务:
   → Setup: 项目初始化、依赖、代码规范
   → Tests: 契约测试、集成测试
   → Core: 模型、服务、端点
   → Integration: 数据库、中间件、日志
   → Polish: 单元测试、性能、文档
4. 应用任务规则:
   → 不同文件 = 标记 [P] 并行执行
   → 相同文件 = 顺序执行 (无 [P])
   → 测试优先于实现 (TDD)
5. 按顺序编号任务 (T001, T002...)
6. 生成依赖关系图
7. 创建并行执行示例
8. 验证任务完整性:
   → 所有契约都有测试?
   → 所有实体都有模型?
   → 所有端点都已实现?
9. 返回: SUCCESS (任务准备就绪)
```

## Format: `[ID] [P?] 描述`
- **[P]**: 可并行运行 (不同文件，无依赖)
- 描述中包含准确的文件路径

## Path Conventions
- **Web 应用**: `sdk/` (后端 Go), `frontend/src/` (前端 Vue)
- 后端路径: `sdk/signature/`, `sdk/server/`
- 前端路径: `frontend/src/components/`, `frontend/src/services/`, `frontend/src/pages/`
- 测试路径: `sdk/*_test.go` (后端契约测试), `frontend/tests/e2e/` (前端烟雾测试)

## Phase 3.1: Setup (设置)
- [ ] T001 根据实现计划检查项目结构 (sdk/ 和 frontend/ 已存在)
- [ ] T002 [P] 在 sdk/ 安装/更新 Go 依赖 (gin, nanoid, crypto)
- [ ] T003 [P] 在 frontend/ 安装/更新前端依赖 (axios, Vue Router, Pinia)
- [ ] T004 [P] 配置后端路由前缀 /sdk/ 与错误响应中间件 (sdk/server/server.go)
- [ ] T005 [P] 配置前端 axios baseURL 动态对齐逻辑 (frontend/src/boot/axios.ts 已存在，验证 UpdateApi 功能)

## Phase 3.2: Tests First (TDD) ⚠️ 必须在 3.3 之前完成
**关键: 这些测试必须先编写且必须失败，然后才能进行任何实现**
**宪章合规: 测试优先与覆盖率标准 - 契约与核心逻辑 TDD**

### 后端契约测试 (Contract Tests)
- [ ] T006 契约测试 GET /sdk/signatures (sdk/signature/signature_test.go)
- [ ] T007 契约测试 POST /sdk/signatures (sdk/signature/signature_test.go)
- [ ] T008 契约测试 DELETE /sdk/signatures/{name} (sdk/signature/signature_test.go)
- [ ] T009 契约测试 POST /sdk/signatures/{name}/export (sdk/signature/signature_test.go)
- [ ] T010 契约测试 POST /sdk/signatures/import (sdk/signature/signature_test.go)
- [ ] T011 [P] 契约测试 GET /sdk/albums/{albumId}/signatures (sdk/server/album_test.go)
- [ ] T012 [P] 契约测试 POST /sdk/albums/{albumId}/sign (sdk/server/album_test.go)
- [ ] T013 [P] 契约测试 POST /sdk/albums/{albumId}/export (sdk/server/album_test.go)
- [ ] T014 [P] 契约测试 GET /ping (sdk/server/server_test.go)

### 前端烟雾测试 (Smoke Tests)
- [ ] T015 [P] Electron 烟雾测试: 应用启动与后端连接 (frontend/tests/e2e/electron-smoke.spec.ts)
- [ ] T016 [P] Electron 烟雾测试: 签名管理对话框可访问 (frontend/tests/e2e/electron-smoke.spec.ts)

## Phase 3.3: Core Implementation (核心实现) - 仅在测试失败后
**宪章合规: 代码质量与架构分离 - 清晰的层次边界**

### 后端实现
- [ ] T017 创建加密工具模块 (sdk/crypto/crypto.go): AES-256-GCM 加密/解密、Base64 编码、SHA-256 哈希
- [ ] T018 [P] DigitalSignature 数据结构与验证 (sdk/signature/model.go)
- [ ] T019 [P] SignatureFile 数据结构 (.ktsign 格式) (sdk/signature/file.go)
- [ ] T020 全局配置签名存储层 (sdk/signature/storage.go): 加密 protectCode 为 key、明文 JSON 为 value
- [ ] T021 专辑配置签名存储层 (sdk/signature/album_storage.go): 双重加密 (key 和 value)
- [ ] T022 签名管理服务 (sdk/signature/service.go): CRUD、导入/导出逻辑
- [ ] T023 GET /sdk/signatures 端点实现 (sdk/server/routes.go)
- [ ] T024 POST /sdk/signatures 端点实现 (sdk/server/routes.go)
- [ ] T025 DELETE /sdk/signatures/{id} 端点实现 (sdk/server/routes.go)
- [ ] T026 POST /sdk/signatures/{id}/export 端点实现 (sdk/server/routes.go)
- [ ] T027 POST /sdk/signatures/import 端点实现 (sdk/server/routes.go)
- [ ] T028 专辑签名服务 (sdk/signature/album_service.go): 签名关联与导出逻辑
- [ ] T029 GET /sdk/albums/{albumId}/signatures 端点实现 (sdk/server/routes.go)
- [ ] T030 POST /sdk/albums/{albumId}/sign 端点实现 (sdk/server/routes.go)
- [ ] T031 POST /sdk/albums/{albumId}/export 端点实现 (sdk/server/routes.go)
- [ ] T032 输入验证与用户友好错误消息 (sdk/server/middleware.go)

### 前端实现
- [ ] T033 [P] 签名类型定义 (frontend/src/types/signature.ts)
- [ ] T034 [P] 签名 API 服务封装 (frontend/src/services/signatureApi.ts): axios 调用封装
- [ ] T035 [P] 专辑签名 API 服务封装 (frontend/src/services/albumSigningApi.ts)
- [ ] T036 签名管理 Pinia Store (frontend/src/stores/signatureStore.ts): 状态管理与缓存
- [ ] T037 签名列表组件 (frontend/src/components/SignatureList.vue): 展示、选择、删除、名片占位符处理
- [ ] T038 签名创建对话框组件 (frontend/src/components/SignatureCreateDialog.vue): 表单与名片图片上传、ARIA 标签
- [ ] T039 签名导入对话框组件 (frontend/src/components/SignatureImportDialog.vue): 文件选择与覆盖确认、ARIA 标签
- [ ] T040 签名导出功能 (frontend/src/components/SignatureList.vue): 调用导出 API 与下载（依赖 T037）
- [ ] T041 签名管理对话框 (frontend/src/components/SignatureManagementDialog.vue): 集成所有签名操作、按钮位于删除按钮右侧、ARIA 标签
- [ ] T042 专辑签名选择对话框 (frontend/src/components/AlbumSigningDialog.vue): 签名选择、筛选/搜索功能、二次导出设置、ARIA 标签
- [ ] T043 专辑导出流程集成 (frontend/src/pages/AlbumExportPage.vue): 集成签名选择与导出调用
- [ ] T044 错误处理与国际化提示 (frontend/src/boot/axios.ts 响应拦截器): 导入失败、网络错误、验证失败、签名冲突等场景

## Phase 3.4: Integration (集成)
**宪章合规: 性能与响应性要求 + 跨平台兼容性**
- [ ] T045 Electron 主进程端口发现逻辑 (frontend/src-electron/electron-main.ts): 捕获 KEYTONE_PORT 并通过 IPC 暴露
- [ ] T046 渲染进程端口获取与 axios 同步 (frontend/src/boot/axios.ts): 调用 getBackendPort 与 UpdateApi
- [ ] T047 [P] 结构化日志集成 (sdk/logger/logger.go): 签名操作审计日志
- [ ] T048 [P] 跨平台文件路径处理 (sdk/signature/path.go): Windows/macOS/Linux 路径标准化
- [ ] T049 资源使用监控验证 (自动化或手动基准测试: 空闲内存 <100MB, CPU <5%)

## Phase 3.5: Polish & Quality Gates (完善与质量门禁)
**宪章合规: 所有原则验证**
- [ ] T050 [P] 签名验证单元测试 (sdk/signature/validation_test.go)
- [ ] T051 [P] 加密/解密单元测试 (sdk/crypto/crypto_test.go)
- [ ] T052 [P] 组件单元测试: SignatureList (frontend/tests/unit/SignatureList.spec.ts)
- [ ] T053 [P] 组件单元测试: SignatureCreateDialog (frontend/tests/unit/SignatureCreateDialog.spec.ts)
- [ ] T054 [P] 代码覆盖率验证 (后端 ≥85% 符合宪章，前端 ≥80% 符合宪章)
- [ ] T055 [P] 国际化与无障碍验证 (frontend/src/i18n/: 无缺失 key; 所有交互元素有 ARIA 标签)
- [ ] T056 [P] 性能基准测试 (签名导入/导出 <1s, 专辑签名 <2s, 应用启动 <3s)
- [ ] T057 跨平台手动测试 (Windows/macOS/Linux: 执行 quickstart.md 全部 6 个场景验证)
- [ ] T058 更新 API 文档 (specs/002-aaaa-aaaa-aaaa/contracts/: 补充示例与错误码)
- [ ] T059 代码审查检查清单验证 (去重、架构分离、REST 边界)
- [ ] T060 运行 quickstart.md 端到端流程验证

## Dependencies (依赖关系)
- Setup (T001-T005) 在所有任务之前
- 契约测试 (T006-T014) 在后端实现 (T017-T032) 之前
- 烟雾测试 (T015-T016) 在前端实现 (T033-T044) 之前
- T017 (加密工具) 阻塞 T020, T021 (存储层)
- T020, T021 (存储层) 阻塞 T022 (签名服务)
- T022 (签名服务) 阻塞 T023-T027 (签名端点)
- T028 (专辑签名服务) 阻塞 T029-T031 (专辑端点)
- T033 (类型定义) 阻塞 T034-T035 (API 服务)
- T034-T035 (API 服务) 阻塞 T036 (Store)
- T036 (Store) 阻塞 T037-T043 (UI 组件)
- T037 (列表组件) 阻塞 T040 (导出功能，同文件)
- 核心实现 (T017-T044) 在集成 (T045-T049) 之前
- 集成 (T045-T049) 在完善 (T050-T060) 之前

## Parallel Example (并行执行示例)
```bash
# 同时启动 T006-T014 (后端契约测试):
Task: "契约测试 GET /sdk/signatures (sdk/signature/signature_test.go)"
Task: "契约测试 POST /sdk/signatures (sdk/signature/signature_test.go)"
Task: "契约测试 DELETE /sdk/signatures/{id} (sdk/signature/signature_test.go)"
Task: "契约测试 POST /sdk/signatures/{id}/export (sdk/signature/signature_test.go)"
Task: "契约测试 POST /sdk/signatures/import (sdk/signature/signature_test.go)"
Task: "契约测试 GET /sdk/albums/{albumId}/signatures (sdk/server/album_test.go)"
Task: "契约测试 POST /sdk/albums/{albumId}/sign (sdk/server/album_test.go)"
Task: "契约测试 POST /sdk/albums/{albumId}/export (sdk/server/album_test.go)"
Task: "契约测试 GET /ping (sdk/server/server_test.go)"

# 同时启动 T018-T019 (数据结构):
Task: "DigitalSignature 数据结构与验证 (sdk/signature/model.go)"
Task: "SignatureFile 数据结构 (sdk/signature/file.go)"

# 同时启动 T033-T035 (前端类型与服务):
Task: "签名类型定义 (frontend/src/types/signature.ts)"
Task: "签名 API 服务封装 (frontend/src/services/signatureApi.ts)"
Task: "专辑签名 API 服务封装 (frontend/src/services/albumSigningApi.ts)"
```

## Notes (注意事项)
- [P] 任务 = 不同文件，无依赖
- 在实现前验证测试失败
- 每个任务后提交代码
- 避免: 模糊任务、同一文件冲突
- 前端 UI 测试采用混合策略: 先实现 → 少量烟雾测试 → 补充关键组件单测
- 后端遵循契约优先: 测试先行 → 实现 → 覆盖率验证

## Task Generation Rules (任务生成规则)
*在 main() 执行期间应用*

1. **从 Contracts**:
   - 每个契约端点 → 契约测试任务 [P]
   - 每个端点 → 实现任务
   
2. **从 Data Model**:
   - 每个实体 → 模型创建任务 [P]
   - 存储层 → 加密与持久化任务
   
3. **从 User Stories**:
   - quickstart.md 场景 → 集成测试 [P]
   - 手动验证步骤 → 自动化验证任务

4. **排序**:
   - Setup → Tests → Models → Services → Endpoints → Integration → Polish
   - 依赖关系阻止并行执行

## Validation Checklist (验证清单)
*门禁: 在 main() 返回前检查*

- [x] 所有契约都有对应的测试 (T006-T014)
- [x] 所有实体都有模型任务 (T018-T019)
- [x] 所有测试都在实现之前 (T006-T016 在 T017-T044 之前)
- [x] 并行任务真正独立 (不同文件，无依赖)
- [x] 每个任务指定准确的文件路径
- [x] 没有任务修改与另一个 [P] 任务相同的文件
