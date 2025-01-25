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

let port = 38888;
// Be careful when using SSR for cross-request state pollution
// due to creating a Singleton instance here;
// If any client changes this (global) instance, it might be a
// good idea to move this instance creation inside of the
// "export default () => {}" function below (which runs individually
// for each client)
// 初始化api实例(默认以端口38888为基准(前后端统一))
let api: AxiosInstance = axios.create({ baseURL: `http://127.0.0.1:${port}` });

// 用于更新api实例(当端口发生变化时)
export const UpdateApi = (newPort: number) => {
  if (newPort !== port) {
    port = newPort;
    api = axios.create({ baseURL: `http://127.0.0.1:${port}` });
  }
};

// 如果您需要访问Object中的参数以链接应用程序的其它部分, 则将内容写在boot的回调中, boot函数会由quasar在vue的main.ts中自动调用。(如果不需要访问Object中的参数, 则无需在boot内部处理(当然, 写进去也无可厚非)。)
export default boot(({ app }) => {
  // 由于前端对于端口的变更时不确定的, 因此需要利用ipc持续监听端口变化，来更新api实例。(TIPS: 虽然本项目不涉及spa, 但后续有必要思考, spa中如何监听端口变化以做到前后端统一的应对。比如使用go启动某个spa时, 是否能做到向spa中传递一些参数这种事情。)
  if (process.env.MODE === 'electron') {
    setInterval(async () => {
      const currentPort = window.myWindowAPI.getBackendPort(); // 这个逻辑放入boot的回调中恰到好处, 因为electron项目中我们无需在nodejs的主进程使用它。
      UpdateApi(currentPort);
    }, 1000);
  }

  // for use inside Vue files (Options API) through this.$axios and this.$api

  app.config.globalProperties.$axios = axios;
  // ^ ^ ^ this will allow you to use this.$axios (for Vue Options API form)
  //       so you won't necessarily have to import axios in each vue file

  app.config.globalProperties.$api = api;
  // ^ ^ ^ this will allow you to use this.$api (for Vue Options API form)
  //       so you can easily perform requests against your app's API
});

export { api };
