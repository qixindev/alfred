const tenant = useRuntimeConfig().public.VITE_APP_TENANT

export const getRoles = async (clientId: number, type: string) => {
  return await useHttp.get(`/${tenant}/iam/clients/${clientId}/types/${type}/roles`)
}

export const saveRole = async (clientId: number, type: string, data: any) => {
  return await useHttp.post(`/${tenant}/iam/clients/${clientId}/types/${type}/roles`, data)
}

export const updateRole = async (clientId: number, type: number, role: string, data: any) => {
  return await useHttp.put(`/${tenant}/iam/clients/${clientId}/types/${type}/roles/${role}`, data)
}

export const delRole = async (clientId: number, type: string, role: string) => {
  return await useHttp.delete(`/${tenant}/iam/clients/${clientId}/types/${type}/roles/${role}`)
}
