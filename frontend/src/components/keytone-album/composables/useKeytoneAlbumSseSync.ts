/**
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
 **/
/**
 * File: useKeytoneAlbumSseSync.ts
 *
 * 目的（Why）
 * - 将 `Keytone_album.vue` 中与 SSE（EventSource）相关的“监听 + 数据映射 + 写入状态”的逻辑抽离为 composable。
 * - 让父组件回归为“状态承载 + 生命周期挂载/卸载”，降低 `Keytone_album.vue` 的认知负担。
 *
 * 职责边界（What / What not）
 * - ✅ 负责：
 *   - 监听 `messageAudioPackage` 事件
 *   - 解析 JSON，筛选 `data.key === 'get_all_value'`
 *   - 将后端配置数据映射写入到父组件持有的 refs/reactive
 *   - 保持原先的 debounce 行为与延迟写入策略（pkgName、embedded test sound）
 * - ❌ 不负责：
 *   - UI 渲染
 *   - ConfigSet/ConfigDelete 写回（写回仍由父组件/其他函数负责）
 *
 * 与哪些文件配合（Integration）
 * - 调用方：`frontend/src/components/Keytone_album.vue`
 *   - 在 onBeforeMount 内调用 `attach(app_store.eventSource)`
 *   - 在 onUnmounted 内调用 attach 返回的 detach 函数
 * - 依赖：`./keytoneAlbumMappers`
 *   - `mapAudioFilesConfigToArray`
 *   - `mapSoundsConfigToList`
 *   - `mapKeySoundsConfigToList`
 *   - `mapSingleKeyConfigToKeysWithSoundEffect`
 * - 依赖：父组件提供的 `convertValue`/`naturalSort`
 *
 * 行为不变约束（Behavior parity）
 * - 重要：这里刻意保留了历史实现中的一个细节：原实现写的是 `debounced.cancel;`（未调用）。
 *   为避免引入“更激进的 cancel”导致时序改变，这里同样不主动 cancel 30ms 的 debounce。
 * - 其它映射与赋值路径保持与原先父组件内实现一致。
 *
 * Debug 指南（Where to debug）
 * - 如果发现 UI 未更新：优先检查 eventSource 是否触发 `messageAudioPackage`。
 * - 如果列表顺序异常：检查 `naturalSort` 的输入与 mapper 输出。
 * - 如果全局/单键联动未刷新：检查 `convertValue` 的映射是否能从最新列表中 find 到目标项。
 */

import { debounce } from 'lodash';
import type { Ref } from 'vue';

import {
  mapAudioFilesConfigToArray,
  mapKeySoundsConfigToList,
  mapSingleKeyConfigToKeysWithSoundEffect,
  mapSoundsConfigToList,
  type NaturalSortFn,
} from './keytoneAlbumMappers';

type EmbeddedTestSound = { down: boolean; up: boolean };

export function useKeytoneAlbumSseSync(params: {
  pkgName: Ref<string>;
  audioFiles: Ref<Array<any>>;
  soundList: Ref<Array<any>>;
  keySoundList: Ref<Array<any>>;
  isEnableEmbeddedTestSound: EmbeddedTestSound;
  keyDownUnifiedSoundEffectSelect: Ref<any>;
  keyUpUnifiedSoundEffectSelect: Ref<any>;
  keysWithSoundEffect: Ref<Map<string, any>>;
  convertValue: (item: any) => any;
  naturalSort: NaturalSortFn;
}) {
  let messageAudioPackageListener: ((e: MessageEvent) => void) | undefined;

  const pkgNameDelayed = debounce(
    (keyTonePkgData: any) => {
      params.pkgName.value = keyTonePkgData.package_name;
    },
    800,
    { trailing: true }
  );

  const isEnableEmbeddedTestSoundDelayed = debounce(
    (val: EmbeddedTestSound) => {
      params.isEnableEmbeddedTestSound.down = val.down;
      params.isEnableEmbeddedTestSound.up = val.up;
    },
    800,
    { trailing: true }
  );

  function sseDataToKeyTonePkgData(keyTonePkgData: any) {
    if (keyTonePkgData.package_name !== undefined) {
      pkgNameDelayed.cancel();
      pkgNameDelayed(keyTonePkgData);
    }

    if (keyTonePkgData.audio_files !== undefined) {
      params.audioFiles.value = mapAudioFilesConfigToArray(keyTonePkgData.audio_files);
    } else {
      params.audioFiles.value = [];
    }

    if (keyTonePkgData.sounds !== undefined) {
      params.soundList.value = mapSoundsConfigToList(keyTonePkgData.sounds, params.naturalSort);
    } else {
      params.soundList.value = [];
    }

    if (keyTonePkgData.key_sounds !== undefined) {
      params.keySoundList.value = mapKeySoundsConfigToList(keyTonePkgData.key_sounds, params.naturalSort);
    } else {
      params.keySoundList.value = [];
    }

    if (keyTonePkgData.key_tone !== undefined) {
      isEnableEmbeddedTestSoundDelayed.cancel();
      isEnableEmbeddedTestSoundDelayed(keyTonePkgData.key_tone.is_enable_embedded_test_sound);
    }

    if (keyTonePkgData.key_tone?.global !== undefined) {
      params.keyDownUnifiedSoundEffectSelect.value = params.convertValue(
        keyTonePkgData.key_tone.global.down ? keyTonePkgData.key_tone.global.down : ''
      );
      params.keyUpUnifiedSoundEffectSelect.value = params.convertValue(
        keyTonePkgData.key_tone.global.up ? keyTonePkgData.key_tone.global.up : ''
      );
    }

    if (keyTonePkgData.key_tone?.single !== undefined) {
      params.keysWithSoundEffect.value = mapSingleKeyConfigToKeysWithSoundEffect(keyTonePkgData.key_tone.single);
    }
  }

  const debouncedSseDataToSettingStore = debounce<(keyTonePkgData: any) => void>(sseDataToKeyTonePkgData, 30, {
    trailing: true,
  });

  function attach(eventSource: EventSource) {
    messageAudioPackageListener = function (e) {
      console.debug('后端钩子函数中的值 = ', e.data);

      const data = JSON.parse(e.data as any);

      if (data.key === 'get_all_value') {
        // 注意：历史实现里这里写的是 `debounced.cancel;`（未调用），为保持行为一致，这里不主动 cancel。
        debouncedSseDataToSettingStore(data.value);
      }
    };

    eventSource.addEventListener('messageAudioPackage', messageAudioPackageListener, false);

    return () => {
      if (messageAudioPackageListener) {
        eventSource.removeEventListener('messageAudioPackage', messageAudioPackageListener);
      }
      messageAudioPackageListener = undefined;
    };
  }

  return {
    attach,
  };
}
