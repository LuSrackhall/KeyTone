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

## 3. SDK（Playback Routing）

- [x] 3.1 设计并实现一个“只读专辑配置加载器”（支持 plain / legacy-hex / core 三种形态）用于播放缓存
- [x] 3.2 新增播放路由状态：键盘配置、鼠标配置、及其对应 albumPath（用于解析 audioFiles 目录）
- [x] 3.3 修改 `keySound.KeySoundHandler`：根据 keycode 的输入类型选择配置源，并保留原 fallback 逻辑
- [x] 3.4 新增后端 API：`POST /keytone_pkg/apply_playback_routing`（一次性加载并缓存，失败时返回可诊断信息）
- [x] 3.5 新增后端 API：`POST /keytone_pkg/set_playback_source_mode`（`editor`/`route`），保证播放来源互斥
- [x] 3.6 当从 `editor` 返回 `route` 时，确保下一次 apply 能生成最新快照（必要时清理旧快照）
- [x] 3.7 新增 `mouse_fallback_to_keyboard` 配置项（默认 false），修改回退逻辑

## 4. Verification

- [ ] 4.1 手动验证：统一模式下行为与旧版本一致
- [ ] 4.2 手动验证：分离模式下键盘触发与鼠标触发分别来自不同专辑
- [ ] 4.3 手动验证：分离模式下鼠标专辑缺失时默认无声（不回退）
- [ ] 4.4 手动验证：开启回退选项后，鼠标专辑缺失回退到键盘专辑
- [ ] 4.5 低延迟验证：播放热路径不做磁盘 IO（只读缓存）
