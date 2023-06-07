// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getRoleActions = async (clientId: number, type: string, roleId: number) => {
  return await useHttp.get(`/${tenant.value}/iam/clients/${clientId}/types/${type}/roles/${roleId}/actions`)
}

export const saveRoleAction = async (clientId: number, type: string, roleId: number, data: any) => {
  return await useHttp.post(`/${tenant.value}/iam/clients/${clientId}/types/${type}/roles/${roleId}/actions`, data)
}

export const updateRoleAction = async (clientId: number, type: number, roleId: number, actionId: number, data: any) => {
  return await useHttp.put(`/${tenant.value}/iam/clients/${clientId}/types/${type}/roles/${roleId}/actions/${actionId}`, data)
}

export const delRoleAction = async (clientId: number, type: string, roleId: number, actionId: number,) => {
  return await useHttp.delete(`/${tenant.value}/iam/clients/${clientId}/types/${type}/roles/${roleId}/actions/${actionId}`)
}
