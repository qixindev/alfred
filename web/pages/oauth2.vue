<script setup lang="ts">
import { ElMessage } from "element-plus";

import { login, getThirdLoginConfigByName, thirdLogin,thirdLoginHandleInfo } from "~/api/user";

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

const phoneLoginHandle = async (params:any, phoneState:any) => {
  let { redirect_uri, client_id,state: tenant } = route.query;
  await thirdLogin(params.code,phoneState);
  navigateTo(
    `${location.origin}${VITE_APP_BASE_API}/${tenant}/oauth2/auth?client_id=${client_id}&scope=profileOpenId&response_type=code&redirect_uri=${redirect_uri}`,
    { external: true }
  );
};

const thirdLoginInfo = async (thirdInfo: any) => {
  const query = route.query;
  const redirect_uri = location.origin + `/redirect`;
  const res=await thirdLoginHandleInfo(thirdInfo.name,query.state,query.redirect_uri,redirect_uri)
  navigateTo(res.location,{ external: true })
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
      @thirdLoginHandle="thirdLoginInfo"
    />
    <Login
      v-else
      @accountLoginHandle="accountLoginHandle"
      @phoneLoginHandle="phoneLoginHandle"
      @thirdLoginHandle="thirdLoginInfo"
    />
  </div>
</template>

<style scoped lang="scss">
.wrap {
  width: 100vw;
  height: 100vh;
}
</style>
