# audioPackageConfig 测试指南

本文档说明如何运行 `audioPackageConfig.go` 的单元测试。

## 测试文件概览

### 文件位置

- **测试文件**: `audioPackageConfig_test.go`
- **被测文件**: `audioPackageConfig.go`

### 测试覆盖范围

测试文件包含 **14 个测试函数**，全面覆盖 `SetValue()`、`GetValue()` 和 `DeleteValue()` 三个核心函数。

## 测试层级详解

### 1. 简单键值对测试 (0-9, a-f)

- `TestSimpleKeysSet_0to9()` - 测试0-9的简单键
- `TestSimpleKeysSet_a_to_f()` - 测试a-f的简单键
- 验证基础的设置和获取功能

### 2. 两层嵌套测试

- `TestNestedKeys_TwoLevels()` - 使用 `.` 分隔符的两层路径
- 测试案例: `0.0`, `1.1`, `4.a`, `5.b`, `6.c` 等
- 涵盖字符串、整数、浮点数、布尔值、数组等多种类型

### 3. 三层嵌套测试

- `TestNestedKeys_ThreeLevels()` - 三层深度的嵌套结构
- 测试案例: `0.0.0`, `3.b.c`, `database.connection.timeout`
- 验证中等复杂度的配置路径

### 4. 四层嵌套测试

- `TestNestedKeys_FourLevels()` - 四层深度的嵌套结构
- 测试案例: `settings.audio.playback.quality`, `system.config.database.host`
- 验证较复杂的配置结构

### 5. 五层嵌套测试 ⭐ 主要测试

- `TestNestedKeys_FiveLevels()` - 五层深度的嵌套结构（约5层左右）
- 测试案例:
  - `0.0.0.0.0`
  - `a.b.c.d.e`
  - `database.connection.pool.settings.timeout`
  - `app.modules.audio.effects.reverb.level`
  - `config.server.api.auth.oauth.token_expiry`
- **重点验证深层嵌套的增/改/查/删功能**

### 6. 值类型覆盖测试

- `TestValueTypes()` - 各种数据类型的支持
- 测试类型: 字符串、整数、浮点数、布尔值、空字符串、负数等
- 验证类型的正确性和完整性

### 7. 删除操作测试（简单键）

- `TestDeleteValue_SimpleKeys()` - 删除简单路径的键
- 测试案例: `del.0`, `del.1`, `del.a`, `del.b`
- 验证删除后键值为nil

### 8. 删除操作测试（嵌套键）

- `TestDeleteValue_NestedKeys()` - 删除嵌套路径的键
- 测试案例: `nested.del.0`, `nested.del.1.sub`, `a.b.c.d.del`
- 验证深层路径的删除功能

### 9. 修改覆盖测试

- `TestUpdateValue_Overwrite()` - 对已存在的键重新赋值
- 测试场景: 字符串→字符串、浮点数→浮点数、类型转换等
- 验证覆盖旧值的正确性

### 10. 边缘情况测试

- `TestEdgeCases_SpecialCharacters()` - 特殊字符和长路径
  - 包含下划线、连字符、混合特殊字符、Unicode字符
  - 很长的键路径（多个段落）
- `TestEdgeCases_ValueContent()` - 特殊值内容
  - 空字符串、仅空格、换行符、制表符、特殊符号
  - JSON格式的字符串值

### 11. 批量操作一致性测试

- `TestBatchOperations_Consistency()` - 多个操作的顺序一致性
- 场景1: 设置→获取→删除→获取的顺序执行
- 场景2: 嵌套键的多次修改和查询
- 验证操作的原子性和顺序性

### 12. 边界值测试

- `TestBoundaryValues()` - 数值的极限范围
- 测试值: 0, -1, 1, 最大int64, 最小int64, 浮点数极限等
- 验证数值边界的正确处理

### 13. 文件持久化测试

- `TestPersistence_FileWrite()` - 验证数据是否正确写入磁盘
- 设置多个键值对
- 重新加载配置并验证数据完整性

### 14. 综合完整序列测试

- `TestComprehensive_FullSequence()` - 完整的A-F序列测试
- 四个测试阶段：
  1. 设置所有简单键 (0-9, a-f)
  2. 验证所有简单键
  3. 设置部分嵌套键
  4. 删除部分键并验证

### 15. JSON文件内容持久化验证 ⭐ 新增

直接验证磁盘上JSON文件中的实际数据内容，确保数据正确写入文件。

**文件验证测试集** (通过读取和解析config.json文件)：

- `TestFileContent_SimpleKeysPersistence()` - 验证简单键的文件持久化（4个测试）
- `TestFileContent_NestedKeysPersistence_TwoLevels()` - 验证两层嵌套键的文件持久化（5个测试）
- `TestFileContent_NestedKeysPersistence_FiveLevels()` - 验证五层深层嵌套键的文件持久化（5个测试）⭐
  - 验证: `database.connection.pool.settings.timeout`
  - 验证: `app.modules.audio.effects.reverb`
  - 验证: `config.server.api.auth.oauth`
  - 验证: `system.storage.cache.redis.host`
  - 验证: `component.ui.theme.dark.opacity`
- `TestFileContent_DeletePersistence()` - 验证删除操作是否正确反映在文件中（3个测试）
- `TestFileContent_UpdatePersistence()` - 验证覆盖更新是否正确反映在文件中（4个测试）
- `TestFileContent_BatchOperationsPersistence()` - 验证批量操作的文件持久化
- `TestFileContent_DataTypeConsistency()` - 验证文件值与内存值的一致性（7个测试）
- `TestFileContent_JSONStructureValidity()` - 验证JSON文件结构的有效性

**工作原理**：
1. 通过 `SetValue()` 设置数据
2. 直接读取 `config.json` 文件（使用 `ioutil.ReadFile`）
3. 解析JSON并通过嵌套路径获取值
4. 验证文件中的值与预期值完全一致
5. 确保JSON结构有效

**验证场景**：
- 简单键: `file_0`, `file_1`, `file_a`
- 两层嵌套: `nested.file.0`, `config.debug`, `app.version`
- 五层嵌套: `database.connection.pool.settings.timeout`, `config.server.api.auth.oauth` 等
- 删除验证: 删除前后的文件内容对比
- 覆盖验证: 修改前后的文件内容对比
- 数据类型: 字符串、浮点数、布尔值、空字符串、零值、负数等

## 运行测试

### 运行所有测试

```bash
cd d:/safe/KeyTone/sdk
go test ./audioPackage/config -v
```

### 运行特定测试

```bash
go test -run TestSimpleKeysSet_0to9 ./audioPackage/config -v
go test -run "TestNestedKeys_Five" ./audioPackage/config -v
go test -run "Delete" ./audioPackage/config -v
```

### 运行测试并输出详细信息

```bash
go test ./audioPackage/config -v
go test ./audioPackage/config -cover
go test ./audioPackage/config -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 测试序列组合说明

### 简单序列 (第1-2层)

```
基础序列: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, a, b, c, d, e, f
```

### 两层组合示例 (第2层)

```
0.0    1.1    2.2    3.3    4.a    5.b    6.c    7.d    8.e    9.f
```

### 三层组合示例 (第3层)

```
0.0.0                          database.connection.timeout
1.1.1                          audio.package.volume
2.a.2                          level.config.debug
```

### 四层组合示例 (第4层)

```
0.0.0.0                        settings.audio.playback.quality
a.b.c.d                        system.config.database.host
```

### 五层组合示例 (第5层) ⭐ 主要关注

```
0.0.0.0.0                      database.connection.pool.settings.timeout
a.b.c.d.e                      app.modules.audio.effects.reverb.level
x.y.z.w.q                      config.server.api.auth.oauth.token_expiry
```

## 关键测试特性

### 1. 支持 `.` 分隔符的路径

所有测试都验证了通过 `.` 来访问嵌套JSON值的功能：

```go
SetValue("database.connection.timeout", 30)      // 创建嵌套结构
value := GetValue("database.connection.timeout") // 获取深层值
DeleteValue("database.connection.timeout")       // 删除深层值
```

### 2. 多种数据类型

- 字符串 (String)
- 整数 (Int - JSON中为浮点数)
- 浮点数 (Float)
- 布尔值 (Boolean)
- 空值和特殊值

### 3. 边缘场景覆盖

- 空字符串、仅空格
- 特殊字符和Unicode
- 数值极限（最大/最小值）
- 很长的键路径

### 4. 操作的一致性

- 设置→获取→删除的完整流程
- 多次修改覆盖
- 文件持久化验证

## 预期结果

所有测试应该通过 (PASS)。示例输出：

```
PASS
ok      KeyTone/audioPackage/config     30.982s
```

如果出现失败 (FAIL)，错误信息会显示：
- 失败的测试函数名
- 具体的键值组合
- 期望值和实际值的对比

## 测试数据持久化

测试使用临时目录来存储配置文件：

```go
tmpDir := setupTestDir(t)      // 创建临时测试目录
defer teardownTestDir(tmpDir)  // 测试结束后清理
```

每个测试都在独立的临时目录中运行，不会影响其他测试或实际配置。

## 故障排查

### 测试失败

如果看到类似的错误：

```
SetValue/GetValue 失败
  键: database.connection.timeout
  期望: 30
  实际: <nil>
```

可能原因：
1. Viper 初始化失败
2. 文件权限问题
3. 嵌套路径解析错误
4. 并发访问导致的竞态条件

### 解决方案

1. 检查临时目录权限
2. 增加 `time.Sleep()` 等待时间
3. 查看日志输出进行诊断
4. 确保Mutex锁正确释放

## 性能考量

- 每个测试用例之间有 50-100ms 的睡眠时间
- 这是为了确保文件监听正确捕捉变更
- 生产环境中可能需要调整这些延迟

## 扩展测试

如需添加更多测试：

1. **添加新的测试函数**

```go
func TestNewFeature(t *testing.T) {
    tmpDir := setupTestDir(t)
    defer teardownTestDir(tmpDir)
    initTestConfig(t, tmpDir)
    // 添加测试逻辑
}
```

2. **遵循命名规则**
- 函数名以 `Test` 开头
- 使用清晰的层级描述：`TestFeatureName_Description()`

3. **使用现有的工具函数**
- `setupTestDir()` - 创建临时目录
- `teardownTestDir()` - 清理目录
- `initTestConfig()` - 初始化配置

## 最佳实践

- ✅ 总是在测试开始时设置、结束时清理
- ✅ 在Set/Delete操作后添加适当的延迟
- ✅ 使用有意义的键名（便于调试）
- ✅ 验证正向和反向操作（Set→Get, Delete→Get nil）
- ✅ 测试边缘情况和特殊值
- ✅ 记录每个测试的目的（description字段）

## 相关函数文档

### SetValue(key string, value any)

- 设置指定键的值
- 自动将值写入JSON配置文件
- 支持 `.` 分隔的嵌套路径

### GetValue(key string) any

- 获取指定键的值
- 如果键不存在返回 nil
- 支持 `.` 分隔的嵌套路径

### DeleteValue(key string)

- 删除指定键的值
- 自动从JSON配置文件中移除
- 支持 `.` 分隔的嵌套路径
- 需要多次写入以确保完全删除（viper的已知行为）

## 测试统计

- **总测试函数数**: 24 (新增10个文件内容验证函数)
- **总测试用例数**: 120+ (新增40+个文件持久化验证用例)
- **涵盖嵌套深度**: 1-5层
- **覆盖数据类型**: 6种以上
- **边缘场景**: 13+
- **文件验证覆盖**: 简单键、两层嵌套、五层深层嵌套、删除验证、覆盖验证、数据类型一致性、JSON有效性

**新增验证范围**：
- ✅ 直接读取config.json文件内容
- ✅ 验证磁盘上的实际数据与预期一致
- ✅ JSON结构的嵌套正确性
- ✅ 文件中的删除操作验证
- ✅ 文件中的覆盖操作验证
- ✅ 文件值与内存值的一致性

---

**最后更新**: 2024年10月31日
**版本**: 2.0 (新增文件内容持久化验证)


## 加密配置调试（新增功能说明）

当专辑在导出流程选择“需要签名”后，目录结构将变为：

- `package.json` 仅保留指示 JSON（`_keytone_encrypted: true` 与 `_keytone_core: "core"` 等元数据）；
- 实际的明文配置会在运行时解密，明文不会直接落盘；
- `core` 文件保存 AES-GCM 的二进制密文（nonce + ciphertext）。

此时：

- Viper 在检测到指示 JSON 后会读取 `core`，解密到临时目录并继续提供对配置的访问；
- 任何对配置的更改会重新加密并原子写回 `core` 文件，同时刷新指示 JSON 更新时间；
- 使用调试工具打印实际配置：

```bash
cd sdk
go run ./audioPackage/cmd/printconfig --path <albumDir>
```

若仅想查看原始密文（不解密）：

```bash
go run ./audioPackage/cmd/printconfig --path <albumDir> --raw
```

注意：本加密方案旨在避免明文直观暴露，属于“防随手窥视”的工程折中，并非强安全设计。

