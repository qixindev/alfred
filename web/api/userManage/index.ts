const tenant = useRuntimeConfig().public.VITE_APP_TENANT

export const getUsers = async () => {
  return await useHttp.get(`/admin/${tenant}/users`)
}

export const saveUser = async (data: any) => {
  return await useHttp.post(`/admin/${tenant}/users`, data)
}

export const updateUser = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant}/users/${id}`, data)
}

export const delUser = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant}/users/${id}`)
}
