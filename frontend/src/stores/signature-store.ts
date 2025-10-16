/*
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { SignatureManager } from 'src/types/signature';

/**
 * 签名管理 Store
 * 用于管理签名数据的状态，支持 SSE 实时同步
 */
export const useSignatureStore = defineStore('signature', () => {
  // 签名管理器数据
  const signatureManager = ref<SignatureManager>({});

  /**
   * 从 SSE 更新签名管理器数据
   * @param data - 来自后端的签名管理器数据
   */
  function updateFromSSE(data: SignatureManager) {
    console.log('[SignatureStore] updateFromSSE 被调用，接收数据:', data);

    // 深度比较，避免不必要的更新
    const currentJson = JSON.stringify(signatureManager.value);
    const newJson = JSON.stringify(data);

    if (currentJson !== newJson) {
      console.log('[SignatureStore] 数据已改变，进行更新');
      signatureManager.value = data;
      console.log('[SignatureStore] 签名数据已更新，当前值:', signatureManager.value);
    } else {
      console.log('[SignatureStore] 数据未改变，跳过更新');
    }
  }

  /**
   * 获取所有签名列表
   */
  function getSignatureList() {
    if (!signatureManager.value) return [];
    return Object.values(signatureManager.value);
  }

  /**
   * 根据 ID 获取签名
   * @param id - 签名 ID
   */
  function getSignatureById(id: string) {
    return signatureManager.value[id];
  }

  /**
   * 检查签名是否存在
   * @param id - 签名 ID
   */
  function hasSignature(id: string) {
    return id in signatureManager.value;
  }

  /**
   * 清空所有签名数据（仅用于测试或重置）
   */
  function reset() {
    signatureManager.value = {};
  }

  return {
    signatureManager,
    updateFromSSE,
    getSignatureList,
    getSignatureById,
    hasSignature,
    reset,
  };
});
