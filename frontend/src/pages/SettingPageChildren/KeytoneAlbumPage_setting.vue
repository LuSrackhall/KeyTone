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

<!--
============================================================================
文件说明: KeytoneAlbumPage_setting.vue - 键音专辑页相关设置
============================================================================

【文件作用】
- 设置页面中"键音专辑页相关"分类的展开面板内容。
- 当前包含：波形滚动行为偏好（paged-jump / edge-push）。
- 后续可在此追加其他键音专辑页相关的配置项。

【UI 规范】
- 遵循项目设置子页面的统一布局模式：
  q-item 容器 + 左侧竖条装饰线 + 右侧控件。
- 不添加额外的 section 标题（父级 q-expansion-item 已有 label/caption）。
============================================================================
-->

<template>
  <!--
    波形滚动行为设置项：
    - 左侧竖条与缩进保持与 Language_setting / MainHome_setting 一致。
    - 使用 q-select（borderless + dense）与语言选择的风格接近。
  -->
  <q-item>
    <!-- 左侧装饰竖条（与其他设置子页面一致） -->
    <div :class="['ml-6 rounded-full border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
    <q-select
      v-model="setting_store.keytoneAlbumScrollBehavior"
      :options="scrollOptions"
      :label="$t('KeyTone.setting.category.scrollBehavior')"
      dense
      borderless
      emit-value
      map-options
      options-dense
      :class="['min-w-[160px] w-[100%]']"
    />
  </q-item>
</template>

<script setup lang="ts">
/**
 * KeytoneAlbumPage_setting.vue
 *
 * 键音专辑页的设置子面板。
 * 当前设置项：
 * - keytoneAlbumScrollBehavior：波形播放跟随滚动行为（paged-jump / edge-push）
 *
 * 持久化链条：
 * 1. 用户在此修改 setting_store.keytoneAlbumScrollBehavior ->
 * 2. setting-store.ts 中的 watch 将新值写入 StoreSet('keytoneAlbumScrollBehavior', val) ->
 * 3. 下次启动时 getConfigFileToUi() 从 StoreGet('get_all_value') 读取并恢复。
 */

import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useSettingStore } from 'src/stores/setting-store';

const { t } = useI18n();
const setting_store = useSettingStore();

/**
 * 下拉选项列表（响应式：随语言切换实时更新 label）
 */
const scrollOptions = computed(() => [
  { label: t('KeyTone.setting.category.scrollOptions.edge-push'), value: 'edge-push' },
  { label: t('KeyTone.setting.category.scrollOptions.paged-jump'), value: 'paged-jump' },
]);
</script>

<style scoped></style>
