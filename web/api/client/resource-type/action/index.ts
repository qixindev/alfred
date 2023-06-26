// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getActions = async (clientId: number, typeId: string) => {
  return await useHttp.get(`/${tenant.value}/iam/clients/${clientId}/types/${typeId}/actions`)
}

export const saveAction = async (clientId: number, typeId: string, data: any) => {
  console.log(typeId,"typeIdurl");
  
  return await useHttp.post(`/${tenant.value}/iam/clients/${clientId}/types/${typeId}/actions`, data)
}

export const updateAction = async (clientId: number, typeId: string, actionId: number, data: any) => {
  return await useHttp.put(`/${tenant.value}/iam/clients/${clientId}/types/${typeId}/actions/${actionId}`, data)
}

export const delAction = async (clientId: number, typeId: string, actionId: number) => {
  return await useHttp.delete(`/${tenant.value}/iam/clients/${clientId}/types/${typeId}/actions/${actionId}`)
}
