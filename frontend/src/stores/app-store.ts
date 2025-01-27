import { defineStore } from 'pinia';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';

// pinia中, 在defineStore的第二个参数中, 不能使用async函数, 否则会在运行过程中产生错误(在构建过程中及运行前不可知)。
// export const useAppStore = defineStore('app', async() => {
export const useAppStore = defineStore('app', () => {
  const q = useQuasar();
  const { t } = useI18n();

  const IPV4 = '127.0.0.1';

  // TIPS: 由于electron主进程中, 相关port的值更新的延迟性, 此处获取的仍有可能是38888这个初始值。(虽机率很小, 但仍需进一步迭代扼杀这个概率, 暂时的做法是在监听到端口变动后, 刷新整个渲染进程以初始化sse链接。)
  // TIPS: sse的链接一旦建立, 是无法动态更改其ip地址和端口的, 因此我们只能在建立之初就使用正确的端口, 否则只能通过刷新整个渲染进程来重新初始化sse链接了。
  //       => 毕竟手动关闭当前sse链接并重新建立新的sse链接涉及到的变更太多了, 几乎所有监听的回调都需重新调用一遍, 而回调的监听逻辑有可能分布在不同的文件中。
  //       => 当然, 如果我们能够将所有监听的回调都集中在一个文件中(或是封装逻辑后再暴露出去的方式将监听集中到某一个函数内方便重新建立), 那么我们只需在该文件中手动关闭当前sse链接并重新建立新的sse链接即可, 但这种方式的可维护性和可读性都不高。
  let port = 38888;
  if (process.env.MODE === 'electron') {
    port = window.myWindowAPI.getBackendPort();
  }

  const eventSource = new EventSource('http://' + IPV4 + `:${port}/stream`, { withCredentials: false });

  eventSource.onerror = function (event) {
    q.notify({
      type: 'info',

      position: 'top',

      message: t('Notify.正在尝试获取配置文件'),

      timeout: 1000,
    });

    setTimeout(() => {
      q.notify({
        type: 'negative',

        position: 'top',

        message: t('Notify.获取配置文件失败'),

        timeout: 1000,
      });
    }, 2800);
  };

  let openIsNotify = false;

  eventSource.onopen = function (event) {
    if (openIsNotify) {
      q.notify({
        type: 'info',

        position: 'bottom',

        message: t('Notify.配置文件获取成功, 正在读取配置'),

        timeout: 500,
      });

      setTimeout(() => {
        q.notify({
          type: 'positive',

          position: 'bottom',

          message: t('Notify.配置读取成功'),

          timeout: 3000,
        });
      }, 1000);
    }

    openIsNotify = true;
  };

  return { eventSource };
});
