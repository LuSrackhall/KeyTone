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

Non-normative: 当用户拖拽播放头或拖拽选区两侧指针时，波形 SHOULD 提供“边缘触发”的自动滚动体验（指针靠近视窗左右边缘才滚动，速度渐进），以便在高 zoom 下持续拖动至视野外区域。
Non-normative: 自动滚动过程中播放头 SHOULD 保持与光标位置一致（避免出现“滚动后播放头跳跃/不跟手”）。
Non-normative: 自动滚动过程中，选区两侧指针 SHOULD 与光标保持一致；必要时可直接更新 region 并回写 start/end，以避免内部拖拽逻辑滞后。
Non-normative: 在边缘渐进区（滚动速度很小）也 SHOULD 持续更新光标对应的时间点，避免指针卡在临界线位置。
Non-normative: 左键拖拽场景下（播放头/选区指针），应在 pointermove 与自动滚动帧内持续更新位置，保证全程跟手。
Non-normative: 为保证拖拽丝滑一致，组件可接管左键拖拽逻辑（基于命中元素判定播放头/选区/指针），避免与第三方库内部拖拽冲突。
Non-normative: 命中识别可基于 composedPath 中的 `part` token（region / region-handle-left/right），保证识别稳定。
Non-normative: 自动滚动速度 SHOULD 进行平滑（例如一阶低通），并使用预测的 scrollLeft 计算光标时间点，以保证跟手与滚动同步。

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
