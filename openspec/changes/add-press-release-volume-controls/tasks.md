# Tasks

## 1. Spec & Design

- [x] 1.1 完成 `playback-album-routing` 的 down/up 音量能力增量 spec
- [x] 1.2 明确配置键命名、默认值和迁移策略（仅补充，不破坏旧键）
- [x] 1.3 确认 UI 门控规则：仅在“按下/抬起单独控制”开启时展示新增模块

## 2. Frontend State & Persistence

- [x] 2.1 在 `setting-store` 新增“按下/抬起单独控制”开关状态（默认 false）
- [x] 2.2 在 `setting-store` 新增全局 down/up 音量状态结构（滑块值、降幅、调试开关、静音状态）
- [x] 2.3 在 `setting-store` 新增分离模式 keyboard/mouse 的 down/up 音量状态结构
- [x] 2.4 在 `getConfigFileToUi` 增加读取与默认值回填（缺失时写默认）
- [x] 2.5 增加 watch 持久化，确保仅新增新键，不影响现有键行为

## 3. Frontend UI (MainHome Setting)

- [x] 3.1 新增“按下/抬起音量单独控制”开关项
- [x] 3.2 在“鼠标回退到键盘”设置项下新增“按下/抬起音量单独控制”独立展开栏
- [x] 3.3 在该独立展开栏中新增全局/键盘/鼠标 down/up 子模块 UI（滑块集中、降幅输入集中、调试开关集中）
- [x] 3.4 新增 UI 展示门控：开关关闭时隐藏所有 down/up 新增模块
- [x] 3.5 保持独立展开项样式与“分离模式音量设置”一致（缩进、竖条、调试滑块布局），且不显示分组小标题

## 4. SDK Config & Runtime

- [x] 4.1 在 `sdk/config/config.go` 新增 down/up 默认配置常量与 `viper.SetDefault` 注册
- [x] 4.2 在 `sdk/keySound/keySound.go` 基于 `keyState` 与设备类型叠加 down/up 音量层
- [x] 4.3 维持原有播放链路优先级，确保未开启开关时行为与当前版本一致

## 5. Localization

- [x] 5.1 新增 zh-CN 文案键：单独控制开关、全局/键盘/鼠标 down/up 标题与 caption
- [x] 5.2 同步 en-US 最低必需文案，避免缺键回退异常

## 6. Validation

- [x] 6.1 运行 `openspec validate add-press-release-volume-controls --strict --no-interactive`
- [ ] 6.2 手动验证：开关关闭时新增模块不展示，开启后按层级展示完整
- [ ] 6.3 手动验证：全局、分离、down/up 三层叠加在键盘/鼠标按下与抬起场景正确生效
- [ ] 6.4 手动验证：旧用户配置升级后无关设置不受影响
