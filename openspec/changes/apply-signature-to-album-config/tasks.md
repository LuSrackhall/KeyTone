# 任务清单

## 实现顺序

### 1. 新增专辑签名加密模块
- [x] 创建`sdk/signature/album.go`文件
  - [x] 定义常量`KeyToneAlbumSignatureEncryptionKey`（32字节固定密钥）
  - [x] 实现`EncryptAlbumSignatureField(signatureJSON string) (string, error)`
  - [x] 实现`DecryptAlbumSignatureField(encryptedHex string) (string, error)`
  - [x] 实现`GenerateQualificationCode(signatureID string) string`（SHA256哈希）
  - [x] 添加详细的函数注释和用途说明

### 2. 新增签名配置应用模块
- [x] 创建`sdk/audioPackage/config/signatureConfig.go`文件
  - [x] 定义专辑签名数据结构`AlbumSignatureEntry`（包含name、intro、cardImagePath、authorization）
  - [x] 定义授权元数据结构`AuthorizationMetadata`（包含requireAuthorization、contactEmail、contactAdditional、authorizedList）
  - [x] 实现核心函数`ApplySignatureToAlbum`
    - [x] 参数验证（albumPath存在性、signatureID有效性、授权信息完整性）
    - [x] 从签名管理配置读取加密签名数据
    - [x] 解密签名ID获取原始UUID
    - [x] 使用动态密钥解密签名Value
    - [x] 生成资格码（SHA256哈希）
    - [x] 处理签名名片图片复制
    - [x] 构建专辑签名对象（包含authorization字段判断）
    - [x] 序列化为JSON并加密
    - [x] 写入专辑配置`signature`字段
    - [x] 添加详细的中文注释说明每个步骤

### 3. 实现图片资源处理
- [x] 在`signatureConfig.go`中添加图片处理函数
  - [x] 实现`copySignatureCardImageToAlbum`函数
    - [x] 从签名配置路径读取原始图片文件
    - [x] 生成新文件名：SHA1(资格码 + 原始文件名 + 时间戳) + 扩展名
    - [x] 复制到`{albumPath}/audioFiles/`目录
    - [x] 返回相对路径`audioFiles/{newFilename}`
  - [x] 错误处理：图片文件不存在时跳过，返回空字符串
  - [x] 添加日志记录图片复制操作

### 4. 完善API端点实现
- [x] 修改`sdk/server/server.go`中的`/keytone_pkg/apply_signature_config`路由
  - [x] 解析请求体：albumPath、signatureId、requireAuthorization、contactEmail、contactAdditional
  - [x] 参数校验：必需字段存在性、授权信息完整性
  - [x] 调用`ApplySignatureToAlbum`函数
  - [x] 返回成功响应：包含qualificationCode
  - [x] 错误处理：返回具体错误信息
  - [x] 添加详细的中文注释

### 5. 添加调试日志输出
- [x] 在`ApplySignatureToAlbum`函数末尾添加调试日志
  - [x] 打印标题：`[专辑签名调试] 签名已成功应用到专辑配置 - 未加密内容：`
  - [x] 格式化输出未加密的签名JSON（使用`json.MarshalIndent`）
  - [x] 包含资格码、签名名称、授权状态等关键信息
  - [x] 使用`fmt.Printf`输出到终端（方便开发者观察）

### 6. 代码审查与优化
- [x] 检查所有新增代码的注释完整性（中文注释）
- [x] 确保错误处理覆盖所有边界情况
- [x] 验证授权字段逻辑正确性
  - [x] 原始作者签名包含authorization对象
  - [x] 非原始作者签名不包含authorization
  - [x] authorizedList初始化为空数组
- [x] 确认图片复制逻辑与现有音频文件管理一致
- [x] 检查资格码生成的唯一性和确定性

### 7. 实现directExportAuthor字段
- [x] 在AuthorizationMetadata结构中添加DirectExportAuthor字段
- [x] 首次导出时设置directExportAuthor为当前导出者资格码
- [x] 再次导出时更新原始作者签名的directExportAuthor
- [x] 非原始作者签名不包含authorization字段
- [x] 添加相关日志记录

### 8. 实现三种导出情况处理
- [x] 删除"无需签名+需要授权"分支逻辑（在规格中明确）
- [x] 情况1：首次导出无需签名 - 跳过签名应用API
- [x] 情况2：首次导出需要签名且需要授权 - 创建authorization（requireAuthorization=true）
- [x] 情况3：首次导出需要签名但无需授权 - 创建authorization（requireAuthorization=false）
- [x] 再次导出时识别原始作者签名并更新directExportAuthor
- [x] 再次导出时添加贡献者签名（无authorization字段）

### 10. 再次导出流程优化 (Frontend)
- [x] 实现再次导出时的提示对话框
  - [x] 检测到专辑已有签名且用户选择"需要签名"时触发
  - [x] 显示提示内容："该键音专辑原始作者明确了该键音包的二次导出必须实施签名..."
- [x] 实现授权导入流程
  - [x] 当requireAuthorization=true时，强制先显示授权导入对话框
  - [x] 验证通过后才允许进入签名选择
- [x] 优化签名选择逻辑
  - [x] 允许选择已存在于专辑中的签名
  - [x] 当选择已存在签名时，弹出二次确认对话框："是否确认更新签名？"
  - [x] 根据用户选择决定是否更新签名内容（name, intro, image）
  - [x] 确保始终更新directExportAuthor
- [x] 适配"无需签名"选项
  - [x] 若专辑无签名，再次导出时允许用户选择是否签名（同首次导出）
  - [x] 若专辑有签名，强制进入签名流程（不可选"无需签名"）
- [x] 修复签名选择列表状态刷新问题
  - [x] 每次打开对话框时强制重新加载签名列表及状态
  - [x] 确保"原始作者"和"贡献作者"标签实时更新

### 11. 集成测试
- [ ] 验证首次导出三种情况
- [ ] 验证再次导出三种情况
- [ ] 验证签名更新逻辑（更新内容 vs 不更新内容）
- [ ] 验证授权限制逻辑
- [x] 移除测试环境配置对话框及相关逻辑

### 10. 前端实现（需求1.2.3）
- [x] 10.1 在types/export-flow.ts中添加类型定义
  - [x] SignatureAuthorInfo - 签名作者信息
  - [x] AlbumSignatureEntry - 专辑签名条目
  - [x] AlbumSignatureInfo - 专辑签名信息
  - [x] AvailableSignature - 可用签名信息
- [x] 10.2 在keytonePkg-query.ts中添加API调用函数
  - [x] GetAlbumSignatureInfo - 获取专辑签名信息
  - [x] CheckSignatureInAlbum - 检查签名是否在专辑中
  - [x] CheckSignatureAuthorization - 检查签名授权状态
  - [x] GetAvailableSignatures - 获取可用签名列表
- [x] 10.3 创建SignatureAuthorsDialog.vue组件（需求4）
  - [x] 位置：frontend/src/components/export-flow/SignatureAuthorsDialog.vue
  - [x] 展示原始作者信息
  - [x] 展示直接导出作者信息
  - [x] 展示历史贡献作者列表
  - [x] 处理无签名情况
  - [x] 加载状态和错误处理
- [x] 10.4 创建SignatureSelectionDialog.vue组件（需求3）
  - [x] 位置：frontend/src/components/export-flow/SignatureSelectionDialog.vue
  - [x] 使用GetAvailableSignatures获取签名列表
  - [x] 标记已在专辑中的签名（蓝色边框）
  - [x] 使能/失能签名选项（未授权签名置灰）
  - [x] 显示签名授权状态徽章
  - [x] 筛选功能（仅显示已授权/已在专辑中）
  - [x] 未授权签名悬停提示

### 11. 前端页面集成
- [x] 11.1 在专辑页面添加"查看签名信息"按钮
  - [x] 添加badge图标按钮
  - [x] 导入SignatureAuthorsDialog组件
  - [x] 添加signatureAuthorsDialogRef引用
  - [x] 实现showAlbumSignatureInfo方法
  - [x] 检查专辑是否选中
- [x] 11.2 集成SignatureAuthorsDialog到专辑页面
  - [x] 传递albumPath属性
  - [x] 通过ref调用open方法
- [x] 11.3 更新useExportSignatureFlow
  - [x] 导入GetAlbumSignatureInfo API
  - [x] 导入AlbumSignatureInfo类型
  - [x] 为后续集成做准备

### 12. Bug修复
- [x] 12.1 修复"无需签名"时仍进入授权对话框的问题
  - [x] 问题：选择"无需签名"后仍然显示授权要求对话框
  - [x] 原因：handleConfirmSignatureSubmit中未检查needSignature标志
  - [x] 修复：在handleConfirmSignatureSubmit中添加条件判断，needSignature=false时直接完成
  - [x] 文件：useExportSignatureFlow.ts
- [x] 12.2 修复SignatureAuthorsDialog尺寸过大导致溢出
  - [x] 问题：对话框宽度600-800px在固定尺寸应用中溢出
  - [x] 修复：调整对话框尺寸为90vw，最大480px，最大高度85vh
  - [x] 优化：调整图片尺寸从100px改为70px
  - [x] 优化：调整字体从text-h6改为text-subtitle2
  - [x] 优化：减少内边距和间距，添加滚动支持
  - [x] 文件：SignatureAuthorsDialog.vue

### 13. 再次导出流程集成（三种情况自动识别）
- [x] 13.1 更新useExportSignatureFlow.ts
  - [x] 修改ExportSignatureFlowOptions接口，添加albumPath参数
  - [x] 更新start方法，使用GetAlbumSignatureInfo获取真实签名状态
  - [x] 实现三种情况自动识别：
    - [x] 情况1：专辑无签名 → 进入"确认签名"对话框
    - [x] 情况2：专辑有签名且需要授权 → 进入"授权门控"对话框
    - [x] 情况3：专辑有签名但不需要授权 → 直接进入"签名选择"对话框
  - [x] 保留旧参数以向后兼容测试代码
  - [x] 添加错误处理，失败时默认按首次导出处理
- [x] 13.2 更新Keytone_album_page_new.vue
  - [x] 修改exportAlbum方法，传递albumPath参数
  - [x] 添加专辑路径验证
  - [x] 保留测试参数以兼容测试对话框
  - [x] 更新注释说明三种情况的自动识别逻辑
- [x] 13.3 更新SignaturePickerDialog.vue
  - [x] 添加albumPath属性
  - [x] 页面调用时传递albumPath
  - [x] 为后续集成GetAvailableSignatures做准备

### 14. 组件文件路径重组
- [x] 14.1 将SignatureAuthorsDialog.vue移动到export-flow目录
  - [x] 从：frontend/src/components/SignatureAuthorsDialog.vue
  - [x] 到：frontend/src/components/export-flow/SignatureAuthorsDialog.vue
  - [x] 原因：该组件专用于专辑导出流程，应与其他导出流程组件放在一起
- [x] 14.2 将SignatureSelectionDialog.vue移动到export-flow目录
  - [x] 从：frontend/src/components/SignatureSelectionDialog.vue
  - [x] 到：frontend/src/components/export-flow/SignatureSelectionDialog.vue
  - [x] 原因：该组件专用于导出时的签名选择，属于导出流程的一部分
- [x] 14.3 更新所有导入路径
  - [x] Keytone_album_page_new.vue：更新导入语句
  - [x] 规格文档：更新所有路径引用

### 15. 集成测试准备
- [ ] 手动测试情况1：无需签名导出
- [ ] 手动测试情况2：首次导出需要签名且需要授权
- [ ] 手动测试情况3：首次导出需要签名但无需授权
- [ ] 手动测试再次导出：验证directExportAuthor更新
- [ ] 手动测试再次导出：验证贡献者签名添加
- [ ] 验证加密后的signature字段可正确解密
- [ ] 验证调试日志输出格式正确
- [ ] 测试授权信息不完整时的错误提示
- [ ] 测试图片文件缺失时的降级处理
- [ ] 测试新增的4个API端点功能
- [ ] 测试"查看签名信息"按钮功能
- [ ] 测试SignatureAuthorsDialog展示效果

## 验收标准

### 12. Bug修复与优化
- [x] 修复：再次导出时更新原始作者签名导致Authorization字段丢失
  - [x] 修改`ApplySignatureToAlbum`函数，增加`updateSignatureContent`参数
  - [x] 实现保留原有Authorization字段的逻辑
  - [x] 实现根据`updateSignatureContent`决定是否更新基本信息
- [x] 修复：选择"不更新"签名时仍然更新了内容
  - [x] 前端传递`updateSignatureContent`标志
  - [x] 后端根据标志跳过内容更新
- [x] 修复：选择"不更新"签名时仍然复制了图片文件
  - [x] 调整`ApplySignatureToAlbum`逻辑，先检查现有签名
  - [x] 若`!updateSignatureContent`且签名存在，跳过图片复制
- [x] 修复：更新签名时未删除旧图片文件
  - [x] 在复制新图片成功后，检查并删除旧图片
  - [x] 避免垃圾文件堆积
- [x] 优化：确保DirectExportAuthor始终更新
  - [x] 无论是否更新内容，始终更新原始作者签名的DirectExportAuthor
- [x] 优化：智能更新检测
  - [x] 修改`CheckSignatureInAlbum`，增加`hasChanges`返回值
  - [x] 后端实现字段比对（Name, Intro）和图片SHA256比对
  - [x] 前端仅在`hasChanges=true`时显示更新确认对话框

### 功能完整性
- [x] 可成功将签名数据写入专辑配置
- [x] 资格码生成正确（SHA256哈希）
- [x] 签名内容使用专用密钥加密
- [x] 图片文件正确复制到专辑目录
- [x] 授权字段根据用户选择正确包含或排除

### 代码质量
- [x] 所有函数包含详细的中文注释
- [x] 错误处理完善，包含有意义的错误信息
- [x] 遵循项目现有代码风格（Go命名约定、模块组织）
- [x] 无重复代码，复用现有加密函数

### 可观测性
- [x] 终端调试日志包含未加密签名内容
- [x] 日志信息清晰，便于开发者排查问题
- [x] API响应包含操作结果（qualificationCode）

## 依赖关系

- 任务1和任务2可并行开始
- 任务3依赖任务2（需要signatureConfig.go文件存在）
- 任务4依赖任务1、2、3完成（调用所有新增函数）
- 任务5依赖任务4完成（在API实现中添加日志）
- 任务6、7依赖所有实现任务完成
