import { app, BrowserWindow, Tray, Menu, ipcMain, shell } from 'electron';
/*
 * 无边框功能实现时, 需借此做到窗口关闭等的控制。
 * 从 @electron/remote/main 模块中导入 initialize 和 enable 函数。
 * initialize 函数用于初始化@electron/remote，而 enable 函数用于启用特定窗口的远程访问。
 */
import { initialize, enable } from '@electron/remote/main';
import path from 'path';
import os from 'os';
import { StoreGet, StoreSet } from 'boot/query/store-query';

// 初始化 @electron/remote 模块，使其可以在主进程和渲染进程之间进行通信。
initialize();

const appDir = path.dirname(app.getAppPath());
// 这里以后支持多平台时, 需要使用, 并在后方path.join的最后一个参数处, 替换为此name变量。
// > 下方代码可能需要使用quasar专用的环境变量来替代
// const key_tone_sdk_name = 'win32' === process.platform ? 'SiYuan-Kernel.exe' : 'SiYuan-Kernel';
const key_tone_sdk_path = path.join(appDir, 'key_tone_sdk', 'key_tone_sdk.exe');

const databasesDir = path.join(
  app.getPath('appData'),
  'KeyToneGoSdk', // 为了和electron原生与前端持久化区域做区分, 我们sdk依赖将使用独立的持久化路径
  'Database'
);
const dbPath = path.join(databasesDir, 'key_tone.db');

// console.log('uuuuuuuuuuuuuuuuuuuuuuuuuuuu=', dbPath);
const logsDir = path.join(app.getPath('home'), '.config', 'KeyToneGoSdk');
const logsDirPath = path.join(logsDir, 'KeyToneSdkLog.jsonl');

// console.log('uuuuuuuuuuuuuuuuuuuuuuuuuuuu=', logsDirPath);

// 确保路径是可用路径(若发现路径不存在, 则递归创建)
// const fs = require('fs');
import fs from 'fs';
if (!fs.existsSync(databasesDir)) {
  fs.mkdirSync(databasesDir, { recursive: true });
}
if (!fs.existsSync(logsDir)) {
  fs.mkdirSync(logsDir, { recursive: true });
}

import cp from 'child_process';
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
  // const sdkProcess = cp.spawn(
  //   './src-electron/key_tone_sdk/key_tone_sdk.exe',
  //   ['-configPath=../sdk'],
  //   {
  //     detached: false,
  //     stdio: 'ignore',
  //   }
  // );
} else {
  // const cp = require('child_process');
  // const sdkProcessParameter = [dbPath, '', logsDirPath];
  // mvp阶段暂时不需要数据库和日志记录
  const sdkProcessParameter = ['-configPath=' + databasesDir];
  const sdkProcess = cp.spawn(key_tone_sdk_path, sdkProcessParameter, {
    detached: false,
    stdio: 'ignore',
  });
}

// needed in case process is undefined under Linux
const platform = process.platform || os.platform();

let mainWindow: BrowserWindow | undefined;
let tray;

function createWindow() {
  /**
   * Initial window options
   */
  mainWindow = new BrowserWindow({
    icon: path.resolve(__dirname, 'icons/icon.png'), // tray icon
    width: 390,
    height: 500,
    useContentSize: true,
    frame: false, // 设置无边框窗口
    resizable: false, // 是否可调整窗口大小, 默认为true
    transparent: true, // 设置透明窗口, 为进一步的毛玻璃窗口做准备 // 由于纯CSS不支持直接透到操作系统桌面的毛玻璃效果, 因此放弃
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

  // 这三行配合上面的`show:false`, 可以使得前端页面加载完毕后再显示窗口。
  // 作用是: 牺牲的窗口启动速度, 来避免窗口启动后白屏加载前端的糟糕体验。
  mainWindow.on('ready-to-show', () => {
    // 注释掉, 包含在后续判断的else中。
    // mainWindow?.show();

    // 这个"可能"用于配合官方库, 来实现自动隐藏启动的功能。(未验证)
    // if (app.getLoginItemSettings().wasOpenedAtLogin) {
    //   mainWindow?.hide();
    // } else {
    //   mainWindow?.show();
    // }

    // 这个用于配合改用第三方库node-auto-launch, 来实现自动隐藏启动的功能。
    if (process.argv[1] == '--hidden') {
      mainWindow?.hide();
    } else {
      mainWindow?.show();
    }
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

  // mainWindow.on('closed', () => {
  //   mainWindow = undefined;
  // });

  // 关闭窗口时隐藏窗口而不是退出应用
  mainWindow.on('close', (event) => {
    if (!(app as any).isQuiting) {
      event.preventDefault();
      mainWindow?.hide();
    }
    return false;
  });
}

function createTray() {
  // 创建托盘图标
  tray = new Tray(path.join(__dirname, 'icons/icon.png')); // 替换为你的图标路径

  // 创建托盘图标的上下文菜单
  const contextMenu = Menu.buildFromTemplate([
    {
      label: '显示',
      click: () => {
        mainWindow?.show();
      },
    },
    {
      label: '退出',
      click: () => {
        (app as any).isQuiting = true;
        app.quit();
      },
    },
  ]);

  // 设置托盘图标的上下文菜单
  tray.setContextMenu(contextMenu);

  // 设置托盘图标的提示文本
  tray.setToolTip('KeyTone');

  // 点击托盘图标时显示窗口
  // tray.on('click', () => {
  //   mainWindow?.show();
  // });

  // 双击托盘图标时显示窗口：
  tray.on('double-click', () => {
    mainWindow?.show();
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
      if (mainWindow.isMinimized()) {
        mainWindow.restore();
        mainWindow.focus();
      }
      // 如果窗口不再前台(或没有展示或隐藏), 则显示窗口。
      if (!mainWindow.isVisible()) {
        mainWindow.show();
      }
    }
  });

  // 创建窗口的代码
  // app.on('ready', () => {                      // 用于监听任意的electron事件, 可以触发任意次。// (与whenReady().then()的作用相同, 不过whenReady提供的是异步的方式)
  // app.once('ready', () => {                    // 可以保证监听的事件仅触发一次。              // (在'ready'事件上 与 app.on('ready',callBack) 没有区别。因为此事件本身就只会触发一次。)
  app.whenReady().then(() => {
    // 创建浏览器窗口。
    createWindow();
    // 当Electron完成初始化并准备创建浏览器窗口时调用此方法
    createTray();
  });
}

//#region    -----<<<<<<<<<<<<<<<<<<<< -- 开机自启动 start ^_^-_-^_^
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// // * 官方提供的开机自动启动的api(由于不支持全平台, 且在macos上, 存在与第三方库相同的问题, 因此本项目弃用此api)
// app.setLoginItemSettings({
//   openAtLogin: true,
//   // openAsHidden: true, // 功能为以隐藏模式打开。(仅在windows中可用, 因为对于macos不支持MAS和从macos13起及更高的版本<基本等于不可用>故弃用此选项)
// });

// * 改用第三方库[node-auto-launch](https://github.com/Teamwork/node-auto-launch)来实现开机自启动的功能。
// const AutoLaunch = require('auto-launch');
import AutoLaunch from 'auto-launch';

// 由于首次启动时, sdk的web服务可能还未能准备好, 因此通过递归的方式, 来保证获取到正确的配置值。
// * 这只是临时取巧方式, 通用解决方案应该是明确架构出整个生命周期, 来保证整个项目的整体执行顺序。
// * * 比如, 可以通过整体的轮询监听, 或是通过sse, 来保证初始化electron时, sdk是已经启动好的, 从而明确生命周期(不过这样可能会对界面的启动速度有影响<主要指客户端ui的创建速度>)。
// * 当然, 也可以通过sse, 仅对对应需要明确顺序逻辑的依赖部分, 进行执行顺序的明确。 (避免影响electron的应用启动速度<主要指客户端ui的创建速度>)
// * 也可以退一步, 加个延时或直接上轮询, 以避免快速的递归造成短期的cpu消耗过度。
function autoRunInit() {
  StoreGet('auto_startup').then((value) => {
    if (value === false) {
      console.log('value的值是false', value);
      autoRunInit();
    } else {
      // 从此可判断, 我们获取到的value, 已经解析过JSON了, 直接是最终对象。
      console.log('解析value前', typeof value.is_auto_run);

      // 创建一个 AutoLaunch 实例
      const autoLauncher = new AutoLaunch({
        name: 'KeyTone',
        // path: app.getPath('exe'), // 此库的官网上说:对于 NW.js 和 Electron 应用程序，您不必指定路径。我们根据 process.execPath 进行猜测。
        isHidden: value.is_hide_windows,
      });

      // 启动时检查并设置自动启动
      autoLauncher
        .isEnabled()
        .then((isEnabled: any) => {
          // 如果应用程序未设置在开机时自启动, 则主动设置, 若已设置, 则跳过。 此判断仅为防止重复开启。
          if (
            (!isEnabled && value.is_auto_run) ||
            (isEnabled && value.is_hide_windows !== value.is_hide_windows_old)
          ) {
            // if (value.is_auto_run) { // 我们必须避免重复开启, 虽然这样可以低成本实现value.is_hide_windows的实时响应, 但每次启动都去触碰敏感操作是不明智的。

            autoLauncher.enable().then(() => {
              // 如果窗口是否隐藏改变了, 则需要更新is_hide_windows_old以记录最新情况。 (由于我们为了防止重复设置, 只有在关闭自启动, 并开启自启动时, 才会触发重新设置自启动)
              if (value.is_hide_windows !== value.is_hide_windows_old) {
                StoreSet(
                  'auto_startup.is_hide_windows_old',
                  value.is_hide_windows
                );
              }
            }); // 开启自启动
          }

          // // 如果应用已设置在开机时自启动, 则主动设置关闭自启动。 此判断仅为防止防止重复关闭。
          if (isEnabled && !value.is_auto_run) {
            // if (!value.is_auto_run) {
            // 我们必须避免重复开启, 虽然这样可以低成本实现value.is_hide_windows的实时响应, 但每次启动都去触碰敏感操作是不明智的。
            autoLauncher.disable(); // 关闭自启动
          }

          // 此部分因重复, 故已集成至首个判断语句
          // if (
          //   isEnabled &&
          //   value.is_hide_windows !== value.is_hide_windows_old
          // ) {
          //   StoreSet('auto_startup.is_hide_windows_old', value.is_hide_windows); // 更新is_hide_windows_old以记录最新情况。
          //   autoLauncher.disable().then(() => {
          //     autoLauncher.enable(); // 开启自启动  <要确确保在then的disable执行之后>
          //   }); // 关闭自启动

          //   // autoLauncher.enable(); // 开启自启动(不用先关闭再开启, 直接重新调用即可完成is_hide_windows的设置)
          // }
        })
        .catch((err: any) => {
          console.error('Error checking auto-launch status:', err);
        });
    }
  });
}

// 后续可轮询此函数, 以避免重启后才能生效的问题。<似乎不起作用, 只能提示用户重启客户端后生效了, 或是直接在客户端帮助用户强制重启客户端>
// 或者, 也可通过sse来触发其重新调用<更为复杂些, 甚至要监听配置文件的变更情况>。<似乎不起作用-<未验证>, 只能提示用户重启客户端后生效了, 或是直接在客户端帮助用户强制重启客户端>
autoRunInit();
// 请放弃不重启就成功的重新设置自启动的不成熟想法, 不然只会再浪费3天的时间。
// setInterval(() => { },500)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//#endregion ----->>>>>>>>>>>>>>>>>>>> -- 开机自启动 end   -_-^_^-_- ^_^-_-^_^-_-
// ...
// ...
// ...
//!endregion ----->>>>>>>>>>>>>>>>>>>> -- 开机自启动 end   -_-^_^-_- ^_^-_-^_^-_-

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

// 处理 IPC 事件
ipcMain.on('open-external', (event, url) => {
  shell.openExternal(url);
});
