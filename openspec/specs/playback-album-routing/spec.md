# playback-album-routing Specification

## Purpose
TBD - created by archiving change add-main-home-random-volume. Update Purpose after archive.
## Requirements
### Requirement: 主页面全局随机音量开关

Normative: The system SHALL provide a global random-volume toggle in Main Home settings. When disabled, runtime MUST NOT apply random attenuation.

#### Scenario: 默认关闭且不影响现有行为

- **GIVEN** 用户未开启“随机音量”
- **WHEN** 任意键盘或鼠标声音被触发播放
- **THEN** 系统 MUST 保持现有确定性音量链路行为
- **AND** MUST NOT apply random attenuation

#### Scenario: 开启后对所有已载入专辑生效

- **GIVEN** 用户开启“随机音量”
- **WHEN** 任意已载入键音专辑中的声音被触发
- **THEN** 系统 MUST 在每次播放时应用一次随机音量衰减

#### Scenario: 参数项展示受开关门控

- **GIVEN** “随机音量”开关关闭
- **WHEN** 用户查看主页面相关设置末尾
- **THEN** 系统 MUST 隐藏“随机降幅上限”参数输入项
- **WHEN** 用户开启“随机音量”开关
- **THEN** 系统 MUST 展示“随机降幅上限”参数输入项

### Requirement: 随机音量采用随机扣减且允许无上限配置

Normative: Random volume attenuation SHALL be applied as a random subtraction in current volume-domain and MUST support non-negative unbounded custom max reduction.

#### Scenario: 每次播放仅做衰减不做放大

- **GIVEN** 用户设置 `max_reduce_ratio = R` 且 `R >= 0`
- **WHEN** 系统为本次播放采样随机降幅 `d` 且 `0 <= d <= R`
- **THEN** 系统 MUST 叠加随机层 `deltaVolume = -d`
- **AND** 随机层 MUST NOT increase effective volume

#### Scenario: 自定义值不设上限

- **GIVEN** 用户在设置页调整随机降幅上限
- **WHEN** 系统持久化并应用该值
- **THEN** 系统 MUST accept any finite value where `max_reduce_ratio >= 0`
- **AND** 系统 MUST NOT apply a configurable upper cap

#### Scenario: 默认值合理可用

- **GIVEN** 用户首次安装或历史配置缺失随机音量字段
- **WHEN** 系统初始化配置
- **THEN** `main_home.random_volume_processing.is_enabled` MUST default to `false`
- **AND** `main_home.random_volume_processing.max_reduce_ratio` MUST default to `3`

### Requirement: 按下/抬起随机音量单独控制

Normative: The system SHALL provide a dedicated press/release random-volume control that can be enabled independently and layered with global random volume.

#### Scenario: 单独控制总开关门控展开项

- **GIVEN** 用户未开启“按下/抬起随机音量单独控制”
- **WHEN** 用户查看主页面相关设置
- **THEN** 系统 MUST NOT 展示按下/抬起随机音量展开项
- **WHEN** 用户开启该总开关
- **THEN** 系统 MUST 展示独立展开项

#### Scenario: 展开项提供六个事件态节点

- **GIVEN** 用户开启“按下/抬起随机音量单独控制”并展开
- **WHEN** 系统渲染设置项
- **THEN** 系统 MUST 提供以下节点开关：
  - 全局按下、全局抬起
  - 键盘按下、键盘抬起
  - 鼠标按下、鼠标抬起
- **AND** 每个节点开关开启时，系统 MUST 展示该节点的随机降幅上限输入框

#### Scenario: 节点默认值

- **GIVEN** 用户首次安装或历史配置缺失按下/抬起随机音量字段
- **WHEN** 系统初始化配置
- **THEN** `main_home.press_release_random_volume_processing.is_enabled` MUST default to `false`
- **AND** 每个节点的 `is_enabled` MUST default to `false`
- **AND** 每个节点的 `max_reduce_ratio` MUST default to `3`

#### Scenario: 叠加两次随机降幅

- **GIVEN** 用户已开启全局随机音量
- **AND** 用户已开启“按下/抬起随机音量单独控制”且命中节点开关开启
- **WHEN** 发生一次非预览播放
- **THEN** 系统 MUST 先应用全局随机层
- **AND** MUST 再应用按下/抬起随机层

### Requirement: 随机音量层叠加顺序固定

Normative: When enabled, random-volume processing MUST be applied after existing deterministic volume layers.

#### Scenario: 运行时叠加顺序

- **GIVEN** 用户已开启随机音量
- **WHEN** 系统处理一次非预览播放
- **THEN** 叠加顺序 MUST 为：
  1. 全局音量处理层
  2. 分离设备音量层（若适用）
  3. 按下/抬起音量层（若适用）
  4. 全局随机音量层（若启用）
  5. 按下/抬起随机音量层（若启用且节点开关开启）

#### Scenario: 预览模式不受随机音量影响

- **GIVEN** 系统处于预览播放模式
- **WHEN** 用户触发预览播放
- **THEN** 系统 MUST NOT apply random-volume attenuation

### Requirement: 随机音量配置持久化键稳定

Normative: The frontend and SDK MUST persist random-volume settings using stable dedicated keys under `main_home.random_volume_processing`.

#### Scenario: 配置写入稳定键名

- **GIVEN** 用户修改随机音量设置
- **WHEN** 系统写入配置
- **THEN** MUST 使用以下键名：
  - `main_home.random_volume_processing.is_enabled`
  - `main_home.random_volume_processing.max_reduce_ratio`
  - `main_home.press_release_random_volume_processing.is_enabled`
  - `main_home.press_release_random_volume_processing.global.down.is_enabled`
  - `main_home.press_release_random_volume_processing.global.down.max_reduce_ratio`
  - `main_home.press_release_random_volume_processing.global.up.is_enabled`
  - `main_home.press_release_random_volume_processing.global.up.max_reduce_ratio`
  - `main_home.press_release_random_volume_processing.split.keyboard.down.is_enabled`
  - `main_home.press_release_random_volume_processing.split.keyboard.down.max_reduce_ratio`
  - `main_home.press_release_random_volume_processing.split.keyboard.up.is_enabled`
  - `main_home.press_release_random_volume_processing.split.keyboard.up.max_reduce_ratio`
  - `main_home.press_release_random_volume_processing.split.mouse.down.is_enabled`
  - `main_home.press_release_random_volume_processing.split.mouse.down.max_reduce_ratio`
  - `main_home.press_release_random_volume_processing.split.mouse.up.is_enabled`
  - `main_home.press_release_random_volume_processing.split.mouse.up.max_reduce_ratio`

### Requirement: 随机音量文案命名清晰区分控制域

Normative: The UI copy for random-volume controls SHALL clearly distinguish global scope and press/release separate scope.

#### Scenario: 全局与单独控制命名区分

- **GIVEN** 用户进入主页面相关设置
- **WHEN** 系统渲染随机音量相关入口与输入项
- **THEN** 全局入口 SHOULD clearly indicate global scope
- **AND** 按下/抬起单独控制入口 SHOULD clearly indicate separate scope
- **AND** 节点输入项 SHOULD use a consistent "max reduction (scenario)" naming format

#### Scenario: 单独控制说明文案强调生效过程叠加

- **GIVEN** 用户查看“随机音量(按下/抬起单独控制)”说明
- **WHEN** 系统渲染说明文案
- **THEN** 文案 SHOULD 明确表达：当全局随机音量开启时，单独控制与全局设置在生效过程中叠加

