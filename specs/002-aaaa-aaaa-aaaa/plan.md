# Implementation Plan: 键音专辑签名系统

**Branch**: `002-aaaa-aaaa-aaaa` | **Date**: 2025-10-01 | **Spec**: D:/safe/KeyTone/specs/002-aaaa-aaaa-aaaa/spec.md
**Input**: Feature specification from `/specs/002-aaaa-aaaa-aaaa/spec.md`

## Execution Flow (/plan command scope)

```text
1. Load feature spec from Input path
   → DONE
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → DONE（无 NEEDS CLARIFICATION 标记；关键不确定性在 spec 中已消除）
3. Fill the Constitution Check section based on the constitution document.
   → DONE（结合用户提供的 TDD 混合策略，标注差异处理）
4. Evaluate Constitution Check section below
   → DONE（记录偏离：严格 TDD 覆盖率阈值暂不作为门禁，转为“关键逻辑先测、UI 小量烟测”）
5. Execute Phase 0 → research.md
   → WILL DO（本命令完成输出 research.md）
6. Execute Phase 1 → contracts, data-model.md, quickstart.md
   → WILL DO（本命令完成输出）
7. Re-evaluate Constitution Check section
   → WILL DO（确保与混合策略一致）
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

 
## Summary
从用户场景出发，实现“签名管理（创建/导入/导出）+ 专辑导出签名嵌入”的最小可用版本。技术路线采用“最小新端点 + 最大复用”：
- 前端：Vue 3 + Quasar（Electron 模式），列表首次通过 GET 加载，后续复用既有 SSE；UI 原型先行，关键逻辑补少量单测/契约测。
- 后端：Go（gin）本地 REST 服务，签名管理的读写持久化复用现有接口：
  - 全局：/store/get、/store/set（viper 配置）
  - 专辑：/keytone_pkg/get、/keytone_pkg/set、/keytone_pkg/delete
  - 推送：/stream SSE（事件名：message、messageAudioPackage；载荷包含 get_all_value）
  - 仅为 .ktsign 导入/导出与“导出时签名”建立最小桥接端点（见“Phase 1: Design & Contracts”）
- 测试策略：保持务实的“保护存量 + 新增最小契约测”姿态：
  - 对新增最小端点与关键数据转换写契约/单元测试
  - UI/端到端保留少量烟雾用例；不强制“测试先行”作为门禁
**Language/Version**: Go 1.21+, TypeScript 4.9 + Vue 3
**Primary Dependencies**: gin（Go）, Quasar/Vue Router/Pinia/Vue I18n、axios（前端）
**Storage**: 本地 JSON 配置 + 资源目录；专辑导出文件内嵌签名结构（第一阶段以结构化完整性为主，不引入复杂加密）

**Testing**:
- 后端：Go testing（契约/单元/少量集成）
- 前端：Playwright（仅烟雾/关键路径），Vitest（可测逻辑/组件）
**Target Platform**: 桌面（Windows 为主，兼顾 macOS/Linux）
**Project Type**: web（frontend + backend + electron 外壳）
**Scale/Scope**: 单机应用；单用户本地签名管理

 
### Constraints

- 复用优先：签名管理 CRUD 不新增后端端点，读写均通过 /store/* 与 /keytone_pkg/*；实时更新复用 /stream。
- 动态端口发现：Electron 主进程输出 KEYTONE_PORT；渲染进程通过 `getBackendPort()` 同步 axios baseURL。
- 自动化成本控制：端到端仅保留烟雾/关键路径；新增能力采用最小契约测试覆盖；覆盖率为指导性指标。
- **固定窗口尺寸（CRITICAL）**: Electron 主窗口固定为 390x500px (`resizable: false`)，页面可用区域：
  - Windows: 379x458.5px
  - macOS: 389.5x458.5px
  - 参考代码: `frontend/src-electron/electron-main.ts` (line 276-289)
- **UI 设计原则**: 所有对话框必须适配固定窗口尺寸
  - 对话框最大尺寸: 宽度 ≤ 360px，高度 ≤ 420px
  - 禁止使用 `maximized` 属性（会导致对话框溢出窗口）
  - 内容溢出时使用 `q-scroll-area` 或 CSS `overflow-y: auto`
  - 参考成功案例: `Main_page.vue`、`Setting_page.vue` 的固定尺寸布局
- **国际化（i18n）约束**:
  - **必须适配语言**: 中文 (zh-CN) 和英文 (en-US)
  - **i18n 路径**: `frontend/src/i18n/{zh-CN,en-US}/index.json`
  - **命名空间**: 签名系统翻译使用 `signature` 命名空间（与 `keyToneAlbumPage` 同级）
  - **覆盖范围**: 所有对话框、按钮、表单标签、提示消息、错误文本
  - **暂不适配**: 其他语言（de-DE, es-ES, fr-FR 等）保持功能优先
  - **实施时机**: 在 UI 组件实现完成后统一添加 i18n keys


 
 
## Constitution Check

（结合项目宪章与用户策略的务实落地）


**I. 代码质量与架构分离**
- [x] 签名持久化复用 /store 与 /keytone_pkg；SSE 复用 /stream
- [x] 仅新增最小桥接端点（.ktsign 导入/导出、导出签名桥）并提供契约
- [x] UI 不强制先测：先实现→烟雾/关键路径 E2E→补组件/单测

- [x] 测试策略（与宪章目标一致的务实落地）：对“既有功能”优先用测试保护与回归；对“新增最小端点/逻辑”提供契约与测试，允许实现与测试交错推进，以“测试通过”作为门禁而非书写时序
- [x] 覆盖率“指导性非刚性门禁”：对新增范围维持合理覆盖（后端新增能力≈60% 或以上；前端关键路径具备用例）
- [x] 继续复用 i18n 基础与错误提示规范

- [x] 构建与发布沿用现有 Quasar/Electron 与 Go 流程


### Documentation (this feature)


```text
├── quickstart.md        # Phase 1 输出（本命令生成/更新）
├── contracts/           # Phase 1 输出（本命令生成/更新）
### Source Code (repository root)


```text

backend: sdk/ (Go, gin)
frontend/: Vue 3 + Quasar + Electron
  ├── src/components/
  ├── src/pages/
  ├── src/services/ (axios 封装)
  └── tests/ (E2E 少量烟测 + 组件/逻辑单测)
```

**Structure Decision**: 采用 Web 应用结构（frontend + backend + electron 壳），契约仅覆盖新增最小端点；其余复用接口无需重复契约。


 
## Phase 0: Outline & Research

1) Unknowns/Decisions 归档到 research.md：

- .ktsign 文件结构、打包/校验的最小需求（第一阶段）
- 新增最小端点的命名、字段与错误码
- Windows 文件对话/路径与 Electron 安全集成边界

1) Best Practices 收集：
- gin 的请求验证/错误响应约定
- 前端 axios 拦截器与错误提示一致性
- Playwright 在 Electron 开发流中的“最小可接受”烟测清单

2) Consolidate：以“Decision/Rationale/Alternatives”结构落盘。

 
## Phase 1: Design & Contracts

1) `data-model.md`：定义实体字段与校验规则（DigitalSignature、SignatureFile、AlbumSignatureRecord）；明确 Stage 1 的“完整性”与“不可变字段”（签名名、保护码、唯一标识不可变），以及资源目录关系（名片图片）。
2) 存储布局（复用接口）：
  - 全局签名管理器存放于全局配置（/store/get|set），约定数据键：`signature_manager`（包含列表与资源引用）。
  - 专辑内签名记录存放于专辑配置（/keytone_pkg/get|set|delete），约定数据键：`album_signatures`（包含签名历史与时间戳）。
  - SSE 复用 `/stream`，事件 `message`（全局）与 `messageAudioPackage`（专辑）承载 `get_all_value` 更新。
3) `contracts/`（仅新增最小端点）：
  - POST `/signature/export`：根据输入签名标识，导出 .ktsign（返回字节流或保存路径）。
  - POST `/signature/import`：导入 .ktsign 并写入全局签名管理器（内部通过 /store/set）。
  - POST `/export/sign-bridge`：导出流程中的签名桥接（将选定签名与导出打包过程衔接；实现上可作为现有导出流程的薄封装）。
  - 其余读写全部通过现有 /store/* 与 /keytone_pkg/* 完成，不在 `contracts/` 重复定义。
4) 合同测试（最小集）：
  - 覆盖上述 2~3 个新增端点（schema/状态码/错误分支）。
  - 回归性断言：/store 与 /keytone_pkg 行为不变（关键键存在与可读写）。
5) `quickstart.md`：
  - 手动验证（UI 操作步骤）包含：首次加载列表→创建/导入→自动刷新→导出流程选择签名→生成包。
  - 自动验证（少量 Playwright 烟测 + 后端契约测命令）。
6) 结构细化（Stage 1 必备）：
  - `album_signatures` 中每个签名需维护导出时间戳数组（多次导出多条），在导出签名桥中合并与持久化。
  - 资源与路径：.ktsign 导入/导出需处理名片图片等资源的复制与引用修正（资源目录位于配置同级的专用文件夹）。
6) **存储层设计要点**（详见 data-model.md Storage Layout）：

- **签名唯一标识**：无独立 id 字段，name 作为唯一标识；protectCode 用于加密/哈希，不作为业务标识
- **全局配置**（明文存储）：
  - key: `encrypt(protectCode)` - 对称加密后的保护码
  - value: 明文 JSON，包含 name、intro、cardImagePath、createdAt（无 id）
  - 便于前端快速索引与渲染
- **专辑配置**（双重加密）：
  - key: `encrypt(sha256(decrypt(protectCode) + name))` - 先解密保护码，与 name 拼接后计算 SHA-256 哈希，最后加密
  - value: `encrypt(JSON_payload)` - 对称加密后的完整签名信息
  - JSON_payload 含 name、intro、cardImagePath、exportTimeList、authorizationBlock（仅原始作者包含此字段, 其它作者的签名不含此字段）
  - 双重加密防止泄露签名逻辑与授权结构
- **AuthorizationBlock**：
  - authCode: 签名授权码（默认为固定写死到代码变量的 sha256 码, 创建是存入配置中, 校验时从配置中读取并与写死的变量做对比, 相同即允许二次导出；不匹配 = 不允许二次导出——想进行二次导出需联系原作者获取授权码方可通过二次导出校验）
  - authorizedList: 资格码列表（存储通过授权的三方签名资格码）

附：.ktsign（二进制）说明（Stage 1 明确）
- .ktsign 的“原始 JSON 内容”与“全局配置文件中的单一签名项”一致，形如 `{ key: encrypt(protectCode), value: <签名明文JSON> }`（可选携带 assets）。
- 该原始 JSON 进行对称加密后写入二进制文件（扩展名 .ktsign）。
- 导入流程：二进制 → 解密 → 解析 JSON（得到 key/value 与可选 assets）→ 写入全局配置（/store/set），并按需落盘资源（名片图片等）。
- 前端/UI 不暴露 protectCode，导入/导出与“导出签名桥”均仅以签名标识（name）进行交互，保护码校验由后端内部完成。

治理与执行（与宪章一致）
- 每个 PR 应附“影响与回归保护说明”（可链接至 quickstart 操作清单或 contracts 的契约测试）。
- 后端新增端点应先补/增契约测试再实现，或与实现交错推进，以“测试通过”作为门禁。
 
 
## Phase 2: Task Planning Approach（描述，不执行）

- 从 contracts & data-model 生成任务分解：
  - 后端：仅实现最小端点（.ktsign 导入/导出、导出签名桥）+ 契约测试
  - 前端：签名管理对话框（首次 GET + SSE 刷新）、导出流程签名选择与桥接调用、axios 服务、少量烟测
- 加注“测试策略标签”：contract-min / ui-smoke / component-unit
