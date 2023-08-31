<script setup lang="ts">
const DELAY = 60
const counting = ref(false)
const countdownSeconds = ref(DELAY)

const buttonText = computed(() => {
  return counting.value ? `${countdownSeconds.value}秒后重新发送` : '发送验证码'
})

function startCountdown() {
  if (counting.value) return;
  countdownSeconds.value = DELAY
  counting.value = true;
  countdown();

  // 模拟倒计时结束后的操作，这里使用setTimeout来模拟实际的倒计时
  setTimeout(() => {
    counting.value = false;
  }, countdownSeconds.value * 1000);
}
function countdown() {
  if (countdownSeconds.value > 0) {
    setTimeout(() => {
      countdownSeconds.value--;
      countdown();
    }, 1000);
  } else {
    isCounting.value = false; // 倒计时结束后恢复按钮状态
  }
}

defineExpose({
  startCountdown
})

</script>

<template>
  <el-button
    :disabled="counting"
    :style="{width: '150px'}"
  >
    {{ buttonText }}
  </el-button>
</template>


<style scoped>

</style>
