// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getActions = async (clientId: number, typeName: string) => {
  return await useHttp.get(`/${tenant.value}/iam/clients/${clientId}/types/${typeName}/actions`)
}

export const saveAction = async (clientId: number, typeName: string, data: any) => {
  return await useHttp.post(`/${tenant.value}/iam/clients/${clientId}/types/${typeName}/actions`, data)
}

export const updateAction = async (clientId: number, typeName: string, actionId: number, data: any) => {
  return await useHttp.put(`/${tenant.value}/iam/clients/${clientId}/types/${typeName}/actions/${actionId}`, data)
}

export const delAction = async (clientId: number, typeName: string, actionId: number) => {
  return await useHttp.delete(`/${tenant.value}/iam/clients/${clientId}/types/${typeName}/actions/${actionId}`)
}
