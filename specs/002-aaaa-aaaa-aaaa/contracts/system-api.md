# Contract: System/Health & Discovery

## Health Check

GET /ping
Response: 200 { message: "ok" } or plain 200

## Discovery (Runtime Port)

- Electron 主进程捕获 SDK 启动日志中的 `KEYTONE_PORT=<port>`（stdout），并通过 `ipcMain` 存储在内存变量。
- 渲染进程通过 `window.myWindowAPI.getBackendPort()` 同步获取端口。
- 前端 axios 通过 `UpdateApi(port)` 对齐 baseURL，避免硬编码。

Notes:
- /ping 仅用于健康检查；端口发现不通过 HTTP 暴露，避免不必要的耦合与安全面。
