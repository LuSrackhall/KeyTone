<template>
  <q-item>
    <!-- select左边的竖线 -->
    <div :class="['ml-6 rounded-full  border-l-solid border-l-5 mr-6 h-6 self-center']"></div>
    <q-select
      v-model="setting_store.languageDefault"
      :options="localeOptions"
      :label="$t('language.setting language')"
      dense
      borderless
      emit-value
      map-options
      options-dense
      :class="['', 'min-w-[160px] w-[100%]']"
    />
  </q-item>
</template>

<script setup lang="ts">
import { useSettingStore } from 'src/stores/setting-store';

interface LocationOption {
  value: string;
  label: string;
}

/**
 * 由开发者LuSrackhall指定本软件目前已支持的语言设置项
 */
const localeOptions: Array<LocationOption> = [
  { value: 'en-US', label: 'English' },
  { value: 'zh-CN', label: '中文-简体' },
];

/**
 * 用于存储用户手动设置的变量值, 自动获取LocationOption对象的 value 值, 附给v-model(也就是setting_store.languageDefault)
 * * 手动更改设置将会引起setting_store.languageDefault变量发生变更, 从而被setting_store中设置的watch监听到
 * * 监听到之后:
 * * 1、 会触发重新设置全局语言的逻辑—— `locale.value = languageDefault.value;`, 语言被重设。
 * * 2、 会触发sqlite的存储逻辑, 将用户的手动设置持久化保存至sqlite的数据库文件中。
 */
const setting_store = useSettingStore();
</script>

<style lang="scss" scoped></style>
