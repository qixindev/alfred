<script setup lang="ts">
import type { User } from '~~/composables/useUser'
const loginVisible: Ref<boolean> = useState('loginVisible')

const user = useState<User>('user')

const showLogin = () => {
  navigateTo('/login')
}

const logout = () => {
  useLogout()
}

// const avatar = ref('https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg')
</script>

<template>
  <div class="navbar">
    <div class="navbar-box">
      <div class="left">
        <nuxt-link to="/">
          Logo
        </nuxt-link>
      </div>
      <div class="center"></div>
      <div class="right">
        <el-dropdown 
          class="avatar-container right-menu-item hover-effect" 
          trigger="click"
          v-if="user"
          >
          <div class="avatar-wrapper">
            <!-- <img v-if="avatar" :src="avatar" class="user-avatar" /> -->
            {{ user.username }}
          </div>

          <template #dropdown>
            <el-dropdown-menu>
              <!-- <nuxt-link to="/profile">
                <el-dropdown-item>个人中心</el-dropdown-item>
              </nuxt-link> -->
              <el-dropdown-item @click="logout">
                退出
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button v-else @click="showLogin">登录</el-button>
      </div>
    </div>
  </div>
  <Login :visible="loginVisible"></Login>
</template>

<style scoped lang="scss">
.navbar {
  display: flex;
  width: 100%;
  justify-content: center;
  background-color: #FFF;
  height: 48px;
  line-height: 48px;
  box-shadow: 0 1px 2px 0 rgb(0 0 0/0.05);

  .navbar-box {
    display: flex;
    width: 100%;
    padding: 0 20px;
    justify-content: space-between;
  }

  .right {
    float: right;
    height: 100%;
    line-height: 64px;
    display: flex;
    justify-content: space-between;
    align-items: center;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background 0.3s;
      }
    }

    .avatar-container {
      margin-right: 30px;

      .avatar-wrapper {
        display: flex;
        align-items: center;
        height: 100%;

        .user-avatar {
          cursor: pointer;
          width: 32px;
          height: 32px;
          border-radius: 50%;
        }

        .icon-msg-expand {
          height: 10px;
          width: 2px;
          margin-left: 8px;
        }
      }
    }
  }
}
</style>
