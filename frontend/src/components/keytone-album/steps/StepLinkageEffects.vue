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
单键声效：为特定按键设置独立声音（优先级高于全键） 内嵌测试音是指在编辑器内测试按键时播放的声音， 可以单独控制
down（按下）和 up（抬起）的开关。 【关联文件】 - ../types.ts : 类型定义 - ../../Keytone_album.vue : 父组件 -
../../DependencyWarning.vue : 依赖警告组件 【当前状态】 ⚠️ 注意：本组件框架已创建，但 Step4 仍保留在父组件中！ Step4
内容复杂（约1500行），包含虚拟键盘和多个嵌套对话框，暂不替换。
============================================================================ -->

<template>
  <q-step
    :name="4"
    :title="ctx.$t('KeyToneAlbum.linkageEffects.title')"
    icon="settings"
    :done="hasAnyEffect"
    :disable="ctx.step.value === 99 && !hasAnyEffect"
    :header-nav="false"
    @click="handleStepClick"
  >
    <!-- 步骤说明 -->
    <div :class="['mb-3', ctx.step_introduce_fontSize.value]">
      {{ ctx.$t('KeyToneAlbum.linkageEffects.description') }}
      <q-icon name="info" color="primary">
        <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
          <div>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltip.linkageExplain') }}</div>
          <div>{{ ctx.$t('KeyToneAlbum.linkageEffects.tooltip.priorityExplain') }}</div>
        </q-tooltip>
      </q-icon>
    </div>

    <!-- 联动声效的业务逻辑 -->
    <div>
      <!-- 内嵌测试音开关 -->
      <div class="mb-4">
        <div class="text-sm text-gray-600 mb-2">
          {{ ctx.$t('KeyToneAlbum.linkageEffects.embeddedTestSound') }}
          <q-icon name="info" color="primary">
            <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
              {{ ctx.$t('KeyToneAlbum.linkageEffects.tooltip.embeddedTestSound') }}
            </q-tooltip>
          </q-icon>
        </div>
        <div class="flex gap-4">
          <q-toggle
            v-model="ctx.isEnableEmbeddedTestSound.down"
            :label="ctx.$t('KeyToneAlbum.linkageEffects.downSound')"
          />
          <q-toggle v-model="ctx.isEnableEmbeddedTestSound.up" :label="ctx.$t('KeyToneAlbum.linkageEffects.upSound')" />
        </div>
      </div>

      <!-- 全键声效设置 -->
      <div class="mb-3">
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.linkageEffects.everyKeyEffect')"
          @click="ctx.showEveryKeyEffectDialog.value = !ctx.showEveryKeyEffectDialog.value"
        />

        <!-- TODO: 全键声效对话框 - 待抽离为独立组件 dialogs/EveryKeyEffectDialog.vue -->
      </div>

      <div :class="['p-2 text-zinc-600']">{{ ctx.$t('KeyToneAlbum.or') }}</div>

      <!-- 单键声效设置 -->
      <div>
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.linkageEffects.singleKeyEffect')"
          @click="ctx.showSingleKeyEffectDialog.value = !ctx.showSingleKeyEffectDialog.value"
        />

        <!-- TODO: 单键声效对话框 - 待抽离为独立组件 dialogs/SingleKeyEffectDialog.vue -->
      </div>

      <!-- 已配置的单键声效列表（如果有） -->
      <div v-if="ctx.keysWithSoundEffect.value.size > 0" class="mt-4">
        <div class="text-sm text-gray-600 mb-2">
          {{ ctx.$t('KeyToneAlbum.linkageEffects.configuredKeys') }}
          ({{ ctx.keysWithSoundEffect.value.size }})
        </div>
        <!-- 显示已配置的按键列表 - 可点击编辑 -->
      </div>
    </div>

    <!-- 导航按钮 -->
    <q-stepper-navigation>
      <q-btn color="primary" :label="ctx.$t('KeyToneAlbum.finish')" @click="handleFinish" />
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
 *
 * 【完成条件】
 * 这是最后一个步骤，用户可以：
 * 1. 点击"完成"按钮结束编辑
 * 2. 或继续返回修改之前的步骤
 *
 * 【内嵌测试音】
 * 控制在编辑器内测试按键时是否播放声音：
 * - down: 按键按下时
 * - up: 按键抬起时
 * 这个开关直接映射到配置文件，实时生效。
 */

import { inject, computed } from 'vue';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// ============================================================================
// 计算属性
// ============================================================================

/**
 * 是否已配置任何联动声效
 * 用于判断步骤的 done 状态
 */
const hasAnyEffect = computed(() => {
  // 检查全键声效或单键声效是否已配置
  const hasUnified = ctx.keyDownUnifiedSoundEffectSelect?.value || ctx.keyUpUnifiedSoundEffectSelect?.value;
  const hasSingleKey = ctx.keysWithSoundEffect.value.size > 0;
  return hasUnified || hasSingleKey;
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

/**
 * 处理"完成"按钮点击
 *
 * 【行为说明】
 * 完成编辑后，折叠所有步骤（step = 99）
 * 用户可以随时点击任意步骤的 header 重新展开编辑
 */
function handleFinish() {
  // 完成编辑，折叠所有步骤
  ctx.step.value = 99;
}
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
