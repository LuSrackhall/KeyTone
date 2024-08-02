import { defineStore } from 'pinia';
import { Quasar, useQuasar } from 'quasar';
import { StoreGet, StoreSet } from 'src/boot/query/store-query';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

export const useSettingStore = defineStore('setting', () => {
  /*------------------------------------------------------------------------------------------------------------------*/
  const q = useQuasar();

  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/

  //#region    -----<<<<<<<<<<<<<<<<<<<< --languageDefault  start ^_^-_-^_^
  /**
   * 此处为, 设置页面定义语言项国际化页面 m 的区域
   * * 即 通过m对接全局相关vm的状态, 从而体现至v视图UI层
   */

  const language = Quasar.lang.getLocale();
  const languageDefault = ref<string>(language ? language : 'en-US');

  const { locale } = useI18n({ useScope: 'global' });

  /**
   * 此处采用立即执行, 是为了当数据库为空时, 也能自动赋值所获取的当地语言值(即使不适配也无妨), 不至于所获得值被刷而被动强制
   * 为英语。
   * * 即使当地语言我没有适配,i18n系统也会在对比不合适后, 反而执行启动文件中, 设置的默认值, 即强制en-us语言为系统语言。
   */
  watch(
    languageDefault,
    () => {
      locale.value = languageDefault.value;
    },
    { immediate: true }
  );
  //#endregion ----->>>>>>>>>>>>>>>>>>>> --languageDefault  end   -_-^_^-_- ^_^-_-^_^-_-
  // ...
  // ...
  // ...
  //!endregion ----->>>>>>>>>>>>>>>>>>>> --languageDefault  end   -_-^_^-_- ^_^-_-^_^-_-

  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/

  //#region    -----<<<<<<<<<<<<<<<<<<<< -- setting持久化 start ^_^-_-^_^

  async function settingInitAndRealTimeStorage() {
    // 优先使用数据库中保存的设置, 即先通过数据库存储, 对内存做初始化
    await StoreGet('get_all_value').then((req) => {
      // console.debug('打印观察获取的值', req);
      if (req === false) {
        // 此时, 说明GetItem_sqlite请求过程中, 出错了, 因此需要错误通知, 并让用户重新启动, 防止用户因继续使用造成的存储设置被初始覆盖
        q.notify({
          type: 'negative',
          position: 'top',
          message: '数据库读取失败,请重启应用',
          timeout: 100000,
        });
        return;
      }

      // TIPS: 由于采取各设置独立的录入即判别方式, 不再依赖整体的JSON字符串, 因此此if判断后续可能没必要存在(目前暂时保留)
      // 第一次进入本应用, 设置本就该是空的, 此时无需对我们的设置项进行任何操作, 也无需做任何通知。
      // 但为防止后续的JSON.parse报错, 因此此处也是必不可少的(因为只要非首次, 就不可能为空, watchEffect是立即执行的, 也就是说至少整体的结构是正常入库的)
      if (req === '' || req === '{}' || req === null) {
        return;
      }

      // // 若有设置数据, 则取出 TIPS: 注意, 这里的设置是直接读出的一个json对象, 而不是需要解析的json字符串
      // const settingStorage = JSON.parse(req);

      const settingStorage = req;

      // 使用从存储取出的设置数据, 对setting-store.ts内的相关变量做初始化
      // TODO: 修改配置名或加入新配置后, 需在此处做相应的初始化处理 (代号 setting)
      if (settingStorage.language_default) {
        languageDefault.value = settingStorage.language_default;
      }
    });

    // realTimeStorageCore(实时存储核心), 用于将用户所做的设置, 实时监听式的存入底层数据库。
    // watchEffect(() => {
    //   const settingStorage = {
    //     // TODO: 修改配置名或加入新配置后, 需在此处做相应的持久化处理 (代号 setting)
    //     languageDefault: languageDefault.value,
    //   };
    //   StoreSet('settingPage', JSON.stringify(settingStorage));
    // });

    watch(languageDefault, () => {
      StoreSet('language_default', languageDefault.value);
    });
  }

  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- setting持久化 end   -_-^_^-_- ^_^-_-^_^-_-
  // ...
  // ...
  // ...
  //!endregion ----->>>>>>>>>>>>>>>>>>>> -- setting持久化 end   -_-^_^-_- ^_^-_-^_^-_-

  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/

  return {
    languageDefault,
    settingInitAndRealTimeStorage,
  };
});
