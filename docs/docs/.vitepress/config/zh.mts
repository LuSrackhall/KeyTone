import { version } from "../../../../frontend/package.json";
import { LocaleSpecificConfig } from "vitepress";

// LocaleSpecificConfig类型接口, 为可以覆盖的属性 <能够被每个locale覆盖(包括 root)>
// * 这里的label, 为当前主题配置, 在Home页面国际化选择组件中, 对应选项的Label。
interface Config extends LocaleSpecificConfig {
  label: string;
  link?: string; // 当前配置不用于root, 因此可以配置link字段来自定义url后缀。
}

export const META_URL = "https://keytone.xuanhall.com";
export const META_TITLE = "KeyTone";
export const META_DESCRIPTION =
  "可以实时模拟键盘敲击的声音, 支持多平台Windows、Mac和Linux系统, 安装简单迅速。按键的按下和抬起, 拥有独立音效, 完美适配按键长按的场景。开箱即用, 注重细节, 在寂静中，让您的键盘唤醒舒适的声音。";

export const zhConfig: Config = {
  label: "中文",
  lang: "zh",
  // link: "/zh/test", // 这个是用于改变默认url的后缀的<默认的url后缀名获取方式, 参考当前项目.vitepress/config.mts中locales字段内容的注释介绍方可获悉>。

  description: META_DESCRIPTION,
  head: [
    ["meta", { property: "og:url", content: META_URL }],
    ["meta", { property: "og:description", content: META_DESCRIPTION }],
    ["meta", { property: "twitter:url", content: META_URL }],
    ["meta", { property: "twitter:title", content: META_TITLE }],
    ["meta", { property: "twitter:description", content: META_DESCRIPTION }],
  ],

  themeConfig: {
    sidebarMenuLabel: "菜单",
    returnToTopLabel: "返回顶部",
    aside: true, // 默认为true, 则页面导航显示在右边, 配置false则关闭显示, 配置"left"则页面导航会显示在左边。如果想对所有页面禁用它，应该使用 outline: false。
    outline: { level: 3, label: "页面导航" }, // level:3 代表只显示3级标题<即会忽略1、2、4、5、6级标题>, 如果想显示更多的标题级别, 可以通过数组来配置level, 如level:[2,3,4,5]这样子配置<不过需要注意的是, 请确保你的文章中真的有这些标题级别, 若是你没有4级标题却有2,3,5级标题, 请配置level:[2,3,5]>
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
        { text: "v0.3.0", link: "/v0.3.0" },
        { text: "v0.2.0", link: "/v0.2.0" },
        { text: "v0.1.0", link: "/v0.1.0" },
      ],
    },
  ];
}
