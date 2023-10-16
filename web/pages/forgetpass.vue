<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import CountdownButton from '@/components/CountdownButton/index.vue'
import { thirdLoginHandle, verifyResetPasswordRequest, getResetPasswordToken, resetPassword } from '~/api/user';
import { ElMessage } from 'element-plus'


const route = useRoute()
let currentTenant = route.query.currentTenant ?? 'default'
const accountForm = reactive({
    phone: '',
    code: '',
    password: '',
    againpassword: ''
})
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
// 自定义密码验证
function validPasswordFn(rule: any, value: string, callback: any) {
    if (!value) {
        return callback(new Error('请输入密码'))
    } else if (!validPassword(value)) {
        return callback(new Error('请输入6 - 16位密码，区分大小写'))
    }
    callback()
}
// 自定义再次密码验证
function validagainPasswordFn(rule: any, value: string, callback: any) {
    if (!value) {
        return callback(new Error('请再次输入您的密码'))
    } else if (value !== accountForm.password) {
        return callback(new Error('两次输入密码不一致'))
    }
    callback()
}
const accountRules = reactive<FormRules>({
    phone: [
        { required: true, validator: validPhoneFn, trigger: 'blur' }
    ],
    code: [
        { required: true, message: '请输入验证码', trigger: 'blur' }
    ],
    password: [
        { required: true, validator: validPasswordFn, trigger: 'blur' }
    ],
    againpassword: [
        { required: true, validator: validagainPasswordFn, trigger: 'blur' }
    ]
})
// 验证码

const countdownButtonRef = ref()
const sendValidCode = async (phone: string) => {
    accountRuleFormRef.value?.validateField('phone', (valid: boolean) => {
        if (valid) {
            phone = '%2B86' + phone
            let formData = new URLSearchParams({ verifyMethod: "phonePassCode", passCodePayload: accountForm.phone, areaCode: "+86" })
            verifyResetPasswordRequest(formData, currentTenant).then((res: any) => {
                if (res) {
                    ElMessage({
                        message: res.message == 'no such user' ? "该用户不存在，去注册" : "请等待，已发送验证码",
                        type: 'warning',
                    })
                } else {
                    countdownButtonRef.value.startCountdown()
                    ElMessage({
                        message: '已发送验证码',
                        type: 'success',
                    })
                }

            }).catch(() => {
                ElMessage({
                    message: "请等待，已发送验证码",
                    type: 'warning',
                })
            })
        }
    })
}

const submit = async (formEl: FormInstance) => {
    await formEl.validate(async (valid) => {
        if (valid) {
            const { phone, code, password, againpassword } = accountForm
            let formData = new URLSearchParams({ code: code })
            getResetPasswordToken(formData, currentTenant).then((res: any) => {
                if (res) {
                    resetPassword(new URLSearchParams({ newPassword: againpassword }), currentTenant, res).then((ram: any) => {
                        if (ram.code == 200) {
                            ElMessage({
                                message: '重置密码成功',
                                type: 'success',
                            })
                            navigateTo('/login')
                        } else {
                            ElMessage({
                                message: '重置密码失败',
                                type: 'warning',
                            })

                        }

                    })
                }
            })

        }
    })
}
// 页面
definePageMeta({
    layout: false
})
</script>

<template>
    <div class="container">
        <div class="login-box">
            <div class="title">重置密码</div>
            <el-form ref="accountRuleFormRef" :model="accountForm" :rules="accountRules">
                <el-form-item prop="phone">
                    <el-input v-model="accountForm.phone" placeholder="请输入手机号">
                        <template #prefix>
                            <svg-icon name="user"></svg-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item prop="code">
                    <div class="verify-box">
                        <el-input v-model="accountForm.code" maxlength="6" placeholder="验证码"
                            :style="{ width: '280px', marginRight: '10px' }">
                            <template #prefix>
                                <svg-icon name="password"></svg-icon>
                            </template>
                        </el-input>
                        <CountdownButton ref="countdownButtonRef" @click="sendValidCode(accountForm.phone)">
                        </CountdownButton>
                    </div>
                </el-form-item>
                <el-form-item prop="password">
                    <el-input v-model="accountForm.password" placeholder="请输入新密码,6 - 16位密码，区分大小写" type="password"
                        show-password>
                        <template #prefix>
                            <svg-icon name="password"></svg-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item prop="againpassword">
                    <el-input v-model="accountForm.againpassword" placeholder="确认新密码" type="password" show-password>
                        <template #prefix>
                            <svg-icon name="password"></svg-icon>
                        </template>
                    </el-input>
                </el-form-item>
            </el-form>
            <el-button class="submit-btn" type="primary" @click="submit(accountRuleFormRef as FormInstance)">重 置
            </el-button>
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

        .protocol-box {
            display: flex;
            align-items: center;
            font-size: 14px;
        }
    }
}
</style>