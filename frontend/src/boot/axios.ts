import { boot } from 'quasar/wrappers';
import axios, { AxiosInstance } from 'axios';

// TIPS: declare module 只影响 TypeScript 的类型检查和开发时的代码提示，完全不会影响实际的编译输出和运行时行为。它属于 TypeScript 的类型系统的一部分，在编译成 JavaScript 后会被完全移除。
// declare module '@vue/runtime-core' { // 是 Vue 2 时代的遗留写法
declare module 'vue' {
  // 是 Vue 3 推荐的方式, 直接扩展 Vue 3 的主模块, 更符合 Vue 3 的模块化设计, 覆盖范围更完整
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
    $api: AxiosInstance;
  }
}

// Be careful when using SSR for cross-request state pollution
// due to creating a Singleton instance here;
// If any client changes this (global) instance, it might be a
// good idea to move this instance creation inside of the
// "export default () => {}" function below (which runs individually
// for each client)
const api = axios.create({ baseURL: 'http://127.0.0.1:38888' });

export default boot(({ app }) => {
  // for use inside Vue files (Options API) through this.$axios and this.$api

  app.config.globalProperties.$axios = axios;
  // ^ ^ ^ this will allow you to use this.$axios (for Vue Options API form)
  //       so you won't necessarily have to import axios in each vue file

  app.config.globalProperties.$api = api;
  // ^ ^ ^ this will allow you to use this.$api (for Vue Options API form)
  //       so you can easily perform requests against your app's API
});

export { api };
