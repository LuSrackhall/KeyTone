## 1. Spec & Design

- [x] 1.1 完成 `playback-album-routing` 的随机音量能力增量 spec
- [x] 1.2 确认随机算法与参数边界（随机扣减、默认值=3、仅非负）
- [x] 1.3 确认 UI 位置与门控规则（主页面相关设置末尾、全局+按下/抬起随机单独控制）

## 2. Frontend State & Persistence

- [x] 2.1 在 `setting-store` 新增 `mainHome.randomVolumeProcessing` 状态结构
- [x] 2.2 将 `mainHome.randomVolumeProcessing.maxReduceRatio` 默认值改为 `3`
- [x] 2.3 新增 `mainHome.pressReleaseRandomVolumeProcessing` 状态结构（总开关 + 6节点）
- [x] 2.4 在 `getConfigFileToUi` 中补充读取与默认值回填（旧配置增量补齐）
- [x] 2.5 增加 watch 持久化：全局随机与按下/抬起随机所有键

## 3. Frontend UI (MainHome Setting)

- [x] 3.1 在主页面相关设置末尾新增“随机音量”开关项
- [x] 3.2 新增“随机降幅上限”输入（仅开关开启时展示，默认重置为 `3`）
- [x] 3.3 新增“按下/抬起随机音量单独控制”开关与展开项
- [x] 3.4 展开项提供六个事件态节点开关，节点开关开启后显示对应输入框
- [x] 3.5 参数边界校验与重置交互（仅限制非负，节点默认重置为 `3`）

## 4. SDK Config & Runtime

- [x] 4.1 在 `sdk/config/config.go` 新增随机音量默认配置常量与 `viper.SetDefault`
- [x] 4.2 将全局随机默认值更新为 `3`
- [x] 4.3 在 `sdk/config/config.go` 新增按下/抬起随机音量单独控制默认键
- [x] 4.4 在 `sdk/keySound/keySound.go` 新增按下/抬起随机叠加函数（节点开关门控）
- [x] 4.5 将链路扩展为“全局随机后再按下/抬起随机”，支持双层叠加

## 5. Localization

- [x] 5.1 新增 zh-CN 文案键（随机音量开关、说明、随机降幅上限）
- [x] 5.2 同步 en-US 最低必需文案，避免缺键回退异常
- [x] 5.3 新增按下/抬起随机音量单独控制相关文案（总开关、展开项、6节点）
- [x] 5.4 文案命名规范化：全局/单独控制域显式区分，节点输入采用“随机降幅上限(场景)”格式
- [x] 5.5 说明文案语义优化：明确“生效过程与全局随机设置叠加”

## 6. Validation

- [x] 6.1 运行 `openspec validate add-main-home-random-volume --strict --no-interactive`
- [x] 6.2 代码路径验证：关闭开关时随机层直接透传，不影响原链路
- [x] 6.3 公式与实现验证：开启后随机层 `deltaVolume<=0`，每次播放不高于当前实际音量
- [x] 6.4 边界验证：`maxReduceRatio` 在前后端仅限制为非负，不设上限
- [x] 6.5 叠加验证：全局随机与按下/抬起随机同时开启时可叠加两次随机降幅
