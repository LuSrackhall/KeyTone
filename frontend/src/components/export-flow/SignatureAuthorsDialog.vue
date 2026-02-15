<template>
  <q-dialog v-model="dialogVisible" persistent>
    <q-card style="width: 90vw; max-width: 520px; max-height: 85vh">
      <!-- 头部 (sticky) -->
      <q-card-section class="row items-center q-pb-none bg-deep-purple-1 bg-opacity-90 sticky top-0 z-10 backdrop-blur-sm">
        <q-icon name="badge" size="24px" color="deep-purple" class="q-mr-sm" />
        <div class="text-h6">{{ t('exportFlow.signatureInfoDialog.title') }}</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <!-- 加载中 -->
      <q-card-section v-if="loading" class="flex flex-center" style="min-height: 200px">
        <q-spinner color="primary" size="3em" />
      </q-card-section>

      <!-- 签名信息内容 -->
      <q-card-section v-else-if="signatureInfo" style="max-height: calc(85vh - 120px); overflow-y: auto">
        <!-- 无签名提示 -->
        <div v-if="!signatureInfo.hasSignature" class="text-center q-pa-lg">
          <q-icon name="info_outline" size="64px" color="grey-5" />
          <div class="text-body1 text-grey-7 q-mt-md">{{ t('exportFlow.signatureInfoDialog.noSignature') }}</div>
          <div class="text-caption text-grey-5 q-mt-sm">
            {{ t('exportFlow.signatureInfoDialog.noSignatureHint') }}
          </div>
        </div>

        <div v-else class="q-gutter-md">
          <!-- ====================== 原始作者区块 ====================== -->
          <q-card v-if="signatureInfo.originalAuthor" flat bordered class="signature-section">
            <q-card-section class="q-pa-sm bg-amber-1">
              <div class="row items-center">
                <q-icon name="star" color="amber-8" size="20px" />
                <span class="text-subtitle2 text-amber-10 q-ml-sm">{{
                  t('exportFlow.signatureInfoDialog.originalAuthor')
                }}</span>
                <q-space />
                <q-badge
                  v-if="signatureInfo.originalAuthor.requireAuthorization"
                  color="orange"
                  text-color="white"
                  class="text-weight-medium"
                >
                  {{ t('exportFlow.signatureInfoDialog.requireAuth') }}
                </q-badge>
                <q-badge v-else color="green" text-color="white" class="text-weight-medium">
                  {{ t('exportFlow.signatureInfoDialog.noAuthRequired') }}
                </q-badge>
              </div>
            </q-card-section>

            <q-card-section class="q-pa-md">
              <!-- 签名卡片 -->
              <div class="row items-start q-gutter-sm q-mb-md">
                <q-img
                  v-if="
                    signatureInfo.originalAuthor.cardImagePath &&
                    getImageUrl(signatureInfo.originalAuthor.cardImagePath)
                  "
                  :src="getImageUrl(signatureInfo.originalAuthor.cardImagePath)"
                  style="width: 64px; height: 64px"
                  class="rounded-borders shadow-1"
                  fit="cover"
                >
                  <template v-slot:error>
                    <div class="flex items-center justify-center bg-grey-3" style="width: 64px; height: 64px">
                      <q-icon name="person" size="32px" color="grey-5" />
                    </div>
                  </template>
                </q-img>
                <div
                  v-else
                  class="flex items-center justify-center bg-grey-2 rounded-borders"
                  style="width: 64px; height: 64px"
                >
                  <q-icon name="person" size="32px" color="grey-5" />
                </div>

                <div class="col" style="min-width: 0">
                  <div
                    class="text-subtitle1 text-weight-bold"
                    :class="scrollableTextClasses.name"
                    style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                  >
                    {{ signatureInfo.originalAuthor.name || t('exportFlow.signatureInfoDialog.noName') }}
                  </div>
                  <div
                    v-if="signatureInfo.originalAuthor.intro"
                    class="text-body2 text-grey-7 q-mt-xs"
                    :class="scrollableTextClasses.intro"
                    style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                  >
                    {{ signatureInfo.originalAuthor.intro }}
                  </div>
                  <div v-else class="text-caption text-grey-5 q-mt-xs">
                    {{ t('exportFlow.signatureInfoDialog.noIntro') }}
                  </div>
                </div>
              </div>

              <!-- 资格码指纹 -->
              <q-item dense class="q-pa-none q-mb-sm">
                <q-item-section avatar style="min-width: 28px">
                  <q-icon name="fingerprint" color="grey-6" size="18px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label caption>{{
                    t('exportFlow.signatureInfoDialog.qualificationFingerprint')
                  }}</q-item-label>
                  <q-item-label
                    class="text-caption"
                    style="font-family: monospace; word-break: break-all; line-height: 1.4"
                  >
                    {{ signatureInfo.originalAuthor.qualificationFingerprint }}
                  </q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-btn
                    flat
                    dense
                    size="xs"
                    icon="content_copy"
                    @click="
                      copyToClipboard(
                        signatureInfo.originalAuthor.qualificationFingerprint,
                        t('exportFlow.signatureInfoDialog.qualificationFingerprint')
                      )
                    "
                  />
                </q-item-section>
              </q-item>

              <!-- 授权信息（从 allSignatures 获取完整的授权元数据） -->
              <template v-if="originalAuthorEntry?.authorization">
                <q-separator class="q-my-sm" />

                <!-- 联系方式 -->
                <div class="text-overline text-grey-7 q-mb-xs">
                  {{ t('exportFlow.signatureInfoDialog.contactSection') }}
                </div>

                <q-item dense class="q-pa-none q-mb-xs">
                  <q-item-section avatar style="min-width: 28px">
                    <q-icon name="email" color="grey-6" size="18px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label caption>{{ t('exportFlow.signatureInfoDialog.email') }}</q-item-label>
                    <q-item-label class="text-body2">
                      {{
                        originalAuthorEntry.authorization.contactEmail || t('exportFlow.signatureInfoDialog.noEmail')
                      }}
                    </q-item-label>
                  </q-item-section>
                  <q-item-section side v-if="originalAuthorEntry.authorization.contactEmail">
                    <q-btn
                      flat
                      dense
                      size="xs"
                      icon="content_copy"
                      @click="
                        copyToClipboard(
                          originalAuthorEntry.authorization.contactEmail,
                          t('exportFlow.signatureInfoDialog.email')
                        )
                      "
                    />
                  </q-item-section>
                </q-item>

                <q-item v-if="originalAuthorEntry.authorization.contactAdditional" dense class="q-pa-none q-mb-xs">
                  <q-item-section avatar style="min-width: 28px">
                    <q-icon name="chat" color="grey-6" size="18px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label caption>{{ t('exportFlow.signatureInfoDialog.additionalContact') }}</q-item-label>
                    <q-item-label class="text-body2" style="white-space: pre-wrap; word-break: break-word">
                      {{ originalAuthorEntry.authorization.contactAdditional }}
                    </q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-btn
                      flat
                      dense
                      size="xs"
                      icon="content_copy"
                      @click="
                        copyToClipboard(
                          originalAuthorEntry.authorization.contactAdditional,
                          t('exportFlow.signatureInfoDialog.additionalContact')
                        )
                      "
                    />
                  </q-item-section>
                </q-item>

                <q-separator class="q-my-sm" />

                <!-- 授权统计 -->
                <div class="text-overline text-grey-7 q-mb-xs">
                  {{ t('exportFlow.signatureInfoDialog.authStatus') }}
                </div>

                <div class="row q-gutter-sm">
                  <q-chip
                    dense
                    :color="originalAuthorEntry.authorization.requireAuthorization ? 'orange' : 'green'"
                    text-color="white"
                    icon="security"
                    size="sm"
                  >
                    {{
                      originalAuthorEntry.authorization.requireAuthorization
                        ? t('exportFlow.signatureInfoDialog.requireAuthorization')
                        : t('exportFlow.signatureInfoDialog.noAuthorization')
                    }}
                  </q-chip>

                  <!-- 已授权签名数量（仅在需要授权时显示） -->
                  <q-chip
                    v-if="originalAuthorEntry.authorization.requireAuthorization"
                    dense
                    color="blue-grey"
                    text-color="white"
                    icon="people"
                    size="sm"
                  >
                    {{
                      t('exportFlow.signatureInfoDialog.authorizedCount', {
                        count: originalAuthorEntry.authorization.authorizedFingerprintList?.length || 0,
                      })
                    }}
                  </q-chip>
                </div>

                <!-- 授权标识UUID -->
                <q-item dense class="q-pa-none q-mt-sm">
                  <q-item-section avatar style="min-width: 28px">
                    <q-icon name="vpn_key" color="grey-6" size="18px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label caption>{{ t('exportFlow.signatureInfoDialog.authUUID') }}</q-item-label>
                    <q-item-label
                      class="text-caption"
                      style="font-family: monospace; word-break: break-all; line-height: 1.4"
                    >
                      {{
                        originalAuthorEntry.authorization.authorizationUUID ||
                        t('exportFlow.signatureInfoDialog.noAuthUUID')
                      }}
                    </q-item-label>
                  </q-item-section>
                  <q-item-section side v-if="originalAuthorEntry.authorization.authorizationUUID">
                    <q-btn
                      flat
                      dense
                      size="xs"
                      icon="content_copy"
                      @click="
                        copyToClipboard(
                          originalAuthorEntry.authorization.authorizationUUID,
                          t('exportFlow.signatureInfoDialog.authUUID')
                        )
                      "
                    />
                  </q-item-section>
                </q-item>

                <!-- 最近导出者资格码指纹 -->
                <q-item dense class="q-pa-none q-mt-xs">
                  <q-item-section avatar style="min-width: 28px">
                    <q-icon name="file_download" color="grey-6" size="18px" />
                  </q-item-section>
                  <q-item-section>
                    <q-item-label caption>{{
                      t('exportFlow.signatureInfoDialog.latestExporterFingerprint')
                    }}</q-item-label>
                    <q-item-label
                      class="text-caption"
                      style="font-family: monospace; word-break: break-all; line-height: 1.4"
                    >
                      {{
                        originalAuthorEntry.authorization.directExportAuthorFingerprint ||
                        t('exportFlow.signatureInfoDialog.noLatestExporter')
                      }}
                    </q-item-label>
                  </q-item-section>
                  <q-item-section side v-if="originalAuthorEntry.authorization.directExportAuthorFingerprint">
                    <q-btn
                      flat
                      dense
                      size="xs"
                      icon="content_copy"
                      @click="
                        copyToClipboard(
                          originalAuthorEntry.authorization.directExportAuthorFingerprint,
                          t('exportFlow.signatureInfoDialog.latestExporterFingerprint')
                        )
                      "
                    />
                  </q-item-section>
                </q-item>

                <!-- 已授权列表（展开/折叠）- 仅在需要授权且有已授权签名时显示 -->
                <q-expansion-item
                  v-if="
                    originalAuthorEntry.authorization.requireAuthorization &&
                    originalAuthorEntry.authorization.authorizedFingerprintList?.length
                  "
                  dense
                  header-class="q-pa-none q-mt-sm"
                  expand-icon-class="text-grey-6"
                >
                  <template v-slot:header>
                    <q-item-section avatar style="min-width: 28px">
                      <q-icon name="checklist" color="grey-6" size="18px" />
                    </q-item-section>
                    <q-item-section>
                      <q-item-label caption>
                        {{ t('exportFlow.signatureInfoDialog.authorizedList') }} ({{
                          originalAuthorEntry.authorization.authorizedFingerprintList.length
                        }})
                      </q-item-label>
                    </q-item-section>
                  </template>

                  <div class="q-pl-lg q-pr-sm">
                    <div
                      v-for="(fingerprint, idx) in originalAuthorEntry.authorization.authorizedFingerprintList"
                      :key="idx"
                      class="text-caption q-py-xs row items-center"
                      style="font-family: monospace; word-break: break-all; border-bottom: 1px dashed #eee"
                    >
                      <span class="col">{{ idx + 1 }}. {{ fingerprint }}</span>
                      <q-btn
                        flat
                        dense
                        size="xs"
                        icon="content_copy"
                        class="q-ml-xs"
                        @click="
                          copyToClipboard(fingerprint, t('exportFlow.signatureInfoDialog.qualificationFingerprint'))
                        "
                      />
                    </div>
                  </div>
                </q-expansion-item>
              </template>
            </q-card-section>
          </q-card>

          <!-- ====================== 直接导出作者区块 ====================== -->
          <q-card
            v-if="signatureInfo.directExportAuthor && !isDirectExportAuthorSameAsOriginal"
            flat
            bordered
            class="signature-section"
          >
            <q-card-section class="q-pa-sm bg-blue-1">
              <div class="row items-center">
                <q-icon name="file_download" color="blue-7" size="20px" />
                <span class="text-subtitle2 text-blue-9 q-ml-sm">{{
                  t('exportFlow.signatureInfoDialog.directExportAuthor')
                }}</span>
                <q-space />
                <q-badge color="blue" text-color="white" class="text-weight-medium">
                  {{ t('exportFlow.signatureInfoDialog.latestExporter') }}
                </q-badge>
              </div>
            </q-card-section>

            <q-card-section class="q-pa-md">
              <div class="row items-start q-gutter-sm">
                <q-img
                  v-if="
                    signatureInfo.directExportAuthor.cardImagePath &&
                    getImageUrl(signatureInfo.directExportAuthor.cardImagePath)
                  "
                  :src="getImageUrl(signatureInfo.directExportAuthor.cardImagePath)"
                  style="width: 56px; height: 56px"
                  class="rounded-borders shadow-1"
                  fit="cover"
                >
                  <template v-slot:error>
                    <div class="flex items-center justify-center bg-grey-3" style="width: 56px; height: 56px">
                      <q-icon name="person" size="28px" color="grey-5" />
                    </div>
                  </template>
                </q-img>
                <div
                  v-else
                  class="flex items-center justify-center bg-grey-2 rounded-borders"
                  style="width: 56px; height: 56px"
                >
                  <q-icon name="person" size="28px" color="grey-5" />
                </div>

                <div class="col" style="min-width: 0">
                  <div
                    class="text-subtitle2 text-weight-bold"
                    :class="scrollableTextClasses.name"
                    style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                  >
                    {{ signatureInfo.directExportAuthor.name || t('exportFlow.signatureInfoDialog.noName') }}
                  </div>
                  <div
                    v-if="signatureInfo.directExportAuthor.intro"
                    class="text-caption text-grey-7 q-mt-xs"
                    :class="scrollableTextClasses.intro"
                    style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                  >
                    {{ signatureInfo.directExportAuthor.intro }}
                  </div>
                </div>
              </div>

              <!-- 资格码指纹 -->
              <q-item dense class="q-pa-none q-mt-sm">
                <q-item-section avatar style="min-width: 28px">
                  <q-icon name="fingerprint" color="grey-6" size="16px" />
                </q-item-section>
                <q-item-section>
                  <q-item-label caption>{{
                    t('exportFlow.signatureInfoDialog.qualificationFingerprint')
                  }}</q-item-label>
                  <q-item-label
                    class="text-caption"
                    style="font-family: monospace; word-break: break-all; line-height: 1.4; font-size: 0.7rem"
                  >
                    {{ signatureInfo.directExportAuthor.qualificationFingerprint }}
                  </q-item-label>
                </q-item-section>
              </q-item>
            </q-card-section>
          </q-card>

          <!-- ====================== 历史贡献作者区块 ====================== -->
          <q-card
            v-if="signatureInfo.contributorAuthors && signatureInfo.contributorAuthors.length > 0"
            flat
            bordered
            class="signature-section"
          >
            <q-card-section class="q-pa-sm bg-green-1">
              <div class="row items-center">
                <q-icon name="group" color="green-7" size="20px" />
                <span class="text-subtitle2 text-green-9 q-ml-sm">
                  {{ t('exportFlow.signatureInfoDialog.contributorAuthors') }} ({{
                    signatureInfo.contributorAuthors.length
                  }})
                </span>
              </div>
            </q-card-section>

            <q-card-section class="q-pa-sm">
              <q-list dense separator>
                <q-item v-for="contributor in signatureInfo.contributorAuthors" :key="contributor.qualificationCode">
                  <q-item-section avatar>
                    <q-avatar size="40px">
                      <q-img
                        v-if="contributor.cardImagePath && getImageUrl(contributor.cardImagePath)"
                        :src="getImageUrl(contributor.cardImagePath)"
                        fit="cover"
                      >
                        <template v-slot:error>
                          <div class="flex items-center justify-center bg-grey-3 full-width full-height">
                            <q-icon name="person" size="20px" color="grey-5" />
                          </div>
                        </template>
                      </q-img>
                      <q-icon v-else name="person" size="20px" color="grey-5" />
                    </q-avatar>
                  </q-item-section>

                  <q-item-section style="min-width: 0">
                    <q-item-label
                      :class="scrollableTextClasses.name"
                      style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                    >
                      {{ contributor.name || t('exportFlow.signatureInfoDialog.noName') }}
                    </q-item-label>
                    <q-item-label
                      caption
                      v-if="contributor.intro"
                      :class="scrollableTextClasses.intro"
                      style="overflow-x: auto; overflow-y: hidden; text-overflow: clip; white-space: nowrap"
                    >
                      {{ contributor.intro }}
                    </q-item-label>
                  </q-item-section>

                  <q-item-section side>
                    <q-btn
                      flat
                      dense
                      size="xs"
                      icon="content_copy"
                      @click="
                        copyToClipboard(
                          contributor.qualificationFingerprint,
                          t('exportFlow.signatureInfoDialog.qualificationFingerprint')
                        )
                      "
                    >
                      <q-tooltip>{{ t('exportFlow.signatureInfoDialog.qualificationFingerprint') }}</q-tooltip>
                    </q-btn>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-card-section>
          </q-card>

          <!-- ====================== 签名统计摘要 ====================== -->
          <q-card flat bordered class="bg-grey-1">
            <q-card-section class="q-pa-sm">
              <div class="text-overline text-grey-7 q-mb-xs">
                {{ t('exportFlow.signatureInfoDialog.signatureStats') }}
              </div>
              <div class="row q-gutter-sm">
                <q-chip dense outline color="grey-7" size="sm" icon="numbers">
                  {{
                    t('exportFlow.signatureInfoDialog.totalSignatures', {
                      count: Object.keys(signatureInfo.allSignatures || {}).length,
                    })
                  }}
                </q-chip>
                <q-chip v-if="signatureInfo.originalAuthor" dense outline color="amber-8" size="sm" icon="star">
                  {{ t('exportFlow.signatureInfoDialog.originalAuthorCount') }}
                </q-chip>
                <q-chip
                  v-if="signatureInfo.contributorAuthors?.length"
                  dense
                  outline
                  color="green-7"
                  size="sm"
                  icon="group"
                >
                  {{
                    t('exportFlow.signatureInfoDialog.contributorCount', {
                      count: signatureInfo.contributorAuthors.length,
                    })
                  }}
                </q-chip>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </q-card-section>

      <!-- 错误提示 -->
      <q-card-section v-else-if="error" class="text-center q-pa-lg">
        <q-icon name="error_outline" size="64px" color="negative" />
        <div class="text-body1 text-negative q-mt-md">{{ error }}</div>
        <q-btn
          flat
          color="primary"
          :label="t('exportFlow.signatureInfoDialog.retry')"
          icon="refresh"
          class="q-mt-md"
          @click="open"
        />
      </q-card-section>

      <!-- 底部操作 (sticky glass) -->
      <q-card-actions align="right" class="q-pa-sm" :class="['sticky bottom-0 z-10 bg-white/30 backdrop-blur-sm']">
        <q-btn flat :label="t('exportFlow.signatureInfoDialog.close')" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import { GetAlbumSignatureInfo, GetAlbumFile } from 'src/boot/query/keytonePkg-query';
import type { AlbumSignatureInfo, AlbumSignatureEntry } from 'src/types/export-flow';

const { t } = useI18n();

interface Props {
  albumPath: string;
}

const props = defineProps<Props>();
const $q = useQuasar();

const dialogVisible = ref(false);
const loading = ref(false);
const error = ref<string | null>(null);
const signatureInfo = ref<AlbumSignatureInfo | null>(null);

/** 图片URL缓存 Map<cardImagePath, blobUrl> */
const imageUrlCache = ref<Map<string, string>>(new Map());

/**
 * 横向滚动条样式类（与签名列表保持一致）
 */
const scrollableTextClasses = {
  name: [
    'max-w-full !overflow-x-auto whitespace-nowrap !text-clip',
    'h-5.5 [&::-webkit-scrollbar]:h-0.4 [&::-webkit-scrollbar-track]:bg-blueGray-400/50 [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
  ],
  intro: [
    'max-w-full !overflow-x-auto whitespace-nowrap',
    'h-4.4 [&::-webkit-scrollbar]:h-0.3 [&::-webkit-scrollbar-track]:bg-blueGray-400/50 [&::-webkit-scrollbar-thumb]:bg-blueGray-500/40 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:hover:bg-blue-400',
  ],
};

/**
 * 从 allSignatures 获取原始作者的完整签名条目（包含授权元数据）
 */
const originalAuthorEntry = computed<AlbumSignatureEntry | null>(() => {
  if (!signatureInfo.value?.originalAuthor || !signatureInfo.value.allSignatures) {
    return null;
  }
  const qualCode = signatureInfo.value.originalAuthor.qualificationCode;
  return signatureInfo.value.allSignatures[qualCode] || null;
});

/**
 * 检查直接导出作者是否与原始作者相同
 */
const isDirectExportAuthorSameAsOriginal = computed(() => {
  if (!signatureInfo.value?.originalAuthor || !signatureInfo.value?.directExportAuthor) {
    return false;
  }
  return (
    signatureInfo.value.originalAuthor.qualificationCode === signatureInfo.value.directExportAuthor.qualificationCode
  );
});

/**
 * 打开对话框并加载签名信息
 */
async function open() {
  dialogVisible.value = true;
  loading.value = true;
  error.value = null;
  signatureInfo.value = null;
  // 清理旧的图片缓存
  clearImageCache();

  try {
    signatureInfo.value = await GetAlbumSignatureInfo(props.albumPath);
    // 加载所有签名图片
    await loadAllImages();
  } catch (err: any) {
    error.value = err.message || '加载签名信息失败';
    console.error('获取专辑签名信息失败:', err);
  } finally {
    loading.value = false;
  }
}

/**
 * 加载所有签名图片
 */
async function loadAllImages() {
  if (!signatureInfo.value) return;

  const imagePaths: string[] = [];

  // 收集所有需要加载的图片路径
  if (signatureInfo.value.originalAuthor?.cardImagePath) {
    imagePaths.push(signatureInfo.value.originalAuthor.cardImagePath);
  }
  if (signatureInfo.value.directExportAuthor?.cardImagePath) {
    imagePaths.push(signatureInfo.value.directExportAuthor.cardImagePath);
  }
  if (signatureInfo.value.contributorAuthors) {
    for (const contributor of signatureInfo.value.contributorAuthors) {
      if (contributor.cardImagePath) {
        imagePaths.push(contributor.cardImagePath);
      }
    }
  }

  // 去重并并行加载
  const uniquePaths = [...new Set(imagePaths)];
  await Promise.all(uniquePaths.map((path) => loadImage(path)));
}

/**
 * 加载单个图片
 */
async function loadImage(cardImagePath: string) {
  if (!cardImagePath || imageUrlCache.value.has(cardImagePath)) return;

  try {
    const blob = await GetAlbumFile(props.albumPath, cardImagePath);
    if (blob) {
      const url = URL.createObjectURL(blob);
      imageUrlCache.value.set(cardImagePath, url);
    }
  } catch (err) {
    console.warn('加载签名图片失败:', cardImagePath, err);
  }
}

/**
 * 获取图片URL（从缓存中获取）
 */
function getImageUrl(cardImagePath: string): string {
  if (!cardImagePath) return '';
  return imageUrlCache.value.get(cardImagePath) || '';
}

/**
 * 清理图片缓存（释放Blob URL）
 */
function clearImageCache() {
  for (const url of imageUrlCache.value.values()) {
    URL.revokeObjectURL(url);
  }
  imageUrlCache.value.clear();
}

/**
 * 复制到剪贴板
 */
function copyToClipboard(text: string, label: string) {
  navigator.clipboard.writeText(text).then(
    () => {
      $q.notify({
        type: 'positive',
        message: t('exportFlow.signatureInfoDialog.copySuccess', { label }),
        position: 'top',
        timeout: 1500,
      });
    },
    () => {
      $q.notify({
        type: 'negative',
        message: t('exportFlow.signatureInfoDialog.copyFailed'),
        position: 'top',
        timeout: 1500,
      });
    }
  );
}

// 对话框关闭时清理资源
watch(dialogVisible, (visible) => {
  if (!visible) {
    clearImageCache();
  }
});

// 暴露方法给父组件
defineExpose({
  open,
});
</script>

<style scoped lang="scss">
.signature-section {
  border-radius: 8px;
  overflow: hidden;

  & + .signature-section {
    margin-top: 12px;
  }
}

:deep(.q-item__section) {
  @apply text-wrap;
  @apply overflow-hidden;
}
</style>
