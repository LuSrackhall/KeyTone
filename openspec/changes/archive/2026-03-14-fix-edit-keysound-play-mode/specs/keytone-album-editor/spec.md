# keytone-album-editor Delta

## MODIFIED Requirements

### Requirement: UI → Config 写回闭环不变

Normative: 在进行组件拆分重构后，键音专辑编辑器（KeyTone Album Editor）在用户侧的视觉、交互与配置读写行为 SHALL 与重构前保持一致。

#### Scenario: UI → Config 写回闭环不变

- **GIVEN** 用户在 UI 中修改名称/声音/键音/联动设置等配置
- **WHEN** watch/按钮触发写回
- **THEN** `ConfigSet`/`ConfigDelete` 的 key 路径、触发时机、以及通知提示 SHALL 与重构前一致

#### Scenario: 编辑已有高级键音播放模式可生效

- **GIVEN** 用户在 Step3 的“编辑已有高级键音”对话框中选择某个按键音
- **WHEN** 用户修改 down/up 的播放模式并点击保存
- **THEN** 配置写回 SHALL 使用最新的播放模式值，重新打开编辑对话框时模式显示与实际配置一致
