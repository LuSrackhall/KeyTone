export default {
  KeyTone: {
    KeyTone: 'KeyTone',
    developer: 'Developer',
    setting: { index: 'Setting', caption: '' },
  },
  setting: {
    setting: 'Setting',
    language: { index: 'Language', caption: '' },
    启动与自动启动: {
      启动与自动启动: { index: 'Startup and Auto Start', caption: '' },
      启动时隐藏窗口: {
        index: 'Hide Window on Startup',
        caption: 'Minimize to system tray when opening the application',
      },
      自动启动: {
        index: 'Auto Start',
        caption:
          'Would you like this application to automatically start when the computer boots up? <Note: This will take effect after restarting the application.>',
      },
      自动启动时隐藏窗口: {
        index: 'Hide Window on Auto Start',
        caption:
          'Minimize this application to the system tray when the computer starts up and automatically launches this application. <Note: This will take effect after restarting the application.>',
      },
    },
  },

  language: {
    'setting language': 'setting language/设置语言',
  },
  Electron: {
    tray: {
      show: 'Show',
      close: 'Close',
    },
  },
};
