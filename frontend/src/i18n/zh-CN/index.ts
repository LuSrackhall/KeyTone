export default {
  KeyTone: {
    KeyTone: 'KeyTone', // KeyTone: '键音', // 如果不需要描述, 则在最后一级直接使用 `字符串`
    developer: '制作者',
    setting: { index: '设置', caption: '' }, // 如果需要describe(caption) , 则在最后一级使用 `对象`{index: '...' , caption: '...'}
  },
  setting: {
    setting: '设置',
    language: { index: '语言', caption: '' }, // 如果需要describe(caption) , 则在最后一级使用 `对象`{index: '...' , caption: '...'}
    启动与自动启动: {
      启动与自动启动: { index: `启动与自动启动`, caption: '' },
      启动时隐藏窗口: { index: '启动时隐藏窗口', caption: '打开应用时最小化到系统托盘' },
      自动启动: { index: '自动启动', caption: '是否要在电脑开机时自动启动本应用<注: 重启本应用后生效>' },
      自动启动时隐藏窗口: {
        index: `自动启动时隐藏窗口`,
        caption: '电脑开机自动启动本应用时最小化到系统托盘<注: 重启本应用后生效>',
      },
    },
  },

  language: {
    'setting language': '设置语言/setting language',
  },
  Electron: {
    tray: {
      show: '展示',
      close: '关闭',
    },
  },
};
