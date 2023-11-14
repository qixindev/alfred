<script lang="ts" setup name="Users">
import type { TabsPaneContext } from "element-plus";
import { ElMessage } from "element-plus";
import { ref } from "vue";
import Style from "./style/index.vue";
import Energy from "./energy/index.vue";
import { getEnergy, putEnergy } from "~/api/energy";
import { getRedirectUris } from "~/api/client/redirect-uri";
import { useRouter } from "vue-router";
const router = useRouter();
const tenant = computed(() => useTenant().value);
const activeName = ref("first");
const top = ref([]);
const bottom = ref([]);
const styleBgcolor = ref<String>("");
const styleLogo = ref("");
const styleName = ref("");
const styleCss = ref("");
const styleLogin = ref(true);
const styleRegion = ref(true);
const stylePass = ref(true);
const styleCode = ref(true);
const styleNumTop = ref(null);
const styleNumLeft = ref(null);
const getInfo = () => {
  getEnergy().then((res: any) => {
    styleName.value = res.styleName;
    styleLogo.value = res.styleLogo;
    styleBgcolor.value = res.styleBgcolor;
    styleCss.value = res.styleCss;
    styleLogin.value = res.styleLogin == undefined ? true : res.styleLogin;
    styleRegion.value = res.styleRegion == undefined ? true : res.styleRegion;
    stylePass.value = res.stylePass == undefined ? true : res.stylePass;
    styleCode.value = res.styleCode == undefined ? true : res.styleCode;
    styleNumTop.value = res.styleNumTop;
    styleNumLeft.value = res.styleNumLeft;
    bottom.value = [...res.bottom];
  });
};
getInfo();
const stylelogin = (value) => {
  styleLogin.value = value;
};
const styleregion = (value) => {
  styleRegion.value = value;
};

const stylepass = (value) => {
  stylePass.value = value;
};
const stylecode = (value) => {
  styleCode.value = value;
};
const styleBg = (value) => {
  styleBgcolor.value = value;
};
const stylelogo = (value) => {
  styleLogo.value = value;
};
const stylename = (value) => {
  styleName.value = value;
};
const stylecss = (value) => {
  styleCss.value = value;
};
const stylenumleft = (value) => {
  styleNumLeft.value = value;
};
const stylenumtop = (value) => {
  styleNumTop.value = value;
};
const zCf = (value) => {
  console.log(value, "value");

  bottom.value = value;
};
const zCp = (value) => {
  top.value = value;
};

const submit = () => {
  putEnergy({
    bottom,
    styleLogin,
    styleBgcolor,
    styleLogo,
    styleName,
    styleCss,
    styleRegion,
    stylePass,
    styleCode,
    styleNumTop,
    styleNumLeft,
  }).finally(() => {
    ElMessage({
      message: "保存成功",
      type: "success",
    });
  });
};
function getList() {
  getRedirectUris("default").then((res: any) => {
    if (res.length != 0) {
      navigateTo(
        `${location.origin}/oauth2/?redirect_uri=${res[0].redirectUri}&response_type=code&client_id=default&scope=openid&prompt=consent&state=${tenant.value}`,
        { external: true }
      );
    } else {
      ElMessage({
        message: "体验失败",
        type: "danger",
      });
    }
  });
}
</script>
<template>
  <h3>
    全局登录配置
    <el-button type="primary" style="float: right; margin-left: 10px" @click="getList"
      >体验登录</el-button
    >
    <el-button type="primary" style="float: right" @click="submit">保存</el-button>
  </h3>
  <h6 style="margin: 10px 0 20px 20px; font-weight: normal">
    支持自定义一个美观漂亮的登录页面
  </h6>
  <el-tabs v-model="activeName" class="demo-tabs">
    <el-tab-pane label="样式配置" name="first"
      ><Style
        @style-login="stylelogin"
        @style-region="styleregion"
        @style-pass="stylepass"
        @style-code="stylecode"
        @style-bgColor="styleBg"
        @style-logo="stylelogo"
        @style-name="stylename"
        @style-css="stylecss"
        @style-numLeft="stylenumleft"
        @style-numTop="stylenumtop"
        :bottom="bottom"
        :top="top"
    /></el-tab-pane>
    <el-tab-pane label="功能配置" name="second"
      ><Energy @child-click="zCf" @child-primary="zCp"
    /></el-tab-pane>
  </el-tabs>
</template>
<style scoped lang="scss"></style>
