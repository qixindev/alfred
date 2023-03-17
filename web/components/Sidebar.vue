<script setup lang="ts">
const route = useRoute()

const router = ref([
  {
    label: '主页',
    name: 'home',
    path: '/'
  },
  {
    name: 'client',
    label: 'client管理',
    path: '/client'
  },
  {
    name: 'providers',
    label: 'providers管理',
    path: '/providers'
  },
  {
    name: 'userManage',
    label: '用户管理',
    path: '/userManage'
  },
  {
    name: 'device',
    label: '设备管理',
    path: '/device'
  },
  {
    name: 'groups',
    label: '用户组',
    path: '/groups'
  },
  {
    name: 'groups',
    label: '租户管理',
    path: '/tenant'
  },
])

const currentIndex = ref(router.value.findIndex(item => item.path === route.path))

const handleClick = (index: number) => {
  currentIndex.value = index
  navigateTo(router.value[index].path)
}

</script>

<template>
  <div class="sidebar">
    <div class="top">
      <nuxt-link to="/">
        Logo
      </nuxt-link>
    </div>
    <div class="menu">
      <div 
        class="menu-item" 
        :class="{ 'active': currentIndex === index}" 
        v-for="(item,index) in router" 
        :key="item.name"
        @click="handleClick(index)"
        >
        {{ item.label }}
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.sidebar {
  background-color: #FFF;
  min-height: 100vh;
  .menu {

  .menu-item {
    height: 48px;
    line-height: 48px;
    cursor: pointer;
    &.active {
      background-color: #409EFF;
      color: #FFF;
    }
  }
  }
}

</style>
