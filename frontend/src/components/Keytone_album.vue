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
  <q-page :style="{ '--i18n_fontSize': i18n_fontSize }">
    <q-scroll-area :class="[isMacOS ? 'w-[389.5px] h-[458.5px]' : 'w-[379px] h-[458.5px]']">
      <div :class="['flex flex-col gap-5  p-8  scale-102']">
        <q-input
          outlined
          stack-label
          dense
          :error-message="$t('KeyToneAlbum.new.name.errorMessage')"
          :error="pkgName === '' || pkgName === undefined || pkgName === null"
          v-model="pkgName"
          :label="$t('KeyToneAlbum.new.name.name')"
          :placeholder="$t('KeyToneAlbum.new.name.defaultValue')"
        />
        <!-- <div>原始声音文件编辑</div>
        <div>键音</div>
        <div>键音列表, 编辑键音</div>
        <div>全局键音规则</div>
        <div>对某个特定按键单独设置键音</div> -->

        <q-stepper v-model="step" vertical header-nav color="primary" animated class="step-custom">
          <div
            :class="[
              // 字体
              'font-semibold text-lg',
              // 对溢出的情况, 采取滚动策略
              'max-w-66 overflow-auto whitespace-nowrap text-nowrap',
              // 隐藏滚动策略的滚动条。
              '[&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none]',
              // 居中对齐
              'mx-auto',
            ]"
          >
            {{ pkgName }}
          </div>
          <!-- Step 1: 加载音频源文件 (已拆分为独立组件) -->
          <StepLoadAudioFiles />

          <!-- Step 2: 定义声音 (已拆分为独立组件) -->
          <StepDefineSounds />

          <!-- Step 3: 制作按键音 (已拆分为独立组件) -->
          <StepCraftKeySounds />

          <!-- <q-step :name="4" title="对全局按键统一设置键音" icon="settings" :done="step > 3">
            <div>设置一个全局所有按键统一使用的按键声音。</div>
            <div>小提示: 用随机按键声音进行此项设置, 可避免键音太过单调。</div>
            <div>小提示: 如果您需要更加全面的键音定制, 可在下一步骤中处理。</div>
            <q-stepper-navigation>
              <q-btn @click="step = 5" color="primary" label="Continue" />
              <q-btn flat @click="step = 3" color="primary" label="Back" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>

          <q-step :name="5" title="对具体按键单独设置键音" caption="本步骤可选(非必填)" icon="settings">
            <div>如果您希望单独禁用某几个按键的全局键音。</div>
            <div>或是有些情况下, 我们希望独立定义某个按键的键音 。</div>
            <div>甚至, 更极端的情况, 我们希望键盘上所有的按键, 都拥有自己的独立键音。</div>
            <div>来吧!这一定制步骤将满足您的需求!!!</div>
            <div>小提示: 本步骤所做的设置, 优先级高于全局键音设置。</div>
            <div></div>
            <q-stepper-navigation>
              <q-btn flat @click="step = 4" color="primary" label="Back" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step> -->
          <q-step
            :name="4"
            :title="$t('KeyToneAlbum.linkageEffects.title')"
            icon="settings"
            :done="
              !(
                isEnableEmbeddedTestSound.down === true &&
                isEnableEmbeddedTestSound.up === true &&
                !keyDownUnifiedSoundEffectSelect &&
                !keyUpUnifiedSoundEffectSelect &&
                keysWithSoundEffect.size === 0
              )
            "
            :disable="
              step === 99 &&
              isEnableEmbeddedTestSound.down === true &&
              isEnableEmbeddedTestSound.up === true &&
              !keyDownUnifiedSoundEffectSelect &&
              !keyUpUnifiedSoundEffectSelect &&
              keysWithSoundEffect.size === 0
            "
            :header-nav="false"
            @click="
              (event: MouseEvent) => {
                const header = (event.target as HTMLElement).closest('.q-stepper__tab');
                if (header) {
                  step = step === 4 ? 99 : 4;
                }
              }
            "
          >
            <div :class="['mb-3', step_introduce_fontSize]">
              {{ $t('KeyToneAlbum.linkageEffects.description') }}
              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                  <span>{{ $t('KeyToneAlbum.linkageEffects.tooltips.description') }}</span>
                </q-tooltip>
              </q-icon>
            </div>
            <div :class="['flex items-center m-t-2 w-[130%]']">
              <span class="text-gray-500 mr-0.7">•</span>
              <span class="text-nowrap">
                {{ $t('KeyToneAlbum.linkageEffects.enableTestSound') }}:
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
                    <span>{{ $t('KeyToneAlbum.linkageEffects.tooltips.testSound') }}</span>
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
                v-model="isEnableEmbeddedTestSound.down"
                color="primary"
                :label="$t('KeyToneAlbum.linkageEffects.downTestSound')"
                dense
              />
            </div>
            <div
              :class="[
                'flex items-center ml-3',
                setting_store.languageDefault === 'fr' ? 'flex-nowrap text-nowrap' : '',
              ]"
            >
              <span class="text-gray-500 mr-1.5">•</span>
              <q-toggle
                v-model="isEnableEmbeddedTestSound.up"
                color="primary"
                :label="$t('KeyToneAlbum.linkageEffects.upTestSound')"
                dense
              />
            </div>
            <q-stepper-navigation>
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  :label="$t('KeyToneAlbum.linkageEffects.globalSettings')"
                  @click="
                    () => {
                      showEveryKeyEffectDialog = true;
                    }
                  "
                >
                </q-btn>
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
                    <span>{{ $t('KeyToneAlbum.linkageEffects.tooltips.globalPriority') }}</span>
                  </q-tooltip>
                </q-icon>
                <!-- 全键声效设置对话框（已拆分为独立组件） -->
                <EveryKeyEffectDialog />
              </div>
              <div :class="['p-2 text-zinc-600']">{{ $t('KeyToneAlbum.or') }}</div>
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  :label="$t('KeyToneAlbum.linkageEffects.singleKeySettings')"
                  @click="
                    () => {
                      showSingleKeyEffectDialog = true;
                    }
                  "
                >
                </q-btn>
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
                    <span>{{ $t('KeyToneAlbum.linkageEffects.tooltips.singleKeyPriority') }}</span>
                  </q-tooltip>
                </q-icon>
                <!-- 单键声效设置对话框（已拆分为独立组件） -->
                <SingleKeyEffectDialog />

                              </div>
            </q-stepper-navigation>
            <q-stepper-navigation>
              <q-btn @click="step = 5" color="primary" :label="$t('KeyToneAlbum.continue')" />
              <q-btn flat @click="step = 3" color="primary" :label="$t('KeyToneAlbum.back')" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>
        </q-stepper>
      </div>
    </q-scroll-area>
  </q-page>
</template>

<script setup lang="ts">
import { debounce } from 'lodash';
import { nanoid } from 'nanoid';
import { QDialog, QSelect, useQuasar } from 'quasar';
import {
  ConfigGet,
  ConfigSet,
  LoadConfig,
  SendFileToServer,
  SoundFileRename,
  SoundFileDelete,
  PlaySound,
  ConfigDelete,
} from 'src/boot/query/keytonePkg-query';
import { useAppStore } from 'src/stores/app-store';
import { useKeyEventStore } from 'src/stores/keyEvent-store';
import { useMainStore } from 'src/stores/main-store';
import { useSettingStore } from 'src/stores/setting-store';
import { computed, onBeforeMount, ref, watch, useTemplateRef, reactive, nextTick, onUnmounted, provide } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  createDependencyValidator,
  hasItemDependencyIssues,
  type DependencyIssue,
  type AudioFile,
  type Sound,
  type KeySound
} from 'src/utils/dependencyValidator';
import DependencyWarning from 'src/components/DependencyWarning.vue';

// ============================================================================
// 导入 Context 类型和注入 Key
// ============================================================================
import {
  KEYTONE_ALBUM_CONTEXT_KEY,
  type KeytoneAlbumContext
} from './keytone-album/types';

// ============================================================================
// 导入拆分后的子组件
// ============================================================================
import StepLoadAudioFiles from './keytone-album/steps/StepLoadAudioFiles.vue';
import StepDefineSounds from './keytone-album/steps/StepDefineSounds.vue';
import StepCraftKeySounds from './keytone-album/steps/StepCraftKeySounds.vue';
import StepLinkageEffects from './keytone-album/steps/StepLinkageEffects.vue';
import EveryKeyEffectDialog from './keytone-album/dialogs/EveryKeyEffectDialog.vue';
import SingleKeyEffectDialog from './keytone-album/dialogs/SingleKeyEffectDialog.vue';

// console.error("重新载入")   // 用笨方法, 严重组件的重新渲染情况
const q = useQuasar();
const { t } = useI18n();
const $t = t;
const app_store = useAppStore();
const setting_store = useSettingStore();

export interface Props {
  pkgPath: string;
  isCreate: boolean;
}
const props = withDefaults(defineProps<Props>(), {});

// 防止空字符串触发不能为空的提示, 虽然初始化时只有一瞬间, 但也不希望看到
const pkgName = ref<string>($t('KeyToneAlbum.new.name.defaultValue'));

const step = ref(99);
// watch(step, () => {
//   console.log('step-------------=', step.value);
// });

const addNewSoundFile = ref(false);
const files = ref<Array<File>>([]);
watch(files, () => {
  console.debug('观察files=', files.value);
});

const editSoundFile = ref(false);
// 用于初步映射配置文件中的 audio_files 对象, 并将其转换为数组, 并将数组元素转换成对象, 其中包含sha256和value两个key
const audioFiles = ref<Array<any>>([]);
// 用于audioFiles映射后的进一步映射, 主要拆分出value中name的每个值, 并一一对于sha256, 形成ui列表可用的最终数组。
const soundFileList = ref<Array<{ sha256: string; name_id: string; name: string; type: string }>>([]);
const selectedSoundFile = ref<{ sha256: string; name_id: string; name: string; type: string }>({
  sha256: '',
  name_id: '',
  name: '',
  type: '',
});

// 声音制作(制作新的声音)
const createNewSound = ref(false);
const soundName = ref<string>('');
const sourceFileForSound = ref<{ sha256: string; name_id: string; name: string; type: string }>({
  sha256: '',
  name_id: '',
  name: '',
  type: '',
});
const soundStartTime = ref<number>(0);
const soundEndTime = ref<number>(0);
const soundVolume = ref<number>(0.0);

// 修改 confirmAddingSound 函数为更通用的形式, 并将 confirmAddingSound 重构名称为 saveSoundConfig
function saveSoundConfig(params: {
  soundKey?: string; // 可选参数，存在则为修改操作，不存在则为添加操作
  source_file_for_sound: { sha256: string; name_id: string; type: string };
  name: string; // 声音名称, 一般为空, 也可由用户自行定义
  cut: {
    start_time: number;
    end_time: number;
    volume: number;
  };
  onSuccess?: () => void; // 成功后的回调函数, 可用于一些清空操作
}) {
  // 必须选择一个源文件
  if (
    params.source_file_for_sound.sha256 === '' &&
    params.source_file_for_sound.type === '' &&
    params.source_file_for_sound.name_id === ''
  ) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.selectSourceFile'),
      timeout: 5,
    });
    return;
  }
  // 结束时间必须大于开始时间
  if (params.cut.end_time <= params.cut.start_time) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.endTimeGreaterThanStartTime'),
      timeout: 5,
    });
    return;
  }

  // 创建一个新对象,不包含soundKey和onSuccess回调
  const configParams = {
    source_file_for_sound: params.source_file_for_sound,
    name: params.name,
    cut: params.cut,
  };

  // 如果有soundKey则为修改操作，否则为添加操作（生成新的key）
  const key = params.soundKey || nanoid();

  ConfigSet('sounds.' + key, configParams).then((re) => {
    if (re) {
      q.notify({
        type: 'positive',
        position: 'top',
        message: params.soundKey ? $t('KeyToneAlbum.notify.modifySuccess') : $t('KeyToneAlbum.notify.addSuccess'),
        timeout: 5,
      });
      params.onSuccess?.();
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: params.soundKey ? $t('KeyToneAlbum.notify.modifyFailed') : $t('KeyToneAlbum.notify.addFailed'),
        timeout: 5,
      });
    }
  });
}

function deleteSound(params: { soundKey: string; onSuccess?: () => void }) {
  ConfigDelete('sounds.' + params.soundKey).then((re) => {
    if (re) {
      params.onSuccess?.();
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: $t('KeyToneAlbum.notify.deleteFailed'),
        timeout: 5,
      });
    }
  });
}

// 重构previewSound函数,使用相同的参数结构
function previewSound(params: {
  source_file_for_sound: { sha256: string; name_id: string; type: string };
  cut: {
    start_time: number;
    end_time: number;
    volume: number;
  };
}) {
  console.debug('预览声音');
  if (
    params.source_file_for_sound.sha256 === '' &&
    params.source_file_for_sound.type === '' &&
    params.source_file_for_sound.name_id === ''
  ) {
    q.notify({
      type: 'warning',
      position: 'top',
      message: $t('KeyToneAlbum.notify.selectAudioFile'),
      timeout: 5000,
    });
    return;
  }
  // 结束时间必须大于开始时间
  if (params.cut.end_time <= params.cut.start_time) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.endTimeGreaterThanStartTime'),
      timeout: 5,
    });
    return;
  }
  // 时间值不能为负数
  if (params.cut.start_time < 0 || params.cut.end_time < 0) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.timeValueCannotBeNegative'),
      timeout: 5,
    });
  }

  PlaySound(
    params.source_file_for_sound.sha256,
    params.source_file_for_sound.type,
    params.cut.start_time,
    params.cut.end_time,
    params.cut.volume,
    true // 设置 skipGlobalVolume 为 true，使预览不受全局音量影响
  ).then((result) => {
    if (!result) {
      q.notify({
        type: 'negative',
        position: 'top',
        message: $t('KeyToneAlbum.notify.playFailed'),
        timeout: 5000,
      });
    }
  });
}

// 声音编辑(编辑已有声音)
const showEditSoundDialog = ref(false);
const soundList = ref<
  Array<{
    soundKey: string;
    soundValue: {
      cut: { start_time: number; end_time: number; volume: number };
      name: string;
      source_file_for_sound: { sha256: string; name_id: string; type: string };
    };
  }>
>([]); // TIPS: 此处类型一定要指定清楚, 否则在 组件<q-select :option-label="(item)=>{}"> 中, 会发生类型错误(且检测不到原因, 准确指定类型不使用any后, 才解决此问题)--- 我记得之前在做声音源文件的编辑列表时, 就好像遇到过, 但目前找不到注释内容了。
const selectedSound = ref<{
  soundKey: string;
  soundValue: {
    cut: { start_time: number; end_time: number; volume: number };
    name: string;
    source_file_for_sound: { sha256: string; name_id: string; type: string };
  };
}>(); // 此处无需初始化, 但类型一定要指定清楚

watch(selectedSound, () => {
  console.debug('观察selectedSound=', selectedSound.value);
});

// 选项列表(其中包含: 源文件、声音、按键音 三种选项可供选择)
const options = reactive([
  { label: 'KeyToneAlbum.options.audioFile', value: 'audio_files', label_0: 'KeyToneAlbum.options.audioFile_0' },
  { label: 'KeyToneAlbum.options.sound', value: 'sounds', label_0: 'KeyToneAlbum.options.sound_0' },
  { label: 'KeyToneAlbum.options.keySound', value: 'key_sounds', label_0: 'KeyToneAlbum.options.keySound_0' },
]);

// ============================================================================
// 自然排序工具函数：解决下拉列表选项乱序问题
//
// 问题：当选项数量很多的时候，列表的排序没有规律，用户难以快速定位所需内容
// 解决方案：实现字母和数字的自然顺序排序作为默认排序规则
//
// 功能：
// - 支持字母和数字的混合排序（如：sound1, sound2, sound10 而不是 sound1, sound10, sound2）
// - 不区分大小写的字母排序
// - 自动处理中文、英文、数字的混合内容
// ============================================================================
const naturalSort = (a: string, b: string): number => {
  // 将字符串按数字和字母分段处理，支持混合文本数字内容的智能排序
  const segmentize = (str: string) => {
    return str.match(/\d+|\D+/g) || [];
  };

  const aSegments = segmentize(a.toLowerCase());
  const bSegments = segmentize(b.toLowerCase());

  const maxLength = Math.max(aSegments.length, bSegments.length);

  for (let i = 0; i < maxLength; i++) {
    const aSegment = aSegments[i] || '';
    const bSegment = bSegments[i] || '';

    // 如果两个段都是数字，按数值比较（确保 1, 2, 10 的正确顺序）
    if (/^\d+$/.test(aSegment) && /^\d+$/.test(bSegment)) {
      const diff = parseInt(aSegment, 10) - parseInt(bSegment, 10);
      if (diff !== 0) return diff;
    } else {
      // 否则按字符串比较（支持中英文混合）
      if (aSegment < bSegment) return -1;
      if (aSegment > bSegment) return 1;
    }
  }

  return 0;
};

const album_options_select_label = (item: any): any => {
  // console.log('item_1212==', item);
  // console.log('至臻键音列表==', keySoundList.value);
  if (item.type === 'audio_files') {
    return (
      $t(options.find((option) => item.type === option.value)?.label_0 || '') +
      ' § ' +
      soundFileList.value.find(
        (soundFile: any) => soundFile.sha256 === item.value?.sha256 && soundFile.name_id === item.value?.name_id
      )?.name +
      soundFileList.value.find(
        (soundFile: any) => soundFile.sha256 === item.value?.sha256 && soundFile.name_id === item.value?.name_id
      )?.type
    );
  }
  if (item.type === 'sounds') {
    // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
    if (item.value?.soundValue?.name !== '' && item.value?.soundValue?.name !== undefined) {
      return (
        $t(options.find((option) => item.type === option.value)?.label_0 || '') +
        ' § ' +
        soundList.value.find((sound) => sound.soundKey === item.value?.soundKey)?.soundValue.name
      );
    } else {
      return (
        $t(options.find((option) => item.type === option.value)?.label_0 || '') +
        ' § ' +
        (soundFileList.value.find(
          (soundFile: any) =>
            soundFile.sha256 === item.value?.soundValue?.source_file_for_sound?.sha256 &&
            soundFile.name_id === item.value?.soundValue?.source_file_for_sound?.name_id
        )?.name +
          '     - ' +
          ' [' +
          item.value?.soundValue?.cut?.start_time +
          ' ~ ' +
          item.value?.soundValue?.cut?.end_time +
          ']')
      );
    }
  }
  if (item.type === 'key_sounds') {
    // Check if item.value is valid before accessing its properties
    return (
      $t(options.find((option) => item.type === option.value)?.label_0 || 'Error') +
      ' § ' +
      // (item.value.keySoundValue.name || '[Unnamed]')
      keySoundList.value.find((keySound) => keySound.keySoundKey === item.value?.keySoundKey)?.keySoundValue?.name
    );
  }
};

// 按键音
const playModeOptions = ['single', 'random', 'loop'];
const playModeLabels = new Map<string, string>([
  ['single', 'KeyToneAlbum.playMode.single'],
  ['random', 'KeyToneAlbum.playMode.random'],
  ['loop', 'KeyToneAlbum.playMode.loop'],
]);

// 按键音制作
const createNewKeySound = ref(false);

// -- createNewKeySound
const keySoundName = ref<string>($t('KeyToneAlbum.craftKeySounds.keySoundName-placeholder'));
const configureDownSound = ref(false);
const configureUpSound = ref(false);

// -- configureDownSound
const selectedSoundsForDown = ref<Array<any>>([]);
const playModeForDown = ref('random');
const maxSelectionForDown = computed(() => {
  return playModeForDown.value === 'single' ? 1 : Infinity;
});
/* --- 在vue3.5中, 用useTemplateRef方式获取dom元素, 有助于增强可读性
const downSoundSelectDom = ref<QSelect>(); // 在vue3.5中, 用useTemplateRef方式获取dom元素, 有助于增强可读性*/
const downSoundSelectDom = useTemplateRef<QSelect>('downSoundSelectDom');
const downTypeGroup = ref<Array<string>>(['sounds']);
const downSoundList = computed(() => {
  const List: Array<any> = [];
  if (downTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (downTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (downTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  // ===== 应用自然排序：解决混合选项列表乱序问题 =====
  // 对混合类型的选项列表使用显示标签进行自然排序
  List.sort((a, b) => {
    const aLabel = album_options_select_label(a);
    const bLabel = album_options_select_label(b);
    return naturalSort(aLabel, bLabel);
  });
  console.debug('downSoundList=', List);
  return List;
});

// -- configureUpSound
const selectedSoundsForUp = ref<Array<any>>([]);
const playModeForUp = ref('random');
const maxSelectionForUp = computed(() => {
  return playModeForUp.value === 'single' ? 1 : Infinity;
});
/* --- 在vue3.5中, 用useTemplateRef方式获取dom元素, 有助于增强可读性
const upSoundSelectDom = ref<QSelect>();*/
const upSoundSelectDom = useTemplateRef<QSelect>('upSoundSelectDom');
const upTypeGroup = ref<Array<string>>(['sounds']);
const upSoundList = computed(() => {
  const List: Array<any> = [];
  if (upTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (upTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (upTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  // ===== 应用自然排序：解决混合选项列表乱序问题 =====
  // 对混合类型的选项列表使用显示标签进行自然排序
  List.sort((a, b) => {
    const aLabel = album_options_select_label(a);
    const bLabel = album_options_select_label(b);
    return naturalSort(aLabel, bLabel);
  });
  console.debug('upSoundList=', List);
  return List;
});

// 按键音编辑
const editExistingKeySound = ref(false);

// -- editExistingKeySound
const edit_configureDownSound = ref(false);
const edit_configureUpSound = ref(false);

// -- edit_configureDownSound / edit_configureUpSound
const edit_downSoundSelectDom = useTemplateRef<QSelect>('edit_downSoundSelectDom');
const edit_upSoundSelectDom = useTemplateRef<QSelect>('edit_upSoundSelectDom');
const edit_downTypeGroup = ref<Array<string>>(['sounds']);
const edit_upTypeGroup = ref<Array<string>>(['sounds']);
const edit_downSoundList = computed(() => {
  const List: Array<any> = [];
  if (edit_downTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (edit_downTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (edit_downTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  // ===== 应用自然排序：解决混合选项列表乱序问题 =====
  // 对混合类型的选项列表使用显示标签进行自然排序
  List.sort((a, b) => {
    const aLabel = album_options_select_label(a);
    const bLabel = album_options_select_label(b);
    return naturalSort(aLabel, bLabel);
  });
  console.debug('edit_downSoundList=', List);
  return List;
});
const edit_upSoundList = computed(() => {
  const List: Array<any> = [];
  if (edit_upTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (edit_upTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (edit_upTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  // ===== 应用自然排序：解决混合选项列表乱序问题 =====
  // 对混合类型的选项列表使用显示标签进行自然排序
  List.sort((a, b) => {
    const aLabel = album_options_select_label(a);
    const bLabel = album_options_select_label(b);
    return naturalSort(aLabel, bLabel);
  });
  console.debug('edit_upSoundList=', List);
  return List;
});
const keySoundList = ref<Array<any>>([]);
const selectedKeySound = ref<any>();

// 改变selectedKeySound.value.keySoundValue.down.value和selectedKeySound.value.keySoundValue.up.value的类型结构, 使其符合选择输入框组件的使用需求
// 使用深拷贝避免修改原始数据，确保依赖验证准确性
watch(selectedKeySound, (newVal, oldVal) => {
  if (!newVal) {
    return;
  }

  // 检查是否是新选择的KeySound（避免重复处理已转换的数据）
  if (oldVal && newVal.keySoundKey === oldVal.keySoundKey) {
    return;
  }

  console.debug('观察selectedKeySound=', newVal);

  // 创建深拷贝避免修改keySoundList中的原始数据
  const originalKeySound = keySoundList.value.find(ks => ks.keySoundKey === newVal.keySoundKey);
  if (originalKeySound) {
    selectedKeySound.value = JSON.parse(JSON.stringify(originalKeySound));

    // 对拷贝的数据进行UI格式转换
    selectedKeySound.value.keySoundValue.down.value = selectedKeySound.value.keySoundValue.down.value.map((item: any) => {
      /**
       * json中的存储格式分别是
       *  {key:'audio_files', value:{sha256: string, name_id: string, type:string}}
       *  {key:'sounds', value:string} // 此处value, 是soundKey
       *  {key:'key_sounds', value:string} // 此处value, 是keySoundKey
       */
      if (item.type === 'audio_files') {
        return {
          type: 'audio_files',
          value: soundFileList.value.find(
            (soundFile) => item.value && soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id
          ),
        };
      }
      if (item.type === 'sounds') {
        return {
          type: 'sounds',
          value: soundList.value.find((sound) => sound.soundKey === item.value),
        };
      }
      if (item.type === 'key_sounds') {
        return {
          type: 'key_sounds',
          value: keySoundList.value.find((keySound) => keySound.keySoundKey === item.value),
        };
      }
      return item;
    });

    selectedKeySound.value.keySoundValue.up.value = selectedKeySound.value.keySoundValue.up.value.map((item: any) => {
      /**
       * json中的存储格式分别是
       *  {key:'audio_files', value:{sha256: string, name_id: string, type:string}}
       *  {key:'sounds', value:string} // 此处value, 是soundKey
       *  {key:'key_sounds', value:string} // 此处value, 是keySoundKey
       */
      if (item.type === 'audio_files') {
        return {
          type: 'audio_files',
          value: soundFileList.value.find(
            (soundFile) => item.value && soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id
          ),
        };
      }
      if (item.type === 'sounds') {
        return {
          type: 'sounds',
          value: soundList.value.find((sound) => sound.soundKey === item.value),
        };
      }
      if (item.type === 'key_sounds') {
        return {
          type: 'key_sounds',
          value: keySoundList.value.find((keySound) => keySound.keySoundKey === item.value),
        };
      }
      return item;
    });
  }
});

// 按键音api
// -- 保存按键音配置
function saveKeySoundConfig(
  params: {
    key: string;
    name: string;
    down: { mode: string; value: Array<any> };
    up: { mode: string; value: Array<any> };
  },
  onSuccess?: () => void
) {
  let isReturn = false;
  if (params.down.mode === 'single' && params.down.value.length > 1) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.downSoundInvalid'),
      timeout: 3000,
    });
    isReturn = true;
  }
  if (params.up.mode === 'single' && params.up.value.length > 1) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.upSoundInvalid'),
      timeout: 3000,
    });
    isReturn = true;
  }
  if (isReturn) {
    return;
  }

  const configParams = {
    name: params.name,
    down: {
      mode: params.down.mode,
      value: params.down.value.map((item) => {
        if (item.type === 'audio_files') {
          return {
            type: 'audio_files',
            value: { sha256: item.value.sha256, name_id: item.value.name_id, type: item.value.type },
          };
        }
        if (item.type === 'sounds') {
          return { type: 'sounds', value: item.value.soundKey };
        }
        if (item.type === 'key_sounds') {
          return { type: 'key_sounds', value: item.value.keySoundKey };
        }
      }),
    },
    up: {
      mode: params.up.mode,
      value: params.up.value.map((item) => {
        if (item.type === 'audio_files') {
          return {
            type: 'audio_files',
            value: { sha256: item.value.sha256, name_id: item.value.name_id, type: item.value.type },
          };
        }
        if (item.type === 'sounds') {
          return { type: 'sounds', value: item.value.soundKey };
        }
        if (item.type === 'key_sounds') {
          return { type: 'key_sounds', value: item.value.keySoundKey };
        }
      }),
    },
  };

  const key = params.key || nanoid();
  ConfigSet('key_sounds.' + key, configParams).then((re) => {
    if (re) {
      q.notify({
        type: 'positive',
        position: 'top',
        message: params.key ? $t('KeyToneAlbum.notify.modifySuccess') : $t('KeyToneAlbum.notify.addSuccess'),
        timeout: 5,
      });
      if (onSuccess) {
        onSuccess();
      }
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: params.key ? $t('KeyToneAlbum.notify.modifyFailed') : $t('KeyToneAlbum.notify.addFailed'),
        timeout: 5,
      });
    }
  });
}

// -- 删除键音
function deleteKeySound(params: { keySoundKey: string; onSuccess?: () => void }) {
  ConfigDelete('key_sounds.' + params.keySoundKey).then((re) => {
    if (re) {
      params.onSuccess?.();
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: $t('KeyToneAlbum.notify.deleteFailed'),
        timeout: 5,
      });
    }
  });
}

// 按键联动声效

// -- 内嵌测试音是否使能
const isEnableEmbeddedTestSound = reactive({
  down: true,
  up: true,
}); // 该字段"直接"与配置文件相映射

// -- 全键声效
const showEveryKeyEffectDialog = ref(false);

const keyDownUnifiedSoundEffectSelect = ref<any>();
const keyUpUnifiedSoundEffectSelect = ref<any>();
const unifiedTypeGroup = ref<Array<string>>(['sounds']);
const keyUnifiedSoundEffectOptions = computed(() => {
  const List: Array<any> = [];
  if (unifiedTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (unifiedTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (unifiedTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  // ===== 应用自然排序：解决混合选项列表乱序问题 =====
  // 对混合类型的选项列表使用显示标签进行自然排序
  List.sort((a, b) => {
    const aLabel = album_options_select_label(a);
    const bLabel = album_options_select_label(b);
    return naturalSort(aLabel, bLabel);
  });
  console.debug('观察keyUnifiedSoundEffectOptions=', List);
  return List;
});
const isShowUltimatePerfectionKeySoundAnchoring = computed(() => {
  return unifiedTypeGroup.value.includes('key_sounds');
});
const isAnchoringUltimatePerfectionKeySound = ref(true);
watch(keyDownUnifiedSoundEffectSelect, (newVal, oldVal) => {
  console.debug('观察keyDownUnifiedSoundEffectSelect=', keyDownUnifiedSoundEffectSelect.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (isShowUltimatePerfectionKeySoundAnchoring.value && isAnchoringUltimatePerfectionKeySound.value) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyUpUnifiedSoundEffectSelect.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyUpUnifiedSoundEffectSelect.value = keyDownUnifiedSoundEffectSelect.value;
      }
    }
  }
});
watch(keyUpUnifiedSoundEffectSelect, (newVal, oldVal) => {
  console.debug('观察keyUpUnifiedSoundEffectSelect=', keyUpUnifiedSoundEffectSelect.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (isShowUltimatePerfectionKeySoundAnchoring.value && isAnchoringUltimatePerfectionKeySound.value) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyDownUnifiedSoundEffectSelect.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyDownUnifiedSoundEffectSelect.value = keyUpUnifiedSoundEffectSelect.value;
      }
    }
  }
});

function saveUnifiedSoundEffectConfig(params: { down: any; up: any }, onSuccess?: () => void) {
  const keyTone_global = {
    down: {
      type: params.down?.type || '',
      // TIPS: 此处需要注意, 我们需要的是此lambda表达式执行后的返回值赋值给value, 而不是直接将lambda表达式赋值给value。(此处用三元表达式会更为直观)
      value: (() => {
        if (params.down?.type === 'audio_files') {
          return { sha256: params.down.value.sha256, name_id: params.down.value.name_id, type: params.down.value.type };
        }
        if (params.down?.type === 'sounds') {
          return params.down.value.soundKey;
        }
        if (params.down?.type === 'key_sounds') {
          return params.down.value.keySoundKey;
        }
        return '';
      })(),
    },
    up: {
      type: params.up?.type || '',
      // TIPS: 此处需要注意, 我们需要的是此lambda表达式执行后的返回值赋值给value, 而不是直接将lambda表达式赋值给value。(此处用三元表达式会更为直观)
      value: (() => {
        if (params.up?.type === 'audio_files') {
          return { sha256: params.up.value.sha256, name_id: params.up.value.name_id, type: params.up.value.type };
        }
        if (params.up?.type === 'sounds') {
          return params.up.value.soundKey;
        }
        if (params.up?.type === 'key_sounds') {
          return params.up.value.keySoundKey;
        }
        return '';
      })(),
    },
  };

  ConfigSet('key_tone.global', keyTone_global)
    .then((re) => {
      if (re) {
        onSuccess?.();
      } else {
        q.notify({
          type: 'negative',
          position: 'top',
          message: $t('KeyToneAlbum.notify.unifiedSoundEffectConfigFailed'),
          timeout: 5000,
        });
      }
    })
    .catch((err) => {
      console.error('全键声效配置时发生错误:', err);
      q.notify({
        type: 'negative',
        position: 'top',
        message: $t('KeyToneAlbum.notify.unifiedSoundEffectConfigFailed'),
        timeout: 5000,
      });
    });
}

// -- 单键声效

// -- -- 选择按键
const showSingleKeyEffectDialog = ref(false);

const isShowAddOrSettingSingleKeyEffectDialog = ref(false);

const singleKeysSelectRef = useTemplateRef<QSelect>('singleKeysSelectRef');
const selectedSingleKeys = ref<Array<number>>([]);

const keyEvent_store = useKeyEventStore();

const isRecordingSingleKeys = ref(false);

const keyOptions = computed(() => {
  // 将 Map 转换为数组形式的选项
  if (isRecordingSingleKeys.value) {
    return [];
  } else {
    // 默认以系统主映射表的keys为主
    const reArray = Array.from(keyEvent_store.dikCodeToName.keys());

    return reArray;
  }
});
const filterOptions = ref(keyOptions.value); // 用于过滤选项

const isGetsFocused = ref(false);

let first_flag = false; // 用于避免录制按键打开瞬间(即isRecordingSingleKeys由false->true的瞬间)鼠标左键被记录。
let clear_flag = false; // 用于避免录制按键打开瞬间(即isRecordingSingleKeys由false->true的瞬间)鼠标左键被记录。

const setSingleKeyRecordingClearFlag = () => {
  clear_flag = true;
};

const recordingSingleKeysCallback = (keycode: number, keyName: string) => {
  console.debug('keycode=', keycode, 'keyName=', keyName);
  if (!first_flag) {
    first_flag = true;
    return;
  }
  if (clear_flag) {
    clear_flag = false;
    return;
  }

  // 如果按键不在列表中，则添加
  if (!selectedSingleKeys.value.includes(keycode)) {
    selectedSingleKeys.value.push(keycode);
  } else {
    q.notify({
      type: 'info',
      position: 'top',
      message: $t('KeyToneAlbum.notify.keyAlreadySelected'),
      timeout: 1000,
    });
  }

  console.debug('当前已选择的按键:', selectedSingleKeys.value);
};

watch(isShowAddOrSettingSingleKeyEffectDialog, (newVal) => {
  if (!newVal) {
    keyEvent_store.clearKeyStateCallback_Record();
    // 当通过点击对话框外使得对话框关闭时, 不会触发失去焦点的事件(因此此时isGetsFocused的值不会被置为false, 故补充此逻辑)
    isGetsFocused.value = false;
  }
});

watch(isRecordingSingleKeys, (newVal, oldVal) => {
  if (newVal) {
    // 录制单键时, 清空输入框。(由于是录制, 因此需要清空输入框, 防止用户输入内容。)
    // * 如何防止用户输入内容?
    // * * 当然也可以利用updateInputValue。但有更简单的解决思路, 即定义组件特有属性maxlength为0即可阻止用户输入内容。
    singleKeysSelectRef.value?.updateInputValue('');

    if (!oldVal) {
      first_flag = false;
    }

    keyEvent_store.setKeyStateCallback_Record(recordingSingleKeysCallback);
  } else {
    keyEvent_store.clearKeyStateCallback_Record();
  }
});

watch(isGetsFocused, (newVal) => {
  if (newVal && isRecordingSingleKeys.value) {
    keyEvent_store.setKeyStateCallback_Record(recordingSingleKeysCallback);
  } else {
    keyEvent_store.clearKeyStateCallback_Record();
  }
});

//   - completed(已完成)   FIXME: 修复'Backspace'按键, 在录制过程中, 删除已选择列表中最后一项的bug
let oldSelectedSingleKeys: Array<any> = [];

function preventDefaultKeyBehaviorWhenRecording(event: KeyboardEvent) {
  if (isRecordingSingleKeys.value) {
    // TIPS: 这里被打印两次的原因可能是以下, 不过不用担心和处理, 因为没有bug。
    // 1. 首先触发 input 元素的 keydown 事件
    // 2. 然后冒泡到 q-select 组件
    // 3. 最后最后冒泡到全局监听器
    console.debug('event.key=', event.key);

    //   - completed(已完成)   FIXME: 修复'Enter'按键无法被录制的bug
    if (event.key === 'Enter') {
      // 虽然无法录制'Enter'事件的原因就是select组件阻止了默认的'Enter'事件的冒泡行为,
      // * 但为防止quasar后续更新改变它, 便再次手动阻止一次, 以防止本次修复被quasar的更新影响。
      event.stopPropagation(); // 阻止事件冒泡

      //   ↓    - completed(已完成)   增加location字段即可。原理:这是手动构建事件时, 缺少了构建UUID必要的location字段, 导致小数字键盘中的'Enter'被按下时, 在前端事件状态集中, 创建了一个并不存在的 UUID 的及对应的按下状态 的假按键。
      // FIXME: 在数字键盘'enter'按键按下后, 再次按下任何按键都 报多个按键同时按下的消息, 这是一个bug。
      // 手动创建并分发一个新的键盘事件 TIPS: 没必要(虽然修复了, 但我并不打算继续这样用)
      // const newEvent = new KeyboardEvent('keydown', {
      //   key: event.key,
      //   code: event.code,
      //   keyCode: event.keyCode,
      //   which: event.which,
      //   altKey: event.altKey,
      //   location: event.location, // 增加此字段, 以修复FIXME。
      //   ctrlKey: event.ctrlKey,
      //   shiftKey: event.shiftKey,
      //   metaKey: event.metaKey,
      //   bubbles: true,
      //   cancelable: true,
      // });
      // document.dispatchEvent(newEvent);

      // 其实干脆直接使用当前事件创建新事件也行( TIPS: 不要直接使用event, 不然会因重复分发相同引用的事件而报错)
      const newEvent = new KeyboardEvent('keydown', event);
      document.dispatchEvent(newEvent);

      // TIPS: 这样是不对的-> 直接将当前事件, 原封不动的, 给到全局。 会引发意外的报错。
      // document.dispatchEvent(event); // 报错原因是由于这个事件是已经分发过的事件。
    }

    //   - completed(已完成)   FIXME: 修复'Backspace'按键, 在录制过程中, 删除已选择列表中最后一项的bug
    if (event.key === 'Backspace') {
      // nextTick是为了确保selectedSingleKeys.value = oldSelectedSingleKeys的逻辑, 发生在元素被删除之后。
      nextTick(() => {
        selectedSingleKeys.value = oldSelectedSingleKeys;
      });
    }
  }

  // 更新oldSelectedSingleKeys, 以备下次使用
  oldSelectedSingleKeys = selectedSingleKeys.value.slice();
}
const preventDefaultMouseWhenRecording = (event: MouseEvent) => {
  // TIPS: 需要在mouseup事件中阻止鼠标按键4、5的前进后退功能。
  //                  sdk的button值  |  前端的event.button值
  // 'MouseLeft'            1                    0
  // 'MouseRight'           2                    2
  // 'MouseMiddle'          3                    1
  // 'MouseBack'            4                    3
  // 'MouseForward'         5                    4
  // console.log(event.button);

  // 鼠标按钮4是后退，按钮5是前进
  if (event.button === 3 || event.button === 4) {
    // 虽然无法录制'Enter'事件的原因就是select组件阻止了默认的'Enter'事件的冒泡行为,
    // * 但为防止quasar后续更新改变它, 便再次手动阻止一次, 以防止本次修复被quasar的更新影响。
    event.preventDefault(); // 阻止默认行为
    event.stopPropagation(); // 阻止事件冒泡
  }
};

// -- -- 选择声效
const isDownSoundEffectSelectEnabled = ref(true);
const isUpSoundEffectSelectEnabled = ref(true);

const keyDownSingleKeySoundEffectSelect = ref<any>();
const keyUpSingleKeySoundEffectSelect = ref<any>();
// const keyDownSingleKeySoundEffectSelect_diff = keyDownSingleKeySoundEffectSelect.value;
// const keyUpSingleKeySoundEffectSelect_diff = keyUpSingleKeySoundEffectSelect.value;
const singleKeyTypeGroup = ref<Array<string>>(['sounds']);
const keySingleKeySoundEffectOptions = computed(() => {
  const List: Array<any> = [];
  if (singleKeyTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (singleKeyTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (singleKeyTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  // ===== 应用自然排序：解决混合选项列表乱序问题 =====
  // 对混合类型的选项列表使用显示标签进行自然排序
  List.sort((a, b) => {
    const aLabel = album_options_select_label(a);
    const bLabel = album_options_select_label(b);
    return naturalSort(aLabel, bLabel);
  });
  console.debug('观察keySingleKeySoundEffectOptions=', List);
  return List;
});
const isShowUltimatePerfectionKeySoundAnchoring_singleKey = computed(() => {
  return singleKeyTypeGroup.value.includes('key_sounds');
});
const isAnchoringUltimatePerfectionKeySound_singleKey = ref(true);
watch(keyDownSingleKeySoundEffectSelect, (newVal, oldVal) => {
  console.debug('观察keyDownSingleKeySoundEffectSelect=', keyDownSingleKeySoundEffectSelect.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (
      isShowUltimatePerfectionKeySoundAnchoring_singleKey.value &&
      isAnchoringUltimatePerfectionKeySound_singleKey.value
    ) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyUpSingleKeySoundEffectSelect.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyUpSingleKeySoundEffectSelect.value = keyDownSingleKeySoundEffectSelect.value;
      }
    }
  }
});
watch(keyUpSingleKeySoundEffectSelect, (newVal, oldVal) => {
  console.debug('观察keyUpSingleKeySoundEffectSelect=', keyUpSingleKeySoundEffectSelect.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (
      isShowUltimatePerfectionKeySoundAnchoring_singleKey.value &&
      isAnchoringUltimatePerfectionKeySound_singleKey.value
    ) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyDownSingleKeySoundEffectSelect.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyDownSingleKeySoundEffectSelect.value = keyUpSingleKeySoundEffectSelect.value;
      }
    }
  }
});

function saveSingleKeySoundEffectConfig(
  params: { singleKeys: Array<number>; down: any; up: any },
  onSuccess?: () => void
) {
  const keyTone_single = {
    down: {
      type: params.down?.type || '',
      // TIPS: 此处需要注意, 我们需要的是此lambda表达式执行后的返回值赋值给value, 而不是直接将lambda表达式赋值给value。(此处用三元表达式会更为直观)
      value: (() => {
        if (params.down?.type === 'audio_files') {
          return { sha256: params.down.value.sha256, name_id: params.down.value.name_id, type: params.down.value.type };
        }
        if (params.down?.type === 'sounds') {
          return params.down.value.soundKey;
        }
        if (params.down?.type === 'key_sounds') {
          return params.down.value.keySoundKey;
        }
        return '';
      })(),
    },
    up: {
      type: params.up?.type || '',
      // TIPS: 此处需要注意, 我们需要的是此lambda表达式执行后的返回值赋值给value, 而不是直接将lambda表达式赋值给value。(此处用三元表达式会更为直观)
      value: (() => {
        if (params.up?.type === 'audio_files') {
          return { sha256: params.up.value.sha256, name_id: params.up.value.name_id, type: params.up.value.type };
        }
        if (params.up?.type === 'sounds') {
          return params.up.value.soundKey;
        }
        if (params.up?.type === 'key_sounds') {
          return params.up.value.keySoundKey;
        }
        return '';
      })(),
    },
  };

  // 需要保证只有在对应的down或up声效设置被使能时, 才会修改对应的设置。(避免意外修改到不希望修改的配置)
  let downOrUpIfEnable: string;
  if (isDownSoundEffectSelectEnabled.value && !isUpSoundEffectSelectEnabled.value) {
    downOrUpIfEnable = '.down';
  } else if (!isDownSoundEffectSelectEnabled.value && isUpSoundEffectSelectEnabled.value) {
    downOrUpIfEnable = '.up';
  } else if (isDownSoundEffectSelectEnabled.value && isUpSoundEffectSelectEnabled.value) {
    downOrUpIfEnable = '';
  } else {
    // !isDownSoundEffectSelectEnabled.value && !isUpSoundEffectSelectEnabled.value
    // 此时, 说明用户未选择任何声效, 因此无需配置, 直接跳过。没必要继续执行后续的保存至 配置文件的步骤。
    q.notify({
      type: 'info',
      position: 'top',
      message: $t('KeyToneAlbum.notify.noSoundEffectSelected'),
      timeout: 5000,
    });
    return;
  }

  let allSuccess = true;
  const promises = params.singleKeys.map((item) => {
    return ConfigSet(
      'key_tone.single.' + item + downOrUpIfEnable,
      // `downOrUpIfEnable.slice(1)` 是去掉字符串开头的 . 如  '.down' -> 'down'
      // `as keyof typeof keyTone_single` 是 TypeScript 类型断言，确保属性名是 `keyTone_single` 对象的有效键名
      downOrUpIfEnable ? keyTone_single[downOrUpIfEnable.slice(1) as keyof typeof keyTone_single] : keyTone_single
    )
      .then((re) => {
        if (!re) {
          allSuccess = false;
          q.notify({
            type: 'negative',
            position: 'top',
            message: $t('KeyToneAlbum.notify.singleKeySoundEffectConfigFailed', {
              key: keyEvent_store.dikCodeToName.get(item) || 'Dik-{' + item + '}',
            }),
            timeout: 5000,
          });
        }
        return re;
      })
      .catch((err) => {
        allSuccess = false;
        console.error(
          $t('KeyToneAlbum.notify.singleKeySoundEffectConfigError', {
            key: keyEvent_store.dikCodeToName.get(item) || 'Dik-{' + item + '}',
          }),
          err
        );
        q.notify({
          type: 'negative',
          position: 'top',
          message: $t('KeyToneAlbum.notify.singleKeySoundEffectConfigFailed', {
            key: keyEvent_store.dikCodeToName.get(item) || 'Dik-{' + item + '}',
          }),
          timeout: 5000,
        });
      });
  });

  Promise.all(promises).then(() => {
    if (allSuccess) {
      onSuccess?.();
    } else {
      q.notify({
        type: 'warning',
        position: 'top',
        message: $t('KeyToneAlbum.notify.partialSingleKeySoundEffectConfigSuccess'),
        timeout: 5000,
      });
    }
  });
}

// -- -- 查看编辑声效
// const keysWithSoundEffect = ref<string[]>([]);
// watch(
//   // TIPS: 注意, 如果要监听的ref对象是数组或对象等js/ts中的默认引用类型, 要使用此种方式才可触发监听。(或者也可以弃用ref, 改用reactive。)
//   () => keysWithSoundEffect.value,
//   (newVal) => {
//     console.debug('观察keysWithSoundEffect=', keysWithSoundEffect.value);
//   }
// );

const keysWithSoundEffect = ref<Map<string, any>>(new Map());
watch(
  // TIPS: 注意, 如果要监听的ref对象是数组或对象等js/ts中的默认引用类型, 要使用此种方式才可触发监听。(或者也可以弃用ref, 改用reactive。)
  () => keysWithSoundEffect.value,
  (newVal) => {
    console.debug('观察keysWithSoundEffect=', keysWithSoundEffect.value);
  }
);

// Dependency validation logic
const dependencyIssues = ref<DependencyIssue[]>([]);

// Computed property to get all dependency issues
const allDependencyIssues = computed(() => {
  const audioFiles = soundFileList.value as AudioFile[];
  const sounds = soundList.value as Sound[];

  // Create a copy of keySoundList with original string format for validation
  const keySounds = keySoundList.value.map(keySound => {
    // Create a deep copy to avoid modifying the original
    const keySoundCopy = JSON.parse(JSON.stringify(keySound));

    // Restore original string format for dependency validation
    keySoundCopy.keySoundValue.down.value = keySoundCopy.keySoundValue.down.value.map((item: any) => {
      if (item.type === 'sounds' && item.value && typeof item.value === 'object' && item.value.soundKey) {
        return {
          type: 'sounds',
          value: item.value.soundKey
        };
      }
      if (item.type === 'key_sounds' && item.value && typeof item.value === 'object' && item.value.keySoundKey) {
        return {
          type: 'key_sounds',
          value: item.value.keySoundKey
        };
      }
      return item;
    });

    keySoundCopy.keySoundValue.up.value = keySoundCopy.keySoundValue.up.value.map((item: any) => {
      if (item.type === 'sounds' && item.value && typeof item.value === 'object' && item.value.soundKey) {
        return {
          type: 'sounds',
          value: item.value.soundKey
        };
      }
      if (item.type === 'key_sounds' && item.value && typeof item.value === 'object' && item.value.keySoundKey) {
        return {
          type: 'key_sounds',
          value: item.value.keySoundKey
        };
      }
      return item;
    });

    return keySoundCopy;
  }) as KeySound[];

  if (audioFiles.length === 0 && sounds.length === 0 && keySounds.length === 0) {
    return [];
  }

  const validator = createDependencyValidator(audioFiles, sounds, keySounds);

  // Only validate actual saved dependencies, not current UI selections
  // The global binding and single key bindings should come from saved album data, not current UI state

  // For now, we don't include global binding validation since it represents current UI selections
  // TODO: If there are actual saved global bindings, they should be included here
  const globalBinding = undefined;

  // Convert keysWithSoundEffect Map to the format expected by validator
  const singleKeyBindings = keysWithSoundEffect.value.size > 0
    ? keysWithSoundEffect.value
    : undefined;

  return validator.validateAllDependencies(globalBinding, singleKeyBindings);
});

// Update dependency issues when data changes
watch([soundFileList, soundList, keySoundList, keysWithSoundEffect], () => {
  dependencyIssues.value = allDependencyIssues.value;
}, { deep: true });

// Helper function to check if an item has dependency issues
const checkItemDependencyIssues = (itemType: 'audio_files' | 'sounds' | 'key_sounds', itemId: string) => {
  return hasItemDependencyIssues(itemType, itemId, dependencyIssues.value);
};
function convertValue(item: any) {
  if (item.type === 'audio_files') {
    return {
      type: 'audio_files',
      value: soundFileList.value.find(
        (soundFile) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id
      ),
    };
  }
  if (item.type === 'sounds') {
    return {
      type: 'sounds',
      // value: soundList.value.find((sound) => sound.soundKey === item.value.soundKey),
      value: soundList.value.find((sound) => sound.soundKey === item.value),
    };
  }
  if (item.type === 'key_sounds') {
    return {
      type: 'key_sounds',
      // value: keySoundList.value.find((keySound) => keySound.keySoundKey === item.value.keySoundKey),
      value: keySoundList.value.find((keySound) => keySound.keySoundKey === item.value),
    };
  }
  return null;
}

const isShowSingleKeySoundEffectEditDialog = ref(false);

const currentEditingKey = ref<number | null>(null);
let currentEditingKey_old = currentEditingKey.value; // TIPS: 此处旧值的记录, 在实际的使用逻辑中进行。
// TIPS: 如果currentEditingKey不变, 则watch就不会触发, 或者说这种做法无法实时记录正确的old值, 留此注释是警示下->不是所有场景都适合使用watch正确记录旧值的。
//       * 3 -> 2 后 old为3 ,  2 -> 3 后, old 为2。(不相同是符合预期的。)
//       * 3 -> 2 后 old为3 ,  2 -> 2 后, old 仍为3(但应该是2相同才对, 但只能判别不相同, 不符合预期, 或者说这种做法无法实时记录正确的old值)。
// watch(currentEditingKey, (newVal, oldVal) => {
//   currentEditingKey_old = oldVal;
// });
const currentEditingKeyOfName = computed(() => {
  return currentEditingKey.value !== null
    ? keyEvent_store.dikCodeToName.get(currentEditingKey.value) || 'Dik-{' + currentEditingKey.value + '}'
    : '';
});

// -- -- -- 编辑声效(重新选择声效)
const keyDownSingleKeySoundEffectSelect_edit = ref<any>();
const keyUpSingleKeySoundEffectSelect_edit = ref<any>();
let keyDownSingleKeySoundEffectSelect_edit_old = keyDownSingleKeySoundEffectSelect_edit.value;
let keyUpSingleKeySoundEffectSelect_edit_old = keyUpSingleKeySoundEffectSelect_edit.value;

const singleKeyTypeGroup_edit = ref<Array<string>>(['sounds']);
const keySingleKeySoundEffectOptions_edit = computed(() => {
  const List: Array<any> = [];
  if (singleKeyTypeGroup_edit.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (singleKeyTypeGroup_edit.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (singleKeyTypeGroup_edit.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  // ===== 应用自然排序：解决混合选项列表乱序问题 =====
  // 对混合类型的选项列表使用显示标签进行自然排序
  List.sort((a, b) => {
    const aLabel = album_options_select_label(a);
    const bLabel = album_options_select_label(b);
    return naturalSort(aLabel, bLabel);
  });
  console.debug('观察keySingleKeySoundEffectOptions_edit=', List);
  return List;
});
const isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit = computed(() => {
  return singleKeyTypeGroup_edit.value.includes('key_sounds');
});
const isAnchoringUltimatePerfectionKeySound_singleKey_edit = ref(true);
watch(keyDownSingleKeySoundEffectSelect_edit, (newVal, oldVal) => {
  console.debug('观察keyDownSingleKeySoundEffectSelect_edit=', keyDownSingleKeySoundEffectSelect_edit.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (
      isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit.value &&
      isAnchoringUltimatePerfectionKeySound_singleKey_edit.value
    ) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyUpSingleKeySoundEffectSelect_edit.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyUpSingleKeySoundEffectSelect_edit.value = keyDownSingleKeySoundEffectSelect_edit.value;
      }
    }
  }
});
watch(keyUpSingleKeySoundEffectSelect_edit, (newVal, oldVal) => {
  console.debug('观察keyUpSingleKeySoundEffectSelect_edit=', keyUpSingleKeySoundEffectSelect_edit.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (
      isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit.value &&
      isAnchoringUltimatePerfectionKeySound_singleKey_edit.value
    ) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyDownSingleKeySoundEffectSelect_edit.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyDownSingleKeySoundEffectSelect_edit.value = keyUpSingleKeySoundEffectSelect_edit.value;
      }
    }
  }
});

// 存储事件监听器的引用，以便后续移除
let messageAudioPackageListener: (e: MessageEvent) => void;

onBeforeMount(async () => {
  // 此时由于是新建键音包, 因此是没有对应配置文件, 需要我们主动去创建的。 故第二个参数设置为true
  // 这也是我们加载页面前必须确定的事情, 否则无法进行后续操作, 一切以配置文件为前提。
  const audioPkgPath = (await LoadConfig(props.pkgPath, props.isCreate)).audioPkgPath;

  // 如果是创建键音包, 则需要执行一定的初始化工作。
  if (props.isCreate) {
    await ConfigSet('package_name', $t('KeyToneAlbum.new.name.defaultValue'));

    await ConfigSet('audio_pkg_uuid', props.pkgPath);

    await ConfigSet('key_tone.is_enable_embedded_test_sound', isEnableEmbeddedTestSound);

    main_store.GetKeyToneAlbumList(); // 更新键音包选择列表的名称。

    setting_store.mainHome.selectedKeyTonePkg = audioPkgPath;
  }
  // 数据初始化
  await initData();
  // 将初始化数据的操作封装成一个函数, 并设置为异步函数, 以便使用await调用
  async function initData() {
    await ConfigGet('get_all_value').then((req) => {
      // console.debug('打印观察获取的值', req);
      if (req === false) {
        // 此时, 说明GetItem_sqlite请求过程中, 出错了, 因此需要错误通知, 并让用户重新启动, 防止用户因继续使用造成的存储设置被初始覆盖
        q.notify({
          type: 'negative',
          position: 'top',
          message: $t('KeyToneAlbum.notify.configFileReadFailed'),
          timeout: 100000,
        });
        return;
      }

      // TIPS: 由于采取各设置独立的录入即判别方式, 不再依赖整体的JSON字符串, 因此此if判断后续可能没必要存在(目前暂时保留)
      // 第一次进入本应用, 设置本就该是空的, 此时无需对我们的设置项进行任何操作, 也无需做任何通知。
      // 但为防止后续的JSON.parse报错, 因此此处也是必不可少的(因为只要非首次, 就不可能为空, watchEffect是立即执行的, 也就是说至少整体的结构是正常入库的)
      if (req === '' || req === '{}' || req === null) {
        return;
      }

      // // 若有设置数据, 则取出 TIPS: 注意, 这里的设置是直接读出的一个json对象, 而不是需要解析的json字符串
      // const settingStorage = JSON.parse(req);

      const data = req;

      // 键音包名称初始化。 (不过由于这里是新建键音包, 这个不出意外的话一开始是undefined
      if (data.package_name !== undefined) {
        pkgName.value = data.package_name;
      }

      // 已载入的声音文件列表初始化。  (不过由于这里是新建键音包, 这个不出意外的话一开始是undefine)
      if (data.audio_files !== undefined) {
        // keyTonePkgData.audio_files 是一个从后端获取的对象, 通过此方式可以简便的将其转换为数组, 数组元素为原对象中的key和value(增加了这两个key)
        const audioFilesArray = Object.entries(data.audio_files).map(([key, value]) => ({
          sha256: key,
          value: value,
        }));
        audioFiles.value = audioFilesArray;
        const tempSoundFileList: Array<any> = [];

        audioFiles.value.forEach((item) => {
          // 此处必须判断其是否存在, 否则会引起Object.entries报错崩溃, 影响后续流程执行。
          if (item.value.name !== undefined && item.value.name !== null) {
            Object.entries(item.value.name).forEach(([name_id, name]) => {
              tempSoundFileList.push({ sha256: item.sha256, name_id: name_id, name: name, type: item.value.type });
            });
          }
        });
        // ===== 应用自然排序：解决声音文件列表乱序问题 =====
        // 按文件名+扩展名进行自然排序（初始化时）
        tempSoundFileList.sort((a, b) => naturalSort(a.name + a.type, b.name + b.type));
        soundFileList.value = tempSoundFileList;
      }

      // 已载入的声音列表初始化。  (不过由于这里是新建键音包, 这个不出意外的话一开始是undefine)
      if (data.sound_list !== undefined) {
        const sounds = Object.entries(data.sounds).map(([key, value]) => ({
          soundKey: key,
          soundValue: value,
        }));
        // ===== 应用自然排序：解决声音列表乱序问题 =====
        // 基于声音名称或soundKey进行自然排序（初始化时）
        sounds.sort((a: any, b: any) => {
          const aName = (a.soundValue?.name as string) || a.soundKey;
          const bName = (b.soundValue?.name as string) || b.soundKey;
          return naturalSort(aName, bName);
        });
        soundList.value = sounds as Array<{
          soundKey: string;
          soundValue: {
            cut: { start_time: number; end_time: number; volume: number };
            name: string;
            source_file_for_sound: { sha256: string; name_id: string; type: string };
          };
        }>;
      }

      if (data.key_tone?.is_enable_embedded_test_sound !== undefined) {
        isEnableEmbeddedTestSound.down = data.key_tone.is_enable_embedded_test_sound.down;
        isEnableEmbeddedTestSound.up = data.key_tone.is_enable_embedded_test_sound.up;
      }

      // TODO: 此逻辑未验证, 需要到编辑键音包界面才能验证
      if (data.key_tone?.single !== undefined) {
        keysWithSoundEffect.value.clear();
        Object.entries(data.key_tone.single).forEach(([dikCode, value]) => {
          // 只有 down/up 至少一个被正确设置且value不为空字符串时, 才算作 已设置单键声效的按键。
          if ((value as any)?.down?.value || (value as any)?.up?.value) {
            keysWithSoundEffect.value.set(dikCode, value);
          }
        });
      }
    });
    const updateKeyToneAlbumListName = debounce(
      () => {
        main_store.GetKeyToneAlbumList();
      },
      800,
      { trailing: true }
    );
    watch(pkgName, (newVal) => {
      ConfigSet('package_name', pkgName.value);
      updateKeyToneAlbumListName.cancel();
      updateKeyToneAlbumListName();
    });

    // 2.配置文件中audio_files的进一步映射变更, 获取我们最终需要的结构
    watch(audioFiles, (newVal) => {
      console.debug('观察audioFiles=', audioFiles.value);
      // 为了更容易理解, 故引入audioFiles这一变量, 做初步映射, audioFiles只是过程值, 我们最终需要对此过程值做进一步映射, 形成soundFileList
      const tempSoundFileList: Array<any> = [];

      audioFiles.value.forEach((item) => {
        // 此处必须判断其是否存在, 否则会引起Object.entries报错崩溃, 影响后续流程执行。
        if (item.value.name !== undefined && item.value.name !== null) {
          Object.entries(item.value.name).forEach(([name_id, name]) => {
            tempSoundFileList.push({ sha256: item.sha256, name_id: name_id, name: name, type: item.value.type });
          });
        }
      });
      // ===== 应用自然排序：解决声音文件列表乱序问题 =====
      // 按文件名+扩展名进行自然排序，确保 sound1.wav, sound2.wav, sound10.wav 的正确顺序
      tempSoundFileList.sort((a, b) => naturalSort(a.name + a.type, b.name + b.type));
      soundFileList.value = tempSoundFileList;
    });

    // 3.观察进一步映射变更后, 最终需要的audio_file映射, 即我们的soundFileList。
    watch(soundFileList, (newVal) => {
      console.debug('观察soundFileList=', soundFileList.value);
    });

    //  - completed(已完成)   TODO:
    // 4.观察selectedSoundFile的变化, 当selectedSoundFile变化时,
    //   说明用户做了对应修改, 此时需要向sdk发送请求, 更新配置文件中的对应值, 然后触发sse形成闭环。
    //   当然, 删除时同理, 但删除是独立的按钮点击后手动触发对应函数, 以向sdk发送请求, 不由此处的数据驱动。
    watch(
      // TIPS: 对于ref的响应式变量, 如果直接整体监听, 则内部的某个值变化时, 不会触发监听。需要使用返回值的函数, 对固定字段进行监听。
      () => selectedSoundFile.value.name,
      (newVal) => {
        console.debug('观察selectedSoundFile=', selectedSoundFile.value);
        if (selectedSoundFile.value.sha256 !== '' && selectedSoundFile.value.name_id !== '') {
          SoundFileRename(
            selectedSoundFile.value.sha256,
            selectedSoundFile.value.name_id,
            selectedSoundFile.value.name
          );
        }
      }
    );

    watch(soundList, (newVal) => {
      console.debug('观察soundList=', soundList.value);
    });

    watch(keySoundList, (newVal) => {
      console.debug('观察keySoundList=', keySoundList.value);
    });

    watch(
      isEnableEmbeddedTestSound,
      (newVal) => {
        ConfigSet('key_tone.is_enable_embedded_test_sound', newVal);
      },
      { immediate: true }
    );
  }

  const pkgNameDelayed = debounce(
    (keyTonePkgData: any) => {
      pkgName.value = keyTonePkgData.package_name;
    },
    800,
    { trailing: true }
  );

  const isEnableEmbeddedTestSoundDelayed = debounce(
    (val: { down: boolean; up: boolean }) => {
      isEnableEmbeddedTestSound.down = val.down;
      isEnableEmbeddedTestSound.up = val.up;
    },
    800,
    { trailing: true }
  );

  // 将后端从键音包配置文件中获取的全部数据, 转换前端可用的键音包数据。(只要配置文件变更, 就会触发相关sse发送, 此处就会接收)
  function sseDataToKeyTonePkgData(keyTonePkgData: any) {
    // 键音包名称初始化。 (不过由于这里是新建键音包, 这个不出意外的话一开始是undefined
    if (keyTonePkgData.package_name !== undefined) {
      pkgNameDelayed.cancel();
      pkgNameDelayed(keyTonePkgData);
    }

    // 1. 初步映射配置文件中的audio_files到audioFiles。(只要配置文件变更, 就会触发相关sse发送, 此处就会接收)
    //    使用audioFiles作为中间值, 而不是一步到位的映射, 是为代码的可读性, 后续阅读理解是方便。
    if (keyTonePkgData.audio_files !== undefined) {
      // keyTonePkgData.audio_files 是一个从后端获取的对象, 通过此方式可以简便的将其转换为数组, 数组元素为原对象中的key和value(增加了这两个key)
      const audioFilesArray = Object.entries(keyTonePkgData.audio_files).map(([key, value]) => ({
        sha256: key,
        value: value,
      }));
      audioFiles.value = audioFilesArray;
    } else {
      // 此处else是为防止最后一项的audio_files为undefined, 而导致的删除最后一项音频源文件后, audioFiles值无法清空, 从而导致无法触发soundFileList的变更, 从而ui界面导致无法删除最后一项音频源文件。
      audioFiles.value = [];
    }

    // 映射配置文件中的sounds到ui中的soundList。(只要配置文件变更, 就会触发相关sse发送, 此处就会接收)
    if (keyTonePkgData.sounds !== undefined) {
      const sounds = Object.entries(keyTonePkgData.sounds).map(([key, value]) => ({
        soundKey: key,
        soundValue: value,
      }));
      // ===== 应用自然排序：解决声音列表乱序问题 =====
      // 基于声音名称或soundKey进行自然排序（EventSource更新时）
      sounds.sort((a: any, b: any) => {
        const aName = (a.soundValue?.name as string) || a.soundKey;
        const bName = (b.soundValue?.name as string) || b.soundKey;
        return naturalSort(aName, bName);
      });
      soundList.value = sounds as Array<{
        soundKey: string;
        soundValue: {
          cut: { start_time: number; end_time: number; volume: number };
          name: string;
          source_file_for_sound: { sha256: string; name_id: string; type: string };
        };
      }>;
    } else {
      soundList.value = [];
    }

    // 映射配置文件中的key_sounds到ui中的keySoundList。(只要配置文件变更, 就会触发相关sse发送, 此处就会接收)
    if (keyTonePkgData.key_sounds !== undefined) {
      const keySounds = Object.entries(keyTonePkgData.key_sounds).map(([key, value]) => ({
        keySoundKey: key,
        keySoundValue: value,
      }));
      // ===== 应用自然排序：解决键音列表乱序问题 =====
      // 基于键音名称或keySoundKey进行自然排序
      keySounds.sort((a: any, b: any) => {
        const aName = (a.keySoundValue?.name as string) || a.keySoundKey;
        const bName = (b.keySoundValue?.name as string) || b.keySoundKey;
        return naturalSort(aName, bName);
      });
      keySoundList.value = keySounds;
    } else {
      keySoundList.value = [];
    }

    if (keyTonePkgData.key_tone !== undefined) {
      isEnableEmbeddedTestSoundDelayed.cancel();
      isEnableEmbeddedTestSoundDelayed(keyTonePkgData.key_tone.is_enable_embedded_test_sound);
    }

    if (keyTonePkgData.key_tone?.global !== undefined) {
      keyDownUnifiedSoundEffectSelect.value = convertValue(
        keyTonePkgData.key_tone.global.down ? keyTonePkgData.key_tone.global.down : ''
      );
      keyUpUnifiedSoundEffectSelect.value = convertValue(
        keyTonePkgData.key_tone.global.up ? keyTonePkgData.key_tone.global.up : ''
      );
    }

    if (keyTonePkgData.key_tone?.single !== undefined) {
      keysWithSoundEffect.value.clear();
      Object.entries(keyTonePkgData.key_tone.single).forEach(([dikCode, value]) => {
        // 只有 down/up 至少一个被正确设置且value不为空字符串时, 才算作 已设置单键声效的按键。
        if ((value as any)?.down?.value || (value as any)?.up?.value) {
          keysWithSoundEffect.value.set(dikCode, value);
        }
      });
    }
  }
  const debounced_sseDataToSettingStore = debounce<(keyTonePkgData: any) => void>(sseDataToKeyTonePkgData, 30, {
    trailing: true,
  });

  // 定义事件监听器
  messageAudioPackageListener = function (e) {
    console.debug('后端钩子函数中的值 = ', e.data);

    const data = JSON.parse(e.data);

    if (data.key === 'get_all_value') {
      debounced_sseDataToSettingStore.cancel;
      debounced_sseDataToSettingStore(data.value);
    }
  };

  // 添加事件监听
  app_store.eventSource.addEventListener('messageAudioPackage', messageAudioPackageListener, false);
});

// 在退出创建键音包的页面后, 载入 持久化的 用户选择的 键音包。(在 创建 键音包界面 退出时, 重新加载 用户持久化至 设置 文件中的 键音包。)
const main_store = useMainStore();
onUnmounted(() => {
  // 卸载组件后, 更新键音包列表
  main_store.GetKeyToneAlbumList();
  // // 卸载组件后, 重新载入持久化配置中用户所选的键音包(在新设计的键音专辑页面逻辑中, 不需要此步骤)
  // main_store.LoadSelectedKeyTonePkg();

  // 移除事件监听
  if (messageAudioPackageListener) {
    app_store.eventSource.removeEventListener('messageAudioPackage', messageAudioPackageListener);
  }
});

const isMacOS = ref(getMacOSStatus());
function getMacOSStatus() {
  if (process.env.MODE === 'electron') {
    return window.myWindowAPI.getMacOSStatus();
  }
  return false;
}

const i18n_fontSize = computed(() => {
  return setting_store.languageDefault === 'ru' ||
    setting_store.languageDefault === 'it' ||
    setting_store.languageDefault === 'es' ||
    setting_store.languageDefault === 'pt-BR' ||
    setting_store.languageDefault === 'pl' ||
    setting_store.languageDefault === 'tr' ||
    setting_store.languageDefault === 'id'
    ? '0.66rem'
    : setting_store.languageDefault === 'fr'
    ? '0.63rem'
    : '0.75rem';
});

const step_introduce_fontSize = computed(() => {
  return isMacOS.value
    ? // MacOS
      setting_store.languageDefault === 'ru' ||
      setting_store.languageDefault === 'ja' ||
      setting_store.languageDefault === 'es'
      ? 'text-[0.80rem]'
      : setting_store.languageDefault === 'ko-KR'
      ? 'text-[0.83rem]'
      : setting_store.languageDefault === 'id'
      ? 'text-[0.87]'
      : 'text-[0.85rem]'
    : // windows
    setting_store.languageDefault === 'ru' ||
      setting_store.languageDefault === 'ko-KR' ||
      setting_store.languageDefault === 'pl' ||
      setting_store.languageDefault === 'ar'
    ? 'text-[0.80rem]'
    : setting_store.languageDefault === 'tr'
    ? 'text-[0.83rem]'
    : 'text-[0.85rem]';
});

// ============================================================================
// Context 提供（provide）
// ============================================================================
//
// 【作用】
// 将父组件的所有状态和方法打包成 Context 对象，通过 provide 提供给子组件。
// 子组件通过 inject(KEYTONE_ALBUM_CONTEXT_KEY) 获取这个 Context。
//
// 【为什么在这里定义】
// 需要在所有状态和方法都定义完成后才能构建 Context。
// 这样确保 Context 中引用的所有变量都已存在。
//
// 【子组件如何使用】
// 在子组件中：const ctx = inject<KeytoneAlbumContext>(KEYTONE_ALBUM_CONTEXT_KEY)!;
// 然后通过 ctx.step.value、ctx.$t() 等方式访问状态和方法。
// 注意：在模板中访问 Ref 类型需要用 .value，因为 ctx 是普通对象而非 Ref。
//
// ============================================================================

/**
 * 构建 KeytoneAlbumContext 对象
 *
 * 这里只列出 Step1 (音频源文件) 相关的状态和方法，
 * 其他 Step 的状态会在后续迁移时逐步添加。
 *
 * 当前策略：先提供必要的状态，验证机制正常后再扩展。
 */
const keytoneAlbumContext: KeytoneAlbumContext = {
  // ============ Props ============
  // 从 props 获取，子组件可以读取但不能修改
  pkgPath: props.pkgPath,
  isCreate: props.isCreate,

  // ============ 核心状态 ============
  step,           // 当前步骤（1-4 或 99 表示折叠）
  pkgName,        // 键音包名称

  // ============ Step1: 音频源文件相关 ============
  addNewSoundFile,    // 控制"添加音频文件"对话框的 v-model
  files,              // 待上传的文件列表
  editSoundFile,      // 控制"管理音频文件"对话框的 v-model
  soundFileList,      // 已加载的音频文件列表
  selectedSoundFile,  // 当前选中的音频文件（用于编辑/删除）

  // ============ Step2: 声音定义相关 ============
  createNewSound,     // 控制"创建声音"对话框
  soundName,          // 声音名称
  sourceFileForSound, // 声音的源文件引用
  soundStartTime,     // 声音裁剪开始时间
  soundEndTime,       // 声音裁剪结束时间
  soundVolume,        // 声音音量
  showEditSoundDialog,// 控制"编辑声音"对话框
  soundList,          // 已定义的声音列表
  selectedSound,      // 当前选中的声音

  // ============ Step3: 按键音相关 ============
  createNewKeySound,      // 控制"创建按键音"对话框
  keySoundName,           // 按键音名称
  configureDownSound,     // 是否配置按下声音
  configureUpSound,       // 是否配置抬起声音
  selectedSoundsForDown,  // 按下时选中的声音
  playModeForDown,        // 按下播放模式
  maxSelectionForDown,    // 按下最大选择数
  downTypeGroup,          // 按下类型组
  downSoundList,          // 按下声音列表
  selectedSoundsForUp,    // 抬起时选中的声音
  playModeForUp,          // 抬起播放模式
  maxSelectionForUp,      // 抬起最大选择数
  upTypeGroup,            // 抬起类型组
  upSoundList,            // 抬起声音列表
  editExistingKeySound,   // 控制"编辑按键音"对话框
  edit_configureDownSound,
  edit_configureUpSound,
  edit_downTypeGroup,
  edit_upTypeGroup,
  edit_downSoundList,
  edit_upSoundList,
  keySoundList,           // 已定义的按键音列表
  selectedKeySound,       // 当前选中的按键音

  // ============ Step4: 联动声效相关 ============
  isEnableEmbeddedTestSound,  // 内嵌测试音开关
  showEveryKeyEffectDialog,   // 全键声效对话框
  keyDownUnifiedSoundEffectSelect,
  keyUpUnifiedSoundEffectSelect,
  unifiedTypeGroup,
  keyUnifiedSoundEffectOptions,
  isShowUltimatePerfectionKeySoundAnchoring,
  isAnchoringUltimatePerfectionKeySound,
  showSingleKeyEffectDialog,
  isShowAddOrSettingSingleKeyEffectDialog,
  selectedSingleKeys,
  isRecordingSingleKeys,
  keyOptions,
  filterOptions,
  isGetsFocused,
  isDownSoundEffectSelectEnabled,
  isUpSoundEffectSelectEnabled,
  keyDownSingleKeySoundEffectSelect,
  keyUpSingleKeySoundEffectSelect,
  singleKeyTypeGroup,
  keySingleKeySoundEffectOptions,
  isShowUltimatePerfectionKeySoundAnchoring_singleKey,
  isAnchoringUltimatePerfectionKeySound_singleKey,
  keysWithSoundEffect,
  isShowSingleKeySoundEffectEditDialog,
  currentEditingKey,
  currentEditingKey_old,
  currentEditingKeyOfName,
  keyDownSingleKeySoundEffectSelect_edit,
  keyUpSingleKeySoundEffectSelect_edit,
  keyDownSingleKeySoundEffectSelect_edit_old,
  keyUpSingleKeySoundEffectSelect_edit_old,
  singleKeyTypeGroup_edit,
  keySingleKeySoundEffectOptions_edit,
  isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit,
  isAnchoringUltimatePerfectionKeySound_singleKey_edit,

  // ============ 依赖校验 ============
  dependencyIssues,

  // ============ 工具函数 ============
  album_options_select_label,   // 选项标签显示函数
  naturalSort,                  // 自然排序函数
  preventDefaultKeyBehaviorWhenRecording,  // 防止录制时键盘默认行为
  preventDefaultMouseWhenRecording,        // 防止录制时鼠标默认行为
  setSingleKeyRecordingClearFlag,          // 单键录制：设置 clear_flag
  convertValue,                 // 值转换函数

  // ============ 操作函数 ============
  saveSoundConfig,      // 保存声音配置
  deleteSound,          // 删除声音
  previewSound,         // 预览声音
  saveKeySoundConfig,   // 保存按键音配置
  deleteKeySound,       // 删除按键音
  saveUnifiedSoundEffectConfig,     // 保存全局联动声效
  saveSingleKeySoundEffectConfig,   // 保存单键联动声效

  // ============ i18n ============
  $t,  // 国际化翻译函数

  // ============ 样式相关 ============
  i18n_fontSize,           // 按语言调整的字体大小
  step_introduce_fontSize, // 步骤说明文字大小
  isMacOS,                 // 是否 MacOS 平台

  // ============ 选项常量 ============
  options,           // 类型选项（audio_files/sounds/key_sounds）
  playModeOptions,   // 播放模式选项
  playModeLabels,    // 播放模式标签映射
};

// 提供 Context 给子组件
// 子组件通过 inject(KEYTONE_ALBUM_CONTEXT_KEY) 获取
provide(KEYTONE_ALBUM_CONTEXT_KEY, keytoneAlbumContext);
</script>

<style lang="scss" scoped>
// :deep(.q-stepper__tab) {
//   cursor: pointer;
// }

// :deep(.q-stepper__tab:hover) {
//   background-color: rgb(243 244 246);
// }

// TIPS: 注意 unocss 的默认预设中, 默认情况下是不支持以下这种特定的 tailwindcss 语法的 - 即 @apply 的用法不受支持。
//       * 需要通过手动更改相应的 配置文件uno.config.ts 来支持此转换语法。 -[参考链接](https://unocss.jiangruyi.com/transformers/directives)
:deep(.q-stepper__tab) {
  @apply cursor-pointer hover:bg-gray-100;
}

// 对选择器组件的label溢出情况, 采取滚动策略 (为防止刷新后样式丢失问题, 而加的。)
:deep(.q-field__native) {
  // 对溢出的情况, 采取滚动策略
  @apply max-w-full overflow-auto whitespace-nowrap;

  // 隐藏滚动策略的滚动条。
  // @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];

  // 添加细微滚动条
  @apply h-5.8 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400;
}

// // 对于多选的, 带芯片的选择框, 某芯片内的名称内容过长的情况, 采取溢出滚动的策略。
// :deep(.ellipsis) {
//   // 对溢出的情况, 采取滚动策略
//   @apply max-w-full overflow-auto whitespace-nowrap  text-clip;
//   // 隐藏滚动策略的滚动条。
//   @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];
// }

// 对本组件选择框添加的可清楚图标的大小做设置
:deep(.q-field__focusable-action) {
  @apply text-lg;
}

// 用于修复主页面全局的:global(.q-field__native)中的h-5.8这个样式影响了当前页面中的q-input的问题
:deep(.q-placeholder) {
  // 在这里重置q-input组件的输入样式的高度以修复这个问题
  @apply h-auto;
}

// 为防止刷新后样式丢失问题, 而加的。
:global(.q-item__section) {
  /* 对溢出的情况, 采取滚动策略 */
  @apply max-w-full overflow-auto whitespace-nowrap;

  /* 隐藏滚动策略的滚动条 */
  // @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];

  // 添加细微滚动条
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30  [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}

.zl-ll {
  :deep(.q-field__native) {
    @apply h-auto;
  }
  :deep(.q-field__messages) {
    @apply text-nowrap;
  }
}

:global(.q-card) {
  @apply mr-2.33;
}

// 对于键音专辑组件的选择框, 键音专辑的名称内容过长的情况, 采取溢出滚动的策略。
// 对于多选的, 带芯片的选择框, 某芯片内的名称内容过长的情况, 采取溢出滚动的策略。
:deep(.ellipsis) {
  // 对溢出的情况, 采取滚动策略
  @apply max-w-full overflow-auto whitespace-nowrap  text-clip;
  // // 隐藏滚动策略的滚动条。
  // @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];

  // 添加细微滚动条
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30  [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}

.q-btn {
  @apply text-xs;
  // font-size: 0.66rem;
  // line-height: 1rem;
  font-size: var(--i18n_fontSize);
  @apply p-1.5;
  @apply transition-transform hover:scale-105;
  @apply scale-103;
}

.step-custom {
  :deep(.q-stepper__step-inner) {
    padding: 0 10px 32px 55px;
  }
}
:deep(.q-field__label) {
  @apply overflow-visible -ml-1.5 text-[0.8rem];
}
</style>
