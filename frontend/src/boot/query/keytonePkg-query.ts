import { api } from 'boot/axios';

export async function SendFileToServer(audioPkgUUID: string, file: File) {
  const formData = new FormData();

  formData.append('audioPkgUUID', audioPkgUUID);
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
        return true;
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

export async function SoundFileDelete(audioPkgUUID: string, sha256: string, nameID: string, type: string) {
  return await api
    .post('/keytone_pkg/sound_file_delete', {
      audioPkgUUID: audioPkgUUID, // 目录名/音频包名ID
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
  audioPkgUUID: string,
  sha256: string,
  type: string,
  startTime: number,
  endTime: number,
  volume: number
) {
  return await api
    .post('/keytone_pkg/play_sound', {
      audioPkgUUID: audioPkgUUID,
      sha256: sha256,
      type: type,
      startTime: startTime,
      endTime: endTime,
      volume: volume,
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
