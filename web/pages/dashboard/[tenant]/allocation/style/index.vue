<script lang="ts" setup name="Users">
import {
  Monitor,
  Iphone,
  UploadFilled,
  Delete,
  FullScreen,
} from "@element-plus/icons-vue";
import { genFileId } from "element-plus";
import type { UploadInstance, UploadProps, UploadRawFile, FormRules } from "element-plus";
import { getThirdLoginConfigs, smsAvailable } from "~/api/user";
import dayjs from "dayjs";
import { ref, reactive, onMounted, onUnmounted } from "vue";
import { defineProps } from "vue";
import screenfull from "screenfull";
import { ElMessage } from "element-plus";
import { useRouter, useRoute } from "vue-router";
const tenant = computed(() => useTenant().value);
const route = useRoute();
const { state: tanent } = route.query as any;

let currentTenant =
  route.path.substring(0, 10) == "/dashboard" ? tenant.value : tanent ?? "default";
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
  "style-numTop",
  "style-numLeft",
  "equip",
]);
// 是否全屏
const isFullscreen = ref(false);
// 第三方账号登录
const thirdLoginTypes = ref([]);
// 短信验证
const codePass = ref(false);
// 忘记密码
const isPhone = ref(false);
const numTop = ref(null);
const numLeft = ref(null);
const loginSwitch = ref(false);
const regionSwitch = ref(false);
const passSwitch = ref(false);
const codeSwitch = ref(false);
const styleSwitch = ref("all");
const equip = ref("monitor");
const backG = ref("1");
const cssWrite = ref("");
const inputTitle = ref("");
const logoUpload = ref("");
let backgroundColor = ref("");

const getLoginConfig = async () => {
  const option = ["wecom", "dingtalk"];
  const data = await getThirdLoginConfigs(currentTenant);
  const thirdLoginList = data ? data.filter((item) => option.includes(item.type)) : "";
  thirdLoginTypes.value = thirdLoginList;
  if (data && data.find((item) => item.type === "sms")) {
    checkPhone();
    codePass.value = true;
  } else {
    codePass.value = false;
  }
};
getLoginConfig();
const checkPhone = async () => {
  const res = await smsAvailable(currentTenant);
  if (res) {
    isPhone.value = true;
  } else {
    isPhone.value = false;
  }
};
const cellLogin = () => {
  if (thirdLoginTypes.value.length != 0) {
    emits("style-login", loginSwitch.value);
  } else if (thirdLoginTypes.value.length == 0 && loginSwitch.value) {
    ElMessage({
      message: "还未配置第三方登录方式",
      type: "warning",
    });
  }
};
const cellRegion = () => {
  emits("style-region", regionSwitch.value);
};
const cellPass = () => {
  if (isPhone.value) {
    emits("style-pass", passSwitch.value);
  } else if (!isPhone.value && passSwitch.value) {
    ElMessage({
      message: "还未配置短信验证方式",
      type: "warning",
    });
  }
};
const cellCode = () => {
  if (codePass.value) {
    emits("style-code", codeSwitch.value);
  } else if (!codePass.value && codeSwitch.value) {
    ElMessage({
      message: "还未配置短信验证方式",
      type: "warning",
    });
  }
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
const numTopFn = (e) => {
  emits("style-numTop", numTop.value);
};
const numLeftFn = () => {
  emits("style-numLeft", numLeft.value);
};
const equipTab = () => {
  emits("equip", equip.value);
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
  allInfo: {
    type: Object,
    default: {},
  },
});
watch(
  () => props.allInfo,
  () => {
    inputTitle.value = props.allInfo.styleName;
    logoUpload.value = props.allInfo.styleLogo;
    backgroundColor.value = props.allInfo.styleBgcolor;
    cssWrite.value = props.allInfo.styleCss;
    loginSwitch.value =
      props.allInfo.styleLogin == undefined ? false : props.allInfo.styleLogin;
    regionSwitch.value =
      props.allInfo.styleRegion == undefined ? false : props.allInfo.styleRegion;
    passSwitch.value =
      props.allInfo.stylePass == undefined ? false : props.allInfo.stylePass;
    codeSwitch.value =
      props.allInfo.styleCode == undefined ? false : props.allInfo.styleCode;
    numTop.value = props.allInfo.styleNumTop;
    numLeft.value = props.allInfo.styleNumLeft;
    if (backgroundColor.value && backgroundColor.value.substring(0, 1) != "#") {
      backG.value = "2";
    } else {
      backG.value = "1";
    }
  },
  { immediate: true, deep: true }
);
// 切换事件
const preFn = () => {
  screenfull.toggle(document.getElementById("embedContainer"));
};
// 监听变化
const change = () => {
  isFullscreen.value = screenfull.isFullscreen;
};
// 设置侦听器
onMounted(() => {
  screenfull.on("change", change);
});

// 删除侦听器
onUnmounted(() => {
  screenfull.off("change", change);
});
</script>
<template>
  <div class="content">
    <el-radio-group v-model="styleSwitch" size="large" ref="myImg">
      <el-radio-button label="all">整体样式</el-radio-button>
      <el-radio-button label="convention">常规登录</el-radio-button>
    </el-radio-group>
  </div>
  <div id="wrap">
    <div class="centerMain">
      <div class="changeEquip" style="margin: 10px 0 0 10px">
        <el-tabs
          class="demo-tabs"
          tab-position="left"
          v-model="equip"
          @tab-change="equipTab"
        >
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
        <el-button
          @click="preFn"
          style="margin-top: 10px; box-shadow: 2px 5px 12px rgb(0 0 0/0.2)"
          v-if="equip == 'monitor'"
        >
          <el-icon :size="20"><FullScreen /></el-icon>
        </el-button>
      </div>
      <Login
        v-if="equip == 'monitor'"
        :numTop="numTop"
        :numLeft="numLeft"
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
        id="embedContainer"
        ref="scrollBox"
      />

      <div
        v-else
        style="
          background: #fff;
          display: flex;
          justify-content: center;
          background-color: #f7f8fa;
          height: 100vh;
        "
      >
        <div class="iphoneBorder">
          <LoginIphone
            :numTop="numTop"
            :numLeft="numLeft"
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
            style="border-radius: 30px"
          />
        </div>
      </div>
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
        <p class="bg" v-if="equip == 'monitor'">登录框位置</p>

        <div style="margin: 10px 10px 0 20px" v-if="equip == 'monitor'">
          top<el-input-number
            v-model="numTop"
            style="margin-left: 20px"
            @change="numTopFn"
            @blur="numTopFn"
            :min="-100"
            :max="100"
          />
        </div>
        <div style="margin: 10px 10px 0 20px" v-if="equip == 'monitor'">
          left<el-input-number
            v-model="numLeft"
            style="margin-left: 20px"
            @change="numLeftFn"
            @blur="numLeftFn"
          />
        </div>
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
        <el-switch
          v-model="regionSwitch"
          style="float: right"
          @change="cellRegion"
          disabled
        />
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
#wrap {
  width: 100%;
  height: 100vh;
  overflow-y: auto;
  display: flex;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE10+ */
  &::-webkit-scrollbar {
    display: none; /* ChromeSafari */
  }
}

.centerMain {
  height: 100vh;
  width: 75%;
  position: relative;
  .iphoneBorder {
    margin-top: 8vh;
    width: 375px;
    height: 670px;
    background: #fff;
    border: 10px solid #000;
    border-radius: 40px;
  }
}
.changeEquip {
  position: absolute;
  top: 0;
  left: 0;
  width: 58px;
  z-index: 999;
  .demo-tabs {
    background-color: #fff;
    border: 0;
    border-radius: 5px;
    box-shadow: 2px 5px 12px rgb(0 0 0/0.2);
  }
}
.allmain,
.conMain {
  flex: 1;
  height: 100vh;
  background-color: white;
  border-top: 1px solid #eee;
  overflow-y: auto;
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
