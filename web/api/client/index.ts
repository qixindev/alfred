// const tenant.value =  import.meta.env.VITE_APP_tenant.value

const tenant = computed(() => useTenant().value)
export const getClient = async () => {
  return await useHttp.get(`/admin/${tenant.value}/clients`)
}

export const saveClient = async (data: any) => {
  return await useHttp.post(`/admin/${tenant.value}/clients`, data)
}

export const updateClient = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant.value}/clients/${id}`, data)
}

export const delClient = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant.value}/clients/${id}`)
}

export const getSecret = async (clientId: number) => {
  return await useHttp.get(`/admin/${tenant.value}/clients/${clientId}/secret`)
}

export const setSecret = async (clientId: number, data: any) => {
  return await useHttp.post(`/admin/${tenant.value}/clients/${clientId}/secret`, data)
}
