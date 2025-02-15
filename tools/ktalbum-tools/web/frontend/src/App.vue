<template>
  <div class="min-h-screen bg-gray-100">
    <div class="container mx-auto p-6">
      <header class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800">KeyTone 专辑工具</h1>
        <p class="text-gray-600 mt-2">用于处理 KeyTone 专辑文件的工具</p>
      </header>

      <!-- 文件上传区域 -->
      <div
        class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center"
        @dragover.prevent
        @drop.prevent="handleFileDrop"
      >
        <input type="file" ref="fileInput" accept=".ktalbum" @change="handleFileSelect" class="hidden" />
        <div class="space-y-4">
          <div class="text-gray-600">
            <i class="fas fa-cloud-upload-alt text-4xl mb-2"></i>
            <p>拖放文件到此处，或</p>
          </div>
          <button
            @click="$refs.fileInput.click()"
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
          >
            选择文件
          </button>
        </div>
      </div>

      <!-- 文件信息 -->
      <div v-if="selectedFile" class="mt-6">
        <div class="flex items-center justify-between p-4 bg-gray-50 rounded">
          <div>
            <p class="font-medium">已选择文件：{{ selectedFile.name }}</p>
            <p class="text-sm text-gray-500">大小：{{ formatFileSize(selectedFile.size) }}</p>
          </div>
          <div class="space-x-2">
            <button
              @click="extractFile"
              class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition"
              :disabled="isProcessing"
            >
              <span v-if="isProcessing"> <i class="fas fa-spinner fa-spin mr-2"></i>处理中... </span>
              <span v-else>解包文件</span>
            </button>
            <button
              @click="viewFileInfo"
              class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
              :disabled="isProcessing"
            >
              查看信息
            </button>
          </div>
        </div>
      </div>

      <!-- 状态信息 -->
      <div v-if="message" :class="['p-4 rounded-lg mt-4', messageClass, 'animate-fade-in']">
        {{ message }}
      </div>

      <!-- 文件信息弹窗 -->
      <div
        v-if="showInfoModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center"
        @click.self="showInfoModal = false"
      >
        <div class="bg-white rounded-lg p-6 max-w-2xl w-full mx-4">
          <h2 class="text-2xl font-bold mb-4">文件信息</h2>
          <div v-if="fileInfo" class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div class="text-gray-600">文件名</div>
              <div>{{ fileInfo.name }}</div>
              <div class="text-gray-600">版本</div>
              <div>{{ fileInfo.version }}</div>
              <div class="text-gray-600">创建时间</div>
              <div>{{ new Date(fileInfo.exportTime).toLocaleString() }}</div>
              <div class="text-gray-600">专辑 UUID</div>
              <div>{{ fileInfo.albumUUID }}</div>
            </div>
          </div>
          <div class="mt-6 flex justify-end">
            <button
              @click="showInfoModal = false"
              class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition"
            >
              关闭
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import type { FileInfo, FileInfoResponse, UploadResponse } from "./types";

const selectedFile = ref<File | null>(null);
const isProcessing = ref(false);
const message = ref("");
const messageClass = ref("");
const showInfoModal = ref(false);
const fileInfo = ref<FileInfo | null>(null);

// 文件选择处理
const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (input.files?.length) {
    selectedFile.value = input.files[0];
    message.value = "";
  }
};

// 文件拖放处理
const handleFileDrop = (event: DragEvent) => {
  const files = event.dataTransfer?.files;
  if (files?.length) {
    const file = files[0];
    if (file.name.toLowerCase().endsWith(".ktalbum")) {
      selectedFile.value = file;
      message.value = "";
    } else {
      message.value = "请选择 .ktalbum 格式的文件";
      messageClass.value = "bg-red-100 text-red-800";
    }
  }
};

// 文件大小格式化
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
};

// 查看文件信息
const viewFileInfo = async () => {
  if (!selectedFile.value) return;

  isProcessing.value = true;
  try {
    const formData = new FormData();
    formData.append("file", selectedFile.value);

    const response = await fetch("/api/info", {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      throw new Error(await response.text());
    }

    const result = (await response.json()) as FileInfoResponse;
    if (result.data) {
      fileInfo.value = result.data;
      showInfoModal.value = true;
    } else {
      throw new Error("获取文件信息失败");
    }
  } catch (error: any) {
    message.value = `获取文件信息失败: ${error.message}`;
    messageClass.value = "bg-red-100 text-red-800";
  } finally {
    isProcessing.value = false;
  }
};

// 解包文件
const extractFile = async () => {
  if (!selectedFile.value) return;

  isProcessing.value = true;
  message.value = "";

  try {
    const formData = new FormData();
    formData.append("file", selectedFile.value);

    const response = await fetch("/api/extract", {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      throw new Error(await response.text());
    }

    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = selectedFile.value.name.replace(".ktalbum", ".zip");
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);

    message.value = "解包成功！";
    messageClass.value = "bg-green-100 text-green-800";
  } catch (error: any) {
    message.value = `解包失败: ${error.message}`;
    messageClass.value = "bg-red-100 text-red-800";
  } finally {
    isProcessing.value = false;
  }
};
</script>

<style>
.animate-fade-in {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
