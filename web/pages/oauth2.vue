<script setup lang="ts">
import { ElMessage } from "element-plus";

import { login, getThirdLoginConfigByName, phoneThirdLogin } from "~/api/user";

const VITE_APP_BASE_API = import.meta.env.VITE_APP_BASE_API;
const route = useRoute();
const isPhone = ref(true);
const _isMobile = () => {
  let flag = navigator.userAgent.match(
    /(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i
  );
  return flag;
};
if (_isMobile()) {
  isPhone.value = true;
} else {
  isPhone.value = false;
}
const accountLoginHandle = (formData: any) => {
  let { redirect_uri, client_id, state: tenant } = route.query;
  login(formData, tenant as string).then((res) => {
    if (res == 10000) {
      ElMessage({
        message: "账号或密码错误",
        type: "error",
      });
    } else {
      navigateTo(
        `${location.origin}${VITE_APP_BASE_API}/${tenant}/oauth2/auth?client_id=${client_id}&scope=profileOpenId&response_type=code&redirect_uri=${redirect_uri}`,
        { external: true }
      );
    }
  });
};

const phoneLoginHandle = async (phoneProvider: string, params: any) => {
  let { redirect_uri, client_id, state: tenant } = route.query;
  await phoneThirdLogin(phoneProvider, params, tenant as string);
  navigateTo(
    `${location.origin}${VITE_APP_BASE_API}/${tenant}/oauth2/auth?client_id=${client_id}&scope=profileOpenId&response_type=code&redirect_uri=${redirect_uri}`,
    { external: true }
  );
};

const thirdLogin = async (thirdInfo: any) => {
  const query = route.query;
  const redirect_uri = location.origin + "/redirect";
  const config = await getThirdLoginConfigByName(thirdInfo.name, query.state as string);
  const params = {
    redirect_uri: query.redirect_uri,
    type: thirdInfo.name,
    client_id: query.client_id,
    tenant: query.state,
  };

  switch (thirdInfo.type) {
    case "dingtalk":
      navigateTo(
        `https://login.dingtalk.com/oauth2/auth?redirect_uri=${redirect_uri}&response_type=code&client_id=${
          config.appKey
        }&scope=openid&prompt=consent&state=${encodeURI(JSON.stringify(params))}`,
        { external: true }
      );
      break;
    case "wecom":
      navigateTo(
        `https://login.work.weixin.qq.com/wwlogin/sso/login?appid=${
          config.corpId
        }&redirect_uri=${redirect_uri}&state=${encodeURI(
          JSON.stringify(params)
        )}&agentid=${config.agentId}`,
        { external: true }
      );
      break;

    default:
      break;
  }
};

definePageMeta({
  layout: false,
});
</script>

<template>
  <div class="wrap">
    <LoginIphone
      v-if="isPhone"
      @accountLoginHandle="accountLoginHandle"
      @phoneLoginHandle="phoneLoginHandle"
      @thirdLoginHandle="thirdLogin"
    />
    <Login
      v-else
      @accountLoginHandle="accountLoginHandle"
      @phoneLoginHandle="phoneLoginHandle"
      @thirdLoginHandle="thirdLogin"
    />
  </div>
</template>

<style scoped lang="scss">
.wrap {
  width: 100vw;
  height: 100vh;
}
</style>
