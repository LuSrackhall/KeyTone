<template>
  <div :class="[isMacOS ? 'w-[389.5px] h-[458.5px]' : 'w-[379px] h-[458.5px]', 'overflow-hidden']">
    <!-- 展开按钮 -->
    <div v-if="isCollapsed" class="fixed top-6 transform -translate-x-1/2 z-50" style="left: 50%">
      <q-btn flat class="custom-expand-btn" @click="isCollapsed = false">
        <div class="chevron-down"></div>
      </q-btn>
    </div>

    <!-- 选择器容器使用绝对定位 -->
    <div class="relative">
      <transition name="slide">
        <div
          v-show="!isCollapsed"
          class="selector-container absolute w-[88%] ml-[6.2%] mr-[5.8%]"
          style="top: 0; left: 0; right: 0; z-index: 1"
        >
          <!-- 优化按钮组布局 -->
          <div class="flex items-center gap-1.5 mt-0.5 px-1">
            <q-btn
              flat
              dense
              round
              size="xs"
              icon="add"
              color="primary"
              class="w-6.5 h-6.5 opacity-60 transition-all duration-200 ease-[cubic-bezier(0.4,0,0.2,1)] bg-white/10 backdrop-blur hover:opacity-100 hover:-translate-y-px hover:bg-white/15 disabled:opacity-30 disabled:transform-none disabled:cursor-not-allowed"
              @click="createNewAlbum"
            >
              <q-tooltip
                anchor="bottom middle"
                self="top middle"
                :offset="[0, 8]"
                class="rounded-lg text-[0.8rem] px-3 py-1.2"
              >
                {{ $t('keyToneAlbumPage.new') }}
              </q-tooltip>
            </q-btn>

            <q-btn
              flat
              dense
              round
              size="xs"
              icon="save_alt"
              color="primary"
              class="w-6.5 h-6.5 opacity-60 transition-all duration-200 ease-[cubic-bezier(0.4,0,0.2,1)] bg-white/10 backdrop-blur hover:opacity-100 hover:-translate-y-px hover:bg-white/15 disabled:opacity-30 disabled:transform-none disabled:cursor-not-allowed"
              @click="importAlbum"
            >
              <q-tooltip
                anchor="bottom middle"
                self="top middle"
                :offset="[0, 8]"
                class="rounded-lg text-[0.8rem] px-3 py-1.2"
              >
                {{ $t('keyToneAlbumPage.import') }}
              </q-tooltip>
            </q-btn>

            <q-btn
              flat
              dense
              round
              size="xs"
              icon="drive_file_move"
              color="primary"
              :disable="!setting_store.mainHome.selectedKeyTonePkg"
              class="w-6.5 h-6.5 opacity-60 transition-all duration-200 ease-[cubic-bezier(0.4,0,0.2,1)] bg-white/10 backdrop-blur hover:opacity-100 hover:-translate-y-px hover:bg-white/15 disabled:opacity-30 disabled:transform-none disabled:cursor-not-allowed"
              @click="exportAlbum"
            >
              <q-tooltip
                anchor="bottom middle"
                self="top middle"
                :offset="[0, 8]"
                class="rounded-lg text-[0.8rem] px-3 py-1.2"
              >
                {{ $t('keyToneAlbumPage.export') }}
              </q-tooltip>
            </q-btn>

            <q-btn
              flat
              dense
              round
              size="xs"
              icon="delete"
              color="negative"
              :disable="!setting_store.mainHome.selectedKeyTonePkg"
              class="w-6.5 h-6.5 opacity-60 transition-all duration-200 ease-[cubic-bezier(0.4,0,0.2,1)] bg-white/10 backdrop-blur hover:opacity-100 hover:-translate-y-px hover:bg-white/15 disabled:opacity-30 disabled:transform-none disabled:cursor-not-allowed"
              @click="deleteAlbum"
            >
              <q-tooltip
                anchor="bottom middle"
                self="top middle"
                :offset="[0, 8]"
                class="rounded-lg text-[0.8rem] px-3 py-1.2"
              >
                {{ $t('keyToneAlbumPage.delete') }}
              </q-tooltip>
            </q-btn>

            <!-- 新增帮助按钮 -->
            <q-btn
              flat
              dense
              round
              size="xs"
              icon="help"
              color="primary"
              class="w-6.5 h-6.5 opacity-60 transition-all duration-200 ease-[cubic-bezier(0.4,0,0.2,1)] bg-white/10 backdrop-blur hover:opacity-100 hover:-translate-y-px hover:bg-white/15 disabled:opacity-30 disabled:transform-none disabled:cursor-not-allowed"
              @click="openExternal('https://keytone.xuanhall.com/guide/key-package/introduction/')"
            >
              <q-tooltip
                anchor="bottom middle"
                self="top middle"
                :offset="[0, 8]"
                class="rounded-lg text-[0.8rem] px-3 py-1.2"
              >
                {{ $t('keyToneAlbumPage.help') }}
              </q-tooltip>
            </q-btn>
          </div>

          <!-- 键音专辑选择器 -->
          <div class="relative">
            <q-select
              v-model="setting_store.mainHome.selectedKeyTonePkg"
              :options="main_store.keyTonePkgOptions"
              :option-label="(item: any) => main_store.keyTonePkgOptionsName.get(item)"
              :label="$t('keyToneAlbumPage.label')"
              :virtual-scroll-slice-size="999999"
              outlined
              dense
              emit-value
              map-options
              ref="selectedKeyTonePkgRef"
              @popup-hide="blur()"
              :disable="
                (() => {
                  // 在键音专辑创建期间, 应禁止选择器的使用, 避免意外选择其它键音专辑造成创建被中断, 以及其它混乱问题。
                  return keytoneAlbum_store.isCreateNewKeytoneAlbum;
                })()
              "
              popup-content-class="w-[1%] whitespace-normal break-words [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-zinc-200/30 [&::-webkit-scrollbar-thumb]:bg-zinc-900/30 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-zinc-900/50"
            >
              <template v-if="setting_store.mainHome.selectedKeyTonePkg" v-slot:append>
                <q-icon
                  name="cancel"
                  @click.stop.prevent="setting_store.mainHome.selectedKeyTonePkg = ''"
                  class="cursor-pointer text-lg"
                />
              </template>

              <!-- 空状态提示 -->
              <template v-slot:no-option>
                <div class="flex flex-col items-center py-6 text-gray-500">
                  <q-icon name="library_music" size="56px" class="empty-state-icon mb-3" />
                  <div class="text-sm mb-5 opacity-75">{{ $t('keyToneAlbumPage.emptyState.noAlbum') }}</div>
                  <div class="flex gap-4">
                    <q-btn
                      flat
                      dense
                      class="empty-state-btn flex items-center bg-blue-500/10 p-1.5 rounded-lg"
                      @click="createNewAlbum"
                    >
                      <q-icon name="add" size="20px" class="mr-1.5" />
                      <span class="font-medium" :class="[i18n_fontSize]">{{ $t('keyToneAlbumPage.new') }}</span>
                    </q-btn>
                    <q-btn
                      flat
                      dense
                      class="empty-state-btn flex items-center bg-blue-500/10 p-1.5 rounded-lg"
                      @click="importAlbum"
                    >
                      <q-icon name="upload_file" size="20px" class="mr-1.5" />
                      <span class="font-medium" :class="[i18n_fontSize]">{{ $t('keyToneAlbumPage.import') }}</span>
                    </q-btn>
                  </div>
                </div>
              </template>
            </q-select>
          </div>

          <q-btn
            v-if="setting_store.mainHome.selectedKeyTonePkg"
            flat
            round
            color="grey"
            icon="expand_less"
            class="collapse-btn absolute -bottom-6 transform -translate-x-1/2"
            style="left: 50%"
            @click="isCollapsed = true"
          />
        </div>
      </transition>

      <!-- 内容区域添加过渡padding -->
      <div
        class="content-wrapper"
        :style="{
          paddingTop: !isCollapsed ? '68px' : '0',
          transition: `padding-top ${!isCollapsed ? '0.8s' : '1.2s'} ease`,
        }"
      >
        <div :class="{ 'hide-scrollbar': isAtTop }" class="keytone-album-container">
          <div class="relative min-h-[360px]">
            <KeytoneAlbum
              v-if="keytoneAlbum_PathOrUUID"
              v-show="
                (() => {
                  // 如果执行创建, 则创建完毕前, 不展示组件。(此时展示一个可以表示加载中的组件)
                  return !keytoneAlbum_store.isCreateNewKeytoneAlbum;
                })()
              "
              :key="
                (() => {
                  // 为了避免键音包专辑的重复渲染, 这里使用键音包专辑UUID作为key。 (主要造成重复渲染的场景->新建键音专辑时的两个阶段都会更改 keytoneAlbum_PathOrUUID)
                  // 这个表达式就是为了将路径转换为UUID。 此表达式不会对本身就是UUID的情况造成影响。
                  const key_UUID = keytoneAlbum_PathOrUUID.split(q.platform.is.win ? '\\' : '/').pop();
                  console.log('key_UUID', key_UUID);
                  return key_UUID;
                })()
              "
              :pkgPath="keytoneAlbum_PathOrUUID"
              :isCreate="keytoneAlbum_store.isCreateNewKeytoneAlbum"
              ref="keytoneAlbumRef"
            />
            <div
              v-if="keytoneAlbum_store.isCreateNewKeytoneAlbum"
              class="absolute inset-0 flex items-center justify-center bg-transparent"
            >
              <!-- <q-spinner-tail color="primary" size="32px" /> -->
              <!-- <q-spinner-pie color="primary" size="56px" /> -->
              <!-- <q-spinner-gears color="primary" size="56px" /> -->
              <q-spinner-cube color="primary" size="56px" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { QSelect, useQuasar } from 'quasar';
import { useMainStore } from 'src/stores/main-store';
import { useSettingStore } from 'src/stores/setting-store';
import { nanoid } from 'nanoid';
import { computed, nextTick, useTemplateRef } from 'vue';
import KeytoneAlbum from 'src/components/Keytone_album.vue';
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  DeleteAlbum,
  GetAudioPackageName,
  ExportAlbum,
  ImportAlbum,
  ImportAlbumOverwrite,
  ImportAlbumAsNew,
  LoadConfig,
  GetAlbumMeta,
  type AlbumMeta,
} from 'src/boot/query/keytonePkg-query';

// 扩展HTMLInputElement类型以支持webkitdirectory属性
declare global {
  interface HTMLInputElement {
    webkitdirectory: boolean;
  }
}
import { useKeytoneAlbumStore } from 'src/stores/keytoneAlbum-store';

const q = useQuasar();
const { t } = useI18n();
const $t = t;
const main_store = useMainStore();
const setting_store = useSettingStore();
const keytoneAlbum_store = useKeytoneAlbumStore();
const selectedKeyTonePkgRef = useTemplateRef<QSelect>('selectedKeyTonePkgRef');
const keytoneAlbumRef = ref<InstanceType<typeof KeytoneAlbum> | null>(null);
const isCollapsed = ref(false);
let lastScrollTop = 0;
const isAtTop = ref(true);

// 键音包的创建与编辑
// // 键音包的创建分为两个阶段
// // * 第一个阶段只有目录UUID名称, 而没有完整路径的阶段, 此时将UUID传入后端创建api, 并获取完整路径的返回值。
// // * 当创建完毕后, 进入第二阶段, 将完整的路径传递给用户已选择键音包的变量(即setting_store.mainHome.selectedKeyTonePkg)。
// // * 两个阶段完成后, 即创建成功。(实际上第一阶段完成就算创建成功, 第二阶段仅影响前端展示)
// // * 如果第一阶段进行了一般, 即将UUID传入了后端api但未进行获取返回值等后续步骤, 则存在失败的可能。
const keytoneAlbum_PathOrUUID = ref<string>(setting_store.mainHome.selectedKeyTonePkg); // 用于向KeytoneAlbum组件传递键音包的路径或UUID

// 实现删除专辑的逻辑
const deleteAlbum = async () => {
  if (!setting_store.mainHome.selectedKeyTonePkg) return;

  const albumPath = setting_store.mainHome.selectedKeyTonePkg;
  const result = await DeleteAlbum(albumPath);

  if (result) {
    setting_store.mainHome.selectedKeyTonePkg = '';
    q.notify({
      type: 'positive',
      message: $t('keyToneAlbumPage.notify.deleteSuccess'),
    });
  } else {
    q.notify({
      type: 'negative',
      message: $t('keyToneAlbumPage.notify.deleteFailed'),
    });
  }
};

// 实现新建专辑的逻辑
const createNewAlbum = () => {
  if (keytoneAlbum_store.isCreateNewKeytoneAlbum) return; // 若上个专辑创建未完成, 则不允许再次创建(无需提示用户, )
  keytoneAlbum_store.isCreateNewKeytoneAlbum = true;
  keytoneAlbum_PathOrUUID.value = nanoid();
};
watch(
  () => setting_store.mainHome.selectedKeyTonePkg,
  () => {
    keytoneAlbum_store.isCreateNewKeytoneAlbum = false; // 避免递归创建
    keytoneAlbum_PathOrUUID.value = setting_store.mainHome.selectedKeyTonePkg;
  }
);

// 声明 File System Access API 相关类型
declare global {
  interface Window {
    showSaveFilePicker: (options?: {
      suggestedName?: string;
      types?: Array<{
        description: string;
        accept: Record<string, string[]>;
      }>;
    }) => Promise<FileSystemFileHandle>;
  }
}

interface FileSystemCreateWritableOptions {
  keepExistingData?: boolean;
}

interface FileSystemFileHandle {
  createWritable: (options?: FileSystemCreateWritableOptions) => Promise<FileSystemWritableFileStream>;
}

interface FileSystemWritableFileStream extends WritableStream {
  write: (data: BufferSource | Blob | string) => Promise<void>;
  close: () => Promise<void>;
}

// 降级方案 - 使用传统的下载方式
const exportAlbumLegacy = async () => {
  try {
    // 获取专辑名称
    const albumNameResponse = await GetAudioPackageName(setting_store.mainHome.selectedKeyTonePkg);
    if (!albumNameResponse || albumNameResponse.message !== 'ok') {
      throw new Error('获取专辑名称失败');
    }
    const albumName = albumNameResponse.name;

    // 调用导出函数获取zip文件blob
    const blob = await ExportAlbum(setting_store.mainHome.selectedKeyTonePkg);
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `${albumName}.ktalbum`; // 改为 .ktalbum
    document.body.appendChild(link);
    link.click();

    // 清理
    window.URL.revokeObjectURL(url);
    document.body.removeChild(link);

    q.notify({
      type: 'positive',
      message: $t('keyToneAlbumPage.notify.exportSuccess'),
    });
  } catch (error) {
    console.error('导出专辑失败:', error);
    q.notify({
      type: 'negative',
      message:
        $t('keyToneAlbumPage.notify.exportFailed') + ': ' + (error instanceof Error ? error.message : String(error)),
    });
  }
};

// 导出键音专辑 - 使用 File System Access API
const exportAlbum = async () => {
  // 检查 API 是否可用
  if (typeof window.showSaveFilePicker !== 'function') {
    console.log('Browser does not support File System Access API, falling back to legacy export');
    return exportAlbumLegacy();
  }

  try {
    // 获取专辑名称
    const albumNameResponse = await GetAudioPackageName(setting_store.mainHome.selectedKeyTonePkg);
    if (!albumNameResponse || albumNameResponse.message !== 'ok') {
      throw new Error('获取专辑名称失败');
    }
    const albumName = albumNameResponse.name;

    // 获取导出数据
    const blob = await ExportAlbum(setting_store.mainHome.selectedKeyTonePkg);

    try {
      // 打开系统的保存文件对话框
      const handle = await window.showSaveFilePicker({
        suggestedName: `${albumName}.ktalbum`,
        types: [
          {
            description: $t('keyToneAlbumPage.notify.fileDescription'),
            accept: { 'application/octet-stream': ['.ktalbum'] },
          },
        ],
      });

      // 写入文件
      const writable = await handle.createWritable();
      await writable.write(blob);
      await writable.close();

      // 文件成功保存后再通知
      q.notify({
        type: 'positive',
        message: $t('keyToneAlbumPage.notify.exportSuccess'),
      });
    } catch (err) {
      // 用户取消选择文件时不显示错误
      if (err instanceof Error && err.name === 'AbortError') {
        return;
      }
      throw err;
    }
  } catch (error) {
    console.error('导出专辑失败:', error);
    q.notify({
      type: 'negative',
      message:
        $t('keyToneAlbumPage.notify.exportFailed') + ': ' + (error instanceof Error ? error.message : String(error)),
    });
  }
};

const blur = () => {
  setTimeout(() => {
    selectedKeyTonePkgRef?.value?.blur();
    // TIPS: 这里需要延迟后再blur, 以确保blur的正确触发(太早触发blur会不起作用, 经验证, 本人电脑延迟10ms后, 可以正确触发blur使焦点丧失, 为确保适配更多的低性能设备, 这里保险起见设置为66ms)
  }, 66);
};

console.log('main_store.keyTonePkgOptions', main_store.keyTonePkgOptions);

// 监听 KeytoneAlbum 内部的滚动
const handleAlbumScroll = (event: Event) => {
  const scrollableElement = (event.target as HTMLElement).closest('.q-scrollarea__container');
  if (!scrollableElement) return;

  const currentScroll = scrollableElement.scrollTop;
  const maxScroll = scrollableElement.scrollHeight - scrollableElement.clientHeight;

  // 更新是否在顶部的状态
  isAtTop.value = currentScroll === 0;

  // 在顶部继续向上滚动时展开
  if (currentScroll === 0 && event instanceof WheelEvent && event.deltaY < 0 && isCollapsed.value) {
    isCollapsed.value = false;
    return;
  }

  // 向下滚动时立即收起
  if (
    (currentScroll > lastScrollTop ||
      (currentScroll >= maxScroll && event instanceof WheelEvent && event.deltaY > 0)) &&
    !isCollapsed.value
  ) {
    isCollapsed.value = true;
  }

  lastScrollTop = currentScroll;
};

// 添加事件监听器的函数
const setupScrollListeners = () => {
  // 给一点延时确保 DOM 已更新
  setTimeout(() => {
    const scrollContainer = keytoneAlbumRef.value?.$el.querySelector('.q-scrollarea__container');
    if (scrollContainer) {
      // 先移除可能存在的旧监听器
      scrollContainer.removeEventListener('scroll', handleAlbumScroll);
      scrollContainer.removeEventListener('wheel', handleAlbumScroll);
      // 添加新的监听器
      scrollContainer.addEventListener('scroll', handleAlbumScroll, { passive: true });
      scrollContainer.addEventListener('wheel', handleAlbumScroll, { passive: true });
    }
  }, 100);
};

// 监听键音包变化
watch(
  () => setting_store.mainHome.selectedKeyTonePkg,
  () => {
    setupScrollListeners();
  }
);

onMounted(() => {
  setupScrollListeners();
});

onUnmounted(() => {
  const scrollContainer = keytoneAlbumRef.value?.$el.querySelector('.q-scrollarea__container');
  if (scrollContainer) {
    scrollContainer.removeEventListener('scroll', handleAlbumScroll);
    scrollContainer.removeEventListener('wheel', handleAlbumScroll);
  }
});

const importAlbum = async () => {
  // 创建文件输入元素
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = '.ktalbum';

  // 处理文件选择
  input.onchange = async (e) => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;

    // 检查文件类型
    if (!file.name.toLowerCase().endsWith('.ktalbum')) {
      q.notify({
        type: 'negative',
        message: $t('keyToneAlbumPage.notify.invalidFormat'),
      });
      return;
    }

    try {
      const result = await ImportAlbum(file);
      if (result) {
        q.notify({
          type: 'positive',
          message: $t('keyToneAlbumPage.notify.importSuccess'),
        });
        // 刷新专辑列表
        await main_store.GetKeyToneAlbumList();
      }
    } catch (error) {
      // 处理专辑已存在的情况
      if (error instanceof Error && error.message === 'album_exists') {
        // 使用 q.dialog 显示选项对话框
        q.dialog({
          title: $t('keyToneAlbumPage.importDialog.title'),
          message: $t('keyToneAlbumPage.importDialog.message'),
          persistent: true,
          options: {
            type: 'radio',
            model: 'overwrite',
            items: [
              {
                label: $t('keyToneAlbumPage.importDialog.overwrite'),
                value: 'overwrite',
                color: 'negative',
              },
              {
                label: $t('keyToneAlbumPage.importDialog.saveAsNew'),
                value: 'new',
                color: 'primary',
              },
              {
                label: $t('keyToneAlbumPage.importDialog.cancel'),
                value: 'cancel',
                color: 'primary',
              },
            ],
          },
          ok: {
            label: $t('keyToneAlbumPage.importDialog.confirm'),
            push: true,
            color: 'primary',
          },
        }).onOk(async (data) => {
          if (data === 'cancel') {
            return;
          }
          try {
            let result = false;

            if (data === 'overwrite') {
              // 用户选择覆盖
              result = await ImportAlbumOverwrite(file);
              if (result) {
                q.notify({
                  type: 'positive',
                  message: $t('keyToneAlbumPage.notify.importOverwriteSuccess'),
                });
                try {
                  // 刷新专辑列表
                  await main_store.GetKeyToneAlbumList();

                  // 只有当覆盖的是当前选中的专辑时,才需要重新加载
                  const currentAlbumUUID = setting_store.mainHome.selectedKeyTonePkg
                    ?.split(q.platform.is.win ? '\\' : '/')
                    .pop();
                  const meta = await GetAlbumMeta(file);

                  if (currentAlbumUUID === meta.albumUUID) {
                    console.log('检测到覆盖当前选中专辑,正在重新加载...');
                    // 重新加载选中的键音包(触发sdk重新加载键音包)
                    const loadResult = await LoadConfig(setting_store.mainHome.selectedKeyTonePkg, false);
                    if (loadResult) {
                      // 重新设置选中的键音包(触发ui重新渲染)
                      keytoneAlbum_PathOrUUID.value = '';
                      await nextTick();
                      keytoneAlbum_PathOrUUID.value = setting_store.mainHome.selectedKeyTonePkg;
                    } else {
                      throw new Error('重新加载键音包失败');
                    }
                  } else {
                    console.log('覆盖的不是当前选中专辑,无需重新加载');
                  }
                } catch (error) {
                  console.error('刷新专辑数据失败:', error);
                  q.notify({
                    type: 'warning',
                    message: $t('keyToneAlbumPage.notify.refreshFailed'),
                  });
                }
              }
            } else if (data === 'new') {
              // 用户选择保存为新专辑
              const newAlbumId = nanoid();
              result = await ImportAlbumAsNew(file, newAlbumId);
              if (result) {
                q.notify({
                  type: 'positive',
                  message: $t('keyToneAlbumPage.notify.importSaveAsNewSuccess'),
                });
                // 刷新专辑列表
                await main_store.GetKeyToneAlbumList();
              }
            }

            if (!result) {
              throw new Error('导入失败');
            }
          } catch (err) {
            q.notify({
              type: 'negative',
              message:
                $t('keyToneAlbumPage.notify.importFailed') + ': ' + (err instanceof Error ? err.message : String(err)),
            });
          }
        });
      } else {
        // 处理其他错误
        q.notify({
          type: 'negative',
          message:
            $t('keyToneAlbumPage.notify.importFailed') +
            ': ' +
            (error instanceof Error ? error.message : String(error)),
        });
      }
    }
  };

  // 触发文件选择
  input.click();
};

const openExternal = (url: string) => {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.openExternal(url);
  }
};

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
    ? 'text-[0.63rem]'
    : setting_store.languageDefault === 'fr'
    ? 'text-[0.63rem]'
    : 'text-[0.75rem] px-1.5';
});
</script>

<style lang="scss" scoped>
.selector-container {
  transform-origin: top;
}

// 滑动动画
.slide-enter-active,
.slide-leave-active {
  transition: all 0.8s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateY(-100%);
  opacity: 0;
}

.slide-enter-to,
.slide-leave-from {
  transform: translateY(0);
  opacity: 1;
}

.collapse-btn {
  position: absolute;
  z-index: 1;
  min-height: 24px;
  backdrop-filter: blur(4px);
  border-radius: 12px;
}

.custom-expand-btn {
  min-height: 24px;
  width: 64px;
  padding: 0;
  background: rgba(100, 100, 100, 0.3);
  backdrop-filter: blur(4px);
  border-radius: 12px;
  margin-top: 8px;

  &:hover {
    background: rgba(100, 100, 100, 0.3);
  }
}

.chevron-down {
  width: 10px;
  height: 10px;
  border-right: 2px solid white;
  border-bottom: 2px solid white;
  transform: rotate(45deg);
  margin-top: -4px;
}

.content-wrapper {
  will-change: padding-top;
  position: relative;
}

.keytone-album-container {
  &.hide-scrollbar :deep(.q-scrollarea__thumb) {
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  :deep(.q-scrollarea__thumb) {
    transition: opacity 0.3s ease;
  }
}

// 空状态相关样式
:deep(.q-select) {
  .q-field__control {
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(8px);
    transition: all 0.3s ease;

    &:hover {
      background: rgba(255, 255, 255, 0.08);
    }
  }
}

// 空状态动画
:deep(.empty-state-icon) {
  animation: float 3s ease-in-out infinite;
  opacity: 0.5;
  transition: opacity 0.3s ease;

  &:hover {
    opacity: 0.8;
  }
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-6px);
  }
}

// 空状态按钮样式
:deep(.empty-state-btn) {
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;

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

  &:hover {
    transform: translateY(-2px);

    &::before {
      transform: translate(-50%, -50%) scale(2);
      opacity: 1;
    }
  }

  &:active {
    transform: translateY(0);
  }
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

// 对于键音专辑页面键音专辑的选择框, 键音专辑名称内容过长的情况, 采取溢出滚动的策略。
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
