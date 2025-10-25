import { api } from 'boot/axios';
import type { Signature, SignatureStorageEntry } from 'src/types/signature';

/**
 * 获取所有签名列表（加密的key-value对，包含排序元数据）
 * Get all signatures list (encrypted key-value pairs with sort metadata)
 *
 * 返回格式可能为：
 * 1. 新格式：{ encryptedId: { value: encryptedData, sort: { time: 1234567890 } } }
 * 2. 旧格式：{ encryptedId: encryptedData } (需要升级处理)
 */
export async function getSignaturesList(): Promise<{ [key: string]: SignatureStorageEntry } | false> {
  return await api
    .get('/signature/list')
    .then((req) => {
      console.debug('status=', req.status, '->getSignaturesList 请求已成功执行并返回->', req.data);
      if (req.data.success && req.data.data) {
        // 处理新格式和旧格式的兼容
        const result: { [key: string]: SignatureStorageEntry } = {};
        const data = req.data.data;

        for (const [key, value] of Object.entries(data)) {
          if (typeof value === 'string') {
            // 旧格式：直接是加密字符串，升级为新格式
            result[key] = {
              value: value,
              sort: {
                time: 0, // TODO: 应该从本地存储或其他方式获取原始创建时间
              },
            };
          } else if (typeof value === 'object' && value !== null) {
            // 新格式：已是 SignatureStorageEntry
            result[key] = value as SignatureStorageEntry;
          }
        }

        return result;
      } else {
        console.error('Failed to get signatures list:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
      console.group('getSignaturesList 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

/**
 * 解密单个签名数据（value值）
 * Decrypt a single signature data (value)
 *
 * 支持新的加密方案：
 * - 如果提供 encryptedId，使用动态密钥解密（新方案）
 * - 如果不提供 encryptedId，使用旧方式解密（向后兼容）
 *
 * @param encryptedValue - 加密的签名数据
 * @param encryptedId - 可选：已加密的签名ID，用于生成动态密钥
 */
export async function decryptSignatureData(encryptedValue: string, encryptedId?: string): Promise<string | false> {
  return await api
    .post('/signature/decrypt', {
      encryptedValue,
      ...(encryptedId && { encryptedId }), // 仅当提供encryptedId时才包含在请求中
    })
    .then((req) => {
      console.debug('status=', req.status, '->decryptSignatureData 请求已成功执行并返回->', req.data);
      if (req.data.success && req.data.data) {
        // 返回解密后的JSON字符串
        return req.data.data as string;
      } else {
        console.error('Failed to decrypt signature data:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
      console.group('decryptSignatureData 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

/**
 * 获取签名图片
 * Get signature image file
 */
export async function getSignatureImage(imagePath: string): Promise<Blob | false> {
  return await api
    .post(
      '/signature/get-image',
      { imagePath },
      {
        responseType: 'blob',
      }
    )
    .then((req) => {
      console.debug('status=', req.status, '->getSignatureImage 请求已成功执行并返回');
      return req.data as Blob;
    })
    .catch((error) => {
      console.group('getSignatureImage 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}
export async function createSignature(data: Signature): Promise<boolean> {
  // 使用 FormData 以支持文件上传
  const formData = new FormData();
  formData.append('id', data.id);
  formData.append('name', data.name);
  formData.append('intro', data.intro);
  // 只在文件存在且有效时才添加文件
  if (data.cardImage && data.cardImage.size > 0) {
    formData.append('cardImage', data.cardImage);
  }

  return await api
    .post('/signature/create', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    .then((req) => {
      console.debug('status=', req.status, '->createSignature 请求已成功执行并返回->', req.data);
      if (req.data.success && req.data.data) {
        return req.data.data;
      } else {
        console.error('Failed to create signature:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
      console.group('createSignature 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

/**
 * 更新签名
 * Update a signature
 *
 * 参数说明：
 * - data: 包含 id (加密的签名ID)、name、intro、cardImage (File | null) 的 Signature 对象
 * - 使用 FormData 格式以支持文件上传
 * - 如果 cardImage 为 null，则不上传新图片（保留原图片）
 * - 如果需要删除图片，cardImage 应为大小为0的空File，并通过 removeImage 标记告知后端
 * - 如果 imageChanged 为 false，表示图片未变更，后端将跳过图片处理逻辑
 */
export async function updateSignature(data: Signature): Promise<boolean> {
  // 使用 FormData 以支持文件上传
  const formData = new FormData();
  formData.append('encryptedId', data.id);
  formData.append('name', data.name);
  formData.append('intro', data.intro);

  // 检查是否需要删除图片：如果 cardImage 是大小为0的空File
  const isRemovingImage = data.cardImage && data.cardImage.size === 0 && data.cardImage.name === '';

  // 判断图片是否发生变更
  const hasImageChanged = data.imageChanged ?? true; // 默认为 true（向后兼容）

  if (isRemovingImage) {
    // 添加删除图片标记
    formData.append('removeImage', 'true');
  } else if (data.cardImage && data.cardImage.size > 0 && hasImageChanged) {
    // 仅当图片发生变更且不是删除时才上传
    formData.append('cardImage', data.cardImage);
  }

  // 总是传递 imageChanged 标记，让后端了解图片是否变更
  formData.append('imageChanged', String(hasImageChanged));

  return await api
    .post('/signature/update', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    .then((req) => {
      console.debug('status=', req.status, '->updateSignature 请求已成功执行并返回->', req.data);
      if (req.data.success) {
        return true;
      } else {
        console.error('Failed to update signature:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
      console.group('updateSignature 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

/**
 * 删除签名
 * Delete a signature
 */
export async function deleteSignature(id: string): Promise<boolean> {
  return await api
    .post('/signature/delete', { id })
    .then((req) => {
      console.debug('status=', req.status, '->deleteSignature 请求已成功执行并返回->', req.data);
      if (req.data.success) {
        return true;
      } else {
        console.error('Failed to delete signature:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
      console.group('deleteSignature 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

/**
 * 导出签名为 .ktsign 文件
 * Export signature as .ktsign file
 *
 * @param encryptedId 加密的签名ID
 * @returns Promise<Blob | false> 返回二进制文件数据，或 false 表示失败
 */
export async function exportSignature(encryptedId: string): Promise<Blob | false> {
  return await api
    .post(
      '/signature/export',
      { encryptedId },
      {
        responseType: 'arraybuffer', // 获取二进制数据
      }
    )
    .then((req) => {
      console.debug('status=', req.status, '->exportSignature 请求已成功执行并返回');

      // 将 ArrayBuffer 转换为 Blob
      const blob = new Blob([req.data], { type: 'application/octet-stream' });
      return blob;
    })
    .catch((error) => {
      console.group('exportSignature 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

/**
 * 导入签名从 .ktsign 文件
 * Import signature from .ktsign file
 *
 * @param fileData 要导入的 .ktsign 文件
 * @returns { conflict?: boolean; encryptedId?: string; name?: string } | false
 *          返回导入结果。如果 conflict 为 true，表示签名已存在，需要用户确认
 */
export async function importSignature(fileData: File): Promise<any> {
  const formData = new FormData();
  formData.append('file', fileData);

  return await api
    .post('/signature/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    .then((req) => {
      console.debug('status=', req.status, '->importSignature 请求已成功执行并返回->', req.data);

      // 处理冲突情况 (409 Conflict)
      if (req.status === 409 || req.data.conflict) {
        console.warn('Signature already exists, conflict detected');
        return {
          conflict: true,
          encryptedId: req.data.data?.encryptedId,
          name: req.data.data?.name,
        };
      }

      // 处理成功导入
      if (req.data.success && req.data.data) {
        console.debug('Signature imported successfully');
        return {
          success: true,
          encryptedId: req.data.data.encryptedId,
          name: req.data.data.name,
        };
      } else {
        console.error('Failed to import signature:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
      // 检查是否是冲突错误 (409)
      if (error.response?.status === 409) {
        console.warn('Signature conflict detected:', error.response.data);
        return {
          conflict: true,
          encryptedId: error.response.data.data?.encryptedId,
          name: error.response.data.data?.name,
        };
      }

      console.group('importSignature 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

/**
 * 确认导入签名（处理冲突选择）
 * Confirm import signature with overwrite option
 *
 * @param encryptedId 签名的加密ID
 * @param fileContent 文件的加密内容（Base64/十六进制字符串）
 * @param overwrite 是否覆盖现有签名
 */
export async function confirmImportSignature(
  encryptedId: string,
  fileContent: string,
  overwrite: boolean
): Promise<any> {
  return await api
    .post('/signature/import-confirm', {
      file: fileContent,
      overwrite,
    })
    .then((req) => {
      console.debug('status=', req.status, '->confirmImportSignature 请求已成功执行并返回->', req.data);

      if (req.data.success && req.data.data) {
        console.debug('Signature import confirmed successfully');
        return {
          success: true,
          encryptedId: req.data.data.encryptedId,
          name: req.data.data.name,
        };
      } else {
        console.error('Failed to confirm import signature:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
      console.group('confirmImportSignature 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

/**
 * 更新签名排序
 * Update signature sort order (for drag-and-drop reordering)
 *
 * 参数说明：
 * - sortOrder: { id: 签名ID, sortTime: Unix时间戳 } 的数组，表示新的排序顺序
 *
 * 功能：
 * 1. 前端用户拖动排序完成后调用此函数
 * 2. 向后端 POST /signature/update-sort
 * 3. 后端需要根据提供的 sort.time 值更新配置文件
 */
export async function updateSignatureSort(sortOrder: Array<{ id: string; sortTime: number }>): Promise<boolean> {
  return await api
    .post('/signature/update-sort', { sortOrder })
    .then((req) => {
      console.debug('status=', req.status, '->updateSignatureSort 请求已成功执行并返回->', req.data);
      if (req.data.success) {
        return true;
      } else {
        console.error('Failed to update signature sort:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
      console.group('updateSignatureSort 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}
