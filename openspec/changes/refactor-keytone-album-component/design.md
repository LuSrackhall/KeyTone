## Context

`Keytone_album.vue` 当前把“步骤 UI + 多层 dialog + SSE 映射 + 配置写回 + 依赖校验提示 + 样式”全部聚合在单文件中。

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

## Decision: 采用“父组件持有状态 + 子组件承载 UI”的拆分方式

### 选择

- 将现有 `Keytone_album.vue` 的 **状态（refs/computed/watch）与核心副作用（SSE 监听、ConfigSet/ConfigDelete、排序、依赖校验）** 保持在父层。
- 将每个 `q-step` 的模板与局部交互拆到单独的 step 组件。
- 将对话框从 step 内抽离为 `dialogs/` 下的组件，使其可被任意步骤直接拉起。

### 关键原因

- **最小风险**：保留数据流与副作用在一处，减少“拆分导致 reactivity/副作用丢失”的概率。
- **减少 prop 爆炸**：如果完全靠 props 传递大量 refs/方法，会导致 step 组件 API 复杂。

## Decision: 上下文注入（provide/inject）用于共享 ref 与 action（建议）

为避免 step 组件 props 过长，建议父组件提供一个 `KeytoneAlbumContext`：

- 状态：`step`、`pkgName`、各列表与选择项 refs、各种 dialog 的 `v-model` 状态等
- 行为：`saveSoundConfig`、`deleteSound`、`previewSound`、`deleteKeySound`、配置写回等

Step 与 Dialog 组件通过 `inject` 获取上下文，达到：

- 子组件更薄，只关心 UI 逻辑
- Dialog 可在任意 step 触发（满足“在某步骤中拉起其他步骤配置对话框”的诉求）

> 这仍然保持“单一数据源”在父组件，不引入新 store。

## Alternatives considered

1) **全 props 传递**
- 优点：显式依赖
- 缺点：参数爆炸、维护成本高

2) **抽到 Pinia store**
- 优点：跨组件共享更自然
- 缺点：状态迁移更大、更容易引入行为变化；不符合“最小改动”

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

1. 创建新目录结构与上下文类型（不改变功能）。
2. 逐个把 Step1~4 的模板迁移到 step 子组件：
   - 每迁移一个 step：编译/运行，手动点测关键路径。
3. 将可复用对话框从 step 中抽离到 dialogs（保持 v-model 与事件函数一致）。
4. 可选：把“纯函数/映射/校验”逐步迁移到 composables，父层只组装。

## Validation

- 静态检查：`npm -C frontend run lint`
- 构建验证：`npm -C frontend run build`
- 手动回归（最小集）：
  - Step1：上传/管理音频源文件
  - Step2：创建/编辑/删除/预览声音
  - Step3：创建/编辑/删除键音；依赖警告展示
  - Step4：切换 embedded test sound；打开全局设置对话框并选择声效
  - Step header 折叠/展开（step=99）与 Continue/Back 导航
