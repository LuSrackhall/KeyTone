import { api } from 'boot/axios';
import type {
  Signature,
  SignatureManager,
  SignatureFile,
  SignatureCreateRequest,
  SignatureUpdateRequest,
  SignatureResponse,
  SignatureListResponse,
} from 'src/types/signature';

/**
 * 获取所有签名列表
 * Get all signatures
 */
export async function getAllSignatures(): Promise<SignatureManager | false> {
  return await api
    .get<SignatureListResponse>('/signature/list')
    .then((req) => {
      console.debug('status=', req.status, '->getAllSignatures 请求已成功执行并返回->', req.data);
      if (req.data.success && req.data.data) {
        return req.data.data;
      } else {
        console.error('Failed to get signatures:', req.data.message);
        return false;
      }
    })
    .catch((error) => {
      console.group('getAllSignatures 请求执行失败');
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
 * 创建新签名
 * Create a new signature
 */
export async function createSignature(data: SignatureCreateRequest): Promise<Signature | false> {
  return await api
    .post<SignatureResponse>('/signature/create', data)
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
export async function updateSignature(data: SignatureUpdateRequest): Promise<Signature | false> {
  return await api
    .put<SignatureResponse>('/signature/update', data)
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
    .delete<SignatureResponse>(`/signature/delete/${id}`)
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
export async function exportSignature(id: string): Promise<SignatureFile | false> {
  return await api
    .get<{ success: boolean; message?: string; data?: SignatureFile }>(`/signature/export/${id}`)
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
export async function importSignature(fileData: SignatureFile): Promise<Signature | false> {
  return await api
    .post<SignatureResponse>('/signature/import', fileData)
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

/**
 * 获取签名图片 URL
 * Get signature image URL
 */
export function getSignatureImageUrl(filename: string): string {
  if (!filename) return '';

  // Get backend port from window API
  const port = (window as any).myWindowAPI?.getBackendPort() || 38888;
  return `http://127.0.0.1:${port}/signature/image/${filename}`;
}

/**
 * 将文件转换为 Base64
 * Convert file to Base64
 */
export function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      if (typeof reader.result === 'string') {
        resolve(reader.result);
      } else {
        reject(new Error('Failed to read file as base64'));
      }
    };
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });
}
