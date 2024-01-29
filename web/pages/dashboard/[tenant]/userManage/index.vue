<template>
  <div>
    <div class="option">
      <el-button type="primary" icon="Plus" @click="handleAdd">新增Users</el-button>
    </div>
    <el-card>
      <div class="top">
        <el-input
          v-model="serchValue"
          placeholder="Please Input"
          class="topInput"
        /><el-button type="primary" @click="getList">搜索</el-button>
      </div>
      <el-table v-loading="loading" stripe :data="dataList">
        <el-table-column label="ID" align="center" prop="id" />
        <el-table-column label="username" align="center" prop="username" />
        <el-table-column label="displayName" align="center" prop="displayName" />
        <el-table-column label="firstName" align="center" prop="firstName" />
        <el-table-column label="lastName" align="center" prop="lastName" />
        <el-table-column label="phone" align="center" prop="phone" />
        <el-table-column label="phoneVerified" align="center" prop="phoneVerified" />
        <el-table-column label="email" align="center" prop="email" />
        <el-table-column label="emailVerified" align="center" prop="emailVerified" />
        <el-table-column label="disabled" align="center" prop="disabled" />
        <el-table-column
          label="twoFactorEnabled"
          align="center"
          prop="twoFactorEnabled"
        />
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
              @click="viewGroups(row)"
              >group管理
            </el-button>
            <el-button
              size="small"
              type="primary"
              link
              icon="Edit"
              @click="handleUpdate(row)"
              >修改
            </el-button>
            <el-button
              size="small"
              type="primary"
              link
              icon="Edit"
              @click="handleUpdatePass(row)"
              >修改密码
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
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[10, 15, 20, 50, 100]"
        :small="small"
        :disabled="disabled"
        :background="background"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        class="pagination"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

    <!-- 添加或修改对话框 -->
    <el-dialog
      :title="`${
        open === Status.ADD ? '新增' : open === Status.EDIT ? '修改' : '修改密码'
      }`"
      titleIcon="modify"
      v-model="visible"
      width="800px"
      append-to-body
      :before-close="cancel"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="200px">
        <el-form-item
          label="username"
          prop="username"
          v-if="state.open === Status.ADD || state.open === Status.EDIT"
        >
          <el-input v-model="form.username" placeholder="请输入 name" />
        </el-form-item>
        <el-form-item
          label="password"
          prop="passwordHash"
          v-if="state.open === Status.ADD || state.open === Status.PASS"
        >
          <el-input
            v-model="form.passwordHash"
            placeholder="请输入 password"
            show-password
          />
        </el-form-item>
        <el-form-item
          label="displayName"
          prop="displayName"
          v-if="state.open === Status.ADD || state.open === Status.EDIT"
        >
          <el-input v-model="form.displayName" placeholder="请输入 displayName" />
        </el-form-item>
        <el-form-item
          label="firstName"
          prop="firstName"
          v-if="state.open === Status.ADD || state.open === Status.EDIT"
        >
          <el-input v-model="form.firstName" placeholder="请输入 firstName" />
        </el-form-item>
        <el-form-item
          label="lastName"
          prop="lastName"
          v-if="state.open === Status.ADD || state.open === Status.EDIT"
        >
          <el-input v-model="form.lastName" placeholder="请输入 lastName" />
        </el-form-item>
        <el-form-item
          label="phone"
          prop="phone"
          v-if="state.open === Status.ADD || state.open === Status.EDIT"
        >
          <el-input v-model="form.phone" placeholder="请输入 phone" />
        </el-form-item>
        <el-form-item
          label="email"
          prop="email"
          v-if="state.open === Status.ADD || state.open === Status.EDIT"
        >
          <el-input v-model="form.email" placeholder="请输入 email" />
        </el-form-item>
        <el-form-item
          label="disabled"
          prop="disabled"
          v-if="state.open === Status.ADD || state.open === Status.EDIT"
        >
          <el-switch v-model="form.disabled" />
        </el-form-item>
        <el-form-item
          label="twoFactorEnabled"
          prop="twoFactorEnabled"
          v-if="state.open === Status.ADD || state.open === Status.EDIT"
        >
          <el-switch v-model="form.twoFactorEnabled" />
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

import { getUsers, saveUser, updateUser, delUser, passUser } from "~/api/userManage";
const serchValue = ref("");
const currentPage = ref(1);
const pageSize = ref(10);
const small = ref(false);
const background = ref(false);
const disabled = ref(false);
const total = ref(0);
const tenant = useTenant();
interface Form {
  id: undefined | Number;
  username: undefined | string;
  passwordHash: undefined | string;
  displayName: undefined | string;
  firstName: undefined | string;
  lastName: undefined | string;
  phone: undefined | string;
  phoneVerified: false | boolean;
  email: undefined | string;
  emailVerified: false | boolean;
  disabled: false | boolean;
  twoFactorEnabled: false | boolean;
}

enum Status {
  CLOSE = 0,
  ADD = 1,
  EDIT = 2,
  PASS = 3,
}

const state = reactive({
  // 遮罩层
  loading: false,
  dataList: [],
  // 是否显示弹出层
  open: Status.CLOSE, // 0:关闭 1:新增 2:修改 3:修改密码
  // 表单参数
  form: {
    id: undefined,
    username: undefined,
    passwordHash: undefined,
    displayName: undefined,
    firstName: undefined,
    lastName: undefined,
    phone: undefined,
    phoneVerified: false,
    email: undefined,
    emailVerified: false,
    disabled: false,
    twoFactorEnabled: false,
  } as Form,
  // 表单校验
  rules: {
    username: [{ required: true, message: "username 不能为空", trigger: "blur" }],
    passwordHash: [{ required: true, message: "password 不能为空", trigger: "blur" }],
    displayName: [{ required: true, message: "displayName 不能为空", trigger: "blur" }],
  },
});

const { loading, dataList, open, form, rules } = toRefs(state);

const formRef = ref(ElForm);

const visible = computed(() => {
  return !!state.open;
});

const viewDialogVisible = ref(false);
// 改变pageSize
const handleSizeChange = (val) => {
  pageSize.value = val;
  getList();
};
// 改变pageNum
const handleCurrentChange = (val) => {
  currentPage.value = val;
  getList();
};
/** 查询列表 */
function getList() {
  state.loading = true;
  getUsers({
    pageNum: currentPage.value,
    pageSize: pageSize.value,
    search: serchValue.value,
  })
    .then((res: any) => {
      state.dataList = res.data;
      total.value = res.total;
    })
    .finally(() => {
      state.loading = false;
    });
}

// 表单重置
function resetForm() {
  state.form = {
    id: undefined,
    username: undefined,
    passwordHash: undefined,
    displayName: undefined,
    firstName: undefined,
    lastName: undefined,
    phone: undefined,
    phoneVerified: false,
    email: undefined,
    emailVerified: false,
    disabled: false,
    twoFactorEnabled: false,
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
function handleUpdate(row: any) {
  const {
    id,
    username,
    passwordHash,
    displayName,
    firstName,
    lastName,
    phone,
    phoneVerified,
    email,
    emailVerified,
    disabled,
    twoFactorEnabled,
  } = row;
  state.open = Status.EDIT;
  nextTick(() => {
    state.form = {
      id,
      username,
      passwordHash,
      displayName,
      firstName,
      lastName,
      phone,
      phoneVerified,
      email,
      emailVerified,
      disabled,
      twoFactorEnabled,
    };
  });
}
/** 修改密码操作 */
function handleUpdatePass(row: any) {
  const { id, passwordHash } = row;
  state.open = Status.PASS;
  nextTick(() => {
    state.form = {
      id,
      passwordHash,
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
        username,
        passwordHash,
        displayName,
        firstName,
        lastName,
        phone,
        phoneVerified,
        email,
        emailVerified,
        disabled,
        twoFactorEnabled,
      } = state.form;

      const params = {
        id,
        username,
        passwordHash,
        displayName,
        firstName,
        lastName,
        phone,
        phoneVerified,
        email,
        emailVerified,
        disabled,
        twoFactorEnabled,
      };

      if (state.open === Status.EDIT) {
        updateUser(id as number, params)
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
      } else if (state.open === Status.ADD) {
        saveUser(params)
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
      } else {
        passUser(id, { passwordHash: passwordHash })
          .then(() => {
            ElMessage({
              showClose: true,
              message: "修改密码成功",
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
      await delUser(row.id);
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

function viewGroups(row: any) {
  navigateTo(`/dashboard/${tenant.value}/userManage/${row.id}/groups`);
}

onMounted(() => {
  getList();
});
</script>

<style lang="scss" scoped>
.top {
  display: flex;
}
.topInput {
  width: 250px;
  margin-right: 50px;
}
.option {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}
.pagination {
  margin: 10px;
  float: right;
}
</style>
