<template>
  <div v-if="visible" class="dependency-warning">
    <!-- 简单警告图标 -->
    <q-icon 
      v-if="showIcon" 
      name="warning" 
      class="warning-icon"
      :size="iconSize"
    />
    
    <!-- 警告内容 -->
    <div v-if="showTooltip" class="warning-content">
      <q-tooltip 
        anchor="bottom middle" 
        self="top middle" 
        :offset="[0, 8]"
        class="dependency-tooltip"
      >
        <div class="tooltip-content">
          <div class="tooltip-title">
            {{ $t('KeyToneAlbum.notify.dependencyWarning') }}
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
                class="dep-icon"
              />
              <span>{{ getDependencyTypeText(dep.type) }}: {{ dep.name || dep.id }}</span>
            </div>
          </div>
        </div>
      </q-tooltip>
    </div>
    
    <!-- 警告文本 -->
    <span v-if="showText" class="warning-text">
      {{ $t('KeyToneAlbum.notify.dependencyWarning') }}
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

.warning-icon {
  color: #f59e0b; /* yellow-500 */
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

.warning-text {
  color: #f59e0b;
  font-size: 0.75rem;
  font-weight: 500;
}

:deep(.dependency-tooltip) {
  background: rgba(245, 158, 11, 0.95) !important;
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
      
      .dep-icon {
        color: rgba(255, 255, 255, 0.9);
        flex-shrink: 0;
      }
    }
  }
}

.warning-content {
  position: relative;
}
</style>