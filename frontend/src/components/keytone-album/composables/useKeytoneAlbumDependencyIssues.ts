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
 * File: useKeytoneAlbumDependencyIssues.ts
 *
 * 目的（Why）
 * - 将 `Keytone_album.vue` 中“依赖校验（Dependency validation）”相关的 computed + watch 抽离为 composable。
 * - 使父组件只需要拿到 `dependencyIssues` 与 `checkItemDependencyIssues`，用于 UI 展示与标记。
 *
 * 输入/输出（API）
 * - 输入：soundFileList / soundList / keySoundList / keysWithSoundEffect（均来自父组件的响应式状态）
 * - 输出：
 *   - `dependencyIssues`: Ref<DependencyIssue[]> 供 UI 组件展示（例如 DependencyWarning）
 *   - `allDependencyIssues`: computed（便于调试/扩展）
 *   - `checkItemDependencyIssues`: 给列表项快速判定是否存在依赖问题
 *
 * 与哪些文件配合（Integration）
 * - 调用方：`frontend/src/components/Keytone_album.vue`
 * - 依赖：`src/utils/dependencyValidator`
 *
 * 关键实现说明（Important details）
 * - 历史原因：UI 选择器里 keySound 的 down/up value 可能被转换成“对象形态”（包含 soundKey/keySoundKey）。
 *   但依赖校验器期望的是“字符串 key”形态。
 *   因此这里对 keySoundList 做深拷贝，并把对象形态还原为字符串 key，再交给 validator。
 * - 重要：这里的深拷贝是为了避免污染 UI 状态，确保依赖校验是只读的。
 *
 * 行为不变约束（Behavior parity）
 * - `globalBinding` 目前仍为 undefined（与原实现一致），仅校验 singleKeyBindings。
 * - watch 采用 deep: true（与原实现一致）。
 */

import { computed, ref, watch, type Ref } from 'vue';

import {
  createDependencyValidator,
  hasItemDependencyIssues,
  type AudioFile,
  type DependencyIssue,
  type KeySound,
  type Sound,
} from 'src/utils/dependencyValidator';

type KeytoneKeySoundListItem = any;

type Params = {
  soundFileList: Ref<Array<any>>;
  soundList: Ref<Array<any>>;
  keySoundList: Ref<Array<KeytoneKeySoundListItem>>;
  keysWithSoundEffect: Ref<Map<string, any>>;
};

export function useKeytoneAlbumDependencyIssues(params: Params) {
  const dependencyIssues = ref<DependencyIssue[]>([]);

  const allDependencyIssues = computed(() => {
    const audioFiles = params.soundFileList.value as AudioFile[];
    const sounds = params.soundList.value as Sound[];

    const keySounds = params.keySoundList.value
      .map((keySound) => {
        const keySoundCopy = JSON.parse(JSON.stringify(keySound));

        keySoundCopy.keySoundValue.down.value = keySoundCopy.keySoundValue.down.value.map((item: any) => {
          if (item.type === 'sounds' && item.value && typeof item.value === 'object' && item.value.soundKey) {
            return { type: 'sounds', value: item.value.soundKey };
          }
          if (item.type === 'key_sounds' && item.value && typeof item.value === 'object' && item.value.keySoundKey) {
            return { type: 'key_sounds', value: item.value.keySoundKey };
          }
          return item;
        });

        keySoundCopy.keySoundValue.up.value = keySoundCopy.keySoundValue.up.value.map((item: any) => {
          if (item.type === 'sounds' && item.value && typeof item.value === 'object' && item.value.soundKey) {
            return { type: 'sounds', value: item.value.soundKey };
          }
          if (item.type === 'key_sounds' && item.value && typeof item.value === 'object' && item.value.keySoundKey) {
            return { type: 'key_sounds', value: item.value.keySoundKey };
          }
          return item;
        });

        return keySoundCopy;
      })
      .filter(Boolean) as KeySound[];

    if (audioFiles.length === 0 && sounds.length === 0 && keySounds.length === 0) {
      return [];
    }

    const validator = createDependencyValidator(audioFiles, sounds, keySounds);

    const globalBinding = undefined;

    const singleKeyBindings = params.keysWithSoundEffect.value.size > 0 ? params.keysWithSoundEffect.value : undefined;

    return validator.validateAllDependencies(globalBinding, singleKeyBindings);
  });

  watch(
    [params.soundFileList, params.soundList, params.keySoundList, params.keysWithSoundEffect],
    () => {
      dependencyIssues.value = allDependencyIssues.value;
    },
    { deep: true }
  );

  const checkItemDependencyIssues = (itemType: 'audio_files' | 'sounds' | 'key_sounds', itemId: string) => {
    return hasItemDependencyIssues(itemType, itemId, dependencyIssues.value);
  };

  return {
    dependencyIssues,
    allDependencyIssues,
    checkItemDependencyIssues,
  };
}
