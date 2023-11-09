<script lang="ts" setup name="Users">
import { Monitor, Iphone, UploadFilled, Delete } from "@element-plus/icons-vue";
import { genFileId } from "element-plus";
import type { UploadInstance, UploadProps, UploadRawFile, FormRules } from "element-plus";
import Login from "@/components/Login.vue";
import dayjs from "dayjs";
import { ref, reactive } from "vue";
import { defineProps } from "vue";
import { getEnergy, putEnergy, getProto, putProto } from "~/api/energy";
const tenant = computed(() => useTenant().value);
const VITE_APP_BASE_API = import.meta.env.VITE_APP_BASE_API;
const emits = defineEmits([
  "style-bgColor",
  "style-logo",
  "style-name",
  "style-css",
  "style-login",
  "style-region",
  "style-pass",
  "style-code",
]);
const loginSwitch = ref(true);
const regionSwitch = ref(true);
const passSwitch = ref(true);
const codeSwitch = ref(true);
const styleSwitch = ref("all");
const equip = ref("monitor");
const backG = ref("1");
const cssWrite = ref("");
const inputTitle = ref("");
const logoUpload = ref("");
let backgroundColor = ref("");

const getInfo = () => {
  getEnergy()
    .then((res: any) => {
      inputTitle.value = res.styleName;
      logoUpload.value = res.styleLogo;
      backgroundColor.value = res.styleBgcolor;
      cssWrite.value = res.styleCss;
      loginSwitch.value = res.styleLogin;
      regionSwitch.value = res.styleRegion;
      passSwitch.value = res.stylePass;
      codeSwitch.value = res.styleCode;
    })
    .finally(() => {});
};
getInfo();
const cellLogin = () => {
  emits("style-login", loginSwitch.value);
};
const cellRegion = () => {
  emits("style-region", regionSwitch.value);
};
const cellPass = () => {
  emits("style-pass", passSwitch.value);
};
const cellCode = () => {
  emits("style-code", codeSwitch.value);
};
function changeColor() {
  emits("style-bgColor", backgroundColor.value);
}
function mainCss() {
  emits("style-css", cssWrite.value);
}
const uploadFile = (e) => {
  backgroundColor.value = e.url;
  emits("style-bgColor", backgroundColor.value);
};
const uploadlogo = (e) => {
  logoUpload.value = e.url;
  emits("style-logo", logoUpload.value);
};
const cellPri = () => {
  emits("style-name", inputTitle.value);
};

const props = defineProps({
  top: {
    type: Array,
    default: [],
  },
  bottom: {
    type: Array,
    default: [],
  },
});
</script>
<template>
  <div class="content">
    <el-radio-group v-model="styleSwitch" size="large">
      <el-radio-button label="all">整体样式</el-radio-button>
      <el-radio-button label="convention">常规登录</el-radio-button>
    </el-radio-group>
  </div>
  <div class="wrap">
    <div class="centerMain">
      <div class="changeEquip">
        <el-tabs class="demo-tabs" tab-position="left" v-model="equip">
          <el-tab-pane name="monitor">
            <template #label>
              <span class="custom-tabs-label">
                <el-icon :size="20"><Monitor /></el-icon>
              </span>
            </template>
          </el-tab-pane>
          <el-tab-pane name="iphone">
            <template #label>
              <span class="custom-tabs-label">
                <el-icon :size="20"><Iphone /></el-icon>
              </span>
            </template>
          </el-tab-pane>
        </el-tabs>
      </div>

      <Login
        :inputTitle="inputTitle"
        :cssWrite="cssWrite"
        :top="props.top"
        :bottom="props.bottom"
        :passSwitch="passSwitch"
        :loginSwitch="loginSwitch"
        :regionSwitch="regionSwitch"
        :codeSwitch="codeSwitch"
        :logoUpload="logoUpload"
        :backgroundColor="backgroundColor"
      ></Login>
    </div>

    <div class="allmain" v-if="styleSwitch == 'all'">
      <div class="top">
        <p class="bg">自定义背景</p>
        <aside style="font-size: 14px; color: #aeaaaa; margin: 10px 0 0 50px">
          登录页面展示的背景
        </aside>
        <div class="mb-2 flex items-center text-sm">
          <el-radio-group v-model="backG" class="ml-4">
            <el-radio label="1" size="large">纯色背景</el-radio>
            <el-radio label="2" size="large">图片背景</el-radio>
          </el-radio-group>
        </div>
        <div class="demo-color-block" v-if="backG == 1">
          <el-color-picker v-model="backgroundColor" @change="changeColor" size="large" />
        </div>
        <el-upload
          class="upload-demo"
          drag
          :action="`${VITE_APP_BASE_API}/admin/${tenant}/picture/background/upload`"
          :on-success="uploadFile"
          v-else
          method="put"
        >
          <el-icon class="el-icon--upload" size="large"><upload-filled /></el-icon>
          <div class="el-upload__text">拖动文件或者点击上传</div>
          <template #tip>
            <div class="el-upload__tip">jpg/png 文件不超过 2MB</div>
          </template>
        </el-upload>
        <div>
          <p class="bg">自定义LOGO</p>
          <el-upload
            class="upload-demo1"
            :action="`${VITE_APP_BASE_API}/admin/${tenant}/picture/logo/upload`"
            :on-success="uploadlogo"
            :limit="1"
            method="put"
            :show-file-list="false"
          >
            <template #trigger>
              <el-button type="primary">选择文件</el-button>
            </template>
          </el-upload>
        </div>
        <div style="display: flex">
          <p style="font-size: 20px; margin: 20px 20px 0px 20px">平台名称</p>
          <el-input
            v-model="inputTitle"
            placeholder="请输入名称"
            style="width: 50%; margin: 20px 20px 0px 0px"
            @change="cellPri"
          />
        </div>
        <p style="font-size: 20px; margin: 20px 20px 0px 20px">自定义CSS</p>
        <p style="margin: 10px 20px 20px 20px">
          <el-input
            v-model="cssWrite"
            :rows="3"
            type="textarea"
            @change="mainCss"
            placeholder="例：.login-boxL{background:blue!important;}"
          />
        </p>
      </div>
    </div>
    <div class="conMain" v-if="styleSwitch == 'convention'">
      <p class="bg">
        第三方账号登录
        <el-switch v-model="loginSwitch" style="float: right" @change="cellLogin" />
      </p>
      <p style="margin-left: 20px; color: #aeaaaa">对配置了第三方账号登录的应用生效</p>
      <p class="bg">
        注册账户
        <el-switch v-model="regionSwitch" style="float: right" @change="cellRegion" />
      </p>
      <p style="margin-left: 20px; color: #aeaaaa">对配置了注册账户的应用生效</p>
      <p class="bg">
        忘记密码
        <el-switch v-model="passSwitch" style="float: right" @change="cellPass" />
      </p>
      <p style="margin-left: 20px; color: #aeaaaa">对配置了密码登录的应用生效</p>
      <p class="bg">
        验证码登录
        <el-switch v-model="codeSwitch" style="float: right" @change="cellCode" />
      </p>
      <p style="margin-left: 20px; color: #aeaaaa">对配置了手机验证码登录的应用生效</p>
    </div>
  </div>
</template>
<style scoped lang="scss">
.content {
  margin-bottom: 10px;
}
.wrap {
  width: 100%;
  height: 68vh;
  overflow-y: auto;
  display: flex;
}

.centerMain {
  position: relative;
  height: 68vh;
  width: 75%;
  overflow: hidden;
}
.changeEquip {
  position: absolute;
  top: 0;
  left: 0;
}
.allmain,
.conMain {
  flex: 1;
  min-height: 68vh;
  background-color: white;
  border-top: 1px solid #eee;
  .text-sm {
    margin-left: 20px;
  }
  .demo-color-block {
    display: flex;
    margin-left: 30px;
  }

  .el-upload__tip {
    margin-left: 100px;
  }
  .bottomWord {
    width: 200px;
    margin-left: 80px;
  }
  .w-50 {
    width: 150px;
    margin: 2px;
  }
  .upload-demo1 {
    margin: 10px 20px 10px 20px;
  }
}
.bg {
  font-size: 20px;
  margin: 20px 20px 0px 20px;
}
</style>
