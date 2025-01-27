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

  <q-page style="min-height: 0px" class="w-[379px] h-[458.5px]">
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

    <div :class="['w-[58%] flex flex-col ml-[21.8%] mr-[20.2%]']">
      <q-select
        v-model="setting_store.mainHome.selectedKeyTonePkg"
        :options="main_store.keyTonePkgOptions"
        :option-label="(item: any) => {
          return main_store.keyTonePkgOptionsName.get(item)
        }"
        :label="$t('mainHome.selectedKeySoundAlbum')"
        :virtual-scroll-slice-size="999999"
        outlined
        dense
        emit-value
        map-options
        behavior="dialog"
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
      </q-select>
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
      <div v-if="$t('mainHome.license') === 'License'" class="flex justify-center items-center">
        <!-- 关于许可协议的描述,
             * 没加the可在此链接参考https://www.kernel.org/doc/html/latest/process/license-rules.html#:~:text=Module%20is%20licensed%20under%20GPL%20version%202
             * 加the的可在此链接参考https://projects.blender.org/blender/blender?utm_medium=www-footer#:~:text=is%20licensed%20under%20the%20GNU%20General%20Public%20License%2C%20Version%203.
             * 使用released的链接参考https://www.blender.org/about/license/#:~:text=Blender%20is%20released%20under%20the%20GNU%20General%20Public%20License%20(GPL%2C%20or%20%E2%80%9Cfree%20software%E2%80%9D).
             总结,
             * 是否加the并不影响含义和语法。
             * 源码文档处声明尽量使用licensed, 软件下载网站页面和软件本身的ui界面中尽量使用released。
         -->
        KeyTone is released under the&nbsp;
        <div
          class="cursor-pointer hover:bg-gray-100"
          @click="openExternal('https://choosealicense.com/licenses/gpl-3.0/')"
        >
          GNU GPLv3
        </div>
      </div>
      <div v-else class="flex justify-center items-center">
        <!-- 此时相当于 v-if="$t('mainHome.license') === '许可协议'" -->
        <!-- 关于许可协议的描述, 可以在此链接参考https://www.kernel.org/doc/html/latest/translations/zh_CN/process/license-rules.html#:~:text=%E2%80%9CGPL%E2%80%9D-,%E6%A8%A1%E5%9D%97%E6%98%AF%E6%A0%B9%E6%8D%AEGPL%E7%89%88%E6%9C%AC2%E8%AE%B8%E5%8F%AF%E7%9A%84,-%E3%80%82%E8%BF%99%E5%B9%B6%E4%B8%8D%E8%A1%A8%E7%A4%BA -->
        <!-- 总之可在 '根据'  `基于`  `遵循` 等三个词语中选择 -->
        KeyTone 根据&nbsp;
        <div
          class="cursor-pointer hover:bg-gray-100"
          @click="openExternal('https://choosealicense.com/licenses/gpl-3.0/')"
        >
          GNU GPLv3
        </div>
        &nbsp;许可证发布
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
import { QSelect } from 'quasar';
import { useMainStore } from 'src/stores/main-store';
import { useSettingStore } from 'src/stores/setting-store';
import { computed, useTemplateRef, watch, onMounted, ref, onBeforeUnmount } from 'vue';

const setting_store = useSettingStore();
const main_store = useMainStore();

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
    // 更新保存的百分比
    main_store.volumePercentage = 1 - -setting_store.mainHome.audioVolumeProcessing.volumeNormal / min.value;

    // 原有的静音逻辑保持不变
    if (labelValue.value !== '0%') {
      setting_store.mainHome.audioVolumeProcessing.volumeSilent = false;
    } else {
      setting_store.mainHome.audioVolumeProcessing.volumeSilent = true;
    }
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
</script>

<style lang="scss" scoped></style>
