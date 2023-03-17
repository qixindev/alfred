const tenant = useRuntimeConfig().public.VITE_APP_TENANT

export const getActions = async (clientId: number, typeName: string) => {
  return await useHttp.get(`/${tenant}/iam/clients/${clientId}/types/${typeName}/actions`)
}

export const saveAction = async (clientId: number, typeName: string, data: any) => {
  return await useHttp.post(`/${tenant}/iam/clients/${clientId}/types/${typeName}/actions`, data)
}

export const updateAction = async (clientId: number, typeName: string, actionId: number, data: any) => {
  return await useHttp.put(`/${tenant}/iam/clients/${clientId}/types/${typeName}/actions/${actionId}`, data)
}

export const delAction = async (clientId: number, typeName: string, actionId: number) => {
  return await useHttp.delete(`/${tenant}/iam/clients/${clientId}/types/${typeName}/actions/${actionId}`)
}
