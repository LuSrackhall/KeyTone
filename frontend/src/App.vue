<template>
  <router-view />
</template>

<script setup lang="ts">
import { onBeforeMount } from 'vue';
import { useSettingStore } from 'src/stores/setting-store';
import { useAppStore } from './stores/app-store';
import { debounce } from 'lodash';

const app_store = useAppStore();
const setting_store = useSettingStore();

onBeforeMount(async () => {
  setting_store.settingInitAndRealTimeStorage();
  //#region    -----<<<<<<<<<<<<<<<<<<<< -- save setting start ^_^-_-^_^

  function sseDataToSettingStore(settingStorage: any) {
    if (settingStorage.language_default !== undefined) {
      setting_store.languageDefault = settingStorage.language_default;
    }

    // 手动打开应用时的默认设置
    // TIPS: 因为值本身就是boolean类型, 因此不能直接用于判断(最常见的做法时通过判断undefined来实现<因为当对象中不存在某个字段时, 会返回undefined>)。
    //       *  if (typeof settingStorage.startup.is_hide_windows === 'boolean') 虽然这样判断更准确, 但不够通用。 因为我只想简化开发成本, 所以我不用。
    if (settingStorage.startup.is_hide_windows !== undefined) {
      setting_store.startup.isHideWindows = settingStorage.startup.is_hide_windows;
    }

    // 自动启动应用时的默认设置
    if (settingStorage.auto_startup.is_auto_run !== undefined) {
      setting_store.autoStartup.isAutoRun = settingStorage.auto_startup.is_auto_run;
    }

    if (settingStorage.auto_startup.is_hide_windows !== undefined) {
      setting_store.autoStartup.isHideWindows = settingStorage.auto_startup.is_hide_windows;
    }

    // 音频音量处理的默认设置
    // * 用于设置页面 音量提升/缩减 设置
    if (settingStorage.audio_volume_processing.volume_amplify !== undefined) {
      setting_store.audioVolumeProcessing.volumeAmplify = settingStorage.audio_volume_processing.volume_amplify;
    }
    if (settingStorage.audio_volume_processing.volume_amplify_limit !== undefined) {
      setting_store.audioVolumeProcessing.volumeAmplifyLimit =
        settingStorage.audio_volume_processing.volume_amplify_limit;
    }

    // 主页面的默认设置
    if (settingStorage.main_home.audio_volume_processing.volume_normal !== undefined) {
      setting_store.mainHome.audioVolumeProcessing.volumeNormal =
        settingStorage.main_home.audio_volume_processing.volume_normal;
    }
    if (settingStorage.main_home.audio_volume_processing.volume_normal_reduce_scope !== undefined) {
      setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope =
        settingStorage.main_home.audio_volume_processing.volume_normal_reduce_scope;
    }
    if (settingStorage.main_home.audio_volume_processing.volume_silent !== undefined) {
      setting_store.mainHome.audioVolumeProcessing.volumeSilent =
        settingStorage.main_home.audio_volume_processing.volume_silent;
    }
    if (settingStorage.main_home.audio_volume_processing.is_open_volume_debug_slider !== undefined) {
      setting_store.mainHome.audioVolumeProcessing.isOpenVolumeDebugSlider =
        settingStorage.main_home.audio_volume_processing.is_open_volume_debug_slider;
    }
  }
  const debounced_sseDataToSettingStore = debounce<(settingStorage: any) => void>(sseDataToSettingStore, 300, {
    trailing: true,
  });
  app_store.eventSource.addEventListener(
    'message',
    function (e) {
      console.debug('后端钩子函数中AfterDelete中的值 = ', e.data);

      const data = JSON.parse(e.data);

      if (data.key === 'get_all_value') {
        // const settingStorage = data.value;
        // TODO: 修改配置名或加入新配置后, 需在此处做相应的初始化处理 (代号 setting)
        //        * TIPS: 如果涉及到数组或对象类的 配置项, 请不要直接赋值(即使内容一致, 也会引起最外层引用的变更, 从而造成循环触发)。
        //                我们可以使用 JSON序列化进行对比后再决定是否赋值,或直接使用可以保留最外层引用的deepAssign进行深度拷贝式赋值来解决。(这里我选择后者)
        // deepAssign(setting_store.bottomNavigationTabs , settingStorage.bottomNavigationTabs)
        debounced_sseDataToSettingStore.cancel;
        debounced_sseDataToSettingStore(data.value);
      }
    },
    false
  );
  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- save setting end   -_-^_^-_- ^_^-_-^_^-_-
  // ...
  // ...
  // ...
  //!endregion ----->>>>>>>>>>>>>>>>>>>> -- save setting end   -_-^_^-_- ^_^-_-^_^-_-
});
</script>
