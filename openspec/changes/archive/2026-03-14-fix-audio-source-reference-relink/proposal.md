# 变更提案：修复裁剪声音在源文件删除后被重导入“自动复联”

## Why
当前 `sounds.*.source_file_for_sound` 通过 `sha256 + name_id + type` 识别源文件，但后端 `name_id` 分配会复用已删除槽位（从 `0` 向上找空位）。
当用户删除某个源文件引用后，再导入同一文件（同 `sha256`）时，`name_id` 可能被复用，导致历史裁剪声音在未重新选择源文件的情况下恢复依赖状态，与预期不符。

## What Changes
- 将同一 `sha256` 下的 `name_id` 生成策略改为 UUID（随机唯一字符串），不再使用可复用的递增编号。
- 保持字段名 `name_id` 与前后端接口不变，仅改变其取值策略（数字字符串 → UUID 字符串），降低改动面。
- 强化播放链路的引用校验：当 `sounds` 或 `key_tone.*` 引用的 `sha256 + name_id + type` 在 `audio_files` 中不存在时，视为缺失引用，不得按 `sha256 + type` 隐式回连。
- 明确兼容策略：旧配置可继续读取（历史数字 `name_id` 仍有效）；新导入项使用 UUID，不要求一次性迁移全量历史数据。

## Impact
- Affected specs: `keytone-album-editor`
- Affected code (planned):
  - `sdk/server/server.go`（导入时 `name_id` UUID 分配策略）
  - `sdk/keySound/keySound.go`（声音与键音播放前的严格引用校验）
  - `sdk/server/server_nameid_test.go`（UUID 分配测试）
  - `sdk/keySound/keySound_alias_validation_test.go`（严格三元校验测试）
  - `frontend/src/components/keytone-album/**`（必要的缺失引用展示与交互一致性回归）
  - `frontend/src/utils/dependencyValidator.ts`（与新语义的一致性校验）
- Risks:
  - `name_id` 值域变化后，需确认所有拼接/比较逻辑仅依赖字符串相等而非数字语义。
  - 需避免“依赖显示缺失但仍可播放”的语义分裂。