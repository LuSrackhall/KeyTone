import { version } from "../../../../frontend/package.json";
import { LocaleSpecificConfig } from "vitepress";

// LocaleSpecificConfig类型接口, 为可以覆盖的属性 <能够被每个locale覆盖(包括 root)>
// * 这里的label, 为当前主题配置, 在Home页面国际化选择组件中, 对应选项的Label。
interface Config extends LocaleSpecificConfig {
  label: string;
  // link?: string; // 当前配置作为root使用, 因此不能配置link字段来自定义url后缀, 否则会引起主页面异常。
}

export const META_URL = "https://keytone.xuanhall.com";
export const META_TITLE = "KeyTone";
export const META_DESCRIPTION =
  "Can simulate the sound of keyboard strokes in real-time, supports multi-platform Windows, Mac and Linux, easy fast installation. The pressing and releasing of keyboard keys have independent sound effects, perfectly matching scenarios where keys are held down for a long time. Ready to use out of the box, Attention to detail, In silence, let your keyboard bring forth pleasing sounds.";

export const enConfig: Config = {
  label: "English",
  lang: "en",
  // 应用的changelog内容中要注意, 由于个人习惯使用<>, 但<>内容中如果是英文的开头, 则需要一个空格, 否则会造成在页面上解析失败的报错。如<Ni Hao>要改成< Ni Hao>才行。

  description: META_DESCRIPTION,
  head: [
    ["meta", { property: "og:url", content: META_URL }],
    ["meta", { property: "og:description", content: META_DESCRIPTION }],
    ["meta", { property: "twitter:url", content: META_URL }],
    ["meta", { property: "twitter:title", content: META_TITLE }],
    ["meta", { property: "twitter:description", content: META_DESCRIPTION }],
  ],

  themeConfig: {
    // sidebarMenuLabel: "Menu",  // 默认为"Menu", 无需显示配置
    // returnToTopLabel: "Return to top",  // 默认为"Return to top", 无需显示配置
    // aside: true, // 默认为true, 则页面导航显示在右边, 配置false则关闭显示, 配置"left"则页面导航会显示在左边。如果想对所有页面禁用它，应该使用 outline: false。
    outline: { level: 3, label: "On this page" }, // level:3 代表只显示3级标题<即会忽略1、2、4、5、6级标题>, 如果想显示更多的标题级别, 可以通过数组来配置level, 如level:[2,3,4,5]这样子配置<不过需要注意的是, 请确保你的文章中真的有这些标题级别, 若是你没有4级标题却有2,3,5级标题, 请配置level:[2,3,5]>
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
