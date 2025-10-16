# 签名管理系统 - 优先级行动计划

## 当前状态（2025-10-16）

**完成度**：75%（核心功能 95% + 测试 5% + 文档 0%）

**分支**：`openSpec`

## 立即优先的任务（可以现在开始）

### 任务 A：英文国际化翻译（4 小时）
**文件**：`frontend/src/i18n/en-US/index.json`

**需要添加的部分**：签名管理模块的英文翻译

**验收标准**：
- [ ] 所有中文翻译都有对应的英文翻译
- [ ] 翻译语法正确，风格一致
- [ ] 没有遗漏的翻译键
- [ ] 测试两种语言的完整功能

**预计工时**：4 小时

---

### 任务 B：后端 API 集成测试（8 小时）
**文件**：
- `sdk/signature/encryption_test.go`（已完成 90%）
- 新增：`sdk/server/signature_handlers_test.go`（需要创建）
- 新增：`sdk/config/config_signature_test.go`（可选）

**需要覆盖的 API**：
- [ ] POST /signature/create
- [ ] GET /signature/list
- [ ] PUT /signature/update
- [ ] DELETE /signature/delete/:id
- [ ] GET /signature/export/:id
- [ ] POST /signature/import
- [ ] GET /signature/image/:filename

**预计工时**：8 小时

---

### 任务 C：前端单元测试（10 小时）
**文件**：需要创建测试文件

**需要覆盖的模块**：
- [ ] `frontend/src/boot/query/signature-query.ts` - 签名服务层
- [ ] `frontend/src/components/SignatureFormDialog.vue` - 表单对话框
- [ ] `frontend/src/components/SignatureSelectDialog.vue` - 选择对话框
- [ ] `frontend/src/pages/Signature_management_page.vue` - 主页面

**预计工时**：10 小时

**测试框架**：需要根据项目使用的框架（推荐 Vitest）

---

### 任务 D：专辑文件格式扩展（4 小时）
**文件**：`tools/ktalbum-tools/utils/album.go`

**需要实现**：
- [ ] 在 `Album` 结构体中添加 `Signatures` 字段
- [ ] 更新序列化逻辑
- [ ] 更新反序列化逻辑
- [ ] 添加向后兼容性
- [ ] 编写单元测试

**预计工时**：4 小时

---

### 任务 E：专辑导出流程集成（6 小时）
**文件**：
- 修改：`frontend/src/pages/AlbumExport.vue`（或相关导出页面）
- 使用：`frontend/src/components/SignatureSelectDialog.vue`

**需要实现**：
- [ ] 在导出流程中添加签名选择步骤
- [ ] 将选中的签名嵌入专辑文件
- [ ] 更新导出 API 调用
- [ ] 显示签名状态（已签名/未签名）

**预计工时**：6 小时

---

## 执行顺序建议

### 第 1 阶段（2 天）- 快速获胜
1. **任务 A**：英文 i18n（4 小时）
   - 参考中文翻译，直接添加英文版本
   - 快速、低风险

2. **任务 B**：后端测试（8 小时，分 2 天）
   - 第一天：创建基本测试框架（4 小时）
   - 第二天：完成所有 API 测试（4 小时）

### 第 2 阶段（3 天）- 核心功能补全
3. **任务 D**：专辑格式扩展（4 小时）
   - 后端独立任务，不依赖前端
   - 必须完成才能进行任务 E

4. **任务 E**：专辑导出集成（6 小时）
   - 依赖任务 D 完成
   - 使用已完成的 SignatureSelectDialog 组件

### 第 3 阶段（5 天）- 测试完善
5. **任务 C**：前端单元测试（10 小时）
   - 分散在整个测试阶段

---

## 技术细节

### 任务 A：英文国际化翻译

**复制位置**：参考 `frontend/src/i18n/zh-CN/index.json` 中的 `signature` 部分

**关键翻译术语**：
- 签名 = Signature
- 保护码 = Protect Code
- 名片图片 = Card Image
- 个人介绍 = Personal Introduction
- 导出 = Export
- 导入 = Import
- 创建 = Create
- 编辑 = Edit
- 删除 = Delete

---

### 任务 B：后端测试框架

**推荐使用**：Go 的 `testing` 包 + `testify` 库

**测试结构**：
```go
// 创建测试用的临时目录和配置文件
func setupTestConfig() (string, error) { ... }

// 清理测试环境
func teardownTestConfig() { ... }

// 测试创建签名
func TestCreateSignature(t *testing.T) { ... }

// 测试导入导出流程
func TestImportExportSignature(t *testing.T) { ... }
```

---

### 任务 D：专辑格式扩展

**当前 Album 结构**：查看 `tools/ktalbum-tools/utils/album.go`

**需要添加的字段**：
```go
type Album struct {
    // 现有字段...
    Signatures []SignatureInAlbum `json:"signatures"`
}

type SignatureInAlbum struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Timestamp string `json:"timestamp"` // 签名应用时的时间戳
}
```

**向后兼容性**：
- 旧版本文件中没有此字段时，应默认返回空数组
- 序列化时始终包含此字段（即使为空）

---

### 任务 E：专辑导出集成

**集成点**：找到导出流程的主页面

**集成步骤**：
1. 导入 `SignatureSelectDialog` 组件
2. 在导出步骤序列中添加签名选择步骤
3. 收集用户选择的签名 ID
4. 调用导出 API 时传递 `signatures` 参数
5. 在导出结果中显示签名信息

---

## 验收检查清单

### 任务 A 完成标准
- [ ] en-US 翻译文件包含所有签名相关键
- [ ] 中英文测试均能正常显示文本
- [ ] 没有缺失的翻译键提示

### 任务 B 完成标准
- [ ] 所有 API 端点都有对应的测试
- [ ] 测试覆盖成功和失败场景
- [ ] 测试覆盖率 > 80%
- [ ] 所有测试通过

### 任务 C 完成标准
- [ ] 所有签名相关组件和服务都有单元测试
- [ ] 测试覆盖率 > 80%
- [ ] 所有测试通过

### 任务 D 完成标准
- [ ] 新字段在 JSON 中正确序列化/反序列化
- [ ] 旧版本文件能正确加载
- [ ] 单元测试覆盖新功能

### 任务 E 完成标准
- [ ] 导出页面能显示签名选择对话框
- [ ] 签名数据正确嵌入到导出的专辑文件
- [ ] 导入时能正确读取签名数据

---

## 可选的后续任务（非阻塞发布）

- [ ] E2E 测试（20 小时）
- [ ] 用户文档（8 小时）
- [ ] 开发者文档（6 小时）
- [ ] 性能优化（12 小时）
- [ ] 安全审计（8 小时）
- [ ] CHANGELOG 更新（2 小时）
- [ ] 版本发布（4 小时）

---

## 风险和依赖

### 关键依赖
- [ ] 任务 D 必须在任务 E 之前完成
- [ ] 任务 B 和 C 可以并行进行
- [ ] 任务 A 独立，可随时进行

### 潜在风险
1. **测试框架缺失**：需要确认项目使用的 Go 测试框架
2. **前端测试设置**：需要确认 Vue 组件测试框架
3. **专辑导出页面位置**：需要找到正确的集成点

---

## 时间估算总结

| 任务          | 预计工时    | 优先级 | 状态   |
| ------------- | ----------- | ------ | ------ |
| A - 英文 i18n | 4 小时      | ★★★    | 待开始 |
| B - 后端测试  | 8 小时      | ★★★    | 待开始 |
| C - 前端测试  | 10 小时     | ★★★    | 待开始 |
| D - 专辑格式  | 4 小时      | ★★★    | 待开始 |
| E - 专辑集成  | 6 小时      | ★★★    | 待开始 |
| **合计**      | **32 小时** | -      | -      |

**预计完成时间**：5 个工作日（每天 8 小时）

---

## 立即需要的信息

请补充以下信息，以加快实施：

1. **测试框架确认**：
   - 后端使用什么 Go 测试库？
   - 前端使用什么 Vue 测试框架？

2. **专辑导出页面**：
   - 导出功能在哪个文件中？
   - 导出页面结构是什么？

3. **Album 结构**：
   - 能否确认当前 Album 结构中是否已有其他 Slice 字段作参考？

---

## 后续行动

1. **确认上述信息**后，立即开始任务 A（英文翻译）
2. 并行进行任务 B（后端测试框架）
3. 后端测试完成后，立即开始任务 D（专辑格式）
4. 任务 D 完成后，立即开始任务 E（专辑集成）
5. 在整个过程中进行任务 C（前端测试）

---

**文档最后更新**：2025-10-16

**预期发布日期**（完成以上 5 个任务后）：2025-10-21
