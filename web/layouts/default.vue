<script setup lang="ts">
const auth = useCookie("QixinAuth");
auth.value && useGetUserInfo();
</script>

<template>
  <div class="app-wrapper">
    <Sidebar class="sidebar"></Sidebar>
    <div class="main">
      <NavBar></NavBar>
      <div class="container">
        <slot />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.app-wrapper {
  min-height: 100vh;
  background-color: #fff;
  display: flex;
  .sidebar {
    width: 200px;
    flex-shrink: 0;
  }
  .main {
    width: 100%;
    height: 100vh;
    flex: 1;
    overflow-y: auto;
    .container {
      padding: 20px;
    }
  }
}
</style>
