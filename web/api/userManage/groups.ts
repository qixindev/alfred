const tenant = useRuntimeConfig().public.VITE_APP_TENANT

export const getGroups = async (userId: number) => {
  return await useHttp.get(`/admin/${tenant}/users/${userId}/groups`)
}

export const saveGroup = async (userId: number, data: any) => {
  return await useHttp.post(`/admin/${tenant}/users/${userId}/groups`, data)
}

export const updateGroup = async (userId: number, groupId: number, data: any) => {
  return await useHttp.put(`/admin/${tenant}/users/${userId}/groups/${groupId}`, data)
}

export const delGroup = async (userId: number, groupId: number,) => {
  return await useHttp.delete(`/admin/${tenant}/users/${userId}/groups/${groupId}`)
}
