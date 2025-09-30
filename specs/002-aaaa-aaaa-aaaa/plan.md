
# Implementation Plan: 键音专辑签名系统

**Branch**: `002-aaaa-aaaa-aaaa` | **Date**: 2025-09-30 | **Spec**: D:/safe/KeyTone/specs/002-aaaa-aaaa-aaaa/spec.md
**Input**: Feature specification from `/specs/002-aaaa-aaaa-aaaa/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary

为键音专辑提供“数字签名”能力及本地签名管理：
- 在“键音专辑”页面顶部（删除按钮右侧）提供“签名管理”入口，以对话框形式进行签名创建/导入/删除/导出
- 签名由“签名名称、个人介绍、名片图片(可选)”构成；保护码在创建时由系统自动生成且不可编辑
- 支持签名导出为 .ktsign 文件与从 .ktsign 导入（导入无需输入保护码；保护码仅用于防重复、防篡改与校验）
- 导出专辑时可选择签名并写入专辑签名信息；第一阶段记录“允许二次导出”偏好但不做强校验

技术约束（用户显式说明）：
- 仅签名导入/导出使用浏览器前端能力（文件选择/保存）；除此之外的文件交互全部由 Go SDK 负责
- 前端与 SDK 的通信遵循项目既有模式
- 端到端测试使用 Playwright（本地已有相关 MCP 服务器, 可通过`#microsoft/playwright-mcp`来使用此MCP）

## Technical Context

**Language/Version**: Go 1.21+ (SDK), TypeScript + Vue 3 + Quasar (frontend), Electron (桌面)
**Primary Dependencies**: Pinia, Vue Router, Vue I18n, Electron IPC（既有通信模式）、Go SDK（本仓库 sdk/）
**Storage**: 本地配置文件（全局设置文件）与同级资源目录（名片图片等）
**Testing**: Go 单元与集成测试；前端单元测试（可补充）；E2E 使用 Playwright
**Target Platform**: Windows/macOS/Linux（跨平台）
**Project Type**: Web (frontend + backend/SDK)
**Performance Goals**: 遵循宪章：键盘响应 <10ms、音频延迟 <20ms、启动 <3s、空闲内存 <100MB、CPU <5%
**Constraints**:
- 文件交互：除导入/导出使用浏览器文件选择/保存，其它文件读写统一走 Go SDK
- 前端与 SDK 通信沿用既有模式（Electron contextBridge/IPC 等）
- .ktsign 为签名文件扩展名；导入无需输入保护码
**Scale/Scope**: 面向桌面端单用户场景；签名管理规模通常为个位到两位数

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. 代码质量与架构分离**:

- [ ] Go 后端和前端架构边界明确定义
- [ ] API 契约清晰且文档化
- [ ] 无跨层直接系统调用

**II. 测试优先与覆盖率标准**:

- [ ] TDD 方法论已规划（测试先行）
- [ ] 目标覆盖率已设定（后端 85%+，前端 80%+）
- [ ] 关键功能集成测试已计划

**III. 用户体验一致性**:

- [ ] UI 设计系统一致性已考虑
- [ ] 国际化支持已规划
- [ ] 跨平台行为一致性已设计

**IV. 性能与响应性要求**:

- [ ] 性能指标已定义（10ms 键盘响应，20ms 音频延迟）
- [ ] 资源使用限制已设定
- [ ] 性能测试策略已规划

**V. 跨平台兼容性**:

- [ ] 目标平台支持已明确
- [ ] 平台特定功能抽象已设计
- [ ] 构建和发布策略已考虑

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->
```text
sdk/                         # Go SDK（文件交互、签名存取、校验）
├── keySound/                # 现有模块（示例）
├── server/                  # 若有 RPC/IPC 适配层
└── ...

frontend/                    # Quasar + Vue 3 + Electron
├── src/
│  ├── pages/
│  │  └── AlbumPage.vue      # 键音专辑页面（签名按钮位于删除按钮右侧）
│  ├── components/
│  │  └── SignatureDialog.vue# 签名管理对话框（创建/导入/导出/删除）
│  ├── stores/
│  │  └── signature.ts       # 前端签名列表状态（从 SDK 拉取/推送）
│  └── services/
│     └── sdk-bridge.ts      # 与 SDK 的通信封装（沿用既有模式）
└── tests/
   └── e2e/                  # Playwright E2E（导入/导出/签名选择流程）

specs/002-aaaa-aaaa-aaaa/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── contracts/
   ├── signature-api.md
   └── album-signing-api.md
```

**Structure Decision**: 采用 Web（frontend + SDK）结构。前端负责交互与浏览器文件选择/保存；SDK 负责其余文件交互与签名校验、持久化。Electron 负责桥接通信，沿用既有项目模式与目录布局。

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved（本规划基于用户显式 override，允许在无 Clarifications 节的前提下继续）

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - 计划执行：`.specify/scripts/powershell/update-agent-context.ps1 -AgentType copilot`
   - 由于当前默认 shell 为 bash，将在实现阶段通过 PowerShell 终端按原样执行
   - 仅追加本计划中的新增技术点（Playwright E2E、.ktsign 扩展、SDK 统一文件交互约束）

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each contract → contract test task [P]
- Each entity → model creation task [P] 
- Each user story → integration test task
- Implementation tasks to make tests pass

**Ordering Strategy**:
- TDD order: Tests before implementation 
- Dependency order: Models before services before UI
- Mark [P] for parallel execution (independent files)

**Estimated Output**: 25-30 numbered, ordered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation                  | Why Needed         | Simpler Alternative Rejected Because |
| -------------------------- | ------------------ | ------------------------------------ |
| [e.g., 4th project]        | [current need]     | [why 3 projects insufficient]        |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient]  |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [ ] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved（用户提供 override 与技术约束）
- [ ] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
