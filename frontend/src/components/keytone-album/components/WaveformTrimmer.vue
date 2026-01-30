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
      ref="containerEl"
      class="waveform"
      @pointerdown.capture="onWaveformPointerDown"
      @mousedown.capture="onWaveformMouseDown"
      @contextmenu.prevent
    />

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
const zoomMinPxPerSec = ref(200);

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
  if (instance && regionsApi && isWaveSurferInteractive(instance) && rightDragStartSec.value !== null) {
    const currentSec = pointerEventToTimeSec(instance, ev);
    const endSec = Math.max(rightDragStartSec.value, currentSec);
    const normalized = normalizeSelectionSec(instance, rightDragStartSec.value, endSec);
    ensureTrimRegion(regionsApi, instance, normalized.startSec, normalized.endSec);

    emit('update:startMs', Math.max(0, Math.round(normalized.startSec * 1000)));
    emit('update:endMs', Math.max(0, Math.round(normalized.endSec * 1000)));
  }

  isRightDragSelecting.value = false;
  rightDragStartSec.value = null;
  window.removeEventListener('pointermove', onWaveformPointerMove);
  window.removeEventListener('pointerup', onWaveformPointerUp);
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

/**
 * 将 SDK 的 cut.volume（Base=1.6 的指数音量）映射为前端试听用的线性 gain。
 *
 * SDK：amplitude = 1.6 ^ volume
 *
 * 前端：使用 WebAudio backend 时，可以设置 gain > 1 来模拟“增益”。
 * 为避免爆音与设备差异，这里加一个保守上限（可按需再调）。
 */
function sdkVolumeToFrontendGain(volume: number): number {
  const base = 1.6;
  const gain = Math.pow(base, volume);
  // 保护：避免极端值导致爆音/失真
  return Math.max(0, Math.min(gain, 4));
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
  // Ctrl + 滚轮缩放：属于加成体验，不影响 slider 主入口
  const el = containerEl.value;
  if (!el) return;

  const onWheel = (ev: WheelEvent) => {
    if (!ev.ctrlKey) return;
    ev.preventDefault();

    // deltaY > 0 往下滚：缩小；deltaY < 0 往上滚：放大
    const direction = ev.deltaY > 0 ? -1 : 1;
    const next = Math.max(zoomMin, Math.min(zoomMax, zoomMinPxPerSec.value + direction * 50));
    zoomMinPxPerSec.value = next;
  };

  // passive:false 才能 preventDefault
  // 注意：为避免 TS 在 add/remove 的 listener 类型上产生不必要的噪音，这里显式转为 EventListener。
  const listener = onWheel as unknown as EventListener;
  el.addEventListener('wheel', listener, { passive: false });
  return () => el.removeEventListener('wheel', listener);
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
      normalize: true,
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
        isRightDragSelecting.value = false;
        rightDragStartSec.value = null;
        rightDragCapturedEl.value = null;
        rightDragCapturedPointerId.value = null;
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
</style>
