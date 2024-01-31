<script setup lang="ts">
import type { FormInstance, FormRules } from "element-plus";
import { User, Lock, Iphone } from "@element-plus/icons-vue";

import { ElMessage } from "element-plus";
import CountdownButton from "@/components/CountdownButton/index.vue";
import { defineProps } from "vue";
import { getEnergy, getProto } from "~/api/energy";
import {
  getThirdLoginConfigs,
  getThirdLoginConfigByName,
  thirdLoginHandle,
  smsAvailable,
} from "~/api/user";
import { useRoute, useRouter } from "vue-router";
interface ThirdLoginType {
  id: number;
  name: string;
  type: string;
}
const tenant = computed(() => useTenant().value);

const route = useRoute();
const { state: tanent } = route.query as any;

let currentTenant =
  route.path.substring(0, 10) == "/dashboard" ? tenant.value : tanent ?? "default";

const router = useRoute();
const routerPath = useRouter();
const info = ref({});
const bottomTitle = ref([]);
const newPrimaryWord = ref([]);
const newTop = ref([]);
const login_top = ref(0);
const login_left = ref(0);
const navigatePage = ref(false);
const getInfo = () => {
  navigatePage.value = true;
  getEnergy(currentTenant).then((res: any) => {
    //  解决 is not iterable
    navigatePage.value = false;
    if (JSON.stringify(res) !== "{}") {
      info.value = { ...res };
      bottomTitle.value = [...res.bottom];
      login_top.value = res.styleNumTop;
      login_left.value = res.styleNumLeft;
    }
  });
  getProto(currentTenant).then((res: any) => {
    newPrimaryWord.value = res.filter((item: any) => {
      return item.loginSwitch;
    });
  });
};
getInfo();
const protocol = ref(false);
const props = defineProps({
  numTop: {
    type: Number,
    default: null,
  },
  numLeft: {
    type: Number,
    default: null,
  },
  loginSwitch: {
    type: Boolean,
    default: false,
  },
  regionSwitch: {
    type: Boolean,
    default: false,
  },
  passSwitch: {
    type: Boolean,
    default: false,
  },
  codeSwitch: {
    type: Boolean,
    default: false,
  },
  inputTitle: {
    type: String,
    default: "",
  },
  cssWrite: {
    type: String,
    default: "",
  },
  backgroundColor: {
    type: String,
    default: "",
  },
  logoUpload: {
    type: String,
    default: "",
  },
  top: {
    type: Array,
    default: [],
  },
  bottom: {
    type: Array,
    default: [],
  },
  accountLoad: {
    type: Boolean,
    default: false,
  },
});
const style = document.createElement("style");
watch(
  () => [info.value.styleCss, props.cssWrite],
  () => {
    const newCSS = `
    ${props.cssWrite ? props.cssWrite : info && info.value && info.value.styleCss}
`;
    style.textContent = newCSS;
    const reStyle = document.body.querySelector("style");
    if (reStyle) {
      document.body.removeChild(reStyle);
    }
    document.body.appendChild(style);
  },
  { immediate: true, deep: true }
);
watch(
  () => props.top,
  () => {
    newTop.value = props.top.filter((item: any) => {
      return item.loginSwitch;
    });
    if (newTop.value.length == 0) {
      newPrimaryWord.value = [];
    }
  },
  { immediate: true, deep: true }
);
watch(
  () => props.bottom,
  () => {
    if (props.bottom.length == 0) {
      bottomTitle.value = [];
    }
  },
  { immediate: true, deep: true }
);
watch(
  () => [props.numTop, props.numLeft],
  () => {
    if (!props.numTop) {
      login_top.value = 0;
    }
    if (!props.numLeft) {
      login_left.value = 0;
    }
  },
  { immediate: true, deep: true }
);
const emit = defineEmits(["accountLoginHandle", "phoneLoginHandle", "thirdLoginHandle"]);

const phoneForm = reactive({
  phone: "",
  code: "",
});

const accountForm = reactive({
  login: "",
  password: "",
});
// ？？？？？？？？
const hasRegister = computed(() => {
  return route.query?.platform === "tenant";
});

let thirdLoginTypes = ref<ThirdLoginType[]>([]);
let thirdLoginTypesLength = ref();
const phoneRuleFormRef = ref<FormInstance>();
const accountRuleFormRef = ref<FormInstance>();

// 自定义手机号验证
function validPhoneFn(rule: any, value: string, callback: any) {
  if (!value) {
    return callback(new Error("手机号不能为空"));
  } else if (!validPhone(value)) {
    return callback(new Error("请输入正确的手机号"));
  }
  callback();
}

const phoneRules = reactive<FormRules>({
  phone: [{ required: true, validator: validPhoneFn, trigger: "blur" }],
  code: [{ required: true, message: "请输入验证码", trigger: "blur" }],
});

const accountRules = reactive<FormRules>({
  login: [{ required: true, message: "请输入账号", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
});

const state = reactive({
  activeName: "login",
});

const { activeName } = toRefs(state);

const submit = async (formEl: FormInstance) => {
  switch (state.activeName) {
    case "login":
      accountLogin(formEl);
      break;
    case "phone":
      phoneLogin(formEl);
      break;
    default:
      break;
  }
};

function accountLogin(formEl: FormInstance) {
  formEl.validate(async (valid) => {
    if (valid) {
      if (newPrimaryWord.value.length != 0 && !protocol.value) {
        ElMessage.warning("请阅读并勾选同意协议");
        return;
      }
      let formData = new URLSearchParams(accountForm);
      emit("accountLoginHandle", formData);
    }
  });
}
let phoneState = ref("");
function phoneLogin(formEl: FormInstance) {
  formEl.validate(async (valid) => {
    if (valid) {
      if (newPrimaryWord.value.length != 0 && !protocol.value) {
        ElMessage.warning("请阅读并勾选同意协议");
        return;
      }
      const params = { ...phoneForm, phone: "+86" + phoneForm.phone };
      emit("phoneLoginHandle", params, phoneState.value);
    }
  });
}

const handleNavigate = (url: string) => {
  routerPath.push({
    path: "/protocol",
    state: { url },
  });
};

const thirdLogin = async (params: any) => {
  emit("thirdLoginHandle", params);
};
const isPhone = ref(false);

let phoneProvider = ref("");
const getLoginConfig = async () => {
  const option = ["wecom", "dingtalk", "wechat"];
  const data = (await getThirdLoginConfigs(currentTenant)) as ThirdLoginType[];
  const thirdLoginList = data ? data.filter((item) => option.includes(item.type)) : "";

  thirdLoginTypes.value = thirdLoginList;
  thirdLoginTypesLength.value = thirdLoginTypes.value.length;
  phoneProvider.value = data ? data.find((item) => item.type === "sms")!.name : "";
  if (data && data.find((item) => item.type === "sms")) {
    checkPhone();
  }
};

const countdownButtonRef = ref();
const sendValidCode = async (phone: string) => {
  phoneRuleFormRef.value?.validateField("phone", (valid: boolean) => {
    if (valid) {
      countdownButtonRef.value.startCountdown();
      phone = "%2B86" + phone;
      thirdLoginHandle(phoneProvider.value, phone, currentTenant).then((res) => {
        phoneState.value = res.state;
      });
    }
  });
};

getLoginConfig();
// 验证有手机号
const checkPhone = async () => {
  const res = await smsAvailable(currentTenant);
  if (res) {
    isPhone.value = true;
  } else {
    isPhone.value = false;
  }
};

const navigateToRegister = async () => {
  navigateTo({
    path: "/oauth2Register",
    query: route.query,
  });
};
const navigateRegister = async () => {
  navigateTo({
    path: "/register",
  });
};
const forgetPass = async () => {
  navigateTo({
    path: "/forgetpass",
    query: { currentTenant },
  });
};
const handleSubmit = (formEl: FormInstance) => {
  submit(formEl);
};
definePageMeta({
  layout: false,
});
</script>

<template>
  <div v-if="navigatePage"></div>
  <div
    v-else
    class="containerL"
    :style="{
      background: backgroundColor
        ? backgroundColor.substring(0, 1) != '#'
          ? `url(${backgroundColor}) center no-repeat`
          : backgroundColor
        : info && info.styleBgcolor && info.styleBgcolor.substring(0, 1) != '#'
        ? `url(${info && info.styleBgcolor && info.styleBgcolor}) center no-repeat`
        : info && info.styleBgcolor && info.styleBgcolor,
    }"
  >
    <div
      class="login-boxL"
      :style="{
        marginTop: `${numTop ? numTop : login_top}%`,
        marginLeft: `${numLeft ? numLeft : login_left}%`,
      }"
    >
      <div class="titleL">
        <span
          class="logoL"
          :style="{
            background: `url(${logoUpload ? logoUpload : info && info.styleLogo}) center`,
          }"
        ></span>
        {{ inputTitle ? inputTitle : info && info.styleName }}
      </div>
      <el-tabs v-model="activeName">
        <el-tab-pane label="密码登录" name="login">
          <el-form ref="accountRuleFormRef" :model="accountForm" :rules="accountRules">
            <el-form-item prop="login">
              <el-input
                v-model="accountForm.login"
                placeholder="请输入手机号/用户名"
                @keyup.enter="handleSubmit(accountRuleFormRef as FormInstance)"
              >
                <template #prefix>
                  <el-icon class="icon-userL"><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                v-model="accountForm.password"
                placeholder="请输入登录密码"
                type="password"
                show-password
                @keyup.enter="handleSubmit(accountRuleFormRef as FormInstance)"
              >
                <template #prefix>
                  <el-icon class="icon-passL"><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-form>

          <el-link
            @click="router.path.substring(0, 10) != '/dashboard' ? forgetPass() : ''"
            style="cursor: pointer; font-size: small; color: #409eff"
            v-if="
              router.path.substring(0, 10) == '/dashboard'
                ? isPhone && passSwitch
                : isPhone && info && info.stylePass
            "
            :underline="false"
            class="forgetL"
            >忘记密码</el-link
          >
          <el-button
            :loading="accountLoad"
            class="submit-btnL"
            type="primary"
            @click="router.path.substring(0, 10)!='/dashboard'?submit(accountRuleFormRef as FormInstance):''"
            >登 录</el-button
          >
        </el-tab-pane>

        <el-tab-pane
          label="验证码登录"
          name="phone"
          v-if="
            router.path.substring(0, 10) == '/dashboard'
              ? isPhone && codeSwitch
              : isPhone && info && info.styleCode
          "
        >
          <el-form ref="phoneRuleFormRef" :model="phoneForm" :rules="phoneRules">
            <el-form-item prop="phone">
              <el-input v-model="phoneForm.phone" placeholder="请输入手机号">
                <template #prefix>
                  <el-icon class="icon-phoneL"><Iphone /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item prop="code">
              <div class="verify-boxL">
                <el-input
                  v-model="phoneForm.code"
                  maxlength="6"
                  placeholder="请输入6位验证码"
                  :style="{ width: '280px', marginRight: '10px' }"
                  @keyup.enter="handleSubmit(phoneRuleFormRef as FormInstance)"
                >
                  <template #prefix>
                    <el-icon class="icon-codeL"><Lock /></el-icon>
                  </template>
                </el-input>
                <CountdownButton
                  ref="countdownButtonRef"
                  @click="
                    router.path.substring(0, 10) != '/dashboard'
                      ? sendValidCode(phoneForm.phone)
                      : ''
                  "
                ></CountdownButton>
              </div>
            </el-form-item>
          </el-form>
          <el-button
            :loading="accountLoad"
            class="submit-btnL"
            type="primary"
            @click="router.path.substring(0, 10) != '/dashboard' ? submit(phoneRuleFormRef as FormInstance):''"
            >登 录</el-button
          >
        </el-tab-pane>
      </el-tabs>
      <div class="protocol-boxL" v-if="newPrimaryWord.length != 0 || newTop.length != 0">
        <div>
          <el-checkbox v-model="protocol" size="large" style="margin-right: 5px" />
          我已同意
          <el-link
            v-for="item in newTop.length != 0 ? newTop : newPrimaryWord"
            @click="
              router.path.substring(0, 10) != '/dashboard'
                ? handleNavigate(item.privacyWrite)
                : ''
            "
            type="primary"
            >《{{ item.privacyWrite }}》</el-link
          >
        </div>
      </div>
      <div class="optionL">
        <div
          class="other-loginL"
          v-if="
            router.path.substring(0, 10) == '/dashboard'
              ? thirdLoginTypesLength > 0 && loginSwitch
              : thirdLoginTypesLength > 0 && info && info.styleLogin
          "
        >
          其它方式登录：
          <svg-icon
            v-for="item in thirdLoginTypes"
            :name="item.type"
            @click="router.path.substring(0, 10) != '/dashboard' ? thirdLogin(item) : ''"
            size="1.5em"
            style="margin-left: 4px"
          ></svg-icon>
        </div>
        <div v-else></div>
        <nuxt-link
          v-if="
            router.path.substring(0, 10) == '/dashboard'
              ? currentTenant == 'default' && regionSwitch
              : currentTenant == 'default' && info && info.styleRegion
          "
          @click="router.path.substring(0, 10) != '/dashboard' ? navigateRegister() : ''"
        >
          <span style="cursor: pointer" class="regionL">注册账户</span>
        </nuxt-link>
        <nuxt-link
          v-if="
            router.path.substring(0, 10) == '/dashboard'
              ? regionSwitch && hasRegister && currentTenant !== 'default'
              : info && info.styleRegion && hasRegister && currentTenant !== 'default'
          "
          @click="
            router.path.substring(0, 10) != '/dashboard' ? navigateToRegister() : ''
          "
        >
          <span style="cursor: pointer" class="regionL">注册账户</span>
        </nuxt-link>
      </div>
    </div>
    <div class="bottomL">
      <a
        :href="
          router.path.substring(0, 10) != '/dashboard'
            ? item.wordlink == ''
              ? 'javascript:;'
              : item.wordlink
            : 'javascript:;'
        "
        v-for="item in bottom.length != 0 ? bottom : bottomTitle"
        :style="{ cursor: item.wordlink == '' ? 'not-allowed' : 'pointer' }"
        >{{ item.wordCen }}</a
      >
    </div>
  </div>
</template>

<style scoped lang="scss">
.containerL {
  overflow-y: auto;
  height: 100%;
  width: 100%;
  background-size: cover !important;
  background-color: #f7f8fa;
  position: relative;
  .login-boxL {
    position: absolute;
    flex: 1;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 400px;
    min-height: 250px;
    background-color: #fff;
    padding: 20px 20px;
    border-radius: 8px;
    box-shadow: 0 10px 15px -3px rgb(0 0 0/0.1), 0 4px 6px -4px rgb(0 0 0/0.1);
    .titleL {
      margin-bottom: 10px;
      font-size: 20px;
      display: flex;
      .logoL {
        margin-right: 10px;
        width: 30px;
        height: 30px;
        background-size: cover !important;
      }
    }

    .send-code-btnL {
      width: 150px;
      margin-left: 10px;
    }

    .submit-btnL {
      width: 100%;
      margin: 10px 0;
    }

    .verify-boxL {
      width: 100%;
      display: flex;
      justify-content: space-between;
    }

    .optionL {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .other-loginL {
        font-size: 14px;
        line-height: 40px;

        .svg-iconL {
          margin-right: 10px;
          cursor: pointer;
        }
      }
    }
    .protocol-boxL {
      display: flex;
      align-items: center;
      font-size: 14px;
    }
  }
  .bottomL {
    width: 100%;
    height: 20px;
    position: absolute;
    bottom: 5%;
    text-align: center;
    a {
      text-decoration: none;
      font-size: 16px;
      color: #000;
      display: inline;
    }
  }
}
</style>
