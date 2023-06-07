// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getDevices = async () => {
  return await useHttp.get(`/admin/${tenant.value}/devices`)
}

export const saveDevice = async (data: any) => {
  return await useHttp.post(`/admin/${tenant.value}/devices`, data)
}

export const updateDevice = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant.value}/devices/${id}`, data)
}

export const delDevice = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant.value}/devices/${id}`)
}
