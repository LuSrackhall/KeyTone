# 签名管理功能规格增量：签名名片图片清理安全性

## MODIFIED Requirements

### Requirement: 签名图片路径初始化

Normative: The backend SHALL create and use the signature image directory under `ConfigPath/signature`; `CleanupOrphanCardImages` SHALL delete image files that are not referenced by any signature **only when the system can reliably enumerate all referenced images from the signature configuration**. If any signature entry cannot be decrypted or parsed (e.g., due to encryption key mismatch), `CleanupOrphanCardImages` MUST NOT delete any files and MUST emit a warning log indicating cleanup was skipped for safety.

#### Scenario: 正常清理孤立图片（可可靠解析）

- **GIVEN** 签名配置中存在一个或多个签名条目，且所有条目均可被成功解密并解析
- **AND** `ConfigPath/signature` 目录下存在一些未被任一签名引用的图片文件
- **WHEN** 执行 `CleanupOrphanCardImages`
- **THEN** 系统删除所有“未被引用”的图片文件
- **AND** 系统保留所有“被引用”的图片文件

#### Scenario: 密钥不兼容导致无法解密时跳过清理（防误删）

- **GIVEN** 签名配置中存在一个或多个签名条目
- **AND** 当前运行实例的 KeyA/动态密钥与写入该配置时使用的密钥不一致
- **WHEN** 执行 `CleanupOrphanCardImages`
- **THEN** 系统无法可靠解密/解析至少一个签名条目
- **AND** 系统 MUST 跳过本次删除操作，不删除 `ConfigPath/signature` 中任何文件
- **AND** 系统 MUST 记录 warning 日志，明确说明“因解密/解析失败，为安全起见跳过清理”

#### Scenario: 部分条目解析失败时同样跳过清理（保守策略）

- **GIVEN** 签名配置包含多个签名条目
- **AND** 其中至少一个条目可解密解析成功、至少一个条目解密或解析失败
- **WHEN** 执行 `CleanupOrphanCardImages`
- **THEN** 系统 MUST 视引用集合为不可信
- **AND** 系统 MUST 跳过本次删除操作，不删除 `ConfigPath/signature` 中任何文件
