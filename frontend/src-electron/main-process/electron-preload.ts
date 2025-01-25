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

// 从 electron 模块中导入 contextBridge，用于在渲染进程和主进程之间安全地暴露API。
// import { contextBridge } from 'electron';
const { contextBridge, ipcRenderer } = require('electron');
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
    // console.log('minimize, 窗口最小化');
  },

  // toggleMaximize：切换当前窗口的最大化状态。
  toggleMaximize() {
    const win = BrowserWindow.getFocusedWindow();

    if (win?.isMaximized()) {
      win.unmaximize();
      // console.log('unmaximize, 窗口取消最大化');
    } else {
      win?.maximize();
      // console.log('maximize, 窗口最大化');
    }
  },

  // close：关闭当前窗口。
  close() {
    BrowserWindow.getFocusedWindow()?.close();
    // console.log('close, 关闭当前窗口');
  },

  // 用于在外部浏览器打开url链接
  openExternal: (url: string) => ipcRenderer.send('open-external', url),

  // 获取windows store状态
  getWindowsStoreStatus() {
    return process.windowsStore || process.env.WINDOWS_STORE;
  },

  // 添加获取后端端口的方法
  getBackendPort: () => ipcRenderer.invoke('get-backend-port'),
});
