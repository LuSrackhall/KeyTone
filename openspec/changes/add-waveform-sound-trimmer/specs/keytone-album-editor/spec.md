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

#### Scenario: 波形不可用时可降级完成流程

- **GIVEN** 波形数据加载/解码失败（或源文件不可访问）
- **WHEN** 用户输入 start/end 并点击预览/保存
- **THEN** 系统 SHALL 继续执行预览/保存逻辑
- **AND** UI SHOULD 明确提示“波形不可用，已降级为手动裁剪”

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
