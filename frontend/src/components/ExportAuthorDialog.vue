<template>
  <q-dialog v-model="dialogVisible" persistent>
    <q-card class="export-dialog" style="width: 90%; max-width: 350px; max-height: 80vh">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">{{ $t('exportDialog.title') || '导出版权信息' }}</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section style="max-height: 60vh; overflow-y: auto;">
        <div class="q-gutter-md">
          <!-- 作者名称 -->
          <q-input
            v-model="form.authorName"
            :label="$t('exportDialog.authorName') || '创作者名称 (可选)'"
            outlined
            dense
            :rules="[validateAuthorField]"
            hint="留空时显示为未提供"
          />

          <!-- 联系方式文本 -->
          <q-input
            v-model="form.authorContact"
            :label="$t('exportDialog.authorContact') || '联系方式 (可选)'"
            outlined
            dense
            type="textarea"
            rows="2"
            :rules="[validateAuthorField]"
            :hint="$t('exportDialog.contactHint') || '可填写邮箱、社交媒体等联系方式'"
          />

          <!-- 联系方式图片 -->
          <div>
            <q-file
              v-model="contactImageFile"
              :label="$t('exportDialog.contactImage') || '联系方式图片 (可选)'"
              outlined
              dense
              accept="image/*"
              @update:model-value="handleImageUpload"
              :hint="$t('exportDialog.imageHint') || '可上传二维码等联系图片'"
            >
              <template v-slot:prepend>
                <q-icon name="image" />
              </template>
            </q-file>
            
            <!-- 图片预览 -->
            <div v-if="imagePreview" class="q-mt-sm">
              <img :src="imagePreview" alt="Contact Image" style="max-width: 100%; max-height: 80px; border-radius: 4px;" />
            </div>
          </div>

          <!-- 历史创作者 -->
          <div v-if="historyAuthors && historyAuthors.length > 0">
            <q-separator class="q-my-md" />
            <div class="text-subtitle2 q-mb-sm">{{ $t('exportDialog.historyAuthors') || '历史创作者' }}</div>
            <div class="row q-gutter-xs">
              <q-chip
                v-for="(author, index) in historyAuthors"
                :key="index"
                :color="selectedHistoryAuthor === author ? 'primary' : undefined"
                :text-color="selectedHistoryAuthor === author ? 'white' : undefined"
                clickable
                size="sm"
                @click="selectHistoryAuthor(author)"
              >
                {{ author }}
              </q-chip>
            </div>
            <div class="text-caption text-grey-6 q-mt-xs">
              {{ $t('exportDialog.historyHint') || '点击选择优先展示的创作者名称' }}
            </div>
          </div>

          <q-separator class="q-my-md" />

          <!-- 二次导出设置 -->
          <div>
            <q-checkbox
              v-model="form.allowReExport"
              :label="$t('exportDialog.allowReExport') || '允许二次导出'"
              :disable="!hasAuthorInfo"
            />
            <div class="text-caption text-grey-6 q-mt-xs">
              {{ $t('exportDialog.reExportHint') || '只有填写名称或联系方式的创作者可以取消勾选' }}
            </div>
          </div>

          <!-- 导出密码 -->
          <div v-if="!form.allowReExport">
            <q-input
              v-model="form.exportPassword"
              :label="$t('exportDialog.password') || '导出密码 (9位以上)'"
              outlined
              dense
              type="password"
              :rules="[validatePassword]"
              :hint="$t('exportDialog.passwordHint') || '密码将以SHA512形式存储，用于后续校验'"
            />
          </div>
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat :label="$t('common.cancel') || '取消'" @click="closeDialog" />
        <q-btn
          unelevated
          color="primary"
          :label="$t('common.export') || '导出'"
          :loading="loading"
          @click="handleExport"
          :disable="!isFormValid"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, defineEmits, defineProps } from 'vue';
import { useI18n } from 'vue-i18n';
import { useQuasar } from 'quasar';

const { t } = useI18n();
const q = useQuasar();

// Props and Emits
const props = defineProps<{
  modelValue: boolean;
  albumPath: string;
  historyAuthors?: string[];
}>();

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  'export': [data: {
    albumPath: string;
    authorName?: string;
    authorContact?: string;
    authorContactImg?: string;
    historyAuthors?: string[];
    allowReExport: boolean;
    exportPassword?: string;
  }];
}>();

// Reactive data
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const form = ref({
  authorName: '',
  authorContact: '',
  allowReExport: true,
  exportPassword: ''
});

const contactImageFile = ref<File | null>(null);
const imagePreview = ref<string>('');
const selectedHistoryAuthor = ref<string>('');
const loading = ref(false);

// Computed properties
const hasAuthorInfo = computed(() => {
  return form.value.authorName.trim() !== '' || form.value.authorContact.trim() !== '';
});

const isFormValid = computed(() => {
  if (!form.value.allowReExport) {
    return validatePassword(form.value.exportPassword) === true;
  }
  return true;
});

// Methods
const validateAuthorField = (value: string) => {
  if (!value) return true; // 允许为空
  
  const lowerValue = value.toLowerCase();
  if (lowerValue.includes('keytone') || lowerValue.includes('lusrackhall')) {
    return t('exportDialog.forbiddenWords') || '不能包含 KeyTone 或 LuSrackhall 字段';
  }
  return true;
};

const validatePassword = (password: string) => {
  if (!password) return t('exportDialog.passwordRequired') || '请输入密码';
  if (password.length < 9) return t('exportDialog.passwordMinLength') || '密码至少需要9位';
  return true;
};

const handleImageUpload = async (file: File | null) => {
  if (!file) {
    imagePreview.value = '';
    return;
  }

  // 验证文件大小 (比如限制在5MB)
  if (file.size > 5 * 1024 * 1024) {
    q.notify({
      type: 'negative',
      message: t('exportDialog.fileTooLarge') || '图片文件过大，请选择小于5MB的图片'
    });
    contactImageFile.value = null;
    return;
  }

  // 生成预览
  const reader = new FileReader();
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);
};

const selectHistoryAuthor = (author: string) => {
  if (selectedHistoryAuthor.value === author) {
    selectedHistoryAuthor.value = '';
    form.value.authorName = '';
  } else {
    selectedHistoryAuthor.value = author;
    form.value.authorName = author;
  }
};

const generateImageMD5 = async (file: File): Promise<string> => {
  // 简化实现：使用文件名和大小作为标识
  // 在实际应用中，可以使用更复杂的哈希算法
  const fileName = file.name;
  const fileSize = file.size;
  const timestamp = Date.now();
  return `${fileName.replace(/[^a-zA-Z0-9]/g, '')}_${fileSize}_${timestamp}`;
};

const handleExport = async () => {
  loading.value = true;
  
  try {
    let authorContactImg = '';
    
    // 处理图片上传
    if (contactImageFile.value) {
      authorContactImg = await generateImageMD5(contactImageFile.value);
      // 这里可以添加实际的图片上传逻辑
      // 为了简化，我们先只保存哈希文件名
    }

    // 构建导出数据
    const exportData = {
      albumPath: props.albumPath,
      authorName: form.value.authorName.trim() || undefined,
      authorContact: form.value.authorContact.trim() || undefined,
      authorContactImg: authorContactImg || undefined,
      historyAuthors: props.historyAuthors,
      allowReExport: form.value.allowReExport,
      exportPassword: form.value.exportPassword || undefined
    };

    emit('export', exportData);
    
  } catch (error) {
    console.error('处理导出数据失败:', error);
    q.notify({
      type: 'negative',
      message: t('exportDialog.exportError') || '处理导出数据失败'
    });
  } finally {
    loading.value = false;
  }
};

const closeDialog = () => {
  dialogVisible.value = false;
};

// Watch for allowReExport changes
watch(() => form.value.allowReExport, (newValue) => {
  if (newValue) {
    form.value.exportPassword = '';
  }
});

// Watch for hasAuthorInfo changes
watch(hasAuthorInfo, (newValue) => {
  if (!newValue) {
    form.value.allowReExport = true;
  }
});
</script>

<style scoped>
.export-dialog {
  border-radius: 12px;
}
</style>