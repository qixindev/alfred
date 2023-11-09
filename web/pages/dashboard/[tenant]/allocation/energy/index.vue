<script lang="ts" setup name="Users">
import { Monitor, Iphone, UploadFilled, Delete, Plus } from "@element-plus/icons-vue";
import { genFileId } from "element-plus";
import type { UploadInstance, UploadProps, UploadRawFile } from "element-plus";
import { getEnergy, putEnergy, getProto, putProto } from "~/api/energy";
import Login from "@/components/Login.vue";
import dayjs from "dayjs";
import { ref } from "vue";
import { ElMessage } from "element-plus";
const emit = defineEmits(["child-click"], ["child-primary"]);
const privacyWrite = ref("");
const privacyWord = ref("");
const wordCen = ref("");
const wordlink = ref("");
const resignSwitch = ref(true);
const loginSwitch = ref(false);
const tablePri = ref([]);
const primaryCon = ref([]);
const addPrivacy = () => {
  tablePri.value.push({
    resignSwitch: false,
    loginSwitch: false,
    privacyWrite: "",
    privacyWord: "",
  });
};
const getInfo = () => {
  getEnergy().then((res: any) => {
    tableData.value = [...res.bottom];
  });
  getProto().then((res: any) => {
    tablePri.value = [...res];
  });
};
getInfo();
function changeColor(e) {
  bgcolor.value = e;
}
function mainCss(e) {
  cssWrite.value = e;
}
const deletePri = (index: number) => {
  tablePri.value.splice(index, 1);
};
const addPri = () => {
  putProto(tablePri.value).finally(() => {
    ElMessage({
      message: "保存协议成功",
      type: "success",
    });
  });
};

const tableData = ref([]);
const deleteRow = (index: number) => {
  tableData.value.splice(index, 1);
};
const onAddItem = () => {
  tableData.value.push({
    wordCen: "",
    wordlink: "",
  });
};
const cellHover = () => {
  emit("child-click", tableData.value);
};
const cellPri = () => {
  emit("child-primary", tablePri.value);
};
</script>
<template>
  <p class="bg">登录注册协议</p>
  <el-button
    class="mt-4"
    style="width: 100px; margin: 0 30px 30px 30px; float: right"
    @click="addPrivacy"
    type="primary"
    link
    ><el-icon><Plus /></el-icon>添加协议</el-button
  >

  <el-table :data="tablePri" style="width: 90%" max-height="200" @change="cellPri">
    <el-table-column label="显示位置" width="230">
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
    <el-table-column label="勾选文字的协议内容" width="300">
      <template #default="scope">
        <el-input v-model="scope.row.privacyWord" :rows="1" type="textarea" />
      </template>
    </el-table-column>
    <el-table-column fixed="right" label="删除" width="120">
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
  <el-button
    type="primary"
    size="large"
    @click.prevent="addPri"
    style="width: 100px; margin: 20px 20px 0 20px"
  >
    保存协议
  </el-button>

  <p class="bg">底部导航栏</p>
  <el-button
    class="mt-4"
    style="width: 100px; margin: 0 30px 30px 30px; float: right"
    @click="onAddItem"
    type="primary"
    link
    ><el-icon><Plus /></el-icon>添加</el-button
  >

  <el-table :data="tableData" style="width: 50%" max-height="140" @change="cellHover">
    <el-table-column prop="wordCen" label="文字内容" width="230">
      <template #default="scope">
        <el-input v-model="scope.row.wordCen"> </el-input>
      </template>
    </el-table-column>
    <el-table-column prop="wordlink" label="链接地址" width="230">
      <template #default="scope">
        <el-input v-model="scope.row.wordlink"> </el-input> </template
    ></el-table-column>
    <el-table-column fixed="right" label="删除" width="120">
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
