---
name: release-version
description: 执行 KeyTone 版本发布工作流：生成 CHANGELOG、确定版本号、更新前端版本、更新多语言文档站点、暂存变更。
applyTo: "**"
---

# 版本发布 Skill

## 适用场景

当用户要求执行以下操作时调用本 Skill：

- "发布新版本" / "release" / "打版"
- "更新版本号"
- "生成 changelog"
- "准备发版"
- "bump version"

## 工作流

### 第 1 步：生成增量 CHANGELOG

在项目根目录执行：

```bash
npx conventional-changelog -p angular -i CHANGELOG.md -s
```

该命令会根据**约定式提交（Conventional Commits）**规范，从上次版本标签至今的提交中提取增量日志，追加到 `CHANGELOG.md`。

> **注意**：`CHANGELOG.md` 顶部标题中的版本号可能为空（`# `），这是 `conventional-changelog` 的默认行为，后续步骤会处理。

生成后，从 `CHANGELOG.md` 顶部读取**增量部分**——即从文件顶部到上一次已知版本标题之间的所有内容。

---

### 第 2 步：确定版本号

#### 2.1 用户显式提供了版本号

直接使用用户提供的版本号（如 `1.2.3`），跳过自动决策。

#### 2.2 用户未提供版本号 → 自动决策

分析第 1 步生成的增量 changelog 内容，按语义化版本（SemVer）规则判定：

| 增量中有无 Breaking Changes / `Reverts` | 增量中有无 `Features` 类型 | 推荐版本 bump | 示例      |
|----------------------------------------|---------------------------|---------------|-----------|
| 有                                     | —                         | major         | 1.0.0 → 2.0.0 |
| 无                                     | 有                        | minor         | 1.0.0 → 1.1.0 |
| 无                                     | 无（仅 Bug Fixes / 其他）  | patch         | 1.0.0 → 1.0.1 |

**判定来源**：读取 `CHANGELOG.md` 顶部的增量内容，查找以下关键词：

- **Breaking Changes**：提交信息或 changelog 中带有 `BREAKING CHANGE`、`!:`（如 `feat!: ...`）、`Reverts` 标题
- **Features**：changelog 中存在 `### Features` 或 `### 新功能` 标题
- **Bug Fixes**：changelog 中存在 `### Bug Fixes` 或 `### 问题修复` 标题

从 `frontend/package.json` 中读取**当前版本号**作为基准，计算新版本号。

**自动决策示例**：

- 当前版本 `1.0.1`，增量包含 Features → 新版本为 `1.1.0`
- 当前版本 `1.0.1`，增量仅 Bug Fixes → 新版本为 `1.0.2`
- 当前版本 `1.0.1`，增量包含 BREAKING CHANGE → 新版本为 `2.0.0`

> ⚠️ 如果增量内容为空或无法判断，询问用户是否仍要发版，以及期望的版本号。

---

### 第 3 步：更新前端版本（UI 侧）

进入 `frontend/` 目录执行 npm version：

```bash
cd frontend && npm version <version>
```

此命令会：
- 更新 `frontend/package.json` 的 `version` 字段
- 自动创建一条 git commit（如果不需要自动 commit，可在后续 `git add` 中统一管理——但 npm version 默认会 commit，可以追加 `--no-git-tag-version` 参数来阻止）

**建议使用 `--no-git-tag-version`**，让所有变更统一暂存：

```bash
cd /Users/srackhalllu/Desktop/资源管理器/safe/KeyTone/frontend && npm version <version> --no-git-tag-version
```

> **注意**：`npm version` 会校验版本号格式（必须是有效的 SemVer，如 `1.2.3`），且不能与已有 git tag 重复（若未使用 `--no-git-tag-version`）。

---

### 第 4 步：更新官网文档站点

文档站点位于 `docs/docs/` 目录，使用 VitePress 构建，支持多语言。

#### 4.1 更新 VitePress 配置文件

对每个语言的配置文件（`docs/docs/.vitepress/config/*.mts`）执行以下更新：

1. **导航栏版本链接**：找到 `text: "v" + version` 或 `text: "v<旧版本>"` 的导航项，将其中的 `version` 引用或硬编码版本号更新为新版本。
   - 在中文配置中，对应的 items 链接为 `"zh/changelog/v" + version`
   - 在英文配置中，对应的 items 链接为 `"/changelog/v" + version`
   - ......其他语言同理，路径根据 base 配置调整。

2. **侧边栏（Sidebar）更新日志列表**：在每个配置文件的 `sidebarChangelog()` 函数中，在数组**最前面**插入新的版本条目：
   ```typescript
   { text: "v<新版本>", link: "/v<新版本>" },  // 英文
   { text: "v<新版本>", link: "/v<新版本>" },  // 中文（link 路径相同，因为 base 配置不同）
   { text: "v<新版本>", link: "/v<新版本>" },  // 其他语言同理
   ```

#### 4.2 创建各语言的 Changelog 版本文件

对于每种语言，在对应目录下创建 `<version>.md` 文件：

- 英文：`docs/docs/changelog/v<version>.md`
- 中文：`docs/docs/zh/changelog/v<version>.md`
- 其他语言(同理)：`docs/docs/<lang>/changelog/v<version>.md`

**文件内容格式模板**（基于项目中已有版本文件提炼）：

````markdown
# `<version>` *(<yyyy>-<mm>-<dd>)*

<英文/中文概述段落 — 从 CHANGELOG.md 增量中提取核心内容概括>

### Bug Fixes

* <从 CHANGELOG.md 中提取对应语言的 Bug Fixes 条目，保留原始格式>
* <每个条目包含 commit hash 链接和 issue 引用>

### Features

* <从 CHANGELOG.md 中提取对应语言的 Features 条目>
* <同上>
````

**转换规则**(以中英为例, 其它语言同理)：

- 从根目录 `CHANGELOG.md` 顶部的增量内容中提取信息
- 英文版（`docs/docs/changelog/`）：保留原始英文条目；若条目为中文，需翻译为英文
- 中文版（`docs/docs/zh/changelog/`）：保留原始中文条目；若条目为英文，需翻译为中文
- 版本标题格式：`` # `<version>` *(<date>)* ``
- 日期格式：`YYYY-MM-DD`（使用当前日期）
- 分类标题映射：
  - `### Features` / `### 新功能` → 英文版用 `### Features`，中文版用 `### 新功能`
  - `### Bug Fixes` / `### 问题修复` → 英文版用 `### Bug Fixes`，中文版用 `### 问题修复`
  - `### Reverts` / `### 回退` → 保留原标题
  - `### Performance Improvements` / `### 性能优化` → 保留原标题
- 每个条目保留原始格式（星号列表、commit 链接、issue 引用）

**如果某个分类在增量中不存在，则省略对应的标题**。

---

### 第 5 步：暂存所有变更

```bash
git add .
```

> 只执行 `git add`，不自动 commit，以便用户审查变更。

---

## 输出总结

完成上述步骤后，向用户展示以下格式的总结：

```
## ✅ 版本发布准备完成

### 版本信息
- **新版本号**：<version>
- **版本类型**：<major|minor|patch>
- **发布日期**：<YYYY-MM-DD>

### 已更新文件
1. 📝 `CHANGELOG.md` — 增量 changelog
2. 📦 `frontend/package.json` — UI 版本号
3. 🌐 `docs/docs/.vitepress/config/<lang>.mts` — 文档站配置（导航 + 侧边栏）
4. 📖 `docs/docs/changelog/v<version>.md` — 英文更新日志
5. 📖 `docs/docs/zh/changelog/v<version>.md` — 中文更新日志
6. 📖 `docs/docs/<lang>/changelog/v<version>.md` — 其他语言更新日志（如适用）

### 增量内容概要
<简要列出本次发版包含的主要变更类别和数量，例如：2 个新功能、5 个问题修复>

### ⏩ 后续建议
请审查以上变更，确认无误后可执行：

```bash
git commit -m "chore(release): <version>"
git tag <version>
git push origin <当前分支名>
git push origin <version>
```
