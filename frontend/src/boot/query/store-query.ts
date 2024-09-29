import { api } from 'boot/axios';
export async function StoreSet(key: string, value: any) {
  return await api
    .post('/store/set', {
      key: key,
      value: value,
    })
    .then((req) => {
      console.debug('status=', req.status, '->StoreSet 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return true;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('StoreSet 请求执行失败');
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

export async function StoreGet(key: string): Promise<any | false> {
  return await api
    .get('/store/get', {
      params: {
        key: key,
      },
    })
    .then((req) => {
      console.debug('status=', req.status, '->StoreGet 请求已成功执行并返回->', req.data);
      if (req.data.message === 'ok') {
        return req.data.value;
      } else {
        return false;
      }
    })
    .catch((error) => {
      console.group('StoreGet 请求执行失败');
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
