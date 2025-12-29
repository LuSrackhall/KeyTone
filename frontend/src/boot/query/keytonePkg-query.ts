import { api } from 'boot/axios';
import type {
  ApplySignatureConfigPayload,
  AlbumSignatureInfo,
  AvailableSignature,
  CheckSignatureInAlbumResult,
} from 'src/types/export-flow';

export async function SendFileToServer(file: File) {
  const formData = new FormData();

  formData.append('file', file);

  return await api
    .post('/keytone_pkg/add_new_sound_file', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    .then((response) => {
      // console.log('文件上传成功', response.data);
      console.debug('status=', response.status, '->SendFileToServer 请求已成功执行并返回->', response.data);
      if (response.data.message === 'ok') {
        return true;
      } else {
        return false;
      }
    })
    .catch((error) => {
      // console.error('文件上传失败', error);
      console.group('SendFileToServer 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function LoadConfig(audioPkgUUID: string, isCreate: boolean) {
  return await api
    .post('/keytone_pkg/load_config', {
      audioPkgUUID: audioPkgUUID,
      isCreate: isCreate,
    })
    .then((req) => {
      console.debug('status=', req.status, '->LoadConfig 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return req.data;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('LoadConfig 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function ConfigSet(key: string, value: any) {
  return await api
    .post('/keytone_pkg/set', {
      key: key,
      value: value,
    })
    .then((req) => {
      console.debug('status=', req.status, '->ConfigSet 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return true;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('ConfigSet 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function ConfigGet(key: string): Promise<any | false> {
  return await api
    .get('/keytone_pkg/get', {
      params: {
        key: key,
      },
    })
    .then((req) => {
      console.debug('status=', req.status, '->ConfigGet 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return req.data.value;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('ConfigGet 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function ConfigDelete(key: string): Promise<boolean> {
  return await api
    .post('/keytone_pkg/delete', {
      key: key,
    })
    .then((req) => {
      console.debug('status=', req.status, '->ConfigDelete 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return true;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('ConfigDelete 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function SoundFileRename(sha256: string, nameID: string, name: string) {
  return await api
    .post('/keytone_pkg/sound_file_rename', {
      sha256: sha256,
      nameID: nameID,
      name: name,
    })
    .then((req) => {
      console.debug('status=', req.status, '->SoundFileRename 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return true;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('SoundFileRename 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function SoundFileDelete(sha256: string, nameID: string, type: string) {
  return await api
    .post('/keytone_pkg/sound_file_delete', {
      sha256: sha256, // 文件名ID(实际文件名)
      nameID: nameID, // 文件名ID(UI端使用, 用于索引虚拟文件名)
      type: type, // 文件类型
    })
    .then((req) => {
      console.debug('status=', req.status, '->SoundFileDelete 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return true;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('SoundFileDelete 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function PlaySound(
  sha256: string,
  type: string,
  startTime: number,
  endTime: number,
  volume: number,
  isPreviewMode = false
) {
  return await api
    .post('/keytone_pkg/play_sound', {
      sha256: sha256,
      type: type,
      startTime: startTime,
      endTime: endTime,
      volume: volume,
      isPreviewMode: isPreviewMode,
    })
    .then((req) => {
      console.debug('status=', req.status, '->PlaySound 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return true;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('PlaySound 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function GetAudioPackageList() {
  return await api
    .get('/keytone_pkg/get_audio_package_list')
    .then((req) => {
      console.debug('status=', req.status, '->GetAudioPackageList 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return req.data;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('GetAudioPackageList 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function GetAudioPackageName(path: string) {
  return await api
    .get('/keytone_pkg/get_audio_package_name', {
      params: {
        path: path,
      },
    })
    .then((req) => {
      console.debug('status=', req.status, '->GetAudioPackageName 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return req.data;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('GetAudioPackageName 请求执行失败');
      if (error.response) {
        // 请求已经发出，但是服务器返回了一个非 2xx 的状态码
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        // 请求已经发出，但是没有收到响应
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        // 发送请求时出了点问题
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      // 通过打印 error.config，可以查看到导致错误的请求的详细配置，这对于调试和解决问题非常有帮助
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function DeleteAlbum(albumPath: string): Promise<boolean> {
  return await api
    .post('/keytone_pkg/delete_album', {
      albumPath: albumPath,
    })
    .then((req) => {
      console.debug('status=', req.status, '->DeleteAlbum 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return true;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('DeleteAlbum 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

// 获取上传的专辑文件的元数据信息
export interface AlbumMeta {
  magicNumber: string;
  version: string;
  exportTime: string;
  albumUUID: string;
  albumName: string;
}

export async function GetAlbumMeta(file: File): Promise<AlbumMeta> {
  const formData = new FormData();
  formData.append('file', file);

  try {
    const response = await api.post('/keytone_pkg/get_album_meta', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });

    if (response.data.message === 'ok') {
      return response.data.meta;
    }
    throw new Error(response.data.message);
  } catch (err: any) {
    console.group('GetAlbumMeta 请求执行失败');
    const error = err as {
      response?: { status: number; data: any };
      request?: any;
      message?: string;
      config?: any;
    };
    if (error.response) {
      console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
      console.error('Error status:', error.response.status);
      console.error('Error data:', error.response.data);
    } else if (error.request) {
      console.error('Error:', '请求已经发出，但是没有收到响应');
      console.error('Error request:', error.request);
    } else {
      console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
      console.error('Error message:', error.message);
    }
    console.error('Error config:', error.config);
    console.groupEnd();
    throw err;
  }
}

// 加密专辑配置（仅在需要签名时调用）
export async function EncryptAlbumConfig(albumPath: string): Promise<{
  message: string;
  encrypted?: boolean;
  already_encrypted?: boolean;
}> {
  if (!albumPath) {
    throw new Error('缺少 albumPath');
  }
  return await api
    .post('/keytone_pkg/encrypt_album_config', { albumPath })
    .then((response) => {
      console.debug('status=', response.status, '->EncryptAlbumConfig 请求已成功执行并返回');
      if (response.status === 200) {
        return response.data;
      }
      throw new Error('加密配置失败');
    })
    .catch((error) => {
      console.group('EncryptAlbumConfig 请求执行失败');
      console.error('加密专辑配置失败:', error);
      if (error.response?.data?.message) {
        console.error('服务器返回:', error.response.data.message);
      }
      console.groupEnd();
      throw error;
    });
}

// 将签名/授权决策发送到 SDK，供后续写入专辑配置
export async function ApplySignatureConfig(payload: ApplySignatureConfigPayload): Promise<boolean> {
  return await api
    .post('/keytone_pkg/apply_signature_config', payload)
    .then((response) => {
      console.debug('status=', response.status, '->ApplySignatureConfig 请求已成功执行并返回->', response.data);
      if (response.status === 200) {
        return true;
      }
      throw new Error('应用签名配置失败');
    })
    .catch((error) => {
      console.group('ApplySignatureConfig 请求执行失败');
      if (error.response) {
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
      throw error;
    });
}

/**
 * 获取专辑签名信息
 * 用于前端需求2（再次导出时的签名识别）和需求4（签名作者信息展示）
 */
export async function GetAlbumSignatureInfo(albumPath: string): Promise<AlbumSignatureInfo> {
  return await api
    .post('/keytone_pkg/get_album_signature_info', { albumPath })
    .then((response) => {
      console.debug('GetAlbumSignatureInfo 成功:', response.data);
      if (response.status === 200 && response.data.message === 'ok') {
        return response.data.data;
      }
      throw new Error('获取专辑签名信息失败');
    })
    .catch((error) => {
      console.error('GetAlbumSignatureInfo 失败:', error);
      throw error;
    });
}

/**
 * 读取专辑内的文件（如签名图片）
 * @param albumPath 专辑路径
 * @param relativePath 相对于专辑目录的文件路径，如 "audioFiles/xxx.jpg"
 * @returns 文件Blob，失败时返回null
 */
export async function GetAlbumFile(albumPath: string, relativePath: string): Promise<Blob | null> {
  try {
    const response = await api.post(
      '/keytone_pkg/get_album_file',
      { albumPath, relativePath },
      { responseType: 'blob' }
    );
    if (response.status === 200) {
      return response.data;
    }
    return null;
  } catch (error) {
    console.warn('GetAlbumFile 失败:', error);
    return null;
  }
}

/**
 * 检查签名是否在专辑中
 * 用于前端需求3（标记已在专辑中的签名）
 */
export async function CheckSignatureInAlbum(
  albumPath: string,
  signatureId: string
): Promise<CheckSignatureInAlbumResult> {
  return await api
    .post('/keytone_pkg/check_signature_in_album', { albumPath, signatureId })
    .then((response) => {
      console.debug('CheckSignatureInAlbum 成功:', response.data);
      if (response.status === 200 && response.data.message === 'ok') {
        return {
          isInAlbum: response.data.isInAlbum,
          qualificationCode: response.data.qualificationCode,
          hasChanges: response.data.hasChanges,
        };
      }
      throw new Error('检查签名失败');
    })
    .catch((error) => {
      console.error('CheckSignatureInAlbum 失败:', error);
      throw error;
    });
}

/**
 * 检查签名授权状态
 * 用于前端需求3（使能/失能签名选项）
 */
export async function CheckSignatureAuthorization(
  albumPath: string,
  signatureId: string
): Promise<{ isAuthorized: boolean; requireAuthorization: boolean; qualificationCode: string }> {
  return await api
    .post('/keytone_pkg/check_signature_authorization', { albumPath, signatureId })
    .then((response) => {
      console.debug('CheckSignatureAuthorization 成功:', response.data);
      if (response.status === 200 && response.data.message === 'ok') {
        return {
          isAuthorized: response.data.isAuthorized,
          requireAuthorization: response.data.requireAuthorization,
          qualificationCode: response.data.qualificationCode,
        };
      }
      throw new Error('检查签名授权失败');
    })
    .catch((error) => {
      console.error('CheckSignatureAuthorization 失败:', error);
      throw error;
    });
}

/**
 * 获取可用于导出的签名列表
 * 用于前端需求3（签名选择页面增强）
 */
export async function GetAvailableSignatures(albumPath: string): Promise<AvailableSignature[]> {
  return await api
    .post('/keytone_pkg/get_available_signatures', { albumPath })
    .then((response) => {
      console.debug('GetAvailableSignatures 成功:', response.data);
      if (response.status === 200 && response.data.message === 'ok') {
        return response.data.signatures;
      }
      throw new Error('获取可用签名列表失败');
    })
    .catch((error) => {
      console.error('GetAvailableSignatures 失败:', error);
      throw error;
    });
}

// 导出专辑，直接返回zip文件内容
export async function ExportAlbum(albumPath: string): Promise<Blob> {
  return await api
    .post('/keytone_pkg/export_album', { albumPath }, { responseType: 'blob' })
    .then((response) => {
      console.debug('status=', response.status, '->ExportAlbum 请求已成功执行并返回');
      if (response.status === 200) {
        return new Blob([response.data], { type: 'application/zip' });
      }
      throw new Error('导出失败');
    })
    .catch((error) => {
      console.group('ExportAlbum 请求执行失败');
      console.error('导出专辑失败:', error);
      console.groupEnd();
      throw error;
    });
}

// 导入专辑
export async function ImportAlbum(file: File): Promise<boolean> {
  const formData = new FormData();
  formData.append('file', file);

  return await api
    .post('/keytone_pkg/import_album', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    .then((response) => {
      console.debug('status=', response.status, '->ImportAlbum 请求已成功执行并返回->', response.data);
      if (response.data.message === 'ok') {
        return true;
      }
      // 如果是专辑已存在的错误，抛出特殊错误以便UI处理
      if (response.data.message === 'album_exists') {
        throw new Error('album_exists');
      }
      return false;
    })
    .catch((error) => {
      console.group('ImportAlbum 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      throw error; // 继续抛出错误以便UI层处理
    });
}

// 添加新的覆盖导入函数
export async function ImportAlbumOverwrite(file: File): Promise<boolean> {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('overwrite', 'true');

  return await api
    .post('/keytone_pkg/import_album', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    .then((response) => {
      console.debug('status=', response.status, '->ImportAlbumOverwrite 请求已成功执行并返回->', response.data);
      if (response.data.message === 'ok') {
        return true;
      }
      return false;
    })
    .catch((error) => {
      console.group('ImportAlbumOverwrite 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

export async function ImportAlbumAsNew(file: File, newAlbumId: string): Promise<boolean> {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('newAlbumId', newAlbumId);

  return await api
    .post('/keytone_pkg/import_album_as_new', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    .then((response) => {
      console.debug('status=', response.status, '->ImportAlbumAsNew 请求已成功执行并返回->', response.data);
      if (response.data.message === 'ok') {
        return true;
      }
      return false;
    })
    .catch((error) => {
      console.group('ImportAlbumAsNew 请求执行失败');
      if (error.response) {
        console.error('Error:', '请求已经发出且收到响应，但是服务器返回了一个非 2xx 的状态码');
        console.error('Error status:', error.response.status);
        console.error('Error data:', error.response.data);
        if (error.response.status >= 400 && error.response.status < 500) {
          console.error('This is a client error.', '(此为服务端的独断, 若有不服可详细分析)');
        } else if (error.response.status >= 500) {
          console.error('This is a server error.', '(此为服务端的独断, 若有不服可详细分析)');
        }
      } else if (error.request) {
        console.error('Error:', '请求已经发出，但是没有收到响应');
        console.error('Error request:', error.request);
      } else {
        console.error('Error:', '请求未正常发出,请检查请求地址是否正确,或其它种类的错误可能');
        console.error('Error message:', error.message);
      }
      console.error('Error config:', error.config);
      console.groupEnd();
      return false;
    });
}

/**
 * 获取可用于导出的签名列表
 */
export async function GetAvailableSignaturesForExport(albumPath: string): Promise<AvailableSignature[]> {
  return await api.post('/keytone_pkg/get_available_signatures', { albumPath }).then((response) => {
    if (response.status === 200 && response.data.message === 'ok') {
      return response.data.signatures;
    }
    throw new Error('获取可用签名列表失败');
  });
}
