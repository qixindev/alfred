import { ElMessage } from 'element-plus'

export default defineNuxtRouteMiddleware((to, from) => {
  const route = useRoute()
  const auth = useCookie('QixinAuth')
    
  if (!auth.value) {
    ElMessage.warning('请先登录')
    return navigateTo("/login?from=" + route.fullPath);
  }
})