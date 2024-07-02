import { version } from "../../../../frontend/package.json";
import { LocaleSpecificConfig } from "vitepress";

// LocaleSpecificConfig类型接口, 为可以覆盖的属性 <能够被每个locale覆盖(包括 root)>
// * 这里的label, 为当前主题配置, 在Home页面国际化选择组件中, 对应选项的Label。
interface Config extends LocaleSpecificConfig {
  label: string;
  // link?: string; // 当前配置作为root使用, 因此不能配置link字段来自定义url后缀, 否则会引起主页面异常。
}

export const enConfig: Config = {
  label: "English",
  lang: "en",
};
