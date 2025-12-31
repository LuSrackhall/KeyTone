# 键音专辑编辑器组件拆分重构规格（增量）

## Purpose

本增量规格描述对 `Keytone_album.vue` 的结构性重构要求：**拆分组件以提升可维护性与复用性，同时不改变任何现有用户可感知行为**。

---

## MODIFIED Requirements

### Requirement: 键音专辑编辑器行为保持不变

Normative: 在进行组件拆分重构后，键音专辑编辑器（KeyTone Album Editor）在用户侧的视觉、交互与配置读写行为 SHALL 与重构前保持一致。

#### Scenario: Stepper 折叠/展开行为不变

- **GIVEN** 编辑器使用 `q-stepper` 组织多个步骤
- **WHEN** 用户点击某个 step 的 header
- **THEN** step SHALL 在“展开该 step”与“折叠为 step=99”之间切换

#### Scenario: Continue/Back 导航行为不变

- **GIVEN** 用户位于任意步骤
- **WHEN** 用户点击 Continue / Back
- **THEN** step 值切换与原实现一致，且不会引入额外的导航路径

#### Scenario: SSE → UI 状态映射不变

- **GIVEN** 前端通过 `messageAudioPackage` 事件接收来自后端的键音包配置数据
- **WHEN** 配置数据发生变化
- **THEN** UI 中与以下领域相关的数据映射与排序 SHALL 与重构前一致：
  - audio_files → soundFileList
  - sounds → soundList
  - key_sounds → keySoundList
  - key_tone（global/single/is_enable_embedded_test_sound）相关字段

#### Scenario: UI → Config 写回闭环不变

- **GIVEN** 用户在 UI 中修改名称/声音/键音/联动设置等配置
- **WHEN** watch/按钮触发写回
- **THEN** `ConfigSet`/`ConfigDelete` 的 key 路径、触发时机、以及通知提示 SHALL 与重构前一致

---

### Requirement: 对话框可复用且可被跨步骤直接拉起

Normative: 拆分后的对话框组件 SHALL 支持被任意步骤直接拉起使用（例如在某一步中打开另一领域的配置对话框），并保持与原行为一致。

#### Scenario: 在 Step4 拉起全局设置对话框

- **GIVEN** 用户处于联动声效步骤
- **WHEN** 用户点击“全局设置”按钮
- **THEN** 全局设置对话框打开，且其字段绑定/保存行为与重构前一致

#### Scenario: 在 Step4 拉起单键设置对话框

- **GIVEN** 用户处于联动声效步骤
- **WHEN** 用户点击“单键设置”按钮
- **THEN** 单键设置对话框打开，且以下行为 SHALL 与重构前一致：
  - 添加单键声效：按键多选/搜索/录制、声效选择、保存后的提示与状态重置
  - 编辑单键声效：切换不同按键时的初始化逻辑、保存/删除提示、无改动提示

> 实现说明：父组件中原先用于对照回滚的 `v-if="false"` 旧单键对话框模板已移除；当前以 `SingleKeyEffectDialog` 作为唯一实现入口。

---

## ADDED Requirements

### Requirement: 组件边界与目录结构清晰

Normative: 为降低单文件复杂度，键音专辑编辑器前端实现 SHALL 采用清晰的目录结构将步骤、对话框与逻辑域拆分管理，并保持对外入口兼容。

#### Scenario: 保持外部引用入口兼容

- **GIVEN** 其他页面/组件可能引用 `frontend/src/components/Keytone_album.vue`
- **WHEN** 本次重构落地
- **THEN** 外部引用 SHALL 继续可用（通过薄壳组件或兼容导出），无需调用方修改行为

### Requirement: 新增拆分文件具备可读性说明

Normative: 本次重构新增的 step/dialog/composable/mapper 文件 SHOULD 在文件头部包含说明注释，用于降低 review 与排错成本。

#### Scenario: 排查行为回归时的定位成本可控

- **GIVEN** 用户在重构后遇到行为回归或数据不同步
- **WHEN** 开发者打开新增的拆分文件
- **THEN** 文件头部注释 SHOULD 明确：文件职责边界、关联的调用方文件、关键行为不变约束、以及首选的 Debug 切入点
