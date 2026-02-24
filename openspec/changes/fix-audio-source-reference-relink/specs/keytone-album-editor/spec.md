## ADDED Requirements

### Requirement: 音频源别名标识使用唯一 UUID

Normative: 对于同一 `sha256` 的音频源文件，系统在新增别名时 SHALL 生成唯一 UUID 作为 `name_id`；删除旧别名后，后续新增 SHALL NOT 复用历史 `name_id`。

#### Scenario: 删除后重导入同一文件不会复用旧标识

- **GIVEN** 用户曾导入某音频文件（`sha256=A`）并产生别名 `name_id=U1`
- **AND** 该别名已被删除
- **WHEN** 用户再次导入同一音频文件（同 `sha256=A`）
- **THEN** 新别名的 `name_id` MUST 为新的 UUID
- **AND** 新别名 SHALL NOT 等于 `U1`

#### Scenario: 连续导入同一文件时标识均唯一

- **GIVEN** 用户连续多次导入同一 `sha256` 文件
- **WHEN** 系统为每次导入创建别名
- **THEN** 每个别名的 `name_id` SHALL 为不同 UUID
- **AND** 这些别名仍可共同引用同一物理音频文件（不额外复制）

### Requirement: 缺失源文件引用不得自动复联

Normative: 当裁剪定义声音或键音配置引用的 `source_file_for_sound`（`sha256 + name_id + type`）在当前 `audio_files` 中不存在时，系统 SHALL 将其视为缺失引用，并保持缺失状态，直到用户显式重新选择并保存新的源文件引用。

#### Scenario: 重导入同文件不自动恢复裁剪声音依赖

- **GIVEN** 某裁剪定义声音引用 `sha256=A, name_id=0, type=.wav`
- **AND** 该引用对应的音频源别名已删除，当前状态显示为缺失依赖
- **WHEN** 用户仅重导入同一文件（同 `sha256=A`）但未在编辑器中重新选择源文件
- **THEN** 该裁剪定义声音的依赖状态 MUST 继续为缺失
- **AND** 系统 MUST NOT 自动将其绑定到新导入的别名

#### Scenario: 仅显式重选后才恢复依赖

- **GIVEN** 裁剪定义声音当前为缺失依赖
- **WHEN** 用户在编辑对话框中手动选择可用源文件并保存
- **THEN** 该声音依赖 SHALL 恢复为可用

### Requirement: 播放链路遵循严格引用一致性

Normative: 运行时播放解析 SHALL 与依赖语义一致；当配置中的 `sha256 + name_id + type` 引用不存在时，播放链路 MUST 按缺失引用处理，SHALL NOT 仅依据 `sha256 + type` 推断为可播放。

#### Scenario: 缺失引用时拒绝隐式回连播放

- **GIVEN** 配置中某声音引用的 `name_id` 不存在于当前 `audio_files.<sha256>.name`
- **WHEN** 触发该声音播放
- **THEN** 播放链路 MUST 识别为缺失引用并拒绝隐式恢复
- **AND** 结果与依赖检查结论保持一致
