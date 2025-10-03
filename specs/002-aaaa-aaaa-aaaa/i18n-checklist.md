# i18n 验收 Checklist (T043a)

## 目的
确保签名系统的所有用户可见文本都有完整的中英文翻译支持。

## 验收标准
所有以下类别的文本必须在 `frontend/src/i18n/{zh-CN,en-US}/index.json` 的 `signature` 命名空间中定义。

---

## 1. 对话框标题 (Dialog Titles)

### SignatureManagementDialog
- [x] `signature.title` - "签名管理" / "Signature Management"
- [x] `signature.dialog.createTitle` - "创建新签名" / "Create Signature"
- [x] `signature.dialog.editTitle` - "编辑签名" / "Edit Signature"
- [x] `signature.dialog.importTitle` - "导入签名" / "Import Signature"
- [x] `signature.dialog.exportTitle` - "导出签名" / "Export Signature"
- [x] `signature.dialog.deleteTitle` - "删除签名" / "Delete Signature"
- [x] `signature.dialog.overwriteTitle` - "覆盖确认" / "Overwrite Confirmation"

### SignatureSelectDialog
- [x] `signature.exportFlow.selectSignature` - "选择签名" / "Select Signature"

---

## 2. 按钮文本 (Button Text)

### 主要操作
- [x] `signature.createSignature` - "创建签名" / "Create Signature"
- [x] `signature.importSignature` - "导入签名" / "Import Signature"
- [x] `signature.exportSignature` - "导出签名" / "Export Signature"
- [x] `signature.deleteSignature` - "删除签名" / "Delete Signature"

### 对话框操作
- [x] `signature.dialog.confirm` - "确定" / "OK"
- [x] `signature.dialog.cancel` - "取消" / "Cancel"
- [x] `signature.dialog.overwrite` - "覆盖" / "Overwrite"
- [x] `signature.dialog.selectFile` - "选择文件" / "Select File"

### 复用按钮 (通过 KeyToneAlbum 命名空间)
- [x] `KeyToneAlbum.close` - "关闭" / "Close"
- [x] `KeyToneAlbum.cancel` - "取消" / "Cancel"
- [x] `KeyToneAlbum.confirm` - "确定" / "Confirm"

---

## 3. 表单标签 (Form Labels)

- [x] `signature.signatureName` - "签名名称" / "Signature Name"
- [x] `signature.signatureIntro` - "签名简介" / "Introduction"
- [x] `signature.cardImage` - "名片图片" / "Card Image"
- [x] `signature.selectImage` - "选择图片" / "Select Image"
- [x] `signature.createdAt` - "创建时间" / "Created At"

---

## 4. 错误提示 (Error Messages)

### 创建/导入错误
- [x] `signature.notify.createFailed` - "签名创建失败" / "Create failed"
- [x] `signature.notify.importFailed` - "签名导入失败" / "Import failed"
- [x] `signature.notify.exportFailed` - "签名导出失败" / "Export failed"
- [x] `signature.notify.deleteFailed` - "签名删除失败" / "Delete failed"

### 验证错误
- [x] `signature.notify.nameRequired` - "签名名称不能为空" / "Name required"
- [x] `signature.notify.invalidFormat` - "请选择 .ktsign 格式的签名文件" / "Select .ktsign file"
- [x] `signature.notify.signatureExists` - "签名已存在" / "Signature exists"

### 导出流程错误
- [x] `signature.exportFlow.signatureRequired` - "专辑已有签名，导出时必须选择签名" / "Album is signed, must select signature for export"

---

## 5. 成功提示 (Success Messages)

- [x] `signature.notify.createSuccess` - "签名创建成功" / "Signature created"
- [x] `signature.notify.importSuccess` - "签名导入成功" / "Signature imported"
- [x] `signature.notify.exportSuccess` - "签名导出成功" / "Signature exported"
- [x] `signature.notify.deleteSuccess` - "签名删除成功" / "Signature deleted"

---

## 6. 占位文案 (Placeholder Text)

### 空状态
- [x] `signature.emptyState.noSignatures` - "暂无签名" / "No signatures"
- [x] `signature.emptyState.createFirst` - "请先创建或导入签名" / "Create or import a signature"

### 文件选择
- [x] `signature.dialog.dragOrClick` - "点击选择或拖放 .ktsign 文件" / "Click or drag .ktsign file"

### 导出流程
- [x] `signature.exportFlow.selectPrompt` - "请选择用于导出的签名" / "Select signature for export"

---

## 7. 工具提示 (Tooltips)

通过按钮的 q-tooltip 实现：
- [x] 签名管理按钮：`signature.title`
- [x] 导出按钮：`signature.exportSignature`
- [x] 删除按钮：`signature.deleteSignature`

---

## 8. 确认对话框 (Confirmation Dialogs)

### 覆盖确认
- [x] `signature.dialog.overwriteMessage` - "签名 \"{name}\" 已存在，是否覆盖？" / "Signature \"{name}\" already exists, overwrite?"

### 删除确认
- 通过 Quasar 的 q.dialog 实现，使用通用的确认消息

---

## 9. 文件描述 (File Descriptions)

- [x] `signature.notify.fileDescription` - "KeyTone 签名文件" / "KeyTone Signature File"

---

## 验收结果

### zh-CN 覆盖率
- ✅ 所有标题、按钮、标签: 100%
- ✅ 所有错误提示: 100%
- ✅ 所有成功提示: 100%
- ✅ 所有占位文案: 100%
- ✅ 所有确认对话: 100%

### en-US 覆盖率
- ✅ 所有标题、按钮、标签: 100%
- ✅ 所有错误提示: 100%
- ✅ 所有成功提示: 100%
- ✅ 所有占位文案: 100%
- ✅ 所有确认对话: 100%

---

## 未覆盖语言

按照 plan.md 的约束，以下语言暂不适配（功能优先）：
- de-DE (德语)
- es-ES (西班牙语)
- fr-FR (法语)
- 其他语言

这些语言将回退到英文或显示未翻译的 key。

---

## 检查方法

### 自动检查
```bash
# 检查 zh-CN
grep -o '"signature\.[^"]*"' frontend/src/i18n/zh-CN/index.json | sort | uniq

# 检查 en-US
grep -o '"signature\.[^"]*"' frontend/src/i18n/en-US/index.json | sort | uniq

# 对比两者是否一致
diff <(grep -o '"signature\.[^"]*"' frontend/src/i18n/zh-CN/index.json | sort) \
     <(grep -o '"signature\.[^"]*"' frontend/src/i18n/en-US/index.json | sort)
```

### 手动检查
1. 启动应用
2. 切换语言到中文
3. 打开签名管理对话框，验证所有文本为中文
4. 执行所有操作（创建、导入、导出、删除）
5. 切换语言到英文
6. 重复步骤3-4，验证所有文本为英文

---

## 结论

✅ **i18n 覆盖率验收通过**

所有签名系统的用户可见文本都已在 `signature` 命名空间中完整定义中英文翻译，符合 spec 的 NFR 要求。
