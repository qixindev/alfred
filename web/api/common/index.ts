const tenant = useRuntimeConfig().public.VITE_APP_TENANT

/**
 * 获取clientUser列表
 */
export const getClientUsers = async (clientId: number) => {
  return await useHttp.get(`/admin/${tenant}/clients/${clientId}/users`)
}