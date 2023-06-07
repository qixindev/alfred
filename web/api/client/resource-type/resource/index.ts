// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getResources = async (clientId: number, type: string) => {
  return await useHttp.get(`/${tenant.value}/iam/clients/${clientId}/types/${type}/resources`)
}

export const saveResource = async (clientId: number, type: string, data: any) => {
  return await useHttp.post(`/${tenant.value}/iam/clients/${clientId}/types/${type}/resources`, data)
}

export const updateResource = async (clientId: number, type: number, resources: string, data: any) => {
  return await useHttp.put(`/${tenant.value}/iam/clients/${clientId}/types/${type}/resources/${resources}`, data)
}

export const delResource = async (clientId: number, type: string, resources: string) => {
  return await useHttp.delete(`/${tenant.value}/iam/clients/${clientId}/types/${type}/resources/${resources}`)
}
