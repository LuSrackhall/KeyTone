import { api } from 'boot/axios';
import type { Signature } from 'src/types/signature';

/**
 * 获取所有签名列表（加密的key-value对）
 * Get all signatures list (encrypted key-value pairs)
 */
export async function getSignaturesList(): Promise<{ [key: string]: string } | false> {
  return await api
    .get('/signature/list')
    .then((req) => {
      console.debug('status=', req.status, '->getSignaturesList 请求已成功执行并返回->', req.data);
      if (req.data.success && req.data.data) {
        // 返回加密的key-value对象
        return req.data.data as { [key: string]: string };
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
 */
export async function decryptSignatureData(encryptedValue: string): Promise<string | false> {
  return await api
    .post('/signature/decrypt', { encryptedValue })
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
 */
export async function updateSignature(data: Signature): Promise<boolean> {
  return await api
    .post('/signature/update', data)
    .then((req) => {
      console.debug('status=', req.status, '->updateSignature 请求已成功执行并返回->', req.data);
      if (req.data.success && req.data.data) {
        return req.data.data;
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
 */
export async function exportSignature(id: string): Promise<boolean> {
  return await api
    .post('/signature/export', { id })
    .then((req) => {
      console.debug('status=', req.status, '->exportSignature 请求已成功执行并返回->', req.data);
      if (req.data.success && req.data.data) {
        return req.data.data;
      } else {
        console.error('Failed to export signature:', req.data.message);
        return false;
      }
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
 */
export async function importSignature(fileData: File): Promise<Signature | false> {
  return await api
    .post('/signature/import', fileData)
    .then((req) => {
      console.debug('status=', req.status, '->importSignature 请求已成功执行并返回->', req.data);
      if (req.data.success && req.data.data) {
        return req.data.data;
      } else {
        console.error('Failed to import signature:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
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
