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
const router = useRoute();
const routerPath = useRouter();
const info = ref({});
const bottomTitle = ref([]);
const newPrimaryWord = ref([]);
const newTop = ref([]);
const getInfo = () => {
  getEnergy().then((res: any) => {
    info.value = { ...res };
    bottomTitle.value = [...res.bottom];
  });
  getProto().then((res: any) => {
    newPrimaryWord.value = res.filter((item: any) => {
      return item.loginSwitch;
    });
  });
};
getInfo();
const protocol = ref(false);
const props = defineProps({
  loginSwitch: {
    type: Boolean,
    default: true,
  },
  regionSwitch: {
    type: Boolean,
    default: true,
  },
  passSwitch: {
    type: Boolean,
    default: true,
  },
  codeSwitch: {
    type: Boolean,
    default: true,
  },
  inputTitle: {
    type: String,
    default: "登录",
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
  top: {
    type: Array,
    default: [],
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
    document.head.appendChild(style);
  },
  { immediate: true, deep: true }
);
watch(
  () => props.top,
  () => {
    newTop.value = props.top.filter((item: any) => {
      return item.loginSwitch;
    });
  },
  { immediate: true, deep: true }
);
const emit = defineEmits(["accountLoginHandle", "phoneLoginHandle", "thirdLoginHandle"]);
const route = useRoute();
const { state: tanent } = route.query as any;
let currentTenant = tanent ?? "default";

const phoneForm = reactive({
  phone: "",
  code: "",
});

const accountForm = reactive({
  login: "",
  password: "",
});

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

function phoneLogin(formEl: FormInstance) {
  formEl.validate(async (valid) => {
    if (valid) {
      if (newPrimaryWord.value.length != 0 && !protocol.value) {
        ElMessage.warning("请阅读并勾选同意协议");
        return;
      }
      const params = { ...phoneForm, phone: "+86" + phoneForm.phone };
      emit("phoneLoginHandle", phoneProvider.value, params);
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

let phoneProvider = ref("");
const getLoginConfig = async () => {
  const option = ["wecom", "dingtalk"];
  const data = (await getThirdLoginConfigs(currentTenant)) as ThirdLoginType[];
  const thirdLoginList = data.filter((item) => option.includes(item.type));
  thirdLoginTypes.value = thirdLoginList;
  thirdLoginTypesLength.value = thirdLoginTypes.value.length;
  phoneProvider.value = data.find((item) => item.type === "sms")!.name;
};

const countdownButtonRef = ref();
const sendValidCode = async (phone: string) => {
  phoneRuleFormRef.value?.validateField("phone", (valid: boolean) => {
    if (valid) {
      countdownButtonRef.value.startCountdown();
      phone = "%2B86" + phone;
      thirdLoginHandle(phoneProvider.value, phone, currentTenant);
    }
  });
};

getLoginConfig();
// 验证有手机号
const isPhone = ref(true);
const checkPhone = async () => {
  const res = await smsAvailable(currentTenant);
  if (res) {
    isPhone.value = true;
  } else {
    isPhone.value = false;
  }
};
checkPhone();
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
definePageMeta({
  layout: false,
});
</script>

<template>
  <div
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
    <div class="login-boxL">
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
        <el-tab-pane label="账户密码登录" name="login">
          <el-form ref="accountRuleFormRef" :model="accountForm" :rules="accountRules">
            <el-form-item prop="login">
              <el-input v-model="accountForm.login" placeholder="账号">
                <template #prefix>
                  <el-icon class="icon-userL"><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                v-model="accountForm.password"
                placeholder="密码"
                type="password"
                show-password
              >
                <template #prefix>
                  <el-icon class="icon-passL"><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-form>

          <nuxt-link
            @click="router.path.substring(0, 10) != '/dashboard' ? forgetPass : ''"
            style="cursor: pointer; font-size: small; color: #409eff; width: 60px"
            v-if="isPhone && passSwitch && info && info.stylePass"
            class="forgetL"
            >忘记密码</nuxt-link
          >
          <el-button
            class="submit-btnL"
            type="primary"
            @click="router.path.substring(0, 10)!='/dashboard'?submit(accountRuleFormRef as FormInstance):''"
            >登 录</el-button
          >
        </el-tab-pane>

        <el-tab-pane
          label="手机号登录"
          name="phone"
          v-if="codeSwitch && info && info && info.styleCode && isPhone"
        >
          <el-form ref="phoneRuleFormRef" :model="phoneForm" :rules="phoneRules">
            <el-form-item prop="phone">
              <el-input v-model="phoneForm.phone" placeholder="手机号">
                <template #prefix>
                  <el-icon class="icon-phoneL"><Iphone /></el-icon>
                  <span>+86</span>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item prop="code">
              <div class="verify-boxL">
                <el-input
                  v-model="phoneForm.code"
                  maxlength="6"
                  placeholder="验证码"
                  :style="{ width: '280px', marginRight: '10px' }"
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
      <div
        class="optionL"
        v-if="
          (loginSwitch && info && info.styleLogin) ||
          (regionSwitch && info && info.styleRegion)
        "
      >
        <div
          class="other-loginL"
          v-if="loginSwitch && info && info.styleLogin && thirdLoginTypesLength > 0"
        >
          其它方式登录：
          <svg-icon
            v-for="item in thirdLoginTypes"
            :name="item.type"
            @click="router.path.substring(0, 10) != '/dashboard' ? thirdLogin(item) : ''"
            size="1.5em"
          ></svg-icon>
        </div>
        <div v-else></div>
        <nuxt-link
          v-if="regionSwitch && info && info.styleRegion && currentTenant == 'default'"
          @click="router.path.substring(0, 10) != '/dashboard' ? navigateRegister() : ''"
        >
          <span style="cursor: pointer" class="regionL">注册账户</span>
        </nuxt-link>
        <nuxt-link
          v-if="
            regionSwitch &&
            info &&
            info.styleRegion &&
            currentTenant !== 'default' &&
            hasRegister
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
        class="linkL"
        >{{ item.wordCen }}</a
      >
    </div>
  </div>
</template>

<style scoped lang="scss">
.containerL {
  overflow: hidden;
  height: 100%;
  width: 100%;
  .login-boxL {
    flex: 1;
    margin: 10%;
    width: 400px;
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
      height: 40px;
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
    width: 80%;
    height: 20px;
    // margin-top: 10%;
    position: absolute;
    bottom: 5%;
    display: flex;
    justify-content: center;
    .linkL {
      text-decoration: none;
      margin-right: 10px;
      font-size: 16px;
      color: #000;
    }
  }
}
</style>
