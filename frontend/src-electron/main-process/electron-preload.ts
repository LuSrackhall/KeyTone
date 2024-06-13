// 从 electron 模块中导入 contextBridge，用于在渲染进程和主进程之间安全地暴露API。
// import { contextBridge } from 'electron';
const { contextBridge } = require('electron');
// 从 @electron/remote 模块中导入 BrowserWindow，用于在渲染进程中控制窗口。
import { BrowserWindow } from '@electron/remote';

/*
 * 使用 contextBridge.exposeInMainWorld 方法在渲染进程的全局上下文中暴露一个名为 myWindowAPI 的对象。这个对象包含三个方法：
 * minimize：最小化当前窗口。
 * toggleMaximize：切换当前窗口的最大化状态。
 * close：关闭当前窗口。
 */
contextBridge.exposeInMainWorld('myWindowAPI', {
  // minimize：最小化当前窗口。
  minimize() {
    BrowserWindow.getFocusedWindow()?.minimize();
    console.log('minimize, 窗口最小化');
  },

  // toggleMaximize：切换当前窗口的最大化状态。
  toggleMaximize() {
    const win = BrowserWindow.getFocusedWindow();

    if (win?.isMaximized()) {
      win.unmaximize();
      console.log('unmaximize, 窗口取消最大化');
    } else {
      win?.maximize();
      console.log('maximize, 窗口最大化');
    }
  },

  // close：关闭当前窗口。
  close() {
    BrowserWindow.getFocusedWindow()?.close();
    console.log('close, 关闭当前窗口');
  },
});
