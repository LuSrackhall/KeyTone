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
文件说明: steps/StepDefineSounds.vue - 步骤2：定义声音
============================================================================

【文件作用】
本组件是键音专辑编辑器的第二个步骤（Step 2），负责：
1. 创建新的声音（从音频源文件中裁剪定义）
2. 编辑已有的声音配置
3. 预览声音效果
4. 删除声音

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(KEYTONE_ALBUM_CONTEXT_KEY, ctx)
        │
        ▼
  ┌─────────────────────────────────────────────────────┐
  │  q-stepper                                          │
  │  ├── StepLoadAudioFiles.vue       (Step 1)          │
  │  ├── StepDefineSounds.vue    <── 当前文件 (Step 2)  │
  │  ├── StepCraftKeySounds.vue       (Step 3)          │
  │  └── StepLinkageEffects.vue       (Step 4)          │
  └─────────────────────────────────────────────────────┘

【数据流】
本组件不持有任何状态，所有状态都通过 inject 从父组件获取：

  父组件状态                        本组件使用方式
  ───────────────────────────────────────────────────────
  ctx.step                     ->
控制当前步骤，点击 header 切换 ctx.soundList -> 已定义的声音列表（用于判断 done 状态） ctx.createNewSound ->
控制"创建声音"对话框显示 ctx.showEditSoundDialog -> 控制"编辑声音"对话框显示 ctx.soundFileList ->
可选择的音频源文件列表 ctx.sourceFileForSound -> 当前选中的源文件 ctx.soundStartTime/EndTime -> 裁剪时间范围
ctx.soundVolume -> 音量调整 ctx.saveSoundConfig() -> 保存声音配置 ctx.previewSound() -> 预览声音 ctx.deleteSound()
-> 删除声音 【关联文件】 - ../types.ts : 类型定义，包含 KEYTONE_ALBUM_CONTEXT_KEY - ../../Keytone_album.vue :
父组件，提供 Context - ../../DependencyWarning.vue : 依赖警告组件

【当前状态】
✅ 本组件已集成到父组件中：父组件 provide Context，并以 `<StepDefineSounds />` 替换原 Step2 模板。

============================================================================
-->

<template>
  <q-step
    :name="2"
    :title="ctx.$t('KeyToneAlbum.defineSounds.title')"
    icon="add_comment"
    :done="ctx.soundList.value.length !== 0"
    :disable="ctx.step.value === 99 && ctx.soundList.value.length === 0"
    :header-nav="false"
    @click="handleStepClick"
  >
    <!-- 步骤说明 -->
    <div :class="['mb-3', ctx.step_introduce_fontSize.value]">
      {{ ctx.$t('KeyToneAlbum.defineSounds.description') }}
      <q-icon name="info" color="primary">
        <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
          <div>{{ ctx.$t('KeyToneAlbum.defineSounds.tooltip.noImpactOnSource') }}</div>
          <div>{{ ctx.$t('KeyToneAlbum.defineSounds.tooltip.multipleSoundsFromSameSource') }}</div>
        </q-tooltip>
      </q-icon>
    </div>

    <!-- 定义声音的业务逻辑 -->
    <div>
      <!-- 制作新的声音 -->
      <div>
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.defineSounds.createNewSound')"
          @click="ctx.createNewSound.value = !ctx.createNewSound.value"
        />

        <!-- 创建声音对话框 -->
        <CreateSoundDialog />
      </div>

      <div :class="['p-2 text-zinc-600']">{{ ctx.$t('KeyToneAlbum.or') }}</div>

      <!-- 编辑已有声音 -->
      <div>
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.defineSounds.editExistingSound')"
          @click="handleEditSound"
        />

        <!-- 编辑声音对话框 -->
        <EditSoundDialog />
      </div>
    </div>

    <!-- 导航按钮 -->
    <q-stepper-navigation>
      <q-btn @click="ctx.step.value = 3" color="primary" :label="ctx.$t('KeyToneAlbum.continue')" />
      <q-btn flat @click="ctx.step.value = 1" color="primary" :label="ctx.$t('KeyToneAlbum.back')" class="q-ml-sm" />
    </q-stepper-navigation>
  </q-step>
</template>

<script setup lang="ts">
/**
 * StepDefineSounds.vue - 步骤2：定义声音
 *
 * 【组件职责】
 * - 只负责 UI 渲染和用户交互
 * - 所有状态和业务逻辑通过 inject 从父组件获取
 * - 不持有任何本地状态（除了纯 UI 状态）
 *
 * 【数据流向】
 * 1. 用户点击按钮 → 调用 ctx.xxx.value = ... 修改状态
 * 2. 状态变化 → 父组件的 watch 触发 → ConfigSet 写入配置
 * 3. SDK 后端处理 → SSE 推送更新 → 父组件更新状态 → UI 自动刷新
 */

import { inject } from 'vue';
import { useQuasar } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import CreateSoundDialog from '../dialogs/CreateSoundDialog.vue';
import EditSoundDialog from '../dialogs/EditSoundDialog.vue';

// 移除 defineAsyncComponent，直接引用以确保组件与 step 生命周期绑定，
// 或者保持现状，但确保 dialog 销毁机制正确。
// 鉴于用户反馈 Dialog 关闭后声音未停，我们需要确保 Dialog 组件能正确捕获关闭事件。
// 这里直接使用导入的组件即可。

const q = useQuasar();

// ============================================================================
// 注入父组件提供的上下文
// ============================================================================
// ctx 包含所有需要的状态和方法，子组件不需要自己定义
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// ============================================================================
// 事件处理函数
// ============================================================================

/**
 * 处理 step header 点击事件
 * 实现步骤的折叠/展开功能
 *
 * 【行为说明】
 * - 点击 header 区域时，在当前步骤 (2) 和折叠状态 (99) 之间切换
 * - 只响应 header 区域的点击，不响应内容区域
 */
function handleStepClick(event: MouseEvent) {
  const header = (event.target as HTMLElement).closest('.q-stepper__tab');
  if (header) {
    ctx.step.value = ctx.step.value === 2 ? 99 : 2;
  }
}

/**
 * 处理"编辑已有声音"按钮点击
 *
 * 【行为说明】
 * - 如果没有可编辑的声音，显示警告通知
 * - 否则打开编辑对话框
 */
function handleEditSound() {
  if (ctx.soundList.value.length === 0) {
    q.notify({
      type: 'warning',
      message: ctx.$t('KeyToneAlbum.notify.noSoundsToEdit'),
      position: 'top',
    });
    return;
  }
  ctx.showEditSoundDialog.value = true;
}
</script>

<style lang="scss" scoped>
/**
 * StepDefineSounds 组件样式
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
