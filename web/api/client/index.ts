// const tenant =  import.meta.env.VITE_APP_TENANT

const tenant = computed(() => useTenant().value).value ||localStorage.getItem('tenantValue')
export const getClient = async () => {
  console.log(tenant,computed(() => useTenant().value),"接口tenant111getClient",localStorage.getItem('tenantValue'));
  console.log(import.meta)
  return await useHttp.get(`/admin/${tenant}/clients`)
}

export const saveClient = async (data: any) => {
  console.log(tenant,computed(() => useTenant().value),"接口tenant111saveClient",localStorage.getItem('tenantValue'));
  return await useHttp.post(`/admin/${tenant}/clients`, data)
}

export const updateClient = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant}/clients/${id}`, data)
}

export const delClient = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant}/clients/${id}`)
}

export const getSecret = async (clientId: number) => {
  return await useHttp.get(`/admin/${tenant}/clients/${clientId}/secret`)
}

export const setSecret = async (clientId: number, data: any) => {
  return await useHttp.post(`/admin/${tenant}/clients/${clientId}/secret`, data)
}
