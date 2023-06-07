// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getTypes = async (clientId: number) => {
  return await useHttp.get(`/${tenant.value}/iam/clients/${clientId}/types`)
}

export const saveType = async (clientId: number, data: any) => {
  return await useHttp.post(`/${tenant.value}/iam/clients/${clientId}/types`, data)
}

export const updateType = async (clientId: number, typeId: number, data: any) => {
  return await useHttp.put(`/${tenant.value}/iam/clients/${clientId}/types/${typeId}`, data)
}

export const delType = async (clientId: number, typeId: string) => {
  return await useHttp.delete(`/${tenant.value}/iam/clients/${clientId}/types/${typeId}`)
}
