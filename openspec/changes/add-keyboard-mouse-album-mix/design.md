# Design: add-keyboard-mouse-album-mix

## Context

当前系统的关键事实：
- 前端两处专辑选择框（主页、键音专辑页）都绑定同一个字段 `selectedKeyTonePkg`，其含义同时承担“用于播放的专辑”与“用于编辑的专辑”。
- SDK 播放侧依赖 `audioPackage/config` 的全局 `Viper` 实例；`keySound.KeySoundHandler` 每次事件触发时从全局配置读取 `key_tone.single/global` 决策。
- 鼠标按键通过负数 keycode 字符串（例如 `-1`）复用同一套 `key_tone.single.<keycode>.<down|up>` 映射。

因此“混搭”问题的根因不是 UI，而是播放侧配置源是单例：一次只能加载一个专辑。

同时，当前 UI 存在另一个结构性问题：

- “主页选择（日常使用/播放）”与“专辑页选择（创作/编辑）”复用了同一个字段，导致用户一旦进入编辑流程，就会不小心改动日常播放专辑。
- 播放路由引入后，如果仍然让两个页面同时影响播放来源，会进一步加剧心智负担与误操作。

## Goals / Non-Goals

### Goals

- 允许用户为键盘与鼠标分别选择播放专辑（支持 A 键盘 + B 鼠标）。
- 播放热路径不做磁盘 IO，不引入明显延迟抖动。
- 默认保持旧行为（统一专辑），升级后无需用户理解新概念也不影响使用。
- 不要求专辑格式升级，避免破坏社区已有专辑生态。

### Non-Goals

- 不提供“把 A 键盘 + B 鼠标 导出成新专辑”的一键合并。
- 不改变专辑编辑器内的 key_tone 配置结构。

## Decisions

### Decision 1: 引入 Playback Routing（键盘/鼠标专辑选择）

- **What**：前端新增统一/分离开关；分离时用户可分别选择键盘专辑与鼠标专辑。
- **Why**：用户心智模型简单（键盘 vs 鼠标），与现有“鼠标 keycode 为负数”实现一致。

### Decision 2: SDK 侧维护两份“只读配置实例”用于播放

- **What**：播放侧不再只依赖全局 `audioPackageConfig.Viper`；新增键盘/鼠标两份只读配置（可由 `viper.New()` + `ReadConfig(bytes)` 构建）。
- **Why**：避免在播放事件触发时动态加载/解密配置，从根本上保证延迟。

### Decision 3: 编辑配置与播放配置解耦

**What**：引入“播放来源模式（Playback Source Mode）”使主页路由播放与专辑编辑互斥：

- **主页（Main Home）**：只负责“日常播放”，其选择器仅绑定播放路由（只读快照）。
- **键音专辑页（Album Editor）**：只负责“创作/修改”，进入页面后临时切换到 `editor` 播放来源模式：
  - 播放时只使用当前“可编辑专辑”的 `LoadConfig`（可写配置 + WatchConfig + messageAudioPackage SSE）。
  - 暂停使用主页的路由快照（route-unified / route-split）。
  - 进入页面时 SHOULD 立即调用 `LoadConfig` 以让编辑试听生效。
- 离开键音专辑页返回主页时，系统 MUST 恢复 `route-*` 播放来源模式，并重新 apply 路由以确保播放快照与路由一致。

**Why**：

- 保证页面语义单一：主页=使用，专辑页=创作。
- 避免“编辑页切换专辑导致日常播放跳变”的误操作。
- 允许编辑页自然复用现有 SSE/写回闭环，并且可在编辑期间直接试听（试听来源与编辑对象一致）。

## Alternatives considered

### Alternative A: 运行时每次鼠标事件临时 LoadConfig

- **Rejected**：磁盘 IO + 解密 + viper 初始化会显著影响 <50ms 延迟目标。

### Alternative B: 把两个专辑合并成一个“临时叠加配置”

- **Rejected**：写回/加密/签名链路会变复杂，且用户很难理解“合并体”的来源与持久化。

## Risks / Trade-offs

- **取舍：编辑期间暂停日常路由播放**：在专辑编辑页中，播放来源切换为 `editor`，因此主页路由配置不会实时生效。
  - **Mitigation**：离开编辑页时强制恢复路由并重新 apply。

- **取舍：新增模式带来更多状态**：需要新增一个明确的“播放来源模式”并在前后端对齐。
  - **Mitigation**：将模式切换做成显式 API（进入/退出编辑页各调用一次），并在 spec 中约束互斥与回退。

## Migration Plan

- 新增设置字段时：
  - 默认 `mode=unified`
  - `keyboardAlbumPath` 与 `mouseAlbumPath` 均回退到历史的 `selectedKeyTonePkg`
- 持久化键名固定为：
  - `playback.routing.mode`
  - `playback.routing.unified_album_path`
  - `playback.routing.keyboard_album_path`
  - `playback.routing.mouse_album_path`
  - `playback.routing.editor_notice_dismissed`
  - `playback.routing.mouse_fallback_to_keyboard`（默认 false：彻底分离，鼠标无专辑则无声）
  - `main_home.split_audio_volume_processing.keyboard.volume_normal`
  - `main_home.split_audio_volume_processing.keyboard.volume_normal_reduce_scope`
  - `main_home.split_audio_volume_processing.keyboard.is_open_volume_debug_slider`
  - `main_home.split_audio_volume_processing.mouse.volume_normal`
  - `main_home.split_audio_volume_processing.mouse.volume_normal_reduce_scope`
  - `main_home.split_audio_volume_processing.mouse.is_open_volume_debug_slider`
- 键音专辑页“编辑专辑选择”沿用历史字段 `main_home.selected_key_tone_pkg`，用于记忆上次编辑专辑。
- UI：默认只显示一个选择器（兼容旧体验）；开启分离时显示第二个选择器。
- SDK：新增 API 供前端提交 routing；并新增“播放来源模式”API，在主页与编辑页之间切换。

### 回退策略设计（新增）

- 默认行为：分离模式下，鼠标专辑缺失时不回退到键盘专辑，鼠标事件使用内嵌测试音或静音。
- 用户可在设置界面的“主页面相关设置”中开启 `mouse_fallback_to_keyboard`，开启后分离模式下鼠标专辑缺失时会回退到键盘专辑。
- 该选项仅影响分离模式，统一模式下无意义（因为键盘/鼠标本就共用一个专辑）。

### 分离模式音量设计（新增）

- 在分离模式下，用户可为键盘与鼠标分别设置独立音量（默认 0，不增不减）。
- 独立音量在播放路径中叠加在“主页面全局音量”之后：
  - `finalVolume = globalVolume * splitDeviceVolume`
- 独立音量上限为 0，避免放大超过原始音量。
- 分离模式下允许对键盘/鼠标进行独立静音（与全局静音叠加）。
- 该配置仅在 `route-split` 播放来源模式下生效；`editor` 与 `route-unified` 不生效。

#### 设置页 UI 规范（补充）

- 分离音量设置使用展开栏展示，不在设置页提供分离开关。
- 展开栏头部使用自定义可点击行（替代默认 expansion item），保留左侧竖条与缩进一致性，但不显示图标（可在代码注释保留图标引用）。
- 展开/收起避免高度动画，使用轻量 opacity + transform 过渡，减少视觉延迟感。
- 展开栏头部背景默认透明，仅在 hover/点击时显示。
- 内容区域右侧需保留更大间距，避免影响滚动条展示。
- 键盘/鼠标独立音量滑块应与主页面音量控件保持一致（含静音图标与滑块宽度）。
- 滑块行仅显示“键盘/鼠标”字样，不额外缩进、也不展示竖条。
- 调试滑块紧随对应滑块下方；降幅输入在两组滑块之后；调试开关在降幅之后。
- 调试滑块与降幅设置保留竖条但不额外缩进。
- 主滑块的百分比标签仅在用户拖动时显示；调试滑块数字默认常显。
- 需要为滑块标签预留纵向空间，避免标签被遮挡或与上下内容重叠。
- 主滑块标签通过放开 overflow 解决裁剪问题（包括 q-item__section 等容器），避免过度增加垂直留白。
  - 由于 scoped 样式可能在不同组件间产生同等特异性覆盖，必要时可使用 `overflow: visible !important` 做可见性兜底。
- 需确保滑块容器允许溢出显示，避免标签被父层裁剪。
- 展开内容背景层级需低于滑块行内容，且展开容器需允许溢出，避免标签被背景覆盖/裁剪。
- 键盘/鼠标字样与滑块需水平对齐（不做刻意错位）。
- 展开内容区建议使用更柔和的背景/边框/轻微阴影，降低“突兀感”。
- 分离音量的静音图标需可点击切换，行为与主页面一致（音量为 0% 时不允许取消静音，并提示）。
- 设置页的滑块百分比标签需始终可见，避免被容器样式遮挡。

### UX Copy（专辑编辑页提示）

- 键音专辑页 SHOULD 通过页面底部通知展示提示文案（可关闭、支持“不再提示”），明确说明：
  - 该页面用于创作/修改键音专辑
  - 已为用户临时切换到“编辑试听模式”，与主页日常使用的播放路由隔离
  - 请创建新专辑或选择要编辑的专辑

示例文案（可在实现阶段调整为更贴合 UI 的短句，但语义 MUST 保持）：

> 键音专辑页面用于创作与编辑。
> 为避免影响主页日常使用，我们已临时切换到“编辑试听模式”。
> 请创建新的专辑，或选择要编辑的专辑。

## Open Questions

本提案为后续实现阶段消除不确定性，明确做出以下决策：

- 键音专辑页的选择器 SHALL 仅用于“编辑专辑选择”。
- 主页的选择器 SHALL 仅用于“播放路由选择”。
- 两者 MUST 互斥生效：进入键音专辑页后 MUST 切换到 `editor` 播放来源模式；返回主页后 MUST 切换回 `route-*`。

- 鼠标滚轮（Wheel）事件在本提案中明确 **Out of Scope**。
  - 本次路由仅覆盖鼠标按键（button）与键盘按键。
  - 未来若纳入 Wheel，将单独创建 change。

## SDK Specification Notes

本节作为实现指引，明确 SDK 侧的接口与线程模型（供实现时对齐）。

### Playback Routing State（SDK 内存态）

- SDK SHALL 维护一个只读的播放状态 `PlaybackState`，包含：
  - `SourceMode`: `route-unified` | `route-split` | `editor`
  - `Routing`：当 SourceMode 为 `route-*` 时生效
    - `UnifiedAlbumPath`（route-unified）
    - `KeyboardAlbumPath`、`MouseAlbumPath`（route-split）
    - 预加载的 `KeyboardSnapshot`、`MouseSnapshot`（只读快照；见下）
  - `EditorAlbumPath`：当 SourceMode 为 `editor` 时生效（用于诊断/展示；真实播放使用 `audioPackageConfig.Viper`）

- 该状态 MUST 可在并发下安全读取：
  - KeyEvent → KeySound 的热路径 MUST 只执行内存读取（RWMutex 或 atomic swap 指针均可）。
  - 更新路由 MUST 以“整块替换”的方式提交，避免读写同一结构体造成竞态。

### Album Config Snapshot（只读快照）

- SDK SHALL 提供“专辑配置快照”结构，用于播放侧读取配置：
  - Snapshot MUST 由 `package.json`（plain / legacy-hex / core）解密/解析得到的 JSON 构建。
  - Snapshot MUST 使用 `viper.New()` + `ReadConfig(io.Reader)` 方式加载 JSON 字节流。
  - Snapshot MUST NOT 启动任何 `WatchConfig()` 或 SSE。
  - Snapshot MUST 能提供 `Get(key string) any` 与 `AllSettings()`（必要时用于调试）。

### Apply Routing API（后端 HTTP）

- 后端 SHALL 提供一个新的 API 供前端提交播放路由：
  - `POST /keytone_pkg/apply_playback_routing`
  - 该 API MUST 在返回前完成：参数校验 → 两份专辑快照的加载（或回退）→ SDK 内存态替换。

- 为保证兼容性，该 API MUST 支持入参既可以是“专辑路径”，也可以是“仅 UUID”：
  - 若入参不包含路径分隔符（如无 `/` 且无 `\\`），则视为 UUID，SDK 使用 `AudioPackagePath/UUID` 组装为真实路径。

### Refresh Semantics（编辑器写回后的自动刷新）

本提案调整后，刷新语义明确如下：

- 当处于 `editor` SourceMode 时：
  - 编辑器写回与 SSE 维持现有行为；播放即时读取 `audioPackageConfig.Viper` 的最新值（无需刷新快照）。

- 当处于 `route-*` SourceMode 时：
  - 如果后台发生专辑配置变更且该专辑是当前路由源之一，SDK SHOULD 自动刷新对应快照以提升一致性。
  - 即使不自动刷新，前端在关键时机（返回主页/切换路由）MUST 调用 apply API 重新生成快照以确保最终一致。
