const tenant =  import.meta.env.VITE_APP_TENANT

export const getResources = async (clientId: number, type: string) => {
  return await useHttp.get(`/${tenant}/iam/clients/${clientId}/types/${type}/resources`)
}

export const saveResource = async (clientId: number, type: string, data: any) => {
  return await useHttp.post(`/${tenant}/iam/clients/${clientId}/types/${type}/resources`, data)
}

export const updateResource = async (clientId: number, type: number, resources: string, data: any) => {
  return await useHttp.put(`/${tenant}/iam/clients/${clientId}/types/${type}/resources/${resources}`, data)
}

export const delResource = async (clientId: number, type: string, resources: string) => {
  return await useHttp.delete(`/${tenant}/iam/clients/${clientId}/types/${type}/resources/${resources}`)
}
