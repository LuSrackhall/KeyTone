# Design

## Context

项目当前已有两层音量调节能力：

1. 全局主页面音量（`main_home.audio_volume_processing.*`）
2. 分离模式键盘/鼠标音量（`main_home.split_audio_volume_processing.{keyboard|mouse}.*`）

但缺少第三层“事件态（down/up）”控制，导致用户无法单独调节按下/抬起触发音量。

本次变更需要在现有架构上平滑扩展，不破坏既有设置项与默认行为，并复用当前分离音量子模块交互模式。

## Goals / Non-Goals

### Goals

- 在设置页引入“按下/抬起音量单独控制”门控开关。
- 为全局、键盘、鼠标三个维度增加 down/up 独立音量参数（滑块、降幅、调试开关）。
- 在 SDK 播放热路径按 `keyState`（down/up）叠加对应音量层。
- 未开启开关时完全保持旧行为。

### Non-Goals

- 不重构键音专辑编辑页与其数据结构。
- 不改变已有“统一/分离路由”模式定义。
- 不引入新的设置页分组体系，仅在现有“主页面相关/分离音量”下扩展。

## Decisions

### Decision 1: 三层音量叠加顺序固定

- 顺序固定为：
  1) 全局音量（现有）
  2) 分离设备音量（现有，split 时）
  3) down/up 音量（新增，开关开启时）
- 理由：
  - 最大限度复用现有计算链路；
  - down/up 层仅做“最后一层微调”，最符合用户感知。
  - 在 SDK 中通过 `PlayKeySound(..., keyState, ...)` 显式传入事件态，避免隐式推断。

### Decision 2: down/up 能力由统一门控开关控制显示与生效

- 新增布尔开关（默认 false）。
- 关闭时：
  - UI 不展示新增 down/up 模块；
  - SDK 不应用 down/up 层（即便配置中存在数值）。
- 开启时：
  - 在“鼠标回退到键盘”设置项下展示独立展开项；
  - 展开项按“全局 → 键盘 → 鼠标”顺序组织 down/up 模块；
  - 展开内容采用三段式：滑块集中、降幅输入集中、调试开关集中；
  - 不显示“全局/键盘/鼠标”分组小标题，滑块左侧直接显示项目名。
- 理由：
  - 保持默认体验简洁；
  - 保证升级后行为不突变。

### Decision 3: 数据结构按“全局 + 分离设备”双域扩展

- 全局域：`main_home.press_release_audio_volume_processing.global.{down|up}.*`
- 分离域：`main_home.press_release_audio_volume_processing.split.{keyboard|mouse}.{down|up}.*`
- 每个节点包含：
  - `volume_normal`
  - `volume_normal_reduce_scope`
  - `is_open_volume_debug_slider`
  - `volume_silent`
- 理由：
  - 与现有 `main_home.audio_volume_processing` 与 `main_home.split_audio_volume_processing` 命名风格一致；
  - 便于后续扩展（例如滚轮/触控类输入）。

### Decision 4: 迁移采用“缺省补齐，不覆盖”

- `getConfigFileToUi` 中若新键缺失则写默认值。
- 若新键已存在则尊重已有值。
- 理由：
  - 避免破坏旧用户设置；
  - 保证版本升级后可预测。

## Alternatives Considered

### Alternative A: 不新增配置键，直接复用现有 split volume 字段

- 放弃原因：
  - 无法表达 down/up 语义；
  - 会导致一个值同时影响按下和抬起，不满足需求。

### Alternative B: 总是展示 down/up 模块，不做门控

- 放弃原因：
  - 主页面相关设置复杂度陡增；
  - 不符合“开启后才展示”的明确需求。

### Alternative C: 将 keyboard/mouse down/up 继续放在“分离模式音量设置”展开项内

- 放弃原因：
  - 与“按下/抬起音量单独控制”语义耦合不清；
  - 不符合“统一在按下/抬起展开项管理”的需求。

## Risks / Trade-offs

- 风险：设置项显著增多，用户理解成本提升。
  - 缓解：以统一门控隐藏高级项，默认保持简洁。
- 风险：多层叠加后调试复杂。
  - 缓解：沿用调试滑块与固定叠加顺序，并在 spec 中明确。
- 风险：键名过长导致维护成本上升。
  - 缓解：保持命名模式一致，集中在 `press_release_audio_volume_processing` 前缀下。

## Migration Plan

1. 增加默认常量与 `SetDefault` 键注册。
2. 前端初始化读取时进行缺失补齐。
3. SDK 播放链路读取新开关，按需叠加。
4. UI 在开关开启时展示模块并允许编辑。
5. 验证旧配置升级后行为不变。

## Open Questions

- 无（本次需求在展示条件、作用范围、子模块形态上已明确）。
