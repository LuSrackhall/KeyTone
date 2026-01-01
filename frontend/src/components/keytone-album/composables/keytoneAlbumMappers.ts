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
 * File: keytoneAlbumMappers.ts
 *
 * 目的（Why）
 * - 收敛“配置数据 → UI 结构”的纯映射逻辑，避免在多个位置重复写 Object.entries + sort。
 * - 让这些逻辑保持纯函数（输入 → 输出），便于单独审阅、复用与定位问题。
 *
 * 使用场景（Where used）
 * - `useKeytoneAlbumSseSync.ts`：处理 SSE 推送的全量配置更新。
 * - `Keytone_album.vue` 的 initData：初始化读取配置时复用相同映射，确保"初始化"和"SSE 更新"一致。
 * - `Keytone_album.vue` 的 watch(audioFiles)：当 audioFiles 变化时，复用 mapAudioFilesArrayToSoundFileList。
 *
 * 设计约束（Constraints）
 * - 不直接依赖 Vue 响应式：这里不操作 ref/reactive，仅返回新数组/Map。
 * - 不做业务裁剪：除非历史实现里有明确的裁剪规则（例如：单键声效必须 down/up 至少一个存在）。
 *
 * Debug 指南
 * - 映射结果不对：优先检查传入的 config（audio_files/sounds/key_sounds）结构是否与后端一致。
 * - 排序不对：检查调用方传入的 naturalSort 是否与原实现一致。
 */

export type NaturalSortFn = (a: string, b: string) => number;

export function mapAudioFilesConfigToArray(audioFilesConfig: any): Array<{ sha256: string; value: any }> {
  if (!audioFilesConfig) return [];
  return Object.entries(audioFilesConfig).map(([sha256, value]) => ({ sha256, value }));
}

export function mapAudioFilesArrayToSoundFileList(
  audioFilesArray: Array<{ sha256: string; value: any }>,
  naturalSort: NaturalSortFn
): Array<{ sha256: string; name_id: string; name: string; type: string }> {
  const tempSoundFileList: Array<{ sha256: string; name_id: string; name: string; type: string }> = [];

  audioFilesArray.forEach((item) => {
    if (item?.value?.name !== undefined && item?.value?.name !== null) {
      Object.entries(item.value.name).forEach(([name_id, name]) => {
        tempSoundFileList.push({
          sha256: item.sha256,
          name_id,
          name: name as string,
          type: item.value.type as string,
        });
      });
    }
  });

  tempSoundFileList.sort((a, b) => naturalSort(a.name + a.type, b.name + b.type));
  return tempSoundFileList;
}

export function mapSoundsConfigToList(
  soundsConfig: any,
  naturalSort: NaturalSortFn
): Array<{ soundKey: string; soundValue: any }> {
  if (!soundsConfig) return [];
  const sounds = Object.entries(soundsConfig).map(([soundKey, soundValue]) => ({ soundKey, soundValue }));

  sounds.sort((a: any, b: any) => {
    const aName = (a.soundValue?.name as string) || a.soundKey;
    const bName = (b.soundValue?.name as string) || b.soundKey;
    return naturalSort(aName, bName);
  });

  return sounds;
}

export function mapKeySoundsConfigToList(
  keySoundsConfig: any,
  naturalSort: NaturalSortFn
): Array<{ keySoundKey: string; keySoundValue: any }> {
  if (!keySoundsConfig) return [];
  const keySounds = Object.entries(keySoundsConfig).map(([keySoundKey, keySoundValue]) => ({
    keySoundKey,
    keySoundValue,
  }));

  keySounds.sort((a: any, b: any) => {
    const aName = (a.keySoundValue?.name as string) || a.keySoundKey;
    const bName = (b.keySoundValue?.name as string) || b.keySoundKey;
    return naturalSort(aName, bName);
  });

  return keySounds;
}

export function mapSingleKeyConfigToKeysWithSoundEffect(singleConfig: any): Map<string, any> {
  const keysWithSoundEffect = new Map<string, any>();
  if (!singleConfig) return keysWithSoundEffect;

  Object.entries(singleConfig).forEach(([dikCode, value]) => {
    if ((value as any)?.down?.value || (value as any)?.up?.value) {
      keysWithSoundEffect.set(dikCode, value);
    }
  });

  return keysWithSoundEffect;
}
