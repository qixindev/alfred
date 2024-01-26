<script setup lang="ts">
import { User, useTenant } from "~~/composables/useUser";
import { userTenant } from "~/composables/getUser";
import { getUser } from "~/api/common";
import { ref } from "vue";
import { ElMessage } from "element-plus";
const popoverRef = ref(null);
const activeUser = ref(null);
const route = useRoute();
const routerTenant = useRouter();
//获取
const user = useState<User>("user");
interface SelectOption {
  name: string;
  id: number;
}
const state = reactive({
  dataList: <SelectOption[]>[],
});
const showLogin = () => {
  navigateTo("/login");
};

const logout = () => {
  useLogout();
};
const tenant = useTenant();
/** 用户列表 */
function getList() {
  state.dataList = [...userTenant.value];
  //默认第一个
  tenant.value = localStorage.getItem("tenantValue")
    ? localStorage.getItem("tenantValue")
    : state.dataList?.[0].name;
  localStorage.setItem("tenantValue", tenant.value);
  // 高亮
  activeUser.value = localStorage.getItem("tenantValue");
}
watch(
  () => userTenant,
  () => {
    getList();
  },
  {
    immediate: true,
    deep: true,
  }
);

// 切换列表
function clickUser(row: any) {
  //点击用户关闭
  popoverRef.value.hide();
  tenant.value = row.name;
  localStorage.setItem("tenantValue", tenant.value as string);
  let arr = route.path.split("/");
  arr.splice(2, 1, tenant.value);
  arr.join("/");
  if (route.path == "/") {
    navigateTo("/");
  } else {
    //非主页
    navigateTo(arr.join("/"));
  }
}
onMounted(() => {
  getUser()
    .then((res: any) => {
      if (!res) {
        ElMessage({
          message: "当前没有租户，请创建租户",
          type: "error",
        });
      }
      userTenant.value = [...res];
    })
    .finally(() => {});
});
// 监听当前路由
watch(
  () => routerTenant.currentRoute.value,
  (newValue: any) => {
    // 高亮
    activeUser.value = localStorage.getItem("tenantValue");
  },
  { immediate: true }
);
// const avatar = ref('https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg')
</script>

<template>
  <div class="navbar">
    <div class="navbar-box">
      <div class="center"></div>
      <div class="right">
        <el-dropdown trigger="click" v-if="user">
          <div>
            <span>{{ user.username }}</span>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>
                <el-popover
                  ref="popoverRef"
                  :hide-after="0"
                  placement="left-start"
                  trigger="hover"
                  :offset="15"
                >
                  <template #reference>用户 </template>
                  <el-menu mode="vertical" :default-active="activeUser">
                    <el-menu-item
                      v-for="(item, index) in state.dataList"
                      :index="item.name"
                      @click="clickUser(item)"
                      >{{ item.name }}</el-menu-item
                    >
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
</template>

<style scoped lang="scss">
.navbar {
  display: flex;
  width: 100%;
  justify-content: center;
  background-color: #f7f8fa;
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
  }
}
</style>
