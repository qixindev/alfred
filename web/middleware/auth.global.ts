export default defineNuxtRouteMiddleware((to, from) => {
  const { authCode, state, code } = to.query

  // if (code) {
  //   useGetToken(to.query.code as string)
  // }

  if (authCode && state) {
    useThirdLogin(state as string, authCode as string)
  }
})