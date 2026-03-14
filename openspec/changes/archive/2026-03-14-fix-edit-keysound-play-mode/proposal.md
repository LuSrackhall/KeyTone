# Change: 修复“编辑已有高级键音”播放模式不生效

## Why

编辑已有高级键音时，切换播放模式（单次/随机/循环）保存后不生效，属于明显回归/缺陷，会导致用户误以为修改成功但实际配置未更新。

## What Changes

- 修复编辑已有高级键音时播放模式保存路径，确保保存后配置正确写回。
- 在编辑对话框中显式兼容播放模式字段的两种结构（字符串/对象），避免历史数据或 UI 状态导致保存值为空。
- 增加关键注释，说明播放模式的真实数据结构与保存策略。

## Impact

- Affected specs: keytone-album-editor
- Affected code: frontend/src/components/keytone-album/dialogs/EditKeySoundDialog.vue, frontend/src/components/Keytone_album.vue
