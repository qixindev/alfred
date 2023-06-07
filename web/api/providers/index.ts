// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getProvider = async (id: number) => {
  return await useHttp.get(`/admin/${tenant.value}/providers/${id}`)
}

export const getProviders = async () => {
  return await useHttp.get(`/admin/${tenant.value}/providers`)
}

export const saveProvider = async (data: any) => {
  return await useHttp.post(`/admin/${tenant.value}/providers`, data)
}

export const updateProvider = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant.value}/providers/${id}`, data)
}

export const delProvider = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant.value}/providers/${id}`)
}
