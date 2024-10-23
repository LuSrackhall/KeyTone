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
                  label="载入新的声音源文件"
                  @click="
                    () => {
                      addNewSoundFile = !addNewSoundFile;
                    }
                  "
                >
                </q-btn>
                <q-dialog v-model="addNewSoundFile" backdrop-filter="invert(70%)">
                  <q-card>
                    <q-card-section class="row items-center q-pb-none text-h6"> 载入新的声音源文件 </q-card-section>

                    <q-card-section> <div>文件类型可以是WAV、MP3、OGG。</div></q-card-section>

                    <q-card-section>
                      <q-file
                        :class="['w-56']"
                        dense
                        v-model="files"
                        label="点此选择文件"
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
                  label="管理已载入的源文件"
                  @click="
                    () => {
                      editSoundFile = !editSoundFile;
                    }
                  "
                ></q-btn>
                <q-dialog v-model="editSoundFile" backdrop-filter="invert(70%)">
                  <q-card :class="['p-x-3']">
                    <q-card-section class="row items-center q-pb-none text-h6"> 管理已载入的源文件 </q-card-section>

                    <!-- <q-card-section> <div>请选择您想要修改或删除的声音源文件并执行对应操作。</div></q-card-section> -->

                    <q-card-section>
                      <q-select
                        outlined
                        stack-label
                        v-model="selectedSoundFile"
                        :options="soundFileList"
                        option-label="name"
                        label="选择要管理的源文件"
                        dense
                      />
                      <!-- option-label="name"
                       如果 :options 的元素类型是对象, 则有必要指定其中某个类型为字符串的字段作为label显示。
                       注意, 列表中显示的label名称, 是根据此字段来显示的。
                       (已选择框内显示的也是此字段)->
                       (
                        已选择框内显示的字段, 是否在重新选择后发生变更,
                        则需要由 :value 指定的字段来决定。就算label字段不同,
                        但只要:value指定的字段值是相同的, 已选择框内就不会发生label名称的变更。
                       )
                       -->
                      <!-- option-value="uuid"
                       如果不手动的显示指定, 默认会以整个对象作为区分标准(我们这里就是这种情况, 因为我们没有能作为uuid的字段),
                       但也可以指定具体字段作为区分标准( 可以类比v-for遍历时指定的 :key 值。),
                       当检测到变化, 才会更新已选择框内的label名。
                       这里我们不手动指定, 默认使用整个对象做区分,
                       因为我们的uuid可能和其它不同sha256中的uuid重复,
                       而sha256就更不用说了, 默认就会有相同sha256的项,
                       至于使用name也更不用说, 我们命名时允许相同的名称。
                       使用type这个字段作为唯一表示就更加没有讨论意义了。
                        -->
                    </q-card-section>

                    <!-- 分割线 -->
                    <q-separator v-if="selectedSoundFile.sha256 !== '' && selectedSoundFile.nameID !== ''" />

                    <!-- 以卡片形式展示选择的音频源文件 -->
                    <q-card-section
                      v-if="selectedSoundFile.sha256 !== '' && selectedSoundFile.nameID !== ''"
                      :class="['flex flex-col m-t-3']"
                    >
                      <q-card :class="['flex flex-col']">
                        <q-badge
                          transparent
                          color="orange"
                          :label="selectedSoundFile.name + selectedSoundFile.type"
                          :class="['absolute  overflow-visible']"
                          :style="{
                            // left: 88% 是刚好在最右边的, 但是随着字符串的增长, 左侧不会动, 仅向右追加溢出。
                            // 通过更改left的值, 可以通过左侧的溢出量, 来控制左侧的余量, 尽可能保持先向左增长, 右侧不动。
                            // 当左侧占整体的比例达到一定程度后, 右侧开始追加, 左侧不动。
                            // 计算长度前, 排除空格, 因为空格会引起长度计算错误
                            // TODO: 这个算法只是临时的, 后续可通过计算其真实的实际宽度, 来决定偏移距离, 以真正做到完全可控的占右侧70%
                            // TODO: 如果需要实现多选'编辑/删除'功能, 则需要将整个卡片, 抽离成独立的组件, 方便后续维护。
                            left:
                              selectedSoundFile.name.replace(/\s/g, '').length <= 20
                                ? 88 - selectedSoundFile.name.replace(/\s/g, '').length * 2 + '%'
                                : 88 - 20 * 2 + '%',
                          }"
                        />
                        <q-card-section
                          v-if="selectedSoundFile.sha256 !== '' && selectedSoundFile.nameID !== ''"
                          :class="['flex flex-col m-t-3']"
                        >
                          <!-- 一个重命名的输入框, 一个删除按钮 -->
                          <q-input
                            outlined
                            stack-label
                            dense
                            error-message="空文件名可能增加后续管理成本"
                            :error="
                              selectedSoundFile.name === '' ||
                              selectedSoundFile.name === undefined ||
                              selectedSoundFile.name === null
                            "
                            v-model="selectedSoundFile.name"
                            label="文件名(可更改)"
                          />

                          <q-btn
                            :class="['w-20 self-center bg-pink-700 text-zinc-50']"
                            dense
                            no-caps
                            label="删除"
                            icon="flight_takeoff"
                            @click="
                              async () => {
                                const re = await SoundFileDelete(
                                  pkgPath,
                                  selectedSoundFile.sha256,
                                  selectedSoundFile.nameID,
                                  selectedSoundFile.type
                                );
                                if (re) {
                                  q.notify({
                                    type: 'positive',
                                    position: 'top',
                                    message: '删除成功',
                                    timeout: 5,
                                  });
                                  // 如果sdk中删除操作执行成功, 则前端清除相关的结构体对象
                                  selectedSoundFile = {
                                    sha256: '',
                                    nameID: '',
                                    name: '',
                                    type: '',
                                  };
                                } else {
                                  q.notify({
                                    type: 'negative',
                                    position: 'top',
                                    message: '删除失败',
                                    timeout: 5,
                                  });
                                }
                              }
                            "
                          >
                          </q-btn>
                        </q-card-section>
                      </q-card>
                    </q-card-section>
                    <q-card-actions align="right">
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

            <!-- ------------------------------------------------------------------------裁剪定义声音的业务逻辑 start -->
            <div>
              <!-- ------------------------------------------------------------------------------ 制作新的声音 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  label="制作新的声音"
                  @click="
                    () => {
                      createNewSound = !createNewSound;
                    }
                  "
                >
                </q-btn>
                <q-dialog v-model="createNewSound" backdrop-filter="invert(70%)">
                  <q-card>
                    <q-card-section class="row items-center q-pb-none text-h6"> 制作新的声音 </q-card-section>
                    <q-card-section>
                      <q-input
                        outlined
                        stack-label
                        dense
                        v-model="soundName"
                        label="为声音命名(非必填)"
                        :placeholder="
                          sourceFileForSound.name + '     - ' + ' [' + soundStartTime + ' ~ ' + soundEndTime + ']'
                        "
                        :input-style="{ textOverflow: 'ellipsis' }"
                        :input-class="'text-truncate'"
                      >
                        <template v-slot:append>
                          <q-icon name="info" color="primary">
                            <q-tooltip :class="['bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                              {{
                                '声音的名称 : \n' +
                                (soundName === ''
                                  ? sourceFileForSound.name + ' - ' + ' [' + soundStartTime + ' ~ ' + soundEndTime + ']'
                                  : soundName)
                              }}
                            </q-tooltip>
                          </q-icon>
                        </template>
                      </q-input>
                    </q-card-section>
                    <q-card-section>
                      <q-select
                        outlined
                        stack-label
                        v-model="sourceFileForSound"
                        :options="soundFileList"
                        option-label="name"
                        label="声音的源文件"
                        dense
                      />
                    </q-card-section>
                    <q-card-section>
                      <div class="text-[10px] text-gray-600">从声音源文件中裁剪定义出我们需要的声音</div>
                      <q-input
                        outlined
                        stack-label
                        dense
                        v-model="soundStartTime"
                        label="声音开始时间(毫秒)"
                        type="number"
                      />
                      <q-input
                        outlined
                        stack-label
                        dense
                        v-model="soundEndTime"
                        label="声音结束时间(毫秒)"
                        type="number"
                      />
                    </q-card-section>
                    <q-card-actions align="right">
                      <q-btn flat label="Close" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
              <!-- -------------------------------------------------------------------------------编辑已有声音 -->
              <div></div>
            </div>
            <!-- ------------------------------------------------------------------------裁剪定义声音的业务逻辑   end -->
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
import { debounce } from 'lodash';
import { nanoid } from 'nanoid';
import { useQuasar } from 'quasar';
import {
  ConfigGet,
  ConfigSet,
  LoadConfig,
  SendFileToServer,
  SoundFileRename,
  SoundFileDelete,
} from 'src/boot/query/keytonePkg-query';
import { useAppStore } from 'src/stores/app-store';
import { onBeforeMount, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const q = useQuasar();
const { t } = useI18n();
const $t = t;
const app_store = useAppStore();

// 此路径没必要进行状态管理, 当用户退出此页面时, 自动清除即符合逻辑。
const pkgPath = nanoid();

// 防止空字符串触发不能为空的提示, 虽然初始化时只有一瞬间, 但也不希望看到
const pkgName = ref<string>($t('KeyTonePackage.new.name.defaultValue'));

const step = ref(1);

const addNewSoundFile = ref(false);
const files = ref<Array<File>>([]);
watch(files, () => {
  console.debug('观察files=', files.value);
});

const editSoundFile = ref(false);
// 用于初步映射配置文件中的 audio_files 对象, 并将其转换为数组, 并将数组元素转换成对象, 其中包含sha256和value两个key
const audioFiles = ref<Array<any>>([]);
// 用于audioFiles映射后的进一步映射, 主要拆分出value中name的每个值, 并一一对于sha256, 形成ui列表可用的最终数组。
const soundFileList = ref<Array<{ sha256: string; nameID: string; name: string; type: string }>>([]);
const selectedSoundFile = ref<{ sha256: string; nameID: string; name: string; type: string }>({
  sha256: '',
  nameID: '',
  name: '',
  type: '',
});

const createNewSound = ref(false);
const soundName = ref<string>('');
const sourceFileForSound = ref<{ sha256: string; nameID: string; name: string; type: string }>({
  sha256: '',
  nameID: '',
  name: '',
  type: '',
});
const soundStartTime = ref<number>(0);
const soundEndTime = ref<number>(0);

onBeforeMount(async () => {
  // 此时由于是新建键音包, 因此是没有对应配置文件, 需要我们主动去创建的。 故第二个参数设置为true
  // 这也是我们加载页面前必须确定的事情, 否则无法进行后续操作, 一切以配置文件为前提。
  await LoadConfig(pkgPath, true);

  // 使用i18n, 初始化键音包名称, 在此处手动设置, 是为了防止后续初始化将sdk中默认的名称, 被初始化到pkgName中。
  // 因此在初始化前按照提前设置好的i18n做为默认名称, 故手动发送请求以在数据初始化前更改sdk中的默认名称。
  await ConfigSet('package_name', $t('KeyTonePackage.new.name.defaultValue'));

  // 数据初始化
  await initData();
  // 将初始化数据的操作封装成一个函数, 并设置为异步函数, 以便使用await调用
  async function initData() {
    await ConfigGet('get_all_value').then((req) => {
      // console.debug('打印观察获取的值', req);
      if (req === false) {
        // 此时, 说明GetItem_sqlite请求过程中, 出错了, 因此需要错误通知, 并让用户重新启动, 防止用户因继续使用造成的存储设置被初始覆盖
        q.notify({
          type: 'negative',
          position: 'top',
          message: '键音包配置文件读取失败',
          timeout: 100000,
        });
        return;
      }

      // TIPS: 由于采取各设置独立的录入即判别方式, 不再依赖整体的JSON字符串, 因此此if判断后续可能没必要存在(目前暂时保留)
      // 第一次进入本应用, 设置本就该是空的, 此时无需对我们的设置项进行任何操作, 也无需做任何通知。
      // 但为防止后续的JSON.parse报错, 因此此处也是必不可少的(因为只要非首次, 就不可能为空, watchEffect是立即执行的, 也就是说至少整体的结构是正常入库的)
      if (req === '' || req === '{}' || req === null) {
        return;
      }

      // // 若有设置数据, 则取出 TIPS: 注意, 这里的设置是直接读出的一个json对象, 而不是需要解析的json字符串
      // const settingStorage = JSON.parse(req);

      const data = req;

      // 键音包名称初始化。 (不过由于这里是新建键音包, 这个不出意外的话一开始是undefined
      if (data.package_name !== undefined) {
        pkgName.value = data.package_name;
      }

      // 以载入的声音文件列表初始化。  (不过由于这里是新建键音包, 这个不出意外的话一开始是undefined
      if (data.audio_files !== undefined) {
        // keyTonePkgData.audio_files 是一个从后端获取的对象, 通过此方式可以简便的将其转换为数组, 数组元素为原对象中的key和value(增加了这两个key)
        const audioFilesArray = Object.entries(data.audio_files).map(([key, value]) => ({
          sha256: key,
          value: value,
        }));
        audioFiles.value = audioFilesArray;
      }
    });

    watch(pkgName, (newVal) => {
      ConfigSet('package_name', pkgName.value);
    });

    // 2.配置文件中audio_files的进一步映射变更, 获取我们最终需要的结构
    watch(audioFiles, (newVal) => {
      console.debug('观察audioFiles=', audioFiles.value);
      // 为了更容易理解, 故引入audioFiles这一变量, 做初步映射, audioFiles只是过程值, 我们最终需要对此过程值做进一步映射, 形成soundFileList
      const tempSoundFileList: Array<any> = [];

      audioFiles.value.forEach((item) => {
        // 此处必须判断其是否存在, 否则会引起Object.entries报错崩溃, 影响后续流程执行。
        if (item.value.name !== undefined && item.value.name !== null) {
          Object.entries(item.value.name).forEach(([nameID, name]) => {
            tempSoundFileList.push({ sha256: item.sha256, nameID: nameID, name: name, type: item.value.type });
          });
        }
      });
      soundFileList.value = tempSoundFileList;
    });

    // 3.观察进一步映射变更后, 最终需要的audio_file映射, 即我们的soundFileList。
    watch(soundFileList, (newVal) => {
      console.debug('观察soundFileList=', soundFileList.value);
    });

    //  - completed(已完成)   TODO:
    // 4.观察selectedSoundFile的变化, 当selectedSoundFile变化时,
    //   说明用户做了对应修改, 此时需要向sdk发送请求, 更新配置文件中的对应值, 然后触发sse形成闭环。
    //   当然, 删除时同理, 但删除是独立的按钮点击后手动触发对应函数, 以向sdk发送请求, 不由此处的数据驱动。
    watch(
      // TIPS: 对于ref的响应式变量, 如果直接整体监听, 则内部的某个值变化时, 不会触发监听。需要使用返回值的函数, 对固定字段进行监听。
      () => selectedSoundFile.value.name,
      (newVal) => {
        console.debug('观察selectedSoundFile=', selectedSoundFile.value);
        if (selectedSoundFile.value.sha256 !== '' && selectedSoundFile.value.nameID !== '') {
          SoundFileRename(selectedSoundFile.value.sha256, selectedSoundFile.value.nameID, selectedSoundFile.value.name);
        }
      }
    );
  }

  const pkgNameDelayed = debounce(
    (keyTonePkgData: any) => {
      pkgName.value = keyTonePkgData.package_name;
    },
    800,
    { trailing: true }
  );

  // 将后端从键音包配置文件中获取的全部数据, 转换前端可用的键音包数据。(只要配置文件变更, 就会触发相关sse发送, 此处就会接收)
  function sseDataToKeyTonePkgData(keyTonePkgData: any) {
    // 键音包名称初始化。 (不过由于这里是新建键音包, 这个不出意外的话一开始是undefined
    if (keyTonePkgData.package_name !== undefined) {
      pkgNameDelayed.cancel();
      pkgNameDelayed(keyTonePkgData);
    }

    // 1. 初步映射配置文件中的audio_files到audioFiles。(只要配置文件变更, 就会触发相关sse发送, 此处就会接收)
    //    使用audioFiles作为中间值, 而不是一步到位的映射, 是为代码的可读性, 后续阅读理解是方便。
    if (keyTonePkgData.audio_files !== undefined) {
      // keyTonePkgData.audio_files 是一个从后端获取的对象, 通过此方式可以简便的将其转换为数组, 数组元素为原对象中的key和value(增加了这两个key)
      const audioFilesArray = Object.entries(keyTonePkgData.audio_files).map(([key, value]) => ({
        sha256: key,
        value: value,
      }));
      audioFiles.value = audioFilesArray;
    } else {
      // 此处else是为防止最后一项的audio_files为undefined, 而导致的删除最后一项音频源文件后, audioFiles值无法清空, 从而导致无法触发soundFileList的变更, 从而ui界面导致无法删除最后一项音频源文件。
      audioFiles.value = [];
    }
  }
  const debounced_sseDataToSettingStore = debounce<(keyTonePkgData: any) => void>(sseDataToKeyTonePkgData, 30, {
    trailing: true,
  });

  app_store.eventSource.addEventListener(
    'messageAudioPackage',
    function (e) {
      console.debug('后端钩子函数中的值 = ', e.data);

      const data = JSON.parse(e.data);

      if (data.key === 'get_all_value') {
        debounced_sseDataToSettingStore.cancel;
        debounced_sseDataToSettingStore(data.value);
      }
    },
    false
  );
});
</script>

<style lang="scss" scoped></style>
