const tenant = useRuntimeConfig().public.VITE_APP_TENANT

export const getTypes = async (clientId: number) => {
  return await useHttp.get(`/${tenant}/iam/clients/${clientId}/types`)
}

export const saveType = async (clientId: number, data: any) => {
  return await useHttp.post(`/${tenant}/iam/clients/${clientId}/types`, data)
}

export const updateType = async (clientId: number, typeId: number, data: any) => {
  return await useHttp.put(`/${tenant}/iam/clients/${clientId}/types/${typeId}`, data)
}

export const delType = async (clientId: number, typeId: string) => {
  return await useHttp.delete(`/${tenant}/iam/clients/${clientId}/types/${typeId}`)
}
