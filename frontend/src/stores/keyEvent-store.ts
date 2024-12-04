import { defineStore } from 'pinia';
import { ref } from 'vue';
export const useKeyEventStore = defineStore('keyEvent', () => {
  // const keyCodeState = ref<Array<{ keycode: number; state: string }>>([]);
  const keyCodeState = ref<Map<number, string>>(new Map<number, string>());

  return { keyCodeState };
});
