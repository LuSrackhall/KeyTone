# Change: 为主页面与分离模式增加按下/抬起独立音量控制

## Why

当前主页面相关设置仅支持“全局音量 + 分离模式键盘/鼠标音量”，无法继续区分“按下（down）”与“抬起（up）”两类触发场景。对于希望压低抬起声、保留按下反馈，或为键盘/鼠标分别设置不同 down/up 手感的用户，现有能力不足。

同时，现有设置页已经形成“分离音量子模块（滑块、降幅、调试滑块开关）”的交互范式。本次变更需要在不破坏既有设置项功能的前提下，沿用该范式扩展到 down/up 维度，并由“按下/抬起音量单独控制”开关统一门控展示。

## What Changes

- 新增“按下/抬起音量单独控制”总开关（默认关闭）。
- 在“鼠标回退到键盘”设置项下，新增“按下/抬起音量单独控制”设置项与独立展开项。
- 展开项统一承载 down/up 子模块，不再放入“分离模式音量设置”展开项。
- 展开项按“分离模式音量设置”同款组织方式展示：
  - 所有 down/up 音量滑块集中展示；
  - 所有 down/up 音量降幅输入集中展示；
  - 所有 down/up 调试滑块开关集中展示。
- 展开项内容顺序为：全局（按下/抬起）→ 键盘（按下/抬起）→ 鼠标（按下/抬起）。
- 不展示“全局/键盘/鼠标”小标题，滑块左侧直接显示对应项目名。
- UI 展示门控规则：
  - 仅当“按下/抬起音量单独控制”开关开启时，独立展开项才展示。
  - 关闭时保持当前版本展示与行为（不影响已有全局/分离音量设置）。
- 播放链路叠加规则扩展为：
  - 全局音量处理 → 分离模式键盘/鼠标音量（若启用 split）→ down/up 音量层（若启用单独控制）。
- 持久化与默认值：
  - 新增 down/up 对应配置键名，默认值保持“0（不增不减）/降幅 5/调试开关 false”。
  - 对旧配置兼容迁移：缺失新字段时写入默认值，不覆盖用户已有无关配置。

## Impact

- Affected specs:
  - `playback-album-routing`（新增 down/up 音量控制 requirement）
- Affected code (planned):
  - `frontend/src/pages/SettingPageChildren/MainHome_setting.vue`
  - `frontend/src/stores/setting-store.ts`
  - `frontend/src/i18n/*/index.json`（至少 zh-CN 与 en-US）
  - `sdk/config/config.go`
  - `sdk/keySound/keySound.go`
  - `sdk/server/server.go`（适配 `PlayKeySound` 参数）
- Non-goals:
  - 不改动键音专辑编辑逻辑。
  - 不改动与主页面相关设置无关的既有功能与交互。
