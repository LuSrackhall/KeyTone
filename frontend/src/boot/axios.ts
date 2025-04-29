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
    // 这个逻辑我们在启动时先即时的调用一次, 以避免端口不一致时造成启动瞬间仍使用旧的38888端口。(后续的setInterval()存在1s后才执行第一次的问题(尤其在macos中), 不会即时调用, 这会影响获取实际端口的及时性。)(而且, 由于新架构下是一定可以第一时间及时获取到真实端口的, 因此后续setInterval中的6s监听也没有必要, 但保留它做个保障也行, 多几行无用代码也无可厚非。)
    const currentPort = window.myWindowAPI.getBackendPort(); // 这个逻辑放入boot的回调中恰到好处, 因为electron项目中我们无需在nodejs的主进程使用它。
    UpdateApi(currentPort);
    if (currentPort !== port) {
      // 刷新重启整个应用(以刷新sse的端口)
      window.location.reload(); // 在此处可用
    }

    // TIPS: 由于即使端口变化, 也只会变化一次, 或是无需变化(99%的概率无需变化, 即使端口占用, sdk返回的端口也极大概率在渲染进程启动前返回)。
    //       * 因此通过intervalId, 让其至多触发一次。(或是通过count, 让无需变化的情况下6s后自动停止, 避免持续监听造成的资源浪费)
    let count = 0;
    const intervalId = setInterval(async () => {
      count++;
      const currentPort = window.myWindowAPI.getBackendPort(); // 这个逻辑放入boot的回调中恰到好处, 因为electron项目中我们无需在nodejs的主进程使用它。
      UpdateApi(currentPort);
      if (currentPort !== port) {
        // 刷新重启整个应用(以刷新sse的端口)
        window.location.reload(); // 在此处可用
        clearInterval(intervalId);
      }
      if (count > 6) {
        clearInterval(intervalId);
      }
      // console.log('count', count);
    }, 1000);
    // ^^^^^^ TIPS: 上述是最终选择的方案, 下面介绍被放弃的方案
    //  1. 在主进程中主动触发当前渲染进程中更新api实例的逻辑。
    //       * electron的ipc提供了on/once的api, 可以监听对应事件并触发更新。(不过循环监听事件的方法不一定消耗比这个小)
    //         * 可以利用on持续监听, 也可以利用once监听一次。(并且在electron的preload对应的on或once中, 也设置6s后自动停止, 避免持续监听造成的资源浪费)
    //  2. 在主进程中, 控制生命周期, 已在开始启动渲染进程之前, 完成port的最终确认。
    //     * 可以在 createWindow(); 调用之前, 确认端口好是否亦有sdk返回。
    //     * 也可精确的 在createWindow内部逻辑的 mainWindow.loadURL(process.env.APP_URL); 或 mainWindow.loadFile('index.html'); 调用之前, 确认端口好是否亦有sdk返回。
    //     * 如果未有sdk返回, 则不启动渲染进程, 并循环调用, 直到有sdk而端口数据返回。 端口数据返回, 则启动渲染进程, 并停止循环调用。
    //     * 缺点: 需要精确控制生命周期, 且需要精确控制渲染进程的启动时机。(可能会影响窗口的启动速度)
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
