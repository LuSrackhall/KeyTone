import { defineStore } from 'pinia';
import { useQuasar } from 'quasar';
import { reactive, ref, watch } from 'vue';
export const useKeyEventStore = defineStore('keyEvent', () => {
  const q = useQuasar();

  // const keyCodeState = ref<Array<{ keycode: number; state: string }>>([]);
  const keyCodeState = reactive<Map<number, string>>(new Map<number, string>());

  let previous_keyCodeState: Map<number, string>;

  // TIPS: 这里我们不使用newVal和oldVal, 因为涉及到引用对象时, 它们是相等的(vue官网中有相关说明)。
  watch(keyCodeState, (newVal, oldVal) => {
    // 触发手动开启记录, 否则不会执行这些涉及记录的逻辑
    if (isOpeningTheRecord.value) {
      //// console.log('newVal为', newVal); // 不使用newVal和oldVal, 因为涉及到引用对象时, 它们是相等的(vue官网中有相关说明)。
      //// console.log('oldVal为', oldVal); // 不使用newVal和oldVal, 因为涉及到引用对象时, 它们是相等的(vue官网中有相关说明)。

      // 使用keyCodeState, 和自己定义的previous_keyCodeState, 来代替newVal和oldVal。
      // console.log('新值为', keyCodeState);
      // console.log('旧值为', previous_keyCodeState);

      let count_down = 0;
      // 遍历新值, 若发现有多个按键被同时按下, 则触发提示并退出。
      keyCodeState.forEach((state, keycode) => {
        if (state === 'down') {
          count_down++;
        }
      });
      if (count_down > 1) {
        // 清空keyDownStateName_ui, 防止 dik 与 name 的错误映射。
        keyDownStateName_ui.value = '';
        // 不希望这通知, 在应用未获取焦点时被意外触发, 因此将此部分移动至前端键盘事件的相关逻辑中去。 但 清空keyDownStateName_ui 的逻辑仍需保留, 以保证映射的准确性。
        // q.notify({
        //   message: '按键名称识别时, 检测到有多个按键被同时按下。',
        //   color: 'red',
        // });
        // q.notify({
        //   message: '为保证识别准确性, 请确保仅单个按键从按下至抬起。',
        //   color: 'red',
        // });
        return;
      }

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
            console.debug('[debug]: keycode为', keycode, '的按键从down变为up。其名称可能是', keyDownStateName_ui.value);

            // 使用keyDownStateName_ui, 来记录按键名称

            if (keyDownStateName_ui.value) {
              console.info('[info]:keycode为', keycode, '的按键从down变为up。其名称可能是', keyDownStateName_ui.value);
              // TODO: 用于记录的逻辑写在此处
              // ...

              // 记录行为结束后, 清空keyDownStateName_ui
              keyDownStateName_ui.value = '';
            }
          }
        }
      });
    }
    previous_keyCodeState = new Map(keyCodeState);
  });

  //////////////////////////////////frontend----ui//////////////////////////////////////////////

  // 由于需要时间差, 因此对应ui的按键状态, 不需要监听按下并抬起, 而是仅监听按下, 并记录按键名称。
  // const keyCodeState_ui = reactive<Map<string, string>>(new Map<string, string>());

  // 用于记录前端按键状态的bool值, 主要用于防止持续按下的按键被重复触发。
  // * 也可以在其它逻辑中, 当作keyCodeState_ui来的使用。(使用过程中请注意确保只读式使用, 以免破坏它的主要用途)
  //   > 只读使用的过程中, 当其为true时代表此按键为抬起状态, 为false时代表按键为按下状态。
  const frontendKeyEventStateBool = reactive<Map<string, boolean>>(new Map<string, boolean>());

  watch(frontendKeyEventStateBool, (newVal, oldVal) => {
    if (isOpeningTheRecord.value) {
      let count_down = 0;
      // 遍历新值, 若发现有多个按键被同时按下, 则触发提示并退出。
      frontendKeyEventStateBool.forEach((state, keycode) => {
        if (state === false) {
          count_down++;
        }
      });
      if (count_down > 1) {
        // 不希望这通知, 在应用未获取焦点时触发, 因此将此部分移动至此, 此处仅作为通知触发, 不参与 按键Dik码 与 name 的实时映射逻辑。
        q.notify({
          message: '按键名称识别时, 检测到有多个按键被同时按下。',
          color: 'red',
        });
        q.notify({
          message: '为保证识别准确性, 请确保仅单个按键从按下至抬起。',
          color: 'red',
        });
        return;
      }
    }
  });

  const keyDownStateName_ui = ref<string>('');

  // TIPS: 作为是否要启用'按键Dik码与name实时映射 与 持久化记录功能'的开关。(当某些ui组件需要时, 主动的开启它即可)
  const isOpeningTheRecord = ref<boolean>(false);

  return { keyCodeState, frontendKeyEventStateBool, keyDownStateName_ui, isOpeningTheRecord };
});
