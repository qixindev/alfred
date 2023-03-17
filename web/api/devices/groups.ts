const tenant =  import.meta.env.VITE_APP_TENANT

export const getGroups = async (deviceId: number) => {
  return await useHttp.get(`/admin/${tenant}/devices/${deviceId}/groups`)
}

export const saveGroup = async (deviceId: number, data: any) => {
  return await useHttp.post(`/admin/${tenant}/devices/${deviceId}/groups`, data)
}

export const updateGroup = async (deviceId: number, groupId: number, data: any) => {
  return await useHttp.put(`/admin/${tenant}/devices/${deviceId}/groups/${groupId}`, data)
}

export const delGroup = async (deviceId: number, groupId: number) => {
  return await useHttp.delete(`/admin/${tenant}/devices/${deviceId}/groups/${groupId}`)
}
