import enUS from './en-US/index.json';
import zhCN from './zh-CN/index.json';
import zhHK from './zh-HK/index.json';

// 一下键值对中的key的名称, 以quasar框架(或其它的项目框架)的语言本地化api的列表为准。(如quasar的列表默认值的参考链接https://github.com/quasarframework/quasar/tree/dev/ui/lang)
export default {
  'en-US': enUS,
  'zh-CN': zhCN,
  'zh-HK': zhHK, // 由于quasar相关语言本地化api的列表默认值中没有zh-HK相关的, 因此此行代码理论上不起任何作用。(万一quasar以后支持了, 就省的适配了。)(设置列表中对中文繁体的相关key选择为zh-HK或zh-TW都行, 都代表中文繁体, 但只能选择其中一个, 目前用的是zh-TW, 毕竟quasar目前只支持zh-TW。)
  'zh-TW': zhHK,
};
