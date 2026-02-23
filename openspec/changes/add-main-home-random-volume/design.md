## Context

当前运行时音量链路已具备三层确定性叠加：全局、分离设备、按下/抬起。需求是在此基础上增加随机衰减，并支持“全局随机 + 按下/抬起随机单独控制”两层叠加：
- 任何随机结果都不超过实际音量；
- 用户无需改音效包即可对所有专辑生效；
- 默认体验安全，不出现“随机后几乎听不到”的明显劣化。

## Goals / Non-Goals

- Goals:
  - 提供主页面一键开启的全局随机音量体验。
  - 提供按下/抬起随机音量单独控制，并可与全局随机叠加。
  - 随机范围可控且有合理默认值。
  - 与现有音量体系（Base=1.6、负值降音）兼容。
- Non-Goals:
  - 不新增“智能策略模式”或多算法切换。
  - 不改动专辑编辑器的 `cut.volume` 语义。

## Decisions

- Decision 1: 使用“volume 值随机扣减”模型
  - 定义随机降幅 `d ∈ [0, maxReduceRatio]`。
  - 随机层体积值为 `deltaVolume = -d`。
  - 原因：该模型与当前 Base=1.6 的 volume 叠加机制直接一致，且天然只衰减不放大。

- Decision 2: 提供单一可配置项 `maxReduceRatio`
  - 默认 `3`；仅要求 `>= 0`，不设置上限。
  - 前端以“开关 + 输入框”呈现：
    - 全局随机：总开关开启后展示输入框；
    - 按下/抬起随机：总开关开启后展示展开项；每个节点独立开关开启后展示节点输入框。
  - 原因：
    - 3 在当前音量体系下可提供明显但可控的随机感作为默认体验；
    - 无上限允许高级用户按偏好获得更强随机降幅体验。

- Decision 3: 随机层位于现有三层之后，且支持双随机层叠加
  - 顺序：全局 → 分离设备 → 按下/抬起 → 全局随机 → 按下/抬起随机。
  - 原因：
    - 全局随机可覆盖所有触发场景；
    - 按下/抬起随机可进一步细分触发态并叠加个性化随机手感。

- Decision 4: 仅在真实播放链路生效，预览模式不生效
  - 与现有 `isPreviewMode` 逻辑对齐，避免编辑/预览被随机扰动。

- Decision 5: 文案命名显式区分全局与按下/抬起控制域
  - 全局入口命名为“随机音量(全局)”与“随机降幅上限(全局)”。
  - 单独控制入口命名为“随机音量(按下/抬起单独控制)”。
  - 单独控制说明文案明确“生效过程与全局随机音量设置叠加”。
  - 节点输入命名统一为“随机降幅上限(场景)”，例如“随机降幅上限(全局按下)”。
  - 原因：降低用户对“全局随机”和“单独控制随机”的混淆成本。

## Runtime Formula

- 随机采样降幅：
  - `d ~ Uniform(0, maxReduceRatio)`
- 随机层：
  - `deltaVolume = -d`

实现时通过追加随机层达成：
- 全局随机层：按 `main_home.random_volume_processing` 计算；
- 按下/抬起随机层：按 `main_home.press_release_random_volume_processing` 的事件态节点计算（仅节点启用时生效）。

## Risks / Trade-offs

- 风险：双随机层叠加会使波动更明显，可能接近静音。
  - 缓解：默认节点开关均关闭、默认值为 3，并通过文案明确“会与全局随机叠加”。
- 风险：随机策略会降低可复现性（同一按键每次响度不同）。
  - 缓解：由开关明确门控，默认关闭。

## Migration Plan

- 新增并默认写入：
  - `main_home.random_volume_processing.is_enabled = false`
  - `main_home.random_volume_processing.max_reduce_ratio = 3`
  - `main_home.press_release_random_volume_processing.is_enabled = false`
  - `main_home.press_release_random_volume_processing.{global|split}.{...}.is_enabled = false`
  - `main_home.press_release_random_volume_processing.{global|split}.{...}.max_reduce_ratio = 3`
- 对已有配置执行增量补齐，不覆盖历史已存在值。

## Open Questions

- 无。当前方案满足“简单可用 + 安全可控”，后续如需“智能模式”可另开变更。
