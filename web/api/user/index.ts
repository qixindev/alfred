const tenant =  import.meta.env.VITE_APP_TENANT as string

export const getUserInfo = async () => {
  return await useHttp.get(`/${tenant}/me`)
}

export const login = async (data: any, curTenant: string = tenant) => {
  return await useHttp.post(`/${curTenant}/login`, data, {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    }
  })
}

export const register = async (data: any, curTenant: string = tenant) => {
  console.log(curTenant)
  return await useHttp.post(`/${curTenant}/register`, data, {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    }
  })
}

export const auth = async (data: any) => {
  return await useHttp.get(`/${tenant}/oauth2/auth`, data)
}

export const getToken = async (data: any) => {
  return await useHttp.get(`/${tenant}/oauth2/token`, data)
}

export const jwks = async () => {
  return await useHttp.get(`/${tenant}/.well-known/jwks.json`)
}

export const getProviders = async () => {
  return await useHttp.get(`/${tenant}/providers`)
}

export const getProvidersById = async (providers: string) => {
  return await useHttp.get(`/${tenant}/providers/${providers}`)
}

export const getThirdLoginConfigs = async (currentTenant = tenant) => {
  return await useHttp.get(`/${currentTenant}/providers`)
}

export const getThirdLoginConfigByName = async (providerName: string, currentTenant = tenant) => {
  return await useHttp.get(`/${currentTenant}/providers/${providerName}`)
}

export const thirdLogin = async (providerName: string, data: any, currentTenant=tenant) => {
  return await useHttp.get(`/${currentTenant}/logged-in/${providerName}`, data)
}

