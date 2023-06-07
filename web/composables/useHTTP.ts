import { ElMessage } from 'element-plus'
import { UseFetchOptions } from 'nuxt/dist/app/composables'
const tenant =  import.meta.env.VITE_APP_TENANT

// 指定后端返回的基本数据类型
export interface ResOptions {
    code: number,
    data: any
}
export type HttpOption = UseFetchOptions<ResOptions>

const defaultOption = {
  // 默认request header配置
}

const fetch = (url: string , option: HttpOption) => {
  const baseUrl = import.meta.env.VITE_APP_BASE_API
  url = baseUrl + url
  option = mergeObj(option,{
    defaultOption
  })

  return new Promise((resolve, reject) => {
    useFetch(url, {
      // 请求拦截器
      onRequest({ options }) {
                 
      },
      // 响应拦截
      onResponse({ response }) {
        // 处理 302 重定向
        if (response.redirected) {
          resolve(response.url)
        }
        const { status, _data: data } = response
        if (status === 200) {
          resolve(data)
        } 
        if (status === 204) {
          resolve(data)
        }
        else if ( status === 401 ) {
          // TODO: 临时处理 拦截登录账号密码错误
          if (response.url.includes('/login')) {
            resolve(10000)
          } else if(data.code=='1000'){
            navigateTo('/login')
          }
        }
        // console.log(response)
        // const { code, msg } = response._data;
        // console.log(response);
        
        // if (code === 200) {
        //   resolve(response._data)
        // } else if (code) {
        //   resolve(response._data)
        // } else {
        //   ElMessage({
        //     message: msg || '未知错误',
        //     type: 'error'
        //   });
        //   return Promise.reject(new Error(msg || 'Error'));
        // }
      },
      // 错误处理
      onResponseError({ response }) {
        return Promise.reject(response?._data ?? null)
      },
      // 合并参数
      ...option,
    }).catch((err: any) => {
      reject(err)
    })
  })
}


// 自动导出
export const useHttp = {
  get: (url: string, params?: any, option?: HttpOption) => {
    return fetch(url, { method: 'get', params, ...option })
  },

  post: (url: string, body?: any, option?: HttpOption) => {
    return fetch(url, { method: 'post', body, ...option })
  },

  put: (url: string, body?: any, option?: HttpOption) => {
    return fetch(url, { method: 'put', body, ...option })
  },

  delete: (url: string, body?: any, option?: HttpOption) => {
    return fetch(url, { method: 'delete', body, ...option })
  },
}