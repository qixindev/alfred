const tenant =  import.meta.env.VITE_APP_TENANT

export const getGroups = async () => {
  return await useHttp.get(`/admin/${tenant}/groups`)
}

export const saveGroup = async (data: any) => {
  return await useHttp.post(`/admin/${tenant}/groups`, data)
}

export const updateGroup = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant}/groups/${id}`, data)
}

export const delGroup = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant}/groups/${id}`)
}
