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
    <div v-if="!hasSource" class="p-2">
      <div class="text-[12px] text-gray-500">请选择音频源文件以显示波形。</div>
    </div>

    <div v-else-if="!isReady && !hasError" class="p-2">
      <q-skeleton type="rect" height="92px" />
      <div class="text-[12px] text-gray-500 mt-2">正在加载波形...</div>
    </div>

    <div v-if="hasError" class="p-2">
      <q-banner dense rounded class="bg-orange-50 text-orange-900"> 波形不可用，已降级为手动裁剪。 </q-banner>
    </div>

    <div v-show="hasSource" ref="containerEl" class="waveform" />

    <div v-if="isReady && durationMs" class="text-[12px] text-gray-500 mt-1">
      音频时长：{{ Math.round(durationMs) }} ms
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import WaveSurfer from 'wavesurfer.js';
// wavesurfer.js v7 插件（ESM）
import RegionsPlugin from 'wavesurfer.js/dist/plugins/regions.esm.js';
import { api } from 'boot/axios';

type Region = {
  id: string;
  start: number;
  end: number;
  remove?: () => void;
  update?: (payload: { start: number; end: number }) => void;
};

type RegionsApi = {
  enableDragSelection?: (payload: { color: string }) => void;
  on?: (event: string, cb: (region: Region) => void) => void;
  getRegions?: () => unknown;
  addRegion?: (payload: {
    id: string;
    start: number;
    end: number;
    color: string;
    drag: boolean;
    resize: boolean;
  }) => void;
};

const props = withDefaults(
  defineProps<{
    sha256: string;
    fileType: string;
    startMs: number;
    endMs: number;
    height?: number;
  }>(),
  {
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

const audioUrl = computed(() => {
  if (!props.sha256 || !props.fileType) return '';
  const baseURL = api.defaults.baseURL || '';
  const sha256 = encodeURIComponent(props.sha256);
  const type = encodeURIComponent(props.fileType);
  return `${baseURL}/keytone_pkg/get_audio_stream?sha256=${sha256}&type=${type}`;
});

const hasSource = computed(() => !!audioUrl.value);

const isSyncingFromProps = ref(false);

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
      normalize: true,
      interact: true,
      plugins: [regions],
    });

    ws.value = instance;

    instance.on('error', () => {
      hasError.value = true;
    });

    instance.on('ready', () => {
      isReady.value = true;
      hasError.value = false;
      const dur = instance.getDuration();
      if (Number.isFinite(dur) && dur > 0) {
        durationMs.value = dur * 1000;
        emit('loaded', { durationMs: durationMs.value });
      }

      // 允许拖拽创建 region
      try {
        regions.enableDragSelection?.({
          color: 'rgba(14, 165, 233, 0.18)',
        });
      } catch {
        // ignore
      }

      // 如果外部已经有合法的裁剪区间，则创建/对齐 region
      syncRegionFromProps();

      // 监听 region 变化并回写到 props
      regions.on?.('region-created', (region) => {
        if (region.id !== regionId) {
          // 只保留一个 region
          try {
            const current = regions.getRegions?.();
            const list: Region[] = Array.isArray(current)
              ? (current as Region[])
              : current && typeof current === 'object'
              ? (Object.values(current as Record<string, unknown>) as Region[])
              : [];
            list.forEach((r) => {
              if (r.id !== regionId) r.remove?.();
            });
          } catch {
            // ignore
          }
          region.id = regionId;
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

  // start/end 不合法时，不主动创建 region（避免改变旧行为）
  if (!(endMs > startMs) || startMs < 0 || endMs < 0) return;

  const startSec = startMs / 1000;
  const endSec = endMs / 1000;

  isSyncingFromProps.value = true;
  try {
    const current = regionsApi.getRegions?.();
    const list: Region[] = Array.isArray(current)
      ? (current as Region[])
      : current && typeof current === 'object'
      ? (Object.values(current as Record<string, unknown>) as Region[])
      : [];
    const existing = list.find((r) => r.id === regionId);
    if (existing) {
      existing.update?.({ start: startSec, end: endSec });
    } else {
      regionsApi.addRegion?.({
        id: regionId,
        start: startSec,
        end: endSec,
        color: 'rgba(14, 165, 233, 0.18)',
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
  @apply overflow-hidden;
}
</style>
