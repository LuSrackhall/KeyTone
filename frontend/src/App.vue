<!--
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<template>
  <router-view />
</template>

<script setup lang="ts">
import { onBeforeMount } from 'vue';
import { useSettingStore } from 'src/stores/setting-store';
import { useAppStore } from './stores/app-store';
import { debounce } from 'lodash';
import { useKeyEventStore } from './stores/keyEvent-store';

const app_store = useAppStore();
const setting_store = useSettingStore();
const keyEvent_store = useKeyEventStore();

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
  const debounced_sseDataToSettingStore = debounce<(settingStorage: any) => void>(sseDataToSettingStore, 30, {
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

  //#region    -----<<<<<<<<<<<<<<<<<<<< -- keyEvent start ^_^-_-^_^

  app_store.eventSource.addEventListener('messageKeyEvent', function (e) {
    // console.group('[Debug] 键盘事件SSE消息处理');
    // console.debug('后端钩子函数中的值 = ', e.data);
    try {
      const data = JSON.parse(e.data);
      // console.debug('后端钩子函数中的值(解析后) = ', data);

      keyEvent_store.keyCodeState.set(data.keycode, data.state);

      // console.group('[Debug] 键盘事件状态更新');
      // console.debug('keycode为', data.keycode, '的按键的当前状态 ->  ', data.state);
      // console.debug('keyEvent_store.keyCodeState = ', keyEvent_store.keyCodeState);
      // console.groupEnd();
    } catch (err) {
      console.error('键盘事件处理失败:', {
        error: err,
        rawData: e.data,
      });
    } finally {
      // console.groupEnd();
    }
  });

  const frontendKeyEventStateBool = new Map<string, boolean>();
  window.addEventListener('keydown', function (event) {
    // 如以此按下(不抬起) `ctrl`+`alt`+`任意字母键如j`, 会意外的触发出一个Dik为0, 名称为'Unidentified'的键
    // * 因此我们利用代码逻辑, 禁止此触发。
    const identified = event.code ? event.code : event.key;
    if (identified === 'Unidentified' || !identified) {
      return;
    }

    const frontendKeyUUID = identified + event.location;
    if (!frontendKeyEventStateBool.has(frontendKeyUUID)) {
      frontendKeyEventStateBool.set(frontendKeyUUID, true);
    }
    if (frontendKeyEventStateBool.get(frontendKeyUUID) === true) {
      // console.log('按键', event.code ? event.code : event.key, '的触发状态是: ', event.type);
      // console.debug('录制的输入key:', event.key);
      // console.debug('录制的输入code:', event.code);
      // console.debug('录制的输入All:', event);
      // console.debug('按键的位置:', event.location);
      // TIPS: 在浏览器标准中keyCode已被弃用, 现今由code代替它。
      // console.debug('录制的输入keyCode:', event.keyCode);

      // 为保证时间差, 因此在down状态下, 便对齐进行赋值
      // * 进一步的, 为防止在识别按键名称时, 出现错误映射, 我们仅在字符串为空时, 才进行赋值。
      //   > 否则, 以 'alt + j' 为例(前提是在系统中 主动 被快捷键映射工具 映射为 方向键 'down' 的情况下), 会错误的将 'alt'对应的Dik码, 映射为 'down' 的名称。
      if (!keyEvent_store.keyDownStateName_ui) {
        keyEvent_store.keyDownStateName_ui = identified;
      }

      frontendKeyEventStateBool.set(frontendKeyUUID, false);
    }
  });
  window.addEventListener('keyup', function (event) {
    // 如以此按下(不抬起) `ctrl`+`alt`+`任意字母键如j`, 会意外的触发出一个Dik为0, 名称为'Unidentified'的键
    // * 因此我们利用代码逻辑, 禁止此触发。
    const identified = event.code ? event.code : event.key;
    if (identified === 'Unidentified' || !identified) {
      return;
    }
    const frontendKeyUUID = identified + event.location;
    // console.log('按键', event.code ? event.code : event.key, '的触发状态是: ', event.type);
    // console.debug('录制的输入key:', event.key);
    // console.debug('录制的输入code:', event.code);
    // console.debug('录制的输入All:', event);
    // console.debug('按键的位置:', event.location);
    // TIPS: 在浏览器标准中keyCode已被弃用, 现今由code代替它。
    // console.debug('录制的输入keyCode:', event.keyCode);

    frontendKeyEventStateBool.set(frontendKeyUUID, true);
  });

  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- keyEvent end   -_-^_^-_- ^_^-_-^_^-_-
});
</script>
