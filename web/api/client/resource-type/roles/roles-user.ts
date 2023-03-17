const tenant = useRuntimeConfig().public.VITE_APP_TENANT

export const getUsers = async (clientId: number, type: string, resource: string, role: string) => {
  return await useHttp.get(`/${tenant}/iam/clients/${clientId}/types/${type}/resources/${resource}/roles/${role}/users`)
}

export const saveUser = async (clientId: number, type: string, resource: string, role: string, data: any) => {
  return await useHttp.post(`/${tenant}/iam/clients/${clientId}/types/${type}/resources/${resource}/roles/${role}/users`, data)
}

export const updateUser = async (clientId: number, type: number, resource: string, role: string, data: any) => {
  return await useHttp.put(`/${tenant}/iam/clients/${clientId}/types/${type}/resources/${resource}/roles/${role}/users`, data)
}

export const delUser = async (clientId: number, type: string, resource: string, role: string, user: string) => {
  return await useHttp.delete(`/${tenant}/iam/clients/${clientId}/types/${type}/resources/${resource}/roles/${role}/users/${user}`)
}
