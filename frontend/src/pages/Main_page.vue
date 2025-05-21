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
  <!-- <div>Silent Typist's Friend</div> -->
  <!--
   q-layout 有个 style="min-height: 803.2px" 的样式, 会造成滚动条的出现 进而无法展示窗口底部 进而无法通过css实现圆角
   因此, 我们这里主动设置 style="min-height: 0px"
  -->

  <q-page style="min-height: 0px" :class="[isMacOS ? 'w-[389.5px] h-[458.5px]' : 'w-[379px] h-[458.5px]']">
    <div
      :class="[
        '',
        'w-auto h-50 flex flex-col items-center',
        // 使页面文本无法选择, 免得影响界面体验。
        'select-none',
      ]"
    >
      <q-avatar
        :class="[
          '',
          // 设置整体头像的尺寸大小。 至于内部相关文字的尺寸大小(如果有的话  可通过设置text-[8rem] leading-[10rem]实现)
          'size-56',
          '',
        ]"
      >
        <!-- draggable="false"使得图片无法拖动, 免得影响界面体验 -->
        <img :src="logoUrl" draggable="false" />
      </q-avatar>
    </div>

    <div :class="['flex flex-col']">
      <div class="flex flex-col items-center">
        <q-select
          :class="['w-[216px]', 'select-component-label-show']"
          v-model="setting_store.mainHome.selectedKeyTonePkg"
          :options="main_store.keyTonePkgOptions"
          :option-label="(item: any) => main_store.keyTonePkgOptionsName.get(item)"
          :label="$t('mainHome.selectedKeySoundAlbum')"
          :virtual-scroll-slice-size="999999"
          outlined
          dense
          emit-value
          map-options
          behavior="dialog"
          popup-content-class="w-[100%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
          ref="selectedKeyTonePkgRef"
          @popup-hide="blur()"
        >
          <template v-if="setting_store.mainHome.selectedKeyTonePkg" v-slot:append>
            <!-- 由于直接使用默认的clearable, 会使得mode=null, 而我希望点击清楚按钮时mode=""即空字符串。因此使用插槽来实现。 -->
            <q-icon
              name="cancel"
              @click.stop.prevent="setting_store.mainHome.selectedKeyTonePkg = ''"
              class="cursor-pointer text-lg"
            />
          </template>

          <!-- 空状态提示 -->
          <template v-slot:no-option>
            <div class="flex flex-col items-center py-4 text-gray-500">
              <q-icon name="library_music" size="40px" class="mb-2 opacity-50" />
              <div class="text-sm mb-3">{{ $t('mainHome.emptyState.noAlbum') }}</div>
              <q-btn
                flat
                dense
                class="empty-state-btn flex items-center bg-blue-500/10 px-4 py-1.5 rounded-lg"
                @click="goToAlbumPage"
              >
                <q-icon name="keyboard_arrow_right" size="20px" class="mr-1" />
                <span class="text-sm">{{ $t('mainHome.emptyState.goToAlbumPage') }}</span>
              </q-btn>
            </div>
          </template>
        </q-select>

        <!-- 空状态额外提示 -->
        <transition name="fade">
          <div v-if="!main_store.keyTonePkgOptions.length" class="text-center mt-2">
            <div class="text-xs text-gray-400 -ml-0.6">{{ $t('mainHome.emptyState.createOrImportTip') }}</div>
          </div>
        </transition>
      </div>
    </div>
    <!--
      * TIPS: 对于主页的音量调整, 理念是最终原始音频, 即最大音量等于原始音频的正常音量(或是最大音量), 也就是说保持为0。
              > 在此基础上, 其最小值始终保持, 可调整至正常音量的对数 g <目前默认g为5, 写死, 后续如果有必要, 会在设置界面对其新增相关的设置项>。
              > > 当我们通过设置界面的音量强来增大原始音频音量时, 其音量的原始正常值将不再为0, 而是变为对音量增强的指数取相同基数后的对数值。
              > > * 即实际上音频的原始音量对应的数值是 `-setting_store.audioVolumeProcessing.volumeAmplify`
              > > 因此, 此时我们的最小值, 要想保持正常音量的对数 g , 需要 使用`-setting_store.audioVolumeProcessing.volumeAmplify-5`这个数值
              > >   - completed(已完成)   FIXME:
              > > * 当setting_store.audioVolumeProcessing.volumeAmplify > 0 时, 当然没问题不用过分讨论。
              > > * 当setting_store.audioVolumeProcessing.volumeAmplify < 0 且  > -5 时, 使用这个算法就发生了问题(或者说永远保持正常音量的对数g的方案发生了问题)
              > >   * 此时的缩小是仍旧在以缩小的基础上进行的(即此时的0), 而此时仍保证真实正常音量的对数g的话, 是不现实的, 因为主页的音量调整是一个只能缩小音量的滑块, 这要保持只能缩短滑块范围。
              > >   * 拿-2来举例子, 貌似得出的最小值-3是更靠近正常音量了, 但实际上得出的这个最小值在处理时, 并不是一个更靠近正常音量的值(因为 -2 这个值, 在Amplify时已经被实际处理过了, 而这里在添一个-3, 仍旧是-5)
              > >     * 或者说, 此时的0对应的已经是缩小后的音量了(而不是像正数时放大的音量), 我们无法通过只能缩减音量的滑块获得正常的音量, 此时要想进一步保持对数5, 就只能缩短滑块进度条的实际范围了。
              > > * 当setting_store.audioVolumeProcessing.volumeAmplify < -5时, 问题就更大了。
              > >   * 由于要保持真实正常音量的对数g, 我们不得不缩短进度条, 但当这个 setting_store.audioVolumeProcessing.volumeAmplify = -5 时, 这个进度条实际已经不存在了, 此时只能是正常音量的对数g。
              > >   * 而当其进一步缩小的<-5时, 不但最小音量无法保持正常音量的对数g, 甚至滑块无法表示目前真实缩小的正常音量的对数值。(算成百分比时, 甚至出现了大于100%的离谱情况, 实际上这些只是更小的负的百分数)
              > >  FIXME的完成小记: 放弃最小音量值, 始终要想保持正常音量的对数 g的方案。因为当 volumeAmplify 小于0 后, 这个方案是无意义的
              > >                  * 此时, 我们最多只能保证其是 当前已被缩小音量<即此时的0对应的音量>的对数g。或者说此时可缩小的范围固定为g<因为我们的最大值固定为0这个数>。
              > >                    * 即 当 volumeAmplify 小于0 后 的最新算法是 -5, 或者说-g
     -->
    <div
      :class="[
        '',
        'w-full flex flex-col items-center pr-5 mt-10',
        // 使页面文本无法选择, 免得影响界面体验。
        'select-none',
      ]"
    >
      <div :class="['w-56 flex justify-between items-center']" v-show="isInitialized">
        <q-btn
          dense
          round
          flat
          :icon="setting_store.mainHome.audioVolumeProcessing.volumeSilent ? 'volume_off' : 'volume_up'"
          @click="isSilent"
        >
        </q-btn>

        <q-slider
          :class="['w-[80%]']"
          v-model="setting_store.mainHome.audioVolumeProcessing.volumeNormal"
          :max="0"
          :min="-min"
          :step="0"
          label
          :label-value="labelValue"
          color="light-green"
        />
      </div>
      <div :class="['w-56  flex justify-end items-center mt-5']">
        <q-slider
          :class="['w-[80%]']"
          v-if="setting_store.mainHome.audioVolumeProcessing.isOpenVolumeDebugSlider"
          v-model="setting_store.mainHome.audioVolumeProcessing.volumeNormal"
          :max="0"
          :min="-min"
          :step="0"
          :markers="markersDebug"
          marker-labels
          label
          label-always
          :label-value="labelValueDebug"
          color="light-green"
        />
      </div>
    </div>
    <div
      :class="[
        'text-center text-xs text-gray-500 absolute bottom-8.5 w-full',
        // 使页面文本无法选择, 免得影响界面体验。
        'select-none',
      ]"
    >
      <div class="flex justify-center items-center m-r-0.3">
        {{ $t('mainHome.copyright') }}{{ new Date().getFullYear() }}&nbsp;
        <div class="cursor-pointer hover:bg-gray-100" @click="openExternal('https://github.com/LuSrackhall')">
          LuSrackhall
        </div>
      </div>
      <div class="flex justify-center items-center">
        <!-- Using i18n component with license slot interpolation -->
        <i18n-t keypath="mainHome.licenseText" tag="div" class="flex justify-center items-center">
          <template #license>
            &nbsp;
            <div
              class="cursor-pointer hover:bg-gray-100"
              @click="openExternal('https://choosealicense.com/licenses/gpl-3.0/')"
            >
              GNU GPLv3
            </div>
            &nbsp;
          </template>
        </i18n-t>
      </div>
    </div>
    <div
      :class="[
        'text-center text-2.8 text-gray-500 absolute bottom-2.8 w-full',
        // 使页面文本无法选择, 免得影响界面体验。
        'select-none',
      ]"
    >
      <div class="flex justify-center items-center mb-1.3 m-r-0.6">
        <div
          class="cursor-pointer hover:bg-gray-100"
          @click="openExternal('https://keytone.xuanhall.com/guide/other/privacy-policy/')"
        >
          {{ $t('mainHome.privacyPolicy') }}
        </div>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
        <div
          class="cursor-pointer hover:bg-gray-100"
          @click="openExternal('https://keytone.xuanhall.com/guide/other/user-agreement/')"
        >
          {{ $t('mainHome.userAgreement') }}
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import logoUrl from 'assets/img/KeyTone.png?url';
import { QSelect, useQuasar } from 'quasar';
import { useMainStore } from 'src/stores/main-store';
import { useSettingStore } from 'src/stores/setting-store';
import { computed, useTemplateRef, watch, onMounted, ref, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

const q = useQuasar();
const router = useRouter();
const { t } = useI18n();
const $t = t;
const setting_store = useSettingStore();
const main_store = useMainStore();

// 导航到键音专辑页面
const goToAlbumPage = () => {
  router.push('/keytone_album');
};

// 添加初始化状态控制(用于控制音量百分比初始化完成后再显示, 防止增减音量调整后回到主页面的一瞬间, 音量百分比为保持历史状态而调整的过程的显示)
const isInitialized = ref(false);

// 修改onMounted钩子
onMounted(async () => {
  // 初始化音量百分比
  await main_store.initVolumePercentage();

  // 使用保存的百分比计算初始volumeNormal
  const currentAmplify = setting_store.audioVolumeProcessing.volumeAmplify;
  const newMin =
    currentAmplify > 0
      ? currentAmplify + setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope
      : setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope;

  if (main_store.volumeNormalReduceScope === setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope) {
    setting_store.mainHome.audioVolumeProcessing.volumeNormal = -(newMin * (1 - main_store.volumePercentage));
  } else {
    // 记录百分比(由于此时min变化引起了百分比变化, 因此也需要主动记录, 以更新保存的百分比)
    // TIPS: (准确的说是min所依赖的volumeNormalReduceScope变化引起的,强调这一点是因为, 如果是volumeAmplify变化引起的min及百分比变化, 我们千万不用记录, 因为我们本质需求要保留的就是这个百分比)
    // TIPS: 由于min是计算属性, 只要我在此处调用, 他就一定会计算后在返回给我, 因此不用担心min的值不是最新的问题, 这里是安全的。
    main_store.volumePercentage = 1 - -setting_store.mainHome.audioVolumeProcessing.volumeNormal / min.value;
  }

  // 初始化完成后再显示
  isInitialized.value = true;
});

// 修改volumeAmplify的监听
watch(
  () => setting_store.audioVolumeProcessing.volumeAmplify,
  (newAmplify, oldAmplify) => {
    if (oldAmplify === undefined) return; // 忽略初始化时的变化，由onMounted处理

    // 使用旧的amplify计算旧的min值
    const oldMin =
      oldAmplify > 0
        ? oldAmplify + setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope
        : setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope;

    // 计算并保存当前的百分比
    const currentPercentage = 1 - -setting_store.mainHome.audioVolumeProcessing.volumeNormal / oldMin;
    main_store.volumePercentage = currentPercentage;

    // 使用新的amplify计算新的min值
    const newMin =
      newAmplify > 0
        ? newAmplify + setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope
        : setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope;

    // 使用新的min值和保存的百分比计算新的volumeNormal值
    setting_store.mainHome.audioVolumeProcessing.volumeNormal = -(newMin * (1 - currentPercentage));
  }
);

onBeforeUnmount(() => {
  // 在组件卸载时执行的逻辑
  // 卸载前记录volumeNormalReduceScope的值, 方便对比
  main_store.volumeNormalReduceScope = setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope;
});

const min = computed(() => {
  if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
    return (
      setting_store.audioVolumeProcessing.volumeAmplify +
      setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope
    );
  } else {
    //  FIXME的完成小记: 对应步骤小记, 当setting_store.audioVolumeProcessing.volumeAmplify<0时, 我们使用固定的g值作为缩小幅度(或者说已被缩小的值<即0>为参考的对数g, 因为我们的最大值固定为0, 因此也是固定的g值)
    //                   * g值 = setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope。
    return setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope;
  }
});

const labelValue = computed(() => {
  const percentage = ((1 - -setting_store.mainHome.audioVolumeProcessing.volumeNormal / min.value) * 100)
    .toFixed(2)
    .split('.');
  return percentage[1] === '00' ? percentage[0] + '%' : percentage[0] + '.' + percentage[1] + '%';
});

// 监听volumeNormal的变化，更新保存的百分比
watch(
  () => setting_store.mainHome.audioVolumeProcessing.volumeNormal,
  () => {
    // TIPS: 在引发多个相关变量的写入持久化的操作时, 具体问题具体分析, 通过在相关源头处对次要变量进行手动延后来解决新的变更被错误改回旧变更的问题。
    setTimeout(() => {
      // 更新保存的百分比
      main_store.volumePercentage = 1 - -setting_store.mainHome.audioVolumeProcessing.volumeNormal / min.value;
    }, 30);

    // TIPS: 在引发多个相关变量的写入持久化的操作时, 具体问题具体分析, 通过在相关源头处对次要变量进行手动延后来解决新的变更被错误改回旧变更的问题。
    //       * 对每个涉及到的变量, 都单独进行手动延后, 并且延后的时间不同(保持30ms以上的差距), 以确保其不会因同时触发而造成问题->即同时触发造成的sse的返回数据不可靠的问题->造成所涉及的相关的一切数据项(包括主要变量和次要变量)最终一致性变的随机和不可靠->而这些同时沟通底端与ui端渲染的重要数据变量的不可靠是非常影响用户体验的。
    setTimeout(() => {
      // 原有的静音逻辑保持不变
      if (labelValue.value !== '0%') {
        setting_store.mainHome.audioVolumeProcessing.volumeSilent = false;
      } else {
        setting_store.mainHome.audioVolumeProcessing.volumeSilent = true;
      }
    }, 60);
  }
);

// 监听volumeNormalReduceScope的变化, 更新保存的百分比(确保实时刷新时(实际上实时刷新的场景仅存在与代码调试过程)不会出现百分比变化不符合预期的问题)
watch(
  () => setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope,
  () => {
    // TIPS: 在引发多个相关变量的写入持久化的操作时, 具体问题具体分析, 通过在相关源头处对次要变量进行手动延后来解决新的变更被错误改回旧变更的问题。
    setTimeout(() => {
      // 更新保存的百分比
      main_store.volumePercentage = 1 - -setting_store.mainHome.audioVolumeProcessing.volumeNormal / min.value;
    }, 30);
  }
);

watch(
  min,
  () => {
    if (-setting_store.mainHome.audioVolumeProcessing.volumeNormal > min.value) {
      setting_store.mainHome.audioVolumeProcessing.volumeNormal = -min.value;
    }
  },
  // TIPS: 必须添加立即执行, 因为min是基于setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope的计算值,
  //      而后者是用户设置的值, 且是在其它界面设置的(此时Main_page是未加载的), 因此需要回到Main_page界面后立即执行以确保正确性
  { immediate: true }
);

const labelValueDebug = computed(() => {
  return setting_store.mainHome.audioVolumeProcessing.volumeNormal.toFixed(2);
});
const markersDebug = computed(() => {
  if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
    return (
      (setting_store.audioVolumeProcessing.volumeAmplify +
        setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope) /
      1
    );
  } else {
    return setting_store.mainHome.audioVolumeProcessing.volumeNormalReduceScope / 1;
  }
});

const isSilent = (event: any) => {
  if (event.detail === 0) {
    // 由键盘触发，不处理
    return;
  }
  if (labelValue.value === '0%') {
    // 由于音量为0时，打开声音是没意义的，所以直接返回。(了解原理的话不难理解, 此时打开声音后可能造成0%的音量调状态, 仍会发出声音的bug)
    q.notify({
      message: $t('Notify.音量0%时无法打开声音'),
      color: 'warning',
      position: 'top',
      timeout: 1200,
    });
    // 这种操作的优点是, 可以让用户通过音量图标的状态, 直观感受到音量是否为0, 可以可靠的知道此时KeyTone当前一定是没有声音的。
    // TODO: 当然, 这并不是最好的解决方式, 另一种解决方式是, 对0%的音量状态, 单独设置一个静音变量, 以确保在音量为0时, 仍然可以打开声音。但打开声音也不用担心0%状态会发出声音, 因单独的静音变量作用, 将仍不会有声音, 直到不为0%后声音才能恢复。
    // TIPS: 这种操作的唯一缺点是(也不能算作缺点, 因为当图标关闭时, 仍能确定其当前是一定没有声音的), 某些情况下无法通过音量图标的状态, 直观感受到音量是否为0。(因为音量为0时, 图标也有可能是开启状态)

    return;
  }
  setting_store.mainHome.audioVolumeProcessing.volumeSilent =
    !setting_store.mainHome.audioVolumeProcessing.volumeSilent;
};

// 每次用户的主动选择, 都会触发实际选择的键音包重新进行加载。
watch(
  () => setting_store.mainHome.selectedKeyTonePkg,
  () => {
    // 即使多次调用此函数也无妨, 相同uuid的键音包动作, 不会影响已加载并使用中的键音包。
    // TIPS: 在sse回调中也再次调用此函数, 以保证用户的选择能够最大程度上被可靠的加载, 并且无需担心, 重复调用此函数也不会引发重复加载相同的键音包
    main_store.LoadSelectedKeyTonePkg();
  }
  // TIPS: 放弃下方立即执行的相关代码, 而是将相关逻辑移动到它该有的地方, 以保证逻辑的清晰性。(也就是App.vue中, 或是boot中。)
  //       * 对于 从 创建/编辑键音包界面返回的场景也无需担心。 只需在对应界面稀释逻辑中, 触发加载用户所选的键音包的逻辑即可。(这样可以简化逻辑, 也可以保证逻辑的清晰性)
  // ~~立即执行, 使得每次进入主界面时, 都会加载用户所选的键音包。(主要是软件启动时, 适配加载用户所选的键音包)~~
  // ~~{ immediate: true }~~
);

const selectedKeyTonePkgRef = useTemplateRef<QSelect>('selectedKeyTonePkgRef');

const blur = () => {
  setTimeout(() => {
    selectedKeyTonePkgRef?.value?.blur();
    // TIPS: 这里需要延迟后再blur, 以确保blur的正确触发(太早触发blur会不起作用, 经验证, 本人电脑延迟10ms后, 可以正确触发blur使焦点丧失, 为确保适配更多的低性能设备, 这里保险起见设置为66ms)
  }, 66);
};

function openExternal(url: string) {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.openExternal(url);
  }
}

const isMacOS = ref(getMacOSStatus());
function getMacOSStatus() {
  if (process.env.MODE === 'electron') {
    return window.myWindowAPI.getMacOSStatus();
  }
  return false;
}
</script>

<style lang="scss" scoped>
// 淡入淡出动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// 空状态相关样式
.empty-state-btn {
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-1px);
    background-color: rgba(59, 130, 246, 0.15);
  }

  &:active {
    transform: translateY(0);
  }

  &::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 100%;
    height: 100%;
    background: radial-gradient(circle, rgba(255, 255, 255, 0.2) 0%, transparent 70%);
    transform: translate(-50%, -50%) scale(0);
    opacity: 0;
    transition: transform 0.4s ease, opacity 0.3s ease;
  }

  &:hover::before {
    transform: translate(-50%, -50%) scale(2);
    opacity: 1;
  }
}

// 优化下拉选择器的空状态样式
:deep(.q-select) {
  .q-field__control {
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(8px);
    transition: all 0.3s ease;
  }
}

// 实际上, 在此处通过global来更改的.q-field__native, 已经覆盖了:deep(.q-field__native)的样式, 因此上方(包括其它文件中的):deep(.q-field__native)的样式可以删除(如果与这里相同的话)。
// * global(或者说不带scoped的style)的影响范围是全局的, 包括其它vue文件中的内容也将会受此影响。
//   > 不过这种影响是有前提的, 机这个带有global的组件必须至少加载过一次。 比如a.vue和b.vue为main.vue下的同级别组件, a中拥有global样式, b中没有。
//   > * 若 main.vue -> b.vue 的话, b.vue将不会受到a.vue的global样式的影响。
//   > * 若 main.vue -> a.vue -> b.vue 的话, b.vue将会受到a.vue的global样式的影响。
//   > * 若 main.vue -> a.vue -> b.vue 的话, b.vue将会受到a.vue的global样式的影响, 但此时若在b.vue中进行刷新操作, 则由于刷新后重载的过程中没有加载过a.vue, 因此b.vue将不会受到a.vue的global样式的影响。
// * 前面提到若是'当前(包括其它)'组件的scoped中, 拥有与 global(或者说不带scoped的style)相同的样式, 则可以省略 '当前(包括其它)'组件的scoped的这些样式。
//   > 但重要的是, 这种省略是有条件的, 即这些样式必须与global(或者说不带scoped的style)中的样式相同, 否则 scoped 中的样式将会覆盖 global 的样式。
// TIPS: 开头的叙述存在不当, 若您的global样式不再根组件中的话, 推荐您就算样式都相同, 也不要省略 '其它组件' 的scoped中的这些样式。 尽量在每个文件中都来一份, 避免这些样式因刷新或是其它情况丢失, 从而造成的页面故障。
// 此处使用:global(.q-field__native)的原因是:deep(.q-field__native)无法覆盖当前页面中quasar的选择器的菜单, 以对话框展开时对话框上方的已选择框中的对应样式。
// * 但这样做的弊端, 是有可能影响到其它文件。 因此, 如果发现其它文件受影响, 则需要最组件对应的scoped中, 使用相同的类型, 来恢复global造成的影响。
:global(.q-field__native) {
  // 对溢出的情况, 采取滚动策略
  @apply max-w-full overflow-auto whitespace-nowrap;

  // 隐藏滚动策略的滚动条。
  // @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];

  // 添加细微滚动条
  @apply h-5.8 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400;
}

:global(.q-item__section) {
  /* 对溢出的情况, 采取滚动策略 */
  @apply max-w-full overflow-auto whitespace-nowrap;

  /* 隐藏滚动策略的滚动条 */
  // @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];

  // 添加细微滚动条
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30  [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}

// 对于主页键音专辑的选择框, 键音专辑的名称内容过长的情况, 采取溢出滚动的策略。
:deep(.ellipsis) {
  // 对溢出的情况, 采取滚动策略
  @apply max-w-full overflow-auto whitespace-nowrap  text-clip;
  // // 隐藏滚动策略的滚动条。
  // @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];

  // 添加细微滚动条
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30  [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}

:deep(.q-field__label) {
  @apply overflow-visible;
}
</style>

<!-- 可以通过在 含有 scoped 的style中, 通过:global 来使用, 效果与下方相同 -->
<!-- <style lang="scss">
.q-field__native {
  // 对溢出的情况, 采取滚动策略
  @apply max-w-full overflow-auto whitespace-nowrap;
  // 隐藏滚动策略的滚动条。
  @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];
}
</style> -->

<!--
TIPS: :global 和 :deep 的作用范围不同：
      * :global 作用于全局，会跳过 Vue 的 scoped CSS 限制
      * :deep 只能作用于当前组件及其子组件的 DOM 树内(即需要受到当前组件的 scoped CSS 限制)
      所以在 :global 中使用 :deep 是没有意义的，因为 :global 已经跳过了作用域限制。
-->
