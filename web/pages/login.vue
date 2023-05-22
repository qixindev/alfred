<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { login, getThirdLoginConfigs, getThirdLoginConfigByName } from '~/api/user';

interface ThirdLoginType {
  id: number,
  name: string,
  type: string
}

const phoneForm = reactive({
  phone: '',
  code: ''
})

const accountForm = reactive({
  login: '',
  password: ''
})

let thirdLoginTypes= ref<ThirdLoginType[]>([])

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
    { required: true,  validator: validPhoneFn, trigger: 'blur'}
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur'}
  ]
})

const accountRules = reactive<FormRules>({
  login: [
    { required: true, message: '请输入账号', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur'}
  ]
})

const state = reactive({
  activeName: 'login'
})

const {
  activeName,
} = toRefs(state)

const submit = async (formEl: FormInstance) => {
  await formEl.validate(async (valid) => {
    if (valid) {
      let formData = new URLSearchParams(accountForm)
      login(formData).then(res => {
        if (res == 10000) {
          ElMessage({
            message: '账号或密码错误',
            type: 'error'
          })
        } else {
          const route = useRoute()
          navigateTo(route.query.from as string || '/', { replace: true })
        }
      })
    }
  })
}

const handleClick = () => {
}

const thirdLogin = async (params: any) => {
  const config = await getThirdLoginConfigByName(params.name)
  const redirect_uri  = location.origin + '/redirect'
  switch (params.type) {
    case 'dingtalk':
      navigateTo(`https://login.dingtalk.com/oauth2/auth?redirect_uri=${redirect_uri}&response_type=code&client_id=${config.appKey}&scope=openid&prompt=consent&state=${params.name}`, { external: true})
      break;
      case 'wecom':
      navigateTo(`https://login.work.weixin.qq.com/wwlogin/sso/login?appid=${config.corpId}&redirect_uri=${redirect_uri}&state=${encodeURI(JSON.stringify(params))}&agentid=${config.agentId}`, { external: true})
      break;
  
    default:
      break;
  }
}

const getLoginConfig  = async () => {
  thirdLoginTypes.value = await getThirdLoginConfigs() as ThirdLoginType[]
}

getLoginConfig()

definePageMeta({
  layout: false
})

</script>

<template>
  <div class="container">
    <div class="login-box">
      <div class="title">登录</div>
      <el-tabs v-model="activeName" @tab-click="handleClick">
        <!-- <el-tab-pane label="手机号登录" name="phone">
          <el-form ref="phoneRuleFormRef" :model="phoneForm" :rules="phoneRules">
            <el-form-item prop="phone">
              <el-input v-model="phoneForm.phone" placeholder="手机号">
                <template #prefix>
                  <svg-icon name="user"></svg-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="code">
              <div class="verify-box">
                <el-input v-model="phoneForm.code" placeholder="验证码" :style="{width: '280px'}">
                  <template #prefix>
                    <svg-icon name="password"></svg-icon>
                  </template>
                </el-input>
                <el-button>获取验证码</el-button>
              </div>
            </el-form-item>
          </el-form>

          <el-button class="submit-btn" type="primary" @click="submit(phoneRuleFormRef as FormInstance)">登 录/注 册</el-button>

        </el-tab-pane> -->
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
          <el-button class="submit-btn" type="primary" @click="submit(accountRuleFormRef as FormInstance)">登 录</el-button>

        </el-tab-pane>
      </el-tabs>
      
      <div class="option">
        <div class="other-login">其它方式登录： 
          <svg-icon v-for="item in thirdLoginTypes" :name="item.type" @click="thirdLogin(item)" size="1.5em"></svg-icon>
        </div>
        <nuxt-link to="/register" >
          <span>注册账户</span>
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
    height: 300px;
    background-color: #FFF;
    padding: 30px 20px;
    border-radius: 8px;
    box-shadow: 0 10px 15px -3px rgb(0 0 0/0.1),0 4px 6px -4px rgb(0 0 0/0.1);
    .title {
      margin-bottom: 10px;
      font-size: 20px;
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