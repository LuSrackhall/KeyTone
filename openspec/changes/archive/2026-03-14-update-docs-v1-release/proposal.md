# Change: 更新官方文档以迎接 v1.0.0 正式版发布

## Why

当前文档站（VitePress）仍停留在 v0.x 时代的描述，大量自 v0.6.0 以来引入的核心新特性（波形裁剪、鼠标专辑支持、多语言、键音签名/加密体系等）在文档中均无体现，导致：

1. 新用户无法通过文档快速了解 KeyTone 当前的真实能力；
2. 中英文双版本之间部分页面内容不一致，产生信息差；
3. 待发布的 v1.0.0 版本缺少对应的更新日志页面（changelog）。

## What Changes

### 新增内容

- **v1.0.0 更新日志页面**（中 + 英）：`docs/changelog/v1.0.0.md` & `docs/zh/changelog/v1.0.0.md`，以可读性优先的方式汇总所有重大功能与修复。
- **侧边栏配置更新**：在 `en.mts` 与 `zh.mts` 的 `sidebarChangelog()` 中新增 v1.0.0 条目，并将版本号组件指向最新版。

### 更新内容

| 页面                                   | 更新要点                                                         |
| -------------------------------------- | ---------------------------------------------------------------- |
| **首页**（`index.md` / `zh/index.md`） | 更新 feature 列表，体现鼠标支持、多语言、签名保护等新亮点        |
| **快速入门**（`quick-start`）          | 补全音量调节层次、键鼠分离模式入口、多语言切换说明               |
| **安装指南**（`installation`）         | 更新系统要求措辞，说明 macOS 已获得完整支持                      |
| **裁剪定义声音**（Step 2）             | 新增"波形可视化裁剪"入口介绍，说明双向同步、前端试听、音量指针条 |
| **铸造至臻键音**（Step 3）             | 新增"按下/抬起独立音量"能力说明                                  |
| **按键联动声效**（Step 4）             | 新增"鼠标按键绑定"说明及键鼠分离模式说明                         |
| **键音专辑介绍**（Introduction）       | 新增签名/加密体系简介，帮助用户理解 `.ktalbum` 安全机制          |

### 约束

- 中英文双版本必须同步更新，内容完全对应，不允许出现只更新单语言的情况。
- 不涉及任何应用代码修改。
- 不改变现有文档的 URL 结构。
- 文档页面风格与现有页面保持一致。

## Impact

- Affected docs capabilities: `user-documentation`（新增）
- Affected files:
  - `docs/docs/.vitepress/config/en.mts` — changelog 侧边栏 + 版本号
  - `docs/docs/.vitepress/config/zh.mts` — changelog 侧边栏 + 版本号
  - `docs/docs/index.md` — 首页英文
  - `docs/docs/zh/index.md` — 首页中文
  - `docs/docs/guide/getting-started/installation/index.md` — 安装指南英文
  - `docs/docs/zh/guide/getting-started/installation/index.md` — 安装指南中文
  - `docs/docs/guide/getting-started/quick-start/index.md` — 快速入门英文
  - `docs/docs/zh/guide/getting-started/quick-start/index.md` — 快速入门中文
  - `docs/docs/guide/key-package/introduction/index.md` — 专辑介绍英文
  - `docs/docs/zh/guide/key-package/introduction/index.md` — 专辑介绍中文
  - `docs/docs/guide/key-package/裁剪定义声音/index.md` — Step 2 英文
  - `docs/docs/zh/guide/key-package/裁剪定义声音/index.md` — Step 2 中文
  - `docs/docs/guide/key-package/铸造至臻键音/index.md` — Step 3 英文
  - `docs/docs/zh/guide/key-package/铸造至臻键音/index.md` — Step 3 中文
  - `docs/docs/guide/key-package/按键联动声效/index.md` — Step 4 英文
  - `docs/docs/zh/guide/key-package/按键联动声效/index.md` — Step 4 中文
  - `docs/docs/changelog/v1.0.0.md` — 新建英文 changelog
  - `docs/docs/zh/changelog/v1.0.0.md` — 新建中文 changelog

## Out of Scope（本次不做）

- 不新增独立的"设置页面"指南（留作后续迭代）
- 不新增独立的"签名/授权"指南（留作后续迭代）
- 不改变文档主题或 VitePress 版本
- 不涉及截图或录屏素材的制作
