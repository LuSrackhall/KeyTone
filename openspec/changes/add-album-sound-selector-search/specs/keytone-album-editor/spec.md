# keytone-album-editor Spec Delta (add-album-sound-selector-search)

## ADDED Requirements

### Requirement: 键音专辑关键选择器支持输入搜索

Normative: 键音专辑编辑器在用户点名的关键选择器上 SHALL 提供输入搜索能力；搜索行为 SHALL 基于当前选项展示标签进行匹配，并在输入为空时恢复全量候选项。

#### Scenario: Step1 管理源文件选择器支持搜索
- **GIVEN** 用户打开“管理已载入的源文件”对话框
- **WHEN** 用户在“选择要管理的源文件”选择器输入关键字
- **THEN** 候选项 SHALL 按关键字即时过滤
- **AND** 清空输入后 SHALL 恢复全量候选项

#### Scenario: Step2 创建与编辑声音选择器支持搜索
- **GIVEN** 用户位于 Step2 并打开“创建新的声音”或“编辑已有声音”对话框
- **WHEN** 用户在音频源文件选择器或“选择要管理的声音”选择器输入关键字
- **THEN** 对应候选项 SHALL 按关键字过滤并支持快速定位目标项

---

### Requirement: Step3 制作按键音中的四个声音选择器支持搜索

Normative: Step3（制作按键音）中四个“选择声音”选择器（创建 down/up、编辑 down/up）SHALL 提供一致的输入搜索能力。

#### Scenario: Step3 四个声音选择器均可搜索
- **GIVEN** 用户在 Step3 分别进入“创建按键音”与“编辑按键音”流程
- **WHEN** 用户在 down/up 的“选择声音”选择器中输入关键字
- **THEN** 四个选择器 SHALL 均能按关键字过滤候选项
- **AND** 搜索结果 SHALL 不影响已选值与播放模式约束

---

### Requirement: Step4 按键联动生效中的四个声音选择器支持搜索

Normative: Step4（按键联动生效）中四个声音选择器（全键 down/up、单键新增 down/up）SHALL 提供一致的输入搜索能力。

#### Scenario: Step4 四个核心声音选择器均可搜索
- **GIVEN** 用户在 Step4 分别打开“全局设置”与“添加单键声效”
- **WHEN** 用户在 down/up 声音选择器中输入关键字
- **THEN** 四个选择器 SHALL 均能按关键字过滤候选项并快速选择

#### Scenario: 单键编辑路径保持同等搜索体验
- **GIVEN** 用户在 Step4 打开“编辑单键声效”子对话框
- **WHEN** 用户在 down/up 声音选择器中输入关键字
- **THEN** 编辑路径 SHALL 与新增路径保持同等搜索能力与匹配规则

---

### Requirement: 搜索增强不改变既有配置语义

Normative: 选择器搜索能力属于交互增强，SHALL NOT 改变现有配置数据结构、依赖告警展示、保存/删除行为、以及配置写回路径。

#### Scenario: 搜索后保存行为与未搜索一致
- **GIVEN** 用户在任一目标选择器通过搜索完成选择并执行保存
- **WHEN** 系统写回配置并刷新界面
- **THEN** 写回键路径、值结构与提示行为 SHALL 与未启用搜索时一致
- **AND** 非目标功能（如类型筛选、锚定联动、依赖告警）SHALL 保持原行为
