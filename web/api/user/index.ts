const tenant =  import.meta.env.VITE_APP_TENANT

export const getUserInfo = async () => {
  return await useHttp.get(`/${tenant}/me`)
}

export const login = async (data: any) => {
  return await useHttp.post(`/${tenant}/login`, data, {
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    }
  })
}

export const register = async (data: any) => {
  return await useHttp.post(`/${tenant}/register`, data, {
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

export const getThirdLoginConfig = async () => {
  return await useHttp.get(`/${tenant}/providers`)
}

export const thirdLogin = async (providerName: string, data: any) => {
  return await useHttp.get(`/${tenant}/logged-in/${providerName}`, data)
}

