<script lang="ts" setup name="Users">
import { Monitor, Iphone, UploadFilled, Delete, Plus } from "@element-plus/icons-vue";
import { genFileId } from "element-plus";
import type { UploadInstance, UploadProps, UploadRawFile } from "element-plus";
import { getEnergy, getProto } from "~/api/energy";
import dayjs from "dayjs";
import { ref } from "vue";
import { ElMessage } from "element-plus";
import { useRouter, useRoute } from "vue-router";
const emit = defineEmits(["child-click"], ["child-primary"]);
const tenant = computed(() => useTenant().value);
const route = useRoute();
const { state: tanent } = route.query as any;
let currentTenant =
  route.path.substring(0, 10) == "/dashboard" ? tenant.value : tanent ?? "default";
const privacyWrite = ref("");
const privacyWord = ref("");
const wordCen = ref("");
const wordlink = ref("");
const resignSwitch = ref(true);
const loginSwitch = ref(false);
const tablePri = ref([]);
const primaryCon = ref([]);
const tableData = ref([]);

const addPrivacy = () => {
  tablePri.value.push({
    resignSwitch: false,
    loginSwitch: false,
    privacyWrite: "",
    privacyWord: "",
  });
};
const getInfo = () => {
  getProto(currentTenant).then((res: any) => {
    tablePri.value = [...res];
  });
};
getInfo();
const props = defineProps({
  allInfo: {
    type: Object,
    default: {},
  },
});

const deletePri = (index: number) => {
  tablePri.value.splice(index, 1);
  emit("child-primary", tablePri.value);
};

const deleteRow = (index: number) => {
  tableData.value.splice(index, 1);
  emit("child-click", tableData.value);
};
const onAddItem = () => {
  tableData.value.push({
    wordCen: "",
    wordlink: "",
  });
};
watch(
  () => props.allInfo,
  () => {
    // 解决 is not iterable和tableData的类型
    if(props.allInfo && props.allInfo.bottom){
    tableData.value = [...props.allInfo && props.allInfo.bottom]  
    }
  },
  { immediate: true, deep: true }
);
const cellPri = () => {
  emit("child-primary", tablePri.value);
};
const cellHover = () => {
  emit("child-click", tableData.value);
};
const change1 = (word) => {
  let a = "https://";
  let b = "http://";
  if (
    (word.slice(0, 8) == a && word.slice(0, 7) != b) ||
    (word.slice(0, 8) != a && word.slice(0, 7) == b)
  ) {
  } else {
    ElMessage({
      message: "请输入以https://或者http://开头的地址",
      type: "warning",
    });
  }
};
</script>
<template>
  <div style="display: flex">
    <p class="bg">登录注册协议</p>
    <el-button
      style="width: 100px; margin: 11px 0px 30px 10px"
      @click="addPrivacy"
      type="primary"
      link
      ><el-icon><Plus /></el-icon>添加协议</el-button
    >
  </div>

  <el-table
    :data="tablePri"
    style="width: 95%; box-shadow: 2px 5px 12px rgb(0 0 0/0.2); margin-left: 5px"
    max-height="200"
    @change="cellPri"
  >
    <el-table-column label="显示位置" width="300">
      <template #default="scope">
        <el-checkbox v-model="scope.row.resignSwitch" label="注册页面" size="large" />
        <el-checkbox v-model="scope.row.loginSwitch" label="登录页面" size="large" />
      </template>
    </el-table-column>
    <el-table-column label="勾选文字" width="300">
      <template #default="scope">
        <el-input
          v-model="scope.row.privacyWrite"
          placeholder="例如：阅读隐私协议"
          style="width: 100%"
        />
      </template>
    </el-table-column>
    <el-table-column label="勾选文字的协议内容" width="330">
      <template #default="scope">
        <el-input v-model="scope.row.privacyWord" :rows="1" type="textarea" />
      </template>
    </el-table-column>
    <el-table-column fixed="right" label="删除" width="200">
      <template #default="scope">
        <el-button
          link
          type="primary"
          size="small"
          @click.prevent="deletePri(scope.$index)"
        >
          删除
        </el-button>
      </template>
    </el-table-column>
  </el-table>

  <div style="display: flex">
    <p class="bg">底部导航栏</p>
    <el-button
      class="mt-4"
      style="width: 100px; margin: 23px 0px 30px 10px"
      @click="onAddItem"
      type="primary"
      link
      ><el-icon><Plus /></el-icon>添加</el-button
    >
  </div>
  <el-table
    :data="tableData"
    style="width: 95%; box-shadow: 2px 5px 12px rgb(0 0 0/0.2); margin-left: 5px"
    max-height="150"
    @change="cellHover"
  >
    <el-table-column prop="wordCen" label="文字内容" width="430">
      <template #default="scope">
        <el-input v-model="scope.row.wordCen"> </el-input>
      </template>
    </el-table-column>
    <el-table-column prop="wordlink" label="链接地址" width="430">
      <template #default="scope">
        <el-input
          v-model="scope.row.wordlink"
          placeholder="请输入以https://或者http://开头的地址"
          @change="change1(scope.row.wordlink)"
        >
        </el-input> </template
    ></el-table-column>
    <el-table-column fixed="right" label="删除" width="230">
      <template #default="scope">
        <el-button
          link
          type="primary"
          size="small"
          @click.prevent="deleteRow(scope.$index)"
        >
          删除
        </el-button>
      </template>
    </el-table-column>
  </el-table>
</template>
<style scoped lang="scss">
.center {
  width: 90%;
  background: white;
  border-radius: 5px;
  margin-left: 2%;
  padding: 10px;
}
.bg {
  font-size: 20px;
  margin: 20px 20px 10px 0px;
}
</style>
