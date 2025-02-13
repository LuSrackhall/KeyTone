<template>
  <!--
   q-layout 有个 style="min-height: 835.2px" 的样式, 会造成滚动条的出现 进而无法展示窗口底部 进而无法通过css实现圆角
   因此, 我们这里主动设置 style="min-height: 0px"
  -->
  <q-layout view="hHh lpR lFr" style="min-height: 0px" class="sizeChange">
    <q-header elevated class="bg-primary text-white rounded-t-lg" style="width: calc(100% - 10px); right: 5px">
      <!-- TODO: 工具栏后续可以做的更丰富
             - completed(已完成)    1. 可以根据非'/'路由路径, 使back按钮独立至Menu右面
             - 过分设计, 无需解决   2. 可以对主题Label标题, 添加点击事件, 使其点击后可直接跳转至主页面路由路径'/'  (TIPS: 简单应用涉及的路径太浅, 没必要)
             - completed(已完成)    3. 可以对每次页面变更, 在工具栏中心位置展示页面Label (此时可以把back按钮, 放到Label的左边, 感觉会更好)
             - 过分设计, 无需解决   4. 更进一步, 可以整一个更完善的路由Label, 如 `设置`-`语言设置` + 左右滚动的窗口, 以让用户点击路径上任一节点即可跳转。
                                      * 实现步骤也很容易, 只需要将 pageLabel 转换为对象, 其中包含相关对应的 to 属性即可。 ((TIPS: 简单应用涉及的路径太浅, 没必要))
      -->
      <q-bar class="q-electron-drag rounded-t-lg">
        <q-btn
          :disable="
            (() => {
              // 在键音专辑创建期间, 应禁止选择器的使用, 避免意外选择其它键音专辑造成创建被中断, 以及其它混乱问题。
              return keytoneAlbum_store.isCreateNewKeytoneAlbum;
            })()
          "
          flat
          dense
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />
        <!-- <q-btn flat dense round icon="keyboard_alt"></q-btn> -->
        <div class="flex flex-nowrap">
          <q-btn
            :disable="
              (() => {
                // 在键音专辑创建期间, 应禁止选择器的使用, 避免意外选择其它键音专辑造成创建被中断, 以及其它混乱问题。
                return keytoneAlbum_store.isCreateNewKeytoneAlbum;
              })()
            "
            flat
            dense
            no-caps
            stretch
            size="16px"
            :class="[
              'm-0 p-0 text font-normal',

              // 使页面文本无法选择, 免得影响界面体验。
              'select-none',
            ]"
            @click="home"
          >
            {{ $t('KeyTone.KeyTone') }}
          </q-btn>

          <!-- <div class="ml-1.3 text-xs pt-2.8">
            <div class="flex p-l-1.5 p-r-1.5 outline outline-1 rounded outline-slate-300 text-slate-200 text-[11px]">
              <div class="mr-0.8 text-[9px]">V</div>
              {{ version }}
            </div>
          </div> -->
          <div class="text-xs pt-2.6">
            <div class="ml-2 flex text-[12px]">
              {{ version }}
            </div>
          </div>
        </div>

        <q-space />
        <q-btn
          :disable="
            (() => {
              // 在键音专辑创建期间, 应禁止选择器的使用, 避免意外选择其它键音专辑造成创建被中断, 以及其它混乱问题。
              return keytoneAlbum_store.isCreateNewKeytoneAlbum;
            })()
          "
          v-if="router.currentRoute.value.fullPath !== '/'"
          dense
          flat
          icon="keyboard_backspace"
          @click="back"
          class="mt-0.8 -mr-1.8"
        />
        <div class="flex flex-nowrap">
          <div class="text-sm text-nowrap">{{ pageLabel }}</div>
        </div>
        <q-space />

        <q-btn
          :disable="
            (() => {
              // 在键音专辑创建期间, 应禁止选择器的使用, 避免意外选择其它键音专辑造成创建被中断, 以及其它混乱问题。
              return keytoneAlbum_store.isCreateNewKeytoneAlbum;
            })()
          "
          dense
          flat
          icon="horizontal_rule"
          @click="minimize"
        />
        <!-- <q-btn dense flat icon="crop_square" @click="toggleMaximize" /> -->
        <q-btn
          :disable="
            (() => {
              // 在键音专辑创建期间, 应禁止选择器的使用, 避免意外选择其它键音专辑造成创建被中断, 以及其它混乱问题。
              return keytoneAlbum_store.isCreateNewKeytoneAlbum;
            })()
          "
          dense
          flat
          icon="close"
          @click="closeApp"
        />
      </q-bar>
    </q-header>

    <!-- 实际上的抽屉在q-layout中, 因此将.sizeChange这个css放到那边而不是这里 -->
    <q-drawer v-model="leftDrawerOpen" side="left" overlay behavior="desktop" elevated>
      <!-- draggable="false"使得图片/文本类元素组合,无法拖动, 免得影响界面体验 -->
      <q-list
        :class="[
          // 使页面文本无法选择, 免得影响界面体验。
          'select-none',
        ]"
        draggable="false"
      >
        <q-item-label header :class="['flex pl-4 pt-3 pb-1']">
          <!-- 自己名字这里, 到时候可以放置个链接属性, 让用户可以通过点击来启动系统默认浏览器进入我的github主页 -->
          <div :class="['h-5']">{{ $t('KeyTone.developer') }}:</div>
          <q-item
            clickable
            :class="['p-0 m-l-2 shadowSizeChange', 'h-5 min-h-0']"
            @click="openExternal('https://github.com/LuSrackhall')"
            >LuSrackhall</q-item
          >
        </q-item-label>
        <EssentialLink
          v-for="link in essentialLinks"
          :key="link.title"
          v-bind="link"
          draggable="false"
          @click="
            () => {
              // 每次点击某个项就进入对应页面的同时, 默认情况下会关闭drawer, 但当本身就在所需跳转页面时, 再次点击, 则不会关闭drawer。
              // * 推测后: 任为原因可能是此时未发生跳转的路由重置, 因此drawer不会被关闭。
              // 最终解决方案: 无论原因如何, 我们都在点击后, 触发关闭drawer的逻辑。获得预期行为。
              leftDrawerOpen = false;
            }
          "
        />
      </q-list>
      <div class="q-mini-drawer-hide absolute" style="top: 15px; right: -17px">
        <q-btn dense round unelevated color="accent" icon="chevron_left" @click="toggleLeftDrawer" />
      </div>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import EssentialLink, { EssentialLinkProps } from 'components/EssentialLink.vue';
import { useI18n } from 'vue-i18n';

const keytoneAlbum_store = useKeytoneAlbumStore();

const { t } = useI18n();
const $t = t;

const version = process.env.APP_VERSION;

const leftDrawerOpen = ref(false);

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}

function minimize() {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.minimize();
  }
}

function toggleMaximize() {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.toggleMaximize();
  }
}

function closeApp() {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.close();
  }
}

function openExternal(url: string) {
  if (process.env.MODE === 'electron') {
    window.myWindowAPI.openExternal(url);
  }
}

const essentialLinks: Array<EssentialLinkProps> = [
  {
    title: 'KeyTone.KeyToneAlbum.index',
    caption: '',
    icon: 'album',
    to: '/keytone_album',
  },
  {
    title: 'KeyTone.setting.index',
    caption: '',
    icon: 'settings',
    to: '/setting',
  },
];

import { useRouter } from 'vue-router';
import { useKeytoneAlbumStore } from 'src/stores/keytoneAlbum-store';

const router = useRouter();

// 获取当前路由路径, 通过打印观察。
// console.log(router.currentRoute.value.fullPath);

const back = () => {
  router.back();
};

const pageLabel = computed(() => {
  if (router.currentRoute.value.fullPath.split('-')[0] === '/setting') {
    return $t('KeyTone.setting.index');
  } else if (router.currentRoute.value.fullPath.split('-')[0] === '/keytone_album') {
    return $t('KeyTone.KeyToneAlbum.index');
  } else {
    return '';
  }
});

const home = () => {
  router.push('/');
  // 回到主页的操作, 默认关闭抽屉
  leftDrawerOpen.value = false;
};
</script>

<style lang="scss" scoped>
/* 抽屉尺寸适配 */
.sizeChange {
  :deep(.q-drawer) {
    height: calc(100% - 42px);
    border-bottom-left-radius: 0.25rem /* 4px */;
  }
}

/* 悬停阴影适配 */
.shadowSizeChange {
  :deep(.q-focus-helper) {
    height: 18px;
  }
}
</style>
