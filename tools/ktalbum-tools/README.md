# KeyTone Album Tools

KeyTone Album Tools 是一个用于处理 KeyTone 专辑文件（.ktalbum）的工具集。提供命令行和 Web 界面两种使用方式。

## 功能特点

- 解包 .ktalbum 文件到 .zip 格式
- 查看专辑文件信息
- 支持文件拖放
- 提供 Web 界面和命令行两种使用方式
- 跨平台支持（Windows/macOS/Linux）

## 下载和安装

从源码目录选择对应平台的可执行文件：

- Windows:
  - x64: `ktalbum-tools-v0.1.0-windows-amd64.exe`
  - ARM64: `ktalbum-tools-v0.1.0-windows-arm64.exe`
- macOS:
  - Intel: `ktalbum-tools-v0.1.0-darwin-amd64`
  - Apple Silicon: `ktalbum-tools-v0.1.0-darwin-arm64`
- Linux:
  - x64: `ktalbum-tools-v0.1.0-linux-amd64`

## 使用方法

### Web 界面（推荐）

1. 启动 Web 服务：

```bash
# Windows
ktalbum-tools-v0.1.0-windows-amd64.exe web

# macOS/Linux
./ktalbum-tools-v0.1.0-darwin-arm64 web
```

2. 打开浏览器访问 `http://localhost:8080`

3. 通过界面拖放或选择 .ktalbum 文件进行操作

### 命令行

```bash
# 解包文件（显示详细信息）
ktalbum-tools extract -in album.ktalbum -out output.zip -v

# 查看文件信息
ktalbum-tools info -in album.ktalbum

# 启动 Web 服务（指定端口）
ktalbum-tools web -port 8080
```

### 命令行参数说明

#### extract 命令

- `-in`: 输入的 .ktalbum 文件路径（必需）
- `-out`: 输出的 .zip 文件路径（可选，默认使用输入文件名）
- `-v`: 显示详细信息

#### web 命令

- `-port`: Web 服务端口号（可选，默认 8080）

## 开发相关

### 从源码构建

1. 安装依赖：

```bash
# Go 依赖
go mod tidy

# 前端依赖
cd web/frontend
npm install
```

2. 构建：

```bash
# Linux/macOS
chmod +x build.sh
./build.sh

# Windows
build.bat
```

### 目录结构

```
ktalbum-tools/
├── commands/        # 命令实现
├── utils/          # 工具函数
├── web/            # Web 服务
│   ├── frontend/   # 前端代码
│   └── server.go   # 后端服务
├── main.go         # 主程序
└── build.sh        # 构建脚本
```

## 注意事项

1. Web 服务默认只监听 localhost，仅供本地使用
2. 确保有足够的磁盘空间用于临时文件
3. 处理大文件时可能需要较长时间
4. 在 macOS/Linux 上需要给可执行文件添加执行权限：

```bash
chmod +x ktalbum-tools-*
```

## 常见问题

1. 端口被占用：
   - 使用 `-port` 参数指定其他端口
   - 例如：`ktalbum-tools web -port 8081`

2. 文件无法打开：
   - 确保有文件的读取权限
   - 检查文件是否被其他程序占用

3. 解包失败：
   - 确保文件是有效的 .ktalbum 格式
   - 检查文件是否完整（未损坏）

## 更新日志

### v0.1.0

- 初始版本
- 基本的解包功能
- Web 界面支持
- 文件信息查看