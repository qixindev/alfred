const tenant = useRuntimeConfig().public.VITE_APP_TENANT

export const getTenants = async () => {
  return await useHttp.get(`/admin/tenants`)
}

export const saveTenant = async (data: any) => {
  return await useHttp.post(`/admin/tenants`, data)
}

export const updateTenant = async (tenantId: number, data: any) => {
  return await useHttp.put(`/admin/tenants/${tenantId}`, data)
}

export const delTenant = async (tenantId: number) => {
  return await useHttp.delete(`/admin/tenants/${tenantId}`)
}
