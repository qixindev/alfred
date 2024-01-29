// const tenant =  import.meta.env.VITE_APP_TENANT
const tenant = computed(() => useTenant().value)
export const getUsers = async (data?: any) => {
  return await useHttp.get(`/admin/${tenant.value}/users`,data)
}

export const saveUser = async (data: any) => {
  return await useHttp.post(`/admin/${tenant.value}/users`, data)
}

export const updateUser = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant.value}/users/${id}`, data)
}

export const delUser = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant.value}/users/${id}`)
}
// 修改密码
export const passUser = async (id: number,data:any) => {
  return await useHttp.put(`/admin/${tenant.value}/users/${id}/password`,data)
}