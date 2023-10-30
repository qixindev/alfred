<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import CountdownButton from '@/components/CountdownButton/index.vue'

import { getThirdLoginConfigs, getThirdLoginConfigByName, thirdLoginHandle, smsAvailable } from '~/api/user';

interface ThirdLoginType {
  id: number,
  name: string,
  type: string
}

const emit = defineEmits(['accountLoginHandle', 'phoneLoginHandle', 'thirdLoginHandle'])
const route = useRoute()
const { state: tanent } = route.query as any
let currentTenant = tanent ?? 'default'

const phoneForm = reactive({
  phone: '',
  code: ''
})

const accountForm = reactive({
  login: '',
  password: ''
})

const hasRegister = computed(() => {
  return route.query?.platform === 'tenant'
})

let thirdLoginTypes = ref<ThirdLoginType[]>([])
let thirdLoginTypesLength = ref()
const phoneRuleFormRef = ref<FormInstance>()
const accountRuleFormRef = ref<FormInstance>()

// 自定义手机号验证
function validPhoneFn(rule: any, value: string, callback: any) {
  if (!value) {
    return callback(new Error('手机号不能为空'))
  } else if (!validPhone(value)) {
    return callback(new Error('请输入正确的手机号'))
  }
  callback()
}

const phoneRules = reactive<FormRules>({
  phone: [
    { required: true, validator: validPhoneFn, trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' }
  ]
})

const accountRules = reactive<FormRules>({
  login: [
    { required: true, message: '请输入账号', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
})

const state = reactive({
  activeName: 'login'
})

const {
  activeName,
} = toRefs(state)

const submit = async (formEl: FormInstance) => {
  switch (state.activeName) {
    case 'login':
      accountLogin(formEl)
      break;
    case 'phone':
      phoneLogin(formEl)
      break;
    default:
      break;
  }
}

function accountLogin(formEl: FormInstance) {
  formEl.validate(async (valid) => {
    if (valid) {
      let formData = new URLSearchParams(accountForm)
      emit('accountLoginHandle', formData)
    }
  })
}

function phoneLogin(formEl: FormInstance) {
  formEl.validate(async (valid) => {
    if (valid) {
      const params = { ...phoneForm, phone: '+86' + phoneForm.phone }
      emit('phoneLoginHandle', phoneProvider.value, params)
    }
  })
}

const handleNavigate = (url: string) => {
  window.open(url, '_blank')
}

const thirdLogin = async (params: any) => {
  emit('thirdLoginHandle', params)
}

let phoneProvider = ref('');
const getLoginConfig = async () => {
  const option = ['wecom', 'dingtalk']
  const data = await getThirdLoginConfigs(currentTenant) as ThirdLoginType[]
  const thirdLoginList = data.filter(item => option.includes(item.type))
  thirdLoginTypes.value = thirdLoginList
  thirdLoginTypesLength.value = thirdLoginTypes.value.length
  phoneProvider.value = data.find(item => item.type === 'sms')!.name
}

const countdownButtonRef = ref()
const sendValidCode = async (phone: string) => {
  phoneRuleFormRef.value?.validateField('phone', (valid: boolean) => {
    if (valid) {
      countdownButtonRef.value.startCountdown()
      phone = '%2B86' + phone
      thirdLoginHandle(phoneProvider.value, phone, currentTenant)
    }
  })
}

getLoginConfig()
// 验证有手机号
const isPhone = ref(true)
const checkPhone = async () => {
  const res = await smsAvailable(currentTenant)
  if (res) {
    isPhone.value = true
  } else {
    isPhone.value = false
  }
}
checkPhone()
const navigateToRegister = async () => {
  navigateTo({
    path: '/oauth2Register',
    query: route.query
  })
}
const forgetPass = async () => {
  navigateTo({
    path: '/forgetpass',
    query: { currentTenant }
  })
}
definePageMeta({
  layout: false
})

</script>

<template>
  <div class="container">
    <div class="login-box">
      <div class="title">登录</div>
      <el-tabs v-model="activeName">
        <el-tab-pane label="账户密码登录" name="login">
          <el-form ref="accountRuleFormRef" :model="accountForm" :rules="accountRules">
            <el-form-item prop="login">
              <el-input v-model="accountForm.login" placeholder="账号">
                <template #prefix>
                  <svg-icon name="user"></svg-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item prop="password">
              <el-input v-model="accountForm.password" placeholder="密码" type="password" show-password>
                <template #prefix>
                  <svg-icon name="password"></svg-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-form>

          <nuxt-link @click="forgetPass" style="cursor: pointer;font-size: small;color:#409eff ;width:60px;" v-if="isPhone">忘记密码</nuxt-link>
          <el-button class="submit-btn" type="primary" @click="submit(accountRuleFormRef as FormInstance)">登 录</el-button>
        </el-tab-pane>

        <el-tab-pane label="手机号登录" name="phone" v-if="phoneProvider && isPhone">
          <el-form ref="phoneRuleFormRef" :model="phoneForm" :rules="phoneRules">
            <el-form-item prop="phone">
              <el-input v-model="phoneForm.phone" placeholder="手机号">
                <template #prefix>
                  <svg-icon name="phone"></svg-icon>
                  <span>+86</span>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item prop="code">
              <div class="verify-box">
                <el-input v-model="phoneForm.code" maxlength="6" placeholder="验证码"
                  :style="{ width: '280px', marginRight: '10px' }">
                  <template #prefix>
                    <svg-icon name="password"></svg-icon>
                  </template>
                </el-input>
                <CountdownButton ref="countdownButtonRef" @click="sendValidCode(phoneForm.phone)"></CountdownButton>
              </div>
            </el-form-item>
          </el-form>
          <el-button class="submit-btn" type="primary" @click="submit(phoneRuleFormRef as FormInstance)">登 录</el-button>
        </el-tab-pane>
      </el-tabs>

      <div class="option">
        <div class="other-login" v-if="thirdLoginTypesLength > 0">其它方式登录：
          <svg-icon v-for="item in thirdLoginTypes" :name="item.type" @click="thirdLogin(item)" size="1.5em"></svg-icon>
        </div>
        <div v-else></div>
        <nuxt-link v-if="currentTenant == 'default'" to="/register">
          <span>注册账户</span>
        </nuxt-link>
        <nuxt-link v-if="currentTenant !== 'default' && hasRegister" @click="navigateToRegister">
          <span style="cursor: pointer;">注册账户</span>
        </nuxt-link>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.container {
  display: flex;
  justify-content: center;
  height: 100vh;
  background-color: #eee;

  .login-box {
    position: absolute;
    top: 20%;
    width: 400px;
    background-color: #FFF;
    padding: 20px 20px;
    border-radius: 8px;
    box-shadow: 0 10px 15px -3px rgb(0 0 0/0.1), 0 4px 6px -4px rgb(0 0 0/0.1);

    .title {
      margin-bottom: 10px;
      font-size: 20px;
    }

    .send-code-btn {
      width: 150px;
      margin-left: 10px;
    }

    .submit-btn {
      width: 100%;
      margin: 10px 0;
    }

    .verify-box {
      width: 100%;
      display: flex;
      justify-content: space-between;
    }

    .option {
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: space-between;

      .other-login {
        font-size: 14px;
        line-height: 40px;

        .svg-icon {
          margin-right: 10px;
          cursor: pointer;
        }
      }
    }

  }
}
</style>