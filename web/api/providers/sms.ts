const tenant = computed(() => useTenant().value)
export const getSmsById = async (id: number) => {
  return await useHttp.get(`/admin/${tenant.value}/sms/${id}`)
}

export const getSms = async () => {
  return await useHttp.get(`/admin/${tenant.value}/sms`)
}

export const saveSms = async (data: any) => {
  return await useHttp.post(`/admin/${tenant.value}/sms`, data)
}

export const updateSms = async (id: number, data: any) => {
  return await useHttp.put(`/admin/${tenant.value}/sms/${id}`, data)
}

export const delSms = async (id: number) => {
  return await useHttp.delete(`/admin/${tenant.value}/sms/${id}`)
}
