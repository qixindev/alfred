// const tenant =  import.meta.env.VITE_APP_TENANT as string
const tenant = computed(() => useTenant().value)
export const getUserInfo = async () => {
  return await useHttp.get(`/default/me`)
}

export const login = async (data: any, curTenant: string = 'default') => {
  return await useHttp.post(`/${curTenant}/login`, data, {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    }
  })
}

export const register = async (data: any, curTenant: string = 'default') => {
  return await useHttp.post(`/${curTenant}/register`, data, {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    }
  })
}

export const getThirdLoginConfigs = async (currentTenant: string) => {
  return await useHttp.get(`/${currentTenant}/providers`)
}

export const getThirdLoginConfigByName = async (providerName: string, currentTenant = 'default') => {
  return await useHttp.get(`/${currentTenant}/providers/${providerName}`)
}

export const thirdLoginHandle = async (providerName: string, phone: string, currentTenant = 'default') => {
  return await useHttp.get(`/${currentTenant}/providers/${providerName}/login?phone=${phone}`)
}
export const thirdLoginHandleInfo = async (providerName: string, currentTenant = 'default', next: string,callback:string) => {
  return await useHttp.get(`/${currentTenant}/providers/${providerName}/login?next=${next}&callback=${callback}`)
}

export const thirdLogin = async (code:string,state:string) => {
  return await useHttp.get(`/login/providers/callback?code=${code}&state=${state}`)
}

export const phoneThirdLogin = async (providerName: string, params: { phone: string, code: string }, currentTenant = 'default') => {
  return await useHttp.get(`/${currentTenant}/logged-in/${providerName}`, params)
}

export const getProto = async (fileName: string, currentTenant = 'default') => {
  return await useHttp.get(`/${currentTenant}/login/proto/${fileName}`)
}


// 忘记密码
export const getResetPasswordToken = async (data: any, curTenant: string = 'default') => {
  return await useHttp.post(`/${curTenant}/reset/getResetPasswordToken`, data)
}

export const resetPassword = async (data: any, curTenant: string = 'default', token: string) => {
  return await useHttp.post(`/${curTenant}/reset/resetPassword`, data, {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      "Authorization": token
    }
  })
}

export const smsAvailable = async (curTenant: string) => {
  return await useHttp.get(`/${curTenant}/reset/smsAvailable`)
}
export const verifyResetPasswordRequest = async (data: any, curTenant: string = 'default') => {
  return await useHttp.post(`/${curTenant}/reset/verifyResetPasswordRequest`, data)
}