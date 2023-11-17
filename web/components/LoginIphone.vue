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
const getInfo = () => {
  getEnergy(currentTenant).then((res: any) => {
    info.value = { ...res };
    bottomTitle.value = [...res.bottom];
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
// watch(
//   () => tenant.value,
//   () => {
//     if (route.path.substring(0, 10) == "/dashboard") {
//       style.textContent = "";
//     }
//   },
//   { immediate: true, deep: true }
// );
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
const emit = defineEmits(["accountLoginHandle", "phoneLoginHandle", "thirdLoginHandle"]);

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
const isPhone = ref(false);

let phoneProvider = ref("");
const getLoginConfig = async () => {
  const option = ["wecom", "dingtalk"];
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
      thirdLoginHandle(phoneProvider.value, phone, currentTenant);
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

definePageMeta({
  layout: false,
});
</script>

<template>
  <div class="login-boxL" style="padding-top: 0 !important; padding-bottom: 0 !important">
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
            <el-input v-model="accountForm.login" placeholder="请输入手机号/用户名">
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
          v-if="
            router.path.substring(0, 10) == '/dashboard'
              ? isPhone && passSwitch
              : isPhone && info && info.stylePass
          "
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
    <div class="optionL">
      <div></div>
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
        @click="router.path.substring(0, 10) != '/dashboard' ? navigateToRegister() : ''"
      >
        <span style="cursor: pointer" class="regionL">注册账户</span>
      </nuxt-link>
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
        :style="{ cursor: item.wordlink == '' ? 'not-allowed' : 'pointer' }"
        >{{ item.wordCen }}</a
      >
    </div>
  </div>
</template>

<style scoped lang="scss">
.login-boxL {
  position: relative;
  overflow: hidden;
  height: 100%;
  flex: 1;
  margin: auto;
  background-color: #fff;
  padding: 0 20px 0 20px;
  .titleL {
    margin-top: 120px;
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
  height: 30px;
  // margin-top: 10%;
  position: absolute;
  left: 0;
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
</style>
