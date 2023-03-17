<script setup lang="ts">
const user = useState('user')

</script>

<template>
  <el-card class="profile-box">
    <div class="attr-item">
      用户名：{{ user.username }}
    </div>
    <div class="attr-item">
      phone：{{ user.phone }}
    </div>
    <div class="attr-item">
      email：{{ user.email }}
    </div>
  </el-card>
</template>

<style scoped lang="scss">
.profile-box {
  .attr-item {
    margin: 10px 0;
  }
}
</style>
