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
  "可以实时模拟按键敲击的声音, 支持多平台Windows、Mac和Linux系统, 安装简单迅速。按键的按下和抬起, 拥有独立音效, 完美适配按键长按的场景。开箱即用, 注重细节, 在寂静中，让您的按键唤醒舒适的声音。";

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
    outline: { level: [2, 3, 4], label: "页面导航" }, // level:3 代表只显示3级标题<即会忽略1、2、4、5、6级标题>, 如果想显示更多的标题级别, 可以通过数组来配置level, 如level:[2,3,4,5]这样子配置<不过需要注意的是, 请确保你的文章中真的有这些标题级别, 若是你没有4级标题却有2,3,5级标题, 请配置level:[2,3,5]>
    nav: [
      { text: "主页", link: "/zh" },
      { text: "使用指南", link: "/zh/guide/getting-started/installation", activeMatch: "/zh/guide/" },
      //#region 解释一下这些配置字段的作用：
      /*
        在 VitePress 的导航配置中，这行代码定义了顶部导航栏中的一个链接项：

        ```typescript:docs/docs/.vitepress/config/zh.mts
        { text: "使用指南", link: "/zh/guide/getting-started/installation", activeMatch: "/zh/guide/" }
        ```

        这个配置包含三个重要部分：

        1. `text: "使用指南"` - 显示在导航栏上的文本
        2. `link: "/zh/guide/getting-started/installation"` - 点击后跳转的目标路径
        3. `activeMatch: "/zh/guide/"` - 定义何时将此导航项标记为"激活"状态

        其中 `activeMatch` 特别重要：
        - 它使用正则表达式来匹配当前 URL 路径
        - 当用户访问的页面 URL 包含 "/zh/guide/" 时，这个导航项会被标记为激活状态（通常会有特殊的样式，如高亮显示）
        - 这样可以让用户知道他们当前在网站的哪个部分

        如果不设置 `activeMatch`，导航项只会在完全匹配 `link` 指定的路径时才会显示为激活状态。通过设置 `activeMatch`，可以让该导航项在访问其下的所有子页面时都保持激活状态。
      */
      //#endregion
      {
        text: "v" + version,
        items: [{ text: "更新日志", link: "zh/changelog/v" + version, activeMatch: "/changelog/" }],
      },
    ],
    sidebar: {
      "/zh/guide/": { base: "/zh/guide", items: sidebarUserGuide() },
      "zh/changelog/": { base: "zh/changelog", items: sidebarChangelog() },
    },
    footer: {
      message: "KeyTone基于GNU GPLv3许可证发布",
      copyright: `版权所有 (C) 2024-${new Date().getFullYear()} LuSrackhall`,
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
        { text: "v0.4.2", link: "/v0.4.2" },
        { text: "v0.4.1", link: "/v0.4.1" },
        { text: "v0.4.0", link: "/v0.4.0" },
        { text: "v0.3.6", link: "/v0.3.6" },
        { text: "v0.3.5", link: "/v0.3.5" },
        { text: "v0.3.4", link: "/v0.3.4" },
        { text: "v0.3.3", link: "/v0.3.3" },
        { text: "v0.3.2", link: "/v0.3.2" },
        { text: "v0.3.1", link: "/v0.3.1" },
        { text: "v0.3.0", link: "/v0.3.0" },
        { text: "v0.2.0", link: "/v0.2.0" },
        { text: "v0.1.0", link: "/v0.1.0" },
      ],
    },
  ];
}

function sidebarUserGuide() {
  return [
    {
      text: "开始",
      items: [
        { text: "安装", link: "/getting-started/installation/" },
        { text: "快速上手", link: "/getting-started/quick-start/" },
      ],
    },
    {
      text: "键音专辑",
      items: [
        // { text: "介绍", link: "/key-package/introduction/index.md" }, // 如果最终文件名为index, 则需要加index.md后缀, 否则导航无法正常高亮。这是因为activeMatch而默认值是link,也可通过手动配置activeMatch来实现导航高亮。
        { text: "介绍", link: "/key-package/introduction/" }, // 如果最终文件名为index, 则路径可仅导航到index的所属目录, 而且需要加/后缀, 否则导航无法正常高亮。这是因为activeMatch而默认值是link,也可通过手动配置activeMatch来实现导航高亮。
        { text: "载入音频文件", link: "/key-package/载入音频文件/" },
        { text: "裁剪定义声音", link: "/key-package/裁剪定义声音/" },
        { text: "铸造至臻键音", link: "/key-package/铸造至臻键音/" },
        { text: "按键联动声效", link: "/key-package/按键联动声效/" },
      ],
    },
    {
      text: "其他",
      items: [
        { text: "常见问题", link: "/other/faq/" },
        { text: "隐私政策", link: "/other/privacy-policy/" },
        { text: "用户协议", link: "/other/user-agreement/" },
      ],
    },
  ];
}
