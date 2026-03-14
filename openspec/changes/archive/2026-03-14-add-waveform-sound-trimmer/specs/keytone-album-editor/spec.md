# keytone-album-editor Spec Delta (add-waveform-sound-trimmer)

## ADDED Requirements

### Requirement: 在“定义声音裁剪”提供波形可视化选区

Normative: 在键音专辑编辑器 Step2「定义声音」中，系统 SHALL 提供波形可视化组件，以协助用户为音频源文件选择裁剪区间（start/end）。

#### Scenario: 创建声音时可视化选择裁剪区间

- **GIVEN** 用户打开“创建新的声音”对话框
- **AND** 用户选择了一个音频源文件
- **WHEN** 波形数据可用
- **THEN** UI SHALL 展示该源文件的波形图
- **AND** 用户 SHALL 可以通过拖拽选区（region）设置裁剪区间
- **AND** 选区变化 SHALL 同步更新开始/结束时间输入框的值

#### Scenario: 编辑声音时可视化调整裁剪区间

- **GIVEN** 用户打开“编辑已有声音”对话框并选择了一个声音
- **WHEN** 波形数据可用
- **THEN** UI SHALL 展示该声音对应源文件的波形图
- **AND** 初始选区 SHALL 对齐到该声音当前的 start/end
- **AND** 用户调整选区后，修改 SHALL 可被保存为该声音的新 cut 参数

---

### Requirement: 波形组件支持缩放/滚动并提供试听播放条

Normative: 波形组件 SHALL 提供缩放（zoom）与横向滚动能力，并 SHALL 提供前端试听播放条以支持暂停与拖动播放头；默认播放范围 SHALL 为“播放全部”，并允许切换为“播放选区”。

Normative: 波形组件在首次加载成功后，其默认 zoom（minPxPerSec） SHALL 为 `50`。

Non-normative: macOS 平台的快捷缩放提示可显示为“Control/Command + 滚轮或触控板捏合”，并优先使用 getModifierState('Control') 识别外接键盘的 Control 键。
Non-normative: 若 wheel 事件无法反映 ctrlKey，可在 window 层跟踪 Control/Command 按键状态作为兜底。

Non-normative: 当用户拖拽播放头或拖拽选区两侧指针时，波形 SHOULD 提供“边缘触发”的自动滚动体验（指针靠近视窗左右边缘才滚动，速度渐进），以便在高 zoom 下持续拖动至视野外区域。
Non-normative: 自动滚动过程中播放头 SHOULD 保持与光标位置一致（避免出现“滚动后播放头跳跃/不跟手”）。
Non-normative: 自动滚动过程中，选区两侧指针 SHOULD 与光标保持一致；必要时可直接更新 region 并回写 start/end，以避免内部拖拽逻辑滞后。
Non-normative: 在边缘渐进区（滚动速度很小）也 SHOULD 持续更新光标对应的时间点，避免指针卡在临界线位置。
Non-normative: 左键拖拽场景下（播放头/选区指针），应在 pointermove 与自动滚动帧内持续更新位置，保证全程跟手。
Non-normative: 为保证拖拽丝滑一致，组件可接管左键拖拽逻辑（基于命中元素判定播放头/选区/指针），避免与第三方库内部拖拽冲突。
Non-normative: 命中识别可基于 composedPath 中的 `part` token（region / region-handle-left/right），保证识别稳定。
Non-normative: 自动滚动速度 SHOULD 进行平滑（例如一阶低通），并使用预测的 scrollLeft 计算光标时间点，以保证跟手与滚动同步。

Non-normative: 音量 UI SHOULD 以 dB 展示（内部仍使用 `cut.volume` 原始尺度），并提供 ±18 dB 的可视范围：
- 左侧可展示刻度 18/12/6/0/-6/-12/-18（含负刻度）；
- 标尺建议放在波形外侧（不占用波形宽度）。
  - 对话框宽度 MAY 适度扩展以容纳标尺，但必须不超过视口可视范围。
  - 在宽度受限场景下，实现 SHOULD 将标尺布局在对话框内部（例如三栏布局），避免负向定位溢出对话框外。
  - 标尺应靠近波形边界，便于对齐观察；
  - 实现 SHOULD 避免绘制贯穿整个横向区域的刻度/网格线；仅显示刻度数字（及可选的短刻度）即可。
- 右侧展示当前 dB（可超过 ±18）；
- 超出范围时指示线贴边并显示普通提示（非错误/警告）。

#### Scenario: 缩放/滚动辅助精确定位

- **GIVEN** 用户已加载波形
- **WHEN** 用户调整 zoom
- **THEN** 波形视图 SHALL 改变缩放级别以支持更精细的时间观察（例如约 100ms 级别片段）
- **AND** 在 zoom 增大后，UI SHALL 支持横向滚动以浏览完整波形

#### Scenario: 前端试听播放条可暂停与拖动

- **GIVEN** 波形已加载且音频可播放
- **WHEN** 用户点击播放/暂停
- **THEN** 前端 SHALL 播放/暂停音频
- **AND** 用户拖动播放头后，播放位置 SHALL 跳转到对应时间

#### Scenario: 播放范围默认为全部且可切换选区

- **GIVEN** 用户打开对话框且波形可用
- **THEN** 播放范围默认 SHALL 为“播放全部”
- **WHEN** 用户切换为“播放选区”
- **AND** 当前存在有效选区（end > start）
- **THEN** 前端 SHALL 仅播放该选区范围

---

### Requirement: 波形区域提供快捷音量调节条

Normative: 波形区域 SHALL 提供一条横向音量指示线，用户可通过上下拖动快速调整当前裁剪的 `cut.volume`，并与数值输入框保持双向同步；当 `cut.volume == 0` 时，指示线在 UI 上 SHALL 位于中位（unity gain 参考线）。

Non-normative: 为符合剪辑软件直觉，音量调整过程中波形本体 SHOULD 提供适度的视觉反馈（例如纵向高度变化），但该反馈不得影响时间映射与选区可操作性。

Non-normative: 实现应避免对 canvas 施加 `transform: scaleY(...)` 造成的插值伪影（例如静音中位线变粗/发糊）；若使用 wavesurfer.js，优先考虑通过渲染参数（例如 `barHeight`）进行渲染阶段缩放并触发重绘。

#### Scenario: 拖动音量指示线修改音量

- **GIVEN** 用户打开“创建声音/编辑声音”对话框且波形组件可见
- **WHEN** 用户在音量指示线附近按下并向上拖动
- **THEN** 系统 SHALL 提升当前裁剪的 `cut.volume`
- **AND** 数值输入框 SHALL 同步更新

#### Scenario: 拖动命中区稳定且易用

- **GIVEN** 用户打开“创建声音/编辑声音”对话框且波形组件可见
- **WHEN** 用户尝试在音量指示线附近按下并拖动
- **THEN** 系统 SHOULD 提供足够大的命中区以保证拖动稳定
- **AND** 系统 SHOULD 在手势/系统中断（例如 pointercancel）时正确结束拖动状态，避免交互卡死或失效

---

### Requirement: 播放过程波形滚动行为可配置

Normative: 系统 SHALL 在设置页提供“键音专辑页相关”分类，其中包含“波形滚动行为”选项，允许用户在“分页式跳转”和“智能边缘推挤”之间切换。默认值 SHALL 为“分页式跳转”。

Normative: 波形滚动行为设置 SHALL 通过 setting-store 持久化，应用重启后保持用户选择。

#### Scenario: 用户切换波形滚动行为

- **GIVEN** 用户打开设置页并展开“键音专辑页相关”
- **WHEN** 用户将“波形滚动行为”从“分页式跳转”切换为“智能边缘推挤”
- **THEN** 波形组件的播放跟随行为 SHALL 立即变为边缘推挤模式
- **AND** 设置 SHALL 被持久化，重启后仍然生效

#### Scenario: 智能边缘推挤模式下播放头不抖动

- **GIVEN** 用户选择了“智能边缘推挤”模式
- **WHEN** 播放头接近右侧边缘并触发推挤
- **THEN** 视口 SHALL 平滑跟随，播放头不应出现可见的左右抖动

#### Scenario: 智能边缘推挤模式播放结束可见尾部

- **GIVEN** 用户使用“智能边缘推挤”模式播放到音频结尾
- **WHEN** 播放结束
- **THEN** 视口 SHALL 自然滚动至末尾，尾部波形可见且不应出现强制跳转

---

### Requirement: 对话框关闭时立即停止前端试听音频

Normative: 当用户关闭“创建声音/编辑声音”对话框时，前端试听音频 SHALL 立即停止，以符合用户直觉。

#### Scenario: 关闭创建声音对话框时停止播放

- **GIVEN** 用户正在“创建声音”对话框中播放前端试听音频
- **WHEN** 用户关闭对话框
- **THEN** 音频播放 SHALL 立即停止

#### Scenario: 关闭编辑声音对话框时停止播放

- **GIVEN** 用户正在“编辑声音”对话框中播放前端试听音频
- **WHEN** 用户关闭对话框
- **THEN** 音频播放 SHALL 立即停止

---

### Requirement: 关闭动画保持流畅

Normative: 关闭对话框时，波形组件的销毁或停止逻辑 SHALL 不阻塞关闭动画。

#### Scenario: 关闭对话框时动画可见

- **GIVEN** 用户在“创建/编辑声音”对话框中已加载波形组件
- **WHEN** 用户关闭对话框
- **THEN** 对话框关闭动画 SHALL 正常播放且无明显卡顿

#### Scenario: 手动移动播放头时立即可见

- **GIVEN** 用户在波形组件中用左键点击或拖拽播放头
- **WHEN** 播放头位置发生变化
- **THEN** UI SHALL 立即显示新的播放头位置（不应出现“不可见”或延迟）

#### Scenario: 点击停止后播放头回到起点

- **GIVEN** 用户正在播放或已暂停在任意位置
- **WHEN** 用户点击停止按钮
- **THEN** 播放 SHALL 停止且播放头 SHALL 立即回到 0 位置，并滚动到起点

#### Scenario: 手动滚动时播放头同步

- **GIVEN** 用户手动滚动波形视图
- **WHEN** 滚动位置发生变化
- **THEN** 播放头显示 SHALL 同步更新，不应固定在视野某一位置

#### Scenario: 暂停/停止后首次右键拖拽可见

- **GIVEN** 用户刚暂停或停止播放
- **WHEN** 用户首次右键拖拽进行快速选区
- **THEN** 选区 UI SHALL 在拖拽过程中持续可见（不应等到拖拽结束才显示）

#### Scenario: 右键拖拽时播放头持续更新

- **GIVEN** 用户在播放过程中进行右键快捷选区
- **WHEN** 播放继续进行
- **THEN** 播放头 UI SHALL 持续更新，不应冻结或停滞

#### Scenario: 播放头移出视野时有边缘指示

- **GIVEN** 用户手动滚动或播放头自然滚动导致播放头移出可视区域
- **WHEN** 播放头位于视野外
- **THEN** UI SHALL 在边缘显示指示（变短/变浅）提示播放头在屏幕外（不要求变色）

#### Scenario: 手动操作期间视区优先

- **GIVEN** 用户在播放过程中进行手动操作（右键快捷选区、拖动选区、拖动音量条、拖动滚动条等）
- **WHEN** 播放继续进行
- **THEN** 自动滚动（推挤/翻页）SHALL 暂停，以当前手动视区为主
- **AND** 播放头 UI SHALL 持续更新并允许移出视野

#### Scenario: 点击滚动条不触发 Seek

- **GIVEN** 用户点击或拖拽波形组件的横向滚动条（track/thumb）
- **WHEN** 指针按下并开始滚动
- **THEN** 波形内容区的 seek/播放头跳转 SHALL NOT 被触发

#### Scenario: 拖拽滚动条时滚动丝滑

- **GIVEN** 用户拖拽横向滚动条进行浏览
- **WHEN** 滚动持续发生
- **THEN** 滚动 SHOULD 保持足够细腻（不应出现明显“大步进”式跳动）

#### Scenario: 缩放锚点稳定且无可见漂移

- **GIVEN** 用户进行滑块缩放或 Ctrl+滚轮缩放
- **WHEN** 系统根据锚点修正滚动位置
- **THEN** 锚点 SHOULD 保持稳定（不应出现可见抖动）
- **AND** 波形渲染 SHOULD 保持稳定（不应出现明显闪烁/抖动）
- **AND** 快捷缩放过程中不应出现低频闪烁（例如周期性轻微闪一下）
- **AND** 快捷缩放过程中不应出现可感知的轻微抖动（例如滚动位置来回微调）
- **AND** 实现 MAY 在快捷缩放 burst 期间对波形宽度更新做合并/量化（可随 zoom 自适应步长），以减少闪烁
- **AND** 播放头 UI SHOULD 与波形背景保持一致，不应出现视觉偏移
- **AND** 系统 SHALL 在旧/新 zoom 坐标系之间保持一致的锚点换算（不应混用导致偏移）
- **AND** 对 Ctrl+滚轮场景，系统 SHOULD 基于同一 clientX 执行误差补偿（必要时迭代）
- **AND** 对连续 Ctrl+滚轮缩放，系统 SHOULD 在同一 burst 内固定锚点（timeSec + clientX）
- **AND** 系统 MAY 采用采样索引锚定（整数）以进一步确保锚点可复现并锁定同一波峰

#### Scenario: 手动交互结束后平滑追帧

- **GIVEN** 用户在播放过程中进行手动交互并松开（例如释放滚动条/结束拖拽）
- **WHEN** 系统恢复自动跟随（推挤/翻页）
- **THEN** 视区 SHOULD 以短促动画平滑追到实时播放位置（不应瞬间跳回）

#### Scenario: 向下拖动降低音量

- **GIVEN** 用户打开“创建声音/编辑声音”对话框且波形组件可见
- **WHEN** 用户在音量指示线附近按下并向下拖动
- **THEN** 系统 SHALL 降低当前裁剪的 `cut.volume`
- **AND** 数值输入框 SHALL 同步更新

---

---

### Requirement: 为波形渲染提供音频流数据源

Normative: 系统 SHALL 提供可被前端访问的音频流数据源，用于从当前专辑的音频源文件生成波形；该数据源应当能够通过源文件引用信息定位（例如 sha256 + type）。

#### Scenario: 前端请求音频流并渲染波形

- **GIVEN** 用户在 Step2 的创建或编辑声音对话框中选择了一个音频源文件
- **WHEN** 前端请求该源文件的音频流
- **THEN** 后端 SHALL 返回对应音频文件的内容（失败时返回明确错误）
- **AND** 前端 SHALL 能据此解码并渲染波形图

---

### Requirement: 波形组件与数字输入双向同步且保持可用

Normative: 波形组件与 start/end 数字输入 SHALL 双向同步；即使波形不可用，用户仍 SHALL 能仅通过数字输入完成裁剪定义。

#### Scenario: 数字输入更新选区

- **GIVEN** 波形组件已加载并存在一个可编辑选区
- **WHEN** 用户在输入框中修改 start/end
- **THEN** 波形选区边界 SHALL 同步更新到对应时间位置

---

### Requirement: 右键拖拽快速选区（无菜单）

Normative: 波形区域 SHALL 支持右键拖拽快速创建/调整选区：右键按下即设置 start，拖动过程中实时更新 end，松开右键时固定 end；该过程 SHALL 不显示任何菜单。

#### Scenario: 右键按下并拖动快速建立选区

- **GIVEN** 波形组件已加载且用户位于波形区域
- **WHEN** 用户右键按下
- **THEN** 系统 SHALL 立即将按下位置对应时间设置为选区 start
- **AND** 在用户保持右键按下并向右拖动时，系统 SHALL 实时更新选区 end
- **AND** 在用户保持右键按下并向右拖动时，系统 SHOULD 全程保持选区可见（必要时触发轻量重绘并自动滚动以跟随 end）
- **AND** 自动滚动 SHOULD 采用“边缘触发”模型：仅当指针靠近滚动视窗左右边缘时才滚动，且滚动速度可随靠近程度渐进，以避免视图跳动
- **AND** 当用户松开右键时，系统 SHALL 将松开位置对应时间固定为 end
- **AND** 在高缩放（zoom 很大）场景下，系统 SHOULD 确保新选区能立即可见（必要时触发一次重绘并滚动到选区范围），不应要求用户通过轻微调整 zoom 才能看到结果

#### Scenario: 波形不可用时可降级完成流程

- **GIVEN** 波形数据加载/解码失败（或源文件不可访问）
- **WHEN** 用户输入 start/end 并点击预览/保存
- **THEN** 系统 SHALL 继续执行预览/保存逻辑
- **AND** UI SHOULD 明确提示“波形不可用，已降级为手动裁剪”

---

### Requirement: 波形裁剪组件 UI 文案接入 i18n（至少中文与英文）

Normative: 波形裁剪组件新增的 UI 文案（播放/暂停/停止、缩放、播放范围、加载/降级提示、交互提示等） SHALL 接入 i18n，至少覆盖中文与英文；交互提示文案 SHOULD 包含“右键拖拽快速选区”等能力描述，并保持跨平台表述。

#### Scenario: 中文界面显示中文文案

- **GIVEN** 用户将界面语言设置为中文
- **WHEN** 用户打开“创建声音/编辑声音”对话框且波形组件可见
- **THEN** 播放/缩放/提示等文案 SHALL 显示为中文

#### Scenario: 英文界面显示英文文案

- **GIVEN** 用户将界面语言设置为英文
- **WHEN** 用户打开“创建声音/编辑声音”对话框且波形组件可见
- **THEN** 播放/缩放/提示等文案 SHALL 显示为英文

---

### Requirement: SDK 预览播放严格限制为 <= 5000ms

Normative: 当使用 SDK 预览播放模式时，系统 SHALL 严格限制播放区间长度 `<= 5000ms`；超过限制时系统 SHALL 拒绝执行并给出明确提示，同时前端 SHOULD 引导用户使用前端试听播放条。

#### Scenario: SDK 预览播放超时拒绝

- **GIVEN** 用户点击“预览”（SDK 预览模式）
- **WHEN** (end - start) > 5000ms
- **THEN** 后端 SHALL 拒绝该请求并返回明确错误
- **AND** 前端 SHOULD 提示“SDK预览仅支持<=5000ms”并引导使用前端试听播放条

---

### Requirement: 裁剪区间的基础校验与边界处理

Normative: 系统 SHALL 对裁剪区间进行基础校验，避免保存明显无效的时间范围。

#### Scenario: 防止 end <= start

- **GIVEN** 用户选择或输入了裁剪区间
- **WHEN** end <= start
- **THEN** 系统 SHALL 阻止保存并提示用户修正

#### Scenario: 防止负数时间

- **GIVEN** 用户选择或输入了裁剪区间
- **WHEN** start < 0 或 end < 0
- **THEN** 系统 SHALL 提示错误并要求修正

---

### Requirement: 前端试听音量与 SDK 预览一致

Normative: 前端试听播放的音量映射 SHALL 与 SDK 预览模式一致（Base=1.6 的指数曲线），不应额外 clamp 或引入与 SDK 不一致的增益处理。

#### Scenario: 前端试听音量与 SDK 预览一致

- **GIVEN** 用户将 `cut.volume` 设置为任意值（例如 5）
- **WHEN** 用户分别使用前端试听播放条与 SDK 预览进行试听
- **THEN** 两者的音量变化趋势与相对听感 SHOULD 保持一致（以 SDK 预览为准）
