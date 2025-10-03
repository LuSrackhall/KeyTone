# Tasks: 键音专辑签名系统（最小新端点 + 最大复用）

**Input**: 来自 `/specs/002-aaaa-aaaa-aaaa/` 的最新 spec 与 plan
**Prerequisites**: plan.md (required), data-model.md, contracts/, research.md（可选）

## 执行规则（精简务实版）
- 以复用为先：签名列表 CRUD 走 `/store/*` 与 `/keytone_pkg/*`；实时更新走 `/stream`。
- 仅为 .ktsign 导入/导出与导出签名桥新增最小端点，并提供契约与最小测试。
- 测试策略（与宪章目标一致的务实落地）：优先用测试保护“既有功能”（回归/契约），对“新增最小端点/逻辑”提供必要的契约/单测，允许实现与测试交错推进；以“测试通过”为门槛而非书写时序。
- 窗口与 i18n 约束严格遵守 spec 的 NFR。

## Format: `[ID] [P?] 描述`
- [P] 表示与其它任务可并行（不同文件、无直接依赖）。
- 描述尽量包含准确文件路径（相对仓库根）。

## 路径约定
- 后端：`sdk/server/`（路由与处理器），`sdk/signature/`（如需最小模型与文件编解码）。
- 前端：`frontend/src/components/`, `frontend/src/pages/`, `frontend/src/services/`, `frontend/src/stores/`, `frontend/src/i18n/`。
- 测试：`sdk/*_test.go`, `frontend/tests/e2e/`, `frontend/tests/unit/`。

---

## Phase 3.1: Setup & Audit（建立与审计）
- [x] T001 审计后端现有接口：确认存在 `/store/get|set`、`/keytone_pkg/get|set|delete`、`/stream`（sdk/server/server.go）
- [x] T002 审计前端 SSE：确认 Electron 主进程与渲染进程已订阅 `message` 与 `messageAudioPackage` 并处理 `get_all_value`（frontend/src-electron/electron-main.ts, frontend/src/App.vue, frontend/src/components/Keytone_album.vue）
- [x] T003 校验 axios 动态端口逻辑可用（frontend/src/boot/axios.ts 与主进程桥）
- [x] T004 清点并标注需废弃的旧端点或 UI（如存在 /sdk/signatures* 痕迹，先标记不再使用，保留回滚锚点）

## Phase 3.2: Contracts（仅新增最小端点）

- [ ] T010 新建契约：POST `/signature/export`（导出 .ktsign） → `specs/002-aaaa-aaaa-aaaa/contracts/signature-export.md`
- [ ] T011 新建契约：POST `/signature/import`（导入 .ktsign） → `specs/002-aaaa-aaaa-aaaa/contracts/signature-import.md`
- [ ] T012 新建契约：POST `/export/sign-bridge`（导出流程签名桥） → `specs/002-aaaa-aaaa-aaaa/contracts/export-sign-bridge.md`
- [ ] T013 契约测试（最小集）：为 T010~T012 各写 1 条 happy path + 1 条错误路径（sdk/server/server_test.go 或相关 *_test.go）
- [ ] T014 契约补充：在 `export-sign-bridge.md` 明确 `album_signatures` 中“每次导出的时间戳数组”字段结构与合并规则

## Phase 3.3: Backend（最小实现）

- [x] T020 `.ktsign` 文件编解码最小实现（sdk/signature/file.go）：包含签名名称、唯一标识、介绍、名片资源引用、完整性校验字段（Stage 1）
- [x] T021 导出端点处理器：POST `/signature/export`（sdk/server/signature_handlers.go），从 `/store/get` 读取 `signature_manager`，根据请求选择签名，生成 .ktsign 输出（字节流或保存路径）
- [x] T022 导入端点处理器：POST `/signature/import`（sdk/server/signature_handlers.go），解析上传 .ktsign，写回 `/store/set` 的 `signature_manager`，触发全局 SSE 刷新
- [x] T023 导出签名桥：POST `/export/sign-bridge`（sdk/server/album_handlers.go 或新建 export_handlers.go），将前端导出流程中选择的签名写入专辑配置 `album_signatures`（经 `/keytone_pkg/set`）并返回导出继续所需数据
- [ ] T024 错误与边界：空列表、重复签名（覆盖/取消）、坏文件、路径异常、跨平台路径
- [ ] T025 资源落盘与引用维护：.ktsign 导入/导出时，名片图片等资源复制至资源目录，并在 `signature_manager` 与 `album_signatures` 中维护/修正引用路径（sdk/signature/file.go + 相关 handlers）
- [ ] T026 导出时间戳合并：在导出签名桥中为对应签名追加导出时间戳，保证为数组并去重/排序（按时间）

## Phase 3.4: Frontend（复用读写 + SSE）

- [x] T030 签名管理对话框：首次打开 GET `/store/get` 读取 `signature_manager`；创建/导入/删除后由 SSE 自动刷新（frontend/src/components/SignatureManagementDialog.vue）
- [x] T031 创建与删除：通过 `/store/set` 更新 `signature_manager`；保护码由前端创建时自动生成且 UI 不展示（frontend/src/services/ 或 stores/）
- [x] T032 导出签名文件（.ktsign）：调用 POST `/signature/export`，并处理保存（前端桥接到系统对话框，遵循 Electron 安全边界）；注意：该操作为“签名文件管理”，不涉及“专辑导出”
- [x] T033 导入签名：调用 POST `/signature/import`，成功后列表自动刷新（SSE）
- [ ] T034 导出流程集成：在专辑导出步骤弹出签名选择，调用 `/export/sign-bridge` 将签名写入 `album_signatures` 并继续导出（frontend/src/pages/ 或组件）
- [x] T035 i18n：新增 `signature` 命名空间的中英文 key，覆盖所有用户可见文本（frontend/src/i18n/{zh-CN,en-US}/index.json）
- [x] T036 覆盖确认 UI：导入签名遇到“唯一标识相同”的重复时，弹出“覆盖/取消”确认对话框并完成 i18n 覆盖
- [ ] T037 导出前置校验：当专辑已有签名时，导出流程必须要求再选择签名（否则禁用继续/提示）；与后端桥接响应保持一致的错误提示

### 3.4.x 追加（一致性与可用性）

- [x] T037a 溢出滚动策略验证：当对话框内容超长时，验证 q-scroll-area/overflow-y 行为，无窗口外滚动或裁剪（NFR-001）
- [x] T037b 负例用例：签名管理对话框中不出现“签名专辑/导出专辑”入口（FR-017 负例断言）

## Phase 3.5: Tests & Integration（测试与集成）

- [ ] T040 Playwright 烟雾：应用可启动；“签名管理”按钮可见；对话框在 360x420 内完整显示
- [ ] T041 Playwright 烟雾：创建/导入后列表自动刷新（监听 SSE）
- [ ] T042 后端契约/单测：T010~T012 的 happy path 与错误路径通过
- [ ] T043 质量检查：UI 文案 i18n 覆盖；窗口尺寸与滚动策略符合 NFR；错误提示一致
- [ ] T043a i18n 验收 checklist：标题、按钮、表单标签、错误提示、占位文案、工具提示、确认对话
- [x] T043b UI 位置断言：确认“签名管理”按钮位于“删除”按钮右侧
- [ ] T044 Playwright 烟雾：导出流程中“必须选择签名”的校验触发与通过签名选择后可继续导出

### 3.5.x 追加（不可变性与结构校验）

- [ ] T045 不可变字段校验：后端拒绝更新 name/protectCode 字段；前端禁用编辑并在服务层校验（含负例测试）
- [ ] T046 签名历史结构测试：对 `album_signatures` 的 `signedAt[]` 进行合并/去重/排序测试（含边界：空、重复、多次导出）
- [ ] T047 i18n 覆盖检查：错误提示、确认对话、占位文案等新增 key 覆盖度扫描与校验（对照 `signature` 命名空间清单）
- [ ] T048 跨平台窗口可用面积复核：Windows/macOS 对话框在 360x420 内完整显示，无裁剪（NFR-001）

## Phase 3.6: Polish（完善）

- [x] T050 文档：更新 `quickstart.md` 与 `contracts/` 示例请求/响应
- [ ] T051 性能：签名导入/导出 <1s；导出签名桥 <1s；列表筛选 <100ms（可用本地基准或手测）
- [ ] T052 跨平台手测：Windows/macOS 基线验证；Linux 机会性验证

### 3.6.x 追加（治理与命名一致性）

- [ ] T053 术语与命名统一：统一为 signatureManagementDialog / signatureSelectDialog / exportSignBridge，并与 i18n `signature` 命名空间对齐（提供名词表）
- [ ] T054 PR 影响与回归保护说明：为涉及端点/关键流程的 PR 附带“影响范围与回归保护”说明，链接到 quickstart 或 contracts 的契约测试

---

## 依赖关系（精简）

- Setup（T001-T004）先于 Contracts（T010-T014）与 Backend（T020-T026）
- Contracts（T010-T014）先于 Backend 端点实现（T021-T023）与测试（T042）
- Backend（T020-T026）与 Frontend（T030-T037）完成后再做集成与烟测（T040-T044）

## 并行建议

- [P] T001/T002/T003 可并行（不同文件）
- [P] T020 与 T021-T023 分阶段并行（模型/编解码与处理器不同文件）
- [P] T030~T035 中 UI 与服务可与后端并行开发，以 mock/假数据推进

## 备注

- 旧的 `/sdk/signatures*` 类端点若存在，后续统一下线或迁移到“兼容层”，不在本阶段范围内复用。
- 保护码与唯一标识均由系统自动生成；UI 均不展示保护码。
- 列表“首次 GET + 后续 SSE 自动刷新”是第一阶段必须满足的交互契约。
