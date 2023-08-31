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
  console.log(curTenant)
  return await useHttp.post(`/${curTenant}/register`, data, {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    }
  })
}

export const getThirdLoginConfigs = async (currentTenant = 'default') => {
  return await useHttp.get(`/${currentTenant}/providers`)
}

export const getThirdLoginConfigByName = async (providerName: string, currentTenant = 'default') => {
  return await useHttp.get(`/${currentTenant}/providers/${providerName}`)
}

export const thirdLoginHandle = async (providerName: string, phone: string, currentTenant='default') => {
  return await useHttp.get(`/${currentTenant}/login/${providerName}?phone=${phone}`)
}

export const thirdLogin = async (providerName: string, data: any, currentTenant='default') => {
  return await useHttp.get(`/${currentTenant}/logged-in/${providerName}`, data)
}

export const phoneThirdLogin = async (providerName: string, params: {phone: string, code: string},  currentTenant='default') => {
  return await useHttp.get(`/${currentTenant}/logged-in/${providerName}`, params)
}


