<template>
  <div class="dependency-warning-container">
    <!-- Critical dependency warning -->
    <q-icon
      v-if="criticalCount > 0"
      name="error"
      color="warning"
      size="16px"
      class="dependency-warning-icon critical"
    >
      <q-tooltip
        anchor="center right"
        self="center left"
        :offset="[8, 0]"
        class="bg-warning text-black rounded-lg text-xs px-2 py-1"
      >
        {{ $t('KeyToneAlbum.dependencyWarning.tooltip.critical') }}
        <div v-if="showDetails && issues.length > 0" class="mt-2 text-xs">
          <div v-for="issue in criticalIssues" :key="issue.itemId" class="mb-1">
            {{ issue.message }}
          </div>
        </div>
      </q-tooltip>
    </q-icon>

    <!-- Info dependency warning -->
    <q-icon
      v-else-if="infoCount > 0"
      name="warning"
      color="info"
      size="16px"
      class="dependency-warning-icon info"
    >
      <q-tooltip
        anchor="center right"
        self="center left"
        :offset="[8, 0]"
        class="bg-info text-white rounded-lg text-xs px-2 py-1"
      >
        {{ $t('KeyToneAlbum.dependencyWarning.tooltip.info') }}
        <div v-if="showDetails && issues.length > 0" class="mt-2 text-xs">
          <div v-for="issue in infoIssues" :key="issue.itemId" class="mb-1">
            {{ issue.message }}
          </div>
        </div>
      </q-tooltip>
    </q-icon>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import type { DependencyIssue } from 'src/utils/dependencyValidator';

const { t } = useI18n();

interface Props {
  issues: DependencyIssue[];
  itemType: 'audio_files' | 'sounds' | 'key_sounds' | 'global_binding' | 'single_key_binding';
  itemId: string;
  showDetails?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  showDetails: true
});

// Filter issues for this specific item
const itemIssues = computed(() => 
  props.issues.filter(issue => 
    issue.itemType === props.itemType && issue.itemId === props.itemId
  )
);

const criticalIssues = computed(() => 
  itemIssues.value.filter(issue => issue.severity === 'critical')
);

const infoIssues = computed(() => 
  itemIssues.value.filter(issue => issue.severity === 'info')
);

const criticalCount = computed(() => criticalIssues.value.length);
const infoCount = computed(() => infoIssues.value.length);
</script>

<style lang="scss" scoped>
.dependency-warning-container {
  display: inline-flex;
  align-items: center;
  margin-left: 4px;
}

.dependency-warning-icon {
  cursor: help;
  
  &.critical {
    animation: pulse-error 2s infinite;
  }
  
  &.info {
    animation: pulse-warning 3s infinite;
  }
}

@keyframes pulse-error {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.6;
  }
}

@keyframes pulse-warning {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}
</style>