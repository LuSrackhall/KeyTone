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
  <q-page :style="{ '--i18n_fontSize': i18n_fontSize }">
    <q-scroll-area :class="[isMacOS ? 'w-[389.5px] h-[458.5px]' : 'w-[379px] h-[458.5px]']">
      <div :class="['flex flex-col gap-5  p-8  scale-102']">
        <q-input
          outlined
          stack-label
          dense
          :error-message="$t('KeyToneAlbum.new.name.errorMessage')"
          :error="pkgName === '' || pkgName === undefined || pkgName === null"
          v-model="pkgName"
          :label="$t('KeyToneAlbum.new.name.name')"
          :placeholder="$t('KeyToneAlbum.new.name.defaultValue')"
        />
        <!-- <div>原始声音文件编辑</div>
        <div>键音</div>
        <div>键音列表, 编辑键音</div>
        <div>全局键音规则</div>
        <div>对某个特定按键单独设置键音</div> -->

        <q-stepper v-model="step" vertical header-nav color="primary" animated class="step-custom">
          <div
            :class="[
              // 字体
              'font-semibold text-lg',
              // 对溢出的情况, 采取滚动策略
              'max-w-66 overflow-auto whitespace-nowrap text-nowrap',
              // 隐藏滚动策略的滚动条。
              '[&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none]',
              // 居中对齐
              'mx-auto',
            ]"
          >
            {{ pkgName }}
          </div>
          <q-step
            :name="1"
            :title="$t('KeyToneAlbum.loadAudioFile.title')"
            icon="create_new_folder"
            :done="soundFileList.length !== 0"
            :disable="step === 99 && soundFileList.length === 0"
            :header-nav="false"
            @click="
              (event:MouseEvent) => {
                // TIPS: 由于:header-nav=true是此组件默认的行为,
                //       他会为此组件添加默认的点击事件
                //       这个默认事件的作用是会自动将step其设置为1, 且发生在我们当前的click事件之前。
                //       因此若不手动将  :header-nav 设置为false, 我们就无法得到预期效果, 只能在点击关闭后望而却步
                //       故再次之前, 将 :header-nav 手动设置为false, 以禁用默认的事件
                //       而因为禁用默认 header-nav 后, 组件会失去可点击的样式, 因此我们手动添加相关的 class 可点击样式
                //       TODO: 由于我希望其关闭时, 可以简单的校验是否已完成, 因此我们使用disable和error两个组件状态来完成简单校验(目前仅使用disable, 后续完善时再说)
                // step = step === 99 ? 1 : 99;  // 由于还有其它步骤(如step=2,3,4等)要作此操作, 为了不影响, 我们更改判断方式为step===1,而不是现在的step===99。
                // step = step === 1 ? 99 : 1; // 我们只需要在内部组件的标题部分响应此事件, 因此需要先检查判断点击的区域。
                // 检查点击是否发生在标题区域
                const header = (event.target as HTMLElement).closest('.q-stepper__tab');
                if (header) {
                  step = step === 1 ? 99 : 1;
                }
                // q-step 组件是否为展开状态的逻辑很简单。 只要step的值等于其某个q-step的name字段即可。只要不相等则处于关闭状态, 只要相等就处于展开状态。
              }
            "
          >
            <div :class="['mb-3', step_introduce_fontSize]">{{ $t('KeyToneAlbum.loadAudioFile.description') }}</div>
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
                  :label="$t('KeyToneAlbum.loadAudioFile.addNewFile')"
                  @click="
                    () => {
                      addNewSoundFile = !addNewSoundFile;
                    }
                  "
                >
                </q-btn>
                <q-dialog
                  :style="{ '--i18n_fontSize': i18n_fontSize }"
                  v-model="addNewSoundFile"
                  backdrop-filter="invert(70%)"
                  @mouseup="preventDefaultMouseWhenRecording"
                >
                  <q-card>
                    <q-card-section class="row items-center q-pb-none text-h6">
                      {{ $t('KeyToneAlbum.loadAudioFile.addNewFile_1') }}
                    </q-card-section>

                    <!-- <q-card-section> <div>文件类型可以是WAV、MP3、OGG。</div></q-card-section> -->

                    <q-card-section>
                      <div class="text-gray-600 text-xs">{{ $t('KeyToneAlbum.loadAudioFile.dragAndDrop') }}</div>
                      <q-file
                        :class="['w-56', 'zl-ll']"
                        dense
                        v-model="files"
                        :label="$t('KeyToneAlbum.loadAudioFile.audioFile')"
                        outlined
                        use-chips
                        multiple
                        append
                        accept=".wav,.mp3,.ogg"
                        excludeAcceptAllOption
                        style="max-width: 300px"
                        :hint="$t('KeyToneAlbum.loadAudioFile.supportedFormats')"
                      />
                    </q-card-section>

                    <q-card-section>
                      <div>{{ $t('KeyToneAlbum.loadAudioFile.addAsNeeded') }}</div>
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

                                    message: `${$t('KeyToneAlbum.notify.addFailed')} '${file.name}'`,

                                    timeout: 5,
                                  });
                                  return;
                                }
                              } catch (error) {
                                console.error(`Error uploading file ${file.name}:`, error);
                                q.notify({
                                  type: 'negative',

                                  position: 'top',

                                  message: `${$t('KeyToneAlbum.notify.addFailed')} '${file.name}'`,

                                  timeout: 5,
                                });
                                return;
                              }
                            }
                            nextTick(() => {
                              q.notify({
                                type: 'positive',
                                position: 'top',
                                message: $t('KeyToneAlbum.notify.addSuccess'),
                                timeout: 5,
                              });
                            });
                          }
                        "
                        color="primary"
                        :label="$t('KeyToneAlbum.loadAudioFile.confirmAdd')"
                      />
                      <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
              <div :class="['p-2 text-zinc-600']">{{ $t('KeyToneAlbum.or') }}</div>
              <!-- -------------------------------------------------------------------------------编辑已有音频源文件 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  :label="$t('KeyToneAlbum.loadAudioFile.manageExistingFiles')"
                  @click="
                    () => {
                      if (soundFileList.length === 0) {
                        q.notify({
                          type: 'warning',
                          message: $t('KeyToneAlbum.notify.noFilesToManage'),
                          position: 'top',
                        });
                        return;
                      }

                      editSoundFile = !editSoundFile;
                    }
                  "
                ></q-btn>
                <q-dialog
                  :style="{ '--i18n_fontSize': i18n_fontSize }"
                  v-model="editSoundFile"
                  backdrop-filter="invert(70%)"
                  @mouseup="preventDefaultMouseWhenRecording"
                >
                  <q-card :class="['p-x-3  w-[96%]']">
                    <q-card-section class="row items-center q-pb-none text-h6">
                      {{ $t('KeyToneAlbum.loadAudioFile.manageExistingFiles') }}
                    </q-card-section>

                    <!-- <q-card-section> <div>请选择您想要修改或删除的声音源文件并执行对应操作。</div></q-card-section> -->

                    <q-card-section>
                      <q-select
                        outlined
                        :virtual-scroll-slice-size="999999"
                        stack-label
                        v-model="selectedSoundFile"
                        :options="soundFileList"
                        :option-label="(item) => item.name + item.type"
                        :label="$t('KeyToneAlbum.loadAudioFile.selectFileToManage')"
                        dense
                        popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                      >
                        <!-- 添加清除按钮 -->
                        <template
                          v-if="selectedSoundFile.sha256 !== '' && selectedSoundFile.name_id !== ''"
                          v-slot:append
                        >
                          <q-icon
                            name="cancel"
                            @click.stop.prevent="
                              selectedSoundFile = {
                                sha256: '',
                                name_id: '',
                                name: '',
                                type: '',
                              }
                            "
                            class="cursor-pointer text-lg"
                          />
                        </template>
                      </q-select>
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
                      <!-- option-value="uuid" WARN -> 这个理解似乎是不对的
                       如果不手动的显示指定, 默认会以整个对象作为区分标准(我们这里就是这种情况, 因为我们没有能作为uuid的字段),
                       但也可以指定具体字段作为区分标准( 可以类比v-for遍历时指定的 :key 值。),
                       当检测到变化, 才会更新已选择框内的label名。
                       这里我们不手动指定, 默认使用整个对象做区分,
                       因为我们的uuid可能和其它不同sha256中的uuid重复,
                       而sha256就更不用说了, 默认就会有相同sha256的项,
                       至于使用name也更不用说, 我们命名时允许相同的名称。
                       使用type这个字段作为唯一表示就更加没有讨论意义了。
                        -->
                      <!-- option-value所返回的值, 似乎会影响真正的v-model中所绑定的值 WARN -> 与上述的警告互相参考, 待确定
                        也就是说, select选择框采用的 v-model 的值, 就是option-value中的返回值, 这个值在不指定的情况下默认为options数组元素的对象中的value值, 如果没对象内没value字段则默认返回整个对象。
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
                          :label="selectedSoundFile.type"
                          :class="[
                            'absolute  overflow-visible ',
                            // 'left-0',
                            'right-0',
                          ]"
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
                            :error-message="$t('KeyToneAlbum.notify.emptyFileName')"
                            :error="
                              selectedSoundFile.name === '' ||
                              selectedSoundFile.name === undefined ||
                              selectedSoundFile.name === null
                            "
                            v-model="selectedSoundFile.name"
                            :label="$t('KeyToneAlbum.loadAudioFile.renameFile')"
                          />

                          <q-btn
                            :class="['w-20 self-center bg-pink-700 text-zinc-50']"
                            dense
                            no-caps
                            :label="$t('KeyToneAlbum.delete')"
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
                                    message: $t('KeyToneAlbum.notify.deleteSuccess'),
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
                                    message: $t('KeyToneAlbum.notify.deleteFailed'),
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
                      <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
            </div>
            <!-- ------------------------------------------------------------------------载入声音文件的业务逻辑   end -->
            <q-stepper-navigation>
              <q-btn @click="step = 2" color="primary" :label="$t('KeyToneAlbum.continue')" />
            </q-stepper-navigation>
          </q-step>

          <!-- <q-step :name="2" title="键音制作" caption="Optional" icon="create_new_folder" :done="step > 2"> -->
          <q-step
            :name="2"
            :title="$t('KeyToneAlbum.defineSounds.title')"
            icon="add_comment"
            :done="soundList.length !== 0"
            :disable="step === 99 && soundList.length === 0"
            :header-nav="false"
            @click="
              (event: MouseEvent) => {
                // 检查点击是否发生在标题区域
                const header = (event.target as HTMLElement).closest('.q-stepper__tab');
                if (header) {
                  step = step === 2 ? 99 : 2;
                }
              }
            "
          >
            <div :class="['mb-3', step_introduce_fontSize]">
              {{ $t('KeyToneAlbum.defineSounds.description') }}
              <q-icon name="info" color="primary">
                <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                  <div>{{ $t('KeyToneAlbum.defineSounds.tooltip.noImpactOnSource') }}</div>
                  <div>{{ $t('KeyToneAlbum.defineSounds.tooltip.multipleSoundsFromSameSource') }}</div>
                </q-tooltip>
              </q-icon>
            </div>

            <!-- ------------------------------------------------------------------------裁剪定义声音的业务逻辑 start -->
            <div>
              <!-- ------------------------------------------------------------------------------ 制作新的声音 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  :label="$t('KeyToneAlbum.defineSounds.createNewSound')"
                  @click="
                    () => {
                      createNewSound = !createNewSound;
                    }
                  "
                >
                </q-btn>
                <q-dialog
                  :style="{ '--i18n_fontSize': i18n_fontSize }"
                  v-model="createNewSound"
                  backdrop-filter="invert(70%)"
                  @mouseup="preventDefaultMouseWhenRecording"
                >
                  <q-card>
                    <q-card-section class="row items-center q-pb-none text-h6">
                      {{ $t('KeyToneAlbum.defineSounds.createNewSound') }}
                    </q-card-section>
                    <q-card-section :class="['p-b-1']">
                      <q-input
                        outlined
                        stack-label
                        dense
                        v-model="soundName"
                        :label="$t('KeyToneAlbum.defineSounds.soundName')"
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
                                $t('KeyToneAlbum.defineSounds.tooltip.soundName') +
                                ' : \n' +
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
                        :virtual-scroll-slice-size="999999"
                        v-model="sourceFileForSound"
                        :options="soundFileList"
                        :option-label="(item) => item.name + item.type"
                        :label="$t('KeyToneAlbum.defineSounds.sourceFile')"
                        dense
                        popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                      />
                    </q-card-section>

                    <q-card-section :class="['p-b-1']">
                      <div class="text-[13.5px] text-gray-600 p-b-2">
                        {{ $t('KeyToneAlbum.defineSounds.cropSound') }}
                        <q-icon name="info" color="primary">
                          <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                            {{ $t('KeyToneAlbum.defineSounds.tooltip.soundDuration') }}
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
                          :label="$t('KeyToneAlbum.defineSounds.startTime')"
                          type="number"
                          :error-message="$t('KeyToneAlbum.defineSounds.error.negativeTime')"
                          :error="soundStartTime < 0"
                        />
                        <q-input
                          :class="['w-1/2 p-l-1']"
                          outlined
                          stack-label
                          dense
                          v-model.number="soundEndTime"
                          :label="$t('KeyToneAlbum.defineSounds.endTime')"
                          type="number"
                          :error-message="$t('KeyToneAlbum.defineSounds.error.negativeTime')"
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
                        :label="$t('KeyToneAlbum.defineSounds.volume')"
                        type="number"
                        :step="0.1"
                      >
                        <template v-slot:append>
                          <q-icon name="info" color="primary">
                            <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                              {{ $t('KeyToneAlbum.defineSounds.tooltip.volume') }}
                              {{ $t('KeyToneAlbum.defineSounds.tooltip.volume_1') }}
                            </q-tooltip>
                          </q-icon>
                        </template>
                      </q-input>
                    </q-card-section>
                    <q-card-actions align="right">
                      <q-btn
                        class="mt-2"
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
                        :label="$t('KeyToneAlbum.defineSounds.previewSound')"
                        color="secondary"
                      >
                        <q-tooltip
                          :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-xs']"
                          :delay="600"
                        >
                          {{ $t('KeyToneAlbum.defineSounds.tooltip.previewSound') }}
                        </q-tooltip>
                      </q-btn>
                      <q-btn
                        class="mt-2"
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
                        :label="$t('KeyToneAlbum.defineSounds.confirmAdd')"
                        color="primary"
                      />
                      <q-btn class="mt-2" flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>

              <div :class="['p-2 text-zinc-600']">{{ $t('KeyToneAlbum.or') }}</div>

              <!-- -------------------------------------------------------------------------------编辑已有声音 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  :label="$t('KeyToneAlbum.defineSounds.editExistingSound')"
                  @click="
                    () => {
                      if (soundList.length === 0) {
                        q.notify({
                          type: 'warning',
                          message: $t('KeyToneAlbum.notify.noSoundsToEdit'),
                          position: 'top',
                        });
                        return;
                      }
                      showEditSoundDialog = true;
                    }
                  "
                >
                </q-btn>
                <q-dialog
                  :style="{ '--i18n_fontSize': i18n_fontSize }"
                  v-model="showEditSoundDialog"
                  backdrop-filter="invert(70%)"
                  @mouseup="preventDefaultMouseWhenRecording"
                >
                  <q-card class="min-w-[106%]">
                    <q-card-section
                      class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                    >
                      {{ $t('KeyToneAlbum.defineSounds.editExistingSound') }}
                    </q-card-section>
                    <q-card-section>
                      <q-select
                        outlined
                        stack-label
                        :virtual-scroll-slice-size="999999"
                        clearable
                        v-model="selectedSound"
                        popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                        :options="soundList"
                        :option-label="(item: any) => {
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
                        :label="$t('KeyToneAlbum.defineSounds.selectSoundToManage')"
                        dense
                      >
                        <template v-slot:option="scope">
                          <q-item v-bind="scope.itemProps">
                            <q-item-section>
                              <q-item-label>
                                {{
                                  scope.opt.soundValue.name !== '' && scope.opt.soundValue.name !== undefined
                                    ? scope.opt.soundValue.name
                                    : soundFileList.find(
                                        (soundFile) =>
                                          soundFile.sha256 === scope.opt.soundValue.source_file_for_sound.sha256 &&
                                          soundFile.name_id === scope.opt.soundValue.source_file_for_sound.name_id
                                      )?.name + '     - ' + ' [' + scope.opt.soundValue.cut.start_time + ' ~ ' + scope.opt.soundValue.cut.end_time + ']'
                                }}
                              </q-item-label>
                            </q-item-section>
                            <q-item-section side>
                              <DependencyWarning
                                :issues="dependencyIssues"
                                item-type="sounds"
                                :item-id="scope.opt.soundKey"
                                :show-details="false"
                              />
                            </q-item-section>
                          </q-item>
                        </template>
                      </q-select>
                    </q-card-section>
                    <!-- 以卡片形式展示选择的声音 -->
                    <q-card-section
                      :class="['flex flex-col m-t-3']"
                      v-if="selectedSound?.soundKey !== '' && selectedSound !== undefined"
                    >
                      <q-card :class="['flex flex-col pb-3 w-[100%]']" v-if="selectedSound">
                        <q-card-section :class="['p-b-1 mt-3']">
                          <q-input
                            outlined
                            stack-label
                            dense
                            v-model="selectedSound.soundValue.name"
                            :label="$t('KeyToneAlbum.defineSounds.soundName')"
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
                                    $t('KeyToneAlbum.defineSounds.tooltip.soundName') +
                                    ' : \n' +
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
                        <q-card-section :class="['p-b-1 w-68']">
                          <q-select
                            outlined
                            stack-label
                            :virtual-scroll-slice-size="999999"
                            popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
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
                            :label="$t('KeyToneAlbum.defineSounds.sourceFile')"
                            dense
                          />
                        </q-card-section>

                        <q-card-section :class="['p-b-1']">
                          <div class="text-[13.5px] text-gray-600 p-b-2">
                            {{ $t('KeyToneAlbum.defineSounds.cropSound') }}
                            <q-icon name="info" color="primary">
                              <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']">
                                {{ $t('KeyToneAlbum.defineSounds.tooltip.soundDuration') }}
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
                              :label="$t('KeyToneAlbum.defineSounds.startTime')"
                              type="number"
                              :error-message="$t('KeyToneAlbum.defineSounds.error.negativeTime')"
                              :error="selectedSound.soundValue.cut.start_time < 0"
                            />
                            <q-input
                              :class="['w-1/2 p-l-1']"
                              outlined
                              stack-label
                              dense
                              v-model.number="selectedSound.soundValue.cut.end_time"
                              :label="$t('KeyToneAlbum.defineSounds.endTime')"
                              type="number"
                              :error-message="$t('KeyToneAlbum.defineSounds.error.negativeTime')"
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
                            :label="$t('KeyToneAlbum.defineSounds.volume')"
                            type="number"
                            :step="0.1"
                          >
                            <template v-slot:append>
                              <q-icon name="info" color="primary">
                                <q-tooltip
                                  :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words']"
                                >
                                  {{ $t('KeyToneAlbum.defineSounds.tooltip.volume') }}
                                  {{ $t('KeyToneAlbum.defineSounds.tooltip.volume_1') }}
                                </q-tooltip>
                              </q-icon>
                            </template>
                          </q-input>
                        </q-card-section>

                        <!-- 添加按钮组 -->
                        <q-card-section :class="['flex justify-center gap-4']">
                          <q-btn
                            class="pr-2.3"
                            dense
                            color="secondary"
                            icon="play_arrow"
                            :label="$t('KeyToneAlbum.defineSounds.previewSound')"
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
                              {{ $t('KeyToneAlbum.defineSounds.tooltip.previewSound') }}
                            </q-tooltip>
                          </q-btn>

                          <q-btn
                            class="pr-2.3"
                            dense
                            color="primary"
                            icon="save"
                            :label="$t('KeyToneAlbum.confirmEdit')"
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
                            class="pr-2.3"
                            dense
                            color="negative"
                            icon="delete"
                            :label="$t('KeyToneAlbum.delete')"
                            @click="
                              deleteSound({
                                soundKey: selectedSound.soundKey,
                                onSuccess: () => {
                                  selectedSound = undefined;
                                  q.notify({
                                    type: 'positive',
                                    position: 'top',
                                    message: $t('KeyToneAlbum.notify.deleteSuccess'),
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
                      <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
            </div>
            <!-- ------------------------------------------------------------------------裁剪定义声音的业务逻辑   end -->
            <q-stepper-navigation>
              <q-btn @click="step = 3" color="primary" :label="$t('KeyToneAlbum.continue')" />
              <q-btn flat @click="step = 1" color="primary" :label="$t('KeyToneAlbum.back')" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>

          <q-step
            :name="3"
            :title="$t('KeyToneAlbum.craftKeySounds.title')"
            icon="add_comment"
            :done="keySoundList.length !== 0"
            :disable="step === 99 && keySoundList.length === 0"
            :header-nav="false"
            @click="
              (event: MouseEvent) => {
                // if (soundList.length === 0) {
                //   q.notify({
                //     type: 'warning',
                //     position: 'top',
                //     message: '请先定义声音',
                //     timeout: 5,
                //   });
                //   return;
                // }
                // 检查点击是否发生在标题区域
                const header = (event.target as HTMLElement).closest('.q-stepper__tab');
                if (header) {
                  step = step === 3 ? 99 : 3;
                }
              }
            "
          >
            <div :class="['mb-3', step_introduce_fontSize]">
              {{ $t('KeyToneAlbum.craftKeySounds.description') }}
              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                  <span>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.description_0') }}</span>
                  <br />
                  <span>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.description_1') }}</span>
                  <br />
                  <span>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.description_2') }}</span>
                  <br />
                  <span>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.description_3') }}</span>
                </q-tooltip>
              </q-icon>
            </div>
            <!-- ------------------------------------------------------------------------制作按键声音的业务逻辑 start -->
            <div>
              <!-- ------------------------------------------------------------------------------ 制作新的按键音 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  :label="$t('KeyToneAlbum.craftKeySounds.newKeySound')"
                  @click="
                    () => {
                      createNewKeySound = !createNewKeySound;
                    }
                  "
                >
                </q-btn>
                <q-dialog
                  :style="{ '--i18n_fontSize': i18n_fontSize }"
                  v-model="createNewKeySound"
                  backdrop-filter="invert(70%)"
                  @mouseup="preventDefaultMouseWhenRecording"
                >
                  <q-card :class="['min-w-[90%]']">
                    <q-card-section class="row items-center q-pb-none text-h6">
                      {{ $t('KeyToneAlbum.craftKeySounds.newKeySound') }}
                    </q-card-section>
                    <q-card-section :class="['p-b-1']">
                      <q-input
                        outlined
                        stack-label
                        dense
                        v-model="keySoundName"
                        :label="$t('KeyToneAlbum.craftKeySounds.keySoundName')"
                        :placeholder="$t('KeyToneAlbum.craftKeySounds.keySoundName-placeholder')"
                      />
                      <div class="flex flex-col mt-3">
                        <q-btn
                          :class="['bg-zinc-300 my-7 w-88% self-center']"
                          :label="$t('KeyToneAlbum.craftKeySounds.configureDownSound')"
                          @click="configureDownSound = true"
                        />
                        <q-dialog
                          :style="{ '--i18n_fontSize': i18n_fontSize }"
                          v-model="configureDownSound"
                          backdrop-filter="invert(70%)"
                          @mouseup="preventDefaultMouseWhenRecording"
                        >
                          <q-card :class="['min-w-[80%]']">
                            <q-card-section class="row items-center q-pb-none text-h6">
                              {{ $t('KeyToneAlbum.craftKeySounds.configureDownSound') }}
                            </q-card-section>
                            <q-card-section>
                              <!-- 使用选择框选择模式 -->
                              <q-select
                                outlined
                                popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                :virtual-scroll-slice-size="999999"
                                stack-label
                                v-model="playModeForDown"
                                :options="playModeOptions"
                                :option-label="(item: any) => {
                                  return $t(playModeLabels.get(item) || '')
                                }"
                                :label="$t('KeyToneAlbum.craftKeySounds.selectPlayMode')"
                                dense
                              />
                            </q-card-section>
                            <q-card-section>
                              <!-- 选择声音的选项，支持多选 -->
                              <q-select
                                outlined
                                stack-label
                                :virtual-scroll-slice-size="999999"
                                v-model="selectedSoundsForDown"
                                popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                :options="downSoundList"
                                :option-label="album_options_select_label"
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
                                :label="$t('KeyToneAlbum.craftKeySounds.selectSounds')"
                                multiple
                                use-chips
                                :class="['zl-ll']"
                                dense
                                :max-values="maxSelectionForDown"
                                counter
                                :error-message="$t('KeyToneAlbum.craftKeySounds.error.singleMode')"
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
                                  playModeForDown === 'single'
                                    ? selectedSoundsForDown.length > 1
                                      ? true
                                      : false
                                    : false
                                "
                                ref="downSoundSelectDom"
                                @update:model-value="downSoundSelectDom?.hidePopup()"
                              >
                                <template v-slot:option="scope">
                                  <q-item v-bind="scope.itemProps">
                                    <q-item-section>
                                      <q-item-label>{{ album_options_select_label(scope.opt) }}</q-item-label>
                                    </q-item-section>
                                    <q-item-section side>
                                      <DependencyWarning
                                        v-if="scope.opt.type === 'sounds'"
                                        :issues="dependencyIssues"
                                        item-type="sounds"
                                        :item-id="scope.opt.value.soundKey"
                                        :show-details="false"
                                      />
                                      <DependencyWarning
                                        v-else-if="scope.opt.type === 'key_sounds'"
                                        :issues="dependencyIssues"
                                        item-type="key_sounds"
                                        :item-id="scope.opt.value.keySoundKey"
                                        :show-details="false"
                                      />
                                    </q-item-section>
                                  </q-item>
                                </template>
                              </q-select>
                              <div class="h-10">
                                <q-option-group
                                  dense
                                  v-model="downTypeGroup"
                                  :options="options"
                                  type="checkbox"
                                  class="absolute left-8"
                                >
                                  <template #label-0="props">
                                    <q-item-label>
                                      {{ $t(props.label) }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.audioFile') }}</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                  <template v-slot:label-1="props">
                                    <q-item-label>
                                      {{ $t(props.label) }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.soundList') }}</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                  <template v-slot:label-2="props">
                                    <q-item-label>
                                      {{ $t(props.label) }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.keySounds') }}</div>
                                          <div>⬇</div>
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.inheritKeySound') }}</div>
                                          <div>⬇</div>
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.inheritRule') }}</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                </q-option-group>
                              </div>
                            </q-card-section>
                            <q-card-actions align="right">
                              <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                            </q-card-actions>
                          </q-card>
                        </q-dialog>
                        <q-btn
                          :class="['bg-zinc-300  m-b-7 w-88% self-center']"
                          :label="$t('KeyToneAlbum.craftKeySounds.configureUpSound')"
                          @click="configureUpSound = true"
                        />
                        <q-dialog
                          :style="{ '--i18n_fontSize': i18n_fontSize }"
                          v-model="configureUpSound"
                          backdrop-filter="invert(70%)"
                          @mouseup="preventDefaultMouseWhenRecording"
                        >
                          <q-card :class="['min-w-[80%]']">
                            <q-card-section class="row items-center q-pb-none text-h6">
                              {{ $t('KeyToneAlbum.craftKeySounds.configureUpSound') }}
                            </q-card-section>
                            <q-card-section>
                              <!-- 使用选择框选择模式 -->
                              <q-select
                                outlined
                                stack-label
                                :virtual-scroll-slice-size="999999"
                                popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                v-model="playModeForUp"
                                :options="playModeOptions"
                                :option-label="(item: any) => {
                                  return $t(playModeLabels.get(item) || '')
                                }"
                                :label="$t('KeyToneAlbum.craftKeySounds.selectPlayMode')"
                                dense
                              />
                            </q-card-section>
                            <q-card-section>
                              <!-- 选择声音的选项，支持多选 -->
                              <q-select
                                outlined
                                stack-label
                                :virtual-scroll-slice-size="999999"
                                v-model="selectedSoundsForUp"
                                popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                :options="upSoundList"
                                :option-label="album_options_select_label"
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
                                :label="$t('KeyToneAlbum.craftKeySounds.selectSounds')"
                                multiple
                                use-chips
                                :class="['zl-ll']"
                                dense
                                :max-values="maxSelectionForUp"
                                counter
                                :error-message="$t('KeyToneAlbum.craftKeySounds.error.singleMode')"
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
                                  playModeForUp === 'single' ? (selectedSoundsForUp.length > 1 ? true : false) : false
                                "
                                ref="upSoundSelectDom"
                                @update:model-value="upSoundSelectDom?.hidePopup()"
                              />
                              <div class="h-10">
                                <q-option-group
                                  dense
                                  v-model="upTypeGroup"
                                  :options="options"
                                  type="checkbox"
                                  class="absolute left-8"
                                >
                                  <template #label-0="props">
                                    <q-item-label>
                                      {{ $t(props.label) }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.audioFile') }}</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                  <template v-slot:label-1="props">
                                    <q-item-label>
                                      {{ $t(props.label) }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.soundList') }}</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                  <template v-slot:label-2="props">
                                    <q-item-label>
                                      {{ $t(props.label) }}
                                      <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.keySounds') }}</div>
                                          <div>⬇</div>
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.inheritKeySound') }}</div>
                                          <div>⬇</div>
                                          <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.inheritRule') }}</div>
                                        </q-tooltip>
                                      </q-icon>
                                    </q-item-label>
                                  </template>
                                </q-option-group>
                              </div>
                            </q-card-section>
                            <q-card-actions align="right">
                              <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                            </q-card-actions>
                          </q-card>
                        </q-dialog>
                      </div>
                    </q-card-section>
                    <q-card-actions align="right">
                      <q-btn
                        color="primary"
                        :label="$t('KeyToneAlbum.craftKeySounds.confirmAdd')"
                        @click="
                          saveKeySoundConfig(
                            {
                              key: '',
                              name: keySoundName,
                              down: { mode: playModeForDown, value: selectedSoundsForDown },
                              up: { mode: playModeForUp, value: selectedSoundsForUp },
                            },
                            () => {
                              // 关闭对话框
                              createNewKeySound = !createNewKeySound;

                              // 重置表单变量
                              keySoundName = $t('KeyToneAlbum.craftKeySounds.keySoundName-placeholder');
                              selectedSoundsForDown = [];
                              playModeForDown = 'random';
                              selectedSoundsForUp = [];
                              playModeForUp = 'random';
                            }
                          )
                        "
                      />
                      <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>

              <div :class="['p-2 text-zinc-600']">{{ $t('KeyToneAlbum.or') }}</div>

              <!-- -------------------------------------------------------------------------------编辑已有按键音 -->
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  :label="$t('KeyToneAlbum.craftKeySounds.editExistingKeySound')"
                  @click="
                    () => {
                      if (keySoundList.length === 0) {
                        q.notify({
                          type: 'warning',
                          message: $t('KeyToneAlbum.notify.noKeySoundsToEdit'),
                          position: 'top',
                        });
                        return;
                      }
                      editExistingKeySound = true;
                    }
                  "
                >
                </q-btn>
                <q-dialog
                  :style="{ '--i18n_fontSize': i18n_fontSize }"
                  v-model="editExistingKeySound"
                  backdrop-filter="invert(70%)"
                  @mouseup="preventDefaultMouseWhenRecording"
                >
                  <q-card :class="['min-w-[100%]']">
                    <q-card-section
                      class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                    >
                      {{ $t('KeyToneAlbum.craftKeySounds.editExistingKeySound') }}
                    </q-card-section>
                    <q-card-section>
                      <q-select
                        outlined
                        stack-label
                        clearable
                        :virtual-scroll-slice-size="999999"
                        popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                        v-model="selectedKeySound"
                        :options="keySoundList"
                        :label="$t('KeyToneAlbum.craftKeySounds.selectKeySoundToEdit')"
                        :option-label="(item) => item.keySoundValue.name"
                        :option-value="(item) => item.keySoundKey"
                        dense
                      >
                        <template v-slot:option="scope">
                          <q-item v-bind="scope.itemProps">
                            <q-item-section>
                              <q-item-label>{{ scope.opt.keySoundValue.name }}</q-item-label>
                            </q-item-section>
                            <q-item-section side>
                              <DependencyWarning
                                :issues="dependencyIssues"
                                item-type="key_sounds"
                                :item-id="scope.opt.keySoundKey"
                                :show-details="false"
                              />
                            </q-item-section>
                          </q-item>
                        </template>
                      </q-select>
                    </q-card-section>
                    <!-- 以卡片的形式, 展示选择的按键音 -->
                    <q-card-section
                      :class="['flex flex-col -m-t-2']"
                      v-if="selectedKeySound?.keySoundKey !== '' && selectedKeySound !== undefined"
                    >
                      <q-card :class="['flex flex-col pb-3']" v-if="selectedKeySound">
                        <q-card-section :class="['p-b-1 mt-3']">
                          <q-input
                            outlined
                            stack-label
                            dense
                            v-model="selectedKeySound.keySoundValue.name"
                            :label="$t('KeyToneAlbum.craftKeySounds.keySoundName')"
                            :placeholder="$t('KeyToneAlbum.craftKeySounds.keySoundName-placeholder')"
                          />
                          <div class="flex flex-col mt-1">
                            <q-btn
                              :class="['bg-zinc-300 my-7 w-88% self-center']"
                              :label="$t('KeyToneAlbum.craftKeySounds.configureDownSound')"
                              @click="edit_configureDownSound = true"
                            />
                            <q-dialog
                              :style="{ '--i18n_fontSize': i18n_fontSize }"
                              v-model="edit_configureDownSound"
                              backdrop-filter="invert(70%)"
                              @mouseup="preventDefaultMouseWhenRecording"
                            >
                              <q-card :class="['min-w-[80%]']">
                                <q-card-section class="row items-center q-pb-none text-h6">
                                  {{ $t('KeyToneAlbum.craftKeySounds.configureDownSound') }}
                                </q-card-section>
                                <q-card-section>
                                  <q-select
                                    outlined
                                    stack-label
                                    :virtual-scroll-slice-size="999999"
                                    popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                    v-model="selectedKeySound.keySoundValue.down.mode"
                                    :options="playModeOptions"
                                    :option-label="(item: any) => {
                                      return $t(playModeLabels.get(item) || '')
                                    }"
                                    :label="$t('KeyToneAlbum.craftKeySounds.selectPlayMode')"
                                    dense
                                  />
                                </q-card-section>
                                <q-card-section>
                                  <q-select
                                    outlined
                                    stack-label
                                    :virtual-scroll-slice-size="999999"
                                    popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                    v-model="selectedKeySound.keySoundValue.down.value"
                                    :options="edit_downSoundList"
                                    :option-label="album_options_select_label"
                                    :option-value="
                                      (item) => {
                                        /**
                                         * json中的存储格式分别是
                                         * {key:'audio_files', value:{sha256: string, name_id: string, type:string}}
                                         * {key:'sounds', value:string} // 此处value, 是soundKey
                                         * {key:'key_sounds', value:string} // 此处value, 是keySoundKey
                                         */
                                        if (item.type === 'audio_files') {
                                          return item.value?.sha256 + item.value?.name_id;
                                        }
                                        if (item.type === 'sounds') {
                                          return item.value?.soundKey;
                                        }
                                        if (item.type === 'key_sounds') {
                                          return item.value?.keySoundKey;
                                        }
                                      }
                                    "
                                    :label="$t('KeyToneAlbum.craftKeySounds.selectSounds')"
                                    multiple
                                    use-chips
                                    :class="['zl-ll']"
                                    dense
                                    :max-values="
                                      selectedKeySound.keySoundValue.down.mode.mode === 'single' ? 1 : Infinity
                                    "
                                    counter
                                    :error-message="$t('KeyToneAlbum.craftKeySounds.error.singleMode')"
                                    :error="
                                      selectedKeySound.keySoundValue.down.mode.mode === 'single' &&
                                      selectedKeySound.keySoundValue.down.value.length > 1
                                    "
                                    ref="edit_downSoundSelectDom"
                                    @update:model-value="edit_downSoundSelectDom?.hidePopup()"
                                  />
                                  <div class="h-10">
                                    <q-option-group
                                      dense
                                      v-model="edit_downTypeGroup"
                                      :options="options"
                                      type="checkbox"
                                      class="absolute left-8"
                                    >
                                      <template #label-0="props">
                                        <q-item-label>
                                          {{ $t(props.label) }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.audioFile') }}</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-1="props">
                                        <q-item-label>
                                          {{ $t(props.label) }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.soundList') }}</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-2="props">
                                        <q-item-label>
                                          {{ $t(props.label) }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.keySounds') }}</div>
                                              <div>⬇</div>
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.inheritKeySound') }}</div>
                                              <div>⬇</div>
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.inheritRule') }}</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                    </q-option-group>
                                  </div>
                                </q-card-section>
                                <q-card-actions align="right">
                                  <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                                </q-card-actions>
                              </q-card>
                            </q-dialog>
                            <q-btn
                              :class="['bg-zinc-300 m-b-7 w-88% self-center']"
                              :label="$t('KeyToneAlbum.craftKeySounds.configureUpSound')"
                              @click="edit_configureUpSound = true"
                            />
                            <q-dialog
                              :style="{ '--i18n_fontSize': i18n_fontSize }"
                              v-model="edit_configureUpSound"
                              backdrop-filter="invert(70%)"
                              @mouseup="preventDefaultMouseWhenRecording"
                            >
                              <q-card :class="['min-w-[80%]']">
                                <q-card-section class="row items-center q-pb-none text-h6">
                                  {{ $t('KeyToneAlbum.craftKeySounds.configureUpSound') }}
                                </q-card-section>
                                <q-card-section>
                                  <q-select
                                    outlined
                                    stack-label
                                    :virtual-scroll-slice-size="999999"
                                    popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                    v-model="selectedKeySound.keySoundValue.up.mode"
                                    :options="playModeOptions"
                                    :option-label="(item: any) => {
                                      return $t(playModeLabels.get(item) || '')
                                    }"
                                    :label="$t('KeyToneAlbum.craftKeySounds.selectPlayMode')"
                                    dense
                                  />
                                </q-card-section>
                                <q-card-section>
                                  <q-select
                                    outlined
                                    stack-label
                                    :virtual-scroll-slice-size="999999"
                                    popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                    v-model="selectedKeySound.keySoundValue.up.value"
                                    :options="edit_upSoundList"
                                    :option-label="album_options_select_label"
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
                                          return item.value?.sha256 + item.value?.name_id;
                                        }
                                        if (item.type === 'sounds') {
                                          return item.value?.soundKey;
                                        }
                                        if (item.type === 'key_sounds') {
                                          return item.value?.keySoundKey;
                                        }
                                      }
                                    "
                                    :label="$t('KeyToneAlbum.craftKeySounds.selectSounds')"
                                    multiple
                                    use-chips
                                    :class="['zl-ll']"
                                    dense
                                    :max-values="
                                      selectedKeySound.keySoundValue.up.mode.mode === 'single' ? 1 : Infinity
                                    "
                                    counter
                                    :error-message="$t('KeyToneAlbum.craftKeySounds.error.singleMode')"
                                    :error="
                                      selectedKeySound.keySoundValue.up.mode.mode === 'single' &&
                                      selectedKeySound.keySoundValue.up.value.length > 1
                                    "
                                    ref="edit_upSoundSelectDom"
                                    @update:model-value="edit_upSoundSelectDom?.hidePopup()"
                                  />
                                  <div class="h-10">
                                    <q-option-group
                                      dense
                                      v-model="edit_upTypeGroup"
                                      :options="options"
                                      type="checkbox"
                                      class="absolute left-8"
                                    >
                                      <template #label-0="props">
                                        <q-item-label>
                                          {{ $t(props.label) }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.audioFile') }}</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-1="props">
                                        <q-item-label>
                                          {{ $t(props.label) }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.soundList') }}</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-2="props">
                                        <q-item-label>
                                          {{ $t(props.label) }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.keySounds') }}</div>
                                              <div>⬇</div>
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.inheritKeySound') }}</div>
                                              <div>⬇</div>
                                              <div>{{ $t('KeyToneAlbum.craftKeySounds.tooltip.inheritRule') }}</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                    </q-option-group>
                                  </div>
                                </q-card-section>
                                <q-card-actions align="right">
                                  <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                                </q-card-actions>
                              </q-card>
                            </q-dialog>
                          </div>
                        </q-card-section>
                        <q-card-section :class="['flex justify-center gap-4 -mt-3']">
                          <q-btn
                            dense
                            class="pr-2.3"
                            color="primary"
                            icon="save"
                            :label="$t('KeyToneAlbum.confirmEdit')"
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
                            class="pr-2.3"
                            color="negative"
                            icon="delete"
                            :label="$t('KeyToneAlbum.delete')"
                            @click="
                              deleteKeySound({
                                keySoundKey: selectedKeySound.keySoundKey,
                                onSuccess: () => {
                                  selectedKeySound = undefined;
                                  q.notify({
                                    type: 'positive',
                                    position: 'top',
                                    message: $t('KeyToneAlbum.notify.deleteSuccess'),
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
                      <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
            </div>
            <q-stepper-navigation>
              <q-btn @click="step = 4" color="primary" :label="$t('KeyToneAlbum.continue')" />
              <q-btn flat @click="step = 2" color="primary" :label="$t('KeyToneAlbum.back')" class="q-ml-sm" />
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
          <q-step
            :name="4"
            :title="$t('KeyToneAlbum.linkageEffects.title')"
            icon="settings"
            :done="
              !(
                isEnableEmbeddedTestSound.down === true &&
                isEnableEmbeddedTestSound.up === true &&
                !keyDownUnifiedSoundEffectSelect &&
                !keyUpUnifiedSoundEffectSelect &&
                keysWithSoundEffect.size === 0
              )
            "
            :disable="
              step === 99 &&
              isEnableEmbeddedTestSound.down === true &&
              isEnableEmbeddedTestSound.up === true &&
              !keyDownUnifiedSoundEffectSelect &&
              !keyUpUnifiedSoundEffectSelect &&
              keysWithSoundEffect.size === 0
            "
            :header-nav="false"
            @click="
              (event: MouseEvent) => {
                const header = (event.target as HTMLElement).closest('.q-stepper__tab');
                if (header) {
                  step = step === 4 ? 99 : 4;
                }
              }
            "
          >
            <div :class="['mb-3', step_introduce_fontSize]">
              {{ $t('KeyToneAlbum.linkageEffects.description') }}
              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']">
                  <span>{{ $t('KeyToneAlbum.linkageEffects.tooltips.description') }}</span>
                </q-tooltip>
              </q-icon>
            </div>
            <div :class="['flex items-center m-t-2 w-[130%]']">
              <span class="text-gray-500 mr-0.7">•</span>
              <span class="text-nowrap">
                {{ $t('KeyToneAlbum.linkageEffects.enableTestSound') }}:
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
                    <span>{{ $t('KeyToneAlbum.linkageEffects.tooltips.testSound') }}</span>
                  </q-tooltip>
                </q-icon>
              </span>
            </div>
            <div
              :class="[
                'flex items-center ml-3',
                setting_store.languageDefault === 'pt' || setting_store.languageDefault === 'pt-BR'
                  ? 'flex-nowrap text-nowrap'
                  : '',
              ]"
            >
              <span class="text-gray-500 mr-1.5">•</span>
              <q-toggle
                v-model="isEnableEmbeddedTestSound.down"
                color="primary"
                :label="$t('KeyToneAlbum.linkageEffects.downTestSound')"
                dense
              />
            </div>
            <div
              :class="[
                'flex items-center ml-3',
                setting_store.languageDefault === 'fr' ? 'flex-nowrap text-nowrap' : '',
              ]"
            >
              <span class="text-gray-500 mr-1.5">•</span>
              <q-toggle
                v-model="isEnableEmbeddedTestSound.up"
                color="primary"
                :label="$t('KeyToneAlbum.linkageEffects.upTestSound')"
                dense
              />
            </div>
            <q-stepper-navigation>
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  :label="$t('KeyToneAlbum.linkageEffects.globalSettings')"
                  @click="
                    () => {
                      showEveryKeyEffectDialog = true;
                    }
                  "
                >
                </q-btn>
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
                    <span>{{ $t('KeyToneAlbum.linkageEffects.tooltips.globalPriority') }}</span>
                  </q-tooltip>
                </q-icon>
                <q-dialog
                  :style="{ '--i18n_fontSize': i18n_fontSize }"
                  v-model="showEveryKeyEffectDialog"
                  backdrop-filter="invert(70%)"
                  @mouseup="preventDefaultMouseWhenRecording"
                >
                  <q-card>
                    <q-card-section
                      class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                    >
                      {{ $t('KeyToneAlbum.linkageEffects.globalSettings') }}
                    </q-card-section>

                    <q-card-section class="q-pt-none">
                      <div class="text-subtitle1 q-mb-md">
                        {{ $t('KeyToneAlbum.linkageEffects.global.description') }}
                      </div>
                      <!-- 这里添加声效选择等具体设置内容 -->
                    </q-card-section>
                    <q-card-section>
                      <!-- 选择全键按下声效的选项, 仅支持单选 -->
                      <div class="flex flex-row flex-nowrap items-center m-b-3 m-l-5">
                        <div class="flex flex-col space-y-4 w-7/8">
                          <q-select
                            outlined
                            stack-label
                            :virtual-scroll-slice-size="999999"
                            popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                            v-model="keyDownUnifiedSoundEffectSelect"
                            :options="keyUnifiedSoundEffectOptions"
                            :option-label="album_options_select_label"
                            :option-value="(item: any) => {
                                  // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                  if (item.type === 'audio_files'){
                                    return item.value?.sha256 + item.value?.name_id
                                  }
                                  if(item.type === 'sounds'){
                                    return item.value?.soundKey
                                  }
                                  if(item.type === 'key_sounds'){
                                    return item.value?.keySoundKey
                                  }
                                }"
                            :label="$t('KeyToneAlbum.linkageEffects.global.setKeyDownSound')"
                            use-chips
                            :class="['zl-ll']"
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
                          >
                            <template v-slot:option="scope">
                              <q-item v-bind="scope.itemProps">
                                <q-item-section>
                                  <q-item-label>{{ album_options_select_label(scope.opt) }}</q-item-label>
                                </q-item-section>
                                <q-item-section side>
                                  <DependencyWarning
                                    v-if="scope.opt.type === 'audio_files'"
                                    :issues="dependencyIssues"
                                    item-type="audio_files"
                                    :item-id="scope.opt.value?.sha256 + scope.opt.value?.name_id"
                                    :show-details="false"
                                  />
                                  <DependencyWarning
                                    v-else-if="scope.opt.type === 'sounds'"
                                    :issues="dependencyIssues"
                                    item-type="sounds"
                                    :item-id="scope.opt.value?.soundKey"
                                    :show-details="false"
                                  />
                                  <DependencyWarning
                                    v-else-if="scope.opt.type === 'key_sounds'"
                                    :issues="dependencyIssues"
                                    item-type="key_sounds"
                                    :item-id="scope.opt.value?.keySoundKey"
                                    :show-details="false"
                                  />
                                </q-item-section>
                              </q-item>
                            </template>
                          </q-select>
                          <!-- 选择全键抬起声效的选项, 仅支持单选 -->
                          <q-select
                            outlined
                            stack-label
                            :virtual-scroll-slice-size="999999"
                            popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                            v-model="keyUpUnifiedSoundEffectSelect"
                            :options="keyUnifiedSoundEffectOptions"
                            :option-label="album_options_select_label"
                            :option-value="(item: any) => {
                                  // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                  if (item.type === 'audio_files'){
                                    return item.value?.sha256 + item.value?.name_id
                                  }
                                  if(item.type === 'sounds'){
                                    return item.value?.soundKey
                                  }
                                  if(item.type === 'key_sounds'){
                                    return item.value?.keySoundKey
                                  }
                                }"
                            :label="$t('KeyToneAlbum.linkageEffects.global.setKeyUpSound')"
                            use-chips
                            :class="['zl-ll']"
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
                          >
                            <template v-slot:option="scope">
                              <q-item v-bind="scope.itemProps">
                                <q-item-section>
                                  <q-item-label>{{ album_options_select_label(scope.opt) }}</q-item-label>
                                </q-item-section>
                                <q-item-section side>
                                  <DependencyWarning
                                    v-if="scope.opt.type === 'audio_files'"
                                    :issues="dependencyIssues"
                                    item-type="audio_files"
                                    :item-id="scope.opt.value?.sha256 + scope.opt.value?.name_id"
                                    :show-details="false"
                                  />
                                  <DependencyWarning
                                    v-else-if="scope.opt.type === 'sounds'"
                                    :issues="dependencyIssues"
                                    item-type="sounds"
                                    :item-id="scope.opt.value?.soundKey"
                                    :show-details="false"
                                  />
                                  <DependencyWarning
                                    v-else-if="scope.opt.type === 'key_sounds'"
                                    :issues="dependencyIssues"
                                    item-type="key_sounds"
                                    :item-id="scope.opt.value?.keySoundKey"
                                    :show-details="false"
                                  />
                                </q-item-section>
                              </q-item>
                            </template>
                          </q-select>
                        </div>
                        <div class="flex justify-end -m-l-2">
                          <q-icon
                            @click="isAnchoringUltimatePerfectionKeySound = !isAnchoringUltimatePerfectionKeySound"
                            size="2.75rem"
                            v-if="isShowUltimatePerfectionKeySoundAnchoring"
                          >
                            <template v-if="isAnchoringUltimatePerfectionKeySound">
                              <!-- 锚定 -->
                              <q-icon name="svguse:icons.svg#锚定"></q-icon>
                            </template>
                            <template v-else>
                              <!-- 锚定解除 -->
                              <q-icon name="svguse:icons.svg#锚定解除"></q-icon>
                            </template>
                            <q-tooltip
                              :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center']"
                            >
                              <span class="text-sm">{{
                                $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.title')
                              }}</span>
                              <span class="text-sm" v-if="isAnchoringUltimatePerfectionKeySound"
                                >{{ $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.anchored') }}<br
                              /></span>
                              <span class="text-sm" v-else
                                >{{ $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.unanchored') }}<br
                              /></span>
                              <span>{{ $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.tooltip') }}</span>
                            </q-tooltip>
                          </q-icon>
                        </div>
                      </div>
                      <div class="h-16 m-l-5.8">
                        <q-option-group dense v-model="unifiedTypeGroup" :options="options" type="checkbox">
                          <template #label-0="props">
                            <q-item-label>
                              {{ $t(props.label) }}
                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                <q-tooltip
                                  :class="[
                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                  ]"
                                >
                                  <div>{{ $t('KeyToneAlbum.linkageEffects.tooltips.audioFile') }}</div>
                                </q-tooltip>
                              </q-icon>
                            </q-item-label>
                          </template>
                          <template v-slot:label-1="props">
                            <q-item-label :class="[setting_store.languageDefault === 'es' ? 'text-nowrap' : '']">
                              {{ $t(props.label) }}
                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                <q-tooltip
                                  :class="[
                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                  ]"
                                >
                                  <div>{{ $t('KeyToneAlbum.linkageEffects.tooltips.soundList') }}</div>
                                </q-tooltip>
                              </q-icon>
                            </q-item-label>
                          </template>
                          <template v-slot:label-2="props">
                            <q-item-label>
                              {{ $t(props.label) }}
                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                <q-tooltip
                                  :class="[
                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                  ]"
                                >
                                  <div>{{ $t('KeyToneAlbum.linkageEffects.tooltips.keySounds') }}</div>
                                </q-tooltip>
                              </q-icon>
                            </q-item-label>
                          </template>
                        </q-option-group>
                      </div>
                    </q-card-section>
                    <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
                      <q-btn
                        dense
                        :label="$t('KeyToneAlbum.linkageEffects.confirm')"
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
                                message: $t('KeyToneAlbum.notify.configSuccess'),
                                timeout: 2000,
                              });
                            }
                          )
                        "
                      />
                      <q-btn flat dense color="primary" :label="$t('KeyToneAlbum.close')" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
              <div :class="['p-2 text-zinc-600']">{{ $t('KeyToneAlbum.or') }}</div>
              <div>
                <q-btn
                  :class="['bg-zinc-300']"
                  :label="$t('KeyToneAlbum.linkageEffects.singleKeySettings')"
                  @click="
                    () => {
                      showSingleKeyEffectDialog = true;
                    }
                  "
                >
                </q-btn>
                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                  <q-tooltip :class="['text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words ']">
                    <span>{{ $t('KeyToneAlbum.linkageEffects.tooltips.singleKeyPriority') }}</span>
                  </q-tooltip>
                </q-icon>
                <q-dialog
                  :style="{ '--i18n_fontSize': i18n_fontSize }"
                  v-model="showSingleKeyEffectDialog"
                  backdrop-filter="invert(70%)"
                  @mouseup="preventDefaultMouseWhenRecording"
                >
                  <q-card>
                    <q-card-section
                      class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                    >
                      {{ $t('KeyToneAlbum.linkageEffects.singleKeySettings') }}
                    </q-card-section>

                    <q-card-section class="q-pt-none pb-0">
                      <div class="text-subtitle1 q-mb-md leading-tight m-t-1.5">
                        {{ $t('KeyToneAlbum.linkageEffects.single.description') }}
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
                        {{ $t('KeyToneAlbum.linkageEffects.single.addSingleKeyEffect') }}
                        <q-dialog
                          :style="{ '--i18n_fontSize': i18n_fontSize }"
                          :no-esc-dismiss="isRecordingSingleKeys && isGetsFocused"
                          v-model="isShowAddOrSettingSingleKeyEffectDialog"
                          backdrop-filter="invert(70%)"
                          @mouseup="preventDefaultMouseWhenRecording"
                        >
                          <q-card>
                            <q-card-section
                              class="row items-center q-pb-none text-h6 sticky top-0 z-10 bg-white/30 backdrop-blur-sm"
                            >
                              {{ $t('KeyToneAlbum.linkageEffects.single.addSingleKeyEffect') }}
                            </q-card-section>

                            <q-card-section class="q-pt-none">
                              <div class="text-subtitle1 q-mb-md leading-tight m-t-1.5">
                                {{ $t('KeyToneAlbum.linkageEffects.single.dialog.selectKeyAndEffect') }}
                                <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                  <q-tooltip
                                    :class="[
                                      'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                    ]"
                                  >
                                    <div>{{ $t('KeyToneAlbum.linkageEffects.tooltips.fastDelete') }}</div>
                                    <div></div>
                                    <div></div>
                                  </q-tooltip>
                                </q-icon>
                              </div>
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
                                    :label="$t('KeyToneAlbum.linkageEffects.single.dialog.selectSingleKey')"
                                    ref="singleKeysSelectRef"
                                    popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
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
                                    :placeholder="
                                      isRecordingSingleKeys
                                        ? $t(
                                            'KeyToneAlbum.linkageEffects.single.dialog.selectSingleKey-placeholder_record'
                                          )
                                        : $t(
                                            'KeyToneAlbum.linkageEffects.single.dialog.selectSingleKey-placeholder_search'
                                          )
                                    "
                                    use-input
                                    use-chips
                                    :class="['zl-ll']"
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

                                            return false;
                                          });
                                        });
                                      }
                                    "
                                    @remove="
                                      // TIPS: @remove事件是quasar的select组件提供的。用于多选时被选项减少时触发。(另一个@clear事件, 仅作用于单选时的总清楚按钮, 无法触发多选时各个芯片上的清楚按钮)
                                      () => {
                                        clear_flag = true;
                                        singleKeysSelectRef?.focus(); // remove事件会造成所选项被清空时, 选择组件失去焦点的问题。因此通过手动获取焦点来解决此问题。
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
                                          {{
                                            isRecordingSingleKeys
                                              ? $t('KeyToneAlbum.linkageEffects.tooltips.stopRecording')
                                              : $t('KeyToneAlbum.linkageEffects.tooltips.startRecording')
                                          }}
                                        </q-tooltip>
                                      </q-btn>
                                    </template>
                                  </q-select>
                                  <!-- 按键选择 end -->

                                  <!-- 声效选择 start -->
                                  <div class="flex flex-row items-center justify-center gap-x-9 gap-y-2 w-[95%]">
                                    <q-checkbox
                                      dense
                                      :class="[
                                        setting_store.languageDefault === 'ru' || setting_store.languageDefault === 'pl'
                                          ? 'w-[200px] text-nowrap'
                                          : setting_store.languageDefault === 'fr' ||
                                            setting_store.languageDefault === 'pt' ||
                                            setting_store.languageDefault === 'pt-BR' ||
                                            setting_store.languageDefault === 'vi'
                                          ? 'w-[180px] text-nowrap'
                                          : setting_store.languageDefault === 'it' ||
                                            setting_store.languageDefault === 'tr'
                                          ? 'w-[102px] text-nowrap'
                                          : '',
                                      ]"
                                      v-model="isDownSoundEffectSelectEnabled"
                                      :label="$t('KeyToneAlbum.linkageEffects.single.dialog.downSoundEffect')"
                                    />
                                    <q-checkbox
                                      dense
                                      :class="[
                                        setting_store.languageDefault === 'ru' || setting_store.languageDefault === 'pl'
                                          ? 'w-[200px] text-nowrap'
                                          : setting_store.languageDefault === 'fr' ||
                                            setting_store.languageDefault === 'pt' ||
                                            setting_store.languageDefault === 'pt-BR' ||
                                            setting_store.languageDefault === 'vi'
                                          ? 'w-[180px] text-nowrap'
                                          : setting_store.languageDefault === 'it'
                                          ? 'w-[102px] text-nowrap'
                                          : '',
                                      ]"
                                      v-model="isUpSoundEffectSelectEnabled"
                                      :label="$t('KeyToneAlbum.linkageEffects.single.dialog.upSoundEffect')"
                                    />
                                  </div>
                                  <div class="w-full">
                                    <q-card-section>
                                      <div class="flex flex-row flex-nowrap items-center m-b-3 m-l-5">
                                        <div class="flex flex-col space-y-4 w-[223px]">
                                          <!-- 选择单键按下声效的选项, 仅支持单选 -->
                                          <q-select
                                            v-show="isDownSoundEffectSelectEnabled"
                                            outlined
                                            stack-label
                                            :virtual-scroll-slice-size="999999"
                                            popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                            v-model="keyDownSingleKeySoundEffectSelect"
                                            :options="keySingleKeySoundEffectOptions"
                                            :option-label="album_options_select_label"
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
                                            :label="$t('KeyToneAlbum.linkageEffects.single.dialog.setDownSoundEffect')"
                                            use-chips
                                            :class="['zl-ll']"
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
                                          >
                                            <template v-slot:option="scope">
                                              <q-item v-bind="scope.itemProps">
                                                <q-item-section>
                                                  <q-item-label>{{ album_options_select_label(scope.opt) }}</q-item-label>
                                                </q-item-section>
                                                <q-item-section side>
                                                  <DependencyWarning
                                                    v-if="scope.opt.type === 'audio_files'"
                                                    :issues="dependencyIssues"
                                                    item-type="audio_files"
                                                    :item-id="scope.opt.value?.sha256 + scope.opt.value?.name_id"
                                                    :show-details="false"
                                                  />
                                                  <DependencyWarning
                                                    v-else-if="scope.opt.type === 'sounds'"
                                                    :issues="dependencyIssues"
                                                    item-type="sounds"
                                                    :item-id="scope.opt.value?.soundKey"
                                                    :show-details="false"
                                                  />
                                                  <DependencyWarning
                                                    v-else-if="scope.opt.type === 'key_sounds'"
                                                    :issues="dependencyIssues"
                                                    item-type="key_sounds"
                                                    :item-id="scope.opt.value?.keySoundKey"
                                                    :show-details="false"
                                                  />
                                                </q-item-section>
                                              </q-item>
                                            </template>
                                          </q-select>
                                          <!-- 选择单键抬起声效的选项, 仅支持单选 -->
                                          <q-select
                                            v-show="isUpSoundEffectSelectEnabled"
                                            outlined
                                            stack-label
                                            :virtual-scroll-slice-size="999999"
                                            popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                            v-model="keyUpSingleKeySoundEffectSelect"
                                            :options="keySingleKeySoundEffectOptions"
                                            :option-label="album_options_select_label"
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
                                            :label="$t('KeyToneAlbum.linkageEffects.single.dialog.setUpSoundEffect')"
                                            use-chips
                                            :class="['zl-ll']"
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
                                          >
                                            <template v-slot:option="scope">
                                              <q-item v-bind="scope.itemProps">
                                                <q-item-section>
                                                  <q-item-label>{{ album_options_select_label(scope.opt) }}</q-item-label>
                                                </q-item-section>
                                                <q-item-section side>
                                                  <DependencyWarning
                                                    v-if="scope.opt.type === 'audio_files'"
                                                    :issues="dependencyIssues"
                                                    item-type="audio_files"
                                                    :item-id="scope.opt.value?.sha256 + scope.opt.value?.name_id"
                                                    :show-details="false"
                                                  />
                                                  <DependencyWarning
                                                    v-else-if="scope.opt.type === 'sounds'"
                                                    :issues="dependencyIssues"
                                                    item-type="sounds"
                                                    :item-id="scope.opt.value?.soundKey"
                                                    :show-details="false"
                                                  />
                                                  <DependencyWarning
                                                    v-else-if="scope.opt.type === 'key_sounds'"
                                                    :issues="dependencyIssues"
                                                    item-type="key_sounds"
                                                    :item-id="scope.opt.value?.keySoundKey"
                                                    :show-details="false"
                                                  />
                                                </q-item-section>
                                              </q-item>
                                            </template>
                                          </q-select>
                                        </div>
                                        <div
                                          v-show="isDownSoundEffectSelectEnabled && isUpSoundEffectSelectEnabled"
                                          :class="['absolute -right-2']"
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
                                              <q-icon name="svguse:icons.svg#锚定"></q-icon>
                                            </template>
                                            <template v-else>
                                              <!-- 锚定解除 -->
                                              <q-icon name="svguse:icons.svg#锚定解除"></q-icon>
                                            </template>
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <span class="text-sm">{{
                                                $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.title')
                                              }}</span>
                                              <span
                                                class="text-sm"
                                                v-if="isAnchoringUltimatePerfectionKeySound_singleKey"
                                              >
                                                {{ $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.anchored')
                                                }}<br
                                              /></span>
                                              <span class="text-sm" v-else
                                                >{{
                                                  $t(
                                                    'KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.unanchored'
                                                  )
                                                }}<br
                                              /></span>
                                              <span>{{
                                                $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.tooltip')
                                              }}</span>
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
                                              {{ $t(props.label) }}
                                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                                <q-tooltip
                                                  :class="[
                                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                                  ]"
                                                >
                                                  <div>{{ $t('KeyToneAlbum.linkageEffects.tooltips.audioFile') }}</div>
                                                </q-tooltip>
                                              </q-icon>
                                            </q-item-label>
                                          </template>
                                          <template v-slot:label-1="props">
                                            <q-item-label
                                              :class="[
                                                setting_store.languageDefault === 'fr' ||
                                                setting_store.languageDefault === 'it' ||
                                                setting_store.languageDefault === 'pt' ||
                                                setting_store.languageDefault === 'pt-BR'
                                                  ? 'text-nowrap'
                                                  : '',
                                              ]"
                                            >
                                              {{ $t(props.label) }}
                                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                                <q-tooltip
                                                  :class="[
                                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                                  ]"
                                                >
                                                  <div>{{ $t('KeyToneAlbum.linkageEffects.tooltips.soundList') }}</div>
                                                </q-tooltip>
                                              </q-icon>
                                            </q-item-label>
                                          </template>
                                          <template v-slot:label-2="props">
                                            <q-item-label>
                                              {{ $t(props.label) }}
                                              <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                                <q-tooltip
                                                  :class="[
                                                    'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                                  ]"
                                                >
                                                  <div>
                                                    {{ $t('KeyToneAlbum.linkageEffects.tooltips.keySounds') }}
                                                  </div>
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
                                :label="$t('KeyToneAlbum.linkageEffects.confirm')"
                                color="primary"
                                @click="
                                  if (selectedSingleKeys.length !== 0) {
                                    // 保存之前先记录下keysWithSoundEffect的值, 用于保存成功回调中的判断。 保证保存成功的回调中, 判断逻辑的正确性。
                                    const keysWithSoundEffect_old = [...keysWithSoundEffect]; // 注意,这里我们直接将解构的map给到了数组, 因此后续需要使用数组的some函数判断,而不是map的has。

                                    saveSingleKeySoundEffectConfig(
                                      {
                                        singleKeys: selectedSingleKeys,
                                        down: keyDownSingleKeySoundEffectSelect,
                                        up: keyUpSingleKeySoundEffectSelect,
                                      },
                                      () => {
                                        if (!keyDownSingleKeySoundEffectSelect && !keyUpSingleKeySoundEffectSelect) {
                                          const isDeletes: number[] = [];
                                          const isNotSoundEffectSelect: number[] = [];
                                          selectedSingleKeys.forEach((key) => {
                                            // if (keysWithSoundEffect.has(String(key))) { // TIPS: 由于onSuccess发生在保存成功后, 此时的keysWithSoundEffect已经是更改过的值了(判断的结果是不可用的), 我们应该就更改前的值来做判断才是正确的逻辑。
                                            if (keysWithSoundEffect_old.some((item) => item[0] === String(key))) {
                                              isDeletes.push(key);
                                            } else {
                                              isNotSoundEffectSelect.push(key);
                                            }
                                          });
                                          if (isNotSoundEffectSelect.length !== 0) {
                                            if (isDeletes.length === 0) {
                                              q.notify({
                                                type: 'warning',
                                                position: 'top',
                                                message: $t('KeyToneAlbum.notify.selectSoundEffect'),
                                                timeout: 2000,
                                              });
                                              return;
                                            } else {
                                              // 此时代表isDeletes.length !== 0, 即本次操作同时存在删除单键声效配置的操作, 和一些未知目的的错误录制的按键
                                              let deleteString: string = '';
                                              let notSoundEffectSelectString: string = '';
                                              isDeletes.forEach((key) => {
                                                deleteString +=
                                                  '-[' +
                                                  (keyEvent_store.dikCodeToName.get(key) || 'Dik-{' + key + '}') +
                                                  ']';
                                              });
                                              isNotSoundEffectSelect.forEach((key) => {
                                                notSoundEffectSelectString +=
                                                  '-[' +
                                                  (keyEvent_store.dikCodeToName.get(key) || 'Dik-{' + key + '}') +
                                                  ']';
                                              });
                                              q.notify({
                                                type: 'warning',
                                                position: 'top',
                                                message: $t('KeyToneAlbum.notify.selectSoundEffectForKeys', {
                                                  keys: notSoundEffectSelectString,
                                                }),
                                                timeout: 8000,
                                              });
                                              q.notify({
                                                type: 'positive',
                                                position: 'top',
                                                message: $t('KeyToneAlbum.notify.deleteSuccessForKeys', {
                                                  keys: deleteString,
                                                }),
                                                timeout: 3000,
                                              });
                                              selectedSingleKeys = isNotSoundEffectSelect;
                                              return;
                                            }
                                          } else {
                                            // 此时代表isNotSoundEffectSelect.length === 0, 即本次操作为专门的删除按键配置的操作
                                            q.notify({
                                              type: 'positive',
                                              position: 'top',
                                              message: $t('KeyToneAlbum.notify.deleteSuccess'),
                                              timeout: 2000,
                                            });
                                            selectedSingleKeys = [];
                                            keyDownSingleKeySoundEffectSelect = null;
                                            keyUpSingleKeySoundEffectSelect = null;
                                            isShowAddOrSettingSingleKeyEffectDialog = false;
                                            return;
                                          }
                                        }

                                        q.notify({
                                          type: 'positive',
                                          position: 'top',
                                          message: $t('KeyToneAlbum.notify.configSuccess'),
                                          timeout: 2000,
                                        });
                                        selectedSingleKeys = [];
                                        keyDownSingleKeySoundEffectSelect = null;
                                        keyUpSingleKeySoundEffectSelect = null;
                                        isShowAddOrSettingSingleKeyEffectDialog = false;
                                      }
                                    );
                                  } else {
                                    q.notify({
                                      type: 'warning',
                                      position: 'top',
                                      message: $t('KeyToneAlbum.notify.selectKey'),
                                      timeout: 2000,
                                    });
                                  }
                                "
                              />
                              <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
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
                        {{ $t('KeyToneAlbum.linkageEffects.single.dialog.noSingleKeyEffects') }}
                      </div>
                      <div v-else class="text-[1.06rem] pb-2 font-600 text-gray-700 flex flex-row items-center">
                        {{ $t('KeyToneAlbum.linkageEffects.single.dialog.singleKeyEffects') }}
                        <div class="text-[0.88rem] ml-1">
                          ({{ $t('KeyToneAlbum.linkageEffects.single.dialog.clickToView') }})
                        </div>
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
                              currentEditingKey_old = currentEditingKey;
                              currentEditingKey = Number(item[0]);
                              // TODO: 进一步, 需要在此处读取对应单键的声效设置, 并用读取的数据来初始化对话框, 以供用户的后续编辑。
                              if (currentEditingKey !== currentEditingKey_old) {
                                // 如果点击的是同一个按键, 则没必要重新初始值。(即默认持久化刚才的操作记录)
                                singleKeyTypeGroup_edit = ['sounds']; // 为了防止初始有'key_sounds'时触发的锚定, 会影响原数据的初始化(如因锚定错误的同步双方等连锁反应)
                                const down = item[1].down;
                                const up = item[1].up;
                                keyDownSingleKeySoundEffectSelect_edit = convertValue(down ? down : '');
                                keyUpSingleKeySoundEffectSelect_edit = convertValue(up ? up : '');
                                keyDownSingleKeySoundEffectSelect_edit_old = keyDownSingleKeySoundEffectSelect_edit;
                                keyUpSingleKeySoundEffectSelect_edit_old = keyUpSingleKeySoundEffectSelect_edit;
                              }
                            }
                          "
                        >
                          {{ keyEvent_store.dikCodeToName.get(Number(item[0])) || 'Dik-{' + item[0] + '}' }}
                        </q-chip>
                        <q-dialog
                          :style="{ '--i18n_fontSize': i18n_fontSize }"
                          v-model="isShowSingleKeySoundEffectEditDialog"
                          @mouseup="preventDefaultMouseWhenRecording"
                        >
                          <q-card style="min-width: 350px">
                            <q-card-section>
                              <div class="text-base flex flex-row items-center">
                                {{ $t('KeyToneAlbum.linkageEffects.single.dialog.editSingleKey') }} -
                                <div class="text-sm font-bold">
                                  [
                                  {{ currentEditingKeyOfName }}
                                  ]
                                </div>
                                - {{ $t('KeyToneAlbum.linkageEffects.single.dialog.soundEffect') }}
                              </div>
                            </q-card-section>

                            <q-card-section class="q-pt-none pb-1">
                              <!-- 这里之后添加编辑内容 -->
                              <!-- 声效编辑  start -->
                              <div class="w-full">
                                <q-card-section>
                                  <div class="flex flex-row flex-nowrap items-center mb-3">
                                    <div class="flex flex-col space-y-4 w-full">
                                      <!-- 选择单键按下声效的选项, 仅支持单选 [声效编辑]-->
                                      <q-select
                                        outlined
                                        stack-label
                                        :virtual-scroll-slice-size="999999"
                                        popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                        v-model="keyDownSingleKeySoundEffectSelect_edit"
                                        :options="keySingleKeySoundEffectOptions_edit"
                                        :option-label="album_options_select_label"
                                        :option-value="(item: any) => {
                                              // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                              if (item.type === 'audio_files'){
                                                return item.value?.sha256 + item.value?.name_id
                                              }
                                              if(item.type === 'sounds'){
                                                return item.value?.soundKey
                                              }
                                              if(item.type === 'key_sounds'){
                                                return item.value?.keySoundKey
                                              }
                                            }"
                                        :label="`${$t(
                                          'KeyToneAlbum.linkageEffects.single.dialog.editSingleKey'
                                        )} -[ ${currentEditingKeyOfName} ]- ${$t(
                                          'KeyToneAlbum.linkageEffects.single.dialog.soundEffect-down'
                                        )} `"
                                        use-chips
                                        :class="['zl-ll']"
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
                                      >
                                        <template v-slot:option="scope">
                                          <q-item v-bind="scope.itemProps">
                                            <q-item-section>
                                              <q-item-label>{{ album_options_select_label(scope.opt) }}</q-item-label>
                                            </q-item-section>
                                            <q-item-section side>
                                              <DependencyWarning
                                                v-if="scope.opt.type === 'audio_files'"
                                                :issues="dependencyIssues"
                                                item-type="audio_files"
                                                :item-id="scope.opt.value?.sha256 + scope.opt.value?.name_id"
                                                :show-details="false"
                                              />
                                              <DependencyWarning
                                                v-else-if="scope.opt.type === 'sounds'"
                                                :issues="dependencyIssues"
                                                item-type="sounds"
                                                :item-id="scope.opt.value?.soundKey"
                                                :show-details="false"
                                              />
                                              <DependencyWarning
                                                v-else-if="scope.opt.type === 'key_sounds'"
                                                :issues="dependencyIssues"
                                                item-type="key_sounds"
                                                :item-id="scope.opt.value?.keySoundKey"
                                                :show-details="false"
                                              />
                                            </q-item-section>
                                          </q-item>
                                        </template>
                                      </q-select>
                                      <!-- 选择单键抬起声效的选项, 仅支持单选 [声效编辑]-->
                                      <q-select
                                        outlined
                                        stack-label
                                        :virtual-scroll-slice-size="999999"
                                        popup-content-class="w-[50%] [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
                                        v-model="keyUpSingleKeySoundEffectSelect_edit"
                                        :options="keySingleKeySoundEffectOptions_edit"
                                        :option-label="album_options_select_label"
                                        :option-value="(item: any) => {
                                              // 直接设置uuid, 使组件可轻松精确的区分每个选项。
                                              if (item.type === 'audio_files'){
                                                return item.value?.sha256 + item.value?.name_id
                                              }
                                              if(item.type === 'sounds'){
                                                return item.value?.soundKey
                                              }
                                              if(item.type === 'key_sounds'){
                                                return item.value?.keySoundKey
                                              }
                                            }"
                                        :label="`${$t(
                                          'KeyToneAlbum.linkageEffects.single.dialog.editSingleKey'
                                        )} -[ ${currentEditingKeyOfName} ]- ${$t(
                                          'KeyToneAlbum.linkageEffects.single.dialog.soundEffect-up'
                                        )} `"
                                        use-chips
                                        :class="['zl-ll']"
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
                                      >
                                        <template v-slot:option="scope">
                                          <q-item v-bind="scope.itemProps">
                                            <q-item-section>
                                              <q-item-label>{{ album_options_select_label(scope.opt) }}</q-item-label>
                                            </q-item-section>
                                            <q-item-section side>
                                              <DependencyWarning
                                                v-if="scope.opt.type === 'audio_files'"
                                                :issues="dependencyIssues"
                                                item-type="audio_files"
                                                :item-id="scope.opt.value?.sha256 + scope.opt.value?.name_id"
                                                :show-details="false"
                                              />
                                              <DependencyWarning
                                                v-else-if="scope.opt.type === 'sounds'"
                                                :issues="dependencyIssues"
                                                item-type="sounds"
                                                :item-id="scope.opt.value?.soundKey"
                                                :show-details="false"
                                              />
                                              <DependencyWarning
                                                v-else-if="scope.opt.type === 'key_sounds'"
                                                :issues="dependencyIssues"
                                                item-type="key_sounds"
                                                :item-id="scope.opt.value?.keySoundKey"
                                                :show-details="false"
                                              />
                                            </q-item-section>
                                          </q-item>
                                        </template>
                                      </q-select>
                                    </div>
                                    <div class="flex justify-end -m-l-2">
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
                                          <q-icon name="svguse:icons.svg#锚定"></q-icon>
                                        </template>
                                        <template v-else>
                                          <!-- 锚定解除 [声效编辑]-->
                                          <q-icon name="svguse:icons.svg#锚定解除"></q-icon>
                                        </template>
                                        <q-tooltip
                                          :class="[
                                            'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                          ]"
                                        >
                                          <span class="text-sm">{{
                                            $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.title')
                                          }}</span>
                                          <span
                                            class="text-sm"
                                            v-if="isAnchoringUltimatePerfectionKeySound_singleKey_edit"
                                          >
                                            {{ $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.anchored')
                                            }}<br
                                          /></span>
                                          <span class="text-sm" v-else
                                            >{{ $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.unanchored')
                                            }}<br
                                          /></span>
                                          <span>{{
                                            $t('KeyToneAlbum.linkageEffects.tooltips.ultimateKeySound.tooltip')
                                          }}</span>
                                        </q-tooltip>
                                      </q-icon>
                                    </div>
                                  </div>
                                  <div class="h-16 m-l-9">
                                    <q-option-group
                                      dense
                                      v-model="singleKeyTypeGroup_edit"
                                      :options="options"
                                      type="checkbox"
                                    >
                                      <template #label-0="props">
                                        <q-item-label>
                                          {{ $t(props.label) }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>{{ $t('KeyToneAlbum.linkageEffects.tooltips.audioFile') }}</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-1="props">
                                        <q-item-label>
                                          {{ $t(props.label) }}
                                          <q-icon name="info" color="primary" class="p-l-4.5 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>{{ $t('KeyToneAlbum.linkageEffects.tooltips.soundList') }}</div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                      <template v-slot:label-2="props">
                                        <q-item-label>
                                          {{ $t(props.label) }}
                                          <q-icon name="info" color="primary" class="p-l-1 m-b-0.5">
                                            <q-tooltip
                                              :class="[
                                                'text-xs bg-opacity-80 bg-gray-700 whitespace-pre-wrap break-words text-center',
                                              ]"
                                            >
                                              <div>
                                                {{ $t('KeyToneAlbum.linkageEffects.tooltips.keySounds') }}
                                              </div>
                                            </q-tooltip>
                                          </q-icon>
                                        </q-item-label>
                                      </template>
                                    </q-option-group>
                                  </div>
                                  <div class="flex justify-center gap-8 -m-l-3 m-t-5">
                                    <q-btn
                                      class="p-r-2"
                                      dense
                                      :label="$t('KeyToneAlbum.confirmEdit')"
                                      color="primary"
                                      icon="save"
                                      @click="
                                        async () => {
                                          // TIPS: 只是这样简单的对比, 无法在用户操作后, 仍选择相同声效时作出正确的判断, 因此后续可以更精确的对uuid进行比较。
                                          const uuid = (item: any) => {
                                            if (item?.type === 'audio_files') {
                                              return item?.value.sha256 + item?.value.name_id;
                                            }
                                            if (item?.type === 'sounds') {
                                              return item?.value.soundKey;
                                            }
                                            if (item?.type === 'key_sounds') {
                                              return item?.value.keySoundKey;
                                            }
                                          };
                                          if (
                                            uuid(keyDownSingleKeySoundEffectSelect_edit) !==
                                              uuid(keyDownSingleKeySoundEffectSelect_edit_old) ||
                                            uuid(keyUpSingleKeySoundEffectSelect_edit) !==
                                              uuid(keyUpSingleKeySoundEffectSelect_edit_old)
                                          ) {
                                            // TIPS: 此时毕竟有可能存在误操作(即用户可能并不是主观的想要删除此单键声效的配置, 因此最好能有个二次提示的窗口来警告操作后果会删除单键, 并询问用户是否继续操作)
                                            if (
                                              !keyDownSingleKeySoundEffectSelect_edit &&
                                              !keyUpSingleKeySoundEffectSelect_edit
                                            ) {
                                              q.dialog({
                                                title: $t('KeyToneAlbum.notify.confirmDeleteSingleKeyEffect'),
                                                message: $t('KeyToneAlbum.notify.confirmDeleteSingleKeyEffectMessage'),
                                                ok: {
                                                  label: $t('KeyToneAlbum.cancel'),
                                                  color: 'primary',
                                                  flat: true,
                                                },
                                                cancel: {
                                                  label: $t('KeyToneAlbum.confirm'),
                                                  color: 'primary',
                                                  flat: true,
                                                },
                                                persistent: true,
                                                focus: 'cancel',
                                              }).onCancel(() => {
                                                saveSingleKeySoundEffectConfig(
                                                  {
                                                    singleKeys: currentEditingKey ? [currentEditingKey] : [],
                                                    down: keyDownSingleKeySoundEffectSelect_edit,
                                                    up: keyUpSingleKeySoundEffectSelect_edit,
                                                  },
                                                  () => {
                                                    if (
                                                      !keyDownSingleKeySoundEffectSelect_edit &&
                                                      !keyUpSingleKeySoundEffectSelect_edit
                                                    ) {
                                                      q.notify({
                                                        type: 'positive',
                                                        position: 'top',
                                                        message: $t('KeyToneAlbum.notify.deleteSuccess'),
                                                        timeout: 2000,
                                                      });
                                                    } else {
                                                      q.notify({
                                                        type: 'positive',
                                                        position: 'top',
                                                        message: $t('KeyToneAlbum.notify.saveSuccess'),
                                                        timeout: 2000,
                                                      });
                                                    }
                                                    // TIPS: 由于对连续打开同一个按键对话框时, 默认持久化(即不会更新响应的SoundEffectSelect_edit_old值), 因此需要在保存成功后, 主动地更新它们, 才能完善整体逻辑。
                                                    keyDownSingleKeySoundEffectSelect_edit_old =
                                                      keyDownSingleKeySoundEffectSelect_edit;
                                                    keyUpSingleKeySoundEffectSelect_edit_old =
                                                      keyUpSingleKeySoundEffectSelect_edit;
                                                    // 关闭对话框
                                                    isShowSingleKeySoundEffectEditDialog = false;
                                                  }
                                                );
                                              });
                                            } else {
                                              saveSingleKeySoundEffectConfig(
                                                {
                                                  singleKeys: currentEditingKey ? [currentEditingKey] : [],
                                                  down: keyDownSingleKeySoundEffectSelect_edit,
                                                  up: keyUpSingleKeySoundEffectSelect_edit,
                                                },
                                                () => {
                                                  if (
                                                    !keyDownSingleKeySoundEffectSelect_edit &&
                                                    !keyUpSingleKeySoundEffectSelect_edit
                                                  ) {
                                                    q.notify({
                                                      type: 'positive',
                                                      position: 'top',
                                                      message: $t('KeyToneAlbum.notify.deleteSuccess'),
                                                      timeout: 2000,
                                                    });
                                                  } else {
                                                    q.notify({
                                                      type: 'positive',
                                                      position: 'top',
                                                      message: $t('KeyToneAlbum.notify.saveSuccess'),
                                                      timeout: 2000,
                                                    });
                                                  }
                                                  // TIPS: 由于对连续打开同一个按键对话框时, 默认持久化(即不会更新响应的SoundEffectSelect_edit_old值), 因此需要在保存成功后, 主动地更新它们, 才能完善整体逻辑。
                                                  keyDownSingleKeySoundEffectSelect_edit_old =
                                                    keyDownSingleKeySoundEffectSelect_edit;
                                                  keyUpSingleKeySoundEffectSelect_edit_old =
                                                    keyUpSingleKeySoundEffectSelect_edit;
                                                  // 关闭对话框
                                                  isShowSingleKeySoundEffectEditDialog = false;
                                                }
                                              );
                                            }
                                          } else {
                                            q.notify({
                                              type: 'warning',
                                              position: 'top',
                                              message: $t('KeyToneAlbum.notify.noChangesDetected'),
                                              timeout: 2000,
                                            });
                                          }
                                        }
                                      "
                                    >
                                      <q-tooltip>{{ $t('KeyToneAlbum.linkageEffects.tooltips.save') }}</q-tooltip>
                                    </q-btn>
                                    <q-btn
                                      class="p-r-2"
                                      dense
                                      :label="$t('KeyToneAlbum.delete')"
                                      color="negative"
                                      icon="delete"
                                      v-close-popup
                                      @click="
                                        saveSingleKeySoundEffectConfig(
                                          {
                                            singleKeys: currentEditingKey ? [currentEditingKey] : [],
                                            down: null,
                                            up: null,
                                          },
                                          () => {
                                            q.notify({
                                              type: 'positive',
                                              position: 'top',
                                              message: $t('KeyToneAlbum.notify.deleteSuccess'),
                                              timeout: 2000,
                                            });
                                          }
                                        )
                                      "
                                    >
                                      <q-tooltip>{{ $t('KeyToneAlbum.linkageEffects.tooltips.delete') }}</q-tooltip>
                                    </q-btn>
                                  </div>
                                </q-card-section>
                              </div>
                            </q-card-section>
                            <q-card-actions class="pt-0" align="right">
                              <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                            </q-card-actions>
                          </q-card>
                          <!-- 声效编辑  end -->
                        </q-dialog>
                      </div>
                    </q-card-section>
                    <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
                      <q-btn flat :label="$t('KeyToneAlbum.close')" color="primary" v-close-popup />
                    </q-card-actions>
                  </q-card>
                </q-dialog>
              </div>
            </q-stepper-navigation>
            <q-stepper-navigation>
              <q-btn @click="step = 5" color="primary" :label="$t('KeyToneAlbum.continue')" />
              <q-btn flat @click="step = 3" color="primary" :label="$t('KeyToneAlbum.back')" class="q-ml-sm" />
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
import { QDialog, QSelect, useQuasar } from 'quasar';
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
import { useMainStore } from 'src/stores/main-store';
import { useSettingStore } from 'src/stores/setting-store';
import { computed, onBeforeMount, ref, watch, useTemplateRef, reactive, nextTick, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { 
  createDependencyValidator, 
  hasItemDependencyIssues,
  type DependencyIssue,
  type AudioFile,
  type Sound,
  type KeySound 
} from 'src/utils/dependencyValidator';
import DependencyWarning from 'src/components/DependencyWarning.vue';

// console.error("重新载入")   // 用笨方法, 严重组件的重新渲染情况
const q = useQuasar();
const { t } = useI18n();
const $t = t;
const app_store = useAppStore();
const setting_store = useSettingStore();

export interface Props {
  pkgPath: string;
  isCreate: boolean;
}
const props = withDefaults(defineProps<Props>(), {});

// 防止空字符串触发不能为空的提示, 虽然初始化时只有一瞬间, 但也不希望看到
const pkgName = ref<string>($t('KeyToneAlbum.new.name.defaultValue'));

const step = ref(99);
// watch(step, () => {
//   console.log('step-------------=', step.value);
// });

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
      message: $t('KeyToneAlbum.notify.selectSourceFile'),
      timeout: 5,
    });
    return;
  }
  // 结束时间必须大于开始时间
  if (params.cut.end_time <= params.cut.start_time) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.endTimeGreaterThanStartTime'),
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
        message: params.soundKey ? $t('KeyToneAlbum.notify.modifySuccess') : $t('KeyToneAlbum.notify.addSuccess'),
        timeout: 5,
      });
      params.onSuccess?.();
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: params.soundKey ? $t('KeyToneAlbum.notify.modifyFailed') : $t('KeyToneAlbum.notify.addFailed'),
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
        message: $t('KeyToneAlbum.notify.deleteFailed'),
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
      message: $t('KeyToneAlbum.notify.selectAudioFile'),
      timeout: 5000,
    });
    return;
  }
  // 结束时间必须大于开始时间
  if (params.cut.end_time <= params.cut.start_time) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.endTimeGreaterThanStartTime'),
      timeout: 5,
    });
    return;
  }
  // 时间值不能为负数
  if (params.cut.start_time < 0 || params.cut.end_time < 0) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.timeValueCannotBeNegative'),
      timeout: 5,
    });
  }

  PlaySound(
    params.source_file_for_sound.sha256,
    params.source_file_for_sound.type,
    params.cut.start_time,
    params.cut.end_time,
    params.cut.volume,
    true // 设置 skipGlobalVolume 为 true，使预览不受全局音量影响
  ).then((result) => {
    if (!result) {
      q.notify({
        type: 'negative',
        position: 'top',
        message: $t('KeyToneAlbum.notify.playFailed'),
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
const options = reactive([
  { label: 'KeyToneAlbum.options.audioFile', value: 'audio_files', label_0: 'KeyToneAlbum.options.audioFile_0' },
  { label: 'KeyToneAlbum.options.sound', value: 'sounds', label_0: 'KeyToneAlbum.options.sound_0' },
  { label: 'KeyToneAlbum.options.keySound', value: 'key_sounds', label_0: 'KeyToneAlbum.options.keySound_0' },
]);

const album_options_select_label = (item: any): any => {
  // console.log('item_1212==', item);
  // console.log('至臻键音列表==', keySoundList.value);
  if (item.type === 'audio_files') {
    return (
      $t(options.find((option) => item.type === option.value)?.label_0 || '') +
      ' § ' +
      soundFileList.value.find(
        (soundFile: any) => soundFile.sha256 === item.value?.sha256 && soundFile.name_id === item.value?.name_id
      )?.name +
      soundFileList.value.find(
        (soundFile: any) => soundFile.sha256 === item.value?.sha256 && soundFile.name_id === item.value?.name_id
      )?.type
    );
  }
  if (item.type === 'sounds') {
    // 此处的item可以是any , 但其soundList的源类型, 必须是指定准确, 否则此处会发生意外报错, 且无法定位
    if (item.value?.soundValue?.name !== '' && item.value?.soundValue?.name !== undefined) {
      return (
        $t(options.find((option) => item.type === option.value)?.label_0 || '') +
        ' § ' +
        soundList.value.find((sound) => sound.soundKey === item.value?.soundKey)?.soundValue.name
      );
    } else {
      return (
        $t(options.find((option) => item.type === option.value)?.label_0 || '') +
        ' § ' +
        (soundFileList.value.find(
          (soundFile: any) =>
            soundFile.sha256 === item.value?.soundValue?.source_file_for_sound?.sha256 &&
            soundFile.name_id === item.value?.soundValue?.source_file_for_sound?.name_id
        )?.name +
          '     - ' +
          ' [' +
          item.value?.soundValue?.cut?.start_time +
          ' ~ ' +
          item.value?.soundValue?.cut?.end_time +
          ']')
      );
    }
  }
  if (item.type === 'key_sounds') {
    // Check if item.value is valid before accessing its properties
    return (
      $t(options.find((option) => item.type === option.value)?.label_0 || 'Error') +
      ' § ' +
      // (item.value.keySoundValue.name || '[Unnamed]')
      keySoundList.value.find((keySound) => keySound.keySoundKey === item.value?.keySoundKey)?.keySoundValue?.name
    );
  }
};

// 按键音
const playModeOptions = ['single', 'random', 'loop'];
const playModeLabels = new Map<string, string>([
  ['single', 'KeyToneAlbum.playMode.single'],
  ['random', 'KeyToneAlbum.playMode.random'],
  ['loop', 'KeyToneAlbum.playMode.loop'],
]);

// 按键音制作
const createNewKeySound = ref(false);

// -- createNewKeySound
const keySoundName = ref<string>($t('KeyToneAlbum.craftKeySounds.keySoundName-placeholder'));
const configureDownSound = ref(false);
const configureUpSound = ref(false);

// -- configureDownSound
const selectedSoundsForDown = ref<Array<any>>([]);
const playModeForDown = ref('random');
const maxSelectionForDown = computed(() => {
  return playModeForDown.value === 'single' ? 1 : Infinity;
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
const playModeForUp = ref('random');
const maxSelectionForUp = computed(() => {
  return playModeForUp.value === 'single' ? 1 : Infinity;
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
  if (!selectedKeySound.value) {
    return;
  }
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
          (soundFile) => item.value && soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id
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
          (soundFile) => item.value && soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id
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
      message: $t('KeyToneAlbum.notify.downSoundInvalid'),
      timeout: 3000,
    });
    isReturn = true;
  }
  if (params.up.mode === 'single' && params.up.value.length > 1) {
    q.notify({
      type: 'negative',
      position: 'top',
      message: $t('KeyToneAlbum.notify.upSoundInvalid'),
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
        message: params.key ? $t('KeyToneAlbum.notify.modifySuccess') : $t('KeyToneAlbum.notify.addSuccess'),
        timeout: 5,
      });
      if (onSuccess) {
        onSuccess();
      }
    } else {
      q.notify({
        type: 'negative',
        position: 'top',
        message: params.key ? $t('KeyToneAlbum.notify.modifyFailed') : $t('KeyToneAlbum.notify.addFailed'),
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
        message: $t('KeyToneAlbum.notify.deleteFailed'),
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
          message: $t('KeyToneAlbum.notify.unifiedSoundEffectConfigFailed'),
          timeout: 5000,
        });
      }
    })
    .catch((err) => {
      console.error('全键声效配置时发生错误:', err);
      q.notify({
        type: 'negative',
        position: 'top',
        message: $t('KeyToneAlbum.notify.unifiedSoundEffectConfigFailed'),
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

const keyEvent_store = useKeyEventStore();

const isRecordingSingleKeys = ref(false);

const keyOptions = computed(() => {
  // 将 Map 转换为数组形式的选项
  if (isRecordingSingleKeys.value) {
    return [];
  } else {
    // 默认以系统主映射表的keys为主
    const reArray = Array.from(keyEvent_store.dikCodeToName.keys());

    return reArray;
  }
});
const filterOptions = ref(keyOptions.value); // 用于过滤选项

const isGetsFocused = ref(false);

let first_flag = false; // 用于避免录制按键打开瞬间(即isRecordingSingleKeys由false->true的瞬间)鼠标左键被记录。
let clear_flag = false; // 用于避免录制按键打开瞬间(即isRecordingSingleKeys由false->true的瞬间)鼠标左键被记录。

const recordingSingleKeysCallback = (keycode: number, keyName: string) => {
  console.debug('keycode=', keycode, 'keyName=', keyName);
  if (!first_flag) {
    first_flag = true;
    return;
  }
  if (clear_flag) {
    clear_flag = false;
    return;
  }

  // 如果按键不在列表中，则添加
  if (!selectedSingleKeys.value.includes(keycode)) {
    selectedSingleKeys.value.push(keycode);
  } else {
    q.notify({
      type: 'info',
      position: 'top',
      message: $t('KeyToneAlbum.notify.keyAlreadySelected'),
      timeout: 1000,
    });
  }

  console.debug('当前已选择的按键:', selectedSingleKeys.value);
};

watch(isShowAddOrSettingSingleKeyEffectDialog, (newVal) => {
  if (!newVal) {
    keyEvent_store.clearKeyStateCallback_Record();
    // 当通过点击对话框外使得对话框关闭时, 不会触发失去焦点的事件(因此此时isGetsFocused的值不会被置为false, 故补充此逻辑)
    isGetsFocused.value = false;
  }
});

watch(isRecordingSingleKeys, (newVal, oldVal) => {
  if (newVal) {
    // 录制单键时, 清空输入框。(由于是录制, 因此需要清空输入框, 防止用户输入内容。)
    // * 如何防止用户输入内容?
    // * * 当然也可以利用updateInputValue。但有更简单的解决思路, 即定义组件特有属性maxlength为0即可阻止用户输入内容。
    singleKeysSelectRef.value?.updateInputValue('');

    if (!oldVal) {
      first_flag = false;
    }

    keyEvent_store.setKeyStateCallback_Record(recordingSingleKeysCallback);
  } else {
    keyEvent_store.clearKeyStateCallback_Record();
  }
});

watch(isGetsFocused, (newVal) => {
  if (newVal && isRecordingSingleKeys.value) {
    keyEvent_store.setKeyStateCallback_Record(recordingSingleKeysCallback);
  } else {
    keyEvent_store.clearKeyStateCallback_Record();
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
const preventDefaultMouseWhenRecording = (event: MouseEvent) => {
  // TIPS: 需要在mouseup事件中阻止鼠标按键4、5的前进后退功能。
  //                  sdk的button值  |  前端的event.button值
  // 'MouseLeft'            1                    0
  // 'MouseRight'           2                    2
  // 'MouseMiddle'          3                    1
  // 'MouseBack'            4                    3
  // 'MouseForward'         5                    4
  // console.log(event.button);

  // 鼠标按钮4是后退，按钮5是前进
  if (event.button === 3 || event.button === 4) {
    // 虽然无法录制'Enter'事件的原因就是select组件阻止了默认的'Enter'事件的冒泡行为,
    // * 但为防止quasar后续更新改变它, 便再次手动阻止一次, 以防止本次修复被quasar的更新影响。
    event.preventDefault(); // 阻止默认行为
    event.stopPropagation(); // 阻止事件冒泡
  }
};

// -- -- 选择声效
const isDownSoundEffectSelectEnabled = ref(true);
const isUpSoundEffectSelectEnabled = ref(true);

const keyDownSingleKeySoundEffectSelect = ref<any>();
const keyUpSingleKeySoundEffectSelect = ref<any>();
// const keyDownSingleKeySoundEffectSelect_diff = keyDownSingleKeySoundEffectSelect.value;
// const keyUpSingleKeySoundEffectSelect_diff = keyUpSingleKeySoundEffectSelect.value;
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
      message: $t('KeyToneAlbum.notify.noSoundEffectSelected'),
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
            message: $t('KeyToneAlbum.notify.singleKeySoundEffectConfigFailed', {
              key: keyEvent_store.dikCodeToName.get(item) || 'Dik-{' + item + '}',
            }),
            timeout: 5000,
          });
        }
        return re;
      })
      .catch((err) => {
        allSuccess = false;
        console.error(
          $t('KeyToneAlbum.notify.singleKeySoundEffectConfigError', {
            key: keyEvent_store.dikCodeToName.get(item) || 'Dik-{' + item + '}',
          }),
          err
        );
        q.notify({
          type: 'negative',
          position: 'top',
          message: $t('KeyToneAlbum.notify.singleKeySoundEffectConfigFailed', {
            key: keyEvent_store.dikCodeToName.get(item) || 'Dik-{' + item + '}',
          }),
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
        message: $t('KeyToneAlbum.notify.partialSingleKeySoundEffectConfigSuccess'),
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

// Dependency validation logic
const dependencyIssues = ref<DependencyIssue[]>([]);

// Computed property to get all dependency issues
const allDependencyIssues = computed(() => {
  const audioFiles = soundFileList.value as AudioFile[];
  const sounds = soundList.value as Sound[];
  const keySounds = keySoundList.value as KeySound[];
  
  if (audioFiles.length === 0 && sounds.length === 0 && keySounds.length === 0) {
    return [];
  }

  const validator = createDependencyValidator(audioFiles, sounds, keySounds);
  
  // Get global binding if it exists
  const globalBinding = keyDownUnifiedSoundEffectSelect.value || keyUpUnifiedSoundEffectSelect.value 
    ? {
        down: keyDownUnifiedSoundEffectSelect.value ? {
          type: keyDownUnifiedSoundEffectSelect.value.type,
          value: keyDownUnifiedSoundEffectSelect.value.type === 'audio_files' 
            ? keyDownUnifiedSoundEffectSelect.value.value
            : keyDownUnifiedSoundEffectSelect.value.type === 'sounds'
            ? keyDownUnifiedSoundEffectSelect.value.value.soundKey
            : keyDownUnifiedSoundEffectSelect.value.value.keySoundKey
        } : null,
        up: keyUpUnifiedSoundEffectSelect.value ? {
          type: keyUpUnifiedSoundEffectSelect.value.type,
          value: keyUpUnifiedSoundEffectSelect.value.type === 'audio_files'
            ? keyUpUnifiedSoundEffectSelect.value.value
            : keyUpUnifiedSoundEffectSelect.value.type === 'sounds'
            ? keyUpUnifiedSoundEffectSelect.value.value.soundKey
            : keyUpUnifiedSoundEffectSelect.value.value.keySoundKey
        } : null
      }
    : undefined;

  // Convert keysWithSoundEffect Map to the format expected by validator
  const singleKeyBindings = keysWithSoundEffect.value.size > 0 
    ? keysWithSoundEffect.value 
    : undefined;

  return validator.validateAllDependencies(globalBinding, singleKeyBindings);
});

// Update dependency issues when data changes
watch([soundFileList, soundList, keySoundList, keyDownUnifiedSoundEffectSelect, keyUpUnifiedSoundEffectSelect, keysWithSoundEffect], () => {
  dependencyIssues.value = allDependencyIssues.value;
}, { deep: true });

// Helper function to check if an item has dependency issues
const checkItemDependencyIssues = (itemType: 'audio_files' | 'sounds' | 'key_sounds', itemId: string) => {
  return hasItemDependencyIssues(itemType, itemId, dependencyIssues.value);
};
function convertValue(item: any) {
  if (item.type === 'audio_files') {
    return {
      type: 'audio_files',
      value: soundFileList.value.find(
        (soundFile) => item.value && soundFile.sha256 === item.value.sha256 && soundFile.name_id === item.value.name_id
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
  return null;
}

const isShowSingleKeySoundEffectEditDialog = ref(false);

const currentEditingKey = ref<number | null>(null);
let currentEditingKey_old = currentEditingKey.value; // TIPS: 此处旧值的记录, 在实际的使用逻辑中进行。
// TIPS: 如果currentEditingKey不变, 则watch就不会触发, 或者说这种做法无法实时记录正确的old值, 留此注释是警示下->不是所有场景都适合使用watch正确记录旧值的。
//       * 3 -> 2 后 old为3 ,  2 -> 3 后, old 为2。(不相同是符合预期的。)
//       * 3 -> 2 后 old为3 ,  2 -> 2 后, old 仍为3(但应该是2相同才对, 但只能判别不相同, 不符合预期, 或者说这种做法无法实时记录正确的old值)。
// watch(currentEditingKey, (newVal, oldVal) => {
//   currentEditingKey_old = oldVal;
// });
const currentEditingKeyOfName = computed(() => {
  return currentEditingKey.value !== null
    ? keyEvent_store.dikCodeToName.get(currentEditingKey.value) || 'Dik-{' + currentEditingKey.value + '}'
    : '';
});

// -- -- -- 编辑声效(重新选择声效)
const keyDownSingleKeySoundEffectSelect_edit = ref<any>();
const keyUpSingleKeySoundEffectSelect_edit = ref<any>();
let keyDownSingleKeySoundEffectSelect_edit_old = keyDownSingleKeySoundEffectSelect_edit.value;
let keyUpSingleKeySoundEffectSelect_edit_old = keyUpSingleKeySoundEffectSelect_edit.value;

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

// 存储事件监听器的引用，以便后续移除
let messageAudioPackageListener: (e: MessageEvent) => void;

onBeforeMount(async () => {
  // 此时由于是新建键音包, 因此是没有对应配置文件, 需要我们主动去创建的。 故第二个参数设置为true
  // 这也是我们加载页面前必须确定的事情, 否则无法进行后续操作, 一切以配置文件为前提。
  const audioPkgPath = (await LoadConfig(props.pkgPath, props.isCreate)).audioPkgPath;

  // 如果是创建键音包, 则需要执行一定的初始化工作。
  if (props.isCreate) {
    await ConfigSet('package_name', $t('KeyToneAlbum.new.name.defaultValue'));

    await ConfigSet('audio_pkg_uuid', props.pkgPath);

    await ConfigSet('key_tone.is_enable_embedded_test_sound', isEnableEmbeddedTestSound);

    main_store.GetKeyToneAlbumList(); // 更新键音包选择列表的名称。

    setting_store.mainHome.selectedKeyTonePkg = audioPkgPath;
  }
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
          message: $t('KeyToneAlbum.notify.configFileReadFailed'),
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

      if (data.key_tone?.is_enable_embedded_test_sound !== undefined) {
        isEnableEmbeddedTestSound.down = data.key_tone.is_enable_embedded_test_sound.down;
        isEnableEmbeddedTestSound.up = data.key_tone.is_enable_embedded_test_sound.up;
      }

      // TODO: 此逻辑未验证, 需要到编辑键音包界面才能验证
      if (data.key_tone?.single !== undefined) {
        keysWithSoundEffect.value.clear();
        Object.entries(data.key_tone.single).forEach(([dikCode, value]) => {
          // 只有 down/up 至少一个被正确设置且value不为空字符串时, 才算作 已设置单键声效的按键。
          if ((value as any)?.down?.value || (value as any)?.up?.value) {
            keysWithSoundEffect.value.set(dikCode, value);
          }
        });
      }
    });
    const updateKeyToneAlbumListName = debounce(
      () => {
        main_store.GetKeyToneAlbumList();
      },
      800,
      { trailing: true }
    );
    watch(pkgName, (newVal) => {
      ConfigSet('package_name', pkgName.value);
      updateKeyToneAlbumListName.cancel();
      updateKeyToneAlbumListName();
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

    if (keyTonePkgData.key_tone?.global !== undefined) {
      keyDownUnifiedSoundEffectSelect.value = convertValue(
        keyTonePkgData.key_tone.global.down ? keyTonePkgData.key_tone.global.down : ''
      );
      keyUpUnifiedSoundEffectSelect.value = convertValue(
        keyTonePkgData.key_tone.global.up ? keyTonePkgData.key_tone.global.up : ''
      );
    }

    if (keyTonePkgData.key_tone?.single !== undefined) {
      keysWithSoundEffect.value.clear();
      Object.entries(keyTonePkgData.key_tone.single).forEach(([dikCode, value]) => {
        // 只有 down/up 至少一个被正确设置且value不为空字符串时, 才算作 已设置单键声效的按键。
        if ((value as any)?.down?.value || (value as any)?.up?.value) {
          keysWithSoundEffect.value.set(dikCode, value);
        }
      });
    }
  }
  const debounced_sseDataToSettingStore = debounce<(keyTonePkgData: any) => void>(sseDataToKeyTonePkgData, 30, {
    trailing: true,
  });

  // 定义事件监听器
  messageAudioPackageListener = function (e) {
    console.debug('后端钩子函数中的值 = ', e.data);

    const data = JSON.parse(e.data);

    if (data.key === 'get_all_value') {
      debounced_sseDataToSettingStore.cancel;
      debounced_sseDataToSettingStore(data.value);
    }
  };

  // 添加事件监听
  app_store.eventSource.addEventListener('messageAudioPackage', messageAudioPackageListener, false);
});

// 在退出创建键音包的页面后, 载入 持久化的 用户选择的 键音包。(在 创建 键音包界面 退出时, 重新加载 用户持久化至 设置 文件中的 键音包。)
const main_store = useMainStore();
onUnmounted(() => {
  // 卸载组件后, 更新键音包列表
  main_store.GetKeyToneAlbumList();
  // // 卸载组件后, 重新载入持久化配置中用户所选的键音包(在新设计的键音专辑页面逻辑中, 不需要此步骤)
  // main_store.LoadSelectedKeyTonePkg();

  // 移除事件监听
  if (messageAudioPackageListener) {
    app_store.eventSource.removeEventListener('messageAudioPackage', messageAudioPackageListener);
  }
});

const isMacOS = ref(getMacOSStatus());
function getMacOSStatus() {
  if (process.env.MODE === 'electron') {
    return window.myWindowAPI.getMacOSStatus();
  }
  return false;
}

const i18n_fontSize = computed(() => {
  return setting_store.languageDefault === 'ru' ||
    setting_store.languageDefault === 'it' ||
    setting_store.languageDefault === 'es' ||
    setting_store.languageDefault === 'pt-BR' ||
    setting_store.languageDefault === 'pl' ||
    setting_store.languageDefault === 'tr' ||
    setting_store.languageDefault === 'id'
    ? '0.66rem'
    : setting_store.languageDefault === 'fr'
    ? '0.63rem'
    : '0.75rem';
});

const step_introduce_fontSize = computed(() => {
  return isMacOS.value
    ? // MacOS
      setting_store.languageDefault === 'ru' ||
      setting_store.languageDefault === 'ja' ||
      setting_store.languageDefault === 'es'
      ? 'text-[0.80rem]'
      : setting_store.languageDefault === 'ko-KR'
      ? 'text-[0.83rem]'
      : setting_store.languageDefault === 'id'
      ? 'text-[0.87]'
      : 'text-[0.85rem]'
    : // windows
    setting_store.languageDefault === 'ru' ||
      setting_store.languageDefault === 'ko-KR' ||
      setting_store.languageDefault === 'pl' ||
      setting_store.languageDefault === 'ar'
    ? 'text-[0.80rem]'
    : setting_store.languageDefault === 'tr'
    ? 'text-[0.83rem]'
    : 'text-[0.85rem]';
});
</script>

<style lang="scss" scoped>
// :deep(.q-stepper__tab) {
//   cursor: pointer;
// }

// :deep(.q-stepper__tab:hover) {
//   background-color: rgb(243 244 246);
// }

// TIPS: 注意 unocss 的默认预设中, 默认情况下是不支持以下这种特定的 tailwindcss 语法的 - 即 @apply 的用法不受支持。
//       * 需要通过手动更改相应的 配置文件uno.config.ts 来支持此转换语法。 -[参考链接](https://unocss.jiangruyi.com/transformers/directives)
:deep(.q-stepper__tab) {
  @apply cursor-pointer hover:bg-gray-100;
}

// 对选择器组件的label溢出情况, 采取滚动策略 (为防止刷新后样式丢失问题, 而加的。)
:deep(.q-field__native) {
  // 对溢出的情况, 采取滚动策略
  @apply max-w-full overflow-auto whitespace-nowrap;

  // 隐藏滚动策略的滚动条。
  // @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];

  // 添加细微滚动条
  @apply h-5.8 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50  [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400;
}

// // 对于多选的, 带芯片的选择框, 某芯片内的名称内容过长的情况, 采取溢出滚动的策略。
// :deep(.ellipsis) {
//   // 对溢出的情况, 采取滚动策略
//   @apply max-w-full overflow-auto whitespace-nowrap  text-clip;
//   // 隐藏滚动策略的滚动条。
//   @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];
// }

// 对本组件选择框添加的可清楚图标的大小做设置
:deep(.q-field__focusable-action) {
  @apply text-lg;
}

// 用于修复主页面全局的:global(.q-field__native)中的h-5.8这个样式影响了当前页面中的q-input的问题
:deep(.q-placeholder) {
  // 在这里重置q-input组件的输入样式的高度以修复这个问题
  @apply h-auto;
}

// 为防止刷新后样式丢失问题, 而加的。
:global(.q-item__section) {
  /* 对溢出的情况, 采取滚动策略 */
  @apply max-w-full overflow-auto whitespace-nowrap;

  /* 隐藏滚动策略的滚动条 */
  // @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];

  // 添加细微滚动条
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30  [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}

.zl-ll {
  :deep(.q-field__native) {
    @apply h-auto;
  }
  :deep(.q-field__messages) {
    @apply text-nowrap;
  }
}

:global(.q-card) {
  @apply mr-2.33;
}

// 对于键音专辑组件的选择框, 键音专辑的名称内容过长的情况, 采取溢出滚动的策略。
// 对于多选的, 带芯片的选择框, 某芯片内的名称内容过长的情况, 采取溢出滚动的策略。
:deep(.ellipsis) {
  // 对溢出的情况, 采取滚动策略
  @apply max-w-full overflow-auto whitespace-nowrap  text-clip;
  // // 隐藏滚动策略的滚动条。
  // @apply [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none];

  // 添加细微滚动条
  @apply [&::-webkit-scrollbar]:h-0.5 [&::-webkit-scrollbar-track]:bg-zinc-200/30  [&::-webkit-scrollbar-thumb]:bg-blue-500/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-600/50;
}

.q-btn {
  @apply text-xs;
  // font-size: 0.66rem;
  // line-height: 1rem;
  font-size: var(--i18n_fontSize);
  @apply p-1.5;
  @apply transition-transform hover:scale-105;
  @apply scale-103;
}

.step-custom {
  :deep(.q-stepper__step-inner) {
    padding: 0 10px 32px 55px;
  }
}
:deep(.q-field__label) {
  @apply overflow-visible -ml-1.5 text-[0.8rem];
}
</style>
