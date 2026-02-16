## 1. 实施
- [x] 1.1 为 `ManageAudioFilesDialog` 的“选择要管理的源文件”选择器接入输入过滤能力。
- [x] 1.2 为 `CreateSoundDialog` 的音频源文件选择器接入输入过滤能力。
- [x] 1.3 为 `EditSoundDialog` 的“选择要管理的声音”与“音频源文件”两个选择器接入输入过滤能力。
- [x] 1.4 为 `CreateKeySoundDialog` 中 down/up 两个“选择声音”选择器接入输入过滤能力。
- [x] 1.5 为 `EditKeySoundDialog` 中 down/up 两个“选择声音”选择器接入输入过滤能力。
- [x] 1.6 为 `EveryKeyEffectDialog` 中全键 down/up 两个选择器接入输入过滤能力。
- [x] 1.7 为 `AddSingleKeyEffectSubDialog` 中单键 down/up 两个选择器接入输入过滤能力。
- [x] 1.8 为 `EditSingleKeyEffectSubDialog` 中单键 down/up 两个选择器接入输入过滤能力（与新增路径体验一致）。

## 2. 一致性与回归
- [x] 2.1 统一并复核选择器过滤规则（大小写不敏感、空输入回退全量、不改变已选值结构）。
- [x] 2.2 验证依赖告警插槽、类型筛选、锚定联动、播放模式约束、保存/删除提示不受影响。（通过代码路径核对与类型/语法检查）
- [ ] 2.3 验证 Step1~Step4 全流程可完成（创建/编辑/保存）且配置写回路径不变。

## 3. 验证
- [x] 3.1 运行前端静态检查（如 `npm run lint` 或项目既有等价检查）。
- [ ] 3.2 进行手工验收并记录：每个目标选择器均支持搜索、无匹配时展示空态/无结果提示（遵循现有 UI 能力）。
