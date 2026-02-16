# Change: 为键音专辑关键选择器补充搜索能力

## Why
键音专辑编辑流程中多个 `q-select` 选择器在选项数量变大后定位效率明显下降，用户需要通过滚动逐项查找，导致创建与编辑流程耗时增加、误选概率提高。

本次变更目标是在不改变现有配置语义与保存路径的前提下，为用户明确点名的选择器补充“输入即过滤”的搜索能力，提升可用性并降低操作成本。

## What Changes
- 为 Step1「管理已载入的源文件」对话框中的“选择要管理的源文件”选择器增加搜索。
- 为 Step2「裁剪定义声音」与「编辑已有声音」中涉及的音频源文件/声音选择器增加搜索：
  - 创建声音对话框：音频源文件选择器。
  - 编辑声音对话框：要管理的声音选择器、音频源文件选择器。
- 为 Step3「制作按键音」中的四个声音选择器增加搜索（创建 down/up、编辑 down/up）。
- 为 Step4「按键联动生效」中的四个声音选择器增加搜索（全键 down/up、单键新增 down/up）；并要求单键编辑子对话框保持同等搜索体验，避免同域行为分裂。
- 明确兼容性约束：仅增强筛选交互，不改变现有 option 数据结构、`v-model` 值结构、依赖告警展示、保存/删除逻辑与配置写回键路径。

## Impact
- Affected specs: `keytone-album-editor`
- Affected code (planned):
  - `frontend/src/components/keytone-album/dialogs/ManageAudioFilesDialog.vue`
  - `frontend/src/components/keytone-album/dialogs/CreateSoundDialog.vue`
  - `frontend/src/components/keytone-album/dialogs/EditSoundDialog.vue`
  - `frontend/src/components/keytone-album/dialogs/CreateKeySoundDialog.vue`
  - `frontend/src/components/keytone-album/dialogs/EditKeySoundDialog.vue`
  - `frontend/src/components/keytone-album/dialogs/EveryKeyEffectDialog.vue`
  - `frontend/src/components/keytone-album/dialogs/AddSingleKeyEffectSubDialog.vue`
  - `frontend/src/components/keytone-album/dialogs/EditSingleKeyEffectSubDialog.vue`
- Breaking change: 无
