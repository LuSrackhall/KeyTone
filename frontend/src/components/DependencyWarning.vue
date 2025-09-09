<template>
  <div v-if="visible" class="dependency-warning">
    <!-- 简单警告图标 -->
    <q-icon 
      v-if="showIcon" 
      :name="getWarningIcon()" 
      :class="getWarningIconClass()"
      :size="iconSize"
    />
    
    <!-- 警告内容 -->
    <div v-if="showTooltip" class="warning-content">
      <q-tooltip 
        anchor="bottom middle" 
        self="top middle" 
        :offset="[0, 8]"
        :class="getTooltipClass()"
      >
        <div class="tooltip-content">
          <div class="tooltip-title">
            {{ getWarningTitle() }}
          </div>
          <div class="missing-items">
            <div 
              v-for="dep in issue.missingDependencies" 
              :key="dep.id" 
              class="missing-item"
            >
              <q-icon 
                :name="getDependencyIcon(dep.type)" 
                size="14px" 
                :class="['dep-icon', getDependencyIconClass(dep)]"
              />
              <span :class="getDependencyTextClass(dep)">
                {{ getDependencyTypeText(dep.type) }}: {{ dep.name || dep.id }}
                <span v-if="!dep.direct" class="indirect-indicator">（间接问题）</span>
              </span>
            </div>
          </div>
        </div>
      </q-tooltip>
    </div>
    
    <!-- 警告文本 -->
    <span v-if="showText" :class="getWarningTextClass()">
      {{ getWarningTitle() }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import type { DependencyIssue } from 'src/composables/useDependencyValidation';

interface Props {
  issue?: DependencyIssue;
  visible?: boolean;
  showIcon?: boolean;
  showText?: boolean;
  showTooltip?: boolean;
  iconSize?: string;
}

const props = withDefaults(defineProps<Props>(), {
  visible: true,
  showIcon: true,
  showText: false,
  showTooltip: true,
  iconSize: '16px',
});

const { t } = useI18n();

const getWarningIcon = (): string => {
  return props.issue?.severity === 'error' ? 'error' : 'warning';
};

const getWarningIconClass = (): string => {
  const baseClass = 'warning-icon';
  return props.issue?.severity === 'error' 
    ? `${baseClass} error-icon` 
    : `${baseClass} warning-icon-yellow`;
};

const getWarningTitle = (): string => {
  return props.issue?.severity === 'error' 
    ? t('KeyToneAlbum.notify.dependencyDeleted')
    : t('KeyToneAlbum.notify.dependencyIssues');
};

const getWarningTextClass = (): string => {
  return props.issue?.severity === 'error' 
    ? 'warning-text error-text' 
    : 'warning-text warning-text-yellow';
};

const getTooltipClass = (): string => {
  return props.issue?.severity === 'error' 
    ? 'dependency-tooltip error-tooltip' 
    : 'dependency-tooltip warning-tooltip';
};

const getDependencyIconClass = (dep: any): string => {
  return dep.direct ? 'dep-icon-direct' : 'dep-icon-indirect';
};

const getDependencyTextClass = (dep: any): string => {
  return dep.direct ? 'dep-text-direct' : 'dep-text-indirect';
};

const getDependencyIcon = (type: string): string => {
  switch (type) {
    case 'audioFile':
      return 'audio_file';
    case 'sound':
      return 'volume_up';
    case 'keySound':
      return 'keyboard';
    default:
      return 'help_outline';
  }
};

const getDependencyTypeText = (type: string): string => {
  switch (type) {
    case 'audioFile':
      return t('KeyToneAlbum.notify.missingAudioFile');
    case 'sound':
      return t('KeyToneAlbum.notify.missingSound');
    case 'keySound':
      return t('KeyToneAlbum.notify.missingKeySound');
    default:
      return t('KeyToneAlbum.notify.dependencyWarning');
  }
};
</script>

<style lang="scss" scoped>
.dependency-warning {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.warning-icon-yellow {
  color: #f59e0b; /* yellow-500 */
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

.error-icon {
  color: #ef4444; /* red-500 */
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

.warning-text-yellow {
  color: #f59e0b;
  font-size: 0.75rem;
  font-weight: 500;
}

.error-text {
  color: #ef4444;
  font-size: 0.75rem;
  font-weight: 500;
}

:deep(.warning-tooltip) {
  background: rgba(245, 158, 11, 0.95) !important;
  color: white !important;
  border-radius: 8px !important;
  padding: 8px 12px !important;
  font-size: 0.75rem !important;
  max-width: 280px;
  word-wrap: break-word;
}

:deep(.error-tooltip) {
  background: rgba(239, 68, 68, 0.95) !important;
  color: white !important;
  border-radius: 8px !important;
  padding: 8px 12px !important;
  font-size: 0.75rem !important;
  max-width: 280px;
  word-wrap: break-word;
}

.tooltip-content {
  .tooltip-title {
    font-weight: 600;
    margin-bottom: 6px;
    color: white;
  }
  
  .missing-items {
    .missing-item {
      display: flex;
      align-items: center;
      gap: 6px;
      margin-bottom: 4px;
      font-size: 0.70rem;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .dep-icon-direct {
        color: rgba(255, 255, 255, 0.9);
        flex-shrink: 0;
      }

      .dep-icon-indirect {
        color: rgba(255, 255, 255, 0.7);
        flex-shrink: 0;
      }

      .dep-text-direct {
        color: white;
      }

      .dep-text-indirect {
        color: rgba(255, 255, 255, 0.8);
      }

      .indirect-indicator {
        color: rgba(255, 255, 255, 0.6);
        font-style: italic;
      }
    }
  }
}

.warning-content {
  position: relative;
}
</style>