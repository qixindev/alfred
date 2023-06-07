// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getRedirectUris = async (clientId: number) => {
  return await useHttp.get(`/admin/${tenant.value}/clients/${clientId}/redirect-uris`)
}

export const saveRedirectUri = async (clientId: number, data: any) => {
  return await useHttp.post(`/admin/${tenant.value}/clients/${clientId}/redirect-uris`, data)
}

export const updateRedirectUri = async (clientId: number, uriId: number, data: any) => {
  return await useHttp.put(`/admin/${tenant.value}/clients/${clientId}/redirect-uris/${uriId}`, data)
}

export const delRedirectUri = async (clientId: number, uriId: number) => {
  return await useHttp.delete(`/admin/${tenant.value}/clients/${clientId}/redirect-uris/${uriId}`)
}
