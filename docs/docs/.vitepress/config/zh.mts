import { version } from "../../../../frontend/package.json";
import { LocaleSpecificConfig } from "vitepress";

// LocaleSpecificConfig类型接口, 为可以覆盖的属性 <能够被每个locale覆盖(包括 root)>
// * 这里的label, 为当前主题配置, 在Home页面国际化选择组件中, 对应选项的Label。
interface Config extends LocaleSpecificConfig {
  label: string;
  link?: string; // 当前配置不用于root, 因此可以配置link字段来自定义url后缀。
}

export const zhConfig: Config = {
  label: "中文",
  lang: "zh",
  // link: "/zh/test", // 这个是用于改变默认url的后缀的<默认的url后缀名获取方式, 参考当前项目.vitepress/config.mts中locales字段内容的注释介绍方可获悉>。

  themeConfig: {
    nav: [
      { text: "主页", link: "/" },
      // { text: "使用指南", link: "/zh/guide/", activeMatch: "/zh/guide/" },
      {
        text: "v" + version,
        items: [{ text: "更新日志", link: "zh/changelog/v" + version, activeMatch: "/changelog/" }],
      },
    ],
    sidebar: {
      // "/zh/guide/": { base: "/zh/guide", items: sidebarUserGuide() },
      "zh/changelog/": { base: "zh/changelog", items: sidebarChangelog() },
    },
    footer: {
      message: "基于GPL-3.0开源许可协议",
      copyright: "Copyright (C) 2024 LuSrackhall",
    },
  },
};

// function sidebarUserGuide() {}

function sidebarChangelog() {
  return [
    {
      text: "更新日志",
      items: [
        // { text: "v0.2.0", link: "/v0.2.0" }, // 最新的changelog, 往上写就好了
        { text: "v0.1.0", link: "/v0.1.0" },
      ],
    },
  ];
}
