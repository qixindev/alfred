const tenant =  import.meta.env.VITE_APP_TENANT

export const getProviders = async () => {
  return await useHttp.get(`/admin/${tenant}/providers`)
}

export const saveProvider = async (data: any) => {
  return await useHttp.post(`/admin/${tenant}/providers`, data)
}

export const updateProvider = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant}/providers/${id}`, data)
}

export const delProvider = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant}/providers/${id}`)
}
