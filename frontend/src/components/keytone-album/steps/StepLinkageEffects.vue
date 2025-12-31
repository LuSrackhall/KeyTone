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
文件说明: steps/StepLinkageEffects.vue - 步骤4：联动声效
============================================================================

【文件作用】
本组件是键音专辑编辑器的第四个步骤（Step 4），负责：
1. 控制内嵌测试音的开关（down/up）
2. 配置全局按键声效（统一所有按键的声音）
3. 配置单键声效（为特定按键设置独立声音）
4. 管理已配置的单键声效列表

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(KEYTONE_ALBUM_CONTEXT_KEY, ctx)
        │
        ▼
  ┌─────────────────────────────────────────────────────┐
  │  q-stepper                                          │
  │  ├── StepLoadAudioFiles.vue       (Step 1)          │
  │  ├── StepDefineSounds.vue         (Step 2)          │
  │  ├── StepCraftKeySounds.vue       (Step 3)          │
  │  └── StepLinkageEffects.vue  <── 当前文件 (Step 4)  │
  └─────────────────────────────────────────────────────┘

【数据流】
本组件不持有任何状态，所有状态都通过 inject 从父组件获取：

  父组件状态                              本组件使用方式
  ───────────────────────────────────────────────────────────
  ctx.step                           -->
控制当前步骤 ctx.isEnableEmbeddedTestSound --> 内嵌测试音开关 ctx.showEveryKeyEffectDialog --> 全键声效对话框
ctx.showSingleKeyEffectDialog --> 单键声效对话框 ctx.keysWithSoundEffect --> 已配置声效的按键 Map
ctx.saveUnifiedSoundEffectConfig() --> 保存全局声效 ctx.saveSingleKeySoundEffectConfig() --> 保存单键声效
【联动声效概念】 联动声效分为两种： 1. 全键声效：统一设置所有按键的声音（优先级低于单键） 2.
单键声效：为特定按键设置独立声音（优先级高于全键） 内嵌测试音：指在编辑器内测试按键时播放的声音，可分别控制
down（按下）与 up（抬起）。 【关联文件】 - ../types.ts : 类型定义 - ../../Keytone_album.vue : 父组件（提供 Context） -
../dialogs/EveryKeyEffectDialog.vue : 全键声效对话框 - ../dialogs/SingleKeyEffectDialog.vue : 单键声效对话框
【当前状态】 ✅ Step4 已从父组件迁移到本组件，父组件以 `
<StepLinkageEffects />
` 替换原有模板，实现“父组件持有状态 + 子组件承载 UI”。
============================================================================ -->

<template>
  <q-step
    :name="4"
    :title="ctx.$t('KeyToneAlbum.linkageEffects.title')"
    icon="settings"
    :done="!isDefaultState"
    :disable="ctx.step.value === 99 && isDefaultState"
    :header-nav="false"
    @click="handleStepClick"
  >
    <div :class="['mb-3', ctx.step_introduce_fontSize.value]">
      {{ ctx.$t('KeyToneAlbum.linkageEffects.description') }}
      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
        <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
          <span>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.description') }}</span>
        </q-tooltip>
      </q-icon>
    </div>

    <div :class="['flex items-center m-t-2 w-[130%]']">
      <span class="text-gray-500 mr-0.7">•</span>
      <span class="text-nowrap">
        {{ ctx.$t('KeyToneAlbum.linkageEffects.enableTestSound') }}:
        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
          <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
            <span>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.testSound') }}</span>
          </q-tooltip>
        </q-icon>
      </span>
    </div>
    <div
      :class="[
        'flex items-center ml-3',
        setting_store.languageDefault === 'pt' || setting_store.languageDefault === 'pt-BR'
          ? 'flex-nowrap text-nowrap'
          : '',
      ]"
    >
      <span class="text-gray-500 mr-1.5">•</span>
      <q-toggle
        v-model="ctx.isEnableEmbeddedTestSound.down"
        color="primary"
        :label="ctx.$t('KeyToneAlbum.linkageEffects.downTestSound')"
        dense
      />
    </div>
    <div :class="['flex items-center ml-3', setting_store.languageDefault === 'fr' ? 'flex-nowrap text-nowrap' : '']">
      <span class="text-gray-500 mr-1.5">•</span>
      <q-toggle
        v-model="ctx.isEnableEmbeddedTestSound.up"
        color="primary"
        :label="ctx.$t('KeyToneAlbum.linkageEffects.upTestSound')"
        dense
      />
    </div>

    <q-stepper-navigation>
      <div>
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.linkageEffects.globalSettings')"
          @click="() => (ctx.showEveryKeyEffectDialog.value = true)"
        >
        </q-btn>
        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
          <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
            <span>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.globalPriority') }}</span>
          </q-tooltip>
        </q-icon>
        <!-- 全键声效设置对话框（独立组件，内部通过 inject 获取 Context） -->
        <EveryKeyEffectDialog />
      </div>
      <div :class="['p-2 text-zinc-600']">{{ ctx.$t('KeyToneAlbum.or') }}</div>
      <div>
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.linkageEffects.singleKeySettings')"
          @click="() => (ctx.showSingleKeyEffectDialog.value = true)"
        >
        </q-btn>
        <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
          <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
            <span>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltips.singleKeyPriority') }}</span>
          </q-tooltip>
        </q-icon>
        <!-- 单键声效设置对话框（独立组件，内部通过 inject 获取 Context） -->
        <SingleKeyEffectDialog />
      </div>
    </q-stepper-navigation>

    <q-stepper-navigation>
      <q-btn @click="ctx.step.value = 5" color="primary" :label="ctx.$t('KeyToneAlbum.continue')" />
      <q-btn flat @click="ctx.step.value = 3" color="primary" :label="ctx.$t('KeyToneAlbum.back')" class="q-ml-sm" />
    </q-stepper-navigation>
  </q-step>
</template>

<script setup lang="ts">
/**
 * StepLinkageEffects.vue - 步骤4：联动声效
 *
 * 【组件职责】
 * - 只负责 UI 渲染和用户交互
 * - 所有状态和业务逻辑通过 inject 从父组件获取
 * - 不持有任何本地状态（除了纯 UI 状态）
 * - 作为 Step4 的 UI 承载：包含“内嵌测试音开关 + 全键/单键对话框入口 + Continue/Back 导航”
 *
 * 【内嵌测试音】
 * 控制在编辑器内测试按键时是否播放声音：
 * - down: 按键按下时
 * - up: 按键抬起时
 * 这个开关直接映射到配置文件，实时生效。
 */

import { inject, computed } from 'vue';
import { useSettingStore } from 'src/stores/setting-store';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import EveryKeyEffectDialog from '../dialogs/EveryKeyEffectDialog.vue';
import SingleKeyEffectDialog from '../dialogs/SingleKeyEffectDialog.vue';

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

const setting_store = useSettingStore();

// ============================================================================
// 计算属性
// ============================================================================
/**
 * 是否处于“默认/未配置”状态
 * - 该逻辑必须与父组件原 Step4 模板保持一致（用于 done/disable）
 */
const isDefaultState = computed(() => {
  return (
    ctx.isEnableEmbeddedTestSound.down === true &&
    ctx.isEnableEmbeddedTestSound.up === true &&
    !ctx.keyDownUnifiedSoundEffectSelect.value &&
    !ctx.keyUpUnifiedSoundEffectSelect.value &&
    ctx.keysWithSoundEffect.value.size === 0
  );
});

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 处理 step header 点击事件
 */
function handleStepClick(event: MouseEvent) {
  const header = (event.target as HTMLElement).closest('.q-stepper__tab');
  if (header) {
    ctx.step.value = ctx.step.value === 4 ? 99 : 4;
  }
}

// Step4 的折叠/展开由 handleStepClick 控制；继续/返回由模板中的 Continue/Back 按钮控制。
</script>

<style lang="scss" scoped>
/**
 * StepLinkageEffects 组件样式
 *
 * 【样式说明】
 * 本组件使用的样式大部分继承自父组件 Keytone_album.vue 的全局样式。
 * 此处仅定义本组件特有的样式。
 */

// 按钮样式 - 统一按钮外观
.q-btn {
  @apply text-xs;
  font-size: var(--i18n_fontSize);
  @apply p-1.5;
  @apply transition-transform hover:scale-105;
  @apply scale-103;
}
</style>
