// src-electron/global.d.ts
interface MyWindowAPI {
  minimize: () => void;
  toggleMaximize: () => void;
  close: () => void;
  openExternal: (arg0: string) => void;
}

interface Window {
  myWindowAPI: MyWindowAPI;
}
