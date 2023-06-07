// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getRoles = async (clientId: number, type: string) => {
  return await useHttp.get(`/${tenant.value}/iam/clients/${clientId}/types/${type}/roles`)
}

export const saveRole = async (clientId: number, type: string, data: any) => {
  return await useHttp.post(`/${tenant.value}/iam/clients/${clientId}/types/${type}/roles`, data)
}

export const updateRole = async (clientId: number, type: number, role: string, data: any) => {
  return await useHttp.put(`/${tenant.value}/iam/clients/${clientId}/types/${type}/roles/${role}`, data)
}

export const delRole = async (clientId: number, type: string, role: string) => {
  return await useHttp.delete(`/${tenant.value}/iam/clients/${clientId}/types/${type}/roles/${role}`)
}
