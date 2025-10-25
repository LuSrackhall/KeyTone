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
  <div :class="['app-container', isMacOS ? '' : 'app-container_shadow']">
    <div class="content">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import { useSettingStore } from 'src/stores/setting-store';
import { useAppStore } from './stores/app-store';
import { useSignatureStore } from './stores/signature-store';
import { debounce } from 'lodash';
import { useKeyEventStore } from './stores/keyEvent-store';
import { useMainStore } from './stores/main-store';
import { LoadConfig } from './boot/query/keytonePkg-query';

const app_store = useAppStore();
const setting_store = useSettingStore();
const signature_store = useSignatureStore();
const keyEvent_store = useKeyEventStore();
// 在此处调用, 只是为了提前初始化, 从而避免在主页面中, 出现初始化延迟所造成的 已选择的键音包 无法正常显示名字(即 显示空名字) 的问题。
// TIPS: 以上顾虑已通过将main_store内对应的map变量 keyTonePkgOptionsName 设置成 ref响应式变量来解决了, 不过为了加快速度, 仍在此处提前调用下, 而且除了更快的加载, 还起到一定的双重保险提高准确度的作用。
// TIPS: 调试时也可以观察到, 如果此处调用的话, 主页加载后, 已选择的键音包名称是直接显示的有。(否则, 也就是注释掉此调用后, 是可以观察的到 此名称由 无 到 显现 的闪烁过程的。)
const main_store = useMainStore();

onBeforeMount(async () => {
  // 首次加载时, 加载键音包列表
  main_store.GetKeyToneAlbumList();

  await setting_store.settingInitAndRealTimeStorage();

  // 在此处调用, 只是为了提前初始化sdk中用户所选键音包的加载。
  // * 放在setting_store.settingInitAndRealTimeStorage()后面(并对齐施以await), 是为了确保能够加载配置文件中持久化的用户所选的键音包。
  // * > 毕竟此函数内部所依赖的setting_store.mainHome.selectedKeyTonePkg 是由setting_store.settingInitAndRealTimeStorage()调用完成后才给予赋值的。
  // main_store.LoadSelectedKeyTonePkg();// TIPS: 由于函数内部, 还会依赖ConfigGet('audio_pkg_uuid')的返回值, 而首次加载时, sdk中没有任何键音包的加载, 因此会返回错误。
  // 为了防止出现以上TIPS中的报错, 我们首次加载无需判断 ConfigGet('audio_pkg_uuid')的返回值, 也也就是此处直接调用LoadConfig()函数即可。
  //
  // 如果用户所选键音包为空, 则没必要进行加载逻辑
  if (setting_store.mainHome.selectedKeyTonePkg) {
    LoadConfig(setting_store.mainHome.selectedKeyTonePkg, false).then((res) => {
      if (!res) {
        // 如果LoadConfig加载失败, 说明用户所选的键音包在当前环境下(可能已被外力删除), 因此我们将其置空。
        setting_store.mainHome.selectedKeyTonePkg = '';
        console.log(
          '在首次启动时, 加载持久化中用户所选的键音包失败, 此键音包可能已被破坏, 已清空所选键音包以供用户重新选择。'
        );
        return;
      }
      console.log('已在首次启动时, 成功加载持久化中用户所选的键音包');
    });
  }

  //#region    -----<<<<<<<<<<<<<<<<<<<< -- save setting start ^_^-_-^_^

  function sseDataToSettingStore(settingStorage: any) {
    setting_store.getConfigFileToUi();

    // 签名管理器数据同步处理
    signature_store.sseSync();
  }

  // 防抖处理, 避免短时间内多次触发时, 引起多次不必要的赋值操作, 我们的sse回调仅保证ui与配置文件的最终一致性即可。
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

  window.addEventListener('keydown', function (event) {
    // 如以此按下(不抬起) `ctrl`+`alt`+`任意字母键如j`, 会意外的触发出一个Dik为0, 名称为'Unidentified'的键
    // * 因此我们利用代码逻辑, 禁止此触发。
    const identified = event.code ? event.code : event.key;
    if (identified === 'Unidentified' || !identified) {
      return;
    }

    const frontendKeyUUID = identified + event.location;
    if (!keyEvent_store.frontendKeyEventStateBool.has(frontendKeyUUID)) {
      keyEvent_store.frontendKeyEventStateBool.set(frontendKeyUUID, true);
    }
    if (keyEvent_store.frontendKeyEventStateBool.get(frontendKeyUUID) === true) {
      // console.log('应用内按键', event.code ? event.code : event.key, '的触发状态是: ', event.type);

      // console.debug('录制的输入key:', event.key);
      // console.debug('录制的输入code:', event.code);
      // console.debug('录制的输入All:', event);
      // console.debug('按键的位置:', event.location);
      // TIPS: 在浏览器标准中keyCode已被弃用, 现今由code代替它。
      // console.debug('录制的输入keyCode:', event.keyCode);

      keyEvent_store.frontendKeyEventStateBool.set(frontendKeyUUID, false);
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
    // console.log('应用内按键', event.code ? event.code : event.key, '的触发状态是: ', event.type);

    // console.debug('录制的输入key:', event.key);
    // console.debug('录制的输入code:', event.code);
    // console.debug('录制的输入All:', event);
    // console.debug('按键的位置:', event.location);
    // TIPS: 在浏览器标准中keyCode已被弃用, 现今由code代替它。
    // console.debug('录制的输入keyCode:', event.keyCode);

    keyEvent_store.frontendKeyEventStateBool.set(frontendKeyUUID, true);
  });

  //#endregion ----->>>>>>>>>>>>>>>>>>>> -- keyEvent end   -_-^_^-_- ^_^-_-^_^-_-
});

const isMacOS = ref(getMacOSStatus());
function getMacOSStatus() {
  if (process.env.MODE === 'electron') {
    return window.myWindowAPI.getMacOSStatus();
  }
  return false;
}
</script>

<style lang="scss" scoped>
:global(.q-dialog__backdrop) {
  // background-color: rgba(0, 0, 0, 0) !important;
  @apply bg-transparent;
}

:global(.body) {
  margin: 0;
  height: 500px;
  /* overflow: hidden; */
  /* background: #edc0bf; */
  /* background: rgba(255, 255, 255, 0.2); 半透明背景 */
  /* background: linear-gradient(90deg, #edc0bf 0,#c4caef 58%); */
  /* font-family: 'Inter', sans-serif; */
}

.app-container {
  width: 100%;
  height: 100%;

  /* 将背景的设置从body 转移到此处, 以避免对实现圆角的功能产生影响*/
  /* background: #edc0bf; */
  background: linear-gradient(90deg, #edc0bf 0%, #c4caef 58%);
  /* border-radius: 30px; 圆角<这里设置的远大于实际的即可> */
  /* border-radius: 1.0rem; */
  /* 顶部圆角设置 */
  border-top-left-radius: 0.5rem /* 8px */;
  border-top-right-radius: 0.5rem /* 8px */;
  /* 底部圆角设置 */
  border-bottom-right-radius: 0.25rem /* 4px */;
  border-bottom-left-radius: 0.25rem /* 4px */;
}

.app-container_shadow {
  /* 更好且更方便的做法是, 直接在electron-main.ts中的原生级别解决此问题*/
  /* 留出边距以展示阴影 */
  width: calc(100% - 10px);
  height: calc(100% - 10px);

  /* 阴影设置 */
  /* box-shadow: 0 0 5px rgba(0, 0, 0, 0.5); */
  box-shadow: 2px 4px 8px rgba(0, 0, 0, 0.2);
  /* box-shadow: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1),0 20px 25px -5px var(--tw-shadow-color), 0 8px 10px -6px var(--tw-shadow-color) ; */
}

.content {
  position: relative;
  z-index: 1;
  padding: 0px;
  color: black;
}
</style>
