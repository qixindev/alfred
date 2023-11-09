<script setup lang="ts">
import { register } from "~/api/user";
import type { FormInstance, FormRules } from "element-plus";
import { ElMessage } from "element-plus";
import { getProto } from "~/api/energy";
const form = reactive({
  login: "",
  password: "",
  copyPassword: "",
});

const ruleFormRef = ref<FormInstance>();
const newPrimaryWord = ref([]);
const getInfo = () => {
  getProto().then((res: any) => {
    newPrimaryWord.value = res.filter((item: any) => {
      return item.resignSwitch;
    });
  });
};
getInfo();
// 自定义密码确认
function validCopyPasswordFn(rule: any, value: string, callback: any) {
  if (!value) {
    return callback(new Error("请再次输入您的密码"));
  } else if (value !== form.password) {
    return callback(new Error("两次输入密码不一致"));
  }
  callback();
}

// 自定义手机号验证
function validPhoneFn(rule: any, value: string, callback: any) {
  if (!value) {
    return callback(new Error("手机号不能为空"));
  } else if (!validPhone(value)) {
    return callback(new Error("请输入正确的手机号"));
  }
  callback();
}

const rules = reactive<FormRules>({
  login: [{ required: true, message: "请输入账号", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
  copyPassword: [{ required: true, validator: validCopyPasswordFn, trigger: "blur" }],
  phone: [{ required: true, validator: validPhoneFn, trigger: "blur" }],
});

const handleNavigate = (url: string) => {
  window.open(url, "_blank");
};
const submit = async (formEl: FormInstance) => {
  await formEl.validate(async (valid) => {
    if (valid) {
      if (newPrimaryWord.value.length != 0) {
        ElMessage.warning("请阅读并勾选同意协议");
        return;
      }
      const { login, password } = form;
      const param = { login, password };
      let formData = new URLSearchParams(param);
      register(formData).then((res) => {
        navigateTo("/");
      });
    }
  });
};

definePageMeta({
  layout: false,
});
</script>

<template>
  <div class="container">
    <div class="login-box">
      <div class="title">注册</div>
      <el-form ref="ruleFormRef" :model="form" :rules="rules">
        <el-form-item prop="login">
          <el-input v-model="form.login" placeholder="账号">
            <template #prefix>
              <svg-icon name="user"></svg-icon>
            </template>
          </el-input>
        </el-form-item>
        <!-- <el-form-item prop="phone">
          <el-input v-model="form.phone" placeholder="手机号">
            <template #prefix>
                <svg-icon name="phone"></svg-icon>
              </template>
          </el-input>
        </el-form-item> -->
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="6 - 16位密码，区分大小写"
          >
            <template #prefix>
              <svg-icon name="password"></svg-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="copyPassword">
          <el-input v-model="form.copyPassword" type="password" placeholder="确认密码">
            <template #prefix>
              <svg-icon name="password"></svg-icon>
            </template>
          </el-input>
        </el-form-item>
      </el-form>

      <el-button
        class="login-btn"
        type="primary"
        @click="submit(ruleFormRef as FormInstance)"
        >注册</el-button
      >
      <div class="protocol-boxL" v-if="newPrimaryWord.length != 0">
        <div>
          <el-checkbox v-model="protocol" size="large" style="margin-right: 5px" />
          我已同意
          <el-link
            v-for="item in newPrimaryWord"
            @click="handleNavigate('/protocol/userServiceAgreement')"
            type="primary"
            >《{{ item.privacyWrite }}》</el-link
          >
        </div>
      </div>
      <div class="tip">
        <nuxt-link to="/login">
          <span>使用已有账号登录</span>
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
    background-color: #fff;
    padding: 30px 20px;
    border-radius: 8px;
    box-shadow: 0 10px 15px -3px rgb(0 0 0/0.1), 0 4px 6px -4px rgb(0 0 0/0.1);
    .title {
      margin-bottom: 10px;
    }
    .login-btn {
      width: 100%;
    }
    .tip {
      text-align: end;
      font-size: 14px;
      line-height: 40px;
      span {
        cursor: pointer;
      }
    }
  }
}
</style>
