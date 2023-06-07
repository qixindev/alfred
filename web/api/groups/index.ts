// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getGroups = async () => {
  return await useHttp.get(`/admin/${tenant.value}/groups`)
}

export const saveGroup = async (data: any) => {
  return await useHttp.post(`/admin/${tenant.value}/groups`, data)
}

export const updateGroup = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant.value}/groups/${id}`, data)
}

export const delGroup = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant.value}/groups/${id}`)
}
