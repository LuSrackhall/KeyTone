# Quickstart: 键音专辑签名系统

## Goal
完成一次签名创建→导出签名→导入签名→为专辑签名→导出专辑的端到端流程。

## Manual Steps（推荐先走一遍）
1. 打开键音专辑页面，点击右上方“签名管理”按钮。
2. 在弹出的对话框中点击“创建签名”，填写名称，添加可选介绍与名片图片，点击确认。
3. 在签名列表中选择新签名，点击“导出”，保存生成的 .ktsign 文件。
4. （可选）删除该签名，再点击“导入签名”并选择刚刚导出的 .ktsign，验证成功恢复。
5. 返回专辑导出流程，选择“签名导出”，选中签名并确认（记录允许二次导出设定）。
6. 完成导出，检查专辑文件包含签名元数据。

## Automation（可选的最小集）
- 前端烟雾：在 `frontend/` 运行 `npm run test:e2e:electron:smoke`（Electron 调试环境复用）。
- 报告查看：`npm run test:e2e:electron:report` 打开 HTML 报告；失败时查看截图/视频与 `tests/e2e/.logs/electron-dev.log`。
- 后端契约：为每个 /sdk/* 端点添加/运行最小契约测试（Go testing），覆盖 200/400/409 的典型响应。

### 运行时 baseURL 与端口对齐

- Electron 调试模式：主进程会拉起 SDK 并打印 KEYTONE_PORT；渲染进程通过 `window.myWindowAPI.getBackendPort()` 获取端口，并调用 `UpdateApi(port)`（见 `frontend/src/boot/axios.ts`）对齐 axios baseURL；避免在测试/脚本中硬编码端口。
- 如遇 SSE 已用旧端口建立：触发一次 `window.location.reload()` 以重建连接（代码已在 boot 中处理）。

## Validation Criteria

- 创建后列表中出现签名条目。
- 导出生成的文件扩展名为 .ktsign。
- 导入无需输入保护码即可恢复签名。
- 专辑导出时若选择签名，结果文件包含签名记录字段。
- 多次导出同一签名时，签名记录数组长度增加。
