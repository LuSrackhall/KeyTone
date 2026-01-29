# Tasks: add-waveform-sound-trimmer

> 目标：让用户在“定义声音裁剪”时可视化选择片段，减少试错。

## 1. Scope & UI 入口

- [ ] 在 Step2 定义声音的“创建声音/编辑声音”对话框中确认放置位置与布局（不压缩现有表单可用性）。
- [ ] 明确波形组件的最小交互：加载、缩放/滚动、拖拽选区、选区时间显示。

## 2. 波形数据来源（选一条 MVP 路线）

- [ ] MVP 路线：后端提供音频流读取接口（按 sha256+type 定位当前专辑下 audioFiles），前端使用 WebAudio 解码并渲染波形。
- [ ] 后续可选优化（不纳入 MVP）：后端返回“抽样峰值数据”（min/max 或 RMS），前端仅负责绘制与交互，降低解码成本。

## 3. 前端组件实现

- [ ] 引入并封装波形组件（建议使用 wavesurfer.js + Regions/Timeline/Hover 插件，或等价实现）。
- [ ] 与现有 start/end 输入框双向同步：
  - 拖拽选区 -> 更新 start/end
  - 编辑 start/end -> 更新选区
- [ ] 加入边界与校验：start>=0、end>start、end<=duration（如能获取 duration）。
- [ ] 与预览按钮协作：
  - 预览时使用当前选区与音量参数（保持现有 previewSound 行为）。

## 4. 后端支持（如采用路线 A 或 B）

- [ ] 提供按 sha256+type 的音频流接口，并明确错误响应：文件不存在、类型不支持、读取失败。
- [ ] 确保仅访问当前编辑的专辑目录下 audioFiles，避免路径穿越风险。

## 5. 体验打磨

- [ ] 加载态/失败态：显示 skeleton/loading、错误提示与降级（仍可手动输入裁剪）。
- [ ] 缓存：同一音频文件重复打开对话框时，复用波形数据/解码结果（可先做内存级缓存）。

## 6. 验证

- [ ] 手动验收：
  - 选择不同源文件时波形正确切换。
  - 拖拽选区后，start/end 正确更新且可保存。
  - 直接输入 start/end 后，选区正确跟随。
  - 波形不可用时，仍能完成创建/编辑与预览。
- [ ] 运行 `openspec validate add-waveform-sound-trimmer --strict --no-interactive` 通过。
