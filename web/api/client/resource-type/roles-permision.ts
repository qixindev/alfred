// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getUsers = async (clientId: string, type: string, resource: string, role: string) => {
  return await useHttp.get(`/${tenant.value}/iam/clients/${clientId}/types/${type}/resources/${resource}/roles/${role}/users`)
}

export const saveUser = async (clientId: string, type: string, resource: string, role: string, data: any) => {
  return await useHttp.post(`/${tenant.value}/iam/clients/${clientId}/types/${type}/resources/${resource}/roles/${role}/users`, data)
}

export const delUser = async (clientId: string, type: string, resource: string, role: string, user: string) => {
  return await useHttp.delete(`/${tenant.value}/iam/clients/${clientId}/types/${type}/resources/${resource}/roles/${role}/users/${user}`)
}
export const addUser = async (clientId: string, type: string, role: string, data: any) => {
  return await useHttp.post(`/${tenant.value}/iam/clients/${clientId}/types/${type}/roles/${role}/auth`, data)
}