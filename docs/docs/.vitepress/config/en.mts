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
  // 应用的changelog内容中要注意, 由于个人习惯使用<>, 但<>内容中如果是英文的开头, 则需要一个空格, 否则会造成在页面上解析失败的报错。如<Ni Hao>要改成< Ni Hao>才行。

  themeConfig: {
    nav: [
      { text: "Home", link: "/" },
      // { text: "User Guide", link: "/guide/", activeMatch: "/guide/" },
      {
        text: "v" + version,
        items: [{ text: "Changelog", link: "/changelog/v" + version, activeMatch: "/changelog/" }],
      },
    ],
    sidebar: {
      // "/guide/": { base: "/guide", items: sidebarUserGuide() },
      "/changelog/": { base: "/changelog", items: sidebarChangelog() },
    },
    footer: {
      message: "Released under the GPL-3.0 License.",
      copyright: "Copyright (C) 2024 LuSrackhall",
    },
  },
};

// function sidebarUserGuide() {}

function sidebarChangelog() {
  return [
    {
      text: "Changelog",
      items: [
        // { text: "v0.2.0", link: "/v0.2.0" }, // 最新的changelog, 往上写就好了
        { text: "v0.1.0", link: "/v0.1.0" },
      ],
    },
  ];
}
