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

            // TIPS: 用于录制按键的相关的逻辑在此处执行
            if (keyStateCallback_Record) {
              keyStateCallback_Record(keycode, keyDownStateName_ui.value);
            }

            // 使用keyDownStateName_ui, 来记录按键名称
            if (keyDownStateName_ui.value) {
              console.info('[info]:keycode为', keycode, '的按键从down变为up。其名称可能是', keyDownStateName_ui.value);

              // TIPS: 用于记录的数据持久化相关的逻辑在此处执行
              if (keyStateCallback_PersistentData) {
                keyStateCallback_PersistentData(keycode, keyDownStateName_ui.value);
              }

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
        // q.notify({
        //   message: '按键名称识别时, 检测到有多个按键被同时按下。',
        //   color: 'red',
        // });
        // q.notify({
        //   message: '为保证识别准确性, 请确保仅单个按键从按下至抬起。',
        //   color: 'red',
        // });

        // 在录制逻辑不依赖记录逻辑后, 此处通知将不再重要, 因此改为简单的告警提示。
        // * 不移除的原因是, 在项目未成熟的现阶段, 记录逻辑仍可在一定程度上作为更新前保底的最后防线(虽然很小,但还是有点用的)。
        q.notify({
          type: 'warning',
          message: '按键录制期间, 尽量避免同时按下多个按键',
        });
        return;
      }
    }
  });

  const keyDownStateName_ui = ref<string>('');

  // TIPS: 作为是否要启用'按键Dik码与name实时映射 与 持久化记录功能'的开关。(当某些ui组件需要时, 主动的开启它即可)
  const isOpeningTheRecord = ref<boolean>(false);

  /////////////////////////////////////////////////////////////////////////////////////////////////////////////
  // 声明可配置的 '录制' 用的回调函数
  // - (或者说用于录制按键的回调函数)
  let keyStateCallback_Record: ((keycode: number, keyName: string) => void) | null = null;

  /**
   * 此方法的作用是给回调用的函数变量'keyStateCallback_Record', 设置真实逻辑, 以完成其定义。
   * @param callback 回调函数(即要设置执行的真实逻辑)
   */
  function setKeyStateCallback_Record(callback: (keycode: number, keyName: string) => void) {
    isOpeningTheRecord.value = true;
    keyStateCallback_Record = callback;
  }

  /**
   * 此方法的作用是给回调用的函数变量'keyStateCallback_Record', 做清除处理, 使得其不再执行。
   */
  function clearKeyStateCallback_Record() {
    isOpeningTheRecord.value = false;
    keyStateCallback_Record = null;
  }

  /////////////////////////////////////////////////////////////////////////////////////////////////////////////

  // 定义可配置的 '记录' 用的回调函数
  // - (或者说持久化数据用的回调函数)
  let keyStateCallback_PersistentData: ((keycode: number, keyName: string) => void) | null = null;

  /**
   * 此方法的作用是给回调用的函数变量'keyStateCallback_PersistentData', 设置真实逻辑, 以完成其定义。
   * @param callback 回调函数(即要设置执行的真实逻辑)
   */
  function setKeyStateCallback_PersistentData(callback: (keycode: number, keyName: string) => void) {
    // 因为记录的逻辑, 是为了服务与录制按键逻辑的, 因此这里无需重复的开启录制逻辑。
    // isOpeningTheRecord.value = true;
    keyStateCallback_PersistentData = callback;
  }

  /**
   * 此方法的作用是为回调用的函数变量'keyStateCallback_PersistentData', 做清除处理, 使得其不再执行。
   */
  function clearKeyStateCallback_PersistentData(): void {
    // 因为记录的逻辑, 是为了服务与录制按键逻辑的, 因此这里无需重复的关闭录制逻辑。
    // isOpeningTheRecord.value = false;
    keyStateCallback_PersistentData = null;
  }

  /////////////////////////////////////////////////////////////////////////////////////////////////////////////

  // 用于记录按键的Dik码与name的映射(初版由libuiohook库中对应的定义转换而来, 后续会不断的更新, 直到映射出所有按键的名称)
  const dikCodeToName = reactive<Map<number, string>>(new Map<number, string>());

  // TODO: 以下部分不需要新的变量了, 直接复用dikCodeToName变量即可。(在初始化码表赋值后, 更改或添加对应 Dik码与name 的映射赋值。)
  // 对于类似与Meta的按键, 在windows系统中对应名称是'Windows', 在macos系统中对应名称是'Command', 在linux系统中对应名称是'Super'
  // 因此后续考虑将某些在特定系统中名称不同的按键, 进行区分处理。(这也是符合设计图纸的, 不过是否有必要仍待定。)
  // const dikCodeToName_system_windows = reactive<Map<number, string>>(new Map<number, string>());
  // const dikCodeToName_system_macos = reactive<Map<number, string>>(new Map<number, string>());
  // const dikCodeToName_system_linux = reactive<Map<number, string>>(new Map<number, string>());
  // const dikCodeToName_system_custom = reactive<Map<number, string>>(new Map<number, string>());

  // 按照 由libuiohook库中对Dik码与name的宏定义, 并借鉴前端键盘时间的keycode和key字段, 完成的码表, 来给dikCodeToName赋值。
  Object.entries({
    '1': 'Escape',
    '59': 'F1',
    '60': 'F2',
    '61': 'F3',
    '62': 'F4',
    '63': 'F5',
    '64': 'F6',
    '65': 'F7',
    '66': 'F8',
    '67': 'F9',
    '68': 'F10',
    '87': 'F11',
    '88': 'F12',
    '41': '`',
    '2': '1',
    '3': '2',
    '4': '3',
    '5': '4',
    '6': '5',
    '7': '6',
    '8': '7',
    '9': '8',
    '10': '9',
    '11': '0',
    '12': '-',
    '13': '=',
    '14': 'Backspace',
    '15': 'Tab',
    '58': 'CapsLock',
    '30': 'a',
    '48': 'b',
    '46': 'c',
    '32': 'd',
    '18': 'e',
    '33': 'f',
    '34': 'g',
    '35': 'h',
    '23': 'i',
    '36': 'j',
    '37': 'k',
    '38': 'l',
    '50': 'm',
    '49': 'n',
    '24': 'o',
    '25': 'p',
    '16': 'q',
    '19': 'r',
    '31': 's',
    '20': 't',
    '22': 'u',
    '47': 'v',
    '17': 'w',
    '45': 'x',
    '21': 'y',
    '44': 'z',
    '26': '[',
    '27': ']',
    '43': '\\',
    '39': ';',
    '40': "'",
    '28': 'Enter',
    '51': ',',
    '52': '.',
    '53': '/',
    '57': 'Space',
    '3639': 'PrintScreen',
    '70': 'ScrollLock',
    '3653': 'Pause',
    '3666': 'NumpadInsert',
    '3667': 'NumpadDelete',
    '3655': 'NumpadHome',
    '3663': 'NumpadEnd',
    '3657': 'NumpadPageUp',
    '3665': 'NumpadPageDown',
    '57416': 'NumpadUp',
    '57419': 'NumpadLeft',
    '57420': 'NumpadClear',
    '57421': 'NumpadRight',
    '57424': 'NumpadDown',
    '69': 'NumLock',
    '3637': 'NumpadDivide',
    '55': 'NumpadMultiply',
    '74': 'NumpadSubtract',
    '3597': 'NumpadEqual',
    '78': 'NumpadAdd',
    '3612': 'NumpadEnter',
    '83': 'NumpadDecimal',
    '79': 'Numpad1',
    '80': 'Numpad2',
    '81': 'Numpad3',
    '75': 'Numpad4',
    '76': 'Numpad5',
    '77': 'Numpad6',
    '71': 'Numpad7',
    '72': 'Numpad8',
    '73': 'Numpad9',
    '82': 'Numpad0',
    '60927': 'NumpadEnd',
    '60928': 'NumpadArrowDown',
    '60929': 'NumpadPageDown',
    '60931': 'NumpadArrowLeft',
    '60932': 'NumpadClear',
    '60933': 'NumpadArrowRight',
    '60935': 'NumpadHome',
    '60936': 'NumpadArrowUp',
    '60937': 'NumpadPageUp',
    '60930': 'NumpadInsert',
    '60947': 'NumpadDelete',
    '42': 'ShiftLeft',
    '54': 'ShiftRight',
    '29': 'ControlLeft',
    '3613': 'ControlRight',
    '56': 'AltLeft',
    '3640': 'AltRight',
    '3675': 'MetaLeft',
    '3676': 'MetaRight',
    '3677': 'ContextMenu',
    '57438': 'Power',
    '57439': 'Sleep',
    '57443': 'Wake',
    '57378': 'MediaPlay',
    '57380': 'MediaStop',
    '57360': 'MediaPrevious',
    '57369': 'MediaNext',
    '57453': 'MediaSelect',
    '57388': 'MediaEject',
    '57376': 'AudioVolumeMute',
    '57392': 'AudioVolumeUp',
    '57390': 'AudioVolumeDown',
    '57452': 'LaunchMail',
    '57377': 'LaunchCalculator',
    '57404': 'LaunchMusic',
    '57444': 'LaunchPictures',
    '57445': 'BrowserSearch',
    '57394': 'BrowserHome',
    '57450': 'BrowserBack',
    '57449': 'BrowserForward',
    '57448': 'BrowserStop',
    '57447': 'BrowserRefresh',
    '57446': 'BrowserFavorites',
    '112': 'Katakana',
    '115': 'Underscore',
    '119': 'Furigana',
    '121': 'Kanji',
    '123': 'Hiragana',
    '125': 'Yen',
    '126': 'NumpadComma',
    '65397': 'Help',
    '65400': 'Stop',
    '65398': 'Props',
    '65399': 'Front',
    '65396': 'Open',
    '65406': 'Find',
    '65401': 'Again',
    '65402': 'Undo',
    '65404': 'Copy',
    '65405': 'Insert',
    '65403': 'Cut',
    '0': 'Undefined',
    '65535': 'CharUndefined',
  }).forEach(([dikCode, name]) => {
    dikCodeToName.set(Number(dikCode), name as string);
  });

  // custom  TIPS: 用于存储在/不在dik原始码表内, 不确定平台的 dik码与name 的映射。(暂时不确定平台的也放于此处)
  Object.entries({
    //   ↓    - completed(已完成)   : #47的bug在解决的同时, 就无需在处理#50的bug了, 因为temp值将会被覆盖, 即使去解决它也是毫无意义的(或者说在 记录 与 录制 逻辑分开后, [#50]这个问题就不应该存在)。
    // FIXME: [#47](https://github.com/LuSrackhall/KeyTone/issues/47#:~:text=%E5%AF%B9%E4%B8%8A%E8%BF%B0json,%E5%8F%AA%E6%9C%89windows%E6%9C%BA%E5%99%A8)
    // FIXME: [#50](https://github.com/LuSrackhall/KeyTone/issues/50#:~:text=%E7%AE%AD%E5%A4%B4%E5%8C%BA%E5%9F%9F%E4%B8%8D%E8%83%BD,%E4%B8%8D%E5%BA%94%E8%AF%A5%E5%AD%98%E5%9C%A8)
    /**
     *
     * 虽然像 Home, End 之类的原始码表中也有, 但我似乎无法正确的具体区分,
     * 目前仅根据我自己发现的, 力所能及的修改这些码表(包括 原始码表、自定义码表、和 各平台码表)。
     * 因此, 因个人能力问题, 本项目暂时无法完全避免name重复现象, 只能在不断更新中完善dik码所对应的name名称, 以做到最终完全的标准化。
     * 但目前, 个人手中仅拥有win设备, 因此在拥有其它设备前, 仅会针对win平台做适配。
     * 本项目虽然是开源的, 但目前看来除了我本人外, 应该不会有其余贡献者了, 因此其它平台的适配需等到作者我本人获取对应设备之后再进行了, 可能是遥远的将来。
     * 当然, 如果您是用户, 又恰好有能力pr的话, 也欢迎您作出贡献, 但我需要对应平台的实际运行时按键按下和识别的对应视频才会合并到分支, 因为我作为唯一负责人, 要保证其映射的准确性。
     */
    60999: 'Home',
    61001: 'PageUp',

    61011: 'Delete',
    61007: 'End',
    61009: 'PageDown',

    61000: 'Up',
    61008: 'Down',
    61003: 'Left',
    61005: 'Right',
  }).forEach(([dikCode, name]) => {
    dikCodeToName.set(Number(dikCode), name as string);
  });

  // win     TIPS: 用于存储在/不在dik原始码表内, 需要对win平台特定name名称的 dik码与name 的映射。
  if (q.platform.is.win) {
    Object.entries({
      3675: 'WinLeft',
      3676: 'WinRight',
    }).forEach(([dikCode, name]) => {
      dikCodeToName.set(Number(dikCode), name as string);
    });
  }

  // mac   TIPS: 用于存储在/不在dik原始码表内, 需要对mac平台特定name名称的 dik码与name 的映射。
  if (q.platform.is.mac) {
    Object.entries({
      3675: 'CommandLeft',
      3676: 'CommandRight',
    }).forEach(([dikCode, name]) => {
      dikCodeToName.set(Number(dikCode), name as string);
    });
  }

  // linux   TIPS: 用于存储在/不在dik原始码表内, 需要对linux平台特定name名称的 dik码与name 的映射。
  if (q.platform.is.linux) {
    Object.entries({
      3675: 'SuperLeft',
      3676: 'SuperRight',
    }).forEach(([dikCode, name]) => {
      dikCodeToName.set(Number(dikCode), name as string);
    });
  }

  return {
    keyCodeState,
    frontendKeyEventStateBool,
    keyDownStateName_ui,
    isOpeningTheRecord,
    setKeyStateCallback_Record,
    clearKeyStateCallback_Record,
    setKeyStateCallback_PersistentData,
    clearKeyStateCallback_PersistentData,
    dikCodeToName,
  };
});
