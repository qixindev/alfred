
export const getEnergy = async (currentTenant: string) => {
    return await useHttp.get(`/admin/${currentTenant}/page/login`)
}

export const putEnergy = async (currentTenant: string, data: any) => {
    return await useHttp.put(`/admin/${currentTenant}/page/login`, data)
}

export const getProto = async (currentTenant: string) => {
    return await useHttp.get(`/admin/${currentTenant}/proto`)
}

export const putProto = async (currentTenant: string, data: any) => {
    return await useHttp.put(`/admin/${currentTenant}/proto`, data)
}
