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
      class="waveform-host"
    >
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
      <div
        class="volume-overlay"
        aria-hidden="true"
        :style="{ height: `${height}px` }"
      >
        <!--
          0 基准线（永远位于中位）：
          - 这是“UI 对齐认知”的辅助线：让用户一眼知道“0=中位（unity gain）”。
          - 当当前音量刚好为 0 时，蓝色音量线会与该线重合。
        -->
        <div class="volume-zero-line" />
        <div class="volume-line" :style="{ top: volumeLineTop }" />
      </div>

      <!--
        交互命中区必须放在 pointer-events 可用的层级中。

        说明：
        - 我们希望“音量线/手柄”不阻挡波形本体的交互（seek、右键选区等），因此视觉层 .volume-overlay
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
    </div>

    <div v-if="isReady && durationMs" class="text-[12px] text-gray-500 mt-1 flex flex-wrap gap-x-3 gap-y-1">
      <div>{{ t('KeyToneAlbum.defineSounds.waveformTrimmer.current') }}{{ formatMs(currentTimeMs) }}</div>
      <div>{{ t('KeyToneAlbum.defineSounds.waveformTrimmer.total') }}{{ formatMs(Math.round(durationMs)) }}</div>
      <div v-if="selectionDurationMs !== null">
        {{ t('KeyToneAlbum.defineSounds.waveformTrimmer.selection') }}{{ formatMs(selectionDurationMs) }}
      </div>
      <div class="text-gray-400">{{ t('KeyToneAlbum.defineSounds.waveformTrimmer.hint') }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import WaveSurfer from 'wavesurfer.js';
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
// - 使用 cut.volume 的“原始尺度”（Base=1.6 的指数音量），不转换为 dB。
// - 以 SDK 预览为准：前端预览采用同一映射，不再额外 clamp。
//
// 设计：
// - 为避免无穷范围带来的 UI 不可控，这里定义“可视范围”用于拖拽/指示线位置映射。
// - 该范围不限制用户输入的数值：手动输入超出范围时，指示线会贴边显示；
//   用户拖拽时则会把值限定在可视范围内（这是交互层的合理约束）。
// ============================================================================
// “可视范围”的定义必须保证 0 永远在中位。
//
// 用户诉求：UI 上 volume=0 一定处于中位，但不影响实际的 cut.volume 语义。
// 因此这里采用“对称区间”来定义拖拽映射：[-X, +X]。
// - 这只是 UI 映射范围：用于拖拽与指示线定位。
// - 不限制用户通过输入框直接输入更大/更小的值（超出范围时指示线会贴边显示）。
const volumeVisibleAbsMax = 10;
const volumeMin = -volumeVisibleAbsMax;
const volumeMax = volumeVisibleAbsMax;

const isVolumeDragging = ref(false);
const volumeDragPointerId = ref<number | null>(null);
const volumeDragTargetEl = ref<HTMLElement | null>(null);

const volumeLineTop = computed(() => {
  const raw = props.volume ?? 0;
  const clamped = Math.max(volumeMin, Math.min(volumeMax, raw));
  const ratio = (clamped - volumeMin) / (volumeMax - volumeMin); // 0..1
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
  const raw = props.volume ?? 0;
  const clamped = Math.max(volumeMin, Math.min(volumeMax, raw));
  const ratio = (clamped - volumeMin) / (volumeMax - volumeMin);
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

    // 2) 仅在需要时滚动：避免每帧 scroll 造成跳动。
    //    这里不追求“强制居中”，只要选区落在视野里即可。
    if (!isTimeRangeVisible(instance, next.startSec, next.endSec, 24)) {
      // reRender 可能会在下一帧才稳定 scrollWidth，因此滚动放到下一帧做。
      requestAnimationFrame(() => {
        scrollTimeRangeIntoView(instance, next.startSec, next.endSec);
      });
    }
  });
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
  // 仅响应“右键语义”（含 ctrl+click）
  if (!isSecondaryPointerDown(ev)) return;
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

  // 注意：此处属于“用户主动操作 region”，不应该触发 isSyncingFromProps 的防循环。
  // 但 region.setOptions 会触发 region-updated 事件，继而 emit 到外部 startMs/endMs。
  // 这是我们想要的：右键拖拽 = 快速设置输入框数值。
  const normalized = normalizeSelectionSec(instance, startSec, startSec);
  ensureTrimRegion(regionsApi, instance, normalized.startSec, normalized.endSec);

  // 修复（可见性）：在高 zoom 下，右键按下后“选区起点”也应该立刻可见。
  // 这一步是“全过程可见”的起点：让用户从按下瞬间开始就能看到一个最小非零选区。
  scheduleQuickSelectLiveVisibility(instance, normalized.startSec, normalized.endSec);

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

  // 严格按需求：只“向右拖动”决定 end。
  // 若用户向左拖动，则 end 固定为 start（不反向选择）。
  const endSec = Math.max(rightDragStartSec.value, currentSec);
  const normalized = normalizeSelectionSec(instance, rightDragStartSec.value, endSec);
  ensureTrimRegion(regionsApi, instance, normalized.startSec, normalized.endSec);

  // 修复（可见性）：拖动过程中保持选区可见。
  // - rAF 合并避免高频重绘
  // - 必要时自动滚动，让 end 不会“拖到视野外但看不到”
  scheduleQuickSelectLiveVisibility(instance, normalized.startSec, normalized.endSec);

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
  // 将鼠标位置映射为音量值：
  // - 顶部 => volumeMax
  // - 底部 => volumeMin
  //
  // 关键修复（音量线与波形中线不对齐的根因）：
  // - 不能使用 host.getBoundingClientRect().height，因为它包含滚动条高度。
  // - 必须使用 props.height（波形绘制区域的实际高度）作为映射基准。
  // - 这样 volume=0 才会精确对应波形纵向中点。
  const rect = host.getBoundingClientRect();
  const waveformHeight = props.height;
  const y = Math.max(0, Math.min(waveformHeight, ev.clientY - rect.top));
  const ratio = 1 - y / waveformHeight; // 0..1（顶端=1）
  return volumeMin + ratio * (volumeMax - volumeMin);
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

  const onWheel = (ev: WheelEvent) => {
    if (!ev.ctrlKey) return;
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
});

onBeforeUnmount(() => {
  // 右键快捷选区的 live 可见性刷新：避免组件销毁后 RAF 回调仍运行。
  if (quickSelectLiveRafId !== null) cancelAnimationFrame(quickSelectLiveRafId);
  quickSelectLiveRafId = null;
  pendingQuickSelectRange = null;

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
  @apply relative;
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
.volume-overlay {
  @apply absolute top-0 left-0 right-0;
  @apply pointer-events-none;
  @apply z-10;
}

.volume-zero-line {
  // 0 基准线：永远在中位（50%）。
  // 该线仅作为视觉参照，不参与交互。
  @apply absolute left-0 right-0;
  top: 50%;
  @apply h-px;
  @apply bg-zinc-400/40;
}

.volume-line {
  @apply absolute left-0 right-0;
  @apply h-px;
  @apply bg-sky-500/70;
}

// .volume-handle 已删除：用户反馈蓝色小圆点无用，现仅保留横线作为视觉指示。

.volume-line-hit {
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
