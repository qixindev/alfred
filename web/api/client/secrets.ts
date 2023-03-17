const tenant =  import.meta.env.VITE_APP_TENANT

export const getSecrets = async (clientId: number) => {
  return await useHttp.get(`/admin/${tenant}/clients/${clientId}/secrets`)
}

export const saveSecret = async (clientId: number, data: any) => {
  return await useHttp.post(`/admin/${tenant}/clients/${clientId}/secrets`, data)
}

export const updateSecret = async (clientId: number, id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant}/clients/${clientId}/secrets/${id}`, data)
}

export const delSecret = async (clientId: number, id: number) => {
  return await useHttp.delete(`/admin/${tenant}/clients/${clientId}/secret/${id}`)
}
