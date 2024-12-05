import { defineStore } from 'pinia';
import { reactive, ref, watch } from 'vue';
export const useKeyEventStore = defineStore('keyEvent', () => {
  // const keyCodeState = ref<Array<{ keycode: number; state: string }>>([]);
  const keyCodeState = reactive<Map<number, string>>(new Map<number, string>());

  let previous_keyCodeState: Map<number, string>;

  // TIPS: 这里我们不使用newVal和oldVal, 因为涉及到引用对象时, 它们是相等的(vue官网中有相关说明)。
  watch(keyCodeState, (newVal, oldVal) => {
    //// console.log('newVal为', newVal); // 不使用newVal和oldVal, 因为涉及到引用对象时, 它们是相等的(vue官网中有相关说明)。
    //// console.log('oldVal为', oldVal); // 不使用newVal和oldVal, 因为涉及到引用对象时, 它们是相等的(vue官网中有相关说明)。

    // 使用keyCodeState, 和自己定义的previous_keyCodeState, 来代替newVal和oldVal。
    // console.log('新值为', keyCodeState);
    // console.log('旧值为', previous_keyCodeState);

    // 遍历新值, 以找寻按键从down变为up的按键, 并执行相关逻辑。
    keyCodeState.forEach((state, keycode) => {
      // 当新值为up时
      if (state === 'up') {
        // console.log('keycode为', keycode, '的按键当前状态为', state);
        // 检查旧值中对应按键是否为down
        // console.log('旧值为', oldVal);

        if (previous_keyCodeState.get(keycode) === 'down') {
          // TODO: 在此处添加按键从down变为up时的处理逻辑
          // 例如:
          // - 触发按键抬起事件
          // - 记录按键持续时间
          // - 执行相关动画
          // - 播放音效等
          console.log('keycode为', keycode, '的按键从down变为up');
        }
      }
    });

    previous_keyCodeState = new Map(keyCodeState);
  });

  const keyCodeState_ui = reactive<Map<string, string>>(new Map<string, string>());

  return { keyCodeState };
});
