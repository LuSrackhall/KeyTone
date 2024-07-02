import { defineConfig } from "vitepress";
import { enConfig } from "./config/en.mts";
import { zhConfig } from "./config/zh.mts";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  // lang: "en-US", // 这个会被locales配置中root的lang配置给覆盖掉(当然root配置代表默认的配置, 未被root覆盖的配置仍会使用当前配置文件的内容; 以此类推所有国际化配置。 <另外, root配置不会在url中将root这个key名的作为默认link追加, 也不允许自行配置link>)
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

  locales: {
    // 当key名为root时, 配置不会在url中将root这个key名的作为默认link值追加, 也不允许在其value的对象中自行配置link。
    root: {
      ...enConfig,
    },
    // 这里的的key名, 将会作为默认的link值, 在url后进行追加。 但是, 可以在value的对象中通过自行配置link字段来覆盖此默认的link值。
    zh: {
      ...zhConfig,
    },
  },
});
