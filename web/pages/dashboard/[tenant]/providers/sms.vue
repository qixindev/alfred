<template>
  <div>
    <div class="option">
      <el-button type="primary" @click="handleAdd">新增SMS</el-button>
    </div>
    <el-card>
      <el-table v-loading="loading" stripe :data="dataList">
        <el-table-column label="ID" align="center" prop="id" />
        <el-table-column label="name" align="center" prop="name" />
        <el-table-column label="type" align="center" prop="type" />
        <el-table-column
          label="操作"
          align="center"
          class-name="small-padding fixed-width"
        >
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              link
              icon="Edit"
              @click="handleUpdate(row.id)"
              >修改
            </el-button>
            <el-button
              size="small"
              type="primary"
              link
              icon="Delete"
              @click="handleDelete(row)"
              :loading="row.deleteLoading"
              >删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加或修改对话框 -->
    <el-dialog
      :title="`${open === Status.ADD ? '新增' : '修改'}`"
      titleIcon="modify"
      v-model="visible"
      width="500px"
      append-to-body
      :before-close="cancel"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="name" prop="name">
          <el-input v-model="form.name" placeholder="请输入name" />
        </el-form-item>
        <el-form-item label="type" prop="type">
          <el-select v-model="form.type" placeholder="请选择type">
            <el-option label="tclould" value="tcloud" />
            <el-option label="alibaba" value="alibaba" />
          </el-select>
        </el-form-item>
        <el-form-item
          label="accessKeyId"
          prop="accessKeyId"
          v-if="form.type == 'alibaba'"
        >
          <el-input v-model="form.accessKeyId" placeholder="请输入accessKeyId" />
        </el-form-item>
        <el-form-item
          label="accessKeySecret"
          prop="accessKeySecret"
          v-if="form.type == 'alibaba'"
        >
          <el-input v-model="form.accessKeySecret" placeholder="请输入accessKeySecret" />
        </el-form-item>
        <el-form-item label="regionId" prop="regionId" v-if="form.type == 'alibaba'">
          <el-input v-model="form.regionId" placeholder="请输入regionId" />
        </el-form-item>
        <el-form-item label="endpoint" prop="endpoint" v-if="form.type == 'alibaba'">
          <el-input v-model="form.endpoint" placeholder="请输入endpoint" />
        </el-form-item>
        <el-form-item
          label="templateCode"
          prop="templateCode"
          v-if="form.type == 'alibaba'"
        >
          <el-input v-model="form.templateCode" placeholder="请输入templateCode" />
        </el-form-item>
        <el-form-item label="region" prop="region" v-if="form.type == 'tcloud'">
          <el-input v-model="form.region" placeholder="请输入region" />
        </el-form-item>
        <el-form-item label="sdkAppId" prop="sdkAppId" v-if="form.type == 'tcloud'">
          <el-input v-model="form.sdkAppId" placeholder="请输入sdkAppId" />
        </el-form-item>
        <el-form-item label="secretId" prop="secretId" v-if="form.type == 'tcloud'">
          <el-input v-model="form.secretId" placeholder="请输入secretId" />
        </el-form-item>
        <el-form-item label="secretKey" prop="secretKey" v-if="form.type == 'tcloud'">
          <el-input v-model="form.secretKey" placeholder="请输入secretKey" />
        </el-form-item>
        <el-form-item
          label="signName"
          prop="signName"
          v-if="form.type == 'tcloud' || form.type == 'alibaba'"
        >
          <el-input v-model="form.signName" placeholder="请输入signName" />
        </el-form-item>
        <el-form-item label="templateId" prop="templateId" v-if="form.type == 'tcloud'">
          <el-input v-model="form.templateId" placeholder="请输入templateId" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="submitForm" :loading="updateLoading"
          >确 定</el-button
        >
        <el-button @click="cancel">取 消</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup name="Users">
import { ElForm, ElInput, ElMessage, ElMessageBox } from "element-plus";

import { getSms, saveSms, updateSms, delSms, getSmsById } from "~/api/providers/sms";

interface Form {
  id?: Number;
  name: undefined | string;
  region: undefined | string;
  secretId: undefined | string;
  sdkAppId: undefined | string;
  secretKey: undefined | string;
  templateId: undefined | string;
  type: undefined | string;
  accessKeyId: undefined | string;
  accessKeySecret: undefined | string;
  regionId: undefined | string;
  endpoint: undefined | string;
  templateCode: undefined | string;
  signName: undefined | string;
}

enum Status {
  CLOSE = 0,
  ADD = 1,
  EDIT = 2,
}

const state = reactive({
  // 遮罩层
  loading: false,
  dataList: [],
  // 是否显示弹出层
  open: Status.CLOSE, // 0:关闭 1:新增 2:修改
  // 表单参数
  form: {
    id: undefined,
    name: undefined,
    region: undefined,
    secretId: undefined,
    secretKey: undefined,
    templateId: undefined,
    type: undefined,
    accessKeyId: undefined,
    accessKeySecret: undefined,
    regionId: undefined,
    endpoint: undefined,
    templateCode: undefined,
    signName: undefined,
  } as Form,
  // 表单校验
  rules: {
    name: [{ required: true, message: "client name 不能为空", trigger: "blur" }],
    region: [{ required: true, message: "region 不能为空", trigger: "blur" }],
    sdkAppId: [{ required: true, message: "sdkAppId 不能为空", trigger: "blur" }],
    secretId: [{ required: true, message: "secretId 不能为空", trigger: "blur" }],
    secretKey: [{ required: true, message: "secretKey 不能为空", trigger: "blur" }],
    signName: [{ required: true, message: "signName 不能为空", trigger: "blur" }],
    templateId: [{ required: true, message: "templateId 不能为空", trigger: "blur" }],
    type: [{ required: true, message: "请选择type", trigger: "change" }],
    accessKeyId: [{ required: true, message: "请选择accessKeyId", trigger: "change" }],
    accessKeySecret: [
      { required: true, message: "请选择accessKeySecret", trigger: "change" },
    ],
    regionId: [{ required: true, message: "请选择regionId", trigger: "change" }],
    endpoint: [{ required: true, message: "请选择endpoint", trigger: "change" }],
    templateCode: [{ required: true, message: "请选择templateCode", trigger: "change" }],
  },
});

const { loading, dataList, open, form, rules } = toRefs(state);

const formRef = ref(ElForm);

const visible = computed(() => {
  return !!state.open;
});

const viewDialogVisible = ref(false);

/** 查询列表 */
function getList() {
  state.loading = true;
  getSms()
    .then((res: any) => {
      state.dataList = res;
    })
    .finally(() => {
      state.loading = false;
    });
}

// 表单重置
function resetForm() {
  state.form = {
    id: undefined,
    name: undefined,
    region: undefined,
    secretId: undefined,
    sdkAppId: undefined,
    secretKey: undefined,
    signName: undefined,
    templateId: undefined,
    type: undefined,
    accessKeyId: undefined,
    accessKeySecret: undefined,
    regionId: undefined,
    endpoint: undefined,
    templateCode: undefined,
  };
  formRef.value.resetFields();
}
// 取消按钮
function cancel() {
  resetForm();
  state.open = Status.CLOSE;
  viewDialogVisible.value = false;
}
/** 新增按钮操作 */
function handleAdd() {
  state.open = Status.ADD;
}
/** 修改按钮操作 */
async function handleUpdate(id: number) {
  const data: any = await getSmsById(id);
  const {
    smsConnector,
    region,
    secretId,
    sdkAppId,
    secretKey,
    signName,
    templateId,
    accessKeyId,
    accessKeySecret,
    regionId,
    endpoint,
    templateCode,
  } = data;
  state.open = Status.EDIT;
  nextTick(() => {
    state.form = {
      id,
      name: smsConnector.name,
      region,
      secretId,
      sdkAppId,
      secretKey,
      signName,
      templateId,
      type: smsConnector.type,
      accessKeyId,
      accessKeySecret,
      regionId,
      endpoint,
      templateCode,
    };
  });
}

let updateLoading = ref(false);
/** 提交按钮 */
function submitForm() {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      updateLoading.value = true;
      let {
        id,
        name,
        region,
        secretId,
        sdkAppId,
        secretKey,
        signName,
        templateId,
        type,
        accessKeyId,
        accessKeySecret,
        endpoint,
        regionId,
        templateCode,
      } = state.form;

      const paramsTcloud = {
        name,
        region,
        secretId,
        sdkAppId,
        secretKey,
        signName,
        templateId,
        type,
      };
      const paramsAlibaba = {
        name,
        accessKeyId,
        accessKeySecret,
        endpoint,
        regionId,
        templateCode,
        signName,
        type,
      };
      if (type == "tcloud") {
        if (state.open === Status.EDIT) {
          updateSms(id as number, paramsTcloud)
            .then(() => {
              ElMessage({
                showClose: true,
                message: "修改成功",
                type: "success",
              });
              cancel();
              getList();
            })
            .finally(() => {
              updateLoading.value = false;
            });
        } else {
          saveSms(paramsTcloud)
            .then(() => {
              ElMessage({
                showClose: true,
                message: "创建成功",
                type: "success",
              });
              cancel();
              getList();
            })
            .finally(() => {
              updateLoading.value = false;
            });
        }
      } else {
        if (state.open === Status.EDIT) {
          updateSms(id as number, paramsAlibaba)
            .then(() => {
              ElMessage({
                showClose: true,
                message: "修改成功",
                type: "success",
              });
              cancel();
              getList();
            })
            .finally(() => {
              updateLoading.value = false;
            });
        } else {
          saveSms(paramsAlibaba)
            .then(() => {
              ElMessage({
                showClose: true,
                message: "创建成功",
                type: "success",
              });
              cancel();
              getList();
            })
            .finally(() => {
              updateLoading.value = false;
            });
        }
      }
    }
  });
}
/** 删除按钮操作 */
function handleDelete(row: any) {
  ElMessageBox.confirm(`是否确认删除${row.name}"`, "Warning", {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(async function () {
      row.deleteLoading = true;
      await delSms(row.id);
      row.deleteLoading = false;
      getList();
      ElMessage({
        showClose: true,
        message: "删除成功",
        type: "success",
      });
    })
    .catch(() => {});
}

onMounted(() => {
  getList();
});
</script>

<style lang="scss" scoped>
.option {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}
</style>
