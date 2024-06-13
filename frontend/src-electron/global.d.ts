// src-electron/global.d.ts
interface MyWindowAPI {
  minimize: () => void;
  toggleMaximize: () => void;
  close: () => void;
}

interface Window {
  myWindowAPI: MyWindowAPI;
}
