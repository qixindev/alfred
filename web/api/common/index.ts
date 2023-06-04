const tenant =  import.meta.env.VITE_APP_TENANT

/**
 * 获取clientUser列表
 */
export const getClientUsers = async (clientId: number) => {
  return await useHttp.get(`/admin/${tenant}/clients/${clientId}/users`)
}
/**
 * 获取用户列表
 */
export const getUser = async () => {
  return await useHttp.get(`/admin/tenants`)
}