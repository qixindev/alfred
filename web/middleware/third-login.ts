export default defineNuxtRouteMiddleware((to, from) => {
  const { authCode, state, code } = to.query

  if (authCode && state) {
    useThirdLogin(state as string, authCode as string)
  }

  if (code && state) {
    useThirdLogin(state as string, code as string)
  }
})