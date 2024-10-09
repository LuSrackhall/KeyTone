<!--
  1、 用户可以为即将诞生的新的音频包起一个名字。
  2、 要有音频文件的导入项, 以便于用户制作时传入音频文件。
  3、 要有音频文件的裁切功能, 以便于用户制作时精准地选择出一个键音。

-->

<template>
  <q-page>
    <q-scroll-area class="w-[379px] h-[458.5px]">
      <div :class="['flex flex-col gap-5  p-8 ']">
        <q-input
          outlined
          stack-label
          dense
          :error-message="$t('KeyTonePackage.new.name.errorMessage')"
          :error="pkgName === '' || pkgName === undefined || pkgName === null"
          v-model="pkgName"
          :label="$t('KeyTonePackage.new.name.name')"
          :placeholder="$t('KeyTonePackage.new.name.defaultValue')"
        />
        <!-- <div>原始声音文件编辑</div>
        <div>键音</div>
        <div>键音列表, 编辑键音</div>
        <div>全局键音规则</div>
        <div>对某个特定按键单独设置键音</div> -->

        <q-stepper v-model="step" vertical header-nav color="primary" animated>
          <div :class="['text-center font-semibold text-lg text-nowrap']">{{ pkgName }}</div>
          <q-step :name="1" title="载入声音文件" icon="create_new_folder" :done="step > 1">
            <div>为此键音包载入原始的声音文件供后续步骤使用。</div>
            <!-- <div>文件类型可以是WAV、MP3、OGG等。</div> -->
            <!-- <div>原始音频文件的数量不定,可根据您的制作喜好来决定。</div> -->
            <!-- <q-card class="bg-slate-500" :class="['p-2']"> -->
            <!-- <q-btn :class="['bg-zinc-300']" label="添加新的声音源文件"></q-btn>
            <div :class="['p-2 text-zinc-300']">或</div>
            <q-btn :class="['bg-zinc-300']" label="选择声音源文件以进行编辑"></q-btn> -->
            <!-- </q-card> -->
            <!-- ------------------------------------------------------------------------载入声音文件的业务逻辑 start -->
            <div>
              <!-- ------------------------------------------------------------------------------ 添加新的声音源文件 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  label="添加新的声音源文件"
                  @click="
                    () => {
                      addNewSoundFile = !addNewSoundFile;
                    }
                  "
                >
                </q-btn>
                <q-dialog v-model="addNewSoundFile" backdrop-filter="invert(70%)">
                  <q-card>
                    <q-card-section class="row items-center q-pb-none text-h6"> 添加新的声音源文件 </q-card-section>

                    <q-card-section> <div>文件类型可以是WAV、MP3、OGG。</div></q-card-section>

                    <q-card-section>
                      <q-file
                        :class="['w-56']"
                        dense
                        v-model="files"
                        label="Pick files"
                        outlined
                        use-chips
                        multiple
                        append
                        accept=".wav,.mp3,.ogg"
                        excludeAcceptAllOption
                        style="max-width: 300px"
                      />
                    </q-card-section>

                    <q-card-section>
                      <div>数量不定, 跟随制作喜好添加即可</div>
                    </q-card-section>
                    <q-card-actions align="right">
                      <q-btn
                        flat
                        @click="
                          async () => {
                            // 循环files, 并在每次上传成功后, 删除对应file

                            if (!files || files.length === 0) {
                              console.warn('No files selected for upload');
                              return;
                            }

                            // 使用slice()方法创建一个数组的浅拷贝, 避免因遍历过程中修改原始数组而导致的遍历中止
                            // slice也可以只截取数组的一部分, 类似golang的切片, 都是左闭右开区间。 如slice(2,4) 会从[1,2,3,4,5]中, 截取[3,4]
                            for (const file of files.slice()) {
                              try {
                                const re = await SendFileToServer(pkgPath, file);
                                if (re === true) {
                                  console.info(`File ${file.name} uploaded successfully`);
                                  // Remove the file from the list after successful upload
                                  const index = files.indexOf(file);
                                  if (index > -1) {
                                    files.splice(index, 1);
                                  }
                                } else {
                                  console.error(`File ${file.name} uploading error`);
                                  q.notify({
                                    type: 'negative',

                                    position: 'top',

                                    message: `File '${file.name}' addition failed`,

                                    timeout: 5,
                                  });
                                }
                              } catch (error) {
                                console.error(`Error uploading file ${file.name}:`, error);
                                q.notify({
                                  type: 'negative',

                                  position: 'top',

                                  message: `File '${file.name}' addition failed`,

                                  timeout: 5,
                                });
                              }
                            }
                          }
                        "
                        color="primary"
                        label="确认添加"
                      />
                      <q-btn flat label="Close" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
              <div :class="['p-2 text-zinc-600']">或</div>
              <!-- -------------------------------------------------------------------------------编辑已有声音源文件 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  label="编辑已有声音源文件"
                  @click="
                    () => {
                      editSoundFile = !editSoundFile;
                    }
                  "
                ></q-btn>
                <q-dialog v-model="editSoundFile" backdrop-filter="invert(70%)">
                  <q-card>
                    <q-card-section class="row items-center q-pb-none text-h6"> 编辑已有声音源文件 </q-card-section>

                    <q-card-section> <div>请选择您想要修改或删除的声音源文件并执行对应操作。</div></q-card-section>

                    <q-card-section>
                      <!-- <q-select outlined v-model="model" :options="options" label="Outlined" /> -->
                    </q-card-section>

                    <q-card-section>
                      <!-- 一个重命名的输入框, 一个删除按钮 -->
                    </q-card-section>
                    <q-card-actions align="right">
                      <q-btn
                        flat
                        @click="
                          async () => {
                            // 循环files, 并在每次上传成功后, 删除对应file

                            if (!files || files.length === 0) {
                              console.warn('No files selected for upload');
                              return;
                            }

                            // 使用slice()方法创建一个数组的浅拷贝, 避免因遍历过程中修改原始数组而导致的遍历中止
                            // slice也可以只截取数组的一部分, 类似golang的切片, 都是左闭右开区间。 如slice(2,4) 会从[1,2,3,4,5]中, 截取[3,4]
                            for (const file of files.slice()) {
                              try {
                                const re = await SendFileToServer(pkgPath, file);
                                if (re === true) {
                                  console.info(`File ${file.name} uploaded successfully`);
                                  // Remove the file from the list after successful upload
                                  const index = files.indexOf(file);
                                  if (index > -1) {
                                    files.splice(index, 1);
                                  }
                                } else {
                                  console.error(`File ${file.name} uploading error`);
                                  q.notify({
                                    type: 'negative',

                                    position: 'top',

                                    message: `File '${file.name}' addition failed`,

                                    timeout: 5,
                                  });
                                }
                              } catch (error) {
                                console.error(`Error uploading file ${file.name}:`, error);
                                q.notify({
                                  type: 'negative',

                                  position: 'top',

                                  message: `File '${file.name}' addition failed`,

                                  timeout: 5,
                                });
                              }
                            }
                          }
                        "
                        color="primary"
                        label="确认添加"
                      />
                      <q-btn flat label="Close" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
            </div>
            <!-- ------------------------------------------------------------------------载入声音文件的业务逻辑   end -->
            <q-stepper-navigation>
              <q-btn @click="step = 2" color="primary" label="Continue" />
            </q-stepper-navigation>
          </q-step>

          <!-- <q-step :name="2" title="键音制作" caption="Optional" icon="create_new_folder" :done="step > 2"> -->
          <q-step :name="2" title="裁剪定义声音" icon="add_comment" :done="step > 2">
            <div>根据载入的原始声音文件裁剪定义出需要的声音。</div>
            <div>此步骤不会影响原始声音文件。</div>
            <div>用户可针对同一声音文件裁剪定义出独立的多个需求声音。</div>
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

          <q-step :name="3" title="制作按键声音" icon="add_comment" :done="step > 3">
            <div>根据裁剪定义好的声音, 制作按键声音。</div>
            <div>每个按键音, 都包括了按下声音和抬起声音, 制作时需要分别定义它们。</div>
            <div>当然, 如果只需要其中之一, 也可根据需求自由制作。</div>
            <q-stepper-navigation>
              <q-btn @click="step = 4" color="primary" label="Continue" />
              <q-btn flat @click="step = 2" color="primary" label="Back" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>

          <q-step :name="4" title="对全局按键统一设置键音" icon="settings" :done="step > 3">
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
          </q-step>
        </q-stepper>
      </div>
    </q-scroll-area>
  </q-page>
</template>

<script setup lang="ts">
import { nanoid } from 'nanoid';
import { useQuasar } from 'quasar';
import { ConfigSet, LoadConfig, SendFileToServer } from 'src/boot/query/keytonePkg-query';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const q = useQuasar();
const { t } = useI18n();
const $t = t;

// 此路径没必要进行状态管理, 当用户退出此页面时, 自动清除即符合逻辑。
const pkgPath = nanoid();

// 此时由于是新建键音包, 因此是没有对应配置文件, 需要我们主动去创建的。 故第二个参数设置为true
LoadConfig(pkgPath, true);

const pkgName = ref<string>('');
watch(pkgName, (newVal) => {
  ConfigSet('package_name', pkgName.value);
});
pkgName.value = $t('KeyTonePackage.new.name.defaultValue');

const step = ref(1);

const addNewSoundFile = ref(false);
const files = ref<Array<File>>([]);

const editSoundFile = ref(false);

watch(files, () => {
  console.log(files.value);
});
</script>

<style lang="scss" scoped></style>
