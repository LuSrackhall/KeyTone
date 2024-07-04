// .vitepress/theme/index.js
import DefaultTheme from "vitepress/theme";
import { h } from "vue";
import "./custom.css";
import HeroActions from "./components/HeroActions.vue";

export default {
  extends: DefaultTheme,
  Layout: () => {
    return h(DefaultTheme.Layout, null, {
      // https://vitepress.dev/guide/extending-default-theme#layout-slots
      "home-hero-actions-after": () => h(HeroActions),
    });
  },
  enhanceApp({ app, router, siteData }) {
    // ...
  },
};
