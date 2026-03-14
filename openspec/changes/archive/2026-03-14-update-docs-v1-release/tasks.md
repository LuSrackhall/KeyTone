# Tasks: 更新官方文档迎接 v1.0.0 发布

> **约束：** 中英文必须同步完成，每个任务对应的中英文页面一并处理，不允许只完成单语言。

---

## 阶段 A — 基础框架（可并行）

- [x] **A1. 更新版本侧边栏配置（英文）**
  - 文件：`docs/docs/.vitepress/config/en.mts`
  - 在 `sidebarChangelog()` 顶部插入 `{ text: "v1.0.0", link: "/v1.0.0" }`
  - 验证：本地 dev server 的 changelog 侧边栏中可见 v1.0.0 条目

- [x] **A2. 更新版本侧边栏配置（中文）**
  - 文件：`docs/docs/.vitepress/config/zh.mts`
  - 在 `sidebarChangelog()` 顶部插入 `{ text: "v1.0.0", link: "/v1.0.0" }`
  - 验证：中文站点 changelog 侧边栏中可见 v1.0.0 条目

---

## 阶段 B — changelog 页面（依赖 A 完成）

- [x] **B1. 新建 v1.0.0 更新日志（英文）**
  - 新建文件：`docs/docs/changelog/v1.0.0.md`
  - 内容要点：
    - 标题：`# v1.0.0 Release Notes`
    - 新增：鼠标按键音支持 & 键鼠分离播放模式
    - 新增：波形可视化裁剪（waveform trimmer）
    - 新增：按下/抬起独立音量控制
    - 新增：全局随机音量变化（主页面）
    - 新增：音效选择器搜索
    - 新增：专辑签名 & 授权体系
    - 新增：专辑配置加密
    - 新增：13+ 语言国际化（含 RTL 阿拉伯语）
    - 优化：macOS 完整支持（托盘行为、边框、自启动）
    - 多项稳定性修复（内存泄漏、并发、删除逻辑等）
  - 验证：页面可正常访问并渲染

- [x] **B2. 新建 v1.0.0 更新日志（中文）**
  - 新建文件：`docs/docs/zh/changelog/v1.0.0.md`
  - 内容与 B1 完全对应，使用中文
  - 验证：中文站点 changelog 页面可正常访问并渲染

---

## 阶段 C — 首页（可并行）

- [x] **C1. 更新英文首页特性列表**
  - 文件：`docs/docs/index.md`
  - 更新 `features` 数组，体现以下新亮点：
    1. 鼠标 + 键盘双重支持（Mouse & Keyboard）
    2. 多语言国际化（13+ 语言含 RTL）
    3. 波形可视化裁剪
    4. 专辑签名保护（防篡改 `.ktalbum`）
  - 保留已有条目的核心信息，避免大幅删减
  - 验证：页面 feature cards 渲染正常

- [x] **C2. 更新中文首页特性列表**
  - 文件：`docs/docs/zh/index.md`
  - 内容与 C1 完全对应，使用中文
  - 验证：中文首页 feature cards 渲染正常

---

## 阶段 D — 安装指南（可并行）

- [x] **D1. 更新安装指南（英文）**
  - 文件：`docs/docs/guide/getting-started/installation/index.md`
  - 修改"System Requirements"系统要求小节，注明 macOS 已获完整支持
  - 验证：页面内容准确，无死链

- [x] **D2. 更新安装指南（中文）**
  - 文件：`docs/docs/zh/guide/getting-started/installation/index.md`
  - 内容与 D1 完全对应，使用中文
  - 验证：页面内容准确，无死链

---

## 阶段 E — 快速入门（可并行）

- [x] **E1. 更新快速入门页面（英文）**
  - 文件：`docs/docs/guide/getting-started/quick-start/index.md`
  - 内容补充：
    - 音量调节：提及全局音量 / 主页面随机音量变化选项
    - 语言设置：提及 16 语言可选
    - 键鼠分离模式入口说明（分离后可独立选择鼠标专辑）
  - 验证：内容逻辑连贯

- [x] **E2. 更新快速入门页面（中文）**
  - 文件：`docs/docs/zh/guide/getting-started/quick-start/index.md`
  - 内容与 E1 完全对应，使用中文
  - 验证：内容逻辑连贯

---

## 阶段 F — 键音专辑介绍（可并行）

- [x] **F1. 更新键音专辑介绍页面（英文）**
  - 文件：`docs/docs/guide/key-package/introduction/index.md`
  - 在"安全机制"部分或 `.ktalbum` 说明处补充：签名保护体系（防篡改）、配置加密说明
  - 验证：内容与中文版一致

- [x] **F2. 更新键音专辑介绍页面（中文）**
  - 文件：`docs/docs/zh/guide/key-package/introduction/index.md`
  - 内容与 F1 完全对应，使用中文
  - 验证：内容与英文版一致

---

## 阶段 G — Step 2 裁剪定义声音（可并行）

- [x] **G1. 更新"裁剪定义声音"指南（英文）**
  - 文件：`docs/docs/guide/key-package/裁剪定义声音/index.md`
  - 新增章节或扩充现有内容：
    - 波形可视化裁剪（Waveform Trimmer）：拖拽选区、与时间输入双向同步
    - 前端试听播放条（无需 SDK）
    - 音量指针条（dB 可视调节，范围 ±18 dB）
    - 注明：数字输入裁剪能力保留，波形组件为增强选项
  - 验证：新增章节标题层级符合当前文档 outline 配置（h2/h3/h4）

- [x] **G2. 更新"裁剪定义声音"指南（中文）**
  - 文件：`docs/docs/zh/guide/key-package/裁剪定义声音/index.md`
  - 内容与 G1 完全对应，使用中文
  - 验证：与英文版章节结构一致

---

## 阶段 H — Step 3 铸造至臻键音（可并行）

- [x] **H1. 更新"铸造至臻键音"指南（英文）**
  - 文件：`docs/docs/guide/key-package/铸造至臻键音/index.md`
  - 新增内容：
    - 按下/抬起独立音量控制（Press/Release Volume）：可为按下和抬起声音分别设置独立音量（dB）
    - 说明独立音量与全局音量的叠加关系
  - 验证：新增段落不与现有"Diverse Playback Modes"等章节冲突

- [x] **H2. 更新"铸造至臻键音"指南（中文）**
  - 文件：`docs/docs/zh/guide/key-package/铸造至臻键音/index.md`
  - 内容与 H1 完全对应，使用中文
  - 验证：与英文版语义完全一致

---

## 阶段 I — Step 4 按键联动声效（可并行）

- [x] **I1. 更新"按键联动声效"指南（英文）**
  - 文件：`docs/docs/guide/key-package/按键联动声效/index.md`
  - 新增内容：
    - 鼠标按键支持：可将音效绑定到鼠标左键/右键/中键/滚轮的按下/抬起事件
    - 键鼠分离模式：支持分别指定键盘专辑和鼠标专辑（主页面入口介绍）
    - 说明：增加了对鼠标按键音效的独立配置（与键盘键位配置并列，同属键音专辑内容）
  - 验证：新增内容不破坏现有"Flexible Binding Strategies"章节的逻辑

- [x] **I2. 更新"按键联动声效"指南（中文）**
  - 文件：`docs/docs/zh/guide/key-package/按键联动声效/index.md`
  - 内容与 I1 完全对应，使用中文
  - 验证：与英文版语义完全一致

---

## 阶段 J — 验证（依赖所有前置任务完成）

- [x] **J1. 中英文内容一致性人工比对**
  - 逐页核对中英文版本，确认：章节数量一致、核心信息无遗漏、无仅出现在单语言版本的段落
  - 验证清单：A1-A2、B1-B2、C1-C2、D1-D2、E1-E2、F1-F2、G1-G2、H1-H2、I1-I2 ✅

- [ ] **J2. 本地构建验证（可选）**
  - 在 `docs/` 目录运行 `npm run docs:build`（或 `vitepress build`）检查无构建报错
  - 验证：changelog 侧边栏 v1.0.0 条目可点击并正常跳转
  - 验证：所有内部链接无 404

---

## 依赖关系

```
A（config）→ B（changelog）
C / D / E / F / G / H / I（各指南页面，可并行）
所有 → J（验证）
```

> A1/A2、B1/B2 应先行，其余阶段 C~I 可并行推进。


> **约束：** 中英文必须同步完成，每个任务对应的中英文页面一并处理，不允许只完成单语言。

---

## 阶段 A — 基础框架（可并行）

- [ ] **A1. 更新版本侧边栏配置（英文）**
  - 文件：`docs/docs/.vitepress/config/en.mts`
  - 在 `sidebarChangelog()` 顶部插入 `{ text: "v1.0.0", link: "/v1.0.0" }`
  - 验证：本地 dev server 的 changelog 侧边栏中可见 v1.0.0 条目

- [ ] **A2. 更新版本侧边栏配置（中文）**
  - 文件：`docs/docs/.vitepress/config/zh.mts`
  - 在 `sidebarChangelog()` 顶部插入 `{ text: "v1.0.0", link: "/v1.0.0" }`
  - 验证：中文站点 changelog 侧边栏中可见 v1.0.0 条目

---

## 阶段 B — changelog 页面（依赖 A 完成）

- [ ] **B1. 新建 v1.0.0 更新日志（英文）**
  - 新建文件：`docs/docs/changelog/v1.0.0.md`
  - 内容要点：
    - 标题：`# v1.0.0 Release Notes`
    - 新增：鼠标按键音支持 & 键鼠分离播放模式
    - 新增：波形可视化裁剪（waveform trimmer）
    - 新增：按下/抬起独立音量控制
    - 新增：全局随机音量变化（主页面）
    - 新增：音效选择器搜索
    - 新增：专辑签名 & 授权体系
    - 新增：专辑配置加密
    - 新增：13+ 语言国际化（含 RTL 阿拉伯语）
    - 优化：macOS 完整支持（托盘行为、边框、自启动）
    - 多项稳定性修复（内存泄漏、并发、删除逻辑等）
  - 验证：页面可正常访问并渲染

- [ ] **B2. 新建 v1.0.0 更新日志（中文）**
  - 新建文件：`docs/docs/zh/changelog/v1.0.0.md`
  - 内容与 B1 完全对应，使用中文
  - 验证：中文站点 changelog 页面可正常访问并渲染

---

## 阶段 C — 首页（可并行）

- [ ] **C1. 更新英文首页特性列表**
  - 文件：`docs/docs/index.md`
  - 更新 `features` 数组，体现以下新亮点：
    1. 鼠标 + 键盘双重支持（Mouse & Keyboard）
    2. 多语言国际化（13+ 语言含 RTL）
    3. 波形可视化裁剪
    4. 专辑签名保护（防篡改 `.ktalbum`）
  - 保留已有条目的核心信息，避免大幅删减
  - 验证：页面 feature cards 渲染正常

- [ ] **C2. 更新中文首页特性列表**
  - 文件：`docs/docs/zh/index.md`
  - 内容与 C1 完全对应，使用中文
  - 验证：中文首页 feature cards 渲染正常

---

## 阶段 D — 安装指南（可并行）

- [ ] **D1. 更新安装指南（英文）**
  - 文件：`docs/docs/guide/getting-started/installation/index.md`
  - 修改"System Requirements"系统要求小节，注明 macOS 已获完整支持
  - 验证：页面内容准确，无死链

- [ ] **D2. 更新安装指南（中文）**
  - 文件：`docs/docs/zh/guide/getting-started/installation/index.md`
  - 内容与 D1 完全对应，使用中文
  - 验证：页面内容准确，无死链

---

## 阶段 E — 快速入门（可并行）

- [ ] **E1. 更新快速入门页面（英文）**
  - 文件：`docs/docs/guide/getting-started/quick-start/index.md`
  - 内容补充：
    - 音量调节：提及全局音量 / 主页面随机音量变化选项
    - 语言设置：提及 13+ 语言可选
    - 键鼠分离模式入口说明（分离后可独立选择鼠标专辑）
  - 验证：内容逻辑连贯

- [ ] **E2. 更新快速入门页面（中文）**
  - 文件：`docs/docs/zh/guide/getting-started/quick-start/index.md`
  - 内容与 E1 完全对应，使用中文
  - 验证：内容逻辑连贯

---

## 阶段 F — 键音专辑介绍（可并行）

- [ ] **F1. 更新键音专辑介绍页面（英文）**
  - 文件：`docs/docs/guide/key-package/introduction/index.md`
  - 在"安全机制"部分或 `.ktalbum` 说明处补充：签名保护体系（防篡改）、配置加密说明
  - 验证：内容与中文版一致

- [ ] **F2. 更新键音专辑介绍页面（中文）**
  - 文件：`docs/docs/zh/guide/key-package/introduction/index.md`
  - 内容与 F1 完全对应，使用中文
  - 验证：内容与英文版一致

---

## 阶段 G — Step 2 裁剪定义声音（可并行）

- [ ] **G1. 更新"裁剪定义声音"指南（英文）**
  - 文件：`docs/docs/guide/key-package/裁剪定义声音/index.md`
  - 新增章节或扩充现有内容：
    - 波形可视化裁剪（Waveform Trimmer）：拖拽选区、与时间输入双向同步
    - 前端试听播放条（无需 SDK）
    - 音量指针条（dB 可视调节，范围 ±18 dB）
    - 注明：数字输入裁剪能力保留，波形组件为增强选项
  - 验证：新增章节标题层级符合当前文档 outline 配置（h2/h3/h4）

- [ ] **G2. 更新"裁剪定义声音"指南（中文）**
  - 文件：`docs/docs/zh/guide/key-package/裁剪定义声音/index.md`
  - 内容与 G1 完全对应，使用中文
  - 验证：与英文版章节结构一致

---

## 阶段 H — Step 3 铸造至臻键音（可并行）

- [ ] **H1. 更新"铸造至臻键音"指南（英文）**
  - 文件：`docs/docs/guide/key-package/铸造至臻键音/index.md`
  - 新增内容：
    - 按下/抬起独立音量控制（Press/Release Volume）：可为按下和抬起声音分别设置独立音量（dB）
    - 说明独立音量与全局音量的叠加关系
  - 验证：新增段落不与现有"Diverse Playback Modes"等章节冲突

- [ ] **H2. 更新"铸造至臻键音"指南（中文）**
  - 文件：`docs/docs/zh/guide/key-package/铸造至臻键音/index.md`
  - 内容与 H1 完全对应，使用中文
  - 验证：与英文版语义完全一致

---

## 阶段 I — Step 4 按键联动声效（可并行）

- [ ] **I1. 更新"按键联动声效"指南（英文）**
  - 文件：`docs/docs/guide/key-package/按键联动声效/index.md`
  - 新增内容：
    - 鼠标按键支持：可将音效绑定到鼠标左键/右键/中键/滚轮的按下/抬起事件
    - 键鼠分离模式：支持分别指定键盘专辑和鼠标专辑（主页面入口介绍）
    - 说明：增加了对鼠标按键音效的独立配置（与键盘键位配置并列，同属键音专辑内容）
  - 验证：新增内容不破坏现有"Flexible Binding Strategies"章节的逻辑

- [ ] **I2. 更新"按键联动声效"指南（中文）**
  - 文件：`docs/docs/zh/guide/key-package/按键联动声效/index.md`
  - 内容与 I1 完全对应，使用中文
  - 验证：与英文版语义完全一致

---

## 阶段 J — 验证（依赖所有前置任务完成）

- [ ] **J1. 中英文内容一致性人工比对**
  - 逐页核对中英文版本，确认：章节数量一致、核心信息无遗漏、无仅出现在单语言版本的段落
  - 验证清单：A1-A2、B1-B2、C1-C2、D1-D2、E1-E2、F1-F2、G1-G2、H1-H2、I1-I2

- [ ] **J2. 本地构建验证（可选）**
  - 在 `docs/` 目录运行 `npm run docs:build`（或 `vitepress build`）检查无构建报错
  - 验证：changelog 侧边栏 v1.0.0 条目可点击并正常跳转
  - 验证：所有内部链接无 404

---

## 依赖关系

```
A（config）→ B（changelog）
C / D / E / F / G / H / I（各指南页面，可并行）
所有 → J（验证）
```

> A1/A2、B1/B2 应先行，其余阶段 C~I 可并行推进。
