<script setup lang="ts">
import  { ElMessage } from 'element-plus'

import { login, getThirdLoginConfigByName, phoneThirdLogin} from '~/api/user';

const route = useRoute()

const accountLoginHandle = (formData: any) => {
  login(formData).then(res => {
    if (res == 10000) {
      ElMessage({
        message: '账号或密码错误',
        type: 'error'
      })
    } else {
      navigateTo(route.query.from as string || '/', { replace: true })
    }
  })
}

const phoneLoginHandle = async (phoneProvider: string, params: any) => {
  await phoneThirdLogin(phoneProvider, params)
  navigateTo(route.query.from as string || '/', { replace: true })
}

const thirdLogin = async (thirdInfo: any) => {
  const config = await getThirdLoginConfigByName(thirdInfo.name)
  const redirect_uri  = location.origin + '/redirect'
  switch (thirdInfo.type) {
    case 'dingtalk':
      navigateTo(`https://login.dingtalk.com/oauth2/auth?redirect_uri=${redirect_uri}&response_type=code&client_id=${config.appKey}&scope=openid&prompt=consent&state=${thirdInfo.name}`, { external: true})
      break;
      case 'wecom':
      navigateTo(`https://login.work.weixin.qq.com/wwlogin/sso/login?appid=${config.corpId}&redirect_uri=${redirect_uri}&state=${encodeURI(JSON.stringify(thirdInfo))}&agentid=${config.agentId}`, { external: true})
      break;
    default:
      break;
  }
}

definePageMeta({
  layout: false
})

</script>

<template>
  <Login
    @accountLoginHandle="accountLoginHandle"
    @phoneLoginHandle="phoneLoginHandle"
    @thirdLoginHandle="thirdLogin"
  />
</template>

<style scoped lang="scss">
</style>