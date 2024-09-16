<!--
  1、 用户可以为即将诞生的新的音频包起一个名字。
  2、 要有音频文件的导入项, 以便于用户制作时传入音频文件。
  3、 要有音频文件的裁切功能, 以便于用户制作时精准地选择出一个键音。

-->

<template>
  <q-page>
    <q-scroll-area class="w-[379px] h-[458.5px]">
      <div :class="['flex flex-col gap-5  p-8 ']">
        <q-input outlined stack-label dense v-model="text" label="键音包名称" />
        <!-- <div>原始声音文件编辑</div>
        <div>键音</div>
        <div>键音列表, 编辑键音</div>
        <div>全局键音规则</div>
        <div>对某个特定按键单独设置键音</div> -->

        <q-stepper v-model="step" vertical header-nav color="primary" animated>
          <div :class="['text-center font-semibold text-lg text-nowrap']">{{ text }}</div>
          <q-step :name="1" title="原始声音文件" icon="settings" :done="step > 1">
            <div>为此键音包载入原始的音频文件供后续步骤使用, 文件类型可以是WAV、MP3、OGG等。</div>
            <div>原始音频文件的数量不定,可根据您的制作喜好来决定。</div>

            <q-card :class="['p-2']"> qq</q-card>
            <q-stepper-navigation>
              <q-btn @click="step = 2" color="primary" label="Continue" />
            </q-stepper-navigation>
          </q-step>

          <!-- <q-step :name="2" title="键音制作" caption="Optional" icon="create_new_folder" :done="step > 2"> -->
          <q-step :name="2" title="键音制作" icon="create_new_folder" :done="step > 2">
            <div>根据原始音频文件制作键音。</div>
            <!-- <div>
              若您载入的原始音频文件本身就是一个独立完善的键音, 则性能更好。换言之, 原始键音文件越接近最终键音性能越好。
            </div> -->
            <!-- <div>
              若您载入的原始音频文件本身就是一个独立完善的键音, 则此步骤可跳过(因为其会在第一步时,
              被自动添加到键音列表中)。
            </div> -->
            <q-stepper-navigation>
              <q-btn @click="step = 3" color="primary" label="Continue" />
              <q-btn flat @click="step = 1" color="primary" label="Back" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>

          <q-step :name="3" title="键音包全局设置" icon="add_comment" :done="step > 3">
            <div>现在有了一定数量的键音, 开始在全局范围内配置它们的用途吧!</div>
            <q-stepper-navigation>
              <q-btn @click="step = 4" color="primary" label="Continue" />
              <q-btn flat @click="step = 2" color="primary" label="Back" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>
          <q-step :name="4" title="对某个特定按键单独设置键音" caption="本步骤可选(非必填)" icon="add_comment">
            <div>有些情况下, 我们希望独立定义某个按键的键音, 而不是全局随机。</div>
            <div>甚至, 更极端的情况, 我们希望键盘上所有的按键, 都拥有自己的独立键音。</div>
            <div>来吧!这一定制步骤将满足您的需求!!!</div>
            <q-stepper-navigation>
              <q-btn flat @click="step = 3" color="primary" label="Back" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>
        </q-stepper>
      </div>
    </q-scroll-area>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const text = ref('');
const step = ref(1);
</script>

<style lang="scss" scoped></style>
