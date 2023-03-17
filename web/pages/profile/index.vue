<script setup lang="ts">
const activeIndex = ref(0)

const tabs = [
  {
    label: '基本资料',
    type: 'Profile',
    component: markRaw(resolveComponent('Profile') as Component)
  },
  {
    label: '我的站点',
    type: 'Tenant',
    component: markRaw(resolveComponent('Tenant') as Component)
  }
]

const currentComponent = ref(tabs[0].component)


const tabChange = (e: any) => {
  const { index } = e.target.dataset
  const currentTab = tabs[index]
  activeIndex.value = index
  currentComponent.value = currentTab.component
}

definePageMeta({
    middleware: ["login"]
})
</script>

<template>
  <div class="container">
    <div class="aside">
      <div class="tabs" @click="tabChange">
        <div class="tab" 
          v-for="(item,index) in tabs" 
          :class="activeIndex == index ? 'active' :''"
          :data-index="index"
          >
          {{item.label}}
        </div>
      </div>
    </div>
    <div class="main">
      <!-- <Tenant v-if="activeName === 'tenant'" /> -->
      <component  v-if="currentComponent" :is="currentComponent" />
    </div>
  </div>
</template>

<style scoped lang="scss">
.container {
  display: flex;
  justify-content: center;
  padding: 30px 0;
  .aside {
    width: 200px;
    height: 300px;
    border-radius: 4px;
    background-color: #FFF;
    margin-right: 20px;
    .tabs {
      .tab {
        padding: 12px;
        cursor: pointer;
        &:hover {
          background-color: rgba(238, 245, 255);
        }
        &.active {
          color: #60A5FA;
          background-color: #e5e7eb;
        }
      }
    }
  }
  .main {
    width: 1000px;
  }
}
</style>
