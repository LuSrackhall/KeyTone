// .vitepress/theme/index.js
import DefaultTheme from "vitepress/theme";
import { h } from "vue";
import "./custom.css";
import HeroActions from "./components/HeroActions.vue";
import { useRouter, useRoute } from "vitepress";

// 用于标记是否已经执行过语言重定向，避免重复重定向而造成的语言切换功能无法正常使用的问题
let hasRedirected = false;

export default {
  // 继承 VitePress 默认主题
  extends: DefaultTheme,
  // 自定义布局组件
  Layout: () => {
    return h(DefaultTheme.Layout, null, {
      // 在首页 hero 区域的 actions 部分之后插入自定义组件
      // 参考文档：https://vitepress.dev/guide/extending-default-theme#layout-slots
      "home-hero-actions-after": () => h(HeroActions),
    });
  },
  // 增强应用程序的功能
  enhanceApp({ app, router, siteData }) {
    // 使用 Vue 的 mixin 在所有组件挂载时执行语言检查
    app.mixin({
      // 在组件挂载时执行
      mounted() {
        // 优先检查用户手动选择的语言
        const storedLang = sessionStorage.getItem("keytone-lang");
        if (storedLang) {
          hasRedirected = true;
          return;
        }

        // 如果已经执行过重定向，则直接返回，避免重复执行
        if (hasRedirected) return;

        // 获取用户浏览器的语言设置
        const userLanguage = navigator.language;
        // 获取当前路由信息
        const route = useRoute();
        // 获取路由实例，用于执行导航
        const router = useRouter();

        // 检查并执行重定向逻辑
        if (
          (userLanguage.startsWith("zh") && !route.path.startsWith("/zh")) ||
          (storedLang === "zh" && !route.path.startsWith("/zh"))
        ) {
          // 重定向到中文路径
          hasRedirected = true;
          sessionStorage.setItem("keytone-lang", "zh");
          router.go("/zh" + route.path);
        } else if (
          (!userLanguage.startsWith("zh") && route.path.startsWith("/zh")) ||
          (storedLang === "en" && route.path.startsWith("/zh"))
        ) {
          // 重定向到英文路径
          hasRedirected = true;
          sessionStorage.setItem("keytone-lang", "en");
          router.go(route.path.replace("/zh", ""));
        }
      },
    });
  },
};
