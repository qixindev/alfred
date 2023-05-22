import { ElMessage } from 'element-plus'

export default defineNuxtRouteMiddleware((to, from) => {
  const route = useRoute()
  const auth = useCookie('QixinAuth')
    
  if (!auth.value) {
    return navigateTo("/login?from=" + route.fullPath);
  }
})