import { computed, Ref } from 'vue';

export interface DependencyIssue {
  type: 'sound' | 'keySound';
  id: string;
  name: string;
  missingDependencies: {
    type: 'audioFile' | 'sound' | 'keySound';
    id: string;
    name?: string;
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

export function useDependencyValidation(
  soundFileList: Ref<AudioFile[]>,
  soundList: Ref<Sound[]>,
  keySoundList: Ref<KeySound[]>
) {
  /**
   * 检查声音片段是否引用了已删除的音频文件
   */
  const validateSoundDependencies = (sound: Sound): DependencyIssue | null => {
    const missingDependencies: DependencyIssue['missingDependencies'] = [];
    
    // 检查声音片段引用的音频文件是否存在
    const referencedAudioFile = soundFileList.value.find(
      (audioFile) =>
        audioFile.sha256 === sound.soundValue.source_file_for_sound.sha256 &&
        audioFile.name_id === sound.soundValue.source_file_for_sound.name_id
    );

    if (!referencedAudioFile) {
      missingDependencies.push({
        type: 'audioFile',
        id: `${sound.soundValue.source_file_for_sound.sha256}-${sound.soundValue.source_file_for_sound.name_id}`,
        name: sound.soundValue.source_file_for_sound.name || '未知音频文件',
      });
    }

    if (missingDependencies.length > 0) {
      return {
        type: 'sound',
        id: sound.soundKey,
        name: sound.soundValue.name || '未命名声音',
        missingDependencies,
      };
    }

    return null;
  };

  /**
   * 检查至臻键音是否引用了已删除的依赖
   */
  const validateKeySoundDependencies = (keySound: KeySound): DependencyIssue | null => {
    const missingDependencies: DependencyIssue['missingDependencies'] = [];

    // 检查按下和抬起声音的依赖
    const checkDependencies = (dependencies: Array<{ type: string; value: any }>) => {
      dependencies.forEach((dep) => {
        if (dep.type === 'audio_files') {
          const audioFile = soundFileList.value.find(
            (af) => af.sha256 === dep.value?.sha256 && af.name_id === dep.value?.name_id
          );
          if (!audioFile) {
            missingDependencies.push({
              type: 'audioFile',
              id: `${dep.value?.sha256}-${dep.value?.name_id}`,
              name: dep.value?.name || '未知音频文件',
            });
          }
        } else if (dep.type === 'sounds') {
          const sound = soundList.value.find((s) => s.soundKey === dep.value?.soundKey);
          if (!sound) {
            missingDependencies.push({
              type: 'sound',
              id: dep.value?.soundKey || 'unknown',
              name: dep.value?.name || '未知声音片段',
            });
          }
        } else if (dep.type === 'key_sounds') {
          const ks = keySoundList.value.find((ks) => ks.keySoundKey === dep.value?.keySoundKey);
          if (!ks) {
            missingDependencies.push({
              type: 'keySound',
              id: dep.value?.keySoundKey || 'unknown',
              name: dep.value?.name || '未知至臻键音',
            });
          }
        }
      });
    };

    // 检查按下声音的依赖
    if (keySound.keySoundValue.down.value) {
      checkDependencies(keySound.keySoundValue.down.value);
    }

    // 检查抬起声音的依赖
    if (keySound.keySoundValue.up.value) {
      checkDependencies(keySound.keySoundValue.up.value);
    }

    if (missingDependencies.length > 0) {
      return {
        type: 'keySound',
        id: keySound.keySoundKey,
        name: keySound.keySoundValue.name || '未命名至臻键音',
        missingDependencies,
      };
    }

    return null;
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
   * 获取特定项目的依赖问题详情
   */
  const getDependencyIssue = (type: 'sound' | 'keySound', id: string): DependencyIssue | undefined => {
    return dependencyIssues.value.find((issue) => issue.type === type && issue.id === id);
  };

  return {
    dependencyIssues,
    hasDependencyIssues,
    dependencyIssuesCount,
    hasSoundDependencyIssue,
    hasKeySoundDependencyIssue,
    getDependencyIssue,
  };
}