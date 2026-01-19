import { defineStore } from 'pinia';
import { Quasar, useQuasar } from 'quasar';
import { StoreGet, StoreSet } from 'src/boot/query/store-query';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import messages from 'src/i18n';
import { useMainStore } from './main-store';

export const useSettingStore = defineStore('setting', () => {
  /*------------------------------------------------------------------------------------------------------------------*/
  const q = useQuasar();
  const { t } = useI18n();
  const main_store = useMainStore();

  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/

  //#region    -----<<<<<<<<<<<<<<<<<<<< -- 页面状态管理(无需持久化) start ^_^-_-^_^

  const settingItemsOpenedState = ref([]);
  // watch(settingItemsOpenedState.value, () => {
  //   console.log('settingItemsOpenedState.value', settingItemsOpenedState.value);
  // });

  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- 页面状态管理(无需持久化) end   -_-^_^-_- ^_^-_-^_^-_-
  // ...
  // ...
  // ...
  //!endregion ----->>>>>>>>>>>>>>>>>>>> -- 页面状态管理(无需持久化) end   -_-^_^-_- ^_^-_-^_^-_-

  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/

  //#region    -----<<<<<<<<<<<<<<<<<<<< --languageDefault  start ^_^-_-^_^
  /**
   * 此处为, 设置页面定义语言项国际化页面 m 的区域
   * * 即 通过m对接全局相关vm的状态, 从而体现至v视图UI层
   */

  const language = Quasar.lang.getLocale();

  type LanguageDefault = keyof typeof messages;
  const languageDefault = ref<LanguageDefault>(
    language && (language as LanguageDefault) in messages ? (language as LanguageDefault) : 'en-US'
  );

  const { locale } = useI18n({ useScope: 'global' });

  // const modules = import.meta.glob('../../node_modules/quasar/lang/(de|en-US|es).js') // WARN: quasar官方提供的太麻烦, 还要自己去指定支持哪些 语言(如果遇到quasar语言包内不支持的简写, 会引起应用崩溃--这是我不希望的(毕竟只是语言问题))。(由于我对于此类限制, 可以随后在其它地方完成, 因此直接使用 通配符 * 最为方便。)
  const modules = import.meta.glob('../../node_modules/quasar/lang/*.js'); //   - completed(已完成)    根源上解决了 languageDefault 为 quasar内部存在, 但未写入联合类型时的报错问题。(但为了避免quasar内部不存在的类型(几乎不可能) 或是 quasar内部存在但我的i18n中没有配置的语言包被误设置, 故仍需要在languageDefault的类型上, 作出相关限制-这个我会在随后的提交中做处理。)

  /**
   * 此处采用立即执行, 是为了当数据库为空时, 也能自动赋值所获取的当地语言值(即使不适配也无妨), 不至于所获得值被刷而被动强制
   * 为英语。
   * * 即使当地语言我没有适配,i18n系统也会在对比不合适后, 反而执行启动文件中, 设置的默认值, 即强制en-us语言为系统语言。
   */
  watch(
    languageDefault,
    () => {
      languageDefault.value =
        languageDefault.value && (languageDefault.value as LanguageDefault) in messages
          ? (languageDefault.value as LanguageDefault)
          : 'en-US';

      locale.value = languageDefault.value;

      modules[`../../node_modules/quasar/lang/${languageDefault.value}.js`]().then((item) => {
        q.lang.set(item.default);
      });
      console.log('fff111=', Quasar.lang.getLocale());
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

  //#region    -----<<<<<<<<<<<<<<<<<<<< -- startup start ^_^-_-^_^

  // 个人习惯使用ref。 ( 不管是 ref 还是 reactive, 都不影响后续watch时, 对其子元素的单独监听。 )
  const startup = ref({
    isHideWindows: false,
  });
  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- startup end   -_-^_^-_- ^_^-_-^_^-_-
  // ...
  // ...
  // ...
  //!endregion ----->>>>>>>>>>>>>>>>>>>> -- startup end   -_-^_^-_- ^_^-_-^_^-_-

  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/

  //#region    -----<<<<<<<<<<<<<<<<<<<< -- auto startup start ^_^-_-^_^

  // 个人习惯使用ref。 ( 不管是 ref 还是 reactive, 都不影响后续watch时, 对其子元素的单独监听。 )
  const autoStartup = ref({
    isAutoRun: false,
    isHideWindows: true,
  });

  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- auto startup end   -_-^_^-_- ^_^-_-^_^-_-
  // ...
  // ...
  // ...
  //!endregion ----->>>>>>>>>>>>>>>>>>>> -- auto startup end   -_-^_^-_- ^_^-_-^_^-_-

  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/
  //#region    -----<<<<<<<<<<<<<<<<<<<< -- volume start ^_^-_-^_^
  const audioVolumeProcessing = ref({
    volumeAmplify: 0,
    volumeAmplifyLimit: 10,
  });
  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- volume end   -_-^_^-_- ^_^-_-^_^-_-
  // ...
  // ...
  // ...
  //!endregion ----->>>>>>>>>>>>>>>>>>>> -- volume end   -_-^_^-_- ^_^-_-^_^-_-
  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/
  //#region    -----<<<<<<<<<<<<<<<<<<<< -- mainHome start ^_^-_-^_^
  const mainHome = ref({
    audioVolumeProcessing: {
      volumeNormal: 0,
      volumeNormalReduceScope: 5,
      volumeSilent: false,
      isOpenVolumeDebugSlider: false,
    },
    selectedKeyTonePkg: '',
  });

  const playbackRouting = ref({
    // 播放路由模式：
    // - unified：键盘/鼠标共用一个播放专辑（默认，兼容旧行为）
    // - split：键盘/鼠标分别指定播放专辑
    mode: 'unified',
    // unifiedAlbumPath：统一模式下的“播放专辑”。
    // 说明：这里存的是“专辑路径或 UUID”，后端会做 resolve。
    unifiedAlbumPath: '',
    // keyboardAlbumPath / mouseAlbumPath：分离模式下的“播放专辑”。
    // 注意：这两个字段只影响主页日常播放，不影响编辑页（编辑页使用 editor 模式）。
    keyboardAlbumPath: '',
    mouseAlbumPath: '',
    // editorNoticeDismissed：编辑页提示是否“不再显示”。
    // 该提示以底部通知展示，不应影响页面布局；用户可手动关闭或选择不再提示。
    editorNoticeDismissed: false,
    // mouseFallbackToKeyboard：分离模式下，鼠标专辑缺失时是否回退到键盘专辑。
    // 默认 false：彻底分离，鼠标无专辑则无声（使用内嵌测试音）。
    // 用户可在设置页面开启此选项以复用键盘专辑处理鼠标事件。
    mouseFallbackToKeyboard: false,
  });

  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- mainHome end   -_-^_^-_- ^_^-_-^_^-_-
  // ...
  // ...
  // ...
  //!endregion ----->>>>>>>>>>>>>>>>>>>> -- mainHome end   -_-^_^-_- ^_^-_-^_^-_-

  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/
  //#region    -----<<<<<<<<<<<<<<<<<<<< -- 首次启动获取及同步获取配置文件内容至ui start ^_^-_-^_^
  async function getConfigFileToUi() {
    await StoreGet('get_all_value').then((req) => {
      // console.debug('打印观察获取的值', req);
      if (req === false) {
        // 此时, 说明GetItem_sqlite请求过程中, 出错了, 因此需要错误通知, 并让用户重新启动, 防止用户因继续使用造成的存储设置被初始覆盖
        q.notify({
          type: 'negative',
          position: 'top',
          message: t('Notify.配置文件读取失败,请重启应用'),
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

      // 语言设置(仅在前端配置的配置项,sdk对此无依赖)
      // 使用从存储取出的设置数据, 对setting-store.ts内的相关变量做初始化
      // TODO: 修改配置名或加入新配置后, 需在此处做相应的初始化处理 (代号 setting)
      // TIPS: 这里只是为了判断是否从配置文件中读到了这个内容。为防止内容本身就为bool类型, 最常见的做法时通过判断undefined来实现<因为当对象中不存在某个字段时, 会返回undefined>。
      if (settingStorage.language_default !== undefined) {
        languageDefault.value = settingStorage.language_default;
      } else {
        StoreSet('language_default', languageDefault.value);
      }

      // 手动打开应用时的默认设置
      // TIPS: 因为值本身就是boolean类型, 因此不能直接用于判断(最常见的做法时通过判断undefined来实现<因为当对象中不存在某个字段时, 会返回undefined>)。
      //       *  if (typeof settingStorage.startup.is_hide_windows === 'boolean') 虽然这样判断更准确, 但不够通用。 因为我只想简化开发成本, 所以我不用。
      if (settingStorage.startup.is_hide_windows !== undefined) {
        startup.value.isHideWindows = settingStorage.startup.is_hide_windows;
      }

      // 自动启动应用时的默认设置
      if (settingStorage.auto_startup.is_auto_run !== undefined) {
        autoStartup.value.isAutoRun = settingStorage.auto_startup.is_auto_run;
      }

      if (settingStorage.auto_startup.is_hide_windows !== undefined) {
        autoStartup.value.isHideWindows = settingStorage.auto_startup.is_hide_windows;
      }

      // 音频音量处理的默认设置
      // * 用于设置页面 音量提升/缩减 设置
      if (settingStorage.audio_volume_processing.volume_amplify !== undefined) {
        audioVolumeProcessing.value.volumeAmplify = settingStorage.audio_volume_processing.volume_amplify;
      }
      if (settingStorage.audio_volume_processing.volume_amplify_limit !== undefined) {
        audioVolumeProcessing.value.volumeAmplifyLimit = settingStorage.audio_volume_processing.volume_amplify_limit;
      }

      // 主页面的默认设置
      if (settingStorage.main_home.audio_volume_processing.volume_normal !== undefined) {
        mainHome.value.audioVolumeProcessing.volumeNormal =
          settingStorage.main_home.audio_volume_processing.volume_normal;
      }
      if (settingStorage.main_home.audio_volume_processing.volume_normal_reduce_scope !== undefined) {
        mainHome.value.audioVolumeProcessing.volumeNormalReduceScope =
          settingStorage.main_home.audio_volume_processing.volume_normal_reduce_scope;
      }
      if (settingStorage.main_home.audio_volume_processing.volume_silent !== undefined) {
        mainHome.value.audioVolumeProcessing.volumeSilent =
          settingStorage.main_home.audio_volume_processing.volume_silent;
      }
      if (settingStorage.main_home.audio_volume_processing.is_open_volume_debug_slider !== undefined) {
        mainHome.value.audioVolumeProcessing.isOpenVolumeDebugSlider =
          settingStorage.main_home.audio_volume_processing.is_open_volume_debug_slider;
      }
      if (settingStorage.main_home.selected_key_tone_pkg !== undefined) {
        mainHome.value.selectedKeyTonePkg = settingStorage.main_home.selected_key_tone_pkg;
        // 每次用户的主动选择, 都会在SSE中触发实际选择的键音包重新进行加载。
        // 也就是说, 不止在主页面中通过watch监听触发此函数, 在sse回调中也再次调用此函数, 以保证用户的选择能够最大程度上被可靠的加载。
        // 并且无需担心, 重复调用此函数也不会引发重复加载相同的键音包。
        main_store.LoadSelectedKeyTonePkg();
      }

      // legacySelected 用于迁移历史字段到新的路由字段。
      // 迁移目标：升级后默认保持旧行为（一个专辑同时用于键盘+鼠标播放）。
      const legacySelected = settingStorage.main_home?.selected_key_tone_pkg ?? '';
      const routingStorage = settingStorage.playback?.routing ?? {};

      if (routingStorage.mode !== undefined) {
        playbackRouting.value.mode = routingStorage.mode || 'unified';
      } else {
        // 初次迁移：默认 unified
        playbackRouting.value.mode = 'unified';
        StoreSet('playback.routing.mode', playbackRouting.value.mode);
      }

      if (routingStorage.unified_album_path !== undefined) {
        playbackRouting.value.unifiedAlbumPath = routingStorage.unified_album_path || '';
      } else {
        // 迁移：使用历史 selected_key_tone_pkg
        playbackRouting.value.unifiedAlbumPath = legacySelected;
        StoreSet('playback.routing.unified_album_path', playbackRouting.value.unifiedAlbumPath);
      }

      if (routingStorage.keyboard_album_path !== undefined) {
        playbackRouting.value.keyboardAlbumPath = routingStorage.keyboard_album_path || '';
      } else {
        playbackRouting.value.keyboardAlbumPath = legacySelected;
        StoreSet('playback.routing.keyboard_album_path', playbackRouting.value.keyboardAlbumPath);
      }

      if (routingStorage.mouse_album_path !== undefined) {
        playbackRouting.value.mouseAlbumPath = routingStorage.mouse_album_path || '';
      } else {
        playbackRouting.value.mouseAlbumPath = legacySelected;
        StoreSet('playback.routing.mouse_album_path', playbackRouting.value.mouseAlbumPath);
      }

      if (routingStorage.editor_notice_dismissed !== undefined) {
        playbackRouting.value.editorNoticeDismissed = !!routingStorage.editor_notice_dismissed;
      } else {
        // 新增字段：默认展示提示（false）。
        playbackRouting.value.editorNoticeDismissed = false;
        StoreSet('playback.routing.editor_notice_dismissed', playbackRouting.value.editorNoticeDismissed);
      }

      // mouseFallbackToKeyboard：分离模式下鼠标专辑缺失时的回退策略
      if (routingStorage.mouse_fallback_to_keyboard !== undefined) {
        playbackRouting.value.mouseFallbackToKeyboard = !!routingStorage.mouse_fallback_to_keyboard;
      } else {
        // 新增字段：默认不回退（彻底分离）
        playbackRouting.value.mouseFallbackToKeyboard = false;
        StoreSet('playback.routing.mouse_fallback_to_keyboard', playbackRouting.value.mouseFallbackToKeyboard);
      }
    });
  }

  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- 首次启动获取及同步获取配置文件内容至ui end -_-^_^-_- ^_^-_-^_^-_-
  // ...
  // ...
  // ...
  //!endregion ----->>>>>>>>>>>>>>>>>>>> -- 首次启动获取及同步获取配置文件内容至ui end -_-^_^-_- ^_^-_-^_^-_-
  /*------------------------------------------------------------------------------------------------------------------*/
  /*------------------------------------------------------------------------------------------------------------------*/
  //#region    -----<<<<<<<<<<<<<<<<<<<< -- setting持久化 start ^_^-_-^_^

  async function settingInitAndRealTimeStorage() {
    // 优先使用数据库中保存的设置, 即先通过数据库存储, 对内存做初始化
    await getConfigFileToUi();
    // realTimeStorageCore(实时存储核心), 用于将用户所做的设置, 实时监听式的存入底层数据库。
    // watchEffect(() => {
    //   const settingStorage = {
    //     // TODO: 修改配置名或加入新配置后, 需在此处做相应的持久化处理 (代号 setting)
    //     languageDefault: languageDefault.value,
    //   };
    //   StoreSet('settingPage', JSON.stringify(settingStorage));
    // });

    // 语言设置(仅在前端配置的配置项,sdk对此无依赖)
    watch(languageDefault, () => {
      StoreSet('language_default', languageDefault.value);
    });

    // 手动打开应用时的默认设置
    watch(
      () => startup.value.isHideWindows,
      () => {
        StoreSet('startup.is_hide_windows', startup.value.isHideWindows);
      }
    );

    // 自动启动应用时的默认设置
    watch(
      () => autoStartup.value.isAutoRun,
      () => {
        StoreSet('auto_startup.is_auto_run', autoStartup.value.isAutoRun);
      }
    );
    watch(
      () => autoStartup.value.isHideWindows,
      () => {
        StoreSet('auto_startup.is_hide_windows', autoStartup.value.isHideWindows);
      }
    );

    // 音频音量处理的默认设置
    // * 用于设置页面 音量提升/缩减 设置
    watch(
      () => audioVolumeProcessing.value.volumeAmplify,
      () => {
        StoreSet('audio_volume_processing.volume_amplify', audioVolumeProcessing.value.volumeAmplify);
      }
    );
    watch(
      () => audioVolumeProcessing.value.volumeAmplifyLimit,
      () => {
        if (
          audioVolumeProcessing.value.volumeAmplifyLimit > 0 &&
          audioVolumeProcessing.value.volumeAmplifyLimit < 100000000
        ) {
          StoreSet('audio_volume_processing.volume_amplify_limit', audioVolumeProcessing.value.volumeAmplifyLimit);
        }
      }
    );

    // 主页面的默认设置
    watch(
      () => mainHome.value.audioVolumeProcessing.volumeNormal,
      () => {
        StoreSet('main_home.audio_volume_processing.volume_normal', mainHome.value.audioVolumeProcessing.volumeNormal);
      }
    );
    watch(
      () => mainHome.value.audioVolumeProcessing.volumeNormalReduceScope,
      () => {
        if (
          mainHome.value.audioVolumeProcessing.volumeNormalReduceScope >= 5 &&
          mainHome.value.audioVolumeProcessing.volumeNormalReduceScope < 100000000
        ) {
          StoreSet(
            'main_home.audio_volume_processing.volume_normal_reduce_scope',
            mainHome.value.audioVolumeProcessing.volumeNormalReduceScope
          );
        }
      }
    );
    watch(
      () => mainHome.value.audioVolumeProcessing.volumeSilent,
      () => {
        StoreSet('main_home.audio_volume_processing.volume_silent', mainHome.value.audioVolumeProcessing.volumeSilent);
      }
    );

    watch(
      () => mainHome.value.audioVolumeProcessing.isOpenVolumeDebugSlider,
      () => {
        StoreSet(
          'main_home.audio_volume_processing.is_open_volume_debug_slider',
          mainHome.value.audioVolumeProcessing.isOpenVolumeDebugSlider
        );
      }
    );

    watch(
      () => mainHome.value.selectedKeyTonePkg,
      () => {
        StoreSet('main_home.selected_key_tone_pkg', mainHome.value.selectedKeyTonePkg);
      }
    );

    watch(
      () => playbackRouting.value.mode,
      () => {
        StoreSet('playback.routing.mode', playbackRouting.value.mode);
      }
    );

    watch(
      () => playbackRouting.value.unifiedAlbumPath,
      () => {
        StoreSet('playback.routing.unified_album_path', playbackRouting.value.unifiedAlbumPath);
      }
    );

    watch(
      () => playbackRouting.value.keyboardAlbumPath,
      () => {
        StoreSet('playback.routing.keyboard_album_path', playbackRouting.value.keyboardAlbumPath);
      }
    );

    watch(
      () => playbackRouting.value.mouseAlbumPath,
      () => {
        StoreSet('playback.routing.mouse_album_path', playbackRouting.value.mouseAlbumPath);
      }
    );

    watch(
      () => playbackRouting.value.editorNoticeDismissed,
      () => {
        StoreSet('playback.routing.editor_notice_dismissed', playbackRouting.value.editorNoticeDismissed);
      }
    );

    // 监听鼠标回退到键盘开关的变化，实时持久化
    watch(
      () => playbackRouting.value.mouseFallbackToKeyboard,
      () => {
        StoreSet('playback.routing.mouse_fallback_to_keyboard', playbackRouting.value.mouseFallbackToKeyboard);
      }
    );
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
    settingItemsOpenedState,
    startup,
    autoStartup,
    audioVolumeProcessing,
    mainHome,
    playbackRouting,
    getConfigFileToUi,
    settingInitAndRealTimeStorage,
  };
});
