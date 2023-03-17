const tenant = useRuntimeConfig().public.VITE_APP_TENANT

export const getClient = async () => {
  return await useHttp.get(`/admin/${tenant}/clients`)
}

export const saveClient = async (data: any) => {
  return await useHttp.post(`/admin/${tenant}/clients`, data)
}

export const updateClient = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant}/clients/${id}`, data)
}

export const delClient = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant}/clients/${id}`)
}

export const getSecret = async (clientId: number) => {
  return await useHttp.get(`/admin/${tenant}/clients/${clientId}/secret`)
}

export const setSecret = async (clientId: number, data: any) => {
  return await useHttp.post(`/admin/${tenant}/clients/${clientId}/secret`, data)
}
