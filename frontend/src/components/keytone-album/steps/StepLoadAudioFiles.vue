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
文件说明: steps/StepLoadAudioFiles.vue - 步骤1：加载音频源文件
============================================================================

【文件作用】
本组件是键音专辑编辑器的第一个步骤（Step 1），负责：
1. 显示"添加新的音频源文件"按钮，打开添加对话框
2. 显示"管理已有音频文件"按钮，打开管理对话框
3. 处理 Step header 的点击事件（折叠/展开）
4. 提供"继续"按钮跳转到下一步

【在整体架构中的位置】

  Keytone_album.vue (父组件)
        │
        │ provide(KEYTONE_ALBUM_CONTEXT_KEY, ctx)
        │
        ▼
  ┌─────────────────────────────────────────────────────┐
  │  q-stepper                                          │
  │  ├── StepLoadAudioFiles.vue  <── 当前文件 (Step 1)  │
  │  ├── StepDefineSounds.vue         (Step 2)          │
  │  ├── StepCraftKeySounds.vue       (Step 3)          │
  │  └── StepLinkageEffects.vue       (Step 4)          │
  └─────────────────────────────────────────────────────┘

【数据流】
本组件不持有任何状态，所有状态都通过 inject 从父组件获取：

  父组件状态                    本组件使用方式
  ─────────────────────────────────────────────────
  ctx.step                 ->
控制当前步骤，点击 header 切换 ctx.soundFileList -> 判断是否已有文件（done 状态） ctx.addNewSoundFile ->
控制"添加文件"对话框显示 ctx.editSoundFile -> 控制"管理文件"对话框显示 【关联文件】 - ../types.ts : 类型定义，包含
KEYTONE_ALBUM_CONTEXT_KEY - ../dialogs/AddAudioFileDialog.vue : 本组件内嵌的"添加音频文件"对话框 -
../dialogs/ManageAudioFilesDialog.vue : 本组件内嵌的"管理音频文件"对话框 - ../../Keytone_album.vue : 父组件，提供 Context

【当前状态】
✅ 本组件已集成到父组件中：父组件 provide Context，并以 `<StepLoadAudioFiles />` 替换原 Step1 模板。

============================================================================
-->

<template>
  <q-step
    :name="1"
    :title="ctx.$t('KeyToneAlbum.loadAudioFile.title')"
    icon="create_new_folder"
    :done="ctx.soundFileList.value.length !== 0"
    :disable="ctx.step.value === 99 && ctx.soundFileList.value.length === 0"
    :header-nav="false"
    @click="handleStepClick"
  >
    <div :class="['mb-3', ctx.step_introduce_fontSize.value]">
      {{ ctx.$t('KeyToneAlbum.loadAudioFile.description') }}
    </div>

    <!-- 载入音频文件的业务逻辑 -->
    <div>
      <!-- 添加新的音频源文件 -->
      <div>
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.loadAudioFile.addNewFile')"
          @click="ctx.addNewSoundFile.value = !ctx.addNewSoundFile.value"
        />

        <!-- 添加音频文件对话框 -->
        <AddAudioFileDialog />
      </div>

      <div :class="['p-2 text-zinc-600']">{{ ctx.$t('KeyToneAlbum.or') }}</div>

      <!-- 编辑已有音频源文件 -->
      <div>
        <q-btn
          :class="['bg-zinc-300']"
          :label="ctx.$t('KeyToneAlbum.loadAudioFile.manageExistingFiles')"
          @click="handleManageFiles"
        />

        <!-- 管理音频文件对话框 -->
        <ManageAudioFilesDialog />
      </div>
    </div>

    <!-- 导航按钮 -->
    <q-stepper-navigation>
      <q-btn @click="ctx.step.value = 2" color="primary" :label="ctx.$t('KeyToneAlbum.continue')" />
    </q-stepper-navigation>
  </q-step>
</template>

<script setup lang="ts">
import { inject } from 'vue';
import { useQuasar } from 'quasar';
import { KEYTONE_ALBUM_CONTEXT_KEY, type KeytoneAlbumContext } from '../types';
import AddAudioFileDialog from '../dialogs/AddAudioFileDialog.vue';
import ManageAudioFilesDialog from '../dialogs/ManageAudioFilesDialog.vue';

const q = useQuasar();

// 注入父组件提供的上下文
const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;

// 处理 step header 点击
function handleStepClick(event: MouseEvent) {
  const header = (event.target as HTMLElement).closest('.q-stepper__tab');
  if (header) {
    ctx.step.value = ctx.step.value === 1 ? 99 : 1;
  }
}

// 处理管理文件按钮点击
function handleManageFiles() {
  if (ctx.soundFileList.value.length === 0) {
    q.notify({
      type: 'warning',
      message: ctx.$t('KeyToneAlbum.notify.noFilesToManage'),
      position: 'top',
    });
    return;
  }
  ctx.editSoundFile.value = !ctx.editSoundFile.value;
}
</script>

<style lang="scss" scoped>
/**
 * StepLoadAudioFiles 组件样式
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
