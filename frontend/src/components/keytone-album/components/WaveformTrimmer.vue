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
文件说明: components/WaveformTrimmer.vue - 波形裁剪组件（MVP）
============================================================================

【文件作用】
- 展示选中音频源文件的波形（前端解码渲染）。
- 支持通过拖拽 region 选择裁剪区间。
- 与外部 startMs/endMs（毫秒）双向同步。

【数据来源】
- 后端提供音频流接口：GET /keytone_pkg/get_audio_stream?sha256=...&type=...
- 本组件通过 boot/axios 的 baseURL 构造请求 URL。

【注意】
- 本组件不负责保存/预览，只负责“选区选择 + 同步”。
- 若音频加载/解码失败，会显示错误并保持数字输入可用（降级）。
============================================================================
-->

<template>
  <div class="waveform-wrapper">
    <!--
      体验增强区：播放条 + 缩放
      - 这里的试听完全在前端完成，不依赖 SDK 播放（因此支持暂停与拖动播放头）
      - SDK 的“预览按钮”仍然保留，用于 <=5s 的链路验真
    -->
    <div v-if="hasSource" class="flex items-center justify-between gap-2 mb-2">
      <div class="flex items-center gap-2">
        <q-btn
          dense
          color="primary"
          :disable="!isReady"
          :icon="isPlaying ? 'pause' : 'play_arrow'"
          :label="isPlaying ? t('KeyToneAlbum.defineSounds.waveformTrimmer.pause') : t('KeyToneAlbum.defineSounds.waveformTrimmer.play')"
          @click="togglePlay"
        />

        <q-btn
          dense
          flat
          color="primary"
          icon="stop"
          :disable="!isReady"
          :label="t('KeyToneAlbum.defineSounds.waveformTrimmer.stop')"
          @click="stop"
        />

        <q-btn-toggle
          dense
          unelevated
          toggle-color="primary"
          color="grey-3"
          text-color="grey-8"
          v-model="playScope"
          :options="playScopeOptions"
          :disable="!isReady"
        />
      </div>

      <div class="flex items-center gap-2 min-w-[260px]">
        <div class="text-[12px] text-gray-600 whitespace-nowrap">{{ t('KeyToneAlbum.defineSounds.waveformTrimmer.zoom') }}</div>
        <q-slider
          v-model="zoomMinPxPerSec"
          :min="zoomMin"
          :max="zoomMax"
          :step="10"
          dense
          :disable="!isReady"
          class="flex-1"
        />
        <div class="text-[12px] text-gray-600 whitespace-nowrap">{{ zoomMinPxPerSec }}</div>
      </div>
    </div>

    <div v-if="!hasSource" class="p-2">
      <div class="text-[12px] text-gray-500">{{ t('KeyToneAlbum.defineSounds.waveformTrimmer.noSource') }}</div>
    </div>

    <div v-else-if="!isReady && !hasError" class="p-2">
      <q-skeleton type="rect" height="92px" />
      <div class="text-[12px] text-gray-500 mt-2">{{ t('KeyToneAlbum.defineSounds.waveformTrimmer.loading') }}</div>
    </div>

    <div v-if="hasError" class="p-2">
      <q-banner dense rounded class="bg-orange-50 text-orange-900">
        {{ t('KeyToneAlbum.defineSounds.waveformTrimmer.unavailableFallback') }}
      </q-banner>
    </div>

    <div
      v-show="hasSource"
      class="waveform-shell"
      :style="{ '--volumeScaleW': `${volumeScaleWidthPx}px`, '--volumeScaleOffset': `${volumeScaleOffsetPx}px` }"
    >
      <!--
        重要：用户明确要求“新增刻度不能压缩波形区域宽度”。
        因此这里采用 overlay 方案：
        - 左右刻度/当前 dB 作为覆盖层贴在波形左右边缘（不占用布局宽度）。
        - 波形 DOM 宽度保持与改造前一致（100% 宽度，不额外减去侧栏）。
        - 对话框宽度不足时，可通过对话框宽度上限（项目固定宽度）适度扩展；
          但实现本身不会通过“分栏”挤压波形。
      -->

      <!-- 波形主区域（保持原有交互/渲染；宽度不被刻度挤压） -->
      <div class="waveform-host">
        <div
          ref="containerEl"
          class="waveform"
          @pointerdown.capture="onWaveformPointerDown"
          @mousedown.capture="onWaveformMouseDown"
          @contextmenu.prevent
        />

      <!--
        音量快捷调节条（剪辑软件风格）：
        - 一条横向指示线代表当前音量（0 位于中线）。
        - 在指示线附近按下并拖动：
          * 向上拖动 => 音量增加
          * 向下拖动 => 音量降低
        - 该控件仅影响当前音频裁剪的 cut.volume（与输入框双向同步）。
      -->
      <!-- 音量线 overlay（覆盖在波形上，不阻挡交互） -->

      <!--
        交互命中区必须放在 pointer-events 可用的层级中。

        说明：
        - 我们希望“音量线/手柄”不阻挡波形本体的交互（seek、右键选区等），因此视觉层 .volume-overlay-shell
          设置为 pointer-events:none。
        - 但在部分浏览器/Electron 组合下，父节点 pointer-events:none 会导致后代元素无法成为命中目标，
          即使子元素显式 pointer-events:auto 也可能无效。
        - 因此把命中区单独作为 overlay 的兄弟节点：既不会挡住波形整体交互，又能稳定接收拖动事件。
      -->
      <div
        class="volume-line-hit"
        :style="{ top: volumeLineTopPx }"
        @pointerdown="onVolumePointerDown"
      />

      <!--
        超出范围提示（波形区域内）：
        - 用户要求“波形组件处也有提示”。
        - 放在右下角，避免遮挡刻度或主交互区域。
        - 仅作轻量提示，不影响交互。
      -->
      <div v-if="isDbOutOfRange" class="volume-overflow-inline">
        {{ t('KeyToneAlbum.defineSounds.waveformTrimmer.volumeOutOfRange') }}
      </div>

        <!--
          横向音量线 overlay（仅覆盖波形区域）：
          - 不再追求“贯穿整个对话框”，避免引入额外布局复杂度。
          - 用户的核心诉求是“刻度贴近波形 + 不压缩波形宽度”。
        -->
        <div class="volume-overlay-shell" aria-hidden="true" :style="{ height: `${height}px` }">
          <div class="volume-zero-line" />
          <div class="volume-line" :style="{ top: volumeLineTop }" />
        </div>
      </div>

      <!-- 左侧 dB 标尺（含负刻度；无单位，仅数字） -->
      <div class="volume-scale-overlay volume-scale-left" :style="{ height: `${height}px` }">
        <div
          v-for="mark in volumeDbScaleMarks"
          :key="`db-${mark}`"
          :class="['volume-scale-mark', { 'is-zero': mark === 0 }]"
          :style="{ top: dbToTopPercent(mark) }"
        >
          <span class="volume-scale-text">{{ formatDbMark(mark) }}</span>
        </div>
      </div>

      <!-- 右侧当前音量显示（对齐蓝线，显示 dB 单位） -->
      <div class="volume-scale-overlay volume-scale-right" :style="{ height: `${height}px` }">
        <div class="volume-current-db" :style="{ top: volumeLineTop }">{{ currentDbLabel }}</div>
      </div>
    </div>

    <div v-if="isReady && durationMs" class="text-[12px] text-gray-500 mt-1 flex flex-wrap gap-x-3 gap-y-1">
      <div>{{ t('KeyToneAlbum.defineSounds.waveformTrimmer.current') }}{{ formatMs(currentTimeMs) }}</div>
      <div>{{ t('KeyToneAlbum.defineSounds.waveformTrimmer.total') }}{{ formatMs(Math.round(durationMs)) }}</div>
      <div v-if="selectionDurationMs !== null">
        {{ t('KeyToneAlbum.defineSounds.waveformTrimmer.selection') }}{{ formatMs(selectionDurationMs) }}
      </div>
      <div class="text-gray-400">{{ zoomHintText }}</div>
      <!--
        超范围提示移到信息行（避免与标尺/当前值重叠）。
        语义：普通提示（非错误/警告），仅告知“显示被贴边”，实际数值仍有效。
      -->
      <div v-if="isDbOutOfRange" class="text-gray-400">
        {{ t('KeyToneAlbum.defineSounds.waveformTrimmer.volumeOutOfRange') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import WaveSurfer from 'wavesurfer.js';
import { Platform } from 'quasar';
// wavesurfer.js v7 插件（ESM）
import RegionsPlugin from 'wavesurfer.js/dist/plugins/regions.esm.js';
import { api } from 'boot/axios';

type Region = {
  id: string;
  start: number;
  end: number;
  remove?: () => void;
  /**
   * wavesurfer.js v7 Regions：更新 region 的标准方式。
   *
   * 重要：这里不能用 update()/set() 之类的“猜测 API”。
   * 之前输入框修改无法反映到选区上，根因就是使用了不存在的 update()，导致调用被可选链吞掉。
   */
  setOptions?: (options: {
    start?: number;
    end?: number;
    id?: string;
    drag?: boolean;
    resize?: boolean;
    color?: string;
  }) => void;
};

type RegionsApi = {
  enableDragSelection?: (payload: { color: string }) => void;
  on?: (event: string, cb: (region: Region) => void) => void;
  getRegions?: () => Region[];
  addRegion?: (payload: {
    id: string;
    start: number;
    end: number;
    color: string;
    drag: boolean;
    resize: boolean;
  }) => Region;
};

const props = withDefaults(
  defineProps<{
    sha256: string;
    fileType: string;
    startMs: number;
    endMs: number;
    /**
     * 与 SDK 一致的 cut.volume。
     * SDK 侧使用 beep/effects.Volume (Base=1.6)，因此该值不是 0~1 的线性音量。
     * 前端试听会尝试做等价映射以尽量贴近听感。
     */
    volume?: number;
    height?: number;
  }>(),
  {
    volume: 0,
    height: 80,
  }
);

const emit = defineEmits<{
  (e: 'update:startMs', value: number): void;
  (e: 'update:endMs', value: number): void;
  (e: 'update:volume', value: number): void;
  (e: 'loaded', payload: { durationMs: number }): void;
}>();

const containerEl = ref<HTMLElement | null>(null);
const ws = ref<WaveSurfer | null>(null);
const regionsPlugin = ref<RegionsApi | null>(null);
const regionId = 'trim-region';

const isReady = ref(false);
const hasError = ref(false);
const durationMs = ref<number | null>(null);

// 播放条状态
const isPlaying = ref(false);
const currentTimeMs = ref(0);
const playScope = ref<'all' | 'selection'>('all');
const playScopeOptions = computed(() => [
  { label: t('KeyToneAlbum.defineSounds.waveformTrimmer.playScopeAll'), value: 'all' },
  { label: t('KeyToneAlbum.defineSounds.waveformTrimmer.playScopeSelection'), value: 'selection' },
]);

// 缩放（zoom）
// minPxPerSec 越大，波形越“拉长”，越容易观察 100ms 级别细节。
const zoomMin = 50;
const zoomMax = 1000;
// 默认缩放：50
// - 之前的默认值 200 在多数音频上“过度放大”，初始打开时会显得很拥挤。
// - 用户期望默认更“放得开”，先从整体观察，再按需放大。
const zoomMinPxPerSec = ref(50);

const selectionDurationMs = computed(() => {
  if (!(props.endMs > props.startMs) || props.startMs < 0 || props.endMs < 0) return null;
  return Math.round(props.endMs - props.startMs);
});

// i18n
// - 本项目同时支持多语言（至少中英文）。
// - 这个组件新增了较多 UI 文案（播放/缩放/提示/错误提示等），必须全部走 i18n，避免后续漏翻译与难维护。
const { t } = useI18n();

// 缩放提示文案：根据平台展示合适的快捷键说明
// - macOS：Control + 滚轮 / 触控板捏合
// - 其他：Ctrl + 滚轮
const zoomHintText = computed(() => (Platform.is.mac ? t('KeyToneAlbum.defineSounds.waveformTrimmer.hintMac') : t('KeyToneAlbum.defineSounds.waveformTrimmer.hint')));

// macOS 外接键盘修饰键状态（兜底）：
// - 某些外接机械键盘的 Control 键在 wheel 事件里不会稳定反映到 ctrlKey。
// - 因此在 window 层记录 keydown/keyup 作为备用判断。
const macModifierState = {
  control: false,
  meta: false,
};

function onMacModifierKeyDown(ev: KeyboardEvent) {
  if (!Platform.is.mac) return;
  if (ev.key === 'Control') macModifierState.control = true;
  if (ev.key === 'Meta') macModifierState.meta = true;
}

function onMacModifierKeyUp(ev: KeyboardEvent) {
  if (!Platform.is.mac) return;
  if (ev.key === 'Control') macModifierState.control = false;
  if (ev.key === 'Meta') macModifierState.meta = false;
}

function resetMacModifierState() {
  macModifierState.control = false;
  macModifierState.meta = false;
}

const audioUrl = computed(() => {
  if (!props.sha256 || !props.fileType) return '';
  const baseURL = api.defaults.baseURL || '';
  const sha256 = encodeURIComponent(props.sha256);
  const type = encodeURIComponent(props.fileType);
  return `${baseURL}/keytone_pkg/get_audio_stream?sha256=${sha256}&type=${type}`;
});

const hasSource = computed(() => !!audioUrl.value);

const isSyncingFromProps = ref(false);

// ============================================================================
// 音量快捷调节（水平指示线 + 垂直拖拽）
//
// 需求：模拟剪辑软件的“水平音量指针条”。
// - 0 音量对应中线；向上拖动增加音量，向下拖动降低音量。
// - 内部仍使用 cut.volume（Base=1.6 的指数音量）；UI 显示为 dB。
// - 以 SDK 预览为准：前端预览采用同一映射，不再额外 clamp。
//
// 设计：
// - 为避免无穷范围带来的 UI 不可控，这里定义“可视范围”用于拖拽/指示线位置映射。
// - 该范围不限制用户输入的数值：手动输入超出范围时，指示线会贴边显示；
//   用户拖拽时则会把值限定在可视范围内（这是交互层的合理约束）。
// ============================================================================
// “可视范围”的定义必须保证 0dB 永远在中位。
//
// 用户诉求：UI 上 volume=0 一定处于中位，但不影响实际的 cut.volume 语义。
// 因此这里采用“对称区间”来定义拖拽映射：[-X, +X]。
// - 这只是 UI 映射范围：用于拖拽与指示线定位。
// - 不限制用户通过输入框直接输入更大/更小的值（超出范围时指示线会贴边显示）。
// dB 可视范围（用户确认：±18 dB）
const volumeDbVisibleAbsMax = 18;
const volumeDbMin = -volumeDbVisibleAbsMax;
const volumeDbMax = volumeDbVisibleAbsMax;

// dB <-> cut.volume（Base=1.6）换算：
// gain = 1.6 ^ volume
// dB = 20 * log10(gain) = 20 * volume * log10(1.6)
const dbPerVolume = 20 * Math.log10(1.6);
const volumeToDb = (volume: number) => volume * dbPerVolume;
const dbToVolume = (db: number) => db / dbPerVolume;

const isVolumeDragging = ref(false);
const volumeDragPointerId = ref<number | null>(null);
const volumeDragTargetEl = ref<HTMLElement | null>(null);

const currentDb = computed(() => volumeToDb(props.volume ?? 0));
const isDbOutOfRange = computed(() => currentDb.value < volumeDbMin || currentDb.value > volumeDbMax);

const currentDbLabel = computed(() => {
  const val = currentDb.value;
  const sign = val > 0 ? '+' : '';
  return `${sign}${val.toFixed(1)} dB`;
});

const dbToTopPercent = (db: number) => {
  const clamped = Math.max(volumeDbMin, Math.min(volumeDbMax, db));
  const ratio = (clamped - volumeDbMin) / (volumeDbMax - volumeDbMin); // 0..1
  return `${(1 - ratio) * 100}%`;
};

const volumeDbScaleMarks = [18, 12, 6, 0, -6, -12, -18];

// 标尺与波形的“外侧间距”控制：
// - 这里设为 0，使刻度紧贴波形边界外侧（符合你的要求）。
// - 负刻度文本通过固定宽度容纳。
const volumeScaleWidthPx = 24;
const volumeScaleOffsetPx = 0;

// 刻度显示文本：正分贝带 + 号，0 保持 0。
const formatDbMark = (mark: number) => (mark > 0 ? `+${mark}` : `${mark}`);

const volumeLineTop = computed(() => {
  const clampedDb = Math.max(volumeDbMin, Math.min(volumeDbMax, currentDb.value));
  const ratio = (clampedDb - volumeDbMin) / (volumeDbMax - volumeDbMin); // 0..1
  // ratio=1 => top=0（最大音量在最上）; ratio=0 => top=100%（最小音量在最下）
  // 注意：这里不要做 Math.round。
  // - 用户明确要求：volume=0 必须与波形中位重合。
  // - 若做整数百分比四舍五入，会在某些高度下造成 1px 级别的偏移，让用户产生“没对齐”的感知。
  // - 保留小数百分比可以保证 0 映射为精确 50%。
  return `${(1 - ratio) * 100}%`;
});

const volumeLineTopPx = computed(() => {
  // 计算音量指示线的纵向位置（像素值，用于命中区的精确定位）。
  //
  // 为何需要像素值：
  // - 命中区（volume-line-hit）需要基于波形绘制区域高度（props.height）进行定位，
  //   而非整个 .waveform 容器高度（可能包含滚动条）。
  // - 使用百分比时，若命中区位于与 overlay 不同高度的容器中，会导致位置偏移。
  // - 因此这里直接计算像素值：top = (1 - ratio) * props.height。
  const clampedDb = Math.max(volumeDbMin, Math.min(volumeDbMax, currentDb.value));
  const ratio = (clampedDb - volumeDbMin) / (volumeDbMax - volumeDbMin);
  return `${(1 - ratio) * props.height}px`;
});

const waveformYScale = computed(() => {
  // 体验增强：音量调整时，波形“看起来”也应该发生变化（剪辑软件直觉）。
  //
  // 关键澄清（库能力边界）：
  // - wavesurfer.js 的 waveform 是根据“音频采样数据”渲染出来的。
  // - playback volume（setVolume / cut.volume）只影响播放增益，不会自动影响波形渲染。
  // - wavesurfer.js 提供了用于渲染阶段的“纵向缩放参数”（barHeight），可改变渲染出来的波形高度。
  //
  // 为什么不用 canvas 的 transform(scaleY)：
  // - transform 会对像素进行缩放插值，容易出现你反馈的“静音中位线变粗/发糊/加醋”的伪影。
  // - 特别是在缩放到最大值瞬间，插值会更明显。
  //
  // 因此我们把该值用于 wavesurfer 的 barHeight（渲染阶段缩放），并触发重绘：
  // - 渲染出来的线条更干净
  // - 不会产生 transform 插值导致的粗线伪影
  //
  // 映射策略：与 SDK 听感曲线同源，但更温和，并限制最大最小，避免极端时过于“顶满”。
  // - SDK 听感：gain = 1.6 ^ volume
  // - 视觉缩放：scale = 1.6 ^ (volume * 0.08)
  //   => volume=0  => 1
  //   => volume=10 => ~1.48（明显但不会把整体顶满得太夸张）
  //   => volume=-10=> ~0.68
  const volume = props.volume ?? 0;
  const scale = Math.pow(1.6, volume * 0.08);
  return Math.max(0.65, Math.min(1.5, scale));
});

// =============================================================================
// 波形“音量反馈”实现（更美观的方式）
//
// 目标：音量变化时波形高度随之变化，但不产生 transform 插值伪影。
// 实现：使用 wavesurfer.js 的渲染参数 barHeight，并在音量变化时 setOptions 触发重绘。
//
// 性能注意：
// - 用户拖动音量线会产生高频更新（pointermove）。
// - 直接在每次事件里 setOptions 会导致过多重绘。
// - 因此这里用 requestAnimationFrame 做“每帧最多一次”的合并更新。
// =============================================================================
let barHeightRafId: number | null = null;
let pendingBarHeight: number | null = null;

function scheduleApplyWaveformBarHeight() {
  const instance = ws.value;
  if (!instance) return;

  // 记录本次希望应用的值（会被后续更新覆盖，确保最后一次赢）
  pendingBarHeight = waveformYScale.value;

  if (barHeightRafId !== null) return;
  barHeightRafId = requestAnimationFrame(() => {
    barHeightRafId = null;
    const next = pendingBarHeight;
    pendingBarHeight = null;
    if (next === null) return;

    // 核心：通过渲染参数改变波形高度（而不是对 canvas 做 transform）。
    // setOptions 会触发内部 reRender，从而更新波形绘制结果。
    instance.setOptions({ barHeight: next });
  });
}

// 视觉：选区使用半透明填充。
// 注意：本组件只保留“可双端拖动”的选区版本。
// 因此即使 start==end（例如右键按下瞬间或用户只修改了一个输入框导致暂时 end<=start），
// 也会被渲染为一个“极小的非零区间”，以确保左右两侧的拖动指针始终可用。
const regionFillColor = 'rgba(14, 165, 233, 0.18)';

// 选区的最小非零长度（秒）：用于在 start==end 时仍能展示两侧可拖动指针。
// - 这是渲染层的“非零化”，不等价于产品层的“最小裁剪长度”限制。
// - 数值取 1ms：足够让 Regions 进入“region”形态，而不会引入可感知的最小长度。
const minNonZeroSelectionSec = 0.001;

// ============================================================================
// 右键拖拽快速选区（无菜单）
//
// 用户需求（严格按提案实现）：
// 1) 右键按下的瞬间即设置起点(start)
// 2) 右键按下后向右拖动，实时更新终点(end)
// 3) 松开右键的瞬间，将当前位置作为终点(end)
// 4) 全程不弹出任何菜单
//
// 备注：
// - 这里不引入“最小长度”规则（用户明确不需要）。
// - 由于浏览器默认右键会弹出 context menu，因此必须在波形区域阻止默认行为。
// - 为了让拖拽过程即使离开波形区域也能持续更新，我们使用 window 级别的 pointermove/pointerup。
// ============================================================================
const isRightDragSelecting = ref(false);
const rightDragStartSec = ref<number | null>(null);

// 左键拖拽类型：
// - playhead：播放头拖拽
// - region-move：拖拽整个选区
// - region-resize-start / region-resize-end：拖拽选区左右指针
//
// 说明：
// - 我们在这里显式区分拖拽模式，以避免 wavesurfer 内部拖拽状态机与自动滚动“打架”。
// - 由我们统一更新播放头/选区位置，确保“跟手、丝滑”。
type LeftDragMode = 'none' | 'playhead' | 'region-move' | 'region-resize-start' | 'region-resize-end';
const leftDragMode = ref<LeftDragMode>('none');

// 左键拖拽的基准信息（用于计算 delta 或固定另一端）：
const leftDragStartSec = ref<number | null>(null);
let leftDragRegionStartSec: number | null = null;
let leftDragRegionEndSec: number | null = null;

// pointer capture 的释放需要拿到同一个 element + pointerId。
// 注意：我们的 pointerup 监听是挂在 window 上的，因此 pointerup 的 currentTarget 并不是当初按下的元素。
// 所以这里显式记录按下时的 element 与 pointerId，用于在结束时成对释放。
const rightDragCapturedEl = ref<HTMLElement | null>(null);
const rightDragCapturedPointerId = ref<number | null>(null);

function clampTimeSec(instance: WaveSurfer, t: number): number {
  const dur = instance.getDuration();
  if (!Number.isFinite(dur) || dur <= 0) return 0;
  return Math.max(0, Math.min(dur, t));
}

function isWaveSurferInteractive(instance: WaveSurfer): boolean {
  const dur = instance.getDuration();
  return Number.isFinite(dur) && dur > 0;
}

function resolveWaveformScrollElement(instance: WaveSurfer): HTMLElement {
  // wavesurfer 的 DOM 结构在不同版本/配置下会变化：
  // - 有时外层 container 本身可滚动
  // - 有时真正滚动的是内部 wrapper 或 wrapper 的 parent
  //
  // 如果选错元素，会出现：
  // - 分母使用了“可视区域宽度”（clientWidth）而不是“总波形宽度”（scrollWidth）
  // - 从而把 x/width 的比例放大，导致 time 被算得非常靠后、非常大（用户反馈：“点一下时间就不对”）
  const wrapper = instance.getWrapper() as HTMLElement | null;
  const innerWrapper = (containerEl.value?.querySelector('.wavesurfer-wrapper') as HTMLElement | null) || null;

  const candidates: Array<HTMLElement | null | undefined> = [
    // 最优：wavesurfer 内部常见的 wrapper
    innerWrapper,
    // 次优：外层容器（我们绑定事件的地方）
    containerEl.value,
    // 再次：wrapper 的 parent 通常就是滚动容器
    wrapper?.parentElement as HTMLElement | null,
    // 最后兜底：wrapper 自身
    wrapper,
  ];

  for (const el of candidates) {
    if (!el) continue;
    // scrollWidth 明显大于 clientWidth 时，说明它是“总波形宽度”的滚动视窗
    if (el.scrollWidth > el.clientWidth + 1) return el;
  }

  // 如果没有任何一个表现出可滚动，则退回到 container（至少坐标系一致）
  return containerEl.value || wrapper || (document.body as HTMLElement);
}

function pointerEventToTimeSec(instance: WaveSurfer, ev: PointerEvent): number {
  // 关键：在有横向滚动/zoom 的情况下，时间映射必须加上 scrollLeft 偏移。
  //
  // 这里不使用 wavesurfer 的 getWidth/getScroll：
  // - 在 wavesurfer v7 不同渲染模式/后端下，getWidth/getScroll 的语义并不总是“像素”。
  // - 一旦把非像素单位当像素使用，会出现严重偏差（例如点击靠前位置却被识别成很靠后时间）。
  //
  // 因此我们完全基于 DOM 计算：
  // - 使用“真实滚动视窗”元素的 rect + scrollLeft + scrollWidth
  // - time = (xInView + scrollLeft) / scrollWidth * duration
  // 这样不依赖 wavesurfer 的内部实现细节，且与用户视觉位置一致。
  const scrollEl = resolveWaveformScrollElement(instance);
  const rect = scrollEl.getBoundingClientRect();
  const xInView = Math.max(0, Math.min(rect.width, ev.clientX - rect.left));

  const dur = instance.getDuration();
  if (!Number.isFinite(dur) || dur <= 0) return 0;

  const scrollLeft = scrollEl.scrollLeft || 0;
  const totalWidth = scrollEl.scrollWidth || scrollEl.clientWidth || 0;
  if (!Number.isFinite(totalWidth) || totalWidth <= 0) return 0;

  const xGlobal = xInView + scrollLeft;
  return clampTimeSec(instance, (xGlobal / totalWidth) * dur);
}

/**
 * 使用指定 scrollLeft 计算 clientX 对应时间点（用于自动滚动帧内预测）。
 */
function clientXToTimeSecWithScrollLeft(
  instance: WaveSurfer,
  clientX: number,
  scrollLeftOverride?: number
): number {
  const scrollEl = resolveWaveformScrollElement(instance);
  const rect = scrollEl.getBoundingClientRect();
  const xInView = Math.max(0, Math.min(rect.width, clientX - rect.left));

  const dur = instance.getDuration();
  if (!Number.isFinite(dur) || dur <= 0) return 0;

  const scrollLeft = scrollLeftOverride ?? scrollEl.scrollLeft ?? 0;
  const totalWidth = scrollEl.scrollWidth || scrollEl.clientWidth || 0;
  if (!Number.isFinite(totalWidth) || totalWidth <= 0) return 0;

  const xGlobal = xInView + scrollLeft;
  return clampTimeSec(instance, (xGlobal / totalWidth) * dur);
}

/**
 * 将给定时间范围滚动到可见区域内（尽量少移动视图）。
 *
 * 为什么需要：
 * - 在 zoom 很大（minPxPerSec 很高）时，波形的总宽度（scrollWidth）会非常大。
 * - 用户通过“右键快捷选区”创建/更新选区后，如果 Regions 的 DOM 刷新/重绘延迟，
 *   或者选区恰好贴近视窗边缘，用户可能会主观感知为“选区没有显示”。
 * - 手动轻微调整 zoom 会触发 wavesurfer reRender，选区又立刻出现（用户反馈的现象）。
 *
 * 本函数负责其中一半问题：
 * - 让选区范围落在当前滚动视窗的可见范围内（带少量 padding），避免“选区其实在视野外”。
 */
function scrollTimeRangeIntoView(instance: WaveSurfer, startSec: number, endSec: number) {
  const scrollEl = resolveWaveformScrollElement(instance);
  const dur = instance.getDuration();
  if (!Number.isFinite(dur) || dur <= 0) return;

  const totalWidth = scrollEl.scrollWidth || scrollEl.clientWidth || 0;
  const viewWidth = scrollEl.clientWidth || 0;
  if (!Number.isFinite(totalWidth) || totalWidth <= 0) return;
  if (!Number.isFinite(viewWidth) || viewWidth <= 0) return;

  // 统一排序：避免出现 start>end 的输入。
  const s = Math.min(startSec, endSec);
  const e = Math.max(startSec, endSec);

  // 映射到“全局像素坐标”（以 scrollWidth 为基准）。
  const startX = (s / dur) * totalWidth;
  const endX = (e / dur) * totalWidth;

  // 可视区留白：让选区不要贴边，观感更像剪辑软件。
  const padding = 24;

  const currentLeft = scrollEl.scrollLeft || 0;
  const currentRight = currentLeft + viewWidth;

  // 如果整个选区已经完全可见，则不动。
  if (startX >= currentLeft + padding && endX <= currentRight - padding) return;

  // 若选区比视窗还宽：优先让 start 出现在可见区左侧附近。
  if (endX - startX > viewWidth - padding * 2) {
    const nextLeft = Math.max(0, Math.min(totalWidth - viewWidth, startX - padding));
    scrollEl.scrollLeft = nextLeft;
    return;
  }

  // 选区较短：尽量让整个选区可见。
  if (startX < currentLeft + padding) {
    const nextLeft = Math.max(0, Math.min(totalWidth - viewWidth, startX - padding));
    scrollEl.scrollLeft = nextLeft;
    return;
  }

  if (endX > currentRight - padding) {
    const nextLeft = Math.max(0, Math.min(totalWidth - viewWidth, endX - (viewWidth - padding)));
    scrollEl.scrollLeft = nextLeft;
  }
}

/**
 * 判断某个时间范围是否已在当前滚动视窗内“基本可见”。
 *
 * 之所以单独抽出来：
 * - 右键拖拽选区时，我们希望做到“全过程可见”，并在必要时自动滚动。
 * - 但如果每一帧都强行 scroll，会带来视图抖动/跳动。
 * - 因此先判断是否真的需要滚动：只有选区跑出视野时才滚。
 */
function isTimeRangeVisible(instance: WaveSurfer, startSec: number, endSec: number, padding = 24): boolean {
  const scrollEl = resolveWaveformScrollElement(instance);
  const dur = instance.getDuration();
  if (!Number.isFinite(dur) || dur <= 0) return true;

  const totalWidth = scrollEl.scrollWidth || scrollEl.clientWidth || 0;
  const viewWidth = scrollEl.clientWidth || 0;
  if (!Number.isFinite(totalWidth) || totalWidth <= 0) return true;
  if (!Number.isFinite(viewWidth) || viewWidth <= 0) return true;

  const s = Math.min(startSec, endSec);
  const e = Math.max(startSec, endSec);
  const startX = (s / dur) * totalWidth;
  const endX = (e / dur) * totalWidth;

  const left = scrollEl.scrollLeft || 0;
  const right = left + viewWidth;

  return startX >= left + padding && endX <= right - padding;
}

// =============================================================================
// 右键快捷选区：高 zoom 下“全过程不可见”的修复
//
// 现象（用户反馈）：
// - zoom 很大时，右键按下/拖动过程中选区经常“不显示”。
// - 松开右键时选区又出现，或轻微调整 zoom 后立即出现。
//
// 根因推断：
// - wavesurfer/regions 的 region DOM 更新在某些渲染时序下会延迟，直到下一次 reRender。
// - zoom 变化会触发 reRender，所以用户“微调 zoom”能把选区“唤醒”。
//
// 修复目标：
// - 从右键按下的瞬间开始，到拖拽结束，全过程都应当看到选区。
// - 同时避免在 pointermove 的高频触发下造成过多重绘。
//
// 实现策略：
// - 使用 requestAnimationFrame 合并刷新：每帧最多触发一次“no-op zoom”来促使 reRender。
// - 在必要时（选区跑出视野）自动滚动，让选区保持可见。
// =============================================================================
let quickSelectLiveRafId: number | null = null;
let pendingQuickSelectRange: { startSec: number; endSec: number } | null = null;

// 右键拖拽时的“边缘自动滚动”状态：
// - 体验目标：像剪辑软件一样，鼠标靠近左右边缘时才自动滚动，且速度随靠近程度渐进。
// - 重要：自动滚动过程中，即使鼠标不再移动（不产生 pointermove），选区 end 也应持续更新。
//   因为 scrollLeft 变化会改变“同一 clientX 对应的时间点”。
const quickSelectAutoScrollEdgePx = 48;
const quickSelectAutoScrollMaxSpeedPxPerSec = 1400;
let quickSelectAutoScrollRafId: number | null = null;
let quickSelectAutoScrollLastTs: number | null = null;
let quickSelectLastClientX: number | null = null;

// 左键拖拽相关（播放头拖拽 / 选区两侧指针拖拽）：
// - 复用“边缘触发”自动滚动手感
// - 与右键快捷选区逻辑分离，避免互相干扰
const leftDragAutoScrollEdgePx = quickSelectAutoScrollEdgePx;
const leftDragAutoScrollMaxSpeedPxPerSec = quickSelectAutoScrollMaxSpeedPxPerSec;
const leftDragMoveThresholdPx = 3; // 防止点击误触发拖拽滚动

let leftDragActive = false;
let leftDragMoved = false;
let leftDragStartClientX: number | null = null;
let leftDragLastClientX: number | null = null;
let leftDragAutoScrollRafId: number | null = null;
let leftDragAutoScrollLastTs: number | null = null;
let leftDragCapturedEl: HTMLElement | null = null;
let leftDragCapturedPointerId: number | null = null;
// 左键自动滚动速度（平滑用）：避免抖动/卡顿
let leftDragScrollVelocity = 0; // px/s

/**
 * 左键拖拽时“跟手更新”的统一入口。
 *
 * 设计目标：
 * - 播放头拖拽：光标在哪，播放头就在哪（始终跟随）。
 * - 选区指针拖拽：光标在哪，对应指针就在哪（start/end 其中一侧跟随）。
 * - 在自动滚动与非自动滚动场景下都保持一致。
 */
function updateLeftDragTarget(instance: WaveSurfer, clientX: number, targetSecOverride?: number) {
  const syntheticEv = { clientX, button: 0 } as unknown as PointerEvent;
  const targetSec = targetSecOverride ?? pointerEventToTimeSec(instance, syntheticEv);

  // 播放头拖拽：直接 setTime，保证跟手。
  if (leftDragMode.value === 'playhead') {
    instance.setTime(targetSec);
    return;
  }

  // 选区相关拖拽：直接更新 region。
  const regionsApi = regionsPlugin.value;
  if (!regionsApi) return;
  const region = getTrimRegion(regionsApi);
  if (!region) return;

  const baseStart = leftDragRegionStartSec ?? region.start;
  const baseEnd = leftDragRegionEndSec ?? region.end;

  let nextStart = baseStart;
  let nextEnd = baseEnd;

  if (leftDragMode.value === 'region-resize-start') {
    nextStart = targetSec;
  } else if (leftDragMode.value === 'region-resize-end') {
    nextEnd = targetSec;
  } else if (leftDragMode.value === 'region-move') {
    // 拖拽整个选区：按光标相对位移整体平移。
    // 使用 leftDragStartSec 作为“拖拽起点”，计算 delta。
    if (leftDragStartSec.value !== null) {
      const delta = targetSec - leftDragStartSec.value;
      nextStart = baseStart + delta;
      nextEnd = baseEnd + delta;
    }
  }

  const normalized = normalizeSelectionSec(instance, nextStart, nextEnd);
  ensureTrimRegion(regionsApi, instance, normalized.startSec, normalized.endSec);

  emit('update:startMs', Math.max(0, Math.round(normalized.startSec * 1000)));
  emit('update:endMs', Math.max(0, Math.round(normalized.endSec * 1000)));

  // best-effort：仍派发 synthetic pointermove，兼容内部状态机
  try {
    const wrapper = instance.getWrapper() as HTMLElement | null;
    if (wrapper) {
      const wrapRect = wrapper.getBoundingClientRect();
      const synthetic = new PointerEvent('pointermove', {
        clientX,
        clientY: wrapRect.top + wrapRect.height / 2,
        button: 0,
        buttons: 1,
        bubbles: true,
      });
      wrapper.dispatchEvent(synthetic);
    }
  } catch {
    // ignore
  }
}

function scheduleQuickSelectLiveVisibility(instance: WaveSurfer, startSec: number, endSec: number) {
  // 记录最新区间（拖动过程中会被持续覆盖，最后一次赢）
  pendingQuickSelectRange = { startSec, endSec };

  if (quickSelectLiveRafId !== null) return;
  quickSelectLiveRafId = requestAnimationFrame(() => {
    quickSelectLiveRafId = null;

    const next = pendingQuickSelectRange;
    pendingQuickSelectRange = null;
    if (!next) return;

    // 1) 强制轻量重绘：no-op zoom
    // 说明：zoom 值不变，但会走 renderer.reRender；用于解决 region DOM 更新延迟。
    try {
      instance.zoom(zoomMinPxPerSec.value);
    } catch {
      // ignore
    }

    // 2) 这里不做“自动滚动到可见”。
    // 原因：
    // - 拖动过程中的自动滚动，应该采用“边缘触发”模型（靠近边缘才滚），否则会出现视图莫名跳动。
    // - finalize（松开右键）后我们仍会调用 ensureQuickSelectRegionVisible 做兜底滚动。
  });
}

/**
 * 右键拖拽的边缘自动滚动（剪辑软件式）。
 *
 * 行为定义：
 * - 当鼠标指针靠近波形滚动视窗的左/右边缘时，自动滚动。
 * - 指针越靠近边缘（甚至越过边缘），滚动速度越快（渐进）。
 * - 指针回到中间区域，立即停止自动滚动。
 *
 * 为什么要做一个独立的 RAF 循环：
 * - pointermove 事件不一定会持续触发（例如用户把鼠标停在边缘不动）。
 * - 但在自动滚动时，scrollLeft 每帧变化，end 对应的时间点也应该每帧更新。
 */
function scheduleQuickSelectEdgeAutoScroll(instance: WaveSurfer) {
  if (quickSelectAutoScrollRafId !== null) return;

  const step = (ts: number) => {
    quickSelectAutoScrollRafId = null;

    // 只有在右键拖拽进行中才允许自动滚动。
    if (!isRightDragSelecting.value) {
      quickSelectAutoScrollLastTs = null;
      return;
    }

    // 必须有 start 与 clientX，才能根据滚动更新 end。
    if (rightDragStartSec.value === null || quickSelectLastClientX === null) {
      quickSelectAutoScrollLastTs = null;
      return;
    }

    const scrollEl = resolveWaveformScrollElement(instance);
    const rect = scrollEl.getBoundingClientRect();

    // 计算边缘“渗透程度”并映射为速度：
    // - 进入边缘区：速度从 0 平滑增长
    // - 越过边缘：速度饱和到 max
    const leftEdge = rect.left + quickSelectAutoScrollEdgePx;
    const rightEdge = rect.right - quickSelectAutoScrollEdgePx;
    const x = quickSelectLastClientX;

    let direction = 0; // -1 向左滚，+1 向右滚
    let intensity = 0; // 0..1
    if (x < leftEdge) {
      direction = -1;
      const dist = Math.min(quickSelectAutoScrollEdgePx, leftEdge - x);
      intensity = dist / quickSelectAutoScrollEdgePx;
    } else if (x > rightEdge) {
      direction = 1;
      const dist = Math.min(quickSelectAutoScrollEdgePx, x - rightEdge);
      intensity = dist / quickSelectAutoScrollEdgePx;
    }

    // 不在边缘区：不滚动也不继续循环。
    if (direction === 0 || intensity <= 0) {
      quickSelectAutoScrollLastTs = null;
      return;
    }

    // dt：使用高精度时间戳，避免不同帧率下速度不一致。
    const lastTs = quickSelectAutoScrollLastTs ?? ts;
    quickSelectAutoScrollLastTs = ts;
    const dtSec = Math.max(0, Math.min(0.05, (ts - lastTs) / 1000));

    // 使用二次曲线让速度更“柔和”：靠近边缘一点点不会突然很快。
    const eased = intensity * intensity;
    const speed = eased * quickSelectAutoScrollMaxSpeedPxPerSec; // px/s
    const deltaPx = direction * speed * dtSec;

    const maxScrollLeft = Math.max(0, (scrollEl.scrollWidth || 0) - scrollEl.clientWidth);
    const before = scrollEl.scrollLeft;
    const next = Math.max(0, Math.min(maxScrollLeft, before + deltaPx));
    if (Math.abs(next - before) >= 0.5) {
      scrollEl.scrollLeft = next;

      // scrollLeft 改变后，需要用“同一 clientX”重新映射 end。
      // 这里复用 pointerEventToTimeSec：它会读取最新的 scrollLeft/scrollWidth。
      const syntheticEv = { clientX: x, button: 2 } as unknown as PointerEvent;
      const currentSec = pointerEventToTimeSec(instance, syntheticEv);
      const endSec = Math.max(rightDragStartSec.value, currentSec);
      const normalized = normalizeSelectionSec(instance, rightDragStartSec.value, endSec);

      const regionsApi = regionsPlugin.value;
      if (regionsApi) {
        ensureTrimRegion(regionsApi, instance, normalized.startSec, normalized.endSec);
        emit('update:startMs', Math.max(0, Math.round(normalized.startSec * 1000)));
        emit('update:endMs', Math.max(0, Math.round(normalized.endSec * 1000)));
      }

      // 同步做一次轻量重绘合并（解决高 zoom 下 region DOM 更新延迟）。
      scheduleQuickSelectLiveVisibility(instance, normalized.startSec, normalized.endSec);
    }

    // 继续下一帧：只要仍在边缘区且仍在拖拽。
    quickSelectAutoScrollRafId = requestAnimationFrame(step);
  };

  quickSelectAutoScrollRafId = requestAnimationFrame(step);
}

/**
 * 左键拖拽时的边缘自动滚动（用于播放头拖拽、选区拖拽/缩放）。
 *
 * 关键点：
 * - 只有在“确定是拖拽”（移动距离超过阈值）时才启用滚动，避免点击播放头时视图乱动。
 * - 自动滚动时，需要让 wavesurfer/regions 的拖拽逻辑“继续收到 pointermove”。
 *   否则会出现“滚动了，但播放头/选区没跟着动”的割裂感。
 * - 这里通过派发 synthetic pointermove 到 wrapper 尝试触发内部拖拽更新。
 *   （这是 best-effort，不依赖私有 API，失败也不会影响正常拖拽）。
 */
function scheduleLeftDragEdgeAutoScroll(instance: WaveSurfer) {
  if (leftDragAutoScrollRafId !== null) return;

  const step = (ts: number) => {
    leftDragAutoScrollRafId = null;

    if (!leftDragActive || !leftDragMoved) {
      leftDragAutoScrollLastTs = null;
      return;
    }

    if (leftDragLastClientX === null) {
      leftDragAutoScrollLastTs = null;
      return;
    }

    const scrollEl = resolveWaveformScrollElement(instance);
    const rect = scrollEl.getBoundingClientRect();

    const leftEdge = rect.left + leftDragAutoScrollEdgePx;
    const rightEdge = rect.right - leftDragAutoScrollEdgePx;
    const x = leftDragLastClientX;

    let direction = 0;
    let intensity = 0;
    if (x <= leftEdge) {
      direction = -1;
      const dist = Math.min(leftDragAutoScrollEdgePx, leftEdge - x);
      intensity = dist / leftDragAutoScrollEdgePx;
    } else if (x >= rightEdge) {
      direction = 1;
      const dist = Math.min(leftDragAutoScrollEdgePx, x - rightEdge);
      intensity = dist / leftDragAutoScrollEdgePx;
    }

    if (direction === 0 || intensity <= 0) {
      leftDragAutoScrollLastTs = null;
      return;
    }

    const lastTs = leftDragAutoScrollLastTs ?? ts;
    leftDragAutoScrollLastTs = ts;
    const dtSec = Math.max(0, Math.min(0.05, (ts - lastTs) / 1000));

    // 目标速度（渐进）：越接近边缘越快。
    const eased = intensity * intensity;
    const targetSpeed = eased * leftDragAutoScrollMaxSpeedPxPerSec; // px/s

    // 速度平滑：一阶低通，避免忽快忽慢导致“卡点感”。
    const alpha = 1 - Math.pow(0.15, dtSec * 60);
    leftDragScrollVelocity = leftDragScrollVelocity + (targetSpeed - leftDragScrollVelocity) * alpha;

    const deltaPx = direction * leftDragScrollVelocity * dtSec;

    const maxScrollLeft = Math.max(0, (scrollEl.scrollWidth || 0) - scrollEl.clientWidth);
    const before = scrollEl.scrollLeft;
    const next = Math.max(0, Math.min(maxScrollLeft, before + deltaPx));

    // 无论滚动是否发生，都必须“每帧”更新光标对应的时间点。
    // 这是“指针跟手”的关键：
    // - 在边缘渐进区，滚动速度可能非常小。
    // - 若仅在 scrollLeft 变化足够大时才更新，就会出现“卡在临界线”。
    const targetSec = clientXToTimeSecWithScrollLeft(instance, x, next);
    updateLeftDragTarget(instance, x, targetSec);

    // 始终写入 scrollLeft（即便很小），避免“速度很小 -> 不写入 -> 视觉卡住”。
    scrollEl.scrollLeft = next;

    leftDragAutoScrollRafId = requestAnimationFrame(step);
  };

  leftDragAutoScrollRafId = requestAnimationFrame(step);
}

/**
 * 修复：高 zoom 下“右键快捷选区”后选区偶发不可见。
 *
 * 用户现象：
 * - zoom 很大时，右键拖拽选区完成后，看不到新选区。
 * - 轻微调整 zoom（哪怕只改一点）后选区立刻出现。
 *
 * 推断根因：
 * - wavesurfer/regions 在某些渲染时序下，region DOM 的定位更新可能延迟到下一次 reRender。
 * - zoom 变化会触发 reRender，所以用户手动改 zoom 能“立刻修复”。
 *
 * 解决策略：
 * 1) 完成右键选区后，做一次“no-op zoom”触发 wavesurfer reRender（zoom 值不变）。
 * 2) 在重绘后，把选区滚动到可见（避免其实在视野外）。
 *
 * 性能说明：
 * - 该修复只在“右键选区完成（pointerup）”时执行一次，不会在拖动过程中高频触发。
 */
function ensureQuickSelectRegionVisible(instance: WaveSurfer, startSec: number, endSec: number) {
  const zoom = zoomMinPxPerSec.value;

  // 第 1 帧：触发 reRender（值不变，但会走 renderer.reRender）。
  requestAnimationFrame(() => {
    try {
      instance.zoom(zoom);
    } catch {
      // ignore
    }

    // 第 2 帧：在 reRender 基本完成后再滚动，避免 scrollWidth 尚未稳定。
    requestAnimationFrame(() => {
      scrollTimeRangeIntoView(instance, startSec, endSec);
    });
  });
}

/**
 * 判断是否为“右键按下”。
 *
 * 注意：提案/规范定义的是“右键拖拽快速选区”。
 * 为了保持行为一致且跨平台可预期，这里严格按 MouseEvent/PointerEvent 的标准语义判断：button === 2。
 */
function isSecondaryPointerDown(ev: PointerEvent): boolean {
  return ev.button === 2;
}

/**
 * 在 composedPath 中查找具有指定 part token 的元素。
 *
 * 背景：
 * - wavesurfer Regions 使用 `part` 标记区域与手柄（region / region-handle-left/right）。
 * - 事件目标可能不是 handle 本身（例如点击到 handle 内部或子节点）。
 * - 使用 composedPath 可以稳定找到真实命中元素。
 */
function findPartInComposedPath(ev: Event, token: string): Element | null {
  const path = (ev.composedPath?.() || []) as Array<EventTarget>;
  for (const node of path) {
    if (!(node instanceof Element)) continue;
    const part = node.getAttribute('part');
    if (!part) continue;
    const parts = part.split(/\s+/);
    if (parts.includes(token)) return node;
  }
  return null;
}

function getTrimRegion(regionsApi: RegionsApi): Region | undefined {
  const list = regionsApi.getRegions?.() || [];
  return list.find((r) => r.id === regionId);
}

function normalizeSelectionSec(instance: WaveSurfer, startSec: number, endSec: number): { startSec: number; endSec: number } {
  // 统一把选区规范化为一个“可操作”的非零区间：end MUST > start。
  // 这样 Regions 插件才会绘制为带两侧拖动手柄的 region，而不是 marker/线条。
  let start = clampTimeSec(instance, startSec);
  let end = clampTimeSec(instance, endSec);

  // 若 end 与 start 重合（或被 clamp 后变为重合），则做最小非零扩展。
  if (end <= start) {
    const dur = instance.getDuration();
    if (!Number.isFinite(dur) || dur <= 0) {
      end = start + minNonZeroSelectionSec;
    } else {
      const proposedEnd = Math.min(dur, start + minNonZeroSelectionSec);
      if (proposedEnd > start) {
        end = proposedEnd;
      } else {
        // 位于音频尾部且无法向右扩展：将 start 向左挪一点，保证 end>start。
        start = Math.max(0, dur - minNonZeroSelectionSec);
        end = dur;
      }
    }
  }

  return { startSec: start, endSec: end };
}

function ensureTrimRegion(regionsApi: RegionsApi, instance: WaveSurfer, startSec: number, endSec: number): Region | undefined {
  // 统一入口：没有就创建，有就更新
  const existing = getTrimRegion(regionsApi);

  const normalized = normalizeSelectionSec(instance, startSec, endSec);
  const color = regionFillColor;

  if (existing) {
    existing.setOptions?.({ start: normalized.startSec, end: normalized.endSec, color });
    return existing;
  }

  return regionsApi.addRegion?.({
    id: regionId,
    start: normalized.startSec,
    end: normalized.endSec,
    color,
    drag: true,
    resize: true,
  });
}

function onWaveformPointerDown(ev: PointerEvent) {
  // 优先处理“右键语义”（含 ctrl+click）
  if (isSecondaryPointerDown(ev)) {
    // 防抖：在 pointerdown + mousedown 同时触发的环境里避免重复启动
    if (isRightDragSelecting.value) return;
    // 阻止默认右键菜单与选择行为。
    // 注意：这里不使用 stopPropagation。
    // - stopPropagation 会让 wavesurfer 内部的一些状态机（如 dragToSeek）无法收到事件，可能引发“UI 消失/不响应”。
    // - 我们只需要阻止浏览器默认行为即可。
    ev.preventDefault();

    const instance = ws.value;
    const regionsApi = regionsPlugin.value;
    if (!instance || !regionsApi) return;
    if (!isWaveSurferInteractive(instance)) return;

    // 关键：使用 pointer capture，确保后续 pointermove/pointerup 一定能收到。
    // macOS/Electron 下右键拖动可能会离开元素区域，如果不 capture 会出现：按下有效，但拖动/松开没有事件 -> “右键没反应”。
    rightDragCapturedEl.value = ev.currentTarget as HTMLElement | null;
    rightDragCapturedPointerId.value = ev.pointerId;
    try {
      rightDragCapturedEl.value?.setPointerCapture?.(ev.pointerId);
    } catch {
      // ignore
    }

    // 右键按下：立即确定 start
    const startSec = pointerEventToTimeSec(instance, ev);
    rightDragStartSec.value = startSec;
    isRightDragSelecting.value = true;
    // 记录指针位置：用于边缘自动滚动（即使后续不再触发 pointermove）。
    quickSelectLastClientX = ev.clientX;

    // 注意：此处属于“用户主动操作 region”，不应该触发 isSyncingFromProps 的防循环。
    // 但 region.setOptions 会触发 region-updated 事件，继而 emit 到外部 startMs/endMs。
    // 这是我们想要的：右键拖拽 = 快速设置输入框数值。
    const normalized = normalizeSelectionSec(instance, startSec, startSec);
    ensureTrimRegion(regionsApi, instance, normalized.startSec, normalized.endSec);

    // 修复（可见性）：在高 zoom 下，右键按下后“选区起点”也应该立刻可见。
    // 这一步是“全过程可见”的起点：让用户从按下瞬间开始就能看到一个最小非零选区。
    scheduleQuickSelectLiveVisibility(instance, normalized.startSec, normalized.endSec);

    // 启动“边缘自动滚动”循环（如果指针不在边缘区，会在第一帧自动退出）。
    scheduleQuickSelectEdgeAutoScroll(instance);

    // 重要：这里必须显式 emit start/end。
    // 原因：wavesurfer v7 的 Regions 在 region.setOptions(...) 时，不一定会触发 region-updated 事件。
    // 如果仅依赖 regions.on('region-updated' ...) 回写输入框，会出现：
    // - 用户右键按下/拖动时 region 可能更新了，但输入框数值不变
    // - 用户会直接感知为“右键没效果/没反应”
    //
    // 因此右键快速选区这一条交互链路，采用“右键事件驱动 -> 直接 emit”来保证 UI 始终同步。
    emit('update:startMs', Math.max(0, Math.round(normalized.startSec * 1000)));
    emit('update:endMs', Math.max(0, Math.round(normalized.endSec * 1000)));

    // 捕获后续 move/up：即使鼠标移出波形，也能持续更新
    window.addEventListener('pointermove', onWaveformPointerMove, { passive: false });
    window.addEventListener('pointerup', onWaveformPointerUp, { passive: false });
    return;
  }

  // 处理左键拖拽（播放头拖拽 / 选区拖拽 / 选区指针拖拽）：
  // - 我们接管拖拽逻辑，避免 wavesurfer 内部拖拽与自动滚动“打架”。
  // - 因此这里对左键拖拽使用 preventDefault + stopPropagation。
  if (ev.button === 0) {
    const instance = ws.value;
    const regionsApi = regionsPlugin.value;
    if (!instance) return;

    // 使用 composedPath 识别命中元素，避免 DOM 结构变化导致识别失败。
    const handleLeft = findPartInComposedPath(ev, 'region-handle-left');
    const handleRight = findPartInComposedPath(ev, 'region-handle-right');
    const regionEl = findPartInComposedPath(ev, 'region');

    // 识别拖拽模式：优先 handle，其次 region，其次播放头。
    if (handleLeft) leftDragMode.value = 'region-resize-start';
    else if (handleRight) leftDragMode.value = 'region-resize-end';
    else if (regionEl) leftDragMode.value = 'region-move';
    else leftDragMode.value = 'playhead';

    // 若是 region 相关拖拽，但当前还没有 region，则退回播放头拖拽。
    const region = regionsApi ? getTrimRegion(regionsApi) : undefined;
    if ((leftDragMode.value === 'region-move' || leftDragMode.value.startsWith('region-')) && !region) {
      leftDragMode.value = 'playhead';
    }

    // 进入左键拖拽模式：接管后续指针流。
    // 说明：
    // - 如果我们无法识别到 region/handle，则视为播放头拖拽。
    // - 仍然由组件接管，以保证与自动滚动逻辑一致。
    ev.preventDefault();
    ev.stopPropagation();

    leftDragActive = true;
    leftDragMoved = false;
    leftDragStartClientX = ev.clientX;
    leftDragLastClientX = ev.clientX;

    // 记录“拖拽起点”时间（用于 region-move 的 delta 计算）
    const startSec = pointerEventToTimeSec(instance, ev);
    leftDragStartSec.value = startSec;

    // 记录 region 初始位置（用于 resize/move）。
    leftDragRegionStartSec = region?.start ?? null;
    leftDragRegionEndSec = region?.end ?? null;

    // 立即更新一次（点击即生效，避免首帧不跟手）。
    updateLeftDragTarget(instance, ev.clientX);

    // 使用 pointer capture，确保移出波形区域仍能拖拽/自动滚动。
    try {
      leftDragCapturedEl = ev.currentTarget as HTMLElement | null;
      leftDragCapturedPointerId = ev.pointerId;
      leftDragCapturedEl?.setPointerCapture?.(ev.pointerId);
    } catch {
      // ignore
    }

    scheduleLeftDragEdgeAutoScroll(instance);

    window.addEventListener('pointermove', onWaveformLeftPointerMove, { passive: false });
    window.addEventListener('pointerup', onWaveformLeftPointerUp, { passive: false });
    window.addEventListener('pointercancel', onWaveformLeftPointerUp, { passive: false });
  }
}

// 左键拖拽：仅用于边缘自动滚动（播放头/选区拖拽的体验优化）。
function onWaveformLeftPointerMove(ev: PointerEvent) {
  if (!leftDragActive) return;
  leftDragLastClientX = ev.clientX;

  if (leftDragStartClientX !== null && !leftDragMoved) {
    if (Math.abs(ev.clientX - leftDragStartClientX) >= leftDragMoveThresholdPx) {
      leftDragMoved = true;
    }
  }

  const instance = ws.value;
  if (instance) {
    // 立即根据光标更新目标（不等待自动滚动），保证“每个像素都跟手”。
    updateLeftDragTarget(instance, ev.clientX);
    // 启动/维持边缘自动滚动
    scheduleLeftDragEdgeAutoScroll(instance);
  }
}

function onWaveformLeftPointerUp() {
  if (!leftDragActive) return;
  // 释放 pointer capture（与 onWaveformPointerDown 左键分支配对）
  try {
    if (leftDragCapturedEl && leftDragCapturedPointerId !== null) {
      leftDragCapturedEl.releasePointerCapture?.(leftDragCapturedPointerId);
    }
  } catch {
    // ignore
  } finally {
    leftDragCapturedEl = null;
    leftDragCapturedPointerId = null;
  }

  leftDragActive = false;
  leftDragMoved = false;
  leftDragStartClientX = null;
  leftDragLastClientX = null;
  leftDragMode.value = 'none';
  leftDragStartSec.value = null;
  leftDragRegionStartSec = null;
  leftDragRegionEndSec = null;

  if (leftDragAutoScrollRafId !== null) cancelAnimationFrame(leftDragAutoScrollRafId);
  leftDragAutoScrollRafId = null;
  leftDragAutoScrollLastTs = null;
  leftDragScrollVelocity = 0;

  window.removeEventListener('pointermove', onWaveformLeftPointerMove);
  window.removeEventListener('pointerup', onWaveformLeftPointerUp);
  window.removeEventListener('pointercancel', onWaveformLeftPointerUp);
}

/**
 * mouse 兜底：部分 Electron/浏览器环境下 pointer 事件可能被禁用或行为不一致。
 *
 * 注意：MouseEvent 没有 pointerId 等字段，但我们这里只需要 clientX / button / ctrlKey。
 */
function onWaveformMouseDown(ev: MouseEvent) {
  // mouse 事件的右键语义同上
  const isSecondary = ev.button === 2;
  if (!isSecondary) return;
  if (isRightDragSelecting.value) return;

  ev.preventDefault();

  const instance = ws.value;
  const regionsApi = regionsPlugin.value;
  if (!instance || !regionsApi) return;
  if (!isWaveSurferInteractive(instance)) return;

  // 复用 pointer 的时间映射逻辑：这里用 PointerEvent 的字段子集（clientX/ctrlKey/button）即可
  const startSec = pointerEventToTimeSec(instance, ev as unknown as PointerEvent);
  rightDragStartSec.value = startSec;
  isRightDragSelecting.value = true;
  const normalized = normalizeSelectionSec(instance, startSec, startSec);
  ensureTrimRegion(regionsApi, instance, normalized.startSec, normalized.endSec);

  // mouse 兜底路径同上：同样显式 emit，避免依赖 regions 的事件。
  emit('update:startMs', Math.max(0, Math.round(normalized.startSec * 1000)));
  emit('update:endMs', Math.max(0, Math.round(normalized.endSec * 1000)));

  const onMove = (moveEv: MouseEvent) => onWaveformPointerMove(moveEv as unknown as PointerEvent);
  const onUp = (upEv: MouseEvent) => {
    onWaveformPointerUp(upEv as unknown as PointerEvent);
    window.removeEventListener('mousemove', onMove);
    window.removeEventListener('mouseup', onUp);
  };

  window.addEventListener('mousemove', onMove, { passive: false });
  window.addEventListener('mouseup', onUp, { passive: false });
}

function onWaveformPointerMove(ev: PointerEvent) {
  if (!isRightDragSelecting.value) return;
  // 拖拽过程中持续阻止默认行为（避免出现选择/拖拽图片等副作用）
  ev.preventDefault();

  const instance = ws.value;
  const regionsApi = regionsPlugin.value;
  if (!instance || !regionsApi) return;
  if (!isWaveSurferInteractive(instance)) return;
  if (rightDragStartSec.value === null) return;

  const currentSec = pointerEventToTimeSec(instance, ev);

  // 更新最后一次指针位置：用于边缘自动滚动。
  quickSelectLastClientX = ev.clientX;

  // 严格按需求：只“向右拖动”决定 end。
  // 若用户向左拖动，则 end 固定为 start（不反向选择）。
  const endSec = Math.max(rightDragStartSec.value, currentSec);
  const normalized = normalizeSelectionSec(instance, rightDragStartSec.value, endSec);
  ensureTrimRegion(regionsApi, instance, normalized.startSec, normalized.endSec);

  // 修复（可见性）：拖动过程中保持选区可见。
  // - rAF 合并避免高频重绘
  // - 必要时自动滚动，让 end 不会“拖到视野外但看不到”
  scheduleQuickSelectLiveVisibility(instance, normalized.startSec, normalized.endSec);

  // 若用户把指针拖到边缘，启动/维持边缘自动滚动。
  scheduleQuickSelectEdgeAutoScroll(instance);

  // 同上：显式回写输入框，确保“拖动过程中 end 实时变化”可见。
  emit('update:startMs', Math.max(0, Math.round(normalized.startSec * 1000)));
  emit('update:endMs', Math.max(0, Math.round(normalized.endSec * 1000)));
}

function onWaveformPointerUp(ev: PointerEvent) {
  if (!isRightDragSelecting.value) return;
  ev.preventDefault();

  // 释放 pointer capture（与 onWaveformPointerDown 配对）
  try {
    if (rightDragCapturedEl.value && rightDragCapturedPointerId.value !== null) {
      rightDragCapturedEl.value.releasePointerCapture?.(rightDragCapturedPointerId.value);
    }
  } catch {
    // ignore
  } finally {
    rightDragCapturedEl.value = null;
    rightDragCapturedPointerId.value = null;
  }

  const instance = ws.value;
  const regionsApi = regionsPlugin.value;

  // 记录最终选区，用于后续“强制刷新 + 滚动到可见”的修复逻辑。
  // 注意：这里不要依赖 props（因为 props 更新在 emit 之后是异步的）。
  let finalized: { startSec: number; endSec: number } | null = null;

  if (instance && regionsApi && isWaveSurferInteractive(instance) && rightDragStartSec.value !== null) {
    const currentSec = pointerEventToTimeSec(instance, ev);
    const endSec = Math.max(rightDragStartSec.value, currentSec);
    const normalized = normalizeSelectionSec(instance, rightDragStartSec.value, endSec);
    ensureTrimRegion(regionsApi, instance, normalized.startSec, normalized.endSec);

    emit('update:startMs', Math.max(0, Math.round(normalized.startSec * 1000)));
    emit('update:endMs', Math.max(0, Math.round(normalized.endSec * 1000)));

    finalized = { startSec: normalized.startSec, endSec: normalized.endSec };
  }

  isRightDragSelecting.value = false;
  rightDragStartSec.value = null;
  window.removeEventListener('pointermove', onWaveformPointerMove);
  window.removeEventListener('pointerup', onWaveformPointerUp);

  // 修复：高倍缩放下，右键快捷选区后选区偶发不可见。
  // - 该问题的一个强信号是：用户轻微调整 zoom 后选区会立刻出现。
  // - 因此这里在 finalize 后主动触发一次“no-op zoom”重绘，并滚动让选区进入视野。
  if (finalized && instance) {
    ensureQuickSelectRegionVisible(instance, finalized.startSec, finalized.endSec);
  }

  // 结束拖动后，清空 live 刷新队列，避免下一帧还在尝试滚动。
  pendingQuickSelectRange = null;

  // 停止边缘自动滚动循环并清理状态。
  if (quickSelectAutoScrollRafId !== null) cancelAnimationFrame(quickSelectAutoScrollRafId);
  quickSelectAutoScrollRafId = null;
  quickSelectAutoScrollLastTs = null;
  quickSelectLastClientX = null;
}

function destroyWaveSurfer() {
  isReady.value = false;
  hasError.value = false;
  durationMs.value = null;

  regionsPlugin.value = null;

  try {
    ws.value?.destroy();
  } catch {
    // ignore
  }
  ws.value = null;
}

function pointerEventToVolume(ev: PointerEvent, host: HTMLElement): number {
  // 将鼠标位置映射为 dB：
  // - 顶部 => +18 dB
  // - 底部 => -18 dB
  // 然后换算为 cut.volume（Base=1.6 的指数音量）。
  //
  // 关键修复（音量线与波形中线不对齐的根因）：
  // - 不能使用 host.getBoundingClientRect().height，因为它包含滚动条高度。
  // - 必须使用 props.height（波形绘制区域的实际高度）作为映射基准。
  // - 这样 volume=0 才会精确对应波形纵向中点。
  const rect = host.getBoundingClientRect();
  const waveformHeight = props.height;
  const y = Math.max(0, Math.min(waveformHeight, ev.clientY - rect.top));
  const ratio = 1 - y / waveformHeight; // 0..1（顶端=1）
  const db = volumeDbMin + ratio * (volumeDbMax - volumeDbMin);
  return dbToVolume(db);
}

function onVolumePointerDown(ev: PointerEvent) {
  // 仅响应主键拖动（避免与右键选区冲突）
  if (ev.button !== 0) return;
  // 关键：阻止波形区域的默认交互（dragToSeek 等）接管该手势。
  // 否则在部分环境下会出现“拖音量线时播放头在动/音量线不动”，用户会感知为“无法拖动”。
  ev.preventDefault();
  ev.stopPropagation();

  const host = containerEl.value;
  if (!host) return;

  isVolumeDragging.value = true;
  volumeDragPointerId.value = ev.pointerId;
  volumeDragTargetEl.value = ev.currentTarget as HTMLElement | null;

  try {
    volumeDragTargetEl.value?.setPointerCapture?.(ev.pointerId);
  } catch {
    // ignore
  }

  const nextVolume = pointerEventToVolume(ev, host);
  emit('update:volume', Number(nextVolume.toFixed(3)));

  window.addEventListener('pointermove', onVolumePointerMove, { passive: false });
  window.addEventListener('pointerup', onVolumePointerUp, { passive: false });
  // pointercancel 在以下场景可能触发：
  // - 触控板/浏览器手势接管了指针流
  // - OS 层中断（弹出菜单、窗口切换等）
  // 如果不处理，用户会感知为“拖动时灵时不灵/有时突然失效”。
  window.addEventListener('pointercancel', onVolumePointerCancel, { passive: false });
}

function onVolumePointerMove(ev: PointerEvent) {
  if (!isVolumeDragging.value) return;
  ev.preventDefault();
  ev.stopPropagation();

  const host = containerEl.value;
  if (!host) return;

  const nextVolume = pointerEventToVolume(ev, host);
  emit('update:volume', Number(nextVolume.toFixed(3)));
}

function onVolumePointerCancel(ev: PointerEvent) {
  // pointercancel 在以下场景可能触发：
  // - 触控板/浏览器手势接管了指针流
  // - OS 层中断（弹出菜单、窗口切换等）
  // 若不统一走清理逻辑，会留下 window listener 与 dragging 状态，导致后续拖动异常。
  onVolumePointerUp(ev);
}

function onVolumePointerUp(ev: PointerEvent) {
  if (!isVolumeDragging.value) return;
  ev.preventDefault();
  ev.stopPropagation();

  try {
    if (volumeDragTargetEl.value && volumeDragPointerId.value !== null) {
      volumeDragTargetEl.value.releasePointerCapture?.(volumeDragPointerId.value);
    }
  } catch {
    // ignore
  } finally {
    isVolumeDragging.value = false;
    volumeDragPointerId.value = null;
    volumeDragTargetEl.value = null;
  }

  window.removeEventListener('pointermove', onVolumePointerMove);
  window.removeEventListener('pointerup', onVolumePointerUp);
  window.removeEventListener('pointercancel', onVolumePointerCancel);
}

/**
 * 将 SDK 的 cut.volume（Base=1.6 的指数音量）映射为前端试听用的线性 gain。
 *
 * SDK：amplitude = 1.6 ^ volume
 *
 * 前端：使用 WebAudio backend 时，可以设置 gain > 1 来模拟“增益”。
 * 本项目以 SDK 预览为准，因此这里不做额外 clamp（与 SDK 的 Base=1.6 行为保持一致）。
 */
function sdkVolumeToFrontendGain(volume: number): number {
  const base = 1.6;
  // SDK 预览模式使用 beep/effects.Volume：
  // amplitude = Base ^ Volume （Base=1.6）。
  // 前端预览以 SDK 预览为准，因此这里不做额外 clamp，保持与 SDK 一致的增益曲线。
  return Math.pow(base, volume);
}

function applyFrontendPreviewVolume() {
  const instance = ws.value;
  if (!instance) return;
  // 注意：wavesurfer 的 setVolume 语义随 backend 不同而不同。
  // - backend=WebAudio：更像 gain
  // - backend=MediaElement：会被浏览器 clamp 到 0..1
  instance.setVolume(sdkVolumeToFrontendGain(props.volume ?? 0));
}

function formatMs(ms: number): string {
  // 便于剪辑定位：同时显示秒与毫秒
  const s = (ms / 1000).toFixed(3);
  return `${s}s (${ms}ms)`;
}

function stop() {
  const instance = ws.value;
  if (!instance) return;
  instance.pause();
  instance.setTime(0);
  isPlaying.value = false;
}

function togglePlay() {
  const instance = ws.value;
  if (!instance) return;
  if (!isReady.value) return;

  // 规则：默认“播放全部”。当用户显式切换为“播放选区”时，才按选区播放。
  // 不做循环播放（用户明确不需要）。
  if (isPlaying.value) {
    instance.pause();
    return;
  }

  if (playScope.value === 'selection' && selectionDurationMs.value !== null) {
    instance.play(props.startMs / 1000, props.endMs / 1000);
  } else {
    instance.play();
  }
}

function bindWheelZoom() {
  // Ctrl + 滚轮缩放：属于加成体验，不影响 slider 主入口。
  //
  // 用户反馈：旧实现“步进跳变大 + 无锚点”，体验很差。
  // 这里改为：
  // 1) 使用指数缩放（更接近剪辑软件/地图缩放的手感）
  // 2) 使用 requestAnimationFrame 节流（避免 trackpad 高频 wheel 造成抖动）
  // 3) 以鼠标所在时间点为锚点（缩放后该时间点仍留在鼠标位置附近，不会“飘”）
  const el = containerEl.value;
  if (!el) return;

  let rafId: number | null = null;
  let accumulatedDeltaY = 0;
  let lastClientX = 0;

  const applyZoomFrame = () => {
    rafId = null;

    const instance = ws.value;
    if (!instance || !isReady.value) {
      accumulatedDeltaY = 0;
      return;
    }

    // 计算锚点时间（缩放前）：用当前 DOM 宽度/滚动来映射，确保与用户视觉一致。
    const scrollEl = resolveWaveformScrollElement(instance);
    const rect = scrollEl.getBoundingClientRect();
    const xInView = Math.max(0, Math.min(rect.width, lastClientX - rect.left));
    const anchorTimeSec = pointerEventToTimeSec(
      instance,
      ({ clientX: lastClientX, button: 0 } as unknown as PointerEvent)
    );

    // 指数缩放系数：deltaY>0(向下滚) => 缩小；deltaY<0(向上滚) => 放大
    // 系数选择：0.002 在触控板与鼠标滚轮上都相对平滑；可按体验继续微调。
    const factor = Math.exp(-accumulatedDeltaY * 0.002);
    accumulatedDeltaY = 0;

    const currentZoom = zoomMinPxPerSec.value;
    const targetZoom = Math.max(zoomMin, Math.min(zoomMax, Math.round(currentZoom * factor)));
    if (targetZoom === currentZoom) return;

    // 通过更新响应式 zoom 值触发后续的 instance.zoom(...)（见下方 watch）。
    zoomMinPxPerSec.value = targetZoom;

    // 锚点滚动校正：在 zoom 生效并重绘后，把 scrollLeft 调整到让锚点时间仍对应鼠标位置。
    // 注意：重绘是异步的，scrollWidth 会在下一帧/下下帧才稳定。
    requestAnimationFrame(() => {
      const dur = instance.getDuration();
      if (!Number.isFinite(dur) || dur <= 0) return;

      const newScrollEl = resolveWaveformScrollElement(instance);
      const newTotalWidth = newScrollEl.scrollWidth || newScrollEl.clientWidth || 0;
      if (!Number.isFinite(newTotalWidth) || newTotalWidth <= 0) return;

      const desiredGlobalX = (anchorTimeSec / dur) * newTotalWidth;
      const desiredScrollLeft = Math.max(0, desiredGlobalX - xInView);
      newScrollEl.scrollLeft = desiredScrollLeft;
    });
  };

  // 判断是否触发“缩放手势”。
  // macOS：
  // - 优先使用 getModifierState('Control') + ev.ctrlKey
  // - 若外接键盘未在 wheel 事件里反映 ctrlKey，则使用 window 记录的修饰键状态兜底
  // - 某些键盘把 Control 映射为 Meta/Command，作为兜底允许
  // 其他平台：仅识别 Ctrl。
  const isZoomGesture = (ev: WheelEvent) => {
    if (Platform.is.mac) {
      return (
        ev.getModifierState?.('Control') ||
        ev.ctrlKey ||
        macModifierState.control ||
        ev.getModifierState?.('Meta') ||
        ev.metaKey ||
        macModifierState.meta
      );
    }
    return ev.ctrlKey;
  };

  const onWheel = (ev: WheelEvent) => {
    if (!isZoomGesture(ev)) return;
    // 方向确认：用户选择 3A（向上滚 = 放大；向下滚 = 缩小）。这与 deltaY 的标准方向一致。
    ev.preventDefault();

    // 聚合 delta，交给 RAF 一次性计算，减少“齿轮感/跳变”。
    accumulatedDeltaY += ev.deltaY;
    lastClientX = ev.clientX;
    if (rafId === null) rafId = requestAnimationFrame(applyZoomFrame);
  };

  // passive:false 才能 preventDefault
  // 注意：为避免 TS 在 add/remove 的 listener 类型上产生不必要的噪音，这里显式转为 EventListener。
  const listener = onWheel as unknown as EventListener;
  el.addEventListener('wheel', listener, { passive: false });
  return () => {
    if (rafId !== null) cancelAnimationFrame(rafId);
    el.removeEventListener('wheel', listener);
  };
}

async function initWaveSurfer() {
  destroyWaveSurfer();

  if (!containerEl.value) return;
  if (!audioUrl.value) {
    // 没选源文件时不报错，只是不初始化
    return;
  }

  try {
    const regions = RegionsPlugin.create() as unknown as RegionsApi;
    regionsPlugin.value = regions;

    const instance = WaveSurfer.create({
      container: containerEl.value,
      url: audioUrl.value,
      height: props.height,
      waveColor: '#94a3b8',
      progressColor: '#64748b',
      cursorColor: '#0ea5e9',
      cursorWidth: 2,
      // 纵向高度缩放（视觉反馈）：
      // - wavesurfer 的波形渲染支持 barHeight（渲染阶段缩放）。
      // - 这比对 canvas 做 transform 更干净，不会出现“静音中位线变粗/糊”的伪影。
      // - 初始值使用当前 props.volume 映射出来的 scale。
      barHeight: waveformYScale.value,
      // 波形显示不做 normalize：
      // - normalize 会把全曲最大峰值拉满高度。
      // - 在高缩放（minPxPerSec 较大）时，底噪/静音段的细小波动会被“视觉夸张”，
      //   用户会误以为静音段也有明显波形（反馈为“缩放>=100后波形异常”）。
      // - 关闭 normalize 后，波形高度更接近真实相对振幅：小噪声仍会显示，但不会被夸张放大。
      normalize: false,
      interact: true,
      dragToSeek: true,
      autoScroll: true,
      autoCenter: false,
      minPxPerSec: zoomMinPxPerSec.value,
      backend: 'WebAudio',
      plugins: [regions],
    });

    ws.value = instance;

    instance.on('error', () => {
      hasError.value = true;
    });

    // 播放状态与时间：用于“剪辑软件式”的播放条显示
    instance.on('play', () => {
      isPlaying.value = true;
    });
    instance.on('pause', () => {
      isPlaying.value = false;
    });
    instance.on('finish', () => {
      isPlaying.value = false;
    });
    instance.on('timeupdate', (t: number) => {
      currentTimeMs.value = Math.max(0, Math.round(t * 1000));
    });

    instance.on('ready', () => {
      isReady.value = true;
      hasError.value = false;
      const dur = instance.getDuration();
      if (Number.isFinite(dur) && dur > 0) {
        durationMs.value = dur * 1000;
        emit('loaded', { durationMs: durationMs.value });
      }

      // 音量：尽量贴近 SDK 的 cut.volume 听感（映射为前端 gain）
      applyFrontendPreviewVolume();
      // 初始化波形纵向缩放（渲染参数 barHeight）。
      // 说明：虽然 create() 里已传入 barHeight，但在某些情况下（例如 ready 后首次完整渲染）
      // 再做一次 schedule 可以确保视觉状态与当前 volume 完全一致。
      scheduleApplyWaveformBarHeight();

      // 允许拖拽创建 region
      try {
        regions.enableDragSelection?.({
          color: regionFillColor,
        });
      } catch {
        // ignore
      }

      // 如果外部已经有合法的裁剪区间，则创建/对齐 region
      syncRegionFromProps();

      // 右键拖拽快速选区（无菜单）：
      // 事件绑定在模板的 .waveform 容器上（@pointerdown/@mousedown/@contextmenu），
      // 避免 wavesurfer 内部 DOM 结构变化导致绑定失效。

      // 监听 region 变化并回写到 props
      regions.on?.('region-created', (region) => {
        if (region.id !== regionId) {
          // 只保留一个 region
          try {
            const list = regions.getRegions?.() || [];
            list.forEach((r) => {
              if (r.id !== regionId) r.remove?.();
            });
          } catch {
            // ignore
          }
          // 使用 setOptions 更新 id（比直接赋值更符合插件语义）
          region.setOptions?.({ id: regionId });
        }

        if (isSyncingFromProps.value) return;
        emitStartEndFromRegion(region);
      });

      regions.on?.('region-updated', (region) => {
        if (region.id !== regionId) return;
        if (isSyncingFromProps.value) return;
        emitStartEndFromRegion(region);
      });

      regions.on?.('region-update-end', (region) => {
        if (region.id !== regionId) return;
        if (isSyncingFromProps.value) return;
        emitStartEndFromRegion(region);
      });

      // wavesurfer destroy 时也需要解绑 listener
      instance.on('destroy', () => {
        window.removeEventListener('pointermove', onWaveformPointerMove);
        window.removeEventListener('pointerup', onWaveformPointerUp);
        window.removeEventListener('pointermove', onVolumePointerMove);
        window.removeEventListener('pointerup', onVolumePointerUp);
        isRightDragSelecting.value = false;
        rightDragStartSec.value = null;
        rightDragCapturedEl.value = null;
        rightDragCapturedPointerId.value = null;
        isVolumeDragging.value = false;
        volumeDragPointerId.value = null;
        volumeDragTargetEl.value = null;
      });
    });

    // 初始绑定 Ctrl+滚轮缩放
    const unbindWheel = bindWheelZoom();
    // wavesurfer destroy 时，container 也会被重建，确保解绑
    instance.on('destroy', () => {
      unbindWheel?.();
    });
  } catch {
    hasError.value = true;
  }
}

function emitStartEndFromRegion(region: Region) {
  const startMs = Math.max(0, Math.round(region.start * 1000));
  const endMs = Math.max(0, Math.round(region.end * 1000));

  emit('update:startMs', startMs);
  emit('update:endMs', endMs);
}

function syncRegionFromProps() {
  const instance = ws.value;
  if (!instance) return;
  if (!isReady.value) return;

  const regionsApi = regionsPlugin.value;
  if (!regionsApi) return;

  const startMs = props.startMs;
  const endMs = props.endMs;

  // 负数直接认为不合法（输入框已有错误提示，这里不做“猜测修正”）
  if (startMs < 0 || endMs < 0) return;

  const startSec = startMs / 1000;

  // 关键体验修复：
  // - 过去在 end<=start 时直接 return，导致用户“先改 start（或先改 end）”完全看不到任何波形变化。
  // - 这会被用户感知为“没效果”。
  //
  // 现在改为：当 end<=start 时，仍然渲染一个“可操作的非零选区”，以确保两侧拖动指针始终可用。
  // 说明：这属于渲染层的非零化（1ms），用于避免 Regions 退化为 marker。
  // 在用户输入尚未完成时（例如先改 start 或先改 end），让 UI 保持可见且可编辑。
  const rawEndSec = endMs > startMs ? endMs / 1000 : startSec;
  const normalized = normalizeSelectionSec(instance, startSec, rawEndSec);
  const color = regionFillColor;

  isSyncingFromProps.value = true;
  try {
    const list = regionsApi.getRegions?.() || [];
    const existing = list.find((r) => r.id === regionId);
    if (existing) {
      // 关键修复点：wavesurfer v7 的 Region 没有 update()，应使用 setOptions()
      existing.setOptions?.({ start: normalized.startSec, end: normalized.endSec, color });
    } else {
      regionsApi.addRegion?.({
        id: regionId,
        start: normalized.startSec,
        end: normalized.endSec,
        color,
        drag: true,
        resize: true,
      });
    }
  } catch {
    // ignore
  } finally {
    // 避免同步回写产生的事件循环
    setTimeout(() => {
      isSyncingFromProps.value = false;
    }, 0);
  }
}

onMounted(() => {
  initWaveSurfer();

  // macOS 外接键盘修饰键兜底：在 window 级别跟踪 Control/Command 按键状态。
  if (Platform.is.mac) {
    window.addEventListener('keydown', onMacModifierKeyDown);
    window.addEventListener('keyup', onMacModifierKeyUp);
    window.addEventListener('blur', resetMacModifierState);
  }
});

onBeforeUnmount(() => {
  if (Platform.is.mac) {
    window.removeEventListener('keydown', onMacModifierKeyDown);
    window.removeEventListener('keyup', onMacModifierKeyUp);
    window.removeEventListener('blur', resetMacModifierState);
    resetMacModifierState();
  }

  // 右键快捷选区的 live 可见性刷新：避免组件销毁后 RAF 回调仍运行。
  if (quickSelectLiveRafId !== null) cancelAnimationFrame(quickSelectLiveRafId);
  quickSelectLiveRafId = null;
  pendingQuickSelectRange = null;

  // 右键快捷选区的边缘自动滚动：避免组件销毁后仍在滚动。
  if (quickSelectAutoScrollRafId !== null) cancelAnimationFrame(quickSelectAutoScrollRafId);
  quickSelectAutoScrollRafId = null;
  quickSelectAutoScrollLastTs = null;
  quickSelectLastClientX = null;

  // 左键拖拽的边缘自动滚动：避免组件销毁后仍在滚动。
  if (leftDragAutoScrollRafId !== null) cancelAnimationFrame(leftDragAutoScrollRafId);
  leftDragAutoScrollRafId = null;
  leftDragAutoScrollLastTs = null;
  leftDragActive = false;
  leftDragMoved = false;
  leftDragStartClientX = null;
  leftDragLastClientX = null;

  // 避免组件销毁后 RAF 回调仍尝试操作 instance。
  if (barHeightRafId !== null) cancelAnimationFrame(barHeightRafId);
  barHeightRafId = null;
  pendingBarHeight = null;
  destroyWaveSurfer();
});

watch(
  () => audioUrl.value,
  () => {
    initWaveSurfer();
  }
);

// 当音量 slider 改变时，实时更新前端试听音量
watch(
  () => props.volume,
  () => {
    applyFrontendPreviewVolume();
    // 音量变化时同步更新波形纵向缩放（视觉反馈）：
    // - 使用 barHeight（渲染阶段缩放）而非 canvas transform（避免伪影）。
    scheduleApplyWaveformBarHeight();
  }
);

// zoom slider 改变时，调用 wavesurfer.zoom
watch(
  () => zoomMinPxPerSec.value,
  () => {
    const instance = ws.value;
    if (!instance) return;
    if (!isReady.value) return;
    instance.zoom(zoomMinPxPerSec.value);
  }
);

watch(
  () => [props.startMs, props.endMs],
  () => {
    syncRegionFromProps();
  }
);
</script>

<style scoped lang="scss">
.waveform-wrapper {
  @apply mt-2;
}

.waveform-host {
  // 波形主区域需要可伸缩（flex:1），并允许内部滚动条工作。
  @apply relative;
  @apply flex-1;
  @apply min-w-0;
}

.waveform-shell {
  // 关键：波形宽度保持“原样”，不被刻度挤压。
  // 刻度通过 overlay 贴在波形外侧显示，不参与布局宽度。
  @apply relative;
  @apply overflow-visible;
}

.waveform {
  @apply w-full;
  @apply rounded-md;
  @apply bg-zinc-50;
  @apply border border-zinc-200;

  // 关键：zoom 后波形会变“很长”，需要允许横向滚动。
  // 这里不使用 overflow-hidden，避免把滚动条与超出部分裁掉。
  @apply overflow-x-auto;
  @apply overflow-y-hidden;

  // 轻量滚动条样式（仅做基础可用性，不追求强定制）
  @apply [&::-webkit-scrollbar]:h-2;
  @apply [&::-webkit-scrollbar-track]:bg-zinc-200/30;
  @apply [&::-webkit-scrollbar-thumb]:bg-zinc-900/30;
  @apply [&::-webkit-scrollbar-thumb]:rounded-full;
  @apply [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50;
}

// 视觉增强：音量变化时，让波形“高度”产生适度反馈。
//
// 当前实现说明：
// - 不使用对 canvas 做 transform(scaleY) 的方式：
//   该方式会引入插值伪影（例如静音中位线变粗/发糊），观感不佳。
// - 改用 wavesurfer 的渲染参数 barHeight：
//   这是渲染阶段缩放，线条更干净，也更符合用户对“波形高度变化”的直觉。
//
// 代码实现位于脚本区：scheduleApplyWaveformBarHeight()。

// 音量快捷调节条
//
// 定位说明：
// - 使用 top:0 left:0 right:0 而非 inset-0，因为高度由 :style 显式指定（等于波形绘制区域高度）。
// - 这样音量线的百分比定位才能正确映射到波形区域，而不会被滚动条撑大。
// shell 级 overlay：贯穿左右标尺 + 波形，确保蓝线能“碰到”标尺边界
.volume-overlay-shell {
  @apply absolute;
  left: 0;
  right: 0;
  top: 0;
  @apply pointer-events-none;
  @apply z-10;
}

// 左右 dB 标尺与当前值显示（overlay，不占用布局宽度）
.volume-scale-overlay {
  @apply absolute top-0;
  @apply text-[10px] text-gray-600;
  @apply select-none;
  @apply pointer-events-none;
  @apply z-20;
}

.volume-scale-left {
  // 左侧刻度在波形外侧，紧贴边界并向右对齐。
  // 使用 right:100% 保证“右侧对齐到波形左边界”。
  right: 100%;
  margin-right: var(--volumeScaleOffset);
  @apply text-right;
  @apply w-[var(--volumeScaleW)];
}

.volume-scale-right {
  // 当前 dB 在波形右外侧，贴边显示。
  left: 100%;
  margin-left: var(--volumeScaleOffset);
  @apply text-left;
  @apply w-[var(--volumeScaleW)];
}

.volume-scale-mark {
  @apply absolute;
  transform: translateY(-50%);
  @apply flex items-center;
  @apply tabular-nums;
  @apply w-full;
  @apply justify-end;
}


.volume-scale-text {
  // 让刻度数字更“器材感”：
  // - tabular-nums 保证数字宽度一致
  // - 字重略提升，避免过轻导致看不清
  @apply text-[10px];
  @apply font-medium;
  @apply text-right;
  @apply block;
  // 固定最小宽度，保证正负刻度（含 + 号）视觉宽度一致。
  min-width: 3ch;
}

.volume-scale-mark.is-zero .volume-scale-text {
  // 0dB（unity gain）是用户最常对齐的参考点，视觉上略强调。
  @apply text-gray-700;
}



.volume-current-db {
  @apply absolute;
  @apply text-[11px] text-gray-700;
  transform: translateY(-50%);
}

.volume-overflow-inline {
  // 波形区域内的超范围提示：轻量、靠右上角（避免遮挡内容与刻度）。
  @apply absolute;
  right: 6px;
  top: 2px;
  @apply text-[10px] text-gray-500;
  @apply pointer-events-none;
}

.volume-zero-line {
  // 0 基准线：永远在中位（50%）。
  // 该线仅作为视觉参照，不参与交互。
  @apply absolute left-0 right-0;
  top: 50%;
  @apply h-px;
  @apply bg-zinc-400/35;
}

.volume-line {
  @apply absolute left-0 right-0;
  @apply h-px;
  @apply bg-sky-500/70;
}

// .volume-handle 已删除：用户反馈蓝色小圆点无用，现仅保留横线作为视觉指示。

.volume-line-hit {
  // 命中区位于 waveform-host 内部，因此 left/right 直接撑满即可。
  @apply absolute left-0 right-0;
  // 命中区做大：这是“拖动时灵时不灵”的主要根因。
  // 之前仅 3px 很难点中，用户稍微偏一点就点到波形本体（被 wavesurfer 接管），导致看起来像没响应。
  @apply h-6;
  @apply -translate-y-1/2;
  @apply pointer-events-auto;
  @apply cursor-ns-resize;
  @apply bg-transparent;

  // 命中区必须处于波形本体之上，确保能吃到 pointerdown。
  @apply z-20;

  // 防止浏览器把拖动当成滚动/选择手势（尤其是触控板 + Electron 环境）
  touch-action: none;
  user-select: none;
}
</style>
