<template>
  <q-dialog v-model="dialogVisible" persistent>
    <q-card style="width: 90vw; max-width: 480px; max-height: 85vh">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">专辑签名信息</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section v-if="loading" class="flex flex-center" style="min-height: 200px">
        <q-spinner color="primary" size="3em" />
      </q-card-section>

      <q-card-section v-else-if="signatureInfo" style="max-height: calc(85vh - 100px); overflow-y: auto">
        <!-- 无签名提示 -->
        <div v-if="!signatureInfo.hasSignature" class="text-center q-pa-md">
          <q-icon name="info" size="48px" color="grey-6" />
          <div class="text-body1 text-grey-7 q-mt-md">此专辑尚未包含任何签名</div>
        </div>

        <div v-else>
          <!-- 原始作者 -->
          <div v-if="signatureInfo.originalAuthor" class="author-section q-mb-lg">
            <div class="section-title">
              <q-icon name="star" color="amber-7" size="20px" />
              <span class="text-subtitle1 q-ml-sm">原始作者</span>
            </div>
            <q-card flat bordered class="author-card q-mt-sm">
              <q-card-section horizontal>
                <q-img
                  v-if="signatureInfo.originalAuthor.cardImagePath"
                  :src="getImagePath(signatureInfo.originalAuthor.cardImagePath)"
                  style="width: 70px; height: 70px"
                  class="rounded-borders"
                />
                <q-card-section class="col q-pa-sm">
                  <div class="text-subtitle2">{{ signatureInfo.originalAuthor.name }}</div>
                  <div class="text-body2 text-grey-7 q-mt-xs">
                    {{ signatureInfo.originalAuthor.intro }}
                  </div>
                  <div class="q-mt-sm">
                    <q-badge v-if="signatureInfo.originalAuthor.requireAuthorization" color="orange" class="q-mr-sm">
                      需要授权导出
                    </q-badge>
                    <q-badge color="amber-7">原始作者</q-badge>
                  </div>
                  <div
                    v-if="
                      signatureInfo.originalAuthor.authorizedList &&
                      signatureInfo.originalAuthor.authorizedList.length > 0
                    "
                    class="text-caption text-grey-6 q-mt-sm"
                  >
                    已授权 {{ signatureInfo.originalAuthor.authorizedList.length }} 个签名导出
                  </div>
                </q-card-section>
              </q-card-section>
            </q-card>
          </div>

          <!-- 直接导出作者 -->
          <div v-if="signatureInfo.directExportAuthor" class="author-section q-mb-md">
            <div class="section-title">
              <q-icon name="file_download" color="blue-7" size="20px" />
              <span class="text-subtitle1 q-ml-sm">直接导出作者</span>
            </div>
            <q-card flat bordered class="author-card q-mt-sm">
              <q-card-section horizontal>
                <q-img
                  v-if="signatureInfo.directExportAuthor.cardImagePath"
                  :src="getImagePath(signatureInfo.directExportAuthor.cardImagePath)"
                  style="width: 70px; height: 70px"
                  class="rounded-borders"
                />
                <q-card-section class="col q-pa-sm">
                  <div class="text-subtitle2">{{ signatureInfo.directExportAuthor.name }}</div>
                  <div class="text-caption text-grey-7 q-mt-xs">
                    {{ signatureInfo.directExportAuthor.intro }}
                  </div>
                  <div class="q-mt-sm">
                    <q-badge color="blue-7">直接导出作者</q-badge>
                  </div>
                </q-card-section>
              </q-card-section>
            </q-card>
          </div>

          <!-- 历史贡献作者 -->
          <div
            v-if="signatureInfo.contributorAuthors && signatureInfo.contributorAuthors.length > 0"
            class="author-section"
          >
            <div class="section-title">
              <q-icon name="group" color="green-7" size="20px" />
              <span class="text-subtitle1 q-ml-sm">历史贡献作者 ({{ signatureInfo.contributorAuthors.length }})</span>
            </div>
            <q-card
              v-for="contributor in signatureInfo.contributorAuthors"
              :key="contributor.qualificationCode"
              flat
              bordered
              class="author-card q-mt-sm"
            >
              <q-card-section horizontal>
                <q-img
                  v-if="contributor.cardImagePath"
                  :src="getImagePath(contributor.cardImagePath)"
                  style="width: 70px; height: 70px"
                  class="rounded-borders"
                />
                <q-card-section class="col q-pa-sm">
                  <div class="text-subtitle2">{{ contributor.name }}</div>
                  <div class="text-caption text-grey-7 q-mt-xs">
                    {{ contributor.intro }}
                  </div>
                  <div class="q-mt-sm">
                    <q-badge color="green-7">贡献作者</q-badge>
                  </div>
                </q-card-section>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-card-section>

      <q-card-section v-else-if="error" class="text-center">
        <q-icon name="error" size="48px" color="negative" />
        <div class="text-body1 text-negative q-mt-md">{{ error }}</div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="关闭" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { GetAlbumSignatureInfo } from 'src/boot/query/keytonePkg-query';
import type { AlbumSignatureInfo } from 'src/types/export-flow';

interface Props {
  albumPath: string;
}

const props = defineProps<Props>();
const dialogVisible = ref(false);
const loading = ref(false);
const error = ref<string | null>(null);
const signatureInfo = ref<AlbumSignatureInfo | null>(null);

/**
 * 打开对话框并加载签名信息
 */
async function open() {
  dialogVisible.value = true;
  loading.value = true;
  error.value = null;
  signatureInfo.value = null;

  try {
    signatureInfo.value = await GetAlbumSignatureInfo(props.albumPath);
  } catch (err: any) {
    error.value = err.message || '加载签名信息失败';
    console.error('获取专辑签名信息失败:', err);
  } finally {
    loading.value = false;
  }
}

/**
 * 获取图片路径
 */
function getImagePath(cardImagePath: string): string {
  if (!cardImagePath) return '';
  // 如果是相对路径，拼接专辑路径
  if (cardImagePath.startsWith('audioFiles/')) {
    return `file://${props.albumPath}/${cardImagePath}`;
  }
  return cardImagePath;
}

// 暴露方法给父组件
defineExpose({
  open,
});
</script>

<style scoped lang="scss">
.author-section {
  .section-title {
    display: flex;
    align-items: center;
    padding-bottom: 6px;
    border-bottom: 1px solid $grey-4;
  }

  & + .author-section {
    margin-top: 1rem;
  }
}

.author-card {
  transition: all 0.3s;

  &:hover {
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  }
}
</style>
