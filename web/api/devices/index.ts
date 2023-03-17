const tenant =  import.meta.env.VITE_APP_TENANT

export const getDevices = async () => {
  return await useHttp.get(`/admin/${tenant}/devices`)
}

export const saveDevice = async (data: any) => {
  return await useHttp.post(`/admin/${tenant}/devices`, data)
}

export const updateDevice = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant}/devices/${id}`, data)
}

export const delDevice = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant}/devices/${id}`)
}
