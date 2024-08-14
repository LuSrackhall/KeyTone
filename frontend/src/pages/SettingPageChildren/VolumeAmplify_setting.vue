<template>
  <!-- <q-item>
    <div :class="['ml-0 rounded-full  mr-6 h-6 self-center']"></div>
    <div :class="['w-full flex justify-between items-center']">
      <q-badge :class="['bg-cyan-700 h-5']"> Volume: {{ setting_store.audioVolumeProcessing.volumeAmplify }} </q-badge>

      <q-btn :class="['w-15 h-5 mr-5']" color="primary" size="10px" label="重置" @click="returnToNormal()" />
    </div>
  </q-item> -->

  <q-item :class="['pt-8 h-50']">
    <!-- select左边的竖线 -->
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-[80%] self-center']"></div>

    <div :class="['w-[100%] grid justify-items-center']">
      <div :class="['w-[90%] flex justify-between items-center ml-3 mb-10']">
        <q-badge :class="['bg-cyan-700 h-5 ml-0']">
          {{ $t('setting.原始音量增减调节.音量') }} : {{ setting_store.audioVolumeProcessing.volumeAmplify }}
        </q-badge>
        <q-btn
          :class="['w-15 h-5 ml-0 mr-2.3']"
          color="primary"
          size="10px"
          :label="$t('setting.原始音量增减调节.重置')"
          @click="returnToNormal()"
        />
      </div>

      <div :class="['w-[100%] grid justify-items-center']">
        <q-slider
          :class="['w-[80%]']"
          v-model="setting_store.audioVolumeProcessing.volumeAmplify"
          :max="max"
          :min="min"
          :step="0"
          :markers="markers"
          marker-labels
          label
          label-always
          :label-value="labelValue"
          color="light-green"
        />
      </div>
    </div>
  </q-item>
  <q-item :class="['h-15 mb-6']">
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-[80%] self-center']"></div>

    <div :class="['w-[100%] grid justify-items-center']">
      <div :class="['w-[88%] flex justify-between items-center']">
        <q-input
          dense
          hide-bottom-space
          :class="['w-[50%] h-10.5 ']"
          v-model.number="setting_store.audioVolumeProcessing.volumeAmplifyLimit"
          type="number"
          filled
          :rules="[(val: number) => { return val > 0 && val<100000000 || '请输入一个大于0且小于100000000的数字'; }]"
        />

        <q-btn
          :class="['w-15 h-5']"
          color="primary"
          size="10px"
          :label="$t('setting.原始音量增减调节.重置')"
          @click="returnToNormalLimit()"
        />
      </div>
    </div>
  </q-item>
</template>

<script setup lang="ts">
import { useSettingStore } from 'src/stores/setting-store';
import { ref, watch } from 'vue';
import { debounce } from 'lodash';

const setting_store = useSettingStore();

const max = ref(setting_store.audioVolumeProcessing.volumeAmplifyLimit);
const min = ref(setting_store.audioVolumeProcessing.volumeAmplifyLimit * -1);
const markers = ref(setting_store.audioVolumeProcessing.volumeAmplifyLimit / 1);
const labelValue = ref(setting_store.audioVolumeProcessing.volumeAmplify.toFixed(2));

const debounced = debounce(
  () => {
    if (
      setting_store.audioVolumeProcessing.volumeAmplifyLimit <
      Math.abs(setting_store.audioVolumeProcessing.volumeAmplify)
    ) {
      if (setting_store.audioVolumeProcessing.volumeAmplify > 0) {
        setting_store.audioVolumeProcessing.volumeAmplify = setting_store.audioVolumeProcessing.volumeAmplifyLimit;
      } else {
        setting_store.audioVolumeProcessing.volumeAmplify = setting_store.audioVolumeProcessing.volumeAmplifyLimit * -1;
      }
    }
  },
  800,
  { trailing: true }
);

watch(
  () => setting_store.audioVolumeProcessing.volumeAmplifyLimit,
  () => {
    if (
      setting_store.audioVolumeProcessing.volumeAmplifyLimit > 0 &&
      setting_store.audioVolumeProcessing.volumeAmplifyLimit < 100000000
    ) {
      markers.value = setting_store.audioVolumeProcessing.volumeAmplifyLimit / 1;
      max.value = setting_store.audioVolumeProcessing.volumeAmplifyLimit;
      min.value = setting_store.audioVolumeProcessing.volumeAmplifyLimit * -1;

      debounced.cancel;
      debounced();
    }
  }
);

// TIPS: 这里的watch处理, 是解决控件的一个bug。(涉及到多个控件的结合使用, 错误原因不明, 不过这样解藕后, 便可解决)
//       > 输入框控件中持续执行`backspace按键`的删除操作时,会引发的`q-slider 控件`的 `:label-value="labelValue"`报错。
watch(
  () => setting_store.audioVolumeProcessing.volumeAmplify,
  () => {
    if (
      setting_store.audioVolumeProcessing.volumeAmplifyLimit > 0 &&
      setting_store.audioVolumeProcessing.volumeAmplifyLimit < 100000000
    ) {
      labelValue.value = setting_store.audioVolumeProcessing.volumeAmplify.toFixed(2);
    } else {
      // 如果setting_store.audioVolumeProcessing.volumeAmplifyLimit不符合规范, 则在用户操作音量条时将最后一次符合规范的值, 重新给到它
      setting_store.audioVolumeProcessing.volumeAmplifyLimit = max.value;
    }
  }
);

function returnToNormal() {
  setting_store.audioVolumeProcessing.volumeAmplify = 0.0;
}

function returnToNormalLimit() {
  setting_store.audioVolumeProcessing.volumeAmplifyLimit = 10.0;
}
</script>

<style lang="scss" scoped></style>
