<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="min-width: 350px">
      <q-card-section class="sticky top-0 z-10 bg-white/30 backdrop-blur-sm">
        <div class="text-h6">再次导出确认</div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        该键音专辑原始作者明确了该键音包的二次导出必须实施签名，请确认后选择要使用的签名或取消。
      </q-card-section>

      <q-card-actions align="right" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat label="取消" color="primary" @click="onCancel" />
        <q-btn flat label="确认" color="primary" @click="onConfirm" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  visible: boolean;
}>();

const emit = defineEmits<{
  (e: 'confirm'): void;
  (e: 'cancel'): void;
}>();

const isVisible = computed({
  get: () => props.visible,
  set: (val) => {
    if (!val) emit('cancel');
  },
});

const onConfirm = () => {
  emit('confirm');
};

const onCancel = () => {
  emit('cancel');
};
</script>
