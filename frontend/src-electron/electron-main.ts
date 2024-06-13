import { app, BrowserWindow } from 'electron';
/*
 * 无边框功能实现时, 需借此做到窗口关闭等的控制。
 * 从 @electron/remote/main 模块中导入 initialize 和 enable 函数。
 * initialize 函数用于初始化@electron/remote，而 enable 函数用于启用特定窗口的远程访问。
 */
import { initialize, enable } from '@electron/remote/main';
import path from 'path';
import os from 'os';

// 初始化 @electron/remote 模块，使其可以在主进程和渲染进程之间进行通信。
initialize();

const appDir = path.dirname(app.getAppPath());
// 这里以后支持多平台时, 需要使用, 并在后方path.join的最后一个参数处, 替换为此name变量。
// > 下方代码可能需要使用quasar专用的环境变量来替代
// const key_tone_sdk_name = 'win32' === process.platform ? 'SiYuan-Kernel.exe' : 'SiYuan-Kernel';
const key_tone_sdk_path = path.join(appDir, 'key_tone_sdk', 'key_tone_sdk.exe');

const databasesDir = path.join(app.getPath('appData'), 'KeyTone', 'Database');
const dbPath = path.join(databasesDir, 'key_tone.db');

// console.log('uuuuuuuuuuuuuuuuuuuuuuuuuuuu=', dbPath);
const logsDir = path.join(app.getPath('home'), '.config', 'KeyTone');
const logsDirPath = path.join(logsDir, 'log.jsonl');

// console.log('uuuuuuuuuuuuuuuuuuuuuuuuuuuu=', logsDirPath);

// 确保路径是可用路径(若发现路径不存在, 则递归创建)
const fs = require('fs');
if (!fs.existsSync(databasesDir)) {
  fs.mkdirSync(databasesDir, { recursive: true });
}
if (!fs.existsSync(logsDir)) {
  fs.mkdirSync(logsDir, { recursive: true });
}

// TIPS: 在最前面进行, 是因为防止启动后因子进程未能及时的成功启动, 而触发强制用户选择ip的模态窗口。
if (process.env.DEBUGGING) {
  // TIPS: win平台下, 想要执行二进制文件, 请注意其只认.exe为扩展名的二进制文件(或其它如.bat .cmd等之类的非二进制文件)。
  // TIPS: 还需要注意的是: 我们的路径因按照src下的来。
  //       > 必须确保子进程的文件路径, 是相对于项目的主要JavaScript文件所在目录的。
  //       > 换句话说，这应该是相对于你运行 electron.的地方的路径。(即我们package按照electron的地方, 也就是最原始的src目录。)
  // WARN: 开发模式(dev)下的路径,与发布构建模式(build)下的路径, 是不通用的。
  //       > 我们需要通过electron提供的api来获取最终打包后的路径。 以保证在发布的软件版本中, 可以正常执行相关路径下的文件与资源。
  // WARN: 在 Electron 的打包过程中，会将应用的所有资源文件，包括 JavaScript、HTML、CSS，还有其他的脚本文件和二进制文件，
  //       一起打包到最终的应用程序中。
  //       > 但是, 使用 cp.spawn 启动的子进程不会被包含在打包的 Electron 应用内。如果你的应用需要在运行时启动其他的可执行
  //         文件或脚本，你需要确保这些文件在每台用户的机器上都能被正确地访问。对于这样的需求，一种常见的做法是在 Electron
  //         应用的安装程序中包含这些额外的可执行文件，并且在安装过程中将它们放置在合适的位置(或通过makefile让其在编译后直
  //         接到这个目录中去)。然后在 Electron 应用中用正确的路径来引用这些文件。(也因此, 注释了dev模式下的这个)
  // const cp = require('child_process');
  // const sdkProcess = cp.spawn('./src-electron/key_tone_sdk/key_tone_sdk.exe', [], {
  //   detached: false,
  //   stdio: 'ignore',
  // });
} else {
  const cp = require('child_process');
  // const sdkProcessParameter = [dbPath, '', logsDirPath];
  // mvp阶段暂时不需要数据库和日志记录
  const sdkProcessParameter = [''];
  const sdkProcess = cp.spawn(key_tone_sdk_path, sdkProcessParameter, {
    detached: false,
    stdio: 'ignore',
  });
}

// needed in case process is undefined under Linux
const platform = process.platform || os.platform();

let mainWindow: BrowserWindow | undefined;

function createWindow() {
  /**
   * Initial window options
   */
  mainWindow = new BrowserWindow({
    icon: path.resolve(__dirname, 'icons/icon.png'), // tray icon
    width: 390,
    height: 500,
    useContentSize: true,
    frame: false,
    resizable: false, // 是否可调整窗口大小, 默认为true
    // autoHideMenuBar: true,  // 此方式只是自动隐藏菜单栏, 但仍可通过 'alt' 键打开。
    show: false,
    webPreferences: {
      sandbox: false, // 能够在预加载脚本中导入@electronic/remote
      contextIsolation: true,
      // More info: https://v2.quasar.dev/quasar-cli-vite/developing-electron-apps/electron-preload-script
      preload: path.resolve(__dirname, process.env.QUASAR_ELECTRON_PRELOAD),
    },
  });

  // 作用：启用指定窗口的远程访问，使得渲染进程可以通过 @electron/remote 模块访问主进程的功能。
  enable(mainWindow.webContents);

  mainWindow.loadURL(process.env.APP_URL);
  mainWindow.on('ready-to-show', () => {
    mainWindow?.show();
  });

  if (process.env.DEBUGGING) {
    // if on DEV or Production with debug enabled
    mainWindow.webContents.openDevTools();
  } else {
    mainWindow.setMenu(null); // 此方式比较彻底(等于彻底放弃了菜单栏, 甚至让开发工具快捷键都失效), 但是此api在macOS下无效。
    mainWindow.setMenuBarVisibility(false); // 此方式是彻底隐藏, 不受alt键影响。(TIPS: 主要用于解决MacOS下无法彻底放弃菜单栏的问题, 两个都开。)
    // we're on production; no access to devtools pls
    mainWindow.webContents.on('devtools-opened', () => {
      mainWindow?.webContents.closeDevTools();
    });
  }

  mainWindow.on('closed', () => {
    mainWindow = undefined;
  });
}

// 实现单利模式 (命令行 传参 自定义/动态 端口的功能, 以后有必要了再说)
const gotTheLock = app.requestSingleInstanceLock();
if (!gotTheLock) {
  app.quit();
} else {
  app.on('second-instance', (event, commandLine, workingDirectory) => {
    // 当运行第二个实例时,将会聚焦到myWindow这个窗口
    if (mainWindow) {
      if (mainWindow.isMinimized()) mainWindow.restore();
      mainWindow.focus();
    }
  });

  // 创建窗口的代码
  // app.on('ready', () => {                      // 用于监听任意的electron事件, 可以触发任意次。// (与whenReady().then()的作用相同, 不过whenReady提供的是异步的方式)
  // app.once('ready', () => {                    // 可以保证监听的事件仅触发一次。              // (在'ready'事件上 与 app.on('ready',callBack) 没有区别。因为此事件本身就只会触发一次。)
  app.whenReady().then(() => {
    // 创建浏览器窗口。
    createWindow();
  });
}

app.on('window-all-closed', () => {
  if (platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (mainWindow === undefined) {
    createWindow();
  }
});
