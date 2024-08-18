import { defineStore } from 'pinia';
import { useQuasar } from 'quasar';

export const useAppStore = defineStore('app', () => {
  const q = useQuasar();

  const IPV4 = '127.0.0.1';

  const eventSource = new EventSource('http://' + IPV4 + ':38888/stream', { withCredentials: false });

  eventSource.onerror = function (event) {
    q.notify({
      type: 'info',

      position: 'top',

      message: '正在尝试获取配置文件...',

      timeout: 1000,
    });

    setTimeout(() => {
      q.notify({
        type: 'negative',

        position: 'top',

        message: '获取配置文件失败',

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

        message: '配置文件获取成功, 正在读取配置...',

        timeout: 500,
      });

      setTimeout(() => {
        q.notify({
          type: 'positive',

          position: 'bottom',

          message: '配置读取成功',

          timeout: 3000,
        });
      }, 1000);
    }

    openIsNotify = true;
  };

  return { eventSource };
});
