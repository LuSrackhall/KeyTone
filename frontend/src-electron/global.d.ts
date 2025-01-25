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

// src-electron/global.d.ts
interface MyWindowAPI {
  minimize: () => void;
  toggleMaximize: () => void;
  close: () => void;
  openExternal: (arg0: string) => void;
  getWindowsStoreStatus: () => any;
  getBackendPort: () => Promise<number>;
}

interface Window {
  myWindowAPI: MyWindowAPI;
}
