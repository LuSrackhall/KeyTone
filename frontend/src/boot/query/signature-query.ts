import { api } from 'boot/axios';
import type { Signature } from 'src/types/signature';

/**
 * 创建新签名
 * Create a new signature
 */
export async function createSignature(data: Signature): Promise<boolean> {
  return await api
    .post('/signature/create', data)
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
