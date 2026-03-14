# playback-album-routing Specification (Delta)

## ADDED Requirements

### Requirement: 按下/抬起音量单独控制开关

Normative: The system SHALL provide a dedicated toggle for press/release volume separation. When the toggle is off, press/release-specific UI and runtime adjustments MUST NOT be applied.

#### Scenario: 开关关闭时保持旧行为

- **GIVEN** 用户未开启“按下/抬起音量单独控制”
- **WHEN** 用户触发键盘或鼠标的按下/抬起事件
- **THEN** 系统 MUST only apply existing global and split device volume layers
- **AND** MUST NOT apply any press/release-specific volume layer
- **AND** 设置页 MUST NOT 展示新增的 down/up 子模块

#### Scenario: 开关开启时启用 down/up 能力

- **GIVEN** 用户开启“按下/抬起音量单独控制”
- **WHEN** 用户进入“设置 > 主页面相关”并查看“鼠标回退到键盘”设置项之后
- **THEN** 系统 MUST 展示“按下/抬起音量单独控制”设置项及其独立展开项
- **AND** 展开项 MUST 按“全局 → 键盘 → 鼠标”顺序展示

---

### Requirement: 全局按下/抬起独立音量

Normative: The system SHALL support separate global volume adjustments for key press (down) and key release (up), each with slider value, reduce-scope input, and debug-slider toggle.

#### Scenario: 全局按下与抬起音量分别生效

- **GIVEN** 用户已开启“按下/抬起音量单独控制”
- **AND** 用户设置了“全局按下音量”与“全局抬起音量”为不同值
- **WHEN** 触发 down 事件
- **THEN** 系统 MUST 叠加全局 down 音量参数
- **WHEN** 触发 up 事件
- **THEN** 系统 MUST 叠加全局 up 音量参数

#### Scenario: 全局 down/up 默认值与约束

- **GIVEN** 新增 down/up 全局配置尚未被用户修改
- **THEN** `volume_normal` 默认值 MUST be 0
- **AND** `volume_normal_reduce_scope` 默认值 MUST be 5
- **AND** `is_open_volume_debug_slider` 默认值 MUST be false
- **AND** `volume_silent` 默认值 MUST be false
- **AND** `volume_normal` MUST NOT exceed 0

---

### Requirement: 键盘/鼠标按下抬起独立音量

Normative: The system SHALL support press/release-specific volume adjustments for both keyboard and mouse sources in a unified press/release panel.

#### Scenario: 键盘 down/up 独立生效

- **GIVEN** 用户开启“按下/抬起音量单独控制”
- **AND** 用户设置 keyboard.down 与 keyboard.up 为不同值
- **WHEN** 触发键盘 down 事件
- **THEN** 系统 MUST 叠加 keyboard.down 音量层
- **WHEN** 触发键盘 up 事件
- **THEN** 系统 MUST 叠加 keyboard.up 音量层

#### Scenario: 鼠标 down/up 独立生效

- **GIVEN** 用户开启“按下/抬起音量单独控制”
- **AND** 用户设置 mouse.down 与 mouse.up 为不同值
- **WHEN** 触发鼠标 down 事件
- **THEN** 系统 MUST 叠加 mouse.down 音量层
- **WHEN** 触发鼠标 up 事件
- **THEN** 系统 MUST 叠加 mouse.up 音量层

#### Scenario: 统一展开项的 UI 子模块组成一致

- **GIVEN** 用户在设置页展开“按下/抬起音量单独控制”独立展开项
- **WHEN** 系统渲染全局、键盘、鼠标的 down/up 子模块
- **THEN** 界面 MUST 先集中展示所有 down/up 音量滑块
- **AND** 之后 MUST 集中展示所有 down/up 音量降幅输入
- **AND** 最后 MUST 集中展示所有 down/up 调试滑块开关
- **AND** MUST NOT 显示“全局/键盘/鼠标”分组小标题
- **AND** 滑块左侧 MUST 直接显示对应项目名

---

### Requirement: down/up 音量叠加顺序固定

Normative: The runtime MUST apply volume layers in a deterministic order to ensure predictable sound output.

#### Scenario: split 模式下的完整叠加链路

- **GIVEN** 用户启用 split 模式并开启“按下/抬起音量单独控制”
- **WHEN** 任一输入事件触发播放
- **THEN** 系统 MUST 按以下顺序叠加：
  1. 全局主页面音量层
  2. 分离设备音量层（keyboard 或 mouse）
  3. down/up 音量层（按事件态选择）

#### Scenario: unified 模式下的完整叠加链路

- **GIVEN** 用户处于 unified 模式并开启“按下/抬起音量单独控制”
- **WHEN** 任一输入事件触发播放
- **THEN** 系统 MUST 按以下顺序叠加：
  1. 全局主页面音量层
  2. down/up 音量层（全局域）

---

### Requirement: down/up 配置持久化键名

Normative: The frontend MUST persist press/release volume settings with stable dedicated keys.

#### Scenario: 持久化键名固定

- **GIVEN** 用户修改了按下/抬起音量相关设置
- **WHEN** 系统写入设置存储
- **THEN** MUST 使用稳定键名集合（至少包括）：
  - `main_home.press_release_audio_volume_processing.is_enabled`
  - `main_home.press_release_audio_volume_processing.global.down.volume_normal`
  - `main_home.press_release_audio_volume_processing.global.down.volume_normal_reduce_scope`
  - `main_home.press_release_audio_volume_processing.global.down.is_open_volume_debug_slider`
  - `main_home.press_release_audio_volume_processing.global.down.volume_silent`
  - `main_home.press_release_audio_volume_processing.global.up.volume_normal`
  - `main_home.press_release_audio_volume_processing.global.up.volume_normal_reduce_scope`
  - `main_home.press_release_audio_volume_processing.global.up.is_open_volume_debug_slider`
  - `main_home.press_release_audio_volume_processing.global.up.volume_silent`
  - `main_home.press_release_audio_volume_processing.split.keyboard.down.volume_normal`
  - `main_home.press_release_audio_volume_processing.split.keyboard.down.volume_normal_reduce_scope`
  - `main_home.press_release_audio_volume_processing.split.keyboard.down.is_open_volume_debug_slider`
  - `main_home.press_release_audio_volume_processing.split.keyboard.down.volume_silent`
  - `main_home.press_release_audio_volume_processing.split.keyboard.up.volume_normal`
  - `main_home.press_release_audio_volume_processing.split.keyboard.up.volume_normal_reduce_scope`
  - `main_home.press_release_audio_volume_processing.split.keyboard.up.is_open_volume_debug_slider`
  - `main_home.press_release_audio_volume_processing.split.keyboard.up.volume_silent`
  - `main_home.press_release_audio_volume_processing.split.mouse.down.volume_normal`
  - `main_home.press_release_audio_volume_processing.split.mouse.down.volume_normal_reduce_scope`
  - `main_home.press_release_audio_volume_processing.split.mouse.down.is_open_volume_debug_slider`
  - `main_home.press_release_audio_volume_processing.split.mouse.down.volume_silent`
  - `main_home.press_release_audio_volume_processing.split.mouse.up.volume_normal`
  - `main_home.press_release_audio_volume_processing.split.mouse.up.volume_normal_reduce_scope`
  - `main_home.press_release_audio_volume_processing.split.mouse.up.is_open_volume_debug_slider`
  - `main_home.press_release_audio_volume_processing.split.mouse.up.volume_silent`

#### Scenario: 旧配置升级兼容

- **GIVEN** 用户历史配置中不存在 `press_release_audio_volume_processing` 字段
- **WHEN** 应用初始化读取配置
- **THEN** 系统 MUST 以默认值补齐新增字段
- **AND** MUST NOT 修改与本能力无关的已有设置值
