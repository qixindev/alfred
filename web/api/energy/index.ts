const tenant = computed(() => useTenant().value)
export const getEnergy = async () => {
    return await useHttp.get(`/admin/${tenant.value}/page/login`)
}

export const putEnergy = async (data: any) => {
    return await useHttp.put(`/admin/${tenant.value}/page/login`, data)
}

export const getProto = async () => {
    return await useHttp.get(`/admin/${tenant.value}/proto`)
}

export const putProto = async (data: any) => {
    return await useHttp.put(`/admin/${tenant.value}/proto`, data)
}
