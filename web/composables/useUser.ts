//用户状态信息模块
import { getUserInfo, thirdLogin } from '~/api/user';

export interface User {
  username: string
}
// export interface Tenant {
//   // name: string
// }
export interface Path {
  name: string,
  path:string,
  list: Array<SelectOption>
}
interface SelectOption {
  label: string,
  name: string,
  path: string
}

export const useUser = () => useState<User | undefined>("user", () => undefined);
export const useTenant = () => useState<String>("tenant", () => localStorage.getItem('tenantValue') as string);
const VITE_APP_BASE_API = import.meta.env.VITE_APP_BASE_API
/**
 *  获取用户信息（昵称、头像、角色集合、权限集合）
 */
export async function useGetUserInfo() {
  getUserInfo().then((res: any) => {
    const user = useUser();
    user.value = res
  })
}

/**
 *  注销
 */
export async function useLogout() {
  useRemoveToken()
  navigateTo('/login')
}

/**
 * 获取token
 */
export async function useGetToken(code: string) {
  const params = {
    client_id: '1',
    code: code,
    redirect_uri: 'http://10.1.0.135:3002',
    grant_type: 'authorization_code',
    client_secret: 'abcdefg'
  }
  // getToken(params).then((res: any) => {
  //   const token = useCookie('token')
  //   token.value = res.access_token
  //   useGetUserInfo()
  //   // 处理登录成功后页面跳转
  //   const route = useRoute()
  //   navigateTo(route.query.from as string || '/', { replace: true })
  // })
}

/**
 * 第三方登录
 */
export async function useThirdLogin(state: string, code: string) {
  // 处理cookie冲突
  useRemoveToken()
  const route = useRoute()

  if (route.query.tenant) {
    const res = await thirdLogin(route.query.tenant as string, code, state)
    const { clientId, redirect, tenant } = res as any
    navigateTo(`${location.origin}${VITE_APP_BASE_API}/${tenant}/oauth2/auth?client_id=${clientId}&scope=profileOpenId&response_type=code&redirect_uri=${redirect}`, { external: true })
  } else {
    await thirdLogin("default",code,state)
    await useGetUserInfo()
    // code使用完后删除url参数
    navigateTo(route.query.from as string || '/', { replace: true })
  }
}

/**
 * 清除本地数据
 */
export async function useRemoveToken() {
  // 清除cookie
  const auth = useCookie('QixinAuth')
  auth.value = null
  localStorage.removeItem('tenantValue')
  const tenant = useTenant()
  tenant.value = ''
  // 重置用户信息
  const user = useState('user')
  user.value = null
}