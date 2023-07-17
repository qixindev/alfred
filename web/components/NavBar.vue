<script setup lang="ts">
import { User, useTenant } from '~~/composables/useUser'
import { getUser } from '~/api/common'
import { ref} from 'vue'
import { ElMessage } from 'element-plus';
const popoverRef = ref(null);
const activeUser = ref(null);
const loginVisible: Ref<boolean> = useState('loginVisible')
const route = useRoute();
const routerTenant = useRouter()
//获取
const user = useState<User>('user')
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
    if(!res){
      ElMessage({
        message:'当前没有租户，请创建租户',
        type:'error'
      })
    }
    state.dataList = res
    //默认第一个
    tenant.value = localStorage.getItem('tenantValue') ? localStorage.getItem('tenantValue') : state.dataList?.[0].name;
    localStorage.setItem("tenantValue", tenant.value)
    // 高亮
    activeUser.value = localStorage.getItem('tenantValue');
  }).finally(() => {
  })
}
// 切换列表
function clickUser(row: any) {
  //点击用户关闭
  popoverRef.value.hide()
  tenant.value = row.name;
  localStorage.setItem("tenantValue", tenant.value as string)
  let arr = route.path.split('/')
  arr.splice(2, 1, tenant.value);
  arr.join("/")
  if (route.path == '/') {
    navigateTo('/')
  } else { //非主页
    navigateTo(arr.join("/"))
  }
};
onMounted(() => {
  getList()
})
// 监听当前路由
watch(
  () => routerTenant.currentRoute.value,
  (newValue: any) => {
    // 高亮
    activeUser.value = localStorage.getItem('tenantValue');
  },
  { immediate: true }
)
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
        <el-dropdown trigger="click" v-if="user">
          <div>
            <span>{{ user.username }}</span>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>
                <el-popover ref="popoverRef" :hide-after="0" placement="left-start" trigger="hover" :offset="15">
                  <template #reference>
                    <span>用户</span>
                  </template>
                  <el-menu mode="vertical" :default-active="activeUser" >
                    <el-menu-item v-for="(item, index) in state.dataList" :index="item.name"  @click="clickUser(item)">{{
                      item.name }}</el-menu-item>
                  </el-menu>
                </el-popover>
              </el-dropdown-item>
              <el-dropdown-item @click="logout">退出</el-dropdown-item>
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
}
</style>
