export default {
  KeyTone: {
    KeyTone: 'KeyTone', // KeyTone: '键音', // 如果不需要描述, 则在最后一级直接使用 `字符串`
    developer: '制作者',
    setting: { index: '设置', caption: '' }, // 如果需要describe(caption) , 则在最后一级使用 `对象`{index: '...' , caption: '...'}
  },
  setting: {
    setting: '设置',
    language: { index: '语言', caption: '' }, // 如果需要describe(caption) , 则在最后一级使用 `对象`{index: '...' , caption: '...'}
  },

  language: {
    'setting language': '设置语言/setting language',
  },
};
