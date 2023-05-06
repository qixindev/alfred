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

export const getThirdLoginConfigs = async () => {
  return await useHttp.get(`/${tenant}/providers`)
}

export const getThirdLoginConfigByName = async (providerName: string) => {
  return await useHttp.get(`/${tenant}/providers/${providerName}`)
}

export const thirdLogin = async (providerName: string, data: any) => {
  return await useHttp.get(`/${tenant}/logged-in/${providerName}`, data)
}

