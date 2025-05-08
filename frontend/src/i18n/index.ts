// 个人对i18n文件的明明规范, 以中文简体为例:
// * zh的字母简写, 来自ISO 639-2: https://www.loc.gov/standards/iso639-2/php/code_list.php
// * CH的字母简写, 来自ISO 3166: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes
// * 当然最后, 还需要通过搜索引擎或ai搜索"i18n中 '中文简体'的简写字母是", 以判断最终正确性。 (当然, 这里还有常用的i18n简写字母介绍 https://github.com/teojs/vue-resume/issues/5)
// * quasar框架内的语言本地化api的支持列: https://github.com/quasarframework/quasar/tree/dev/ui/lang (我们的标准不参考这个, 这个主要用于下方导出键值对中的key的名称)
import enUS from './en-US/index.json';
import zhCN from './zh-CN/index.json';
import zhHK from './zh-HK/index.json';
import jaJP from './ja-JP/index.json';
import koKR from './ko-KR/index.json';
import deDE from './de-DE/index.json';
import ruRU from './ru-RU/index.json';
import frFR from './fr-FR/index.json';
import itIT from './it-IT/index.json';
import esES from './es-ES/index.json';
import ptPT from './pt-PT/index.json';
import ptBR from './pt-BR/index.json';
import plPL from './pl-PL/index.json';
import trTR from './tr-TR/index.json';
import viVN from './vi-VN/index.json';

// 一下键值对中的key的名称, 以quasar框架(或其它的项目框架)的语言本地化api的列表为准。(如quasar的列表默认值的参考链接https://github.com/quasarframework/quasar/tree/dev/ui/lang)
export default {
  'en-US': enUS,
  'zh-CN': zhCN,
  // 'zh-HK': zhHK, // 由于quasar相关语言本地化api的列表默认值中没有zh-HK相关的, 因此此行代码理论上不起任何作用。(万一quasar以后支持了, 就省的适配了。)(设置列表中对中文繁体的相关key选择为zh-HK或zh-TW都行, 都代表中文繁体, 但只能选择其中一个, 目前用的是zh-TW, 毕竟quasar目前只支持zh-TW。)
  'zh-TW': zhHK,
  ja: jaJP,
  'ko-KR': koKR,
  de: deDE,
  'de-DE': deDE,
  'de-CH': deDE,
  ru: ruRU,
  fr: frFR,
  it: itIT,
  es: esES,
  pt: ptPT,
  'pt-BR': ptBR,
  pl: plPL,
  tr: trTR,
  vi: viVN,
};
