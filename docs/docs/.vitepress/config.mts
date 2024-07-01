import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  lang: "en-US",
  title: "KeyTone",
  // description: "In silence, let your keyboard bring forth pleasing sounds.",
  description: "Can simulate the sound of keyboard strokes in real-time",
  cleanUrls: true, // WARN: 设置成true, 可以使得url后面没有.html, 不过对于部署来说, 对不同的部署方式需要额外配置, 才能支持此功能。

  // assetsDir: "static",  // 这个是build后,spa页面静态资源的目录, 默认为assetsDir。

  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config

    // logo: "/KeyTone.png", // 目前这个logo不适合在这里使用, 不美观。
    // siteTitle: false, // 当为字符串时, 可以自定义此项目以替换导航中的默认网站标题（应用程序配置中的 title ）。当设置为 false 时，导航中的标题将被禁用。当您的 logo 已包含网站标题文本时很有用。

    // 后续要增加国际化, 此导航需要在各个国际化的config配置文件中独立配置, 而不是在此处。
    // nav: [
    //   { text: "Home", link: "/" },
    //   { text: "Examples", link: "/markdown-examples" },
    // ],

    // 后续要增加国际化, 此导航需要在各个国际化的config配置文件中独立配置, 而不是在此处。
    // sidebar: [
    //   {
    //     text: "Examples",
    //     items: [
    //       { text: "Markdown Examples", link: "/markdown-examples" },
    //       { text: "Runtime API Examples", link: "/api-examples" },
    //     ],
    //   },
    // ],

    socialLinks: [{ icon: "github", link: "https://github.com/LuSrackhall/KeyTone" }],
  },
});
