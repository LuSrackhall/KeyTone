<template>
  <!-- TIPS: quasar的item组件中的target会造成路由过程中, 不保留历史记录, 他是仅与href这种<a>标签属性配合使用的, 起效优先与路由, 请不要与路由一起用 -->
  <!-- <q-item clickable tag="a" target="_blank" to="props.to"> -->
  <q-item clickable tag="a" :to="props.to" :class="['text-black']">
    <q-item-section v-if="props.icon" avatar>
      <q-icon :name="props.icon" />
    </q-item-section>

    <q-item-section>
      <q-item-label>{{ $t(props.title) }}</q-item-label>
      <q-item-label caption v-if="props.caption !== ''">{{ $t(props.caption) }}</q-item-label>
    </q-item-section>
  </q-item>
</template>

<script setup lang="ts">
export interface EssentialLinkProps {
  title: string;
  caption?: string;
  to?: string;
  icon?: string;
}

const props = withDefaults(defineProps<EssentialLinkProps>(), {
  caption: '',
  link: '#',
  icon: '',
});
</script>
