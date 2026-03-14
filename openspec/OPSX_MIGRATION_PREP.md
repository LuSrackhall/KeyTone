# OpenSpec 1.x 升级前置准备

本文档用于 KeyTone 从当前 OpenSpec 旧工作流迁移到 OpenSpec 1.x OPSX 工作流前的人工准备。

目标不是现在就执行升级，而是把升级前必须确认的事实、风险和顺序整理清楚，避免在真正运行 `openspec update` 或 `openspec init` 时踩坑。

## 1. 已核实的官方事实

基于 Fission-AI/OpenSpec 官方仓库、README、CLI 文档、OPSX 工作流文档、supported-tools 文档、migration guide、CHANGELOG 与相关实现代码，已经确认以下事实：

1. 当前官方包名是 `@fission-ai/openspec`，不是 `openspec-cli`。
2. 当前本机安装版本是 `0.20.0`，npm 最新版本是 `1.2.0`。
3. OpenSpec 1.0 引入了 OPSX 工作流，旧的 `/openspec:*` 命令被 `/opsx:*` 体系替代。
4. 迁移不会删除 `openspec/changes/`、`openspec/specs/`、归档变更，也不会删除用户自己写在 `AGENTS.md` 等文件中的非 OpenSpec 内容。
5. `openspec init` 和 `openspec update` 都会检测 legacy 文件，并在升级时执行清理流程。
6. `openspec/project.md` 不会被自动删除，官方要求人工迁移到 `openspec/config.yaml` 的 `context:` 字段。
7. `openspec/AGENTS.md` 属于会被清理的 legacy OpenSpec 产物。
8. GitHub Copilot 在 1.x 下的推荐产物路径是：
   - `.github/skills/openspec-*/SKILL.md`
   - `.github/prompts/opsx-<id>.prompt.md`
9. 1.2.0 引入了 profile 与 delivery 体系：
   - `core` 默认工作流：`propose`、`explore`、`apply`、`archive`
   - `custom` 可自定义 workflow 子集
   - `openspec update` 会按照当前 profile 和 delivery 同步项目文件
10. 1.2.0 会删除“未被当前 profile 选中”的 workflow 文件，因此升级前必须先想清楚要保留哪些 workflow。

## 2. 当前仓库状态快照

以下状态已在 2026-03-14 实际核实：

### 2.1 本机环境

- `openspec` 可执行文件来源：`@fission-ai/openspec`
- 当前版本：`0.20.0`
- Node 版本：`v24.11.0`
- 全局配置文件路径：`~/.config/openspec/config.json`
- 当前全局配置只包含 `featureFlags` 与 `telemetry`，尚未包含 `profile`、`delivery`、`workflows`

这意味着升级到 1.x 后，profile 迁移行为会直接影响本仓库生成哪些新文件。

### 2.2 仓库内已存在的 OpenSpec 相关产物

- `openspec/AGENTS.md`
- `openspec/project.md`
- `.github/prompts/openspec-proposal.prompt.md`
- `.github/prompts/openspec-apply.prompt.md`
- `.github/prompts/openspec-archive.prompt.md`

当前仓库不存在以下 1.x 目标产物：

- `openspec/config.yaml`
- `.github/skills/openspec-*/SKILL.md`
- `.github/prompts/opsx-*.prompt.md`

### 2.3 当前 OpenSpec 基线可读性

旧版 CLI 当前仍能正常读取本仓库：

- `openspec list` 可列出 active changes
- `openspec list --specs` 可列出主 specs

说明 `openspec/changes/` 与 `openspec/specs/` 主体结构仍是可迁移的。

## 3. 升级前必须注意的风险点

### 3.1 GitHub Copilot 指令文件会发生形态变化

当前仓库使用的是旧版 GitHub Copilot prompt：

- `.github/prompts/openspec-proposal.prompt.md`
- `.github/prompts/openspec-apply.prompt.md`
- `.github/prompts/openspec-archive.prompt.md`

升级后会转向 OPSX 命名和 skills 体系。最直接的变化包括：

1. prompt 文件名从 `openspec-*` 变为 `opsx-*`
2. 可能新增 `explore`，如果采用 `core` profile，则至少会有 4 个 commands
3. 会新增 `.github/skills/openspec-*/SKILL.md`
4. legacy prompts 会被识别为待清理对象

因此，在真正升级前，必须先保留一份这 3 个旧 prompt 的当前内容，以免历史上的项目化提示词被新模板覆盖后难以回溯。

### 3.2 `openspec/project.md` 需要人工抽取上下文

官方迁移文档明确要求把 `openspec/project.md` 中“真正需要长期注入给 AI 的内容”迁移到 `openspec/config.yaml` 的 `context:` 字段，而不是整份 Markdown 原封不动搬过去。

对本仓库来说，建议保留到 `context:` 的核心信息有：

1. 项目目标：KeyTone 是键盘音效模拟软件
2. 主要技术栈：Go SDK、Vue 3、Quasar、Electron、TypeScript、VitePress
3. 核心架构：Go 后端 + Electron/Vue 前端 + HTTP/SSE 通信
4. 关键业务概念：键音专辑、音频源文件、裁剪定义声音、高级声音、签名与导出流
5. 关键约束：GPL-3.0、不能分发受版权限制的音频、性能目标、主要支持 Windows
6. 项目习惯：规范文档尽量使用中文

不建议直接塞进 `context:` 的内容：

1. 过长的背景说明
2. 大段外链列表
3. 细粒度工具版本罗列
4. 将来容易频繁变化的冗余说明

### 3.3 严格校验当前并非全绿

升级前执行过：

```bash
openspec validate --all --strict --no-interactive
```

结果：21 项中 18 项通过，3 项失败。

失败项如下：

1. `add-album-config-encryption`
   - 文件：`openspec/changes/add-album-config-encryption/specs/album-config-encryption/spec.md`
   - 原因：使用了中文旧式标题，如“## 新增需求”，未使用官方 delta 头 `## ADDED Requirements`
2. `add-optional-contact-for-no-auth`
   - 文件：`openspec/changes/add-optional-contact-for-no-auth/specs/export-flow/spec.md`
   - 原因：文件是说明性增量文档，未使用官方 delta 头
3. `refactor-signature-info-dialog`
   - 文件：`openspec/changes/refactor-signature-info-dialog/specs/export-flow/spec.md`
   - 原因：文件是说明性增量文档，未使用官方 delta 头

这 3 个问题不会破坏 `changes/` 目录本身，但会让升级后的严格工作流继续报错。既然最终 spec 规范修订由你自己负责，那么这 3 个文件至少要被列为“升级后第一批手工规范化对象”。

### 3.4 当前全局配置没有 profile，升级行为需要人工选择

由于当前全局配置里没有：

- `profile`
- `delivery`
- `workflows`

所以升级到 1.2 后，必须主动决定两件事：

1. 是否接受官方默认 `core` profile
2. 是否保留 commands、skills，还是两者都保留

如果不先决定，直接执行更新，后续再切 profile 时，`openspec update` 可能会根据新 profile 删除未选中的 workflow 文件。

## 4. 建议的迁移前人工准备顺序

以下步骤建议在你真正执行 CLI 更新前完成。

### Step 1. 备份 legacy GitHub Copilot prompts

至少保留以下文件的内容快照：

- `.github/prompts/openspec-proposal.prompt.md`
- `.github/prompts/openspec-apply.prompt.md`
- `.github/prompts/openspec-archive.prompt.md`

建议备份目的：

1. 保留原有仓库自定义 guardrails
2. 对比新版 `opsx-*` 模板差异
3. 避免升级后忘记旧流程中的项目特定约束

### Step 2. 从 `openspec/project.md` 提炼 `config.yaml` 草稿内容

在真正运行升级前，先人工整理一份你准备放进 `openspec/config.yaml` 的内容草稿。

建议至少准备：

```yaml
schema: spec-driven

context: |
  KeyTone 是键盘音效模拟软件。
  项目采用 Go SDK + Vue 3 + Quasar + Electron 的混合架构。
  前端通过 HTTP 与 SSE 与 Go 后端通信。
  核心领域包括键音专辑、音频源文件、裁剪定义声音、高级声音、签名与导出流。
  规格文档尽量使用中文。
  必须遵守 GPL-3.0，不得引入受版权限制的音频资源。

rules:
  proposal:
    - 明确影响的 capability 与已有 change 冲突检查
    - 涉及导出、签名、加密、音频处理时写清回滚与兼容性影响
  specs:
    - 优先在现有 capability 上增量修改，避免重复 capability
    - 每条 requirement 必须包含至少一个 Scenario
    - 规格文档尽量使用中文，但保留 OpenSpec 要求的 delta 头格式
  design:
    - 涉及前后端协作、状态机、导出链路时写清数据流和失败处理
  tasks:
    - 任务拆分要可验证，包含校验与回归步骤
```

注意：这只是迁移前草稿，不要在现在就把它当作最终规范。

### Step 3. 先决定 profile / delivery 策略

建议先在脑中确定以下策略，再执行 1.2：

1. 如果你只想最小化迁移成本：先用 `core`
2. 如果你想保留更细的工作流控制：切到 `custom`，显式选择你要的 workflow
3. 对 GitHub Copilot，通常建议保留 `both`，这样同时拥有 `.github/skills/` 和 `.github/prompts/`

对本仓库的实际建议：

1. 首次迁移先采用 `core` 或 “接近 core 的 custom”
2. 完成一轮验证后，再决定是否启用 `verify`、`sync`、`onboard` 等扩展 workflow

原因很简单：当前仓库还带有若干 legacy spec 格式问题，先把基础工作流跑通，比一次性把所有 workflow 装满更稳。

### Step 4. 把 3 个严格校验失败项列入迁移后第一批修整

以下文件建议在升级后尽快按官方 delta 头规范修整：

1. `openspec/changes/add-album-config-encryption/specs/album-config-encryption/spec.md`
2. `openspec/changes/add-optional-contact-for-no-auth/specs/export-flow/spec.md`
3. `openspec/changes/refactor-signature-info-dialog/specs/export-flow/spec.md`

至少要完成：

1. 把旧式“新增需求/Requirement/Purpose”结构转为 OpenSpec 认可的 delta 头
2. 保留原有 requirement 语义，不要为了过校验丢信息
3. 继续维持 `#### Scenario:` 结构

### Step 5. 真正升级时优先用官方推荐顺序

官方文档确认：legacy 项目可以通过 `openspec update` 或 `openspec init` 完成迁移。

对本仓库，更推荐：

1. 先升级全局包
2. 再配置 profile
3. 最后在仓库内运行 `openspec update`

推荐顺序：

```bash
npm install -g @fission-ai/openspec@latest
openspec config profile
cd <repo>
openspec update
```

如果你想重新选择工具或强制重建，也可以考虑 `openspec init`，但对已经存在 `.github/prompts/` 的项目，`update` 更贴近“升级并刷新”的意图。

## 5. 本仓库执行升级时的重点观察项

当你真正运行升级命令时，建议重点观察输出里是否出现以下内容：

1. 检测到 `.github/prompts/openspec-*.prompt.md`
2. 检测到 `openspec/AGENTS.md`
3. 检测到 `openspec/project.md`，并给出“Needs your attention”提示
4. 是否生成 `.github/skills/openspec-*/SKILL.md`
5. 是否生成 `.github/prompts/opsx-*.prompt.md`
6. 是否提示重启 IDE 让新指令生效
7. 是否出现“Removed ... deselected workflows”之类信息

如果出现第 7 项，必须复查是不是 profile 选错，而不是直接继续。

## 6. 升级后立即做的验证

升级完成后，建议立刻做以下检查：

```bash
openspec --version
openspec config list
openspec list
openspec list --specs
openspec validate --all --strict --no-interactive
```

并人工确认：

1. GitHub Copilot 是否出现新的 `opsx-*` prompts
2. `.github/skills/` 是否存在
3. `openspec/config.yaml` 是否已生成
4. `openspec/project.md` 是否仍保留，等待人工迁移后再删除
5. 那 3 个已知失败项是否仍然是唯一失败源

## 7. 结论

KeyTone 当前并不是“不能迁移”，而是“迁移前需要先把上下文迁移策略、GitHub Copilot 旧 prompt 备份、以及 3 个 delta spec 历史遗留问题”明确下来。

只要按本文档顺序执行，真正的升级动作可以保持可控：

1. 旧 changes/specs 不会丢
2. legacy prompt 变更有据可查
3. `project.md` 的有效信息不会在清理流程里被忽略
4. 你后续做 spec 规范修订时，范围已经被锁定在 3 个明确文件上