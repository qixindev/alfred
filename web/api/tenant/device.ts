export const getDevices = async () => {
  return await useHttp.get(`/admin/tenants/devices`)
}

export const saveDevice = async (data: any) => {
  return await useHttp.post(`/admin/tenants/devices`, data)
}

export const updateDevice = async (id: number, data: any) => {
  return await useHttp.put(`/admin/tenants/devices/${id}`, data)
}

export const delDevice = async (id: number) => {
  return await useHttp.delete(`/admin/tenants/devices/${id}`)
}
