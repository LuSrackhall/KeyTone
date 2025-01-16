/**
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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

// 未解决但于本项目已无影响FIXME: 只要引入 vue-i18n , 并使用它,  就会造成调试对话框无法独立打开。(猜测可能是影响了`process.env.DEBUGGING`的正常获取<如果真的是这样, 那就太严重了。(测试结果并没有影响这个环境变量的获取, 可以确定只是影响了`mainWindow.webContents.openDevTools({ mode: 'detach' });`这个api的功能效果)>)
// import { createI18n } from 'vue-i18n';  // 直接复用前端boot中的i18n文件的导出即可, 没必要重复写代码
// import messages from 'src/i18n';
// const i18n = createI18n({
//   locale: 'zh-CN',
//   legacy: false,
//   messages,
// });
// console.log('fffffffffffffffffffffff', i18n.global.t('setting.setting'));

// const i = require('boot/i18n');  //node.js对ts的支持有点恶心, 因此不推荐使用这种导入方式
// const i18n = i.i18n
// setInterval(() => {
// console.log('ffffffffffffffffffffffff', i.i18n.global.t('setting.setting'));
// console.log('ffffffffffffffffffffffff', process.env.DEBUGGING);
// }, 1000);

import { i18n } from 'src/boot/i18n'; //node.js对ts的支持有点恶心, 所以推荐使用这种导入方式, 复用es, 使ts可以更好的支持
// setInterval(() => {
//   console.log('ffffffffffffffffffffffff', i18n.global.t('setting.setting'));
// }, 1000);

// 初始化 @electron/remote 模块，使其可以在主进程和渲染进程之间进行通信。
initialize();

const appDir = path.dirname(app.getAppPath());
// 这里以后支持多平台时, 需要使用, 并在后方path.join的最后一个参数处, 替换为此name变量。
// > TODO: 下方代码可能需要使用quasar专用的环境变量来替代。(正式支持Linux或MacOS之前; 或之后<因为是否有.exe的后缀名,不影响其在Linux或MacOS平台上运行, 即可有可无; .exe后缀名, 仅在win下是必须有的。>)
// >        TIPS: 只需适配此位置, 以及修改MakeFile中对应平台的最终生成名为没有.exe即可。
// const key_tone_sdk_name = 'win32' === process.platform ? 'KeyTone.exe' : 'Keytone';
const key_tone_sdk_path = path.join(appDir, 'key_tone_sdk', 'KeyTone.exe');

const appGetPath = app.getPath('appData');

const configDir = path.join(
  appGetPath,
  'KeyToneConfig', // 为了和electron原生与前端持久化区域做区分, 我们sdk依赖将使用独立的持久化路径
  'Config'
);
const dbPath = path.join(configDir, 'key_tone.db');

const audioPackageDir = path.join(appGetPath, 'KeyToneConfig', 'AudioPackage');

// console.log('uuuuuuuuuuuuuuuuuuuuuuuuuuuu=', dbPath);
const logsDir = path.join(app.getPath('home'), '.config', 'KeyToneGoSdk');
const logsDirPath = path.join(logsDir, 'KeyToneSdkLog.jsonl');

// console.log('uuuuuuuuuuuuuuuuuuuuuuuuuuuu=', logsDirPath);

// 确保路径是可用路径(若发现路径不存在, 则递归创建)
// const fs = require('fs');
import fs from 'fs';
if (!fs.existsSync(configDir)) {
  fs.mkdirSync(configDir, { recursive: true });
}
if (!fs.existsSync(audioPackageDir)) {
  fs.mkdirSync(audioPackageDir, { recursive: true });
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
  const sdkProcessParameter = [
    '-configPath=' + configDir,
    '-audioPackagePath=' + audioPackageDir,
    '-logPathAndName=' + logsDirPath,
  ];
  // inherit: 子进程将继承父进程的标准输入、输出和错误流。这意味着子进程的输出会直接显示在父进程的终端中。这种方式不会对子进程的执行产生其他影响，只是改变了输出的显示方式。
  // pipe: 子进程的标准输入、输出和错误流会被重定向到父进程中。你可以通过监听这些流来捕获子进程的输出。这种方式也不会对子进程的执行产生其他影响，但需要你在父进程中处理这些流
  // ignore: 忽略子进程的标准输入。
  const sdkProcess = cp.spawn(key_tone_sdk_path, sdkProcessParameter, {
    detached: false,
    stdio: ['pipe', 'pipe', 'pipe'],
  });
  // 监听子进程的 stdout
  sdkProcess.stdout.on('data', (data) => {
    // console.log(`[SDK] stdout: ${data}`); // 如果输出内容设计多行文本, 则这种简单的方式只能在第一行前添加前缀
    const lines = data.toString().split('\n');
    lines.forEach((line: any) => {
      if (line.trim()) {
        process.stdout.write(`[SDK] ${line}\n`);
      }
    });
  });

  // 监听子进程的 stderr
  sdkProcess.stderr.on('data', (data) => {
    // console.error(`[SDK] stderr: ${data}`); // 如果输出内容设计多行文本, 则这种简单的方式只能在第一行前添加前缀
    const lines = data.toString().split('\n');
    lines.forEach((line: any) => {
      if (line.trim()) {
        process.stderr.write(`[SDK] ${line}\n`);
      }
    });
  });

  // 监听子进程的关闭事件
  sdkProcess.on('close', (code) => {
    console.log(`[SDK] exited with code ${code}`);
  });
}

// needed in case process is undefined under Linux
const platform = process.platform || os.platform();

let mainWindow: BrowserWindow | undefined;
let tray: Tray;

// let sdkIsRun: boolean = false; // eslint有类型推断, 主动设置会报错, 但我又懒得关闭此推断。
let sdkIsRun = false;

const iconPath = process.env.DEBUGGING
  ? path.join(process.cwd(), 'src-electron', 'icons', 'icon.png') // 开发环境路径
  : path.join(__dirname, 'icons', 'icon.png'); // 生产环境路径

function createWindow() {
  /**
   * Initial window options
   */
  mainWindow = new BrowserWindow({
    icon: path.resolve(__dirname, iconPath),
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
      // 虽然一定会因此降低启动速度, 但是我只想降低开发成本。<懒得直接使用nodejs来加载文件了>
      (function startupSetting() {
        // 由于此部分仅开机首次运行时调用, 因此不受sdk中go依赖的viper的bug的影响。(即首次调用时可以获得真实情况, 若内部某字段被基于最终字段变更, 则会使得父字段为null的bug)
        // 而且, 这里还利用了StoreGet的false返回值来实现了递归轮询。(若直接基于最终字段来递归轮询, 则因其值本身就是boolean, 会无法达到实际的递归效果。)
        // TIPS: 以后再开新项目, 这类restful请求失败后的返回值, 不再使用false, 而是使用一个固定的字符串常量(具有绝对uuid特质的--似乎空对象之类的引用常量也行), 用于判断是否请求成功。
        StoreGet('startup').then((value) => {
          if (value === false) {
            startupSetting();
          } else {
            // 有些操作需要保证在sdk运行后再执行。此处利用了这一点。
            sdkIsRun = true;

            // TODO: 此部分窗口是否展示, 实际上不必依赖sdk的启动( 目前依赖的原因是->获取配置文件的对应值的StoreGet, 必须依赖sdk的启动。
            //       * 因此, 后续想要提升未设置窗口启动时因此情况下的窗口展示速度。需要直接获取配置文件的对应值, 并用此值设置窗口的展示状态。
            //       * * 方法1, 在sdk中, 优先获取此对应值, 并打印到标准输入, 然后基于进程调用相关返回值来设置窗口的展示状态。
            //       * * 方法2, 封装go的viper模块, 专门用于nodejs调用, 以通过nodejs调用go的viper模块, 来获取此对应值, 然后设置窗口的展示状态。(可以用进程调用或wasm调用等方式实现)
            //       * * 通过nodejs的相关配置包, 来获取此对应值, 然后设置窗口的展示状态。(由于需求较为单一, 可以不使用第三方的json包, 而是可以直接读取整个json文件, 并解析后直接使用这个值)
            if (value.is_hide_windows) {
              // 如果'is_hide_windows' 为 true , 则隐藏窗口
              mainWindow?.hide();
            } else {
              mainWindow?.show();
            }
          }
        });
      })();
    }
  });

  if (process.env.DEBUGGING) {
    // if on DEV or Production with debug enabled
    mainWindow.webContents.openDevTools({ mode: 'detach' });
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

//  建立一个全局可知的原始菜单模板(用于后续功能实现)<只需将此数组传入Menu.buildFromTemplate()中, 即可得到托盘菜单>
//  * TIPS: 需要注意的是, Menu.buildFromTemplate()在处理原始模板数组时, 不论外部的数组引用如何, 其只会根据数组的元素做判断
//          > 这个判断方式直接排除了相同引用的数组元素对象(即使你更改了某个数组元素对象的某个字段值), 也不会为你更新托盘菜单
//          > * 因此, 我们在需要修改某个托盘菜单时, 必须对相关元素对象彻底的解引用后, 再修改重构这个对象<只有破坏这个元素对象的引用才能得到托盘菜单的更新>。
//  * FIXME: 上方的也仅仅是猜测, 在通过array.map遍历数组时, 除了解构的方式中使用 i18n.global.t(item.label) 可以获得国际化返回值之外, 其它的方式只能得到原始字符串。
//           * 因此, 没有更新托盘菜单的原有也可能是i18n.global.t(item.label) 每次在遍历时, 以非解构的方式使用时, 并为获得正常返回值, 因此原始值都没变, 谈何更新。
const menuTemplate = [
  {
    label: 'Electron.tray.show',
    click: () => {
      mainWindow?.show();
    },
  },
  {
    label: 'Electron.tray.mute',
    click: () => {
      StoreSet('main_home.audio_volume_processing.volume_silent', true);
    },
  },
  {
    label: 'Electron.tray.quit',
    click: () => {
      (app as any).isQuiting = true;
      app.quit();
    },
  },
];

// 这个在每次使用Menu.buildFromTemplate()时, 都要调用, 以防止传入未经i18n运算的原始菜单模板
function menuTemplateI18n() {
  return menuTemplate.map((item) => {
    return {
      ...item,
      label: i18n.global.t(item.label),
    };
  });
}

// 在菜单模板中, 搜索所需菜单项的数组序号, 供后续变更使用
function searchItemIndexInMenuTemplate(label: string): number {
  for (let i = 0; i < menuTemplate.length; i++) {
    if (menuTemplate[i].label === label) {
      return i;
    }
  }
  return -1;
}

// 先通过轮询验证下可行性, 之后再引入sse或websocket。 (等引入sse或websocket时, 再来改动此部分代码, 目前能用就行。)
let history_language_default: string;

let history_volume_silent: boolean;

import AutoLaunch from 'auto-launch';

// 创建一个存储  AutoLaunch 实例的全局变量, 用于后续的自动启动设置
let autoLauncher: AutoLaunch;

setInterval(async () => {
  if (sdkIsRun) {
    // 托盘菜单的语言的设置
    StoreGet('language_default').then((req) => {
      if (req !== history_language_default) {
        // 希望只是没有响应式, 而不是无法用 (已验证, 希望是正确的, 可以使用。<即使在nodejs环境下, 也该遵守其类型去对其赋值>)
        // i18n.global.locale = req; // 错误用法, 未遵守其类型。 //   - completed(已完成)   FIXME: 此设置, 并未根本的更改国际化(由此可知vue-i18n 无法在nodejs中适配使用)
        i18n.global.locale.value = req; // 正确用法。
        // 先验证i18n是否生效
        // console.log('req', req);
        // console.log('i18n.global.t(Electron.tray.show)', i18n.global.t('Electron.tray.show'));
        // console.log('i18n.global.t(Electron.tray.quit)', i18n.global.t('Electron.tray.quit'));
        const contextMenu = Menu.buildFromTemplate(menuTemplateI18n());

        tray.setContextMenu(contextMenu);
        history_language_default = req;
      }
    });

    // 托盘菜单的静音按钮的设置
    StoreGet('main_home.audio_volume_processing.volume_silent').then((req) => {
      if (req != history_volume_silent) {
        if (req === true) {
          const index = searchItemIndexInMenuTemplate('Electron.tray.mute');
          menuTemplate[index] = {
            label: 'Electron.tray.unmute',
            click: () => {
              StoreSet('main_home.audio_volume_processing.volume_silent', false);
            },
          };
        } else if (req === false) {
          const index = searchItemIndexInMenuTemplate('Electron.tray.unmute');
          menuTemplate[index] = {
            label: 'Electron.tray.mute',
            click: () => {
              StoreSet('main_home.audio_volume_processing.volume_silent', true);
            },
          };
        }

        const contextMenu = Menu.buildFromTemplate(menuTemplateI18n());

        tray.setContextMenu(contextMenu);
        history_volume_silent = req;
      }
    });

    // 开机自启动的设置
    const is_hide_windows = await StoreGet('auto_startup.is_hide_windows');
    const is_auto_run = await StoreGet('auto_startup.is_auto_run');
    const is_hide_windows_old = await StoreGet('auto_startup.is_hide_windows_old');

    const isWindowsStore = process.windowsStore || process.env.WINDOWS_STORE;

    // 创建新的 AutoLaunch 实例
    autoLauncher = new AutoLaunch({
      name: 'KeyTone',
      // path: app.getPath('exe'), // 此库的官网上说:对于 NW.js 和 Electron 应用程序，您不必指定路径。我们根据 process.execPath 进行猜测。
      path: isWindowsStore ? 'LuSrackhall.KeyTone_yxzta3pw8j0pp!LuSrackhall' : undefined,
      isHidden: is_hide_windows,
    });

    // 检查并设置自动启动
    autoLauncher
      .isEnabled()
      .then((isEnabled: any) => {
        // 如果应用程序未设置在开机时自启动, 则主动设置, 若已设置, 则跳过。 此判断仅为防止重复开启。
        if ((!isEnabled && is_auto_run) || (isEnabled && is_hide_windows !== is_hide_windows_old)) {
          // if (value.is_auto_run) { // 我们必须避免重复开启, 虽然这样可以低成本实现value.is_hide_windows的实时响应, 但每次启动都去触碰敏感操作是不明智的。

          autoLauncher.enable().then(() => {
            // 如果窗口是否隐藏改变了, 则需要更新is_hide_windows_old以记录最新情况。 (由于我们为了防止重复设置, 只有在关闭自启动, 并开启自启动时, 才会触发重新设置自启动)
            if (is_hide_windows !== is_hide_windows_old) {
              StoreSet('auto_startup.is_hide_windows_old', is_hide_windows);
            }
          }); // 开启自启动
        }

        // // 如果应用已设置在开机时自启动, 则主动设置关闭自启动。 此判断仅为防止防止重复关闭。
        if (isEnabled && !is_auto_run) {
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
}, 1000);

function createTray() {
  // 创建托盘图标(开发环境也是可以创建托盘图标的, 之前失败的原因是图标路径的错误)
  tray = new Tray(iconPath); // 替换为你的图标路径

  // 创建托盘图标的上下文菜单
  const contextMenu = Menu.buildFromTemplate(menuTemplateI18n());

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
    // 虽然一定会因此降低启动速度, 但是我只想降低开发成本。<懒得直接使用nodejs来加载文件了>
    (function startupSetting() {
      StoreGet('startup').then((value) => {
        if (value === false) {
          startupSetting();
        } else {
          if (value.is_hide_windows) {
            // 如果默认隐藏, 则不作任何展示操作, 仅聚焦
            // 当运行第二个实例时,将会聚焦到myWindow这个窗口
            if (mainWindow) {
              if (mainWindow.isMinimized()) {
                mainWindow.restore();
                mainWindow.focus();
              }
            }
          } else {
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
          }
        }
      });
    })();
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
