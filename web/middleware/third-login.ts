export default defineNuxtRouteMiddleware((to, from) => {
  const { authCode, state, code } = to.query

  if ((authCode || code) && state) {
    useThirdLogin(state as string, (authCode || code) as string)
  }
})