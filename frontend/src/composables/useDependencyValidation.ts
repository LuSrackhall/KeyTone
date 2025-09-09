import { computed, Ref } from 'vue';

export interface DependencyIssue {
  type: 'sound' | 'keySound' | 'globalBinding' | 'singleKeyBinding';
  id: string;
  name: string;
  severity: 'error' | 'warning'; // 'error' for direct deletions, 'warning' for indirect issues
  missingDependencies: {
    type: 'audioFile' | 'sound' | 'keySound';
    id: string;
    name?: string;
    direct: boolean; // true if directly missing, false if indirectly problematic
  }[];
}

export interface AudioFile {
  sha256: string;
  name_id: string;
  name: string;
  type: string;
}

export interface Sound {
  soundKey: string;
  soundValue: {
    name: string;
    source_file_for_sound: {
      sha256: string;
      name_id: string;
      type: string;
    };
    cut: {
      start_time: number;
      end_time: number;
      volume: number;
    };
  };
}

export interface KeySound {
  keySoundKey: string;
  keySoundValue: {
    name: string;
    down: {
      mode: string;
      value: Array<{
        type: 'audio_files' | 'sounds' | 'key_sounds';
        value: any;
      }>;
    };
    up: {
      mode: string;
      value: Array<{
        type: 'audio_files' | 'sounds' | 'key_sounds';
        value: any;
      }>;
    };
  };
}

export interface GlobalBinding {
  down: any;
  up: any;
}

export interface SingleKeyBinding {
  key: string;
  binding: {
    down: any;
    up: any;
  };
}

export function useDependencyValidation(
  soundFileList: Ref<AudioFile[]>,
  soundList: Ref<Sound[]>,
  keySoundList: Ref<KeySound[]>,
  globalBinding?: Ref<GlobalBinding | undefined>,
  singleKeyBindings?: Ref<Map<string, any>>
) {
  /**
   * 检查音频文件是否存在
   */
  const isAudioFileExists = (sha256: string, nameId: string): boolean => {
    return soundFileList.value.some(
      (audioFile) => audioFile.sha256 === sha256 && audioFile.name_id === nameId
    );
  };

  /**
   * 检查声音是否存在
   */
  const isSoundExists = (soundKey: string): boolean => {
    return soundList.value.some((sound) => sound.soundKey === soundKey);
  };

  /**
   * 检查键音是否存在
   */
  const isKeySoundExists = (keySoundKey: string): boolean => {
    return keySoundList.value.some((keySound) => keySound.keySoundKey === keySoundKey);
  };

  /**
   * 检查声音片段的依赖问题（直接和间接）
   */
  const validateSoundDependencies = (sound: Sound): DependencyIssue | null => {
    const missingDependencies: DependencyIssue['missingDependencies'] = [];
    let severity: 'error' | 'warning' = 'error';
    
    // 检查声音片段直接引用的音频文件是否存在
    if (!isAudioFileExists(
      sound.soundValue.source_file_for_sound.sha256,
      sound.soundValue.source_file_for_sound.name_id
    )) {
      missingDependencies.push({
        type: 'audioFile',
        id: `${sound.soundValue.source_file_for_sound.sha256}-${sound.soundValue.source_file_for_sound.name_id}`,
        name: '音频源文件',
        direct: true,
      });
    }

    if (missingDependencies.length > 0) {
      return {
        type: 'sound',
        id: sound.soundKey,
        name: sound.soundValue.name || '未命名声音',
        severity,
        missingDependencies,
      };
    }

    return null;
  };

  /**
   * 检查至臻键音的依赖问题（直接和间接）
   */
  const validateKeySoundDependencies = (keySound: KeySound): DependencyIssue | null => {
    const missingDependencies: DependencyIssue['missingDependencies'] = [];
    let hasDirectIssues = false;

    // 检查依赖项的函数
    const checkDependencies = (dependencies: Array<{ type: string; value: any }>) => {
      dependencies.forEach((dep) => {
        if (dep.type === 'audio_files') {
          if (!isAudioFileExists(dep.value?.sha256, dep.value?.name_id)) {
            missingDependencies.push({
              type: 'audioFile',
              id: `${dep.value?.sha256}-${dep.value?.name_id}`,
              name: '音频源文件',
              direct: true,
            });
            hasDirectIssues = true;
          }
        } else if (dep.type === 'sounds') {
          const soundKey = typeof dep.value === 'string' ? dep.value : dep.value?.soundKey;
          if (!isSoundExists(soundKey)) {
            missingDependencies.push({
              type: 'sound',
              id: soundKey || 'unknown',
              name: '声音片段',
              direct: true,
            });
            hasDirectIssues = true;
          } else {
            // 检查间接依赖问题 - 声音存在但它引用的音频文件被删除
            const sound = soundList.value.find((s) => s.soundKey === soundKey);
            if (sound && !isAudioFileExists(
              sound.soundValue.source_file_for_sound.sha256,
              sound.soundValue.source_file_for_sound.name_id
            )) {
              missingDependencies.push({
                type: 'sound',
                id: soundKey,
                name: '声音片段（间接依赖问题）',
                direct: false,
              });
            }
          }
        } else if (dep.type === 'key_sounds') {
          const keySoundKey = typeof dep.value === 'string' ? dep.value : dep.value?.keySoundKey;
          if (!isKeySoundExists(keySoundKey)) {
            missingDependencies.push({
              type: 'keySound',
              id: keySoundKey || 'unknown',
              name: '至臻键音',
              direct: true,
            });
            hasDirectIssues = true;
          } else {
            // 检查间接依赖问题 - 键音存在但它的依赖有问题
            const referencedKeySound = keySoundList.value.find((ks) => ks.keySoundKey === keySoundKey);
            if (referencedKeySound) {
              const indirectIssue = validateKeySoundDependencies(referencedKeySound);
              if (indirectIssue) {
                missingDependencies.push({
                  type: 'keySound',
                  id: keySoundKey,
                  name: '至臻键音（间接依赖问题）',
                  direct: false,
                });
              }
            }
          }
        }
      });
    };

    // 检查按下和抬起声音的依赖
    if (keySound.keySoundValue.down.value) {
      checkDependencies(keySound.keySoundValue.down.value);
    }
    if (keySound.keySoundValue.up.value) {
      checkDependencies(keySound.keySoundValue.up.value);
    }

    if (missingDependencies.length > 0) {
      return {
        type: 'keySound',
        id: keySound.keySoundKey,
        name: keySound.keySoundValue.name || '未命名至臻键音',
        severity: hasDirectIssues ? 'error' : 'warning',
        missingDependencies,
      };
    }

    return null;
  };

  /**
   * 检查单个绑定项的依赖问题
   */
  const validateBindingDependency = (binding: any): Array<DependencyIssue['missingDependencies'][0]> => {
    const missingDependencies: Array<DependencyIssue['missingDependencies'][0]> = [];

    if (!binding || !binding.type) {
      return missingDependencies;
    }

    if (binding.type === 'audio_files') {
      if (!isAudioFileExists(binding.value?.sha256, binding.value?.name_id)) {
        missingDependencies.push({
          type: 'audioFile',
          id: `${binding.value?.sha256}-${binding.value?.name_id}`,
          name: '音频源文件',
          direct: true,
        });
      }
    } else if (binding.type === 'sounds') {
      const soundKey = typeof binding.value === 'string' ? binding.value : binding.value?.soundKey;
      if (!isSoundExists(soundKey)) {
        missingDependencies.push({
          type: 'sound',
          id: soundKey || 'unknown',
          name: '声音片段',
          direct: true,
        });
      } else {
        // 检查间接依赖问题
        const sound = soundList.value.find((s) => s.soundKey === soundKey);
        if (sound && !isAudioFileExists(
          sound.soundValue.source_file_for_sound.sha256,
          sound.soundValue.source_file_for_sound.name_id
        )) {
          missingDependencies.push({
            type: 'sound',
            id: soundKey,
            name: '声音片段（间接依赖问题）',
            direct: false,
          });
        }
      }
    } else if (binding.type === 'key_sounds') {
      const keySoundKey = typeof binding.value === 'string' ? binding.value : binding.value?.keySoundKey;
      if (!isKeySoundExists(keySoundKey)) {
        missingDependencies.push({
          type: 'keySound',
          id: keySoundKey || 'unknown',
          name: '至臻键音',
          direct: true,
        });
      } else {
        // 检查间接依赖问题
        const keySound = keySoundList.value.find((ks) => ks.keySoundKey === keySoundKey);
        if (keySound) {
          const indirectIssue = validateKeySoundDependencies(keySound);
          if (indirectIssue) {
            missingDependencies.push({
              type: 'keySound',
              id: keySoundKey,
              name: '至臻键音（间接依赖问题）',
              direct: false,
            });
          }
        }
      }
    }

    return missingDependencies;
  };

  /**
   * 检查全局绑定的依赖问题
   */
  const validateGlobalBindingDependencies = (): DependencyIssue | null => {
    if (!globalBinding?.value) {
      return null;
    }

    const missingDependencies: DependencyIssue['missingDependencies'] = [];
    let hasDirectIssues = false;

    // 检查按下绑定
    if (globalBinding.value.down) {
      const downIssues = validateBindingDependency(globalBinding.value.down);
      missingDependencies.push(...downIssues);
      if (downIssues.some(issue => issue.direct)) {
        hasDirectIssues = true;
      }
    }

    // 检查抬起绑定
    if (globalBinding.value.up) {
      const upIssues = validateBindingDependency(globalBinding.value.up);
      missingDependencies.push(...upIssues);
      if (upIssues.some(issue => issue.direct)) {
        hasDirectIssues = true;
      }
    }

    if (missingDependencies.length > 0) {
      return {
        type: 'globalBinding',
        id: 'global',
        name: '全局绑定',
        severity: hasDirectIssues ? 'error' : 'warning',
        missingDependencies,
      };
    }

    return null;
  };

  /**
   * 检查单键绑定的依赖问题
   */
  const validateSingleKeyBindingDependencies = (): DependencyIssue[] => {
    const issues: DependencyIssue[] = [];

    if (!singleKeyBindings?.value) {
      return issues;
    }

    for (const [key, binding] of singleKeyBindings.value.entries()) {
      const missingDependencies: DependencyIssue['missingDependencies'] = [];
      let hasDirectIssues = false;

      // 检查按下绑定
      if (binding.down) {
        const downIssues = validateBindingDependency(binding.down);
        missingDependencies.push(...downIssues);
        if (downIssues.some((issue: any) => issue.direct)) {
          hasDirectIssues = true;
        }
      }

      // 检查抬起绑定
      if (binding.up) {
        const upIssues = validateBindingDependency(binding.up);
        missingDependencies.push(...upIssues);
        if (upIssues.some((issue: any) => issue.direct)) {
          hasDirectIssues = true;
        }
      }

      if (missingDependencies.length > 0) {
        issues.push({
          type: 'singleKeyBinding',
          id: key,
          name: `单键绑定 ${key}`,
          severity: hasDirectIssues ? 'error' : 'warning',
          missingDependencies,
        });
      }
    }

    return issues;
  };

  /**
   * 获取所有依赖问题
   */
  const dependencyIssues = computed<DependencyIssue[]>(() => {
    const issues: DependencyIssue[] = [];

    // 检查声音片段的依赖问题
    soundList.value.forEach((sound) => {
      const issue = validateSoundDependencies(sound);
      if (issue) {
        issues.push(issue);
      }
    });

    // 检查至臻键音的依赖问题
    keySoundList.value.forEach((keySound) => {
      const issue = validateKeySoundDependencies(keySound);
      if (issue) {
        issues.push(issue);
      }
    });

    // 检查全局绑定的依赖问题
    const globalIssue = validateGlobalBindingDependencies();
    if (globalIssue) {
      issues.push(globalIssue);
    }

    // 检查单键绑定的依赖问题
    const singleKeyIssues = validateSingleKeyBindingDependencies();
    issues.push(...singleKeyIssues);

    return issues;
  });

  /**
   * 是否有依赖问题
   */
  const hasDependencyIssues = computed(() => dependencyIssues.value.length > 0);

  /**
   * 依赖问题总数
   */
  const dependencyIssuesCount = computed(() => dependencyIssues.value.length);

  /**
   * 错误级别问题总数（直接删除）
   */
  const errorIssuesCount = computed(() => 
    dependencyIssues.value.filter(issue => issue.severity === 'error').length
  );

  /**
   * 警告级别问题总数（间接问题）
   */
  const warningIssuesCount = computed(() => 
    dependencyIssues.value.filter(issue => issue.severity === 'warning').length
  );

  /**
   * 检查特定声音是否有依赖问题
   */
  const hasSoundDependencyIssue = (soundKey: string): boolean => {
    return dependencyIssues.value.some((issue) => issue.type === 'sound' && issue.id === soundKey);
  };

  /**
   * 检查特定至臻键音是否有依赖问题
   */
  const hasKeySoundDependencyIssue = (keySoundKey: string): boolean => {
    return dependencyIssues.value.some((issue) => issue.type === 'keySound' && issue.id === keySoundKey);
  };

  /**
   * 检查全局绑定是否有依赖问题
   */
  const hasGlobalBindingDependencyIssue = (): boolean => {
    return dependencyIssues.value.some((issue) => issue.type === 'globalBinding');
  };

  /**
   * 检查特定单键绑定是否有依赖问题
   */
  const hasSingleKeyBindingDependencyIssue = (key: string): boolean => {
    return dependencyIssues.value.some((issue) => issue.type === 'singleKeyBinding' && issue.id === key);
  };

  /**
   * 获取特定项目的依赖问题详情
   */
  const getDependencyIssue = (type: DependencyIssue['type'], id: string): DependencyIssue | undefined => {
    return dependencyIssues.value.find((issue) => issue.type === type && issue.id === id);
  };

  return {
    dependencyIssues,
    hasDependencyIssues,
    dependencyIssuesCount,
    errorIssuesCount,
    warningIssuesCount,
    hasSoundDependencyIssue,
    hasKeySoundDependencyIssue,
    hasGlobalBindingDependencyIssue,
    hasSingleKeyBindingDependencyIssue,
    getDependencyIssue,
  };
}