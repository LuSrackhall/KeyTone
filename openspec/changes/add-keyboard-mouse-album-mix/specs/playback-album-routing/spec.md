# playback-album-routing Specification (Delta)

## ADDED Requirements

### Requirement: 播放路由支持键盘/鼠标混搭

Normative: The system SHALL allow the user to independently choose a KeyTone album for keyboard events and a KeyTone album for mouse button events, so that keyboard sounds and mouse sounds can originate from different albums.

#### Scenario: 统一模式下保持旧行为

- **GIVEN** 用户未启用“键盘/鼠标分离”模式
- **WHEN** 用户选择一个专辑作为当前播放专辑
- **THEN** 键盘事件与鼠标按键事件 SHALL 都使用该专辑的映射来播放声音

#### Scenario: 分离模式下键盘与鼠标来自不同专辑

- **GIVEN** 用户启用了“键盘/鼠标分离”模式
- **AND** 用户选择 A 专辑作为键盘播放专辑
- **AND** 用户选择 B 专辑作为鼠标播放专辑
- **WHEN** 用户触发任意键盘按键
- **THEN** 系统 SHALL 使用 A 专辑的配置解析并播放
- **WHEN** 用户触发任意鼠标按键
- **THEN** 系统 SHALL 使用 B 专辑的配置解析并播放

---

### Requirement: 事件分类规则必须稳定

Normative: The SDK MUST deterministically classify an input event as either “keyboard” or “mouse button” for routing purposes.

#### Scenario: 负数 keycode 归类为鼠标按键

- **GIVEN** KeySoundHandler 收到 `keycode` 字符串以 `-` 开头
- **WHEN** 进行路由选择
- **THEN** 该事件 MUST 被视为“鼠标按键事件”，并优先使用 Mouse 路由配置

#### Scenario: 非负 keycode 归类为键盘按键

- **GIVEN** KeySoundHandler 收到 `keycode` 字符串不以 `-` 开头
- **WHEN** 进行路由选择
- **THEN** 该事件 MUST 被视为“键盘按键事件”，并优先使用 Keyboard/Unified 路由配置

---

### Requirement: 播放热路径不做磁盘 IO

Normative: In split mode, the SDK playback path MUST NOT perform disk IO (including reading or decrypting album config files) during input event handling; it MUST rely on preloaded in-memory configuration instances.

#### Scenario: 应用路由后事件触发仍保持低延迟

- **GIVEN** 系统已成功应用（apply）键盘/鼠标播放路由
- **WHEN** 发生键盘或鼠标按键事件
- **THEN** 事件处理路径 MUST only read in-memory state to resolve the sound

---

### Requirement: 回退规则与错误容忍

Normative: When an album is missing, corrupted, or cannot be loaded for routing, the system MUST fall back to the unified album if available; if no album is available, it MUST fall back to the embedded test sound behavior.

#### Scenario: 鼠标专辑加载失败时回退（用户启用了回退策略）

- **GIVEN** 用户启用分离模式并设置了鼠标播放专辑
- **AND** 鼠标专辑因损坏/缺失导致加载失败
- **AND** 用户在设置中启用了 `mouse_fallback_to_keyboard`
- **WHEN** 用户触发鼠标按键事件
- **THEN** 系统 MUST fall back to the unified/keyboard album if present
- **AND** 若仍不可用则 MUST fall back to embedded test sounds

#### Scenario: 鼠标专辑缺失时不回退（默认行为）

- **GIVEN** 用户启用分离模式
- **AND** 鼠标专辑未选择或加载失败
- **AND** 用户未启用 `mouse_fallback_to_keyboard`（默认 false）
- **WHEN** 用户触发鼠标按键事件
- **THEN** 系统 MUST fall back to embedded test sounds directly
- **AND** MUST NOT use keyboard album

#### Scenario: 统一专辑加载失败时回退到内嵌测试音

- **GIVEN** 用户处于统一模式
- **AND** 统一专辑路径缺失或专辑无法加载
- **WHEN** 发生键盘或鼠标按键事件
- **THEN** 系统 MUST fall back to embedded test sounds

---

### Requirement: 前端选择与 SDK 路由提交解耦于编辑器

Normative: The system SHALL treat “Album Editor” as an exclusive mode that is decoupled from playback routing selection, so editing never implicitly changes the home playback routing, and home routing never implicitly changes the album being edited.

#### Scenario: 只改变播放路由不影响编辑器 SSE

- **GIVEN** 用户正在键音专辑编辑器中编辑某个专辑
- **WHEN** 用户仅调整播放路由（键盘/鼠标播放专辑选择）
- **THEN** 编辑器的配置 SSE 映射与写回闭环 MUST remain unchanged
- **AND** 编辑器当前选择的“编辑专辑” MUST NOT be changed implicitly

#### Scenario: 编辑专辑切换不改变主页路由

- **GIVEN** 用户在主页已配置播放路由（统一或分离）
- **WHEN** 用户进入键音专辑页并切换“编辑专辑”
- **THEN** 主页的播放路由设置 MUST NOT be modified

---

### Requirement: 提供明确的后端 Apply Routing API

Normative: The SDK server SHALL provide a stable HTTP API for the frontend to apply playback routing, and it MUST validate inputs and report per-source load status.

#### Scenario: 成功应用分离路由

- **GIVEN** 前端提交 `mode=split`、`keyboardAlbumPath=A`、`mouseAlbumPath=B`
- **WHEN** 调用 `POST /keytone_pkg/apply_playback_routing`
- **THEN** 服务端 MUST 预加载 A 与 B 的只读配置快照并写入 SDK 内存态
- **AND** 返回体 MUST 包含 keyboard/mouse 的最终解析路径与加载结果

#### Scenario: 入参为 UUID 时自动解析为路径

- **GIVEN** 前端提交 `keyboardAlbumPath` 值不包含 `/` 且不包含 `\\`
- **WHEN** 服务端处理该请求
- **THEN** 服务端 MUST 将其视为 UUID，并用 `AudioPackagePath/UUID` 组装真实路径后再加载

---

### Requirement: 播放来源模式互斥（route vs editor）

Normative: The system MUST ensure that only one playback source mode is active at a time: either routing-based playback (`route-unified`/`route-split`) or editor-based playback (`editor`).

#### Scenario: 进入键音专辑页自动切换到 editor 模式

- **GIVEN** 用户当前处于主页，且已应用任意路由播放模式
- **WHEN** 用户进入键音专辑页
- **THEN** 系统 MUST 切换到 `editor` 播放来源模式
- **AND** MUST suspend routing-based playback for the duration of the editor session

#### Scenario: 返回主页自动恢复 route 模式

- **GIVEN** 用户处于键音专辑页（`editor` 播放来源模式）
- **WHEN** 用户返回主页
- **THEN** 系统 MUST 恢复到 `route-unified` 或 `route-split`（以用户持久化的路由设置为准）
- **AND** MUST re-apply routing to ensure snapshots are up-to-date

---

### Requirement: 提供播放来源模式切换 API

Normative: The SDK server SHALL provide an HTTP API to switch the playback source mode, so the frontend can deterministically enter/exit editor mode.

#### Scenario: 前端进入编辑页时切换模式

- **GIVEN** 前端将进入键音专辑页
- **WHEN** 调用 `POST /keytone_pkg/set_playback_source_mode` 并提交 `mode=editor`
- **THEN** 服务端 MUST set the SDK playback source mode to `editor`
- **AND** KeySound playback MUST use the currently loaded editable album config (`audioPackageConfig.Viper`) for sound resolution

#### Scenario: 进入编辑页立即加载可编辑专辑

- **GIVEN** 用户进入键音专辑页且已存在上次编辑专辑
- **WHEN** 前端进入页面生命周期
- **THEN** 前端 MUST 立即调用 `LoadConfig(editedAlbumPath, false)` 以使编辑播放来源生效

#### Scenario: 前端回到主页时切换模式

- **GIVEN** 前端将返回主页
- **WHEN** 调用 `POST /keytone_pkg/set_playback_source_mode` 并提交 `mode=route`
- **THEN** 服务端 MUST set the SDK playback source mode back to routing
- **AND** 前端 MUST call `POST /keytone_pkg/apply_playback_routing` (unified/split) to populate snapshots

---

### Requirement: 路由状态持久化键名与迁移规则

Normative: The frontend MUST persist playback routing using dedicated keys, and migration from the legacy single-album selection MUST be deterministic.

#### Scenario: 从历史 `main_home.selected_key_tone_pkg` 迁移

- **GIVEN** 用户升级到支持播放路由的版本
- **AND** 新键名尚未写入
- **WHEN** 应用初始化读取设置
- **THEN** 系统 MUST 默认 `mode=unified`
- **AND** MUST 将 `unifiedAlbumPath` 初始化为历史 `main_home.selected_key_tone_pkg`

#### Scenario: 路由持久化键名固定

- **GIVEN** 系统需要持久化播放路由
- **WHEN** 写入设置存储
- **THEN** MUST 使用以下键名：
  - `playback.routing.mode`
  - `playback.routing.unified_album_path`
  - `playback.routing.keyboard_album_path`
  - `playback.routing.mouse_album_path`
  - `playback.routing.editor_notice_dismissed`
  - `playback.routing.mouse_fallback_to_keyboard`

---

### Requirement: 编辑专辑选择应独立持久化

Normative: The Album Editor MUST persist the last edited album independently from playback routing.

#### Scenario: 进入编辑页时保持上次编辑专辑

- **GIVEN** 用户上次在键音专辑页编辑过专辑 A
- **WHEN** 用户再次进入键音专辑页
- **THEN** 编辑页 SHOULD 默认选中并加载专辑 A（除非专辑不可用）
- **AND** 该行为 MUST NOT affect home playback routing settings

---

### Requirement: 编辑页提示使用可关闭通知

Normative: The album editor page MUST present its “editor mode” notice as a bottom notification that can be dismissed, and the “don’t show again” choice MUST be persisted.

#### Scenario: 用户关闭提示或选择不再提示

- **GIVEN** 用户进入键音专辑页
- **WHEN** 底部通知展示“编辑试听模式”提示
- **THEN** 用户 MAY 关闭通知
- **AND** 用户若选择“不再提示” MUST persist `playback.routing.editor_notice_dismissed=true`

---

### Requirement: 编辑器写回后自动刷新播放快照

Normative: In `editor` playback source mode, playback MUST follow the editable album config directly; routing snapshots are not used. When returning to routing mode, snapshots MUST be refreshed via apply.

#### Scenario: 编辑器修改了正在作为键盘播放源的专辑

- **GIVEN** 当前键盘播放专辑为 A
- **AND** 用户在编辑器中对专辑 A 执行配置写回
- **WHEN** 写回成功
- **THEN** 编辑器内的播放 MUST reflect the editable album config
- **AND** 当用户返回主页并恢复 route 模式时，系统 MUST re-apply routing to refresh snapshots
