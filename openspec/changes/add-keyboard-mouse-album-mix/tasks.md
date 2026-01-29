# Tasks: add-keyboard-mouse-album-mix

## 1. Specs

- [x] 1.1 完成 `playback-album-routing` 的 delta spec（ADDED Requirements + 场景）
- [x] 1.2 `openspec validate add-keyboard-mouse-album-mix --strict --no-interactive`

## 2. Frontend（UI + Store）

- [x] 2.1 在设置存储中新增播放路由字段（统一/分离、键盘专辑、鼠标专辑）并持久化
- [x] 2.2 主页面增加“分离键盘/鼠标”开关与鼠标专辑选择器（默认隐藏）
- [x] 2.3 主页面的选择器绑定到“播放路由”状态（而非仅 `selected_key_tone_pkg`）
- [x] 2.4 键音专辑页面：将选择器语义收敛为“编辑专辑选择”，并展示提示文案（编辑模式与主页路由隔离）
- [x] 2.4a 编辑提示使用底部通知，并支持“不再提示”持久化
- [x] 2.5 进入键音专辑页面时：调用 `POST /keytone_pkg/set_playback_source_mode` 设置 `mode=editor`
- [x] 2.6 返回主页时：调用 `POST /keytone_pkg/set_playback_source_mode` 设置 `mode=route`，并调用 `POST /keytone_pkg/apply_playback_routing` 刷新快照
- [x] 2.7 在路由选择变更时，调用 `POST /keytone_pkg/apply_playback_routing` 提交路由（避免在播放热路径动态加载）
- [x] 2.8 设置页面新增“分离模式音量设置”展开栏（不提供分离开关）
- [x] 2.9 分离模式下新增键盘/鼠标独立音量设置（滑块、降幅、调试滑块开关）
- [x] 2.10 分离音量设置持久化与默认值迁移（默认 0 / 5 / false）
- [x] 2.11 设置页分离音量 UI 与主页一致（音量图标、滑块宽度、子项缩进与竖条高度）
- [x] 2.12 分离模式下键盘/鼠标独立静音开关与图标交互

## 3. SDK（Playback Routing）

- [x] 3.1 设计并实现一个“只读专辑配置加载器”（支持 plain / legacy-hex / core 三种形态）用于播放缓存
- [x] 3.2 新增播放路由状态：键盘配置、鼠标配置、及其对应 albumPath（用于解析 audioFiles 目录）
- [x] 3.3 修改 `keySound.KeySoundHandler`：根据 keycode 的输入类型选择配置源，并保留原 fallback 逻辑
- [x] 3.4 新增后端 API：`POST /keytone_pkg/apply_playback_routing`（一次性加载并缓存，失败时返回可诊断信息）
- [x] 3.5 新增后端 API：`POST /keytone_pkg/set_playback_source_mode`（`editor`/`route`），保证播放来源互斥
- [x] 3.6 当从 `editor` 返回 `route` 时，确保下一次 apply 能生成最新快照（必要时清理旧快照）
- [x] 3.7 新增 `mouse_fallback_to_keyboard` 配置项（默认 false），修改回退逻辑
- [x] 3.8 新增分离模式键盘/鼠标独立音量处理（叠加在全局音量之后）

## 4. Verification

- [ ] 4.1 手动验证：统一模式下行为与旧版本一致
- [ ] 4.2 手动验证：分离模式下键盘触发与鼠标触发分别来自不同专辑
- [ ] 4.3 手动验证：分离模式下鼠标专辑缺失时默认无声（不回退）
- [ ] 4.4 手动验证：开启回退选项后，鼠标专辑缺失回退到键盘专辑
- [ ] 4.5 低延迟验证：播放热路径不做磁盘 IO（只读缓存）
- [ ] 4.6 分离模式下：键盘/鼠标独立音量生效且叠加全局音量
- [ ] 4.7 设置页面：分离音量设置通过自定义展开栏展示（头部保留竖条，不显示图标；opacity + transform 过渡；背景仅 hover/点击时显示）
- [ ] 4.17 设置页面：内容区域右侧留更大间距，避免遮挡滚动条
- [ ] 4.8 设置页面：滑块行仅显示键盘/鼠标字样且不缩进；调试/降幅项保留竖条
- [ ] 4.11 设置页面：键盘/鼠标滑块优先展示，调试滑块紧随其后，降幅后置，调试开关最后
- [ ] 4.9 设置页面：键盘/鼠标独立静音功能可用（音量为 0% 时不允许取消静音）
- [ ] 4.10 设置页面：主滑块标签仅拖动时显示；调试滑块常显；并预留纵向空间避免遮挡
- [ ] 4.13 设置页面：标签无遮挡，键盘/鼠标字样与滑块水平对齐
- [ ] 4.14 设置页面：滑块标签不被父容器裁剪
- [ ] 4.15 设置页面：滑块标签不被展开内容背景覆盖或裁剪
- [ ] 4.16 设置页面：主滑块标签不被 overflow 裁剪（含 q-item__section；必要时用 !important 防止 scoped 覆盖）
