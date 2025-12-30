## Context

`Keytone_album.vue` 当前把"步骤 UI + 多层 dialog + SSE 映射 + 配置写回 + 依赖校验提示 + 样式"全部聚合在单文件中。

其核心风险在于：任何微小改动都会触发大范围重编译、难以定位回归，并且跨步骤复用 dialog 的诉求会进一步放大耦合。

## Goals / Non-Goals

### Goals

- **零行为变更**：拆分后用户侧行为保持一致。
- **边界清晰**：按 step 与对话框职责拆分，降低单文件复杂度。
- **可复用**：对话框支持被任意步骤直接拉起（不强绑在某个 step 内部）。
- **渐进迁移**：允许逐步拆分（先拆模板，再拆逻辑），每一步都可编译运行。

### Non-Goals

- 不做功能增强或 UI 重新设计
- 不引入新的状态管理框架
- 不改动后端接口与数据结构

---

## 架构设计详解

### 整体架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Keytone_album.vue (父组件)                            │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        所有状态 & 副作用                                │  │
│  │  - refs: step, pkgName, soundFileList, soundList, keySoundList...     │  │
│  │  - computed: downSoundList, upSoundList, keyOptions...                │  │
│  │  - watch: SSE 监听、配置自动保存                                        │  │
│  │  - methods: saveSoundConfig, deleteSound, previewSound...             │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                    │                                         │
│                    provide(KEYTONE_ALBUM_CONTEXT_KEY, ctx)                   │
│                                    │                                         │
│    ┌───────────────────────────────┼───────────────────────────────┐        │
│    │                               │                               │        │
│    ▼                               ▼                               ▼        │
│ ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐       │
│ │ StepLoadAudio    │    │ StepDefineSounds │    │ StepCraftKey     │  ...  │
│ │ Files.vue        │    │ .vue             │    │ Sounds.vue       │       │
│ │                  │    │                  │    │                  │       │
│ │ inject(ctx)      │    │ inject(ctx)      │    │ inject(ctx)      │       │
│ │ 只负责 UI 渲染    │    │ 只负责 UI 渲染    │    │ 只负责 UI 渲染    │       │
│ └────────┬─────────┘    └────────┬─────────┘    └────────┬─────────┘       │
│          │                       │                       │                  │
│          ▼                       ▼                       ▼                  │
│    ┌───────────┐           ┌───────────┐           ┌───────────┐           │
│    │ Dialog    │           │ Dialog    │           │ Dialog    │           │
│    │ 组件      │           │ 组件      │           │ 组件      │           │
│    │ inject()  │           │ inject()  │           │ inject()  │           │
│    └───────────┘           └───────────┘           └───────────┘           │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 为什么使用 provide/inject 而不是 props

| 方案                 | 优点                            | 缺点                                   |
| -------------------- | ------------------------------- | -------------------------------------- |
| **props 传递**       | 显式依赖，类型检查友好          | 参数爆炸（100+ 个状态/方法），维护困难 |
| **Pinia store**      | 跨组件共享自然                  | 状态迁移成本高，容易改变行为           |
| **provide/inject** ✓ | 避免 props 爆炸，保持单一数据源 | 隐式依赖，需要类型定义                 |

### 文件职责说明

```
frontend/src/components/
├── Keytone_album.vue              # 原始大组件（待改造为父组件）
│                                  # 职责：持有所有状态、副作用、provide Context
│
└── keytone-album/                 # 新的模块目录
    ├── index.ts                   # 模块入口
    │                              # 职责：统一导出，保持向后兼容
    │
    ├── types.ts                   # 类型定义
    │                              # 职责：定义 Context 接口、数据结构、注入 Key
    │
    ├── steps/                     # Step 子组件目录
    │   ├── StepLoadAudioFiles.vue     # Step 1: 加载音频源文件
    │   ├── StepDefineSounds.vue       # Step 2: 定义声音 [待创建]
    │   ├── StepCraftKeySounds.vue     # Step 3: 制作按键音 [待创建]
    │   └── StepLinkageEffects.vue     # Step 4: 联动声效 [待创建]
    │   # 职责：只负责 UI 渲染，通过 inject 获取状态
    │
    ├── dialogs/                   # Dialog 子组件目录
    │   ├── AddAudioFileDialog.vue     # 添加音频文件对话框
    │   ├── ManageAudioFilesDialog.vue # 管理音频文件对话框
    │   └── ... (更多对话框)
    │   # 职责：可复用的对话框 UI，可被任意 Step 调用
    │
    └── composables/               # 可复用逻辑 [Phase 4，可选]
        # 职责：抽离纯逻辑（如 SSE 映射、排序、校验）
```

### 数据流图

```
用户操作
    │
    ▼
┌───────────────────┐
│ Step/Dialog 组件   │  ←── inject(ctx) 获取状态和方法
│ (UI 层)           │
└─────────┬─────────┘
          │ 调用 ctx.xxx() 方法
          ▼
┌───────────────────┐
│ 父组件            │  ←── 状态更新
│ (数据层)          │
└─────────┬─────────┘
          │ watch 触发
          ▼
┌───────────────────┐
│ ConfigSet/Delete  │  ←── 写入配置文件
│ (持久化层)        │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ SDK 后端          │  ←── SSE 推送更新
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ 父组件            │  ←── messageAudioPackage 事件监听
│ 更新状态          │
└─────────┬─────────┘
          │ 响应式传播
          ▼
┌───────────────────┐
│ 子组件 UI 自动更新 │
└───────────────────┘
```

---

## Decision: 采用"父组件持有状态 + 子组件承载 UI"的拆分方式

### 选择

- 将现有 `Keytone_album.vue` 的 **状态（refs/computed/watch）与核心副作用（SSE 监听、ConfigSet/ConfigDelete、排序、依赖校验）** 保持在父层。
- 将每个 `q-step` 的模板与局部交互拆到单独的 step 组件。
- 将对话框从 step 内抽离为 `dialogs/` 下的组件，使其可被任意步骤直接拉起。

### 关键原因

- **最小风险**：保留数据流与副作用在一处，减少"拆分导致 reactivity/副作用丢失"的概率。
- **减少 prop 爆炸**：如果完全靠 props 传递大量 refs/方法，会导致 step 组件 API 复杂。

## Decision: 上下文注入（provide/inject）用于共享 ref 与 action

为避免 step 组件 props 过长，父组件提供一个 `KeytoneAlbumContext`：

- 状态：`step`、`pkgName`、各列表与选择项 refs、各种 dialog 的 `v-model` 状态等
- 行为：`saveSoundConfig`、`deleteSound`、`previewSound`、`deleteKeySound`、配置写回等

Step 与 Dialog 组件通过 `inject` 获取上下文，达到：

- 子组件更薄，只关心 UI 逻辑
- Dialog 可在任意 step 触发（满足"在某步骤中拉起其他步骤配置对话框"的诉求）

> 这仍然保持"单一数据源"在父组件，不引入新 store。

---

## 当前实现状态

### 已完成的文件

| 文件                                               | 说明                            | 状态       |
| -------------------------------------------------- | ------------------------------- | ---------- |
| `keytone-album/types.ts`                           | 定义 Context 接口和所有数据类型 | ✅ 已完成   |
| `keytone-album/index.ts`                           | 模块入口，统一导出              | ✅ 已完成   |
| `keytone-album/steps/StepLoadAudioFiles.vue`       | Step 1 UI 组件                  | ✅ 框架完成 |
| `keytone-album/dialogs/AddAudioFileDialog.vue`     | 添加音频文件对话框              | ✅ 框架完成 |
| `keytone-album/dialogs/ManageAudioFilesDialog.vue` | 管理音频文件对话框              | ✅ 框架完成 |

### 待完成

- [ ] 在父组件 `Keytone_album.vue` 中添加 `provide()`
- [ ] 用新的 Step 组件替换原有模板
- [ ] 创建 Step2/3/4 组件
- [ ] 创建更多 Dialog 组件

---

## Alternatives considered

1) **全 props 传递**
- 优点：显式依赖
- 缺点：参数爆炸、维护成本高

2) **抽到 Pinia store**
- 优点：跨组件共享更自然
- 缺点：状态迁移更大、更容易引入行为变化；不符合"最小改动"

3) **一次性把所有逻辑拆成 composables**
- 优点：最终结构更干净
- 缺点：迁移风险高；因此采用渐进式（先拆模板/对话框，再逐步抽逻辑）

## Risks / Trade-offs

- **风险：副作用触发时机变化**（watch、debounce、SSE listener）
  - 缓解：保留在父组件；拆分仅移动模板与事件绑定。
- **风险：Dialog 的 v-model 绑定与关闭时机**
  - 缓解：每个 dialog 单独抽取时保持 `v-model` 字段与事件处理不变。
- **风险：scoped style / :deep 选择器作用域变化**
  - 缓解：先保留样式在父组件；必要时将通用样式移动到更合适的位置，但需逐步验证。

## Migration Plan（渐进拆分）

1. ✅ 创建新目录结构与上下文类型（不改变功能）。
2. 逐个把 Step1~4 的模板迁移到 step 子组件：
   - 每迁移一个 step：编译/运行，手动点测关键路径。
3. 将可复用对话框从 step 中抽离到 dialogs（保持 v-model 与事件函数一致）。
4. 可选：把"纯函数/映射/校验"逐步迁移到 composables，父层只组装。

## Validation

- 静态检查：`npm -C frontend run lint`
- 构建验证：`npm -C frontend run build`
- 手动回归（最小集）：
  - Step1：上传/管理音频源文件
  - Step2：创建/编辑/删除/预览声音
  - Step3：创建/编辑/删除键音；依赖警告展示
  - Step4：切换 embedded test sound；打开全局设置对话框并选择声效
  - Step header 折叠/展开（step=99）与 Continue/Back 导航
