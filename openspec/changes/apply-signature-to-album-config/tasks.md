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

### 7. 集成测试准备
- [ ] 手动测试：创建签名 → 应用到专辑 → 检查配置文件
- [ ] 验证加密后的signature字段可正确解密
- [ ] 验证调试日志输出格式正确
- [ ] 测试授权信息不完整时的错误提示
- [ ] 测试图片文件缺失时的降级处理

## 验收标准

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
