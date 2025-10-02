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
从用户场景出发，实现“签名管理（创建/导入/导出）+ 专辑导出签名嵌入”的最小可用版本。技术路线采用：
- 前端：Vue 3 + Quasar（Electron 模式可复用），UI 先原型实现，关键组件/逻辑单测或契约测。
- 后端：Go（gin）本地 REST 服务，提供 /sdk/signatures 与专辑导出相关端点，契约优先。
- 测试策略：放弃“全 UI 先行”的严格 TDD，采用混合式：
  - 合同/数据转换/解析 → 测试优先（契约/单元）
  - UI → 先手动实现与演示，补充 1-2 条烟雾用例与关键组件边界测试。
**Language/Version**: Go 1.21+, TypeScript 4.9 + Vue 3
**Primary Dependencies**: gin（Go）, Quasar/Vue Router/Pinia/Vue I18n、axios（前端）
**Storage**: 本地 JSON 配置 + 资源目录；专辑导出文件内嵌签名结构

**Testing**:
- 后端：Go testing（契约/单元/少量集成）
- 前端：Playwright（仅烟雾/关键路径），Vitest（可测逻辑/组件）
**Target Platform**: 桌面（Windows 为主，兼顾 macOS/Linux）
**Project Type**: web（frontend + backend + electron 外壳）
**Scale/Scope**: 单机应用；单用户本地签名管理

 
### Constraints

- REST-first：所有能力尽量通过 REST 与 SDK 交互，弱化 Electron 依赖，为未来 Wails 版本迁移降低成本。
- 动态端口发现：Electron 调试模式下由主进程拉起 SDK 并输出 KEYTONE_PORT（stdout）；渲染进程通过 `window.myWindowAPI.getBackendPort()` 获取并调用 `UpdateApi(port)` 对齐 axios baseURL，避免硬编码端口。
- 自动化成本控制：端到端仅保留烟雾/关键路径；契约与核心逻辑先测，覆盖率目标为指导性指标而非硬门禁。

 
 
## Constitution Check

（结合项目宪章与用户策略的务实落地）


**I. 代码质量与架构分离**
- [x] Go 后端与前端边界通过 REST 契约隔离
- [x] API 契约文档在 `contracts/` 下维护
- [x] 前端不直连系统 API，系统交互经后端中转
- [x] UI 不强制先测：先实现→烟雾/关键路径 E2E→补组件/单测

- [x] 覆盖率目标“指导性而非门禁”：后端>60% 针对新增、前端组件关键路径有测试
- [x] i18n 现有基础继续复用

- [x] 可视化友好的错误提示（Notify）
- [x] 资源使用受 Electron/Go 典型约束；后续再基线化


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

**Structure Decision**: 采用 Web 应用结构（frontend + backend + electron 壳），API 契约放在 specs/ 下，便于跨层对齐。


 
## Phase 0: Outline & Research

1) Unknowns/Decisions 归档到 research.md：

- .ktsign 文件打包/校验信息（第一阶段范围）
- 契约字段命名与错误码规范
- Windows 文件对话/路径与 Electron 安全集成边界

1) Best Practices 收集：
- gin 的请求验证/错误响应约定
- 前端 axios 拦截器与错误提示一致性
- Playwright 在 Electron 开发流中的“最小可接受”烟测清单

2) Consolidate：以“Decision/Rationale/Alternatives”结构落盘。

 
## Phase 1: Design & Contracts

1) `data-model.md`：实体字段与校验规则（签名、签名文件、专辑签名记录）+ **存储布局**（全局配置明文 vs. 专辑配置双重加密，对称加密算法 AES-256-GCM + Base64 编码，密钥管理第一阶段应用级固定密钥）。
2) `contracts/`：从 FR-006~FR-015 推导端点（/sdk/signatures, 导入/导出, 专辑签名相关）。
3) 合同测试（最小集）：为每个端点写 1 个契约断言（schema/状态码）。
4) `quickstart.md`：两条路径——
   - 手动验证（UI 操作步骤）
   - 自动验证（少量 Playwright 烟测 + 后端契约测命令）。
5) 更新 agent context（如模板指引的脚本存在且需要）。

**存储层设计要点**（详见 data-model.md Storage Layout）：

- **签名唯一标识**：无独立 id 字段，name 作为唯一标识；protectCode 用于加密/哈希，不作为业务标识
- **全局配置**（明文存储）：
  - key: `encrypt(protectCode)` - 对称加密后的保护码
  - value: 明文 JSON，包含 name、intro、cardImagePath、createdAt（无 id）
  - 便于前端快速索引与渲染
- **专辑配置**（双重加密）：
  - key: `encrypt(sha256(decrypt(protectCode) + name))` - 先解密保护码，与 name 拼接后计算 SHA-256 哈希，最后加密
  - value: `encrypt(JSON_payload)` - 对称加密后的完整签名信息
  - JSON_payload 含 name、intro、cardImagePath、signedAt、authorizationBlock（仅原始作者包含此字段, 其它作者的签名不含此字段）
  - 双重加密防止泄露签名逻辑与授权结构
- **AuthorizationBlock**（第二阶段强依赖）：
  - authCode: 签名授权码（默认为固定写死的 sha256 码, 从变量读取并与配置中对比, 相同即允许二次导出；不匹配 = 需授权码校验）
  - authorizedList: 资格码列表（存储通过授权的三方签名资格码）

 
 
## Phase 2: Task Planning Approach（描述，不执行）

- 从 contracts & data-model 生成任务分解：
  - 后端：端点实现 + 契约测试驱动
  - 前端：对话框 UI、导出流程、axios 服务、少量烟测
- 加注“测试策略标签”：contract-first / ui-smoke / component-unit
