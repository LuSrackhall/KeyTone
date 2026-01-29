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
  <q-item :class="['h-15 mb-5']">
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-[80%] self-center']"></div>

    <div :class="['w-[100%] grid']">
      <div :class="['w-[92%] flex justify-between items-center flex-nowrap  gap-[12px]']">
        <q-input
          dense
          hide-bottom-space
          :class="['w-[66%] h-10.5 ']"
          v-model.number="setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope"
          type="number"
          filled
          :label="$t('setting.mainHome.音量降幅.index')"
          stack-label
          :rules="[(val: number) => { return val >= 5 && val<100000000 || $t('setting.mainHome.音量降幅.rulesErrorInfo'); }]"
        />
        <q-btn
          :class="['min-w-15 min-h-5']"
          color="primary"
          size="10px"
          :label="$t('setting.mainHome.重置')"
          @click="returnToNormalReduceScope()"
        />
      </div>
    </div>
  </q-item>

  <q-item>
    <!-- 左边的竖线 -->
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
    <q-item-section>
      <q-item-label>{{ $t('setting.mainHome.音量调试滑块.index') }}</q-item-label>
      <q-item-label caption>{{ $t('setting.mainHome.音量调试滑块.caption') }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <q-toggle v-model="setting_store.mainHome.audioVolumeProcessing.isOpenVolumeDebugSlider" />
    </q-item-section>
  </q-item>

  <!-- =============================
       分离模式音量设置（仅在主页面开启分离时生效）
       =============================
       说明：这里不提供分离开关，仅提供分离音量设置
  -->
  <!--
    重新设计的“分离模式音量设置”展开样式：
    - 使用自定义 header + q-slide-transition 代替内置 q-expansion-item
    - 视觉更轻量，避免展开/收起时的延迟感
    - 保留左侧竖条与缩进一致性
  -->
  <!-- NOTE: icon="tune" 很适合该栏目，但根据你的要求此处不显示图标。
       保留注释方便后续快速恢复图标样式。 -->
  <!-- icon="tune" -->
  <!-- 默认不展示白色背景，仅在 hover/点击时呈现 -->
  <q-item
    clickable
    class="rounded-lg border border-zinc-200/60 shadow-sm transition-colors hover:bg-zinc-50/70 active:bg-zinc-50/70"
    @click="splitVolumeOpen = !splitVolumeOpen"
  >
    <!-- 展开栏头部保持与其它设置一致：缩进 + 左侧竖条 -->
    <div :class="['ml-6 rounded-full border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
    <q-item-section>
      <q-item-label>{{ $t('setting.mainHome.splitVolumePanel.title') }}</q-item-label>
      <q-item-label caption>{{ $t('setting.mainHome.splitVolumePanel.caption') }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <!--
        展开箭头：仅做状态提示，不改变任何功能逻辑
        通过切换图标来表达展开/收起
      -->
      <q-btn
        dense
        round
        flat
        :icon="splitVolumeOpen ? 'expand_less' : 'expand_more'"
      />
    </q-item-section>
  </q-item>

  <!--
    展开内容：使用 q-slide-transition
    - 仅控制显示/隐藏
    - 不改变内部滑块与布局逻辑
  -->
  <!-- 使用轻量过渡：只做 opacity + transform，不做高度动画，避免“停顿感” -->
  <transition name="split-fade">
    <div
      v-show="splitVolumeOpen"
      class="split-volume-content ml-6 mr-4 mt-2 px-2 py-2 rounded-lg bg-zinc-50/70 border border-zinc-200/50 shadow-sm relative z-0"
    >
      <!-- 键盘滑块优先展示 -->
      <q-item>
        <!-- 不在模板中使用 `!important` 类名，统一在样式层处理溢出规则 -->
        <q-item-section>
          <!-- 仅保留左侧“键盘”字样；滑块行不使用竖条与额外缩进 -->
          <!-- overflow-visible：避免全局 q-item__section 的 overflow 影响滑块标签显示 -->
          <!--
            统一行内对齐：把预留空间放到整行（pt-3），
            这样“键盘/鼠标”字样与滑块保持同一水平基线。
          -->
          <!--
            顶部留白恢复为轻量值：
            - 只保留必要的空间
            - 避免因为过宽的垂直高度影响布局紧凑度
          -->
          <!-- 行内容提升层级，避免被背景覆盖 -->
          <div :class="['flex items-center gap-3 pt-3 pb-1 min-h-[52px] relative z-10', 'overflow-visible']">
            <div class="text-xs text-gray-600 min-w-8">{{ $t('setting.mainHome.splitVolume.keyboard.shortLabel') }}</div>
            <!--
              为“滑块标签（百分比）”预留纵向空间：
              - 主滑块的 label 仅在拖动时显示（与主页面一致）
              - 标签出现时会占用 slider 上方空间，因此在整行预留 pt-3
            -->
            <div :class="['w-56 flex justify-between items-center', 'overflow-visible']">
              <q-btn dense round flat :icon="keyboardVolumeIcon" @click="toggleKeyboardSilent" />
              <q-slider
                :class="['w-[80%]']"
                v-model="setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormal"
                :max="0"
                :min="-keyboardMin"
                :step="0"
                label
                :label-value="keyboardLabelValue"
                color="light-green"
              />
            </div>
          </div>

          <!-- 键盘调试滑块（紧挨键盘滑块下方） -->
          <!-- 调试滑块同样预留垂直空间，避免数字只显示一半 -->
          <!-- 行内容提升层级，避免被背景覆盖 -->
          <div v-if="keyboardDebugOpen" :class="['flex items-center gap-3 mt-2 pt-4 pb-1 min-h-[56px] relative z-10', 'overflow-visible']">
            <div class="min-w-8"></div>
            <!--
              调试滑块：默认常显数字（label-always），与主页面调试滑块一致。
              同样在整行预留 pt-3 空间，避免 label 覆盖或被上方内容挤压。
            -->
            <div :class="['w-56 flex justify-end items-center', 'overflow-visible']">
              <q-slider
                :class="['w-[80%]']"
                v-model="setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormal"
                :max="0"
                :min="-keyboardMin"
                :step="0"
                :markers="keyboardMarkersDebug"
                marker-labels
                label
                label-always
                :label-value="keyboardLabelValueDebug"
                color="light-green"
              />
            </div>
          </div>
        </q-item-section>
      </q-item>

      <!-- 鼠标滑块优先展示（紧随键盘之后） -->
      <q-item>
        <q-item-section>
          <!-- 仅保留左侧“鼠标”字样；滑块行不使用竖条与额外缩进 -->
          <!-- overflow-visible：避免全局 q-item__section 的 overflow 影响滑块标签显示 -->
          <!--
            同键盘主滑块：顶部留白恢复为轻量值
            - 仅保留必要空间，避免垂直高度过大
          -->
          <!-- 行内容提升层级，避免被背景覆盖 -->
          <div :class="['flex items-center gap-3 pt-3 pb-1 min-h-[52px] relative z-10', 'overflow-visible']">
            <div class="text-xs text-gray-600 min-w-8">{{ $t('setting.mainHome.splitVolume.mouse.shortLabel') }}</div>
            <!-- 同键盘：主滑块 label 仅拖动时显示，但要预留纵向空间避免遮挡 -->
            <div :class="['w-56 flex justify-between items-center', 'overflow-visible']">
              <q-btn dense round flat :icon="mouseVolumeIcon" @click="toggleMouseSilent" />
              <q-slider
                :class="['w-[80%]']"
                v-model="setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormal"
                :max="0"
                :min="-mouseMin"
                :step="0"
                label
                :label-value="mouseLabelValue"
                color="light-green"
              />
            </div>
          </div>

          <!-- 鼠标调试滑块（紧挨鼠标滑块下方） -->
          <!-- 调试滑块同样预留垂直空间，避免数字只显示一半 -->
          <!-- 行内容提升层级，避免被背景覆盖 -->
          <div v-if="mouseDebugOpen" :class="['flex items-center gap-3 mt-2 pt-4 pb-1 min-h-[56px] relative z-10', 'overflow-visible']">
            <div class="min-w-8"></div>
            <!-- 调试滑块：label-always 常显；并在整行预留 pt-3 空间 -->
            <div :class="['w-56 flex justify-end items-center', 'overflow-visible']">
              <q-slider
                :class="['w-[80%]']"
                v-model="setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormal"
                :max="0"
                :min="-mouseMin"
                :step="0"
                :markers="mouseMarkersDebug"
                marker-labels
                label
                label-always
                :label-value="mouseLabelValueDebug"
                color="light-green"
              />
            </div>
          </div>
        </q-item-section>
      </q-item>

      <!-- 降幅设置：保留左侧竖条但不额外缩进 -->
      <q-item :class="['h-15 mb-2']">
        <div :class="['ml-6 rounded-full border-l-solid border-l-5 mr-6 h-[80%] self-center']"></div>
        <div :class="['w-[100%] grid']">
          <div :class="['w-[92%] flex justify-between items-center flex-nowrap gap-[12px]']">
            <q-input
              dense
              hide-bottom-space
              :class="['w-[66%] h-10.5 ']"
              v-model.number="setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormalReduceScope"
              type="number"
              filled
              :label="$t('setting.mainHome.splitVolume.keyboard.reduceScope')"
              stack-label
              :rules="[(val: number) => { return val >= 5 && val<100000000 || $t('setting.mainHome.音量降幅.rulesErrorInfo'); }]"
            />

            <q-btn
              :class="['min-w-15 min-h-5']"
              color="primary"
              size="10px"
              :label="$t('setting.mainHome.重置')"
              @click="returnKeyboardReduceScope()"
            />
          </div>
        </div>
      </q-item>

      <q-item :class="['h-15 mb-2']">
        <div :class="['ml-6 rounded-full border-l-solid border-l-5 mr-6 h-[80%] self-center']"></div>
        <div :class="['w-[100%] grid']">
          <div :class="['w-[92%] flex justify-between items-center flex-nowrap gap-[12px]']">
            <q-input
              dense
              hide-bottom-space
              :class="['w-[66%] h-10.5 ']"
              v-model.number="setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormalReduceScope"
              type="number"
              filled
              :label="$t('setting.mainHome.splitVolume.mouse.reduceScope')"
              stack-label
              :rules="[(val: number) => { return val >= 5 && val<100000000 || $t('setting.mainHome.音量降幅.rulesErrorInfo'); }]"
            />

            <q-btn
              :class="['min-w-15 min-h-5']"
              color="primary"
              size="10px"
              :label="$t('setting.mainHome.重置')"
              @click="returnMouseReduceScope()"
            />
          </div>
        </div>
      </q-item>

      <!-- 调试开关：最后展示，保留竖条 -->
      <q-item>
        <div :class="['ml-6 rounded-full border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
        <q-item-section>
          <q-item-label>{{ $t('setting.mainHome.splitVolume.keyboard.debugSlider') }}</q-item-label>
        </q-item-section>
        <q-item-section side>
          <q-toggle v-model="setting_store.mainHome.splitAudioVolumeProcessing.keyboard.isOpenVolumeDebugSlider" />
        </q-item-section>
      </q-item>

      <q-item>
        <div :class="['ml-6 rounded-full border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
        <q-item-section>
          <q-item-label>{{ $t('setting.mainHome.splitVolume.mouse.debugSlider') }}</q-item-label>
        </q-item-section>
        <q-item-section side>
          <q-toggle v-model="setting_store.mainHome.splitAudioVolumeProcessing.mouse.isOpenVolumeDebugSlider" />
        </q-item-section>
      </q-item>
  </div>
  </transition>

  <!-- =============================
       鼠标回退到键盘开关
       =============================
       功能说明：
       - 仅在"键盘/鼠标分离"模式下生效
       - 默认关闭：分离模式下鼠标无专辑则无专辑, 彻底分离
       - 开启后：鼠标专辑缺失时会回退到键盘专辑（复用键盘配置）
  -->
  <q-item>
    <div :class="['ml-6 rounded-full border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
    <q-item-section>
      <q-item-label>{{ $t('setting.mainHome.mouseFallback.index') }}</q-item-label>
      <q-item-label caption>{{ $t('setting.mainHome.mouseFallback.caption') }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <q-toggle v-model="setting_store.playbackRouting.mouseFallbackToKeyboard" />
    </q-item-section>
  </q-item>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { useSettingStore } from 'src/stores/setting-store';

const q = useQuasar();
const { t } = useI18n();
const setting_store = useSettingStore();

// 自定义展开状态（替代 q-expansion-item 的内置状态）
const splitVolumeOpen = ref(false);

const returnToNormalReduceScope = () => {
  setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope = 5.0;
};

// 计算键盘/鼠标独立音量的最小值（与主页算法一致）
const keyboardMin = computed(() => {
  if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
    return (
      setting_store.audioVolumeProcessing.volumeAmplify +
      setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormalReduceScope
    );
  }
  return setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormalReduceScope;
});

const mouseMin = computed(() => {
  if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
    return (
      setting_store.audioVolumeProcessing.volumeAmplify +
      setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormalReduceScope
    );
  }
  return setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormalReduceScope;
});

const keyboardLabelValue = computed(() => {
  const percentage = (
    (1 - -setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormal / keyboardMin.value) *
    100
  )
    .toFixed(2)
    .split('.');
  return percentage[1] === '00' ? percentage[0] + '%' : percentage[0] + '.' + percentage[1] + '%';
});

const mouseLabelValue = computed(() => {
  const percentage = (
    (1 - -setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormal / mouseMin.value) * 100
  )
    .toFixed(2)
    .split('.');
  return percentage[1] === '00' ? percentage[0] + '%' : percentage[0] + '.' + percentage[1] + '%';
});

const keyboardLabelValueDebug = computed(() => {
  return setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormal.toFixed(2);
});

const mouseLabelValueDebug = computed(() => {
  return setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormal.toFixed(2);
});

// 图标状态：与主页面一致（静音/非静音）
const keyboardVolumeIcon = computed(() => {
  return setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeSilent ? 'volume_off' : 'volume_up';
});

const mouseVolumeIcon = computed(() => {
  return setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeSilent ? 'volume_off' : 'volume_up';
});

// 调试滑块开关（用于控制 UI 显示与竖条高度）
const keyboardDebugOpen = computed(() => {
  return setting_store.mainHome.splitAudioVolumeProcessing.keyboard.isOpenVolumeDebugSlider;
});

const mouseDebugOpen = computed(() => {
  return setting_store.mainHome.splitAudioVolumeProcessing.mouse.isOpenVolumeDebugSlider;
});

const keyboardMarkersDebug = computed(() => {
  if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
    return (
      (setting_store.audioVolumeProcessing.volumeAmplify +
        setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormalReduceScope) /
      1
    );
  }
  return setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormalReduceScope / 1;
});

const mouseMarkersDebug = computed(() => {
  if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
    return (
      (setting_store.audioVolumeProcessing.volumeAmplify +
        setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormalReduceScope) /
      1
    );
  }
  return setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormalReduceScope / 1;
});

// 当 reduce scope 改变时，确保当前音量不会超出范围
watch(keyboardMin, () => {
  if (-setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormal > keyboardMin.value) {
    setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormal = -keyboardMin.value;
  }
});

watch(mouseMin, () => {
  if (-setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormal > mouseMin.value) {
    setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormal = -mouseMin.value;
  }
});

// 与主页面一致：当音量变化时自动同步静音状态
watch(
  () => setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormal,
  () => {
    setTimeout(() => {
      if (keyboardLabelValue.value !== '0%') {
        setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeSilent = false;
      } else {
        setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeSilent = true;
      }
    }, 60);
  }
);

watch(
  () => setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormal,
  () => {
    setTimeout(() => {
      if (mouseLabelValue.value !== '0%') {
        setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeSilent = false;
      } else {
        setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeSilent = true;
      }
    }, 60);
  }
);

const returnKeyboardReduceScope = () => {
  setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeNormalReduceScope = 5.0;
};

const returnMouseReduceScope = () => {
  setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeNormalReduceScope = 5.0;
};

// 静音切换（分离模式下独立静音，与主页面全局静音叠加生效）
const toggleKeyboardSilent = () => {
  // 与主页面一致：当音量为 0% 时，不允许通过静音图标“开启声音”
  if (keyboardLabelValue.value === '0%') {
    q.notify({
      message: t('Notify.音量0%时无法打开声音'),
      color: 'warning',
      position: 'top',
      timeout: 1200,
    });
    return;
  }
  setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeSilent =
    !setting_store.mainHome.splitAudioVolumeProcessing.keyboard.volumeSilent;
};

const toggleMouseSilent = () => {
  // 与主页面一致：当音量为 0% 时，不允许通过静音图标“开启声音”
  if (mouseLabelValue.value === '0%') {
    q.notify({
      message: t('Notify.音量0%时无法打开声音'),
      color: 'warning',
      position: 'top',
      timeout: 1200,
    });
    return;
  }
  setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeSilent =
    !setting_store.mainHome.splitAudioVolumeProcessing.mouse.volumeSilent;
};
</script>

<style lang="scss" scoped>
// 用于修复主页面全局的:global(.q-field__native)中的h-5.8这个样式影响了当前页面中的q-input的问题
:deep(.q-placeholder) {
  // 在这里重置q-input组件的输入样式的高度以修复这个问题
  @apply h-auto;
}

// 展开内容轻量过渡：只做 opacity + transform，不做高度动画
.split-volume-content {
  will-change: opacity, transform;
  transform: translateZ(0);
}

// Vue transition：进入/离开时只改变透明度与轻微位移
.split-fade-enter-active,
.split-fade-leave-active {
  transition: opacity 120ms ease, transform 120ms ease;
}

.split-fade-enter-from,
.split-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px) scaleY(0.98);
}

.split-fade-enter-to,
.split-fade-leave-from {
  opacity: 1;
  transform: translateY(0) scaleY(1);
}

:deep(.q-item__section) {
  // 关键修复：主滑块百分比 label 被截的根因是 overflow: hidden。
  //
  // 为什么这里必须用 !important：
  // - 该页面与其它页面/组件可能都在各自的 scoped 样式中声明了 `.q-item__section { overflow: hidden; }`
  // - scoped 规则最终会变成 `[data-v-xxx] .q-item__section { ... }`，不同组件的选择器特异性相同
  // - 当其它组件的 scoped 样式后加载/优先级更高时，会把本页面的 overflow 覆盖回 hidden
  //
  // 因此此处用 `overflow: visible !important` 作为“可见性兜底”，只修复裁剪，不改 label 的视觉样式。
  text-wrap: wrap;
  overflow: visible !important;
}

// 展开项内容区/容器同样可能隐藏溢出，这里统一放开
:deep(.q-expansion-item__content),
:deep(.q-expansion-item__container) {
  overflow: visible !important;
}

// slider 与 label 容器允许溢出，避免被父层截断（不改变 label 样式）
:deep(.q-slider),
:deep(.q-slider__track-container),
:deep(.q-slider__markers),
:deep(.q-slider__label-container) {
  overflow: visible !important;
}
</style>
