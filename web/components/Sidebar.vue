<script setup lang="ts">
import "element-plus/dist/index.css";
import * as ELIcons from "@element-plus/icons-vue";
// import { Tenant, usePath } from '~~/composables/useUser'
const route = useRoute();

const tenant = useTenant();
// const tenant = useState<Tenant>('tenant')
const routerTenant = useRouter();
const router = ref([
  {
    label: "主页",
    name: "home",
    path: "/",
    icon: "home-filled",
  },
  {
    name: "client",
    label: "client管理",
    path: "/client",
  },
  {
    name: "providers",
    label: "providers管理",
    path: "/providers",
  },
  {
    name: "userManage",
    label: "用户管理",
    path: "/userManage",
  },
  {
    name: "device",
    label: "设备管理",
    path: "/device",
  },
  {
    name: "groups",
    label: "用户组",
    path: "/groups",
  },
  {
    name: "tenant",
    label: "租户管理",
    path: "/tenant",
  },
  {
    name: "allocation",
    label: "配置页面",
    path: "/allocation",
  },
]);

const currentIndex = ref(router.value.findIndex((item) => item.path === route.path));
const handleClick = (index: number, item: any) => {
  currentIndex.value = index;
  if (index == 0) {
    navigateTo(router.value[index].path);
  } else {
    //非主页
    navigateTo(`/dashboard/${tenant.value}${item.path}`);
  }
  // navigateTo(router.value[index].path)
};
// 监听当前路由
watch(
  () => routerTenant.currentRoute.value,
  (newValue: any) => {
    currentIndex.value = router.value.findIndex(
      (item) => item.name === newValue.fullPath.split("/")[3]
    );
  },
  { immediate: true }
);
</script>

<template>
  <div class="sidebar">
    <div class="top">
      <nuxt-link to="/"> Logo </nuxt-link>
    </div>
    <div class="menu">
      <div
        class="menu-item"
        :class="{ active: currentIndex === index }"
        v-for="(item, index) in router"
        :key="item.name"
        @click="handleClick(index, item)"
      >
        {{ item.label }}
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.sidebar {
  background-color: #fff;
  min-height: 100vh;

  .menu {
    .menu-item {
      height: 48px;
      line-height: 48px;
      cursor: pointer;

      &.active {
        background-color: #409eff;
        color: #fff;
      }
    }
  }
}
</style>
