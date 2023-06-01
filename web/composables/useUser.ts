//用户状态信息模块
import { getUserInfo, getToken, thirdLogin } from '~/api/user';

export interface User {
  username: string
}

export const useUser = () => useState<User | undefined>("user", () => undefined);
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
  getToken(params).then((res: any) => {
    const token = useCookie('token')
    token.value = res.access_token
    useGetUserInfo()
    // 处理登录成功后页面跳转
    const route = useRoute()
    navigateTo(route.query.from as string || '/', { replace: true })
  })
}

/**
 * 获取第三方登录配置
 */
export async function getThirdLoginConfig(type: string, code: string) {
  const params = {
    code: code,
  }
  await thirdLogin(type, params)
  await useGetUserInfo()
  // code使用完后删除url参数
  const route = useRoute()  
  navigateTo(route.query.from as string || '/', { replace: true })
}

/**
 * 第三方登录
 */
export async function useThirdLogin(state: string, code: string) {
  const params = {
    code: code,
  }
  if (isJsonString(state)) {
    const {redirect_uri, client_id, type, tenant } = JSON.parse(state)
    await thirdLogin(type, params, tenant)
    navigateTo(`${location.origin}${VITE_APP_BASE_API}/${tenant}/oauth2/auth?client_id=${client_id}&scope=profileOpenId&response_type=code&redirect_uri=${redirect_uri}`,{ external: true })
  }else {
    const params = {
      code: code,
    }
    await thirdLogin(state, params)
    await useGetUserInfo()
    // code使用完后删除url参数
    const route = useRoute()  
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

  // const token = useCookie('token')
  // token.value = null

  // 重置用户信息
  const user = useState('user')
  user.value = null
}