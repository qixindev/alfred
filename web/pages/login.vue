<script setup lang="ts">
import { ElMessage } from "element-plus";

import { login, getThirdLoginConfigByName, phoneThirdLogin,thirdLoginHandleInfo } from "~/api/user";

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
  login(formData).then((res) => {
    if (res == 10000) {
      ElMessage({
        message: "账号或密码错误",
        type: "error",
      });
    } else {
      navigateTo((route.query.from as string) || "/", { replace: true });
    }
  });
};

const phoneLoginHandle = async (phoneProvider: string, params: any) => {
  await phoneThirdLogin(phoneProvider, params);
  navigateTo((route.query.from as string) || "/", { replace: true });
};

const thirdLogin = async (thirdInfo: any) => {
  const redirect_uri = location.origin + "/redirect";
  
  const res = await thirdLoginHandleInfo(thirdInfo.name,"default",location.origin,redirect_uri)
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
