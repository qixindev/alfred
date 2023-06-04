<script setup lang="ts">
import { User, useTenant, Path } from '~~/composables/useUser'
import { getUser } from '~/api/common'
const loginVisible: Ref<boolean> = useState('loginVisible')
//获取
const user = useState<User>('user')
const path = useState<Path>('path')
interface SelectOption {
  name: string,
  id: number,
}
const state = reactive({
  dataList: <SelectOption[]>[],
})
const showLogin = () => {
  navigateTo('/login')
}

const logout = () => {
  useLogout()
}
const tenant = useTenant();
/** 用户列表 */
function getList() {
  getUser().then((res: any) => {
    state.dataList = res
    tenant.value = state.dataList[0].name,
      console.log(res, "res");
  }).finally(() => {
  })
}
// 切换列表
function clickUser(row: any) {

  tenant.value = row.name;
  if (path.value.path == '/') {
    navigateTo(`${path.value.path}`)
  } else {
    navigateTo(`/${row.name}${path.value.path}`)
  }
}

onMounted(() => {
  getList()


})
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


        <el-dropdown class="avatar-container right-menu-item hover-effect" trigger="click" v-if="user">
          <div class="avatar-wrapper">
            <!-- <img v-if="avatar" :src="avatar" class="user-avatar" /> -->
            {{ user.username }}
          </div>

          <template #dropdown>
            <el-menu default-active="2" class="el-menu-vertical-demo">
              <el-menu-item v-for="(item, index) in state.dataList" :index="index" @click="clickUser(item)">
                <el-icon><icon-menu /></el-icon>
                <span>{{ item.name }}</span>
              </el-menu-item>
            </el-menu>
            <el-dropdown-menu>
              <!-- <nuxt-link to="/profile">
                        <el-dropdown-item>个人中心</el-dropdown-item>
                      </nuxt-link> -->

              <el-dropdown-item>
                用户
              </el-dropdown-item>

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
    z-index: 99;

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

  .el-scrollbar {
    overflow: auto;
    :deep(.el-menu-vertical-demo) {
      position: absolute;
      top: 61px;
      right: 112px;
    }

  }
}
</style>
