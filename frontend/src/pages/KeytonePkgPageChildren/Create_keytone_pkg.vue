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
          <q-step :name="1" title="载入音频文件" icon="create_new_folder" :done="step > 1">
            <div>为此键音包载入原始的音频文件供后续步骤使用。</div>
            <!-- <div>文件类型可以是WAV、MP3、OGG等。</div> -->
            <!-- <div>原始音频文件的数量不定,可根据您的制作喜好来决定。</div> -->
            <!-- <q-card class="bg-slate-500" :class="['p-2']"> -->
            <!-- <q-btn :class="['bg-zinc-300']" label="添加新的声音源文件"></q-btn>
            <div :class="['p-2 text-zinc-300']">或</div>
            <q-btn :class="['bg-zinc-300']" label="选择声音源文件以进行编辑"></q-btn> -->
            <!-- </q-card> -->
            <!-- ------------------------------------------------------------------------载入音频文件的业务逻辑 start -->
            <div>
              <!-- ------------------------------------------------------------------------------ 添加新的音频源文件 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  label="载入新的音频源文件"
                  @click="
                    () => {
                      addNewSoundFile = !addNewSoundFile;
                    }
                  "
                >
                </q-btn>
                <q-dialog v-model="addNewSoundFile" backdrop-filter="invert(70%)">
                  <q-card>
                    <q-card-section class="row items-center q-pb-none text-h6"> 载入新的音频源文件 </q-card-section>

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
                                const re = await SendFileToServer(file);
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
              <!-- -------------------------------------------------------------------------------编辑已有音频源文件 -->
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
                        :option-label="(item) => item.name + item.type"
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
                    <q-separator v-if="selectedSoundFile.sha256 !== '' && selectedSoundFile.name_id !== ''" />

                    <!-- 以卡片形式展示选择的音频源文件 -->
                    <q-card-section
                      v-if="selectedSoundFile.sha256 !== '' && selectedSoundFile.name_id !== ''"
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
                          v-if="selectedSoundFile.sha256 !== '' && selectedSoundFile.name_id !== ''"
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
                                  selectedSoundFile.sha256,
                                  selectedSoundFile.name_id,
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
                                    name_id: '',
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
            <div>
              根据载入的原始音频文件裁剪定义出需要的声音。
              <q-icon name="info" color="primary">
                <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                  <div>通过此步骤制作的声音不会影响声音源文件。</div>
                  <div>用户可针对同一声音源文件裁剪定义出多个独立的声音。</div>
                </q-tooltip>
              </q-icon>
            </div>
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
                    <q-card-section :class="['p-b-1']">
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
                            <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
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
                    <q-card-section :class="['p-b-1']">
                      <q-select
                        outlined
                        stack-label
                        v-model="sourceFileForSound"
                        :options="soundFileList"
                        :option-label="(item) => item.name + item.type"
                        label="声音的源文件"
                        dense
                      />
                    </q-card-section>

                    <q-card-section :class="['p-b-1']">
                      <div class="text-[15px] text-gray-600">
                        从声音源文件中裁剪定义出我们需要的声音
                        <q-icon name="info" color="primary">
                          <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                            声音的开始时间和结束时间, 以毫秒为单位。由于是按键使用, 所以时间不宜过长。
                          </q-tooltip>
                        </q-icon>
                      </div>
                      <!-- TIPS: 注意number类型使用时, 需要使用  v-model.number 。 因为使用后可以自动处理 01、 00.55 这种, 会将其自动变更为 1、 0.55 -->
                      <div class="flex flex-row">
                        <q-input
                          :class="['w-1/2 p-r-1']"
                          outlined
                          stack-label
                          dense
                          v-model.number="soundStartTime"
                          label="声音开始时间(毫秒)"
                          type="number"
                          error-message="时间不能为负数"
                          :error="soundStartTime < 0"
                        />
                        <q-input
                          :class="['w-1/2 p-l-1']"
                          outlined
                          stack-label
                          dense
                          v-model.number="soundEndTime"
                          label="声音结束时间(毫秒)"
                          type="number"
                          error-message="时间不能为负数"
                          :error="soundEndTime < 0"
                        />
                      </div>
                    </q-card-section>
                    <q-card-section :class="['p-y-0']">
                      <q-input
                        outlined
                        stack-label
                        dense
                        v-model.number="soundVolume"
                        label="声音音量"
                        type="number"
                        :step="0.1"
                      >
                        <template v-slot:append>
                          <q-icon name="info" color="primary">
                            <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                              0为原始音量, 大于0为提升音量, 小于0为降低音量。常用的步进为0.1,
                              当然您也可以手动输入以做更细腻的调整, 如0.0001等。请在每次音量调整后重新预览以查看效果,
                              防止音量过小或过大。
                            </q-tooltip>
                          </q-icon>
                        </template>
                      </q-input>
                    </q-card-section>
                    <q-card-actions align="right">
                      <q-btn
                        dense
                        @click="
                          previewSound({
                            source_file_for_sound: sourceFileForSound,
                            cut: {
                              start_time: soundStartTime,
                              end_time: soundEndTime,
                              volume: soundVolume,
                            },
                          })
                        "
                        label="预览声音"
                        color="secondary"
                      >
                        <q-tooltip
                          :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-xs']"
                          :delay="600"
                        >
                          按键用的声音通常很短, 因此预览的声音会并发播放, 且不提供进度条和停止按钮。
                        </q-tooltip>
                      </q-btn>
                      <q-btn
                        @click="
                          saveSoundConfig({
                            // source_file_for_sound: sourceFileForSound ,  // 由于js/ts中，传递对象时是通过引用传递的。这意味着函数接收到的是对象的引用，而不是对象的副本。因此，函数内部可以访问和修改对象的所有属性。包括不希望有的name字段。
                            // source_file_for_sound: { ...sourceFileForSound },  // 使用对象解构（{ ...sourceFileForSound }）会复制sourceFileForSound对象的所有可枚举属性到一个新的对象中。因此，如果sourceFileForSound对象中包含name字段，解构操作会将其包含在新对象中, 本质仍是将一个包含name字段的对象, 赋值给所需参数。
                            source_file_for_sound: {
                              sha256: sourceFileForSound.sha256,
                              name_id: sourceFileForSound.name_id,
                              type: sourceFileForSound.type,
                            }, //手动选择字段：通过手动选择字段，可以确保只传递所需的字段，而不会包含任何不需要的字段。

                            name: soundName,
                            cut: {
                              start_time: soundStartTime,
                              end_time: soundEndTime,
                              volume: soundVolume,
                            },
                            onSuccess: () => {
                              // 重置表单状态
                              soundName = '';
                              sourceFileForSound = {
                                sha256: '',
                                name_id: '',
                                name: '',
                                type: '',
                              };
                              soundStartTime = 0;
                              soundEndTime = 0;
                              soundVolume = 0.0;
                            },
                          })
                        "
                        label="确定添加"
                        color="primary"
                      />
                      <q-btn flat label="Close" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>

              <div :class="['p-2 text-zinc-600']">或</div>

              <!-- -------------------------------------------------------------------------------编辑已有声音 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  label="编辑已有声音"
                  @click="
                    () => {
                      if (soundList.length === 0) {
                        q.notify({
                          type: 'warning',
                          message: '当前没有可编辑的声音',
                          position: 'top',
                        });
                        return;
                      }
                      showEditSoundDialog = true;
                    }
                  "
                >
                </q-btn>
                <q-dialog v-model="showEditSoundDialog" backdrop-filter="invert(70%)">
                  <q-card>
                    <q-card-section
                      class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                    >
                      编辑已有声音
                    </q-card-section>
                    <q-card-section>
                      <q-select
                        outlined
                        stack-label
                        v-model="selectedSound"
                        :options="soundList"
                        :option-label="(item: any) => {
                          // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
                          if (item.soundValue.name !== '' && item.soundValue.name !== undefined) {
                            return item.soundValue.name
                          } else {
                            return soundFileList.find(
                              (soundFile) =>
                                soundFile.sha256 === item.soundValue.source_file_for_sound.sha256 &&
                                soundFile.name_id === item.soundValue.source_file_for_sound.name_id
                            )?.name + '     - ' + ' [' + item.soundValue.cut.start_time + ' ~ ' + item.soundValue.cut.end_time + ']'
                          }
                        }"
                        :option-value="(item:any) =>{
                          // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                          return item.soundKey
                        }"
                        label="选择要管理的声音"
                        dense
                      />
                    </q-card-section>
                    <!-- 以卡片形式展示选择的声音 -->
                    <q-card-section
                      :class="['flex flex-col m-t-3']"
                      v-if="selectedSound?.soundKey !== '' && selectedSound !== undefined"
                    >
                      <q-card :class="['flex flex-col px-3 pb-3']">
                        <q-badge
                          transparent
                          color="orange"
                          :label=" selectedSound?.soundValue.name === ''
                                ? soundFileList.find(
                                    (soundFile:any) =>
                                      soundFile.sha256 === selectedSound?.soundValue.source_file_for_sound.sha256 &&
                                      soundFile.name_id === selectedSound?.soundValue.source_file_for_sound.name_id
                                  )?.name +
                                    ' - ' +
                                    ' [' +
                                    selectedSound?.soundValue.cut.start_time +
                                    ' ~ ' +
                                    selectedSound?.soundValue.cut.end_time +
                                    ']'
                                : selectedSound?.soundValue.name"
                          :class="['absolute  overflow-visible']"
                        />
                        <q-card-section :class="['p-b-1 mt-3']">
                          <q-input
                            outlined
                            stack-label
                            dense
                            v-model="selectedSound.soundValue.name"
                            label="为声音命名(非必填)"
                            :placeholder="
                              soundFileList.find(
                                (soundFile) =>
                                  soundFile.sha256 === selectedSound?.soundValue.source_file_for_sound.sha256 &&
                                  soundFile.name_id === selectedSound?.soundValue.source_file_for_sound.name_id
                              )?.name +
                              '     - ' +
                              ' [' +
                              selectedSound.soundValue.cut.start_time +
                              ' ~ ' +
                              selectedSound.soundValue.cut.end_time +
                              ']'
                            "
                            :input-style="{ textOverflow: 'ellipsis' }"
                            :input-class="'text-truncate'"
                          >
                            <template v-slot:append>
                              <q-icon name="info" color="primary">
                                <q-tooltip
                                  :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']"
                                >
                                  {{
                                    '声音的名称 : \n' +
                                    (selectedSound.soundValue.name === ''
                                      ? soundFileList.find(
                                          (soundFile) =>
                                            soundFile.sha256 ===
                                              selectedSound?.soundValue.source_file_for_sound.sha256 &&
                                            soundFile.name_id ===
                                              selectedSound?.soundValue.source_file_for_sound.name_id
                                        )?.name +
                                        ' - ' +
                                        ' [' +
                                        selectedSound.soundValue.cut.start_time +
                                        ' ~ ' +
                                        selectedSound.soundValue.cut.end_time +
                                        ']'
                                      : selectedSound.soundValue.name)
                                  }}
                                </q-tooltip>
                              </q-icon>
                            </template>
                          </q-input>
                        </q-card-section>
                        <q-card-section :class="['p-b-1']">
                          <q-select
                            outlined
                            stack-label
                            v-model="selectedSound.soundValue.source_file_for_sound"
                            :options="soundFileList"
                            :option-label="(item: any) => {
                              // 此处的:options本身就是soundFileList, 因此直接通过find查找并返回自身即可。
                              // 此处仅是为了避免直接使用其中元素值的name字段, 以让selectedSound.soundValue.source_file_for_sound也可享受name变化时的实时更新。
                              const soundFile = soundFileList.find(
                                (soundFile) =>
                                  soundFile.sha256 === item.sha256 &&
                                  soundFile.name_id === item.name_id
                              )
                              return soundFile ? soundFile.name + soundFile.type : ''
                            }"
                            :option-value="(item: any) => {
                              // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                              return item.sha256 + item.name_id
                            }"
                            label="声音的源文件"
                            dense
                          />
                        </q-card-section>

                        <q-card-section :class="['p-b-1']">
                          <div class="text-[15px] text-gray-600">
                            从声音源文件中裁剪定义出我们需要的声音
                            <q-icon name="info" color="primary">
                              <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                                声音的开始时间和结束时间, 以毫秒为单位。由于是按键使用, 所以时间不宜过长。
                              </q-tooltip>
                            </q-icon>
                          </div>
                          <!-- TIPS: 注意number类型使用时, 需要使用  v-model.number 。 因为使用后可以自动处理 01、 00.55 这种, 会将其自动变更为 1、 0.55 -->
                          <div class="flex flex-row">
                            <q-input
                              :class="['w-1/2 p-r-1']"
                              outlined
                              stack-label
                              dense
                              v-model.number="selectedSound.soundValue.cut.start_time"
                              label="声音开始时间(毫秒)"
                              type="number"
                              error-message="时间不能为负数"
                              :error="selectedSound.soundValue.cut.start_time < 0"
                            />
                            <q-input
                              :class="['w-1/2 p-l-1']"
                              outlined
                              stack-label
                              dense
                              v-model.number="selectedSound.soundValue.cut.end_time"
                              label="声音结束时间(毫秒)"
                              type="number"
                              error-message="时间不能为负数"
                              :error="selectedSound.soundValue.cut.end_time < 0"
                            />
                          </div>
                        </q-card-section>
                        <q-card-section :class="['p-y-0']">
                          <q-input
                            outlined
                            stack-label
                            dense
                            v-model.number="selectedSound.soundValue.cut.volume"
                            label="声音音量"
                            type="number"
                            :step="0.1"
                          >
                            <template v-slot:append>
                              <q-icon name="info" color="primary">
                                <q-tooltip
                                  :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']"
                                >
                                  0为原始音量, 大于0为提升音量, 小于0为降低音量。常用的步进为0.1,
                                  当然您也可以手动输入以做更细腻的调整,
                                  如0.0001等。请在每次音量调整后重新预览以查看效果, 防止音量过小或过大。
                                </q-tooltip>
                              </q-icon>
                            </template>
                          </q-input>
                        </q-card-section>

                        <!-- 添加按钮组 -->
                        <q-card-section :class="['flex justify-center gap-4']">
                          <q-btn
                            dense
                            color="secondary"
                            icon="play_arrow"
                            label="预览声音"
                            @click="
                              previewSound({
                                source_file_for_sound: selectedSound.soundValue.source_file_for_sound,
                                cut: selectedSound.soundValue.cut,
                              })
                            "
                          >
                            <q-tooltip
                              :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-xs']"
                              :delay="600"
                            >
                              按键用的声音通常很短, 因此预览的声音会并发播放, 且不提供进度条和停止按钮。
                            </q-tooltip>
                          </q-btn>

                          <q-btn
                            dense
                            color="primary"
                            icon="save"
                            label="确认修改"
                            @click="
                              saveSoundConfig({
                                soundKey: selectedSound.soundKey,
                                source_file_for_sound: {
                                  sha256: selectedSound.soundValue.source_file_for_sound.sha256,
                                  name_id: selectedSound.soundValue.source_file_for_sound.name_id,
                                  type: selectedSound.soundValue.source_file_for_sound.type,
                                },
                                name: selectedSound.soundValue.name,
                                cut: selectedSound.soundValue.cut,
                                onSuccess: () => {
                                  // q.notify({
                                  //   type: 'positive',
                                  //   position: 'top',
                                  //   message: '修改成功',
                                  //   timeout: 5,
                                  // });
                                },
                              })
                            "
                          />

                          <q-btn
                            dense
                            color="negative"
                            icon="delete"
                            label="删除声音"
                            @click="
                              deleteSound({
                                soundKey: selectedSound.soundKey,
                                onSuccess: () => {
                                  selectedSound = undefined;
                                  q.notify({
                                    type: 'positive',
                                    position: 'top',
                                    message: '删除成功',
                                    timeout: 5,
                                  });
                                },
                              })
                            "
                          />
                        </q-card-section>
                      </q-card>
                    </q-card-section>

                    <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
                      <q-btn flat label="Close" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
            </div>
            <!-- ------------------------------------------------------------------------裁剪定义声音的业务逻辑   end -->
            <q-stepper-navigation>
              <q-btn @click="step = 3" color="primary" label="Continue" />
              <q-btn flat @click="step = 1" color="primary" label="Back" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>

          <q-step :name="3" title="铸造至臻键音" icon="add_comment" :done="step > 3">
            <div>
              <span>键音, 实际就是按键声音或称按键音。</span>
              <span>本步骤默认根据裁剪定义好的声音, 制作按键音。</span>
              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                  <span>理论上每个按键音, 都应包括按下和抬起声音。<br /></span>
                  <span>制作时需分别定义它们, 但这并不是强制的。<br /></span>
                  <span>按下声音和抬起声音都可根据需要选择性定义。<br /></span>
                  <span>可以都定义, 也可以只定义其中一个, 或者都不定义。<br /></span>
                </q-tooltip>
              </q-icon>
            </div>
            <!-- ------------------------------------------------------------------------制作按键声音的业务逻辑 start -->
            <div>
              <!-- ------------------------------------------------------------------------------ 制作新的按键音 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  label="制作新的按键音"
                  @click="
                    () => {
                      createNewKeySound = !createNewKeySound;
                    }
                  "
                >
                </q-btn>
                <q-dialog v-model="createNewKeySound" backdrop-filter="invert(70%)">
                  <q-card :class="['min-w-[90%]']">
                    <q-card-section class="row items-center q-pb-none text-h6"> 制作新的按键音 </q-card-section>
                    <q-card-section :class="['p-b-1']">
                      <q-input
                        outlined
                        stack-label
                        dense
                        v-model="keySoundName"
                        label="按键音名称"
                        :placeholder="'新的按键音'"
                      />
                      <div class="flex flex-col mt-3">
                        <q-btn
                          :class="['bg-zinc-300 my-7 w-70% self-center']"
                          label="配置按下声音"
                          @click="configureDownSound = true"
                        />
                        <q-dialog v-model="configureDownSound" backdrop-filter="invert(70%)">
                          <q-card :class="['min-w-[80%]']">
                            <q-card-section class="row items-center q-pb-none text-h6"> 配置按下声音 </q-card-section>
                            <q-card-section>
                              <!-- 使用选择框选择模式 -->
                              <q-select
                                outlined
                                stack-label
                                v-model="playModeForDown"
                                :options="playModeOptions"
                                label="选择播放模式"
                                dense
                              />
                            </q-card-section>
                            <q-card-section>
                              <!-- 选择声音的选项，支持多选 -->
                              <q-select
                                outlined
                                stack-label
                                v-model="selectedSoundsForDown"
                                :options="downSoundList"
                                :option-label="(item: any) => {
                                  if (item.type === 'audio_files') {
                                    return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.name + soundFileList.find( (soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.type
                                  }
                                  if (item.type === 'sounds') {
                                    // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
                                    if (item.value.soundValue?.name !== '' && item.value.soundValue?.name !== undefined) {
                                      return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (soundList.find((sound) => sound.soundKey === item.value.soundKey)?.soundValue.name)
                                    } else {
                                      return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                        soundFileList.find(
                                          (soundFile:any) =>
                                            soundFile.sha256 === item.value.soundValue?.source_file_for_sound?.sha256 &&
                                            soundFile.name_id === item.value.soundValue?.source_file_for_sound?.name_id
                                        )?.name +
                                        '     - ' +
                                        ' [' +
                                        item.value.soundValue?.cut?.start_time +
                                        ' ~ ' +
                                        item.value.soundValue?.cut?.end_time +
                                        ']'
                                      );
                                    }
                                  }
                                  if (item.type === 'key_sounds') {
                                    return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (item.value.keySoundValue?.name);
                                  }
                                }"
                                :option-value="(item: any) => {
                                  // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                  if (item.type === 'audio_files'){
                                    return item.value.sha256 + item.value.name_id
                                  }
                                  if(item.type === 'sounds'){
                                    return item.value.soundKey
                                  }
                                  if(item.type === 'key_sounds'){
                                    return item.value.keySoundKey
                                  }
                                }"
                                label="选择声音 (多选)"
                                multiple
                                use-chips
                                dense
                                :max-values="maxSelectionForDown"
                                counter
                                error-message="独立模式下, 至多选择一个声音"
                                :注释="
                                  () => {
                                    // 这种写法, 实际上是绑定了箭头函数, 而非函数的返回值。
                                    /* :error = '()=>{
                                          return true
                                        }'
                                    */
                                    // 你需要手动调用它(像下方这样)
                                    /* :error = '(()=>{
                                          return true
                                        })()'
                                    */
                                    // 不过, 这也仅调用了一次。 如果函数内部有响应式变量, 它并不会因触发响应式变更而重新被调用
                                    // 因此, 想要用函数的方式并且响应式生效
                                    // * 常用做法是在script中创建一个计算属性, 并在此处绑定它(而不是绑定某个函数)
                                    // * 还有一种做法是, 直接绑定对应的表达式(表达式中如果有响应式变量的变化, 表达式会被重新计算), 就像下面这样
                                  }
                                "
                                :error="
                                  playModeForDown.mode === 'single'
                                    ? selectedSoundsForDown.length > 1
                                      ? true
                                      : false
                                    : false
                                "
                                ref="downSoundSelectDom"
                                @update:model-value="downSoundSelectDom?.hidePopup()"
                              />
                              <div class="h-3">
                                <q-option-group
                                  dense
                                  v-model="downTypeGroup"
                                  :options="options"
                                  type="checkbox"
                                  class="absolute left-8"
                                >
                                  <template #label-0="props">
                                    <q-item-label>
                                      {{ props.label }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>本软件支持将源文件直接作为声音。</div>
                                          <div>前提是</div>
                                          <div>这个源文件本身就就可作为独立的无需裁剪的声音。</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                  <template v-slot:label-1="props">
                                    <q-item-label>
                                      {{ props.label }}
                                      <q-icon name="info" color="primary" class="p-l-4.5 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>默认仅从声音列表中选择。</div>
                                          <div>如有需求也可勾选其它受支持列表。</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                  <template v-slot:label-2="props">
                                    <q-item-label>
                                      {{ props.label }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>本软件支持将其它按键音作为声音。</div>
                                          <div>或者说</div>
                                          <div>本软件支持继承其它按键音的配置。</div>
                                          <div>⬇</div>
                                          <div>原则: 按下继承按下, 抬起继承抬起。</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                </q-option-group>
                              </div>
                            </q-card-section>
                            <q-card-actions align="right">
                              <q-btn flat label="Close" color="primary" v-close-popup />
                            </q-card-actions>
                          </q-card>
                        </q-dialog>
                        <q-btn
                          :class="['bg-zinc-300  m-b-7 w-70% self-center']"
                          label="配置抬起声音"
                          @click="configureUpSound = true"
                        />
                        <q-dialog v-model="configureUpSound" backdrop-filter="invert(70%)">
                          <q-card :class="['min-w-[80%]']">
                            <q-card-section class="row items-center q-pb-none text-h6"> 配置抬起声音 </q-card-section>
                            <q-card-section>
                              <!-- 使用选择框选择模式 -->
                              <q-select
                                outlined
                                stack-label
                                v-model="playModeForUp"
                                :options="playModeOptions"
                                label="选择播放模式"
                                dense
                              />
                            </q-card-section>
                            <q-card-section>
                              <!-- 选择声音的选项，支持多选 -->
                              <q-select
                                outlined
                                stack-label
                                v-model="selectedSoundsForUp"
                                :options="upSoundList"
                                :option-label="(item: any) => {
                                  if (item.type === 'audio_files') {
                                    return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.name + soundFileList.find( (soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.type;
                                  }
                                  if (item.type === 'sounds') {
                                    // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
                                    if (item.value.soundValue?.name !== '' && item.value.soundValue?.name !== undefined) {
                                      return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;'+ (soundList.find((sound) => sound.soundKey === item.value.soundKey)?.soundValue.name)
                                    } else {
                                      return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                        soundFileList.find(
                                          (soundFile:any) =>
                                            soundFile.sha256 === item.value.soundValue?.source_file_for_sound?.sha256 &&
                                            soundFile.name_id === item.value.soundValue?.source_file_for_sound?.name_id
                                        )?.name +
                                        '     - ' +
                                        ' [' +
                                        item.value.soundValue?.cut?.start_time +
                                        ' ~ ' +
                                        item.value.soundValue?.cut?.end_time +
                                        ']'
                                      );
                                    }
                                  }
                                  if (item.type === 'key_sounds') {
                                    return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (item.value.keySoundValue?.name);
                                  }
                                }"
                                :option-value="(item: any) => {
                                  // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                  if (item.type === 'audio_files'){
                                    return item.value.sha256 + item.value.name_id
                                  }
                                  if(item.type === 'sounds'){
                                    return item.value.soundKey
                                  }
                                  if(item.type === 'key_sounds'){
                                    return item.value.keySoundKey
                                  }
                                }"
                                label="选择声音 (多选)"
                                multiple
                                use-chips
                                dense
                                :max-values="maxSelectionForUp"
                                counter
                                error-message="独立模式下, 至多选择一个声音"
                                :注释="
                                  () => {
                                    // 这种写法, 实际上是绑定了箭头函数, 而非函数的返回值。
                                    /* :error = '()=>{
                                          return true
                                        }'
                                    */
                                    // 你需要手动调用它(像下方这样)
                                    /* :error = '(()=>{
                                          return true
                                        })()'
                                    */
                                    // 不过, 这也仅调用了一次。 如果函数内部有响应式变量, 它并不会因触发响应式变更而重新被调用
                                    // 因此, 想要用函数的方式并且响应式生效
                                    // * 常用做法是在script中创建一个计算属性, 并在此处绑定它(而不是绑定某个函数)
                                    // * 还有一种做法是, 直接绑定对应的表达式(表达式中如果有响应式变量的变化, 表达式会被重新计算), 就像下面这样
                                  }
                                "
                                :error="
                                  playModeForUp.mode === 'single'
                                    ? selectedSoundsForUp.length > 1
                                      ? true
                                      : false
                                    : false
                                "
                                ref="upSoundSelectDom"
                                @update:model-value="upSoundSelectDom?.hidePopup()"
                              />
                              <div class="h-3">
                                <q-option-group
                                  dense
                                  v-model="upTypeGroup"
                                  :options="options"
                                  type="checkbox"
                                  class="absolute left-8"
                                >
                                  <template #label-0="props">
                                    <q-item-label>
                                      {{ props.label }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>本软件支持将源文件直接作为声音。</div>
                                          <div>前提是</div>
                                          <div>这个源文件本身就就可作为独立的无需裁剪的声音。</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                  <template v-slot:label-1="props">
                                    <q-item-label>
                                      {{ props.label }}
                                      <q-icon name="info" color="primary" class="p-l-4.5 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>默认仅从声音列表中选择。</div>
                                          <div>如有需求也可勾选其它受支持列表。</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                  <template v-slot:label-2="props">
                                    <q-item-label>
                                      {{ props.label }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>本软件支持将其它按键音作为声音。</div>
                                          <div>或者说</div>
                                          <div>本软件支持继承其它按键音的配置。</div>
                                          <div>⬇</div>
                                          <div>原则: 按下继承按下, 抬起继承抬起。</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                </q-option-group>
                              </div>
                            </q-card-section>
                            <q-card-actions align="right">
                              <q-btn flat label="Close" color="primary" v-close-popup />
                            </q-card-actions>
                          </q-card>
                        </q-dialog>
                      </div>
                    </q-card-section>
                    <q-card-actions align="right">
                      <q-btn
                        color="primary"
                        label="确认添加"
                        @click="
                          saveKeySoundConfig(
                            {
                              key: '',
                              name: keySoundName,
                              down: { mode: playModeForDown.mode, value: selectedSoundsForDown },
                              up: { mode: playModeForUp.mode, value: selectedSoundsForUp },
                            },
                            () => {
                              // 关闭对话框
                              createNewKeySound = !createNewKeySound;

                              // 重置表单变量
                              keySoundName = '新的按键音';
                              selectedSoundsForDown = [];
                              playModeForDown = { label: '随机', mode: 'random' };
                              selectedSoundsForUp = [];
                              playModeForUp = { label: '随机', mode: 'random' };
                            }
                          )
                        "
                      />
                      <q-btn flat label="Close" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>

              <div :class="['p-2 text-zinc-600']">或</div>

              <!-- -------------------------------------------------------------------------------编辑已有按键音 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  label="编辑已有按键音"
                  @click="
                    () => {
                      if (keySoundList.length === 0) {
                        q.notify({
                          type: 'warning',
                          message: '当前没有可编辑的按键音',
                          position: 'top',
                        });
                        return;
                      }
                      editExistingKeySound = true;
                    }
                  "
                >
                </q-btn>
                <q-dialog v-model="editExistingKeySound" backdrop-filter="invert(70%)">
                  <q-card :class="['min-w-[90%]']">
                    <q-card-section
                      class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                    >
                      编辑已有按键音
                    </q-card-section>
                    <q-card-section>
                      <q-select
                        outlined
                        stack-label
                        v-model="selectedKeySound"
                        :options="keySoundList"
                        label="选择要编辑的按键音"
                        :option-label="(item) => item.keySoundValue.name"
                        :option-value="(item) => item.keySoundKey"
                        dense
                      />
                    </q-card-section>
                    <!-- 以卡片的形式, 展示选择的按键音 -->
                    <q-card-section
                      :class="['flex flex-col m-t-3']"
                      v-if="selectedKeySound?.keySoundKey !== '' && selectedKeySound !== undefined"
                    >
                      <q-card :class="['flex flex-col px-3 pb-3']">
                        <q-badge
                          transparent
                          color="orange"
                          :label="selectedKeySound.keySoundValue.name"
                          :class="['absolute overflow-visible']"
                        />
                        <q-card-section :class="['p-b-1 mt-3']">
                          <q-input
                            outlined
                            stack-label
                            dense
                            v-model="selectedKeySound.keySoundValue.name"
                            label="按键音名称"
                            :placeholder="'新的按键音'"
                          />
                          <div class="flex flex-col mt-3">
                            <q-btn
                              :class="['bg-zinc-300 my-7 w-70% self-center']"
                              label="配置按下声音"
                              @click="edit_configureDownSound = true"
                            />
                            <q-dialog v-model="edit_configureDownSound" backdrop-filter="invert(70%)">
                              <q-card :class="['min-w-[80%]']">
                                <q-card-section class="row items-center q-pb-none text-h6">
                                  配置按下声音
                                </q-card-section>
                                <q-card-section>
                                  <q-select
                                    outlined
                                    stack-label
                                    v-model="selectedKeySound.keySoundValue.down.mode"
                                    :options="playModeOptions"
                                    label="选择播放模式"
                                    dense
                                  />
                                </q-card-section>
                                <q-card-section>
                                  <q-select
                                    outlined
                                    stack-label
                                    v-model="selectedKeySound.keySoundValue.down.value"
                                    :options="edit_downSoundList"
                                    :option-label="
                                      (item: any) => {
                                        if (item.type === 'audio_files') {
                                          const soundFile = soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id);
                                          return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFile?.name + soundFile?.type;
                                        }
                                        if (item.type === 'sounds') {

                                          const sound = soundList.find((sound) => sound.soundKey === item.value.soundKey);

                                          if (sound?.soundValue?.name !== '' && sound?.soundValue?.name !== undefined) {
                                            return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;'+ sound?.soundValue.name
                                          } else {
                                            return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                              soundFileList.find(
                                                (soundFile:any) =>
                                                  soundFile.sha256 === sound?.soundValue?.source_file_for_sound?.sha256 &&
                                                  soundFile.name_id === sound?.soundValue?.source_file_for_sound?.name_id
                                              )?.name +
                                              '     - ' +
                                              ' [' +
                                              sound?.soundValue?.cut?.start_time +
                                              ' ~ ' +
                                              sound?.soundValue?.cut?.end_time +
                                              ']'
                                            );
                                          }
                                        }
                                        if (item.type === 'key_sounds') {
                                          const keySound = keySoundList.find((keySound) => keySound.keySoundKey === item.value.keySoundKey);
                                          return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + keySound?.keySoundValue?.name;
                                        }
                                      }
                                    "
                                    :option-value="
                                      (item) => {
                                        /**
                                         * json中的存储格式分别是
                                         * {key:'audio_files', value:{sha256: string, name_id: string, type:string}}
                                         * {key:'sounds', value:string} // 此处value, 是soundKey
                                         * {key:'key_sounds', value:string} // 此处value, 是keySoundKey
                                         */
                                        if (item.type === 'audio_files') {
                                          return item.value.sha256 + item.value.name_id;
                                        }
                                        if (item.type === 'sounds') {
                                          return item.value.soundKey;
                                        }
                                        if (item.type === 'key_sounds') {
                                          return item.value.keySoundKey;
                                        }
                                      }
                                    "
                                    label="选择声音 (多选)"
                                    multiple
                                    use-chips
                                    dense
                                    :max-values="
                                      selectedKeySound.keySoundValue.down.mode.mode === 'single' ? 1 : Infinity
                                    "
                                    counter
                                    error-message="独立模式下, 至多选择一个声音"
                                    :error="
                                      selectedKeySound.keySoundValue.down.mode.mode === 'single' &&
                                      selectedKeySound.keySoundValue.down.value.length > 1
                                    "
                                    ref="edit_downSoundSelectDom"
                                    @update:model-value="edit_downSoundSelectDom?.hidePopup()"
                                  />
                                  <div class="h-3">
                                    <q-option-group
                                      dense
                                      v-model="edit_downTypeGroup"
                                      :options="options"
                                      type="checkbox"
                                      class="absolute left-8"
                                    >
                                      <template #label-0="props">
                                        <q-item-label>
                                          {{ props.label }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>本软件支持将源文件直接作为声音。</div>
                                              <div>前提是</div>
                                              <div>这个源文件本身就就可作为独立的无需裁剪的声音。</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-1="props">
                                        <q-item-label>
                                          {{ props.label }}
                                          <q-icon name="info" color="primary" class="p-l-4.5 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>默认仅从声音列表中选择。</div>
                                              <div>如有需求也可勾选其它受支持列表。</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-2="props">
                                        <q-item-label>
                                          {{ props.label }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>本软件支持将其它按键音作为声音。</div>
                                              <div>或者说</div>
                                              <div>本软件支持继承其它按键音的配置。</div>
                                              <div>⬇</div>
                                              <div>原则: 按下继承按下, 抬起继承抬起。</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                    </q-option-group>
                                  </div>
                                </q-card-section>
                                <q-card-actions align="right">
                                  <q-btn flat label="Close" color="primary" v-close-popup />
                                </q-card-actions>
                              </q-card>
                            </q-dialog>
                            <q-btn
                              :class="['bg-zinc-300 m-b-7 w-70% self-center']"
                              label="配置抬起声音"
                              @click="edit_configureUpSound = true"
                            />
                            <q-dialog v-model="edit_configureUpSound" backdrop-filter="invert(70%)">
                              <q-card :class="['min-w-[80%]']">
                                <q-card-section class="row items-center q-pb-none text-h6">
                                  配置抬起声音
                                </q-card-section>
                                <q-card-section>
                                  <q-select
                                    outlined
                                    stack-label
                                    v-model="selectedKeySound.keySoundValue.up.mode"
                                    :options="playModeOptions"
                                    label="选择播放模式"
                                    dense
                                  />
                                </q-card-section>
                                <q-card-section>
                                  <q-select
                                    outlined
                                    stack-label
                                    v-model="selectedKeySound.keySoundValue.up.value"
                                    :options="edit_upSoundList"
                                    :option-label="
                                      (item: any) => {
                                        if (item.type === 'audio_files') {
                                          const soundFile = soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id);
                                          return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFile?.name + soundFile?.type;
                                        }
                                        if (item.type === 'sounds') {

                                          const sound = soundList.find((sound) => sound.soundKey === item.value.soundKey);

                                          if (sound?.soundValue?.name !== '' && sound?.soundValue?.name !== undefined) {
                                            return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;'+ sound?.soundValue.name
                                          } else {
                                            return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                              soundFileList.find(
                                                (soundFile:any) =>
                                                  soundFile.sha256 === sound?.soundValue?.source_file_for_sound?.sha256 &&
                                                  soundFile.name_id === sound?.soundValue?.source_file_for_sound?.name_id
                                              )?.name +
                                              '     - ' +
                                              ' [' +
                                              sound?.soundValue?.cut?.start_time +
                                              ' ~ ' +
                                              sound?.soundValue?.cut?.end_time +
                                              ']'
                                            );
                                          }
                                        }
                                        if (item.type === 'key_sounds') {
                                          const keySound = keySoundList.find((keySound) => keySound.keySoundKey === item.value.keySoundKey);
                                          return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + keySound?.keySoundValue?.name;
                                        }
                                      }
                                    "
                                    :option-value="
                                      /**
                                       * 虽然json中的存储格式分别是
                                       * - {key:'audio_files', value:{sha256: string, name_id: string, type:string}}
                                       * - {key:'sounds', value:string} // 此处value, 是soundKey
                                       * - {key:'key_sounds', value:string} // 此处value, 是keySoundKey
                                       * 但是, 我们通过watch对当前组件的model做了变更, 使其类型提前由uuid转换成了相关对象
                                       * - 因此, 此处仍按照对应对象处理即可
                                       */
                                      (item) => {
                                        if (item.type === 'audio_files') {
                                          return item.value.sha256 + item.value.name_id;
                                        }
                                        if (item.type === 'sounds') {
                                          return item.value.soundKey;
                                        }
                                        if (item.type === 'key_sounds') {
                                          return item.value.keySoundKey;
                                        }
                                      }
                                    "
                                    label="选择声音 (多选)"
                                    multiple
                                    use-chips
                                    dense
                                    :max-values="
                                      selectedKeySound.keySoundValue.up.mode.mode === 'single' ? 1 : Infinity
                                    "
                                    counter
                                    error-message="独立模式下, 至多选择一个声音"
                                    :error="
                                      selectedKeySound.keySoundValue.up.mode.mode === 'single' &&
                                      selectedKeySound.keySoundValue.up.value.length > 1
                                    "
                                    ref="edit_upSoundSelectDom"
                                    @update:model-value="edit_upSoundSelectDom?.hidePopup()"
                                  />
                                  <div class="h-3">
                                    <q-option-group
                                      dense
                                      v-model="edit_upTypeGroup"
                                      :options="options"
                                      type="checkbox"
                                      class="absolute left-8"
                                    >
                                      <template #label-0="props">
                                        <q-item-label>
                                          {{ props.label }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>本软件支持将源文件直接作为声音。</div>
                                              <div>前提是</div>
                                              <div>这个源文件本身就就可作为独立的无需裁剪的声音。</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-1="props">
                                        <q-item-label>
                                          {{ props.label }}
                                          <q-icon name="info" color="primary" class="p-l-4.5 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>默认仅从声音列表中选择。</div>
                                              <div>如有需求也可勾选其它受支持列表。</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-2="props">
                                        <q-item-label>
                                          {{ props.label }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>本软件支持将其它按键音作为声音。</div>
                                              <div>或者说</div>
                                              <div>本软件支持继承其它按键音的配置。</div>
                                              <div>⬇</div>
                                              <div>原则: 按下继承按下, 抬起继承抬起。</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                    </q-option-group>
                                  </div>
                                </q-card-section>
                                <q-card-actions align="right">
                                  <q-btn flat label="Close" color="primary" v-close-popup />
                                </q-card-actions>
                              </q-card>
                            </q-dialog>
                          </div>
                        </q-card-section>
                        <q-card-section :class="['flex justify-center gap-4']">
                          <q-btn
                            dense
                            color="primary"
                            icon="save"
                            label="确认修改"
                            @click="
                              saveKeySoundConfig({
                                key: selectedKeySound.keySoundKey,
                                name: selectedKeySound.keySoundValue.name,
                                down: {
                                  mode: selectedKeySound.keySoundValue.down.mode.mode,
                                  value: selectedKeySound.keySoundValue.down.value,
                                },
                                up: {
                                  mode: selectedKeySound.keySoundValue.up.mode.mode,
                                  value: selectedKeySound.keySoundValue.up.value,
                                },
                              })
                            "
                          />
                          <q-btn
                            dense
                            color="negative"
                            icon="delete"
                            label="删除键音"
                            @click="
                              deleteKeySound({
                                keySoundKey: selectedKeySound.keySoundKey,
                                onSuccess: () => {
                                  selectedKeySound = undefined;
                                  q.notify({
                                    type: 'positive',
                                    position: 'top',
                                    message: '删除成功',
                                    timeout: 5,
                                  });
                                },
                              })
                            "
                          />
                        </q-card-section>
                      </q-card>
                    </q-card-section>
                    <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
                      <q-btn flat label="Close" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
            </div>
            <q-stepper-navigation>
              <q-btn @click="step = 4" color="primary" label="Continue" />
              <q-btn flat @click="step = 2" color="primary" label="Back" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>

          <!-- ------------------------------------------------------------------------制作按键声音的业务逻辑   end -->

          <!-- <q-step :name="4" title="对全局按键统一设置键音" icon="settings" :done="step > 3">
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
          </q-step> -->
          <q-step :name="4" title="按键联动声效" icon="settings" :done="step > 3">
            <div>
              为按键设置联动声效，按下或抬起按键时，自动播放预设声效。
              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                  <span>"音频文件"、"声音"、"键音" 均可作为声效与按键联动。<br /></span>
                </q-tooltip>
              </q-icon>
            </div>
            <div :class="['flex items-center m-t-2']">
              <span class="text-gray-500 mr-0.7">•</span>
              <span>
                是否启用内嵌测试音:
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
                    <span>KeyTone内嵌了测试用的全键声效。<br /></span>
                    <span>供用户检测软件是否正常运行。<br /></span>
                    <span>默认开启, 如有需要可手动关闭。<br /></span>
                    <span>播放优先级低于全/单键声效设置。<br /></span>
                  </q-tooltip>
                </q-icon>
              </span>
            </div>
            <div :class="['flex items-center ml-3']">
              <span class="text-gray-500 mr-1.5">•</span>
              <q-toggle v-model="isEnableEmbeddedTestSound.down" color="primary" label="按下测试音" dense />
            </div>
            <div :class="['flex items-center ml-3']">
              <span class="text-gray-500 mr-1.5">•</span>
              <q-toggle v-model="isEnableEmbeddedTestSound.up" color="primary" label="抬起测试音" dense />
            </div>
            <q-stepper-navigation>
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  label="全键声效设置"
                  @click="
                    () => {
                      showEveryKeyEffectDialog = true;
                    }
                  "
                >
                </q-btn>
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
                    <span>播放优先级高于内嵌测试音, 但低于单键声效。<br /></span>
                  </q-tooltip>
                </q-icon>
                <q-dialog v-model="showEveryKeyEffectDialog" backdrop-filter="invert(70%)">
                  <q-card>
                    <q-card-section
                      class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                    >
                      全键声效设置
                    </q-card-section>

                    <q-card-section class="q-pt-none">
                      <div class="text-subtitle1 q-mb-md">无需指定任何按键, 直接进行全局统一的声效设置。</div>
                      <!-- 这里添加声效选择等具体设置内容 -->
                    </q-card-section>
                    <q-card-section>
                      <!-- 选择全键按下声效的选项, 仅支持单选 -->
                      <div class="flex flex-row flex-nowrap items-center m-b-3 m-l-5">
                        <div class="flex flex-col space-y-4 w-7/8">
                          <q-select
                            outlined
                            stack-label
                            v-model="keyDownUnifiedSoundEffectSelect"
                            :options="keyUnifiedSoundEffectOptions"
                            :option-label="(item: any) => {
                                  if (item.type === 'audio_files') {
                                    return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.name + soundFileList.find( (soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.type;
                                  }
                                  if (item.type === 'sounds') {
                                    // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
                                    if (item.value.soundValue?.name !== '' && item.value.soundValue?.name !== undefined) {
                                      return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;'+ (soundList.find((sound) => sound.soundKey === item.value.soundKey)?.soundValue.name)
                                    } else {
                                      return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                        soundFileList.find(
                                          (soundFile:any) =>
                                            soundFile.sha256 === item.value.soundValue?.source_file_for_sound?.sha256 &&
                                            soundFile.name_id === item.value.soundValue?.source_file_for_sound?.name_id
                                        )?.name +
                                        '     - ' +
                                        ' [' +
                                        item.value.soundValue?.cut?.start_time +
                                        ' ~ ' +
                                        item.value.soundValue?.cut?.end_time +
                                        ']'
                                      );
                                    }
                                  }
                                  if (item.type === 'key_sounds') {
                                    return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (item.value.keySoundValue?.name);
                                  }
                                }"
                            :option-value="(item: any) => {
                                  // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                  if (item.type === 'audio_files'){
                                    return item.value.sha256 + item.value.name_id
                                  }
                                  if(item.type === 'sounds'){
                                    return item.value.soundKey
                                  }
                                  if(item.type === 'key_sounds'){
                                    return item.value.keySoundKey
                                  }
                                }"
                            label="设置全键按下声效"
                            use-chips
                            dense
                            @popup-hide="
                              () => {
                                if (
                                  // 为避免循环依赖, 此处作为锚定功能选择声效时的判断逻辑; 而删除声效时的判断逻辑, 在watch中书写。
                                  isShowUltimatePerfectionKeySoundAnchoring &&
                                  isAnchoringUltimatePerfectionKeySound &&
                                  // 这里的?是防止在勾选至臻键音的条件下, 仅打开选项菜单且未做任何选择就关闭时, mode的null值内 没有type字段引起报错。
                                  keyDownUnifiedSoundEffectSelect?.type === 'key_sounds'
                                ) {
                                  keyUpUnifiedSoundEffectSelect = keyDownUnifiedSoundEffectSelect;
                                }
                              }
                            "
                            class="max-w-full"
                          />
                          <!-- 选择全键抬起声效的选项, 仅支持单选 -->
                          <q-select
                            outlined
                            stack-label
                            v-model="keyUpUnifiedSoundEffectSelect"
                            :options="keyUnifiedSoundEffectOptions"
                            :option-label="(item: any) => {
                                  if (item.type === 'audio_files') {
                                    return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.name + soundFileList.find( (soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.type;
                                  }
                                  if (item.type === 'sounds') {
                                    // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
                                    if (item.value.soundValue?.name !== '' && item.value.soundValue?.name !== undefined) {
                                      return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;'+ (soundList.find((sound) => sound.soundKey === item.value.soundKey)?.soundValue.name)
                                    } else {
                                      return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                        soundFileList.find(
                                          (soundFile:any) =>
                                            soundFile.sha256 === item.value.soundValue?.source_file_for_sound?.sha256 &&
                                            soundFile.name_id === item.value.soundValue?.source_file_for_sound?.name_id
                                        )?.name +
                                        '     - ' +
                                        ' [' +
                                        item.value.soundValue?.cut?.start_time +
                                        ' ~ ' +
                                        item.value.soundValue?.cut?.end_time +
                                        ']'
                                      );
                                    }
                                  }
                                  if (item.type === 'key_sounds') {
                                    return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (item.value.keySoundValue?.name);
                                  }
                                }"
                            :option-value="(item: any) => {
                                  // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                  if (item.type === 'audio_files'){
                                    return item.value.sha256 + item.value.name_id
                                  }
                                  if(item.type === 'sounds'){
                                    return item.value.soundKey
                                  }
                                  if(item.type === 'key_sounds'){
                                    return item.value.keySoundKey
                                  }
                                }"
                            label="设置全键抬起声效"
                            use-chips
                            dense
                            @popup-hide="
                              () => {
                                // 为避免循环依赖, 此处作为锚定功能选择声效时的判断逻辑; 而删除声效时的判断逻辑, 在watch中书写。
                                if (
                                  isShowUltimatePerfectionKeySoundAnchoring &&
                                  isAnchoringUltimatePerfectionKeySound &&
                                  // 这里的?是防止在勾选至臻键音的条件下, 仅打开选项菜单且未做任何选择就关闭时, mode的null值内 没有type字段引起报错。
                                  keyUpUnifiedSoundEffectSelect?.type === 'key_sounds'
                                ) {
                                  keyDownUnifiedSoundEffectSelect = keyUpUnifiedSoundEffectSelect;
                                }
                              }
                            "
                            class="max-w-full"
                          />
                        </div>
                        <div class="flex justify-end -m-l-2">
                          <q-icon
                            @click="isAnchoringUltimatePerfectionKeySound = !isAnchoringUltimatePerfectionKeySound"
                            size="2.75rem"
                            v-if="isShowUltimatePerfectionKeySoundAnchoring"
                          >
                            <template v-if="isAnchoringUltimatePerfectionKeySound">
                              <!-- 锚定 -->
                              <svg
                                version="1.1"
                                xmlns="http://www.w3.org/2000/svg"
                                viewBox="0 0 27.838698814764143 45.55476707640558"
                                width="27.838698814764143"
                                height="45.55476707640558"
                              >
                                <g
                                  stroke-linecap="round"
                                  transform="translate(10.257971830020665 11.315110817092783) rotate(0 3.661377577361293 4.912529303576154)"
                                >
                                  <path
                                    d="M1.83 0 C3.09 0, 4.34 0, 5.49 0 M1.83 0 C2.76 0, 3.69 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.07, 7.32 6.3, 7.32 7.99 M7.32 1.83 C7.32 3.67, 7.32 5.51, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.52 9.83, 3.56 9.83, 1.83 9.83 M5.49 9.83 C4.3 9.83, 3.12 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.74, 0 5.48, 0 1.83 M0 7.99 C0 6.56, 0 5.12, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                    stroke="#1e1e1e"
                                    stroke-width="1"
                                    fill="none"
                                  ></path>
                                </g>
                                <g stroke-linecap="round">
                                  <g
                                    transform="translate(13.900885269327091 18.0363210366226) rotate(0 -0.0030445053470202765 2.3801245195212495)"
                                  >
                                    <path
                                      d="M0 0 C0 0.79, -0.01 3.97, -0.01 4.76 M0 0 C0 0.79, -0.01 3.97, -0.01 4.76"
                                      stroke="#1e1e1e"
                                      stroke-width="1"
                                      fill="none"
                                    ></path>
                                  </g>
                                </g>
                                <mask></mask>
                                <g
                                  stroke-linecap="round"
                                  transform="translate(10.257971830020665 24.438144439265034) rotate(0 3.661377577361293 4.912529303576154)"
                                >
                                  <path
                                    d="M1.83 0 C2.75 0, 3.67 0, 5.49 0 M1.83 0 C2.98 0, 4.12 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.09, 7.32 6.36, 7.32 7.99 M7.32 1.83 C7.32 4.06, 7.32 6.28, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.11 9.83, 2.73 9.83, 1.83 9.83 M5.49 9.83 C4.08 9.83, 2.67 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.4, 0 4.8, 0 1.83 M0 7.99 C0 6.46, 0 4.92, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                    stroke="#1e1e1e"
                                    stroke-width="1"
                                    fill="none"
                                  ></path>
                                </g>
                                <g stroke-linecap="round">
                                  <g
                                    transform="translate(13.900885269326977 27.541992826887526) rotate(0 -0.0030445053470202765 -2.3801245195212495)"
                                  >
                                    <path
                                      d="M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76 M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76"
                                      stroke="#1e1e1e"
                                      stroke-width="1"
                                      fill="none"
                                    ></path>
                                  </g>
                                </g>
                                <mask></mask>
                                <g
                                  stroke-linecap="round"
                                  transform="translate(10 10) rotate(0 3.9193494073820716 12.77738353820279)"
                                >
                                  <path
                                    d="M1.96 0 C3.5 0, 5.03 0, 5.88 0 M1.96 0 C3.41 0, 4.86 0, 5.88 0 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M7.84 1.96 C7.84 8.03, 7.84 14.11, 7.84 23.6 M7.84 1.96 C7.84 7.33, 7.84 12.7, 7.84 23.6 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M5.88 25.55 C4.57 25.55, 3.25 25.55, 1.96 25.55 M5.88 25.55 C4.56 25.55, 3.25 25.55, 1.96 25.55 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M0 23.6 C0 16.78, 0 9.97, 0 1.96 M0 23.6 C0 15.66, 0 7.72, 0 1.96 M0 1.96 C0 0.65, 0.65 0, 1.96 0 M0 1.96 C0 0.65, 0.65 0, 1.96 0"
                                    stroke="transparent"
                                    stroke-width="1"
                                    fill="none"
                                  ></path>
                                </g>
                              </svg>
                            </template>
                            <template v-else>
                              <!-- 锚定解除 -->
                              <svg
                                version="1.1"
                                xmlns="http://www.w3.org/2000/svg"
                                viewBox="0 0 27.838698814764143 45.55476707640558"
                                width="27.838698814764143"
                                height="45.55476707640558"
                              >
                                <g
                                  stroke-linecap="round"
                                  transform="translate(10.257971830020779 10.375546926705098) rotate(0 3.661377577361293 4.912529303576154)"
                                >
                                  <path
                                    d="M1.83 0 C3.09 0, 4.34 0, 5.49 0 M1.83 0 C2.76 0, 3.69 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.07, 7.32 6.3, 7.32 7.99 M7.32 1.83 C7.32 3.67, 7.32 5.51, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.52 9.83, 3.56 9.83, 1.83 9.83 M5.49 9.83 C4.3 9.83, 3.12 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.74, 0 5.48, 0 1.83 M0 7.99 C0 6.56, 0 5.12, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                    stroke="#1e1e1e"
                                    stroke-width="1"
                                    fill="none"
                                  ></path>
                                </g>
                                <g stroke-linecap="round">
                                  <g
                                    transform="translate(13.900885269327205 17.096757146234918) rotate(0 -0.0030445053470202765 2.3801245195212513)"
                                  >
                                    <path
                                      d="M0 0 C0 0.79, -0.01 3.97, -0.01 4.76 M0 0 C0 0.79, -0.01 3.97, -0.01 4.76"
                                      stroke="#1e1e1e"
                                      stroke-width="1"
                                      fill="none"
                                    ></path>
                                  </g>
                                </g>
                                <mask></mask>
                                <g
                                  stroke-linecap="round"
                                  transform="translate(10.257971830020779 25.379515181404656) rotate(0 3.661377577361293 4.912529303576154)"
                                >
                                  <path
                                    d="M1.83 0 C2.75 0, 3.67 0, 5.49 0 M1.83 0 C2.98 0, 4.12 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.09, 7.32 6.36, 7.32 7.99 M7.32 1.83 C7.32 4.06, 7.32 6.28, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.11 9.83, 2.73 9.83, 1.83 9.83 M5.49 9.83 C4.08 9.83, 2.67 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.4, 0 4.8, 0 1.83 M0 7.99 C0 6.46, 0 4.92, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                    stroke="#1e1e1e"
                                    stroke-width="1"
                                    fill="none"
                                  ></path>
                                </g>
                                <g stroke-linecap="round">
                                  <g
                                    transform="translate(13.900885269327091 28.483363569027148) rotate(0 -0.0030445053470202765 -2.3801245195212495)"
                                  >
                                    <path
                                      d="M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76 M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76"
                                      stroke="#1e1e1e"
                                      stroke-width="1"
                                      fill="none"
                                    ></path>
                                  </g>
                                </g>
                                <mask></mask>
                                <g
                                  stroke-linecap="round"
                                  transform="translate(10 10) rotate(0 3.9193494073820716 12.777383538202791)"
                                >
                                  <path
                                    d="M1.96 0 C3.5 0, 5.03 0, 5.88 0 M1.96 0 C3.41 0, 4.86 0, 5.88 0 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M7.84 1.96 C7.84 8.03, 7.84 14.11, 7.84 23.6 M7.84 1.96 C7.84 7.33, 7.84 12.7, 7.84 23.6 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M5.88 25.55 C4.57 25.55, 3.25 25.55, 1.96 25.55 M5.88 25.55 C4.56 25.55, 3.25 25.55, 1.96 25.55 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M0 23.6 C0 16.78, 0 9.97, 0 1.96 M0 23.6 C0 15.66, 0 7.72, 0 1.96 M0 1.96 C0 0.65, 0.65 0, 1.96 0 M0 1.96 C0 0.65, 0.65 0, 1.96 0"
                                    stroke="transparent"
                                    stroke-width="1"
                                    fill="none"
                                  ></path>
                                </g>
                              </svg>
                            </template>
                            <q-tooltip
                              :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                            >
                              <span class="text-sm">至臻键音</span>
                              <span class="text-sm" v-if="isAnchoringUltimatePerfectionKeySound"> (已锚定)<br /></span>
                              <span class="text-sm" v-else> (已解除锚定)<br /></span>
                              <span>此锚定仅作用于在设置声效为至臻键音时。</span>
                            </q-tooltip>
                          </q-icon>
                        </div>
                      </div>
                      <div class="h-16 m-l-9">
                        <q-option-group dense v-model="unifiedTypeGroup" :options="options" type="checkbox">
                          <template #label-0="props">
                            <q-item-label>
                              {{ props.label }}
                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                <q-tooltip
                                  :class="[
                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                  ]"
                                >
                                  <div>本软件支持将音频源文件直接设置为声效。</div>
                                </q-tooltip>
                              </q-icon>
                            </q-item-label>
                          </template>
                          <template v-slot:label-1="props">
                            <q-item-label>
                              {{ props.label }}
                              <q-icon name="info" color="primary" class="p-l-4.5 m-b-0.5">
                                <q-tooltip
                                  :class="[
                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                  ]"
                                >
                                  <div>本软件支持将裁剪定义好的声音设置为声效。</div>
                                </q-tooltip>
                              </q-icon>
                            </q-item-label>
                          </template>
                          <template v-slot:label-2="props">
                            <q-item-label>
                              {{ props.label }}
                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                <q-tooltip
                                  :class="[
                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                  ]"
                                >
                                  <div>本软件推荐用户使用至臻键音, 以设置更为丰富的声效。</div>
                                  <div>若本选项勾选, 则自动提供可用于至臻键音的锚定功能。</div>
                                  <div>此功能旨在方便您操作, 您也可解除锚定, 更自由的定义声效。</div>
                                </q-tooltip>
                              </q-icon>
                            </q-item-label>
                          </template>
                        </q-option-group>
                      </div>
                    </q-card-section>
                    <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
                      <q-btn
                        flat
                        label="确定"
                        color="primary"
                        v-close-popup
                        @click="
                          saveUnifiedSoundEffectConfig(
                            {
                              down: keyDownUnifiedSoundEffectSelect,
                              up: keyUpUnifiedSoundEffectSelect,
                            },
                            () => {
                              q.notify({
                                type: 'positive',
                                position: 'top',
                                message: '全键声效配置成功',
                                timeout: 2000,
                              });
                            }
                          )
                        "
                      />
                      <q-btn flat label="取消" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
              <div :class="['p-2 text-zinc-600']">或</div>
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  label="单键声效设置"
                  @click="
                    () => {
                      showSingleKeyEffectDialog = true;
                    }
                  "
                >
                </q-btn>
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
                    <span>播放优先级高于全键声效和内嵌测试音。<br /></span>
                  </q-tooltip>
                </q-icon>
                <q-dialog v-model="showSingleKeyEffectDialog" backdrop-filter="invert(70%)">
                  <q-card>
                    <q-card-section
                      class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                    >
                      单键声效设置
                    </q-card-section>

                    <q-card-section class="q-pt-none pb-0">
                      <div class="text-subtitle1 q-mb-md leading-tight m-t-1.5">
                        指定单个或多个按键, 以进行局部独立的声效设置。
                      </div>
                      <!-- 这里添加按键选择和声效设置等具体内容 -->
                      <div class="flex flex-row items-center gap-2 mb-2 ml-2">
                        <q-btn
                          flat
                          round
                          color="primary"
                          icon="add"
                          @click="
                            () => {
                              isShowAddOrSettingSingleKeyEffectDialog = true;
                            }
                          "
                        />
                        快速增设单键声效
                        <q-dialog
                          :no-esc-dismiss="isRecordingSingleKeys && isGetsFocused"
                          v-model="isShowAddOrSettingSingleKeyEffectDialog"
                          backdrop-filter="invert(70%)"
                        >
                          <q-card>
                            <q-card-section
                              class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                            >
                              添加单键声效
                            </q-card-section>

                            <q-card-section class="q-pt-none">
                              <div class="text-subtitle1 q-mb-md leading-tight m-t-1.5">请选择按键和对应的声效设置</div>
                              <div class="flex flex-col gap-4">
                                <div class="flex flex-row items-center gap-2 w-full">
                                  <!--
                                    // 此选择组件增加录制功能且默认启用。(停止录制的应用场景全配列键盘的用户可能用不到, 但当某些用户使用的键盘不是全配列, 无法录制一些常用按键时(如小键盘数字区)会非常有用。)
                                    //   ↓    - completed(已完成)    只不过方法不是太优雅而已 -> 采用了其中提到的比较啰嗦的其余方式。
                                    // TODO: 在使用录制功能录制按键的过程中会存在一些干扰录制过程的功能键, 最影响的就数`BACKSPACE`键了。
                                             但是说实话, 目前我没有很好的解决此问题的方式, 因为即使阻止@keydown事件的按键默认行为也无效。
                                             对于 input 组件上删除事件监听的做法, 由于无法针对性删除, 因此也无济于事。
                                             而其余的方式, 有点太啰嗦了, 并且此bug对用户体验也是有利有弊的(弊大于利吧),索性就暂时不适配了。
                                             何时适配? 等待quasar的更新中提供相关的选项后, 再适配此TODO。
                                    //
                                    // 组件的关键使用注释
                                    // 单纯而new-value-mode满足不了我的需求, 我需要@new-value事件来更进一步的使用。
                                    new-value-mode="add-unique"  // 一开始想通过此属性来实现录制功能, 不过之后使用其它方案了
                                    // 为了在 录制单键 功能启用时, 阻止正常输入的内容。(使用此方式, 可以便捷的阻止输入)
                                    :maxlength="isRecordingSingleKeys ? 0 : Infinity"
                                  -->
                                  <!-- 按键选择 start -->
                                  <q-select
                                    label="选择单键(数量任意)"
                                    ref="singleKeysSelectRef"
                                    v-model="selectedSingleKeys"
                                    :options="filterOptions"
                                    :virtual-scroll-slice-size="
                                      (() => {
                                        // TIPS: 此处如此写, 只是为了添加注释。因此没使用正常的双引号内直接写9999的方式。
                                        // TIPS: select组件默认启用了虚拟滚动行为, 初始值是10, 因此每次只会渲染10个, 超出就会触发虚拟滚动。
                                        // TIPS: 我不喜欢quasar的虚拟滚动, 但quasar不支持关闭。 因此将此值设置的足够大来满足禁用虚拟滚动的需求。
                                        // TIPS: 就算我需要使用虚拟滚动, 也会使用专门的虚拟滚动库, 而不是quasar的虚拟滚动(因为quasar的虚拟滚动影响显示)。
                                        //       不过似乎不现实, 因为即使自定义菜单时, 似乎也被要求加上虚拟滚动类 -> https://quasar.dev/vue-components/select#customizing-menu-options:~:text=Customizing%20menu%20options,the%20first%20one
                                        return 9999;
                                      })()
                                    "
                                    dense
                                    filled
                                    hide-dropdown-icon
                                    multiple
                                    outlined
                                    stack-label
                                    :placeholder="isRecordingSingleKeys ? '在此键入录制单键' : '在此键入搜寻单键'"
                                    use-input
                                    use-chips
                                    class="flex-1"
                                    @focus="
                                      () => {
                                        isGetsFocused = true;
                                      }
                                    "
                                    @blur="
                                      () => {
                                        isGetsFocused = false;
                                      }
                                    "
                                    :maxlength="isRecordingSingleKeys ? 0 : Infinity"
                                    @keydown="preventDefaultKeyBehaviorWhenRecording"
                                    :option-label="
                                      (option) => {
                                        return (
                                          keyEvent_store.dikCodeToName.get(option) || // 优先从项目作者制作的特定于 dik码 与 按键名称 的映射码表中取按键名。
                                          (dikCodeToName_custom.get(option) // 其次从项目作者所编写的 dik码 与 前端keycode/key名 的自动映射逻辑生成的码表中取按键名。(这一步后续可能移除: 因为这一步只是作者为快速制作第一步中的码表, 才编写的特定程序, 以尽可能轻松定义特定的按键名称--毕竟定义每一个按键的名称, 且不能重复的工作量不小, 因此想借助前端已有的keycode或key的定义来加快进度。)
                                            ? 'Temp-{' + dikCodeToName_custom.get(option) + '}' // 防止影响优先级流程
                                            : '') || // 中间这个需要加上括号, 不然 '||' 运算符的优先级大于 '?'
                                          'Dik-{' + option + '}'
                                        ); // 最后, 若按键名仍无法识别, 将按照此处指定的固定格式, 展示出无法失败的Dik码值。(其实有了这一步, 第二步可以删除的, 毕竟它的存在有可能影响测试, 故对其增加大括号以避免造成影响。)
                                      }
                                    "
                                    @filter="
                                      (inputValue, doneFn) => {
                                        if (inputValue === '') {
                                          doneFn(() => {
                                            filterOptions = keyOptions;
                                          });
                                          return;
                                        }

                                        doneFn(() => {
                                          const inputValueLowerCase = inputValue.toLowerCase();
                                          filterOptions = keyOptions.filter((item) => {
                                            // 首先过滤系统 TIPS: 由于结果0也是符合预期的(即第一个字符), 因此不能使用||来处理, 否则会造成第一个字符无法被识别。
                                            let ifre = keyEvent_store.dikCodeToName
                                              .get(item)
                                              ?.toLowerCase()
                                              ?.indexOf(inputValueLowerCase);

                                            if (ifre !== undefined && ifre > -1) {
                                              return true;
                                            }

                                            // 其次过滤temp(由于名称中字符串'temp_'是显示时候加的, 因此不影响实际过滤行为)   TIPS: 由于结果0也是符合预期的(即第一个字符), 因此不能使用||来处理, 否则会造成第一个字符无法被识别。
                                            ifre = dikCodeToName_custom
                                              .get(item)
                                              ?.toLowerCase()
                                              ?.indexOf(inputValueLowerCase);
                                            if (ifre !== undefined && ifre > -1) {
                                              // TIPS: 举例子, 如系统a在temp中是 keyA, 且temp中所有的字母键都有key三个字符, 这就造成了当按下 k 或 e 或 y 时, 被误触发。
                                              //       而此处的if, 就是防止被误触发的。(因为只要这个dik码在主表中存在并过滤过, 就没必要过滤temp表了。因为我们是以主表为主的, 这样判断可避免误触发。)
                                              //       只要在主表中不存在即get结果为undefined时, 我们才启用temp表的过滤。
                                              if (keyEvent_store.dikCodeToName.get(item) === undefined) {
                                                return true;
                                              }
                                            }
                                            return false;
                                          });
                                        });
                                      }
                                    "
                                  >
                                    <template v-slot:append>
                                      <q-btn
                                        dense
                                        flat
                                        :color="isRecordingSingleKeys ? 'primary' : ''"
                                        icon="keyboard"
                                        @click="isRecordingSingleKeys = !isRecordingSingleKeys"
                                      >
                                        <q-tooltip>
                                          {{ isRecordingSingleKeys ? '停止录制单键' : '开始录制单键' }}
                                        </q-tooltip>
                                      </q-btn>
                                    </template>
                                  </q-select>
                                  <!-- 按键选择 end -->

                                  <!-- 声效选择 start -->
                                  <div class="flex flex-row items-center justify-center gap-9 w-[88%]">
                                    <q-checkbox dense v-model="isDownSoundEffectSelectEnabled" label="按下声效" />
                                    <q-checkbox dense v-model="isUpSoundEffectSelectEnabled" label="抬起声效" />
                                  </div>
                                  <div class="w-full">
                                    <q-card-section>
                                      <div class="flex flex-row flex-nowrap items-center m-b-3 m-l-5">
                                        <div class="flex flex-col space-y-4 w-6.5/8">
                                          <!-- 选择单键按下声效的选项, 仅支持单选 -->
                                          <q-select
                                            v-show="isDownSoundEffectSelectEnabled"
                                            outlined
                                            stack-label
                                            v-model="keyDownSingleKeySoundEffectSelect"
                                            :options="keySingleKeySoundEffectOptions"
                                            :option-label="(item: any) => {
                                              if (item.type === 'audio_files') {
                                                return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.name + soundFileList.find( (soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.type;
                                              }
                                              if (item.type === 'sounds') {
                                                // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
                                                if (item.value.soundValue?.name !== '' && item.value.soundValue?.name !== undefined) {
                                                  return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;'+ (soundList.find((sound) => sound.soundKey === item.value.soundKey)?.soundValue.name)
                                                } else {
                                                  return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                                    soundFileList.find(
                                                      (soundFile:any) =>
                                                        soundFile.sha256 === item.value.soundValue?.source_file_for_sound?.sha256 &&
                                                        soundFile.name_id === item.value.soundValue?.source_file_for_sound?.name_id
                                                    )?.name +
                                                    '     - ' +
                                                    ' [' +
                                                    item.value.soundValue?.cut?.start_time +
                                                    ' ~ ' +
                                                    item.value.soundValue?.cut?.end_time +
                                                    ']'
                                                  );
                                                }
                                              }
                                              if (item.type === 'key_sounds') {
                                                return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (item.value.keySoundValue?.name);
                                              }
                                            }"
                                            :option-value="(item: any) => {
                                              // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                              if (item.type === 'audio_files'){
                                                return item.value.sha256 + item.value.name_id
                                              }
                                              if(item.type === 'sounds'){
                                                return item.value.soundKey
                                              }
                                              if(item.type === 'key_sounds'){
                                                return item.value.keySoundKey
                                              }
                                            }"
                                            label="设置按下声效"
                                            use-chips
                                            dense
                                            @popup-hide="
                                              () => {
                                                if (
                                                  // 为避免循环依赖, 此处作为锚定功能选择声效时的判断逻辑; 而删除声效时的判断逻辑, 在watch中书写。
                                                  isShowUltimatePerfectionKeySoundAnchoring_singleKey &&
                                                  isAnchoringUltimatePerfectionKeySound_singleKey &&
                                                  // 这里的?是防止在勾选至臻键音的条件下, 仅打开选项菜单且未做任何选择就关闭时, mode的null值内 没有type字段引起报错。
                                                  keyDownSingleKeySoundEffectSelect?.type === 'key_sounds'
                                                ) {
                                                  keyUpSingleKeySoundEffectSelect = keyDownSingleKeySoundEffectSelect;
                                                }
                                              }
                                            "
                                            class="max-w-full"
                                          />
                                          <!-- 选择单键抬起声效的选项, 仅支持单选 -->
                                          <q-select
                                            v-show="isUpSoundEffectSelectEnabled"
                                            outlined
                                            stack-label
                                            v-model="keyUpSingleKeySoundEffectSelect"
                                            :options="keySingleKeySoundEffectOptions"
                                            :option-label="(item: any) => {
                                              if (item.type === 'audio_files') {
                                                return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.name + soundFileList.find( (soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.type;
                                              }
                                              if (item.type === 'sounds') {
                                                // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
                                                if (item.value.soundValue?.name !== '' && item.value.soundValue?.name !== undefined) {
                                                  return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;'+ (soundList.find((sound) => sound.soundKey === item.value.soundKey)?.soundValue.name)
                                                } else {
                                                  return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                                    soundFileList.find(
                                                      (soundFile:any) =>
                                                        soundFile.sha256 === item.value.soundValue?.source_file_for_sound?.sha256 &&
                                                        soundFile.name_id === item.value.soundValue?.source_file_for_sound?.name_id
                                                    )?.name +
                                                    '     - ' +
                                                    ' [' +
                                                    item.value.soundValue?.cut?.start_time +
                                                    ' ~ ' +
                                                    item.value.soundValue?.cut?.end_time +
                                                    ']'
                                                  );
                                                }
                                              }
                                              if (item.type === 'key_sounds') {
                                                return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (item.value.keySoundValue?.name);
                                              }
                                            }"
                                            :option-value="(item: any) => {
                                              // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                              if (item.type === 'audio_files'){
                                                return item.value.sha256 + item.value.name_id
                                              }
                                              if(item.type === 'sounds'){
                                                return item.value.soundKey
                                              }
                                              if(item.type === 'key_sounds'){
                                                return item.value.keySoundKey
                                              }
                                            }"
                                            label="设置抬起声效"
                                            use-chips
                                            dense
                                            @popup-hide="
                                              () => {
                                                // 为避免循环依赖, 此处作为锚定功能选择声效时的判断逻辑; 而删除声效时的判断逻辑, 在watch中书写。
                                                if (
                                                  isShowUltimatePerfectionKeySoundAnchoring_singleKey &&
                                                  isAnchoringUltimatePerfectionKeySound_singleKey &&
                                                  // 这里的?是防止在勾选至臻键音的条件下, 仅打开选项菜单且未做任何选择就关闭时, mode的null值内 没有type字段引起报错。
                                                  keyUpSingleKeySoundEffectSelect?.type === 'key_sounds'
                                                ) {
                                                  keyDownSingleKeySoundEffectSelect = keyUpSingleKeySoundEffectSelect;
                                                }
                                              }
                                            "
                                            class="max-w-full"
                                          />
                                        </div>
                                        <div
                                          v-show="isDownSoundEffectSelectEnabled && isUpSoundEffectSelectEnabled"
                                          class="flex justify-end -m-l-2"
                                        >
                                          <q-icon
                                            @click="
                                              isAnchoringUltimatePerfectionKeySound_singleKey =
                                                !isAnchoringUltimatePerfectionKeySound_singleKey
                                            "
                                            size="2.75rem"
                                            v-if="isShowUltimatePerfectionKeySoundAnchoring_singleKey"
                                          >
                                            <template v-if="isAnchoringUltimatePerfectionKeySound_singleKey">
                                              <!-- 锚定 -->
                                              <svg
                                                version="1.1"
                                                xmlns="http://www.w3.org/2000/svg"
                                                viewBox="0 0 27.838698814764143 45.55476707640558"
                                                width="27.838698814764143"
                                                height="45.55476707640558"
                                              >
                                                <g
                                                  stroke-linecap="round"
                                                  transform="translate(10.257971830020665 11.315110817092783) rotate(0 3.661377577361293 4.912529303576154)"
                                                >
                                                  <path
                                                    d="M1.83 0 C3.09 0, 4.34 0, 5.49 0 M1.83 0 C2.76 0, 3.69 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.07, 7.32 6.3, 7.32 7.99 M7.32 1.83 C7.32 3.67, 7.32 5.51, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.52 9.83, 3.56 9.83, 1.83 9.83 M5.49 9.83 C4.3 9.83, 3.12 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.74, 0 5.48, 0 1.83 M0 7.99 C0 6.56, 0 5.12, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                                    stroke="#1e1e1e"
                                                    stroke-width="1"
                                                    fill="none"
                                                  ></path>
                                                </g>
                                                <g stroke-linecap="round">
                                                  <g
                                                    transform="translate(13.900885269327091 18.0363210366226) rotate(0 -0.0030445053470202765 2.3801245195212495)"
                                                  >
                                                    <path
                                                      d="M0 0 C0 0.79, -0.01 3.97, -0.01 4.76 M0 0 C0 0.79, -0.01 3.97, -0.01 4.76"
                                                      stroke="#1e1e1e"
                                                      stroke-width="1"
                                                      fill="none"
                                                    ></path>
                                                  </g>
                                                </g>
                                                <mask></mask>
                                                <g
                                                  stroke-linecap="round"
                                                  transform="translate(10.257971830020665 24.438144439265034) rotate(0 3.661377577361293 4.912529303576154)"
                                                >
                                                  <path
                                                    d="M1.83 0 C2.75 0, 3.67 0, 5.49 0 M1.83 0 C2.98 0, 4.12 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.09, 7.32 6.36, 7.32 7.99 M7.32 1.83 C7.32 4.06, 7.32 6.28, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.11 9.83, 2.73 9.83, 1.83 9.83 M5.49 9.83 C4.08 9.83, 2.67 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.4, 0 4.8, 0 1.83 M0 7.99 C0 6.46, 0 4.92, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                                    stroke="#1e1e1e"
                                                    stroke-width="1"
                                                    fill="none"
                                                  ></path>
                                                </g>
                                                <g stroke-linecap="round">
                                                  <g
                                                    transform="translate(13.900885269326977 27.541992826887526) rotate(0 -0.0030445053470202765 -2.3801245195212495)"
                                                  >
                                                    <path
                                                      d="M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76 M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76"
                                                      stroke="#1e1e1e"
                                                      stroke-width="1"
                                                      fill="none"
                                                    ></path>
                                                  </g>
                                                </g>
                                                <mask></mask>
                                                <g
                                                  stroke-linecap="round"
                                                  transform="translate(10 10) rotate(0 3.9193494073820716 12.77738353820279)"
                                                >
                                                  <path
                                                    d="M1.96 0 C3.5 0, 5.03 0, 5.88 0 M1.96 0 C3.41 0, 4.86 0, 5.88 0 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M7.84 1.96 C7.84 8.03, 7.84 14.11, 7.84 23.6 M7.84 1.96 C7.84 7.33, 7.84 12.7, 7.84 23.6 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M5.88 25.55 C4.57 25.55, 3.25 25.55, 1.96 25.55 M5.88 25.55 C4.56 25.55, 3.25 25.55, 1.96 25.55 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M0 23.6 C0 16.78, 0 9.97, 0 1.96 M0 23.6 C0 15.66, 0 7.72, 0 1.96 M0 1.96 C0 0.65, 0.65 0, 1.96 0 M0 1.96 C0 0.65, 0.65 0, 1.96 0"
                                                    stroke="transparent"
                                                    stroke-width="1"
                                                    fill="none"
                                                  ></path>
                                                </g>
                                              </svg>
                                            </template>
                                            <template v-else>
                                              <!-- 锚定解除 -->
                                              <svg
                                                version="1.1"
                                                xmlns="http://www.w3.org/2000/svg"
                                                viewBox="0 0 27.838698814764143 45.55476707640558"
                                                width="27.838698814764143"
                                                height="45.55476707640558"
                                              >
                                                <g
                                                  stroke-linecap="round"
                                                  transform="translate(10.257971830020779 10.375546926705098) rotate(0 3.661377577361293 4.912529303576154)"
                                                >
                                                  <path
                                                    d="M1.83 0 C3.09 0, 4.34 0, 5.49 0 M1.83 0 C2.76 0, 3.69 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.07, 7.32 6.3, 7.32 7.99 M7.32 1.83 C7.32 3.67, 7.32 5.51, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.52 9.83, 3.56 9.83, 1.83 9.83 M5.49 9.83 C4.3 9.83, 3.12 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.74, 0 5.48, 0 1.83 M0 7.99 C0 6.56, 0 5.12, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                                    stroke="#1e1e1e"
                                                    stroke-width="1"
                                                    fill="none"
                                                  ></path>
                                                </g>
                                                <g stroke-linecap="round">
                                                  <g
                                                    transform="translate(13.900885269327205 17.096757146234918) rotate(0 -0.0030445053470202765 2.3801245195212513)"
                                                  >
                                                    <path
                                                      d="M0 0 C0 0.79, -0.01 3.97, -0.01 4.76 M0 0 C0 0.79, -0.01 3.97, -0.01 4.76"
                                                      stroke="#1e1e1e"
                                                      stroke-width="1"
                                                      fill="none"
                                                    ></path>
                                                  </g>
                                                </g>
                                                <mask></mask>
                                                <g
                                                  stroke-linecap="round"
                                                  transform="translate(10.257971830020779 25.379515181404656) rotate(0 3.661377577361293 4.912529303576154)"
                                                >
                                                  <path
                                                    d="M1.83 0 C2.75 0, 3.67 0, 5.49 0 M1.83 0 C2.98 0, 4.12 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.09, 7.32 6.36, 7.32 7.99 M7.32 1.83 C7.32 4.06, 7.32 6.28, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.11 9.83, 2.73 9.83, 1.83 9.83 M5.49 9.83 C4.08 9.83, 2.67 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.4, 0 4.8, 0 1.83 M0 7.99 C0 6.46, 0 4.92, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                                    stroke="#1e1e1e"
                                                    stroke-width="1"
                                                    fill="none"
                                                  ></path>
                                                </g>
                                                <g stroke-linecap="round">
                                                  <g
                                                    transform="translate(13.900885269327091 28.483363569027148) rotate(0 -0.0030445053470202765 -2.3801245195212495)"
                                                  >
                                                    <path
                                                      d="M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76 M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76"
                                                      stroke="#1e1e1e"
                                                      stroke-width="1"
                                                      fill="none"
                                                    ></path>
                                                  </g>
                                                </g>
                                                <mask></mask>
                                                <g
                                                  stroke-linecap="round"
                                                  transform="translate(10 10) rotate(0 3.9193494073820716 12.777383538202791)"
                                                >
                                                  <path
                                                    d="M1.96 0 C3.5 0, 5.03 0, 5.88 0 M1.96 0 C3.41 0, 4.86 0, 5.88 0 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M7.84 1.96 C7.84 8.03, 7.84 14.11, 7.84 23.6 M7.84 1.96 C7.84 7.33, 7.84 12.7, 7.84 23.6 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M5.88 25.55 C4.57 25.55, 3.25 25.55, 1.96 25.55 M5.88 25.55 C4.56 25.55, 3.25 25.55, 1.96 25.55 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M0 23.6 C0 16.78, 0 9.97, 0 1.96 M0 23.6 C0 15.66, 0 7.72, 0 1.96 M0 1.96 C0 0.65, 0.65 0, 1.96 0 M0 1.96 C0 0.65, 0.65 0, 1.96 0"
                                                    stroke="transparent"
                                                    stroke-width="1"
                                                    fill="none"
                                                  ></path>
                                                </g>
                                              </svg>
                                            </template>
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <span class="text-sm">至臻键音</span>
                                              <span
                                                class="text-sm"
                                                v-if="isAnchoringUltimatePerfectionKeySound_singleKey"
                                              >
                                                (已锚定)<br
                                              /></span>
                                              <span class="text-sm" v-else> (已解除锚定)<br /></span>
                                              <span>此锚定仅作用于在设置声效为至臻键音时。</span>
                                            </q-tooltip>
                                          </q-icon>
                                        </div>
                                      </div>
                                      <div
                                        class="h-16 m-l-9"
                                        v-show="isDownSoundEffectSelectEnabled || isUpSoundEffectSelectEnabled"
                                      >
                                        <q-option-group
                                          dense
                                          v-model="singleKeyTypeGroup"
                                          :options="options"
                                          type="checkbox"
                                        >
                                          <template #label-0="props">
                                            <q-item-label>
                                              {{ props.label }}
                                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                                <q-tooltip
                                                  :class="[
                                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                                  ]"
                                                >
                                                  <div>本软件支持将音频源文件直接设置为声效。</div>
                                                </q-tooltip>
                                              </q-icon>
                                            </q-item-label>
                                          </template>
                                          <template v-slot:label-1="props">
                                            <q-item-label>
                                              {{ props.label }}
                                              <q-icon name="info" color="primary" class="p-l-4.5 m-b-0.5">
                                                <q-tooltip
                                                  :class="[
                                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                                  ]"
                                                >
                                                  <div>本软件支持将裁剪定义好的声音设置为声效。</div>
                                                </q-tooltip>
                                              </q-icon>
                                            </q-item-label>
                                          </template>
                                          <template v-slot:label-2="props">
                                            <q-item-label>
                                              {{ props.label }}
                                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                                <q-tooltip
                                                  :class="[
                                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                                  ]"
                                                >
                                                  <div>本软件推荐用户使用至臻键音, 以设置更为丰富的声效。</div>
                                                  <div>若本选项勾选, 则自动提供可用于至臻键音的锚定功能。</div>
                                                  <div>此功能旨在方便您操作, 您也可解除锚定, 更自由的定义声效。</div>
                                                </q-tooltip>
                                              </q-icon>
                                            </q-item-label>
                                          </template>
                                        </q-option-group>
                                      </div>
                                    </q-card-section>
                                  </div>
                                  <!-- 声效选择 end -->
                                </div>

                                <!-- <q-select filled v-model="selectedEffect" :options="effectOptions" label="选择声效" /> -->
                              </div>
                            </q-card-section>

                            <q-card-actions
                              align="right"
                              :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']"
                            >
                              <q-btn
                                flat
                                label="确定"
                                color="primary"
                                v-close-popup
                                @click="
                                  saveSingleKeySoundEffectConfig(
                                    {
                                      singleKeys: selectedSingleKeys,
                                      down: keyDownSingleKeySoundEffectSelect,
                                      up: keyUpSingleKeySoundEffectSelect,
                                    },
                                    () => {
                                      q.notify({
                                        type: 'positive',
                                        position: 'top',
                                        message: '单键声效配置成功',
                                        timeout: 2000,
                                      });
                                    }
                                  )
                                "
                              />
                              <q-btn flat label="取消" color="primary" v-close-popup />
                            </q-card-actions>
                          </q-card>
                        </q-dialog>
                      </div>

                      <!-- <div class="flex flex-col gap-4">
                        <q-select filled v-model="selectedKeys" :options="keyOptions" multiple label="选择按键" />
                        <q-select filled v-model="selectedEffect" :options="effectOptions" label="选择声效" />
                      </div> -->
                      <!-- <div class="flex flex-row items-center gap-2 mb-2">
                        <q-select
                          class="w-1/2"
                          filled
                          v-model="selectedKeys"
                          :options="keyOptions"
                          multiple
                          label="选择按键"
                        />
                        <q-select
                          class="w-1/2"
                          filled
                          v-model="selectedEffect"
                          :options="effectOptions"
                          label="选择声效"
                        />
                        <q-btn flat round color="primary" icon="add" @click="addSingleKeyEffect" />
                      </div>
                      <q-card-section
                        class="row items-center q-pb-none text-sm font-bold sticky top-0 z-10 bg-white/30 backdrop-blur-sm -m-l-4 m-b-2"
                      >
                        查看已设置的单键声效
                      </q-card-section> -->
                    </q-card-section>
                    <q-card-section>
                      <div v-if="keysWithSoundEffect.size === 0" class="text-[1.06rem]">
                        目前没有已设置单键声效的按键
                      </div>
                      <div v-else class="text-[1.06rem] pb-2 font-600 text-gray-700 flex flex-row items-center">
                        已设置单键声效的按键
                        <div class="text-[0.88rem] ml-1">(点击查看声效)</div>
                      </div>
                      <div class="flex flex-wrap gap-0.8">
                        <q-chip
                          v-for="item in keysWithSoundEffect"
                          :key="item[0]"
                          dense
                          square
                          class="p-t-3.25 p-b-3.25 p-x-2.5 bg-gradient-to-b from-gray-50 to-gray-200 border-2 border-gray-300 rounded-[0.18rem] shadow-[1px_2px_1px_3px_rgba(0,0,0,0.2),inset_1px_1px_1px_rgba(255,255,255,0.6)] inset_1px_1px_1px_rgba(255,255,255,0.6)]"
                          clickable
                          @click="
                            () => {
                              // 打开查看声效的对话框
                              isShowSingleKeySoundEffectEditDialog = true;
                              currentEditingKey = Number(item[0]);
                              // TODO: 进一步, 需要在此处读取对应单键的声效设置, 并用读取的数据来初始化对话框, 以供用户的后续编辑。
                            }
                          "
                        >
                          {{
                            keyEvent_store.dikCodeToName.get(Number(item[0])) ||
                            (dikCodeToName_custom.get(Number(item[0]))
                              ? 'Temp-{' + dikCodeToName_custom.get(Number(item[0])) + '}'
                              : '') ||
                            'Dik-{' + item[0] + '}'
                          }}
                        </q-chip>
                        <q-dialog v-model="isShowSingleKeySoundEffectEditDialog">
                          <q-card style="min-width: 350px">
                            <q-card-section>
                              <div class="text-base flex flex-row items-center">
                                编辑单键 -
                                <div class="text-sm font-bold">
                                  [
                                  {{ currentEditingKeyOfName }}
                                  ]
                                </div>
                                - 的声效
                              </div>
                            </q-card-section>

                            <q-card-section class="q-pt-none">
                              <!-- 这里之后添加编辑内容 -->
                              <!-- 声效编辑  start -->
                              <div class="flex flex-row items-center justify-center gap-9 w-[88%]">
                                <q-checkbox dense v-model="isDownSoundEffectSelectEnabled_edit" label="按下声效" />
                                <q-checkbox dense v-model="isUpSoundEffectSelectEnabled_edit" label="抬起声效" />
                              </div>
                              <div class="w-full">
                                <q-card-section>
                                  <div class="flex flex-row flex-nowrap items-center m-b-3 m-l-5">
                                    <div class="flex flex-col space-y-4 w-6.5/8">
                                      <!-- 选择单键按下声效的选项, 仅支持单选 [声效编辑]-->
                                      <q-select
                                        v-show="isDownSoundEffectSelectEnabled_edit"
                                        outlined
                                        stack-label
                                        v-model="keyDownSingleKeySoundEffectSelect_edit"
                                        :options="keySingleKeySoundEffectOptions_edit"
                                        :option-label="(item: any) => {
                                              if (item.type === 'audio_files') {
                                                return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.name + soundFileList.find( (soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.type;
                                              }
                                              if (item.type === 'sounds') {
                                                // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
                                                if (item.value.soundValue?.name !== '' && item.value.soundValue?.name !== undefined) {
                                                  return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;'+ (soundList.find((sound) => sound.soundKey === item.value.soundKey)?.soundValue.name)
                                                } else {
                                                  return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                                    soundFileList.find(
                                                      (soundFile:any) =>
                                                        soundFile.sha256 === item.value.soundValue?.source_file_for_sound?.sha256 &&
                                                        soundFile.name_id === item.value.soundValue?.source_file_for_sound?.name_id
                                                    )?.name +
                                                    '     - ' +
                                                    ' [' +
                                                    item.value.soundValue?.cut?.start_time +
                                                    ' ~ ' +
                                                    item.value.soundValue?.cut?.end_time +
                                                    ']'
                                                  );
                                                }
                                              }
                                              if (item.type === 'key_sounds') {
                                                return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (item.value.keySoundValue?.name);
                                              }
                                            }"
                                        :option-value="(item: any) => {
                                              // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                              if (item.type === 'audio_files'){
                                                return item.value.sha256 + item.value.name_id
                                              }
                                              if(item.type === 'sounds'){
                                                return item.value.soundKey
                                              }
                                              if(item.type === 'key_sounds'){
                                                return item.value.keySoundKey
                                              }
                                            }"
                                        :label="`编辑单键 -[ ${currentEditingKeyOfName} ]- 的按下声效`"
                                        use-chips
                                        dense
                                        @popup-hide="
                                          () => {
                                            if (
                                              // 为避免循环依赖, 此处作为锚定功能选择声效时的判断逻辑; 而删除声效时的判断逻辑, 在watch中书写。
                                              isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit &&
                                              isAnchoringUltimatePerfectionKeySound_singleKey_edit &&
                                              // 这里的?是防止在勾选至臻键音的条件下, 仅打开选项菜单且未做任何选择就关闭时, mode的null值内 没有type字段引起报错。
                                              keyDownSingleKeySoundEffectSelect_edit?.type === 'key_sounds'
                                            ) {
                                              keyUpSingleKeySoundEffectSelect_edit =
                                                keyDownSingleKeySoundEffectSelect_edit;
                                            }
                                          }
                                        "
                                        class="max-w-full"
                                      />
                                      <!-- 选择单键抬起声效的选项, 仅支持单选 [声效编辑]-->
                                      <q-select
                                        v-show="isUpSoundEffectSelectEnabled_edit"
                                        outlined
                                        stack-label
                                        v-model="keyUpSingleKeySoundEffectSelect_edit"
                                        :options="keySingleKeySoundEffectOptions_edit"
                                        :option-label="(item: any) => {
                                              if (item.type === 'audio_files') {
                                                return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + soundFileList.find((soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.name + soundFileList.find( (soundFile:any) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id)?.type;
                                              }
                                              if (item.type === 'sounds') {
                                                // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
                                                if (item.value.soundValue?.name !== '' && item.value.soundValue?.name !== undefined) {
                                                  return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;'+ (soundList.find((sound) => sound.soundKey === item.value.soundKey)?.soundValue.name)
                                                } else {
                                                  return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (
                                                    soundFileList.find(
                                                      (soundFile:any) =>
                                                        soundFile.sha256 === item.value.soundValue?.source_file_for_sound?.sha256 &&
                                                        soundFile.name_id === item.value.soundValue?.source_file_for_sound?.name_id
                                                    )?.name +
                                                    '     - ' +
                                                    ' [' +
                                                    item.value.soundValue?.cut?.start_time +
                                                    ' ~ ' +
                                                    item.value.soundValue?.cut?.end_time +
                                                    ']'
                                                  );
                                                }
                                              }
                                              if (item.type === 'key_sounds') {
                                                return (options.find((option) =>  item.type === option.value)?.label ) + '&nbsp;&nbsp;&sect;&nbsp;&nbsp;&nbsp;' + (item.value.keySoundValue?.name);
                                              }
                                            }"
                                        :option-value="(item: any) => {
                                              // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                              if (item.type === 'audio_files'){
                                                return item.value.sha256 + item.value.name_id
                                              }
                                              if(item.type === 'sounds'){
                                                return item.value.soundKey
                                              }
                                              if(item.type === 'key_sounds'){
                                                return item.value.keySoundKey
                                              }
                                            }"
                                        :label="`编辑单键 -[ ${currentEditingKeyOfName} ]- 的抬起声效`"
                                        use-chips
                                        dense
                                        @popup-hide="
                                          () => {
                                            // 为避免循环依赖, 此处作为锚定功能选择声效时的判断逻辑; 而删除声效时的判断逻辑, 在watch中书写。
                                            if (
                                              isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit &&
                                              isAnchoringUltimatePerfectionKeySound_singleKey_edit &&
                                              // 这里的?是防止在勾选至臻键音的条件下, 仅打开选项菜单且未做任何选择就关闭时, mode的null值内 没有type字段引起报错。
                                              keyUpSingleKeySoundEffectSelect_edit?.type === 'key_sounds'
                                            ) {
                                              keyDownSingleKeySoundEffectSelect_edit =
                                                keyUpSingleKeySoundEffectSelect_edit;
                                            }
                                          }
                                        "
                                        class="max-w-full"
                                      />
                                    </div>
                                    <div
                                      v-show="isDownSoundEffectSelectEnabled_edit && isUpSoundEffectSelectEnabled_edit"
                                      class="flex justify-end -m-l-2"
                                    >
                                      <q-icon
                                        @click="
                                          isAnchoringUltimatePerfectionKeySound_singleKey_edit =
                                            !isAnchoringUltimatePerfectionKeySound_singleKey_edit
                                        "
                                        size="2.75rem"
                                        v-if="isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit"
                                      >
                                        <template v-if="isAnchoringUltimatePerfectionKeySound_singleKey_edit">
                                          <!-- 锚定 [声效编辑]-->
                                          <svg
                                            version="1.1"
                                            xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 27.838698814764143 45.55476707640558"
                                            width="27.838698814764143"
                                            height="45.55476707640558"
                                          >
                                            <g
                                              stroke-linecap="round"
                                              transform="translate(10.257971830020665 11.315110817092783) rotate(0 3.661377577361293 4.912529303576154)"
                                            >
                                              <path
                                                d="M1.83 0 C3.09 0, 4.34 0, 5.49 0 M1.83 0 C2.76 0, 3.69 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.07, 7.32 6.3, 7.32 7.99 M7.32 1.83 C7.32 3.67, 7.32 5.51, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.52 9.83, 3.56 9.83, 1.83 9.83 M5.49 9.83 C4.3 9.83, 3.12 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.74, 0 5.48, 0 1.83 M0 7.99 C0 6.56, 0 5.12, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                                stroke="#1e1e1e"
                                                stroke-width="1"
                                                fill="none"
                                              ></path>
                                            </g>
                                            <g stroke-linecap="round">
                                              <g
                                                transform="translate(13.900885269327091 18.0363210366226) rotate(0 -0.0030445053470202765 2.3801245195212495)"
                                              >
                                                <path
                                                  d="M0 0 C0 0.79, -0.01 3.97, -0.01 4.76 M0 0 C0 0.79, -0.01 3.97, -0.01 4.76"
                                                  stroke="#1e1e1e"
                                                  stroke-width="1"
                                                  fill="none"
                                                ></path>
                                              </g>
                                            </g>
                                            <mask></mask>
                                            <g
                                              stroke-linecap="round"
                                              transform="translate(10.257971830020665 24.438144439265034) rotate(0 3.661377577361293 4.912529303576154)"
                                            >
                                              <path
                                                d="M1.83 0 C2.75 0, 3.67 0, 5.49 0 M1.83 0 C2.98 0, 4.12 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.09, 7.32 6.36, 7.32 7.99 M7.32 1.83 C7.32 4.06, 7.32 6.28, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.11 9.83, 2.73 9.83, 1.83 9.83 M5.49 9.83 C4.08 9.83, 2.67 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.4, 0 4.8, 0 1.83 M0 7.99 C0 6.46, 0 4.92, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                                stroke="#1e1e1e"
                                                stroke-width="1"
                                                fill="none"
                                              ></path>
                                            </g>
                                            <g stroke-linecap="round">
                                              <g
                                                transform="translate(13.900885269326977 27.541992826887526) rotate(0 -0.0030445053470202765 -2.3801245195212495)"
                                              >
                                                <path
                                                  d="M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76 M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76"
                                                  stroke="#1e1e1e"
                                                  stroke-width="1"
                                                  fill="none"
                                                ></path>
                                              </g>
                                            </g>
                                            <mask></mask>
                                            <g
                                              stroke-linecap="round"
                                              transform="translate(10 10) rotate(0 3.9193494073820716 12.77738353820279)"
                                            >
                                              <path
                                                d="M1.96 0 C3.5 0, 5.03 0, 5.88 0 M1.96 0 C3.41 0, 4.86 0, 5.88 0 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M7.84 1.96 C7.84 8.03, 7.84 14.11, 7.84 23.6 M7.84 1.96 C7.84 7.33, 7.84 12.7, 7.84 23.6 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M5.88 25.55 C4.57 25.55, 3.25 25.55, 1.96 25.55 M5.88 25.55 C4.56 25.55, 3.25 25.55, 1.96 25.55 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M0 23.6 C0 16.78, 0 9.97, 0 1.96 M0 23.6 C0 15.66, 0 7.72, 0 1.96 M0 1.96 C0 0.65, 0.65 0, 1.96 0 M0 1.96 C0 0.65, 0.65 0, 1.96 0"
                                                stroke="transparent"
                                                stroke-width="1"
                                                fill="none"
                                              ></path>
                                            </g>
                                          </svg>
                                        </template>
                                        <template v-else>
                                          <!-- 锚定解除 [声效编辑]-->
                                          <svg
                                            version="1.1"
                                            xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 27.838698814764143 45.55476707640558"
                                            width="27.838698814764143"
                                            height="45.55476707640558"
                                          >
                                            <g
                                              stroke-linecap="round"
                                              transform="translate(10.257971830020779 10.375546926705098) rotate(0 3.661377577361293 4.912529303576154)"
                                            >
                                              <path
                                                d="M1.83 0 C3.09 0, 4.34 0, 5.49 0 M1.83 0 C2.76 0, 3.69 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.07, 7.32 6.3, 7.32 7.99 M7.32 1.83 C7.32 3.67, 7.32 5.51, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.52 9.83, 3.56 9.83, 1.83 9.83 M5.49 9.83 C4.3 9.83, 3.12 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.74, 0 5.48, 0 1.83 M0 7.99 C0 6.56, 0 5.12, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                                stroke="#1e1e1e"
                                                stroke-width="1"
                                                fill="none"
                                              ></path>
                                            </g>
                                            <g stroke-linecap="round">
                                              <g
                                                transform="translate(13.900885269327205 17.096757146234918) rotate(0 -0.0030445053470202765 2.3801245195212513)"
                                              >
                                                <path
                                                  d="M0 0 C0 0.79, -0.01 3.97, -0.01 4.76 M0 0 C0 0.79, -0.01 3.97, -0.01 4.76"
                                                  stroke="#1e1e1e"
                                                  stroke-width="1"
                                                  fill="none"
                                                ></path>
                                              </g>
                                            </g>
                                            <mask></mask>
                                            <g
                                              stroke-linecap="round"
                                              transform="translate(10.257971830020779 25.379515181404656) rotate(0 3.661377577361293 4.912529303576154)"
                                            >
                                              <path
                                                d="M1.83 0 C2.75 0, 3.67 0, 5.49 0 M1.83 0 C2.98 0, 4.12 0, 5.49 0 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M5.49 0 C6.71 0, 7.32 0.61, 7.32 1.83 M7.32 1.83 C7.32 4.09, 7.32 6.36, 7.32 7.99 M7.32 1.83 C7.32 4.06, 7.32 6.28, 7.32 7.99 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M7.32 7.99 C7.32 9.21, 6.71 9.83, 5.49 9.83 M5.49 9.83 C4.11 9.83, 2.73 9.83, 1.83 9.83 M5.49 9.83 C4.08 9.83, 2.67 9.83, 1.83 9.83 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M1.83 9.83 C0.61 9.83, 0 9.21, 0 7.99 M0 7.99 C0 6.4, 0 4.8, 0 1.83 M0 7.99 C0 6.46, 0 4.92, 0 1.83 M0 1.83 C0 0.61, 0.61 0, 1.83 0 M0 1.83 C0 0.61, 0.61 0, 1.83 0"
                                                stroke="#1e1e1e"
                                                stroke-width="1"
                                                fill="none"
                                              ></path>
                                            </g>
                                            <g stroke-linecap="round">
                                              <g
                                                transform="translate(13.900885269327091 28.483363569027148) rotate(0 -0.0030445053470202765 -2.3801245195212495)"
                                              >
                                                <path
                                                  d="M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76 M0 0 C0 -0.79, -0.01 -3.97, -0.01 -4.76"
                                                  stroke="#1e1e1e"
                                                  stroke-width="1"
                                                  fill="none"
                                                ></path>
                                              </g>
                                            </g>
                                            <mask></mask>
                                            <g
                                              stroke-linecap="round"
                                              transform="translate(10 10) rotate(0 3.9193494073820716 12.777383538202791)"
                                            >
                                              <path
                                                d="M1.96 0 C3.5 0, 5.03 0, 5.88 0 M1.96 0 C3.41 0, 4.86 0, 5.88 0 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M5.88 0 C7.19 0, 7.84 0.65, 7.84 1.96 M7.84 1.96 C7.84 8.03, 7.84 14.11, 7.84 23.6 M7.84 1.96 C7.84 7.33, 7.84 12.7, 7.84 23.6 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M7.84 23.6 C7.84 24.9, 7.19 25.55, 5.88 25.55 M5.88 25.55 C4.57 25.55, 3.25 25.55, 1.96 25.55 M5.88 25.55 C4.56 25.55, 3.25 25.55, 1.96 25.55 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M1.96 25.55 C0.65 25.55, 0 24.9, 0 23.6 M0 23.6 C0 16.78, 0 9.97, 0 1.96 M0 23.6 C0 15.66, 0 7.72, 0 1.96 M0 1.96 C0 0.65, 0.65 0, 1.96 0 M0 1.96 C0 0.65, 0.65 0, 1.96 0"
                                                stroke="transparent"
                                                stroke-width="1"
                                                fill="none"
                                              ></path>
                                            </g>
                                          </svg>
                                        </template>
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <span class="text-sm">至臻键音</span>
                                          <span
                                            class="text-sm"
                                            v-if="isAnchoringUltimatePerfectionKeySound_singleKey_edit"
                                          >
                                            (已锚定)<br
                                          /></span>
                                          <span class="text-sm" v-else> (已解除锚定)<br /></span>
                                          <span>此锚定仅作用于在设置声效为至臻键音时。</span>
                                        </q-tooltip>
                                      </q-icon>
                                    </div>
                                  </div>
                                  <div
                                    class="h-16 m-l-9"
                                    v-show="isDownSoundEffectSelectEnabled_edit || isUpSoundEffectSelectEnabled_edit"
                                  >
                                    <q-option-group
                                      dense
                                      v-model="singleKeyTypeGroup_edit"
                                      :options="options"
                                      type="checkbox"
                                    >
                                      <template #label-0="props">
                                        <q-item-label>
                                          {{ props.label }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>本软件支持将音频源文件直接设置为声效。</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-1="props">
                                        <q-item-label>
                                          {{ props.label }}
                                          <q-icon name="info" color="primary" class="p-l-4.5 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>本软件支持将裁剪定义好的声音设置为声效。</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-2="props">
                                        <q-item-label>
                                          {{ props.label }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>本软件推荐用户使用至臻键音, 以设置更为丰富的声效。</div>
                                              <div>若本选项勾选, 则自动提供可用于至臻键音的锚定功能。</div>
                                              <div>此功能旨在方便您操作, 您也可解除锚定, 更自由的定义声效。</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                    </q-option-group>
                                  </div>
                                </q-card-section>
                              </div>
                              <!-- 声效编辑  end -->
                            </q-card-section>

                            <q-card-actions align="right">
                              <q-btn flat label="取消" color="primary" v-close-popup />
                              <q-btn flat label="确定" color="primary" v-close-popup />
                            </q-card-actions>
                          </q-card>
                        </q-dialog>
                      </div>
                    </q-card-section>
                    <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
                      <q-btn flat label="确定" color="primary" v-close-popup />
                      <q-btn flat label="取消" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
            </q-stepper-navigation>
            <q-stepper-navigation>
              <q-btn @click="step = 5" color="primary" label="Continue" />
              <q-btn flat @click="step = 3" color="primary" label="Back" class="q-ml-sm" />
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
import { QSelect, useQuasar } from 'quasar';
import {
  ConfigGet,
  ConfigSet,
  LoadConfig,
  SendFileToServer,
  SoundFileRename,
  SoundFileDelete,
  PlaySound,
  ConfigDelete,
} from 'src/boot/query/keytonePkg-query';
import { useAppStore } from 'src/stores/app-store';
import { useKeyEventStore } from 'src/stores/keyEvent-store';
import { computed, onBeforeMount, ref, watch, useTemplateRef, reactive, nextTick } from 'vue';
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
const soundFileList = ref<Array<{ sha256: string; name_id: string; name: string; type: string }>>([]);
const selectedSoundFile = ref<{ sha256: string; name_id: string; name: string; type: string }>({
  sha256: '',
  name_id: '',
  name: '',
  type: '',
});

// 声音制作(制作新的声音)
const createNewSound = ref(false);
const soundName = ref<string>('');
const sourceFileForSound = ref<{ sha256: string; name_id: string; name: string; type: string }>({
  sha256: '',
  name_id: '',
  name: '',
  type: '',
});
const soundStartTime = ref<number>(0);
const soundEndTime = ref<number>(0);
const soundVolume = ref<number>(0.0);

// 修改 confirmAddingSound 函数为更通用的形式, 并将 confirmAddingSound 重构名称为 saveSoundConfig
function saveSoundConfig(params: {
  soundKey?: string; // 可选参数，存在则为修改操作，不存在则为添加操作
  source_file_for_sound: { sha256: string; name_id: string; type: string };
  name: string; // 声音名称, 一般为空, 也可由用户自行定义
  cut: {
    start_time: number;
    end_time: number;
    volume: number;
  };
  onSuccess?: () => void; // 成功后的回调函数, 可用于一些清空操作
}) {
  // 必须选择一个源文件
  if (
    params.source_file_for_sound.sha256 === '' &&
    params.source_file_for_sound.type === '' &&
    params.source_file_for_sound.name_id === ''
  ) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: '请选择一个源文件',
      timeout: 5,
    });
    return;
  }
  // 结束时间必须大于开始时间
  if (params.cut.end_time <= params.cut.start_time) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: '结束时间必须大于开始时间',
      timeout: 5,
    });
    return;
  }

  // 创建一个新对象,不包含soundKey和onSuccess回调
  const configParams = {
    source_file_for_sound: params.source_file_for_sound,
    name: params.name,
    cut: params.cut,
  };

  // 如果有soundKey则为修改操作，否则为添加操作（生成新的key）
  const key = params.soundKey || nanoid();

  ConfigSet('sounds.' + key, configParams).then((re) => {
    if (re) {
      q.notify({
        type: 'positive',
        position: 'top',
        message: params.soundKey ? '修改成功' : '添加成功',
        timeout: 5,
      });
      params.onSuccess?.();
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: params.soundKey ? '修改失败' : '添加失败',
        timeout: 5,
      });
    }
  });
}

function deleteSound(params: { soundKey: string; onSuccess?: () => void }) {
  ConfigDelete('sounds.' + params.soundKey).then((re) => {
    if (re) {
      params.onSuccess?.();
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: '删除失败',
        timeout: 5,
      });
    }
  });
}

// 重构previewSound函数,使用相同的参数结构
function previewSound(params: {
  source_file_for_sound: { sha256: string; name_id: string; type: string };
  cut: {
    start_time: number;
    end_time: number;
    volume: number;
  };
}) {
  console.debug('预览声音');
  if (
    params.source_file_for_sound.sha256 === '' &&
    params.source_file_for_sound.type === '' &&
    params.source_file_for_sound.name_id === ''
  ) {
    q.notify({
      type: 'warning',
      position: 'top',
      message: '请先选择音频文件',
      timeout: 5000,
    });
    return;
  }
  // 结束时间必须大于开始时间
  if (params.cut.end_time <= params.cut.start_time) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: '结束时间必须大于开始时间',
      timeout: 5,
    });
    return;
  }
  // 时间值不能为负数
  if (params.cut.start_time < 0 || params.cut.end_time < 0) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: '时间值不能为负数',
      timeout: 5,
    });
  }

  PlaySound(
    params.source_file_for_sound.sha256,
    params.source_file_for_sound.type,
    params.cut.start_time,
    params.cut.end_time,
    params.cut.volume
  ).then((result) => {
    if (!result) {
      q.notify({
        type: 'negative',
        position: 'top',
        message: '播放失败',
        timeout: 5000,
      });
    }
  });
}

// 声音编辑(编辑已有声音)
const showEditSoundDialog = ref(false);
const soundList = ref<
  Array<{
    soundKey: string;
    soundValue: {
      cut: { start_time: number; end_time: number; volume: number };
      name: string;
      source_file_for_sound: { sha256: string; name_id: string; type: string };
    };
  }>
>([]); // TIPS: 此处类型一定要指定清楚, 否则在 组件<q-select :option-label="(item)=>{}"> 中, 会发生类型错误(且检测不到原因, 准确指定类型不使用any后, 才解决此问题)--- 我记得之前在做声音源文件的编辑列表时, 就好像遇到过, 但目前找不到注释内容了。
const selectedSound = ref<{
  soundKey: string;
  soundValue: {
    cut: { start_time: number; end_time: number; volume: number };
    name: string;
    source_file_for_sound: { sha256: string; name_id: string; type: string };
  };
}>(); // 此处无需初始化, 但类型一定要指定清楚

watch(selectedSound, () => {
  console.debug('观察selectedSound=', selectedSound.value);
});

// 选项列表(其中包含: 源文件、声音、按键音 三种选项可供选择)
const options = [
  { label: '源文件', value: 'audio_files' },
  { label: '声音', value: 'sounds' },
  { label: '至臻键音', value: 'key_sounds' },
];

// 按键音
type PlayMode = { label: string; mode: string };
const playModeOptions: Array<PlayMode> = [
  { label: '独立', mode: 'single' },
  { label: '随机', mode: 'random' },
  { label: '循环', mode: 'loop' },
];

// 按键音制作
const createNewKeySound = ref(false);

// -- createNewKeySound
const keySoundName = ref<string>('新的按键音');
const configureDownSound = ref(false);
const configureUpSound = ref(false);

// -- configureDownSound
const selectedSoundsForDown = ref<Array<any>>([]);
const playModeForDown = ref<PlayMode>({ label: '随机', mode: 'random' });
const maxSelectionForDown = computed(() => {
  return playModeForDown.value.mode === 'single' ? 1 : Infinity;
});
/* --- 在vue3.5中, 用useTemplateRef方式获取dom元素, 有助于增强可读性
const downSoundSelectDom = ref<QSelect>(); // 在vue3.5中, 用useTemplateRef方式获取dom元素, 有助于增强可读性*/
const downSoundSelectDom = useTemplateRef<QSelect>('downSoundSelectDom');
const downTypeGroup = ref<Array<string>>(['sounds']);
const downSoundList = computed(() => {
  const List: Array<any> = [];
  if (downTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (downTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (downTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  console.debug('downSoundList=', List);
  return List;
});

// -- configureUpSound
const selectedSoundsForUp = ref<Array<any>>([]);
const playModeForUp = ref<PlayMode>({ label: '随机', mode: 'random' });
const maxSelectionForUp = computed(() => {
  return playModeForUp.value.mode === 'single' ? 1 : Infinity;
});
/* --- 在vue3.5中, 用useTemplateRef方式获取dom元素, 有助于增强可读性
const upSoundSelectDom = ref<QSelect>();*/
const upSoundSelectDom = useTemplateRef<QSelect>('upSoundSelectDom');
const upTypeGroup = ref<Array<string>>(['sounds']);
const upSoundList = computed(() => {
  const List: Array<any> = [];
  if (upTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (upTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (upTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  console.debug('upSoundList=', List);
  return List;
});

// 按键音编辑
const editExistingKeySound = ref(false);

// -- editExistingKeySound
const edit_configureDownSound = ref(false);
const edit_configureUpSound = ref(false);

// -- edit_configureDownSound / edit_configureUpSound
const edit_downSoundSelectDom = useTemplateRef<QSelect>('edit_downSoundSelectDom');
const edit_upSoundSelectDom = useTemplateRef<QSelect>('edit_upSoundSelectDom');
const edit_downTypeGroup = ref<Array<string>>(['sounds']);
const edit_upTypeGroup = ref<Array<string>>(['sounds']);
const edit_downSoundList = computed(() => {
  const List: Array<any> = [];
  if (edit_downTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (edit_downTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (edit_downTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  console.debug('edit_downSoundList=', List);
  return List;
});
const edit_upSoundList = computed(() => {
  const List: Array<any> = [];
  if (edit_upTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (edit_upTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (edit_upTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  console.debug('edit_upSoundList=', List);
  return List;
});
const keySoundList = ref<Array<any>>([]);
const selectedKeySound = ref<any>();

// 改变selectedKeySound.value.keySoundValue.down.value和selectedKeySound.value.keySoundValue.up.value的类型结构, 使其符合选择输入框组件的使用需求
watch(selectedKeySound, () => {
  console.debug('观察selectedKeySound=', selectedKeySound.value);
  selectedKeySound.value.keySoundValue.down.value = selectedKeySound.value.keySoundValue.down.value.map((item: any) => {
    /**
     * json中的存储格式分别是
     *  {key:'audio_files', value:{sha256: string, name_id: string, type:string}}
     *  {key:'sounds', value:string} // 此处value, 是soundKey
     *  {key:'key_sounds', value:string} // 此处value, 是keySoundKey
     */
    if (item.type === 'audio_files') {
      return {
        type: 'audio_files',
        value: soundFileList.value.find(
          (soundFile) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id
        ),
      };
    }
    if (item.type === 'sounds') {
      return {
        type: 'sounds',
        // value: soundList.value.find((sound) => sound.soundKey === item.value.soundKey),
        value: soundList.value.find((sound) => sound.soundKey === item.value),
      };
    }
    if (item.type === 'key_sounds') {
      return {
        type: 'key_sounds',
        // value: keySoundList.value.find((keySound) => keySound.keySoundKey === item.value.keySoundKey),
        value: keySoundList.value.find((keySound) => keySound.keySoundKey === item.value),
      };
    }
    return item;
  });
  selectedKeySound.value.keySoundValue.up.value = selectedKeySound.value.keySoundValue.up.value.map((item: any) => {
    /**
     * json中的存储格式分别是
     *  {key:'audio_files', value:{sha256: string, name_id: string, type:string}}
     *  {key:'sounds', value:string} // 此处value, 是soundKey
     *  {key:'key_sounds', value:string} // 此处value, 是keySoundKey
     */
    if (item.type === 'audio_files') {
      return {
        type: 'audio_files',
        value: soundFileList.value.find(
          (soundFile) => soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id
        ),
      };
    }
    if (item.type === 'sounds') {
      return {
        type: 'sounds',
        // value: soundList.value.find((sound) => sound.soundKey === item.value.soundKey),
        value: soundList.value.find((sound) => sound.soundKey === item.value),
      };
    }
    if (item.type === 'key_sounds') {
      return {
        type: 'key_sounds',
        // value: keySoundList.value.find((keySound) => keySound.keySoundKey === item.value.keySoundKey),
        value: keySoundList.value.find((keySound) => keySound.keySoundKey === item.value),
      };
    }
    return item;
  });
});

// 按键音api
// -- 保存按键音配置
function saveKeySoundConfig(
  params: {
    key: string;
    name: string;
    down: { mode: string; value: Array<any> };
    up: { mode: string; value: Array<any> };
  },
  onSuccess?: () => void
) {
  let isReturn = false;
  if (params.down.mode === 'single' && params.down.value.length > 1) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: '失败: "按下声音"不符合要求, 请检查修改后重试',
      timeout: 3000,
    });
    isReturn = true;
  }
  if (params.up.mode === 'single' && params.up.value.length > 1) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: '失败: "抬起声音"不符合要求, 请检查修改后重试',
      timeout: 3000,
    });
    isReturn = true;
  }
  if (isReturn) {
    return;
  }

  const configParams = {
    name: params.name,
    down: {
      mode: params.down.mode,
      value: params.down.value.map((item) => {
        if (item.type === 'audio_files') {
          return {
            type: 'audio_files',
            value: { sha256: item.value.sha256, name_id: item.value.name_id, type: item.value.type },
          };
        }
        if (item.type === 'sounds') {
          return { type: 'sounds', value: item.value.soundKey };
        }
        if (item.type === 'key_sounds') {
          return { type: 'key_sounds', value: item.value.keySoundKey };
        }
      }),
    },
    up: {
      mode: params.up.mode,
      value: params.up.value.map((item) => {
        if (item.type === 'audio_files') {
          return {
            type: 'audio_files',
            value: { sha256: item.value.sha256, name_id: item.value.name_id, type: item.value.type },
          };
        }
        if (item.type === 'sounds') {
          return { type: 'sounds', value: item.value.soundKey };
        }
        if (item.type === 'key_sounds') {
          return { type: 'key_sounds', value: item.value.keySoundKey };
        }
      }),
    },
  };

  const key = params.key || nanoid();
  ConfigSet('key_sounds.' + key, configParams).then((re) => {
    if (re) {
      q.notify({
        type: 'positive',
        position: 'top',
        message: params.key ? '修改成功' : '添加成功',
        timeout: 5,
      });
      if (onSuccess) {
        onSuccess();
      }
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: params.key ? '修改失败' : '添加失败',
        timeout: 5,
      });
    }
  });
}

// -- 删除键音
function deleteKeySound(params: { keySoundKey: string; onSuccess?: () => void }) {
  ConfigDelete('key_sounds.' + params.keySoundKey).then((re) => {
    if (re) {
      params.onSuccess?.();
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: '删除失败',
        timeout: 5,
      });
    }
  });
}

// 按键联动声效

// -- 内嵌测试音是否使能
const isEnableEmbeddedTestSound = reactive({
  down: true,
  up: true,
}); // 该字段"直接"与配置文件相映射

// -- 全键声效
const showEveryKeyEffectDialog = ref(false);

const keyDownUnifiedSoundEffectSelect = ref<any>();
const keyUpUnifiedSoundEffectSelect = ref<any>();
const unifiedTypeGroup = ref<Array<string>>(['sounds']);
const keyUnifiedSoundEffectOptions = computed(() => {
  const List: Array<any> = [];
  if (unifiedTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (unifiedTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (unifiedTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  console.debug('观察keyUnifiedSoundEffectOptions=', List);
  return List;
});
const isShowUltimatePerfectionKeySoundAnchoring = computed(() => {
  return unifiedTypeGroup.value.includes('key_sounds');
});
const isAnchoringUltimatePerfectionKeySound = ref(true);
watch(keyDownUnifiedSoundEffectSelect, (newVal, oldVal) => {
  console.debug('观察keyDownUnifiedSoundEffectSelect=', keyDownUnifiedSoundEffectSelect.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (isShowUltimatePerfectionKeySoundAnchoring.value && isAnchoringUltimatePerfectionKeySound.value) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyUpUnifiedSoundEffectSelect.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyUpUnifiedSoundEffectSelect.value = keyDownUnifiedSoundEffectSelect.value;
      }
    }
  }
});
watch(keyUpUnifiedSoundEffectSelect, (newVal, oldVal) => {
  console.debug('观察keyUpUnifiedSoundEffectSelect=', keyUpUnifiedSoundEffectSelect.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (isShowUltimatePerfectionKeySoundAnchoring.value && isAnchoringUltimatePerfectionKeySound.value) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyDownUnifiedSoundEffectSelect.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyDownUnifiedSoundEffectSelect.value = keyUpUnifiedSoundEffectSelect.value;
      }
    }
  }
});

function saveUnifiedSoundEffectConfig(params: { down: any; up: any }, onSuccess?: () => void) {
  const keyTone_global = {
    down: {
      type: params.down?.type || '',
      // TIPS: 此处需要注意, 我们需要的是此lambda表达式执行后的返回值赋值给value, 而不是直接将lambda表达式赋值给value。(此处用三元表达式会更为直观)
      value: (() => {
        if (params.down?.type === 'audio_files') {
          return { sha256: params.down.value.sha256, name_id: params.down.value.name_id, type: params.down.value.type };
        }
        if (params.down?.type === 'sounds') {
          return params.down.value.soundKey;
        }
        if (params.down?.type === 'key_sounds') {
          return params.down.value.keySoundKey;
        }
        return '';
      })(),
    },
    up: {
      type: params.up?.type || '',
      // TIPS: 此处需要注意, 我们需要的是此lambda表达式执行后的返回值赋值给value, 而不是直接将lambda表达式赋值给value。(此处用三元表达式会更为直观)
      value: (() => {
        if (params.up?.type === 'audio_files') {
          return { sha256: params.up.value.sha256, name_id: params.up.value.name_id, type: params.up.value.type };
        }
        if (params.up?.type === 'sounds') {
          return params.up.value.soundKey;
        }
        if (params.up?.type === 'key_sounds') {
          return params.up.value.keySoundKey;
        }
        return '';
      })(),
    },
  };

  ConfigSet('key_tone.global', keyTone_global)
    .then((re) => {
      if (re) {
        onSuccess?.();
      } else {
        q.notify({
          type: 'negative',
          position: 'top',
          message: '全键声效配置失败',
          timeout: 5000,
        });
      }
    })
    .catch((err) => {
      console.error('全键声效配置时发生错误:', err);
      q.notify({
        type: 'negative',
        position: 'top',
        message: '全键声效配置失败',
        timeout: 5000,
      });
    });
}

// -- 单键声效

// -- -- 选择按键
const showSingleKeyEffectDialog = ref(false);

const isShowAddOrSettingSingleKeyEffectDialog = ref(false);
const singleKeysSelectRef = useTemplateRef<QSelect>('singleKeysSelectRef');
const selectedSingleKeys = ref<Array<number>>([]);
const dikCodeToName_custom = reactive<Map<number, string>>(new Map<number, string>()); // 用于被持久化的实时映射数据的同步工作。(优先级很低, 甚至会被更新的版本覆盖, 仅勉强起个保底作用, 甚至项目完善后有可能删除这个)

const keyEvent_store = useKeyEventStore();

const isRecordingSingleKeys = ref(false);

const keyOptions = computed(() => {
  // 将 Map 转换为数组形式的选项
  if (isRecordingSingleKeys.value) {
    return [];
  } else {
    // 默认以系统主映射表的keys为主
    const reArray = Array.from(keyEvent_store.dikCodeToName.keys());

    // 其次是temp的keys(虽然后续可能会删除相关逻辑。)
    Array.from(dikCodeToName_custom.keys()).forEach((item) => {
      // 在遍历temp映射表的过程中, 检查其每个key是否在主映射表的keys中。如果主映射表中没有, 才会往里加(因为会始终以主映射表为主)。
      if (!reArray.includes(item)) {
        reArray.push(item);
      }
    });

    return reArray;
  }
});
const filterOptions = ref(keyOptions.value); // 用于过滤选项

const isGetsFocused = ref(false);
const recordingSingleKeysCallback = (keycode: number, keyName: string) => {
  console.debug('keycode=', keycode, 'keyName=', keyName);

  // 如果按键不在列表中，则添加
  if (!selectedSingleKeys.value.includes(keycode)) {
    selectedSingleKeys.value.push(keycode);
  } else {
    q.notify({
      type: 'info',
      position: 'top',
      message: '所录制按键已在选择框中',
      timeout: 1000,
    });
  }

  console.debug('当前已选择的按键:', selectedSingleKeys.value);
};
const persistentSingleKeysDataCallback = (keycode: number, keyName: string) => {
  dikCodeToName_custom.set(keycode, keyName);
  ConfigSet('custom_single_keys_name.' + keycode, keyName);
};

watch(isShowAddOrSettingSingleKeyEffectDialog, (newVal) => {
  if (!newVal) {
    keyEvent_store.clearKeyStateCallback_Record();
    keyEvent_store.clearKeyStateCallback_PersistentData();

    // 当通过点击对话框外使得对话框关闭时, 不会触发失去焦点的事件(因此此时isGetsFocused的值不会被置为false, 故补充此逻辑)
    isGetsFocused.value = false;
  }
});

watch(isRecordingSingleKeys, (newVal) => {
  if (newVal) {
    // 录制单键时, 清空输入框。(由于是录制, 因此需要清空输入框, 防止用户输入内容。)
    // * 如何防止用户输入内容?
    // * * 当然也可以利用updateInputValue。但有更简单的解决思路, 即定义组件特有属性maxlength为0即可阻止用户输入内容。
    singleKeysSelectRef.value?.updateInputValue('');

    keyEvent_store.setKeyStateCallback_Record(recordingSingleKeysCallback);
    keyEvent_store.setKeyStateCallback_PersistentData(persistentSingleKeysDataCallback);
  } else {
    keyEvent_store.clearKeyStateCallback_Record();
    keyEvent_store.clearKeyStateCallback_PersistentData();
  }
});

watch(isGetsFocused, (newVal) => {
  if (newVal && isRecordingSingleKeys.value) {
    keyEvent_store.setKeyStateCallback_Record(recordingSingleKeysCallback);
    keyEvent_store.setKeyStateCallback_PersistentData(persistentSingleKeysDataCallback);
  } else {
    keyEvent_store.clearKeyStateCallback_Record();
    keyEvent_store.clearKeyStateCallback_PersistentData();
  }
});

//   - completed(已完成)   FIXME: 修复'Backspace'按键, 在录制过程中, 删除已选择列表中最后一项的bug
let oldSelectedSingleKeys: Array<any> = [];

function preventDefaultKeyBehaviorWhenRecording(event: KeyboardEvent) {
  if (isRecordingSingleKeys.value) {
    // TIPS: 这里被打印两次的原因可能是以下, 不过不用担心和处理, 因为没有bug。
    // 1. 首先触发 input 元素的 keydown 事件
    // 2. 然后冒泡到 q-select 组件
    // 3. 最后最后冒泡到全局监听器
    console.debug('event.key=', event.key);

    //   - completed(已完成)   FIXME: 修复'Enter'按键无法被录制的bug
    if (event.key === 'Enter') {
      // 虽然无法录制'Enter'事件的原因就是select组件阻止了默认的'Enter'事件的冒泡行为,
      // * 但为防止quasar后续更新改变它, 便再次手动阻止一次, 以防止本次修复被quasar的更新影响。
      event.stopPropagation(); // 阻止事件冒泡

      //   ↓    - completed(已完成)   增加location字段即可。原理:这是手动构建事件时, 缺少了构建UUID必要的location字段, 导致小数字键盘中的'Enter'被按下时, 在前端事件状态集中, 创建了一个并不存在的 UUID 的及对应的按下状态 的假按键。
      // FIXME: 在数字键盘'enter'按键按下后, 再次按下任何按键都 报多个按键同时按下的消息, 这是一个bug。
      // 手动创建并分发一个新的键盘事件 TIPS: 没必要(虽然修复了, 但我并不打算继续这样用)
      // const newEvent = new KeyboardEvent('keydown', {
      //   key: event.key,
      //   code: event.code,
      //   keyCode: event.keyCode,
      //   which: event.which,
      //   altKey: event.altKey,
      //   location: event.location, // 增加此字段, 以修复FIXME。
      //   ctrlKey: event.ctrlKey,
      //   shiftKey: event.shiftKey,
      //   metaKey: event.metaKey,
      //   bubbles: true,
      //   cancelable: true,
      // });
      // document.dispatchEvent(newEvent);

      // 其实干脆直接使用当前事件创建新事件也行( TIPS: 不要直接使用event, 不然会因重复分发相同引用的事件而报错)
      const newEvent = new KeyboardEvent('keydown', event);
      document.dispatchEvent(newEvent);

      // TIPS: 这样是不对的-> 直接将当前事件, 原封不动的, 给到全局。 会引发意外的报错。
      // document.dispatchEvent(event); // 报错原因是由于这个事件是已经分发过的事件。
    }

    //   - completed(已完成)   FIXME: 修复'Backspace'按键, 在录制过程中, 删除已选择列表中最后一项的bug
    if (event.key === 'Backspace') {
      // nextTick是为了确保selectedSingleKeys.value = oldSelectedSingleKeys的逻辑, 发生在元素被删除之后。
      nextTick(() => {
        selectedSingleKeys.value = oldSelectedSingleKeys;
      });
    }
  }

  // 更新oldSelectedSingleKeys, 以备下次使用
  oldSelectedSingleKeys = selectedSingleKeys.value.slice();
}

// -- -- 选择声效
const isDownSoundEffectSelectEnabled = ref(true);
const isUpSoundEffectSelectEnabled = ref(true);

const keyDownSingleKeySoundEffectSelect = ref<any>();
const keyUpSingleKeySoundEffectSelect = ref<any>();
const singleKeyTypeGroup = ref<Array<string>>(['sounds']);
const keySingleKeySoundEffectOptions = computed(() => {
  const List: Array<any> = [];
  if (singleKeyTypeGroup.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (singleKeyTypeGroup.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (singleKeyTypeGroup.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  console.debug('观察keyUnifiedSoundEffectOptions=', List);
  return List;
});
const isShowUltimatePerfectionKeySoundAnchoring_singleKey = computed(() => {
  return singleKeyTypeGroup.value.includes('key_sounds');
});
const isAnchoringUltimatePerfectionKeySound_singleKey = ref(true);
watch(keyDownSingleKeySoundEffectSelect, (newVal, oldVal) => {
  console.debug('观察keyDownSingleKeySoundEffectSelect=', keyDownSingleKeySoundEffectSelect.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (
      isShowUltimatePerfectionKeySoundAnchoring_singleKey.value &&
      isAnchoringUltimatePerfectionKeySound_singleKey.value
    ) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyUpSingleKeySoundEffectSelect.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyUpSingleKeySoundEffectSelect.value = keyDownSingleKeySoundEffectSelect.value;
      }
    }
  }
});
watch(keyUpSingleKeySoundEffectSelect, (newVal, oldVal) => {
  console.debug('观察keyUpSingleKeySoundEffectSelect=', keyUpSingleKeySoundEffectSelect.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (
      isShowUltimatePerfectionKeySoundAnchoring_singleKey.value &&
      isAnchoringUltimatePerfectionKeySound_singleKey.value
    ) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyDownSingleKeySoundEffectSelect.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyDownSingleKeySoundEffectSelect.value = keyUpSingleKeySoundEffectSelect.value;
      }
    }
  }
});

function saveSingleKeySoundEffectConfig(
  params: { singleKeys: Array<number>; down: any; up: any },
  onSuccess?: () => void
) {
  const keyTone_single = {
    down: {
      type: params.down?.type || '',
      // TIPS: 此处需要注意, 我们需要的是此lambda表达式执行后的返回值赋值给value, 而不是直接将lambda表达式赋值给value。(此处用三元表达式会更为直观)
      value: (() => {
        if (params.down?.type === 'audio_files') {
          return { sha256: params.down.value.sha256, name_id: params.down.value.name_id, type: params.down.value.type };
        }
        if (params.down?.type === 'sounds') {
          return params.down.value.soundKey;
        }
        if (params.down?.type === 'key_sounds') {
          return params.down.value.keySoundKey;
        }
        return '';
      })(),
    },
    up: {
      type: params.up?.type || '',
      // TIPS: 此处需要注意, 我们需要的是此lambda表达式执行后的返回值赋值给value, 而不是直接将lambda表达式赋值给value。(此处用三元表达式会更为直观)
      value: (() => {
        if (params.up?.type === 'audio_files') {
          return { sha256: params.up.value.sha256, name_id: params.up.value.name_id, type: params.up.value.type };
        }
        if (params.up?.type === 'sounds') {
          return params.up.value.soundKey;
        }
        if (params.up?.type === 'key_sounds') {
          return params.up.value.keySoundKey;
        }
        return '';
      })(),
    },
  };

  // 需要保证只有在对应的down或up声效设置被使能时, 才会修改对应的设置。(避免意外修改到不希望修改的配置)
  let downOrUpIfEnable: string;
  if (isDownSoundEffectSelectEnabled.value && !isUpSoundEffectSelectEnabled.value) {
    downOrUpIfEnable = '.down';
  } else if (!isDownSoundEffectSelectEnabled.value && isUpSoundEffectSelectEnabled.value) {
    downOrUpIfEnable = '.up';
  } else if (isDownSoundEffectSelectEnabled.value && isUpSoundEffectSelectEnabled.value) {
    downOrUpIfEnable = '';
  } else {
    // !isDownSoundEffectSelectEnabled.value && !isUpSoundEffectSelectEnabled.value
    // 此时, 说明用户未选择任何声效, 因此无需配置, 直接跳过。没必要继续执行后续的保存至 配置文件的步骤。
    q.notify({
      type: 'info',
      position: 'top',
      message: '未进行任何单键声效配置',
      timeout: 5000,
    });
    return;
  }

  let allSuccess = true;
  const promises = params.singleKeys.map((item) => {
    return ConfigSet(
      'key_tone.single.' + item + downOrUpIfEnable,
      // `downOrUpIfEnable.slice(1)` 是去掉字符串开头的 . 如  '.down' -> 'down'
      // `as keyof typeof keyTone_single` 是 TypeScript 类型断言，确保属性名是 `keyTone_single` 对象的有效键名
      downOrUpIfEnable ? keyTone_single[downOrUpIfEnable.slice(1) as keyof typeof keyTone_single] : keyTone_single
    )
      .then((re) => {
        if (!re) {
          allSuccess = false;
          q.notify({
            type: 'negative',
            position: 'top',
            message: `单键 ${
              keyEvent_store.dikCodeToName.get(item) ||
              (dikCodeToName_custom.get(item) ? 'Temp-{' + dikCodeToName_custom.get(item) + '}' : '') ||
              'Dik-{' + item + '}'
            } 声效配置失败`,
            timeout: 5000,
          });
        }
        return re;
      })
      .catch((err) => {
        allSuccess = false;
        console.error(`单键 ${item} 声效配置时发生错误:`, err);
        q.notify({
          type: 'negative',
          position: 'top',
          message: `单键 ${
            keyEvent_store.dikCodeToName.get(item) ||
            (dikCodeToName_custom.get(item) ? 'Temp-{' + dikCodeToName_custom.get(item) + '}' : '') ||
            'Dik-{' + item + '}'
          } 声效配置失败`,
          timeout: 5000,
        });
      });
  });

  Promise.all(promises).then(() => {
    if (allSuccess) {
      onSuccess?.();
    } else {
      q.notify({
        type: 'warning',
        position: 'top',
        message: '仅部分单键声效配置成功',
        timeout: 5000,
      });
    }
  });
}

// -- -- 查看编辑声效
// const keysWithSoundEffect = ref<string[]>([]);
// watch(
//   // TIPS: 注意, 如果要监听的ref对象是数组或对象等js/ts中的默认引用类型, 要使用此种方式才可触发监听。(或者也可以弃用ref, 改用reactive。)
//   () => keysWithSoundEffect.value,
//   (newVal) => {
//     console.debug('观察keysWithSoundEffect=', keysWithSoundEffect.value);
//   }
// );

const keysWithSoundEffect = ref<Map<string, any>>(new Map());
watch(
  // TIPS: 注意, 如果要监听的ref对象是数组或对象等js/ts中的默认引用类型, 要使用此种方式才可触发监听。(或者也可以弃用ref, 改用reactive。)
  () => keysWithSoundEffect.value,
  (newVal) => {
    console.debug('观察keysWithSoundEffect=', keysWithSoundEffect.value);
  }
);
const isShowSingleKeySoundEffectEditDialog = ref(false);

const currentEditingKey = ref<number | null>(null);
const currentEditingKeyOfName = computed(() => {
  return currentEditingKey.value !== null
    ? keyEvent_store.dikCodeToName.get(currentEditingKey.value) ||
        (dikCodeToName_custom.get(currentEditingKey.value)
          ? 'Temp-{' + dikCodeToName_custom.get(currentEditingKey.value) + '}'
          : '') ||
        'Dik-{' + currentEditingKey.value + '}'
    : '';
});

// -- -- -- 编辑声效(重新选择声效)
const isDownSoundEffectSelectEnabled_edit = ref(true);
const isUpSoundEffectSelectEnabled_edit = ref(true);

const keyDownSingleKeySoundEffectSelect_edit = ref<any>();
const keyUpSingleKeySoundEffectSelect_edit = ref<any>();
const singleKeyTypeGroup_edit = ref<Array<string>>(['sounds']);
const keySingleKeySoundEffectOptions_edit = computed(() => {
  const List: Array<any> = [];
  if (singleKeyTypeGroup_edit.value.includes('audio_files')) {
    soundFileList.value.forEach((item) => {
      List.push({ type: 'audio_files', value: item });
    });
  }
  if (singleKeyTypeGroup_edit.value.includes('sounds')) {
    soundList.value.forEach((item) => {
      List.push({ type: 'sounds', value: item });
    });
  }
  if (singleKeyTypeGroup_edit.value.includes('key_sounds')) {
    keySoundList.value.forEach((item) => {
      List.push({ type: 'key_sounds', value: item });
    });
  }
  console.debug('观察keyUnifiedSoundEffectOptions=', List);
  return List;
});
const isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit = computed(() => {
  return singleKeyTypeGroup_edit.value.includes('key_sounds');
});
const isAnchoringUltimatePerfectionKeySound_singleKey_edit = ref(true);
watch(keyDownSingleKeySoundEffectSelect_edit, (newVal, oldVal) => {
  console.debug('观察keyDownSingleKeySoundEffectSelect_edit=', keyDownSingleKeySoundEffectSelect_edit.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (
      isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit.value &&
      isAnchoringUltimatePerfectionKeySound_singleKey_edit.value
    ) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyUpSingleKeySoundEffectSelect_edit.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyUpSingleKeySoundEffectSelect_edit.value = keyDownSingleKeySoundEffectSelect_edit.value;
      }
    }
  }
});
watch(keyUpSingleKeySoundEffectSelect_edit, (newVal, oldVal) => {
  console.debug('观察keyUpSingleKeySoundEffectSelect_edit=', keyUpSingleKeySoundEffectSelect_edit.value);
  if (newVal === null && oldVal?.type === 'key_sounds') {
    if (
      isShowUltimatePerfectionKeySoundAnchoring_singleKey_edit.value &&
      isAnchoringUltimatePerfectionKeySound_singleKey_edit.value
    ) {
      // 这里?是为了防止其本身就为null时,访问不存在的type字段引发报错。(不需要处理null:我们本身就是对其做清空操作,没有值正好。)
      if (keyDownSingleKeySoundEffectSelect_edit.value?.type === 'key_sounds') {
        // 如果对方有值, 且值为key_sounds, 则清空。
        keyDownSingleKeySoundEffectSelect_edit.value = keyUpSingleKeySoundEffectSelect_edit.value;
      }
    }
  }
});

onBeforeMount(async () => {
  // 此时由于是新建键音包, 因此是没有对应配置文件, 需要我们主动去创建的。 故第二个参数设置为true
  // 这也是我们加载页面前必须确定的事情, 否则无法进行后续操作, 一切以配置文件为前提。
  await LoadConfig(pkgPath, true);

  // 使用i18n, 初始化键音包名称, 在此处手动设置, 是为了防止后续初始化将sdk中默认的名称, 被初始化到pkgName中。
  // 因此在初始化前按照提前设置好的i18n做为默认名称, 故手动发送请求以在数据初始化前更改sdk中的默认名称。
  await ConfigSet('package_name', $t('KeyTonePackage.new.name.defaultValue'));

  await ConfigSet('audio_pkg_uuid', pkgPath);

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

      // 已载入的声音文件列表初始化。  (不过由于这里是新建键音包, 这个不出意外的话一开始是undefine)
      if (data.audio_files !== undefined) {
        // keyTonePkgData.audio_files 是一个从后端获取的对象, 通过此方式可以简便的将其转换为数组, 数组元素为原对象中的key和value(增加了这两个key)
        const audioFilesArray = Object.entries(data.audio_files).map(([key, value]) => ({
          sha256: key,
          value: value,
        }));
        audioFiles.value = audioFilesArray;
      }

      // 已载入的声音列表初始化。  (不过由于这里是新建键音包, 这个不出意外的话一开始是undefine)
      if (data.sound_list !== undefined) {
        const sounds = Object.entries(data.sounds).map(([key, value]) => ({
          soundKey: key,
          soundValue: value,
        }));
        soundList.value = sounds as Array<{
          soundKey: string;
          soundValue: {
            cut: { start_time: number; end_time: number; volume: number };
            name: string;
            source_file_for_sound: { sha256: string; name_id: string; type: string };
          };
        }>;
      }

      // 后续极大可能会删除它(单键声效的Temp的单键名称)(TODO: 此逻辑未验证, 需要到编辑键音包界面才能验证)
      if (data.custom_single_keys_name !== undefined) {
        // 遍历 custom_single_keys_name 对象的每个键值对并设置到 dikCodeToName_custom 中
        Object.entries(data.custom_single_keys_name).forEach(([dikCode, name]) => {
          dikCodeToName_custom.set(Number(dikCode), name as string);
        });
      }

      // TODO: 此逻辑未验证, 需要到编辑键音包界面才能验证
      if (data.key_tone?.single !== undefined) {
        keysWithSoundEffect.value.clear();
        Object.entries(data.key_tone.single).forEach(([dikCode, value]) => {
          keysWithSoundEffect.value.set(dikCode, value);
        });
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
          Object.entries(item.value.name).forEach(([name_id, name]) => {
            tempSoundFileList.push({ sha256: item.sha256, name_id: name_id, name: name, type: item.value.type });
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
        if (selectedSoundFile.value.sha256 !== '' && selectedSoundFile.value.name_id !== '') {
          SoundFileRename(
            selectedSoundFile.value.sha256,
            selectedSoundFile.value.name_id,
            selectedSoundFile.value.name
          );
        }
      }
    );

    watch(soundList, (newVal) => {
      console.debug('观察soundList=', soundList.value);
    });

    watch(keySoundList, (newVal) => {
      console.debug('观察keySoundList=', keySoundList.value);
    });

    watch(
      isEnableEmbeddedTestSound,
      (newVal) => {
        ConfigSet('key_tone.is_enable_embedded_test_sound', newVal);
      },
      { immediate: true }
    );
  }

  const pkgNameDelayed = debounce(
    (keyTonePkgData: any) => {
      pkgName.value = keyTonePkgData.package_name;
    },
    800,
    { trailing: true }
  );

  const isEnableEmbeddedTestSoundDelayed = debounce(
    (val: { down: boolean; up: boolean }) => {
      isEnableEmbeddedTestSound.down = val.down;
      isEnableEmbeddedTestSound.up = val.up;
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

    // 映射配置文件中的sounds到ui中的soundList。(只要配置文件变更, 就会触发相关sse发送, 此处就会接收)
    if (keyTonePkgData.sounds !== undefined) {
      const sounds = Object.entries(keyTonePkgData.sounds).map(([key, value]) => ({
        soundKey: key,
        soundValue: value,
      }));
      soundList.value = sounds as Array<{
        soundKey: string;
        soundValue: {
          cut: { start_time: number; end_time: number; volume: number };
          name: string;
          source_file_for_sound: { sha256: string; name_id: string; type: string };
        };
      }>;
    } else {
      soundList.value = [];
    }

    // 映射配置文件中的key_sounds到ui中的keySoundList。(只要配置文件变更, 就会触发相关sse发送, 此处就会接收)
    if (keyTonePkgData.key_sounds !== undefined) {
      keySoundList.value = Object.entries(keyTonePkgData.key_sounds).map(([key, value]) => ({
        keySoundKey: key,
        keySoundValue: value,
      }));
    } else {
      keySoundList.value = [];
    }

    if (keyTonePkgData.key_tone !== undefined) {
      isEnableEmbeddedTestSoundDelayed.cancel();
      isEnableEmbeddedTestSoundDelayed(keyTonePkgData.key_tone.is_enable_embedded_test_sound);
    }

    // 后续极大可能会删除它(Temp的单键名称)
    // * 若仅依赖此逻辑, 添加声效时会因读取延迟使得第一时间仅能看到Dik码, 短暂的后续方可看到此项逻辑带来的Temp临时名称。
    // * 若想直接看到, 就改变逻辑中在ConfigSet前注释掉的dikCodeToName_custom.set相关逻辑。(启用)
    if (keyTonePkgData.custom_single_keys_name !== undefined) {
      // 遍历 custom_single_keys_name 对象的每个键值对并设置到 dikCodeToName_custom 中
      Object.entries(keyTonePkgData.custom_single_keys_name).forEach(([dikCode, name]) => {
        dikCodeToName_custom.set(Number(dikCode), name as string);
      });
    }

    if (keyTonePkgData.key_tone?.single !== undefined) {
      keysWithSoundEffect.value.clear();
      Object.entries(keyTonePkgData.key_tone.single).forEach(([dikCode, value]) => {
        keysWithSoundEffect.value.set(dikCode, value);
      });
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
