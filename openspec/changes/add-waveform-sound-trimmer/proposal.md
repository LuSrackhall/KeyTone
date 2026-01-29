# Change: add-waveform-sound-trimmer（为声音裁剪提供波形可视化选择）

## Why

当前“裁剪定义声音”的交互主要依赖开始/结束时间数字输入。用户无法直观看到音频里“哪里有声音”，只能反复试错，学习成本高且效率低。

## What Changes

- 在键音专辑编辑器 Step2「定义声音」的“创建声音 / 编辑声音”对话框中，新增一个“波形可视化裁剪”组件：
  - 展示所选音频源文件的波形图（可缩放/滚动）。
  - 通过拖拽选区（region）可视化地选择裁剪区间。
  - 选区与现有开始/结束时间输入框双向同步（任一侧修改，另一侧即时更新）。
- 保持现有“纯数字输入裁剪”能力不变：波形组件是增强能力，而不是替换（避免破坏老用户工作流）。
- 为波形渲染提供稳定的数据来源：后端仅提供音频流（按 sha256+type 定位当前专辑下的音频源文件），前端负责解码与渲染波形。

## Impact

- Affected specs:
  - keytone-album-editor（新增“波形可视化裁剪”需求）
- Affected code (implementation stage, non-exhaustive):
  - frontend/src/components/keytone-album/dialogs/CreateSoundDialog.vue
  - frontend/src/components/keytone-album/dialogs/EditSoundDialog.vue
  - frontend/src/components/Keytone_album.vue（Context 与预览流程复用）
  - sdk/server/server.go（新增音频流接口）

## Non-goals (本次不做)

- 不做完整 DAW 级编辑（多轨、淡入淡出曲线、频谱编辑等）。
- 不改变 sounds.cut 的数据结构与存储格式。
- 不要求用户必须使用波形组件才能完成裁剪。
