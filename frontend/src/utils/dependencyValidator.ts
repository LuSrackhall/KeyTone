/**
 * Dependency validation utility for KeyTone album creation
 * 
 * This module provides functions to validate dependencies between different components:
 * - Audio source files (音频源文件)
 * - Trimmed/defined sounds (裁剪定义的声音) 
 * - Key sounds (键音/至臻键音/高级键音)
 * - Global bindings (全局绑定)
 * - Single key bindings (单键绑定)
 */

export interface AudioFile {
  sha256: string;
  name_id: string;
  name: string;
  type: string;
}

export interface Sound {
  soundKey: string;
  soundValue: {
    source_file_for_sound: {
      sha256: string;
      name_id: string;
      type: string;
    };
    name: string;
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

export interface DependencyIssue {
  type: 'direct' | 'indirect';
  severity: 'critical' | 'info';
  itemType: 'audio_files' | 'sounds' | 'key_sounds' | 'global_binding' | 'single_key_binding';
  itemId: string;
  itemName: string;
  missingDependencies: Array<{
    type: 'audio_files' | 'sounds' | 'key_sounds';
    id: string;
    name?: string;
  }>;
  message: string;
}

export class DependencyValidator {
  private audioFiles: AudioFile[];
  private sounds: Sound[];
  private keySounds: KeySound[];

  constructor(audioFiles: AudioFile[], sounds: Sound[], keySounds: KeySound[]) {
    this.audioFiles = audioFiles;
    this.sounds = sounds;
    this.keySounds = keySounds;
  }

  /**
   * Check if an audio file exists by sha256 and name_id
   */
  private audioFileExists(sha256: string, name_id: string): boolean {
    return this.audioFiles.some(file => 
      file.sha256 === sha256 && file.name_id === name_id
    );
  }

  /**
   * Check if a sound exists by soundKey
   */
  private soundExists(soundKey: string): boolean {
    return this.sounds.some(sound => sound.soundKey === soundKey);
  }

  /**
   * Check if a key sound exists by keySoundKey
   */
  private keySoundExists(keySoundKey: string): boolean {
    return this.keySounds.some(keySound => keySound.keySoundKey === keySoundKey);
  }

  /**
   * Get audio file name by sha256 and name_id
   */
  private getAudioFileName(sha256: string, name_id: string): string {
    const file = this.audioFiles.find(f => f.sha256 === sha256 && f.name_id === name_id);
    return file ? `${file.name}${file.type}` : `${sha256.substring(0, 8)}...`;
  }

  /**
   * Get sound name by soundKey
   */
  private getSoundName(soundKey: string): string {
    // Ensure soundKey is a string
    if (typeof soundKey !== 'string') {
      return 'Unknown';
    }
    
    const sound = this.sounds.find(s => s.soundKey === soundKey);
    if (sound && sound.soundValue.name) {
      return sound.soundValue.name;
    }
    if (sound) {
      // Fallback to source file name if sound name is empty
      const sourceFile = sound.soundValue.source_file_for_sound;
      return this.getAudioFileName(sourceFile.sha256, sourceFile.name_id);
    }
    return soundKey.substring(0, 8) + '...';
  }

  /**
   * Get key sound name by keySoundKey
   */
  private getKeySoundName(keySoundKey: string): string {
    // Ensure keySoundKey is a string
    if (typeof keySoundKey !== 'string') {
      return 'Unknown';
    }
    
    const keySound = this.keySounds.find(ks => ks.keySoundKey === keySoundKey);
    return keySound ? keySound.keySoundValue.name : keySoundKey.substring(0, 8) + '...';
  }

  /**
   * Validate direct dependencies for sounds (sounds depend on audio files)
   */
  validateSoundsDirectDependencies(): DependencyIssue[] {
    const issues: DependencyIssue[] = [];

    for (const sound of this.sounds) {
      const sourceFile = sound.soundValue.source_file_for_sound;
      
      if (!this.audioFileExists(sourceFile.sha256, sourceFile.name_id)) {
        issues.push({
          type: 'direct',
          severity: 'critical',
          itemType: 'sounds',
          itemId: sound.soundKey,
          itemName: this.getSoundName(sound.soundKey),
          missingDependencies: [{
            type: 'audio_files',
            id: `${sourceFile.sha256}:${sourceFile.name_id}`,
            name: this.getAudioFileName(sourceFile.sha256, sourceFile.name_id)
          }],
          message: '裁剪定义的声音所依赖的音频源文件已被删除'
        });
      }
    }

    return issues;
  }

  /**
   * Validate direct dependencies for key sounds
   */
  validateKeySoundsDirectDependencies(): DependencyIssue[] {
    const issues: DependencyIssue[] = [];

    for (const keySound of this.keySounds) {
      const missingDeps: Array<{type: 'audio_files' | 'sounds' | 'key_sounds'; id: string; name?: string}> = [];

      // Check down dependencies
      for (const dep of keySound.keySoundValue.down.value) {
        if (dep.type === 'audio_files') {
          if (!this.audioFileExists(dep.value.sha256, dep.value.name_id)) {
            missingDeps.push({
              type: 'audio_files',
              id: `${dep.value.sha256}:${dep.value.name_id}`,
              name: this.getAudioFileName(dep.value.sha256, dep.value.name_id)
            });
          }
        } else if (dep.type === 'sounds') {
          if (!this.soundExists(dep.value)) {
            missingDeps.push({
              type: 'sounds',
              id: dep.value,
              name: this.getSoundName(dep.value)
            });
          }
        } else if (dep.type === 'key_sounds') {
          if (!this.keySoundExists(dep.value)) {
            missingDeps.push({
              type: 'key_sounds',
              id: dep.value,
              name: this.getKeySoundName(dep.value)
            });
          }
        }
      }

      // Check up dependencies
      for (const dep of keySound.keySoundValue.up.value) {
        if (dep.type === 'audio_files') {
          if (!this.audioFileExists(dep.value.sha256, dep.value.name_id)) {
            missingDeps.push({
              type: 'audio_files',
              id: `${dep.value.sha256}:${dep.value.name_id}`,
              name: this.getAudioFileName(dep.value.sha256, dep.value.name_id)
            });
          }
        } else if (dep.type === 'sounds') {
          if (!this.soundExists(dep.value)) {
            missingDeps.push({
              type: 'sounds',
              id: dep.value,
              name: this.getSoundName(dep.value)
            });
          }
        } else if (dep.type === 'key_sounds') {
          if (!this.keySoundExists(dep.value)) {
            missingDeps.push({
              type: 'key_sounds',
              id: dep.value,
              name: this.getKeySoundName(dep.value)
            });
          }
        }
      }

      if (missingDeps.length > 0) {
        issues.push({
          type: 'direct',
          severity: 'critical',
          itemType: 'key_sounds',
          itemId: keySound.keySoundKey,
          itemName: keySound.keySoundValue.name,
          missingDependencies: missingDeps,
          message: '键音所直接依赖的依赖项已被删除'
        });
      }
    }

    return issues;
  }

  /**
   * Validate binding dependencies (global or single key bindings)
   */
  validateBindingDirectDependencies(binding: any, bindingType: 'global_binding' | 'single_key_binding', bindingId: string): DependencyIssue[] {
    const issues: DependencyIssue[] = [];
    const missingDeps: Array<{type: 'audio_files' | 'sounds' | 'key_sounds'; id: string; name?: string}> = [];

    // Check down binding
    if (binding.down) {
      if (binding.down.type === 'audio_files') {
        if (!this.audioFileExists(binding.down.value.sha256, binding.down.value.name_id)) {
          missingDeps.push({
            type: 'audio_files',
            id: `${binding.down.value.sha256}:${binding.down.value.name_id}`,
            name: this.getAudioFileName(binding.down.value.sha256, binding.down.value.name_id)
          });
        }
      } else if (binding.down.type === 'sounds') {
        if (!this.soundExists(binding.down.value)) {
          missingDeps.push({
            type: 'sounds',
            id: binding.down.value,
            name: this.getSoundName(binding.down.value)
          });
        }
      } else if (binding.down.type === 'key_sounds') {
        if (!this.keySoundExists(binding.down.value)) {
          missingDeps.push({
            type: 'key_sounds',
            id: binding.down.value,
            name: this.getKeySoundName(binding.down.value)
          });
        }
      }
    }

    // Check up binding
    if (binding.up) {
      if (binding.up.type === 'audio_files') {
        if (!this.audioFileExists(binding.up.value.sha256, binding.up.value.name_id)) {
          missingDeps.push({
            type: 'audio_files',
            id: `${binding.up.value.sha256}:${binding.up.value.name_id}`,
            name: this.getAudioFileName(binding.up.value.sha256, binding.up.value.name_id)
          });
        }
      } else if (binding.up.type === 'sounds') {
        if (!this.soundExists(binding.up.value)) {
          missingDeps.push({
            type: 'sounds',
            id: binding.up.value,
            name: this.getSoundName(binding.up.value)
          });
        }
      } else if (binding.up.type === 'key_sounds') {
        if (!this.keySoundExists(binding.up.value)) {
          missingDeps.push({
            type: 'key_sounds',
            id: binding.up.value,
            name: this.getKeySoundName(binding.up.value)
          });
        }
      }
    }

    if (missingDeps.length > 0) {
      issues.push({
        type: 'direct',
        severity: 'critical',
        itemType: bindingType,
        itemId: bindingId,
        itemName: bindingType === 'global_binding' ? '全局绑定' : `单键绑定(${bindingId})`,
        missingDependencies: missingDeps,
        message: bindingType === 'global_binding' ? '全局绑定所直接依赖的依赖项已被删除' : '单键绑定所直接依赖的依赖项已被删除'
      });
    }

    return issues;
  }

  /**
   * Check if a sound has indirect dependency issues (its audio file is missing)
   */
  private soundHasIndirectIssues(soundKey: string): boolean {
    const sound = this.sounds.find(s => s.soundKey === soundKey);
    if (!sound) return false;
    
    const sourceFile = sound.soundValue.source_file_for_sound;
    return !this.audioFileExists(sourceFile.sha256, sourceFile.name_id);
  }

  /**
   * Check if a key sound has indirect dependency issues
   */
  private keySoundHasIndirectIssues(keySoundKey: string): boolean {
    const keySound = this.keySounds.find(ks => ks.keySoundKey === keySoundKey);
    if (!keySound) return false;

    // Check all dependencies for indirect issues
    const allDeps = [...keySound.keySoundValue.down.value, ...keySound.keySoundValue.up.value];
    
    for (const dep of allDeps) {
      if (dep.type === 'sounds' && this.soundHasIndirectIssues(dep.value)) {
        return true;
      } else if (dep.type === 'key_sounds' && this.keySoundHasIndirectIssues(dep.value)) {
        return true;
      }
    }
    
    return false;
  }

  /**
   * Validate indirect dependencies for key sounds
   */
  validateKeySoundsIndirectDependencies(): DependencyIssue[] {
    const issues: DependencyIssue[] = [];

    for (const keySound of this.keySounds) {
      // Skip if there are direct dependency issues
      const directIssues = this.validateKeySoundsDirectDependencies()
        .filter(issue => issue.itemId === keySound.keySoundKey);
      if (directIssues.length > 0) continue;

      if (this.keySoundHasIndirectIssues(keySound.keySoundKey)) {
        issues.push({
          type: 'indirect',
          severity: 'info',
          itemType: 'key_sounds',
          itemId: keySound.keySoundKey,
          itemName: keySound.keySoundValue.name,
          missingDependencies: [], // Could be expanded to list specific indirect issues
          message: '键音的依赖链路中存在间接删除的情况'
        });
      }
    }

    return issues;
  }

  /**
   * Validate indirect dependencies for bindings
   */
  validateBindingIndirectDependencies(binding: any, bindingType: 'global_binding' | 'single_key_binding', bindingId: string): DependencyIssue[] {
    const issues: DependencyIssue[] = [];

    // Skip if there are direct dependency issues
    const directIssues = this.validateBindingDirectDependencies(binding, bindingType, bindingId);
    if (directIssues.length > 0) return issues;

    let hasIndirectIssues = false;

    // Check down binding for indirect issues
    if (binding.down) {
      if (binding.down.type === 'sounds' && this.soundHasIndirectIssues(binding.down.value)) {
        hasIndirectIssues = true;
      } else if (binding.down.type === 'key_sounds' && this.keySoundHasIndirectIssues(binding.down.value)) {
        hasIndirectIssues = true;
      }
    }

    // Check up binding for indirect issues
    if (binding.up) {
      if (binding.up.type === 'sounds' && this.soundHasIndirectIssues(binding.up.value)) {
        hasIndirectIssues = true;
      } else if (binding.up.type === 'key_sounds' && this.keySoundHasIndirectIssues(binding.up.value)) {
        hasIndirectIssues = true;
      }
    }

    if (hasIndirectIssues) {
      issues.push({
        type: 'indirect',
        severity: 'info',
        itemType: bindingType,
        itemId: bindingId,
        itemName: bindingType === 'global_binding' ? '全局绑定' : `单键绑定(${bindingId})`,
        missingDependencies: [], // Could be expanded to list specific indirect issues
        message: bindingType === 'global_binding' ? '全局绑定的依赖链路中存在间接删除的情况' : '单键绑定的依赖链路中存在间接删除的情况'
      });
    }

    return issues;
  }

  /**
   * Validate all dependencies and return all issues
   */
  validateAllDependencies(globalBinding?: any, singleKeyBindings?: Map<string, any>): DependencyIssue[] {
    const allIssues: DependencyIssue[] = [];

    // Direct dependency validation
    allIssues.push(...this.validateSoundsDirectDependencies());
    allIssues.push(...this.validateKeySoundsDirectDependencies());

    if (globalBinding) {
      allIssues.push(...this.validateBindingDirectDependencies(globalBinding, 'global_binding', 'global'));
    }

    if (singleKeyBindings) {
      singleKeyBindings.forEach((binding, keyId) => {
        allIssues.push(...this.validateBindingDirectDependencies(binding, 'single_key_binding', keyId));
      });
    }

    // Indirect dependency validation
    allIssues.push(...this.validateKeySoundsIndirectDependencies());

    if (globalBinding) {
      allIssues.push(...this.validateBindingIndirectDependencies(globalBinding, 'global_binding', 'global'));
    }

    if (singleKeyBindings) {
      singleKeyBindings.forEach((binding, keyId) => {
        allIssues.push(...this.validateBindingIndirectDependencies(binding, 'single_key_binding', keyId));
      });
    }

    return allIssues;
  }
}

/**
 * Helper function to create a dependency validator instance
 */
export function createDependencyValidator(
  audioFiles: AudioFile[],
  sounds: Sound[],
  keySounds: KeySound[]
): DependencyValidator {
  return new DependencyValidator(audioFiles, sounds, keySounds);
}

/**
 * Helper function to check if an item has dependency issues
 */
export function hasItemDependencyIssues(
  itemType: 'audio_files' | 'sounds' | 'key_sounds',
  itemId: string,
  issues: DependencyIssue[]
): { hasIssues: boolean; criticalCount: number; infoCount: number } {
  const itemIssues = issues.filter(issue => 
    issue.itemType === itemType && issue.itemId === itemId
  );

  const criticalCount = itemIssues.filter(issue => issue.severity === 'critical').length;
  const infoCount = itemIssues.filter(issue => issue.severity === 'info').length;

  return {
    hasIssues: itemIssues.length > 0,
    criticalCount,
    infoCount
  };
}