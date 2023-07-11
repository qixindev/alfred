<template>
  <div>
    <div class="option">
      <el-select v-model="query.role" placeholder="请选择角色" @change="getList">
        <el-option
          v-for="item in roleOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
          />
      </el-select>
      <el-button type="primary" @click="handleAdd">角色分配</el-button>
    </div>
    <el-card>
      <el-table v-loading="loading" stripe :data="dataList">
        <el-table-column label="ID"  align="center" prop="id"/>
        <el-table-column label="用户" align="center" prop="user">
          <template #default="{ row }">
            <!-- {{ userNameFilter(row.id) }} -->
            {{ row.displayName }}
          </template>
        </el-table-column>
        <el-table-column label="角色" align="center" prop="role">
          <template #default="{ row }">
            <!-- {{ roleFilter(row.roleId) }} -->
            {{ row.roleName }}
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
          <template #default="{ row }">
            <el-button size="small" type="primary" link icon="Delete" @click="handleDelete(row)"
              :loading="row.deleteLoading">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加或修改岗位对话框 -->
    <el-dialog title="角色分配" titleIcon="modify" v-model="visible" width="500px" append-to-body :before-close="cancel">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色">
            <el-option v-for="item in roleOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户" prop="name">
          <el-select v-model="form.user"  multiple    placeholder="请选择用户">
            <el-option v-for="item in userOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="submitForm" :loading="updateLoading">确 定</el-button>
        <el-button @click="cancel">取 消</el-button>
      </template>
    </el-dialog>
  </div>
</template>


<script lang="ts" setup name="Users">
import { ElForm, ElMessage, ElMessageBox } from 'element-plus';

import { getRoles } from '~~/api/client/resource-type/roles'
import { getUsers, saveUser, delUser } from '~~/api/client/resource-type/roles-permision'
import { getClientUsers } from '~~/api/common'

const route = useRoute()
const { clientId, type, resource } = route.params as any

interface Form {
  user: string,
  role: string
}

interface SelectOption {
  label: string,
  value: string,
  id?: string,
}

enum Status {
  CLOSE = 0,
  ADD = 1,
  EDIT = 2
}

const state = reactive({
  // 遮罩层
  loading: false,
  dataList: [],
  // 是否显示弹出层
  open: Status.CLOSE, // 0:关闭 1:新增 2:修改
  // 表单参数
  form: {
    user: '',
    role: ''
  } as Form,
  // 表单校验
  rules: {
    user: [
      { required: true, message: '请选择用户', trigger: 'change' }
    ],
    role: [
      { required: true, message: '请选择角色', trigger: 'change' }
    ]
  },
  query: {
    role: ''
  }
})

const {
  loading,
  dataList,
  open,
  form,
  rules,
  query
} = toRefs(state)

const roleOptions = ref<SelectOption[]>([])
const userOptions = ref<SelectOption[]>([])
const formRef = ref(ElForm);

const visible = computed(() => {
  return !!state.open
})

const viewDialogVisible = ref(false)

/** 查询列表 */
function getList() {
  if (!state.query.role) return
  state.loading = true
  getUsers(clientId, type, resource, state.query.role).then((res: any) => {
    state.dataList = res
  }).finally(() => {
    state.loading = false
  })
}

function getRoleOptions() {
  getRoles(clientId, type).then((res: any) => {
    roleOptions.value = res.map((item: any) => ({
      label: item.name,
      value: item.id,
      id: item.id
    }))
  })
}

function getUserOptions() {
  getClientUsers(clientId).then((res: any) => {
    userOptions.value = res.map((item: any) => ({
      label: item.username,
      value: item.id
    }))
  })
}

// 表单重置
function resetForm() {
  state.form = {
    user: '',
    role: ''
  }
  formRef.value.resetFields()
}
// 取消按钮
function cancel() {
  resetForm()
  state.open = Status.CLOSE
  viewDialogVisible.value = false
}
/** 新增按钮操作 */
function handleAdd() {
  state.open = Status.ADD
}


let updateLoading = ref(false);
/** 提交按钮 */
function submitForm() {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      updateLoading.value = true
      let { user, role } = state.form
      // const params = [{ userId: user }]
      const obj = { userId: user };
      const params = Object.entries(obj)
        .flatMap(([key, values]) => values.map((value:any) => ({ [key]: value })));
      saveUser(clientId, type, resource, role, params).then(() => {
        ElMessage({
          showClose: true,
          message: '创建成功',
          type: 'success',
        })
        cancel()
        getList()
      }).finally(() => {
        updateLoading.value = false
      })
    }
  })
}
/** 删除按钮操作 */
function handleDelete(row: any) {
  ElMessageBox.confirm(
    `是否确认删除用户${row.userId}"`,
    'Warning',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async function () {
    row.deleteLoading = true
    await delUser(clientId, type, resource, state.query.role, row.sub)
    row.deleteLoading = false
    getList()
    ElMessage({
      showClose: true,
      message: '删除成功',
      type: 'success',
    })
  }).catch(() => {
  })
}

const userNameFilter = computed(function () {
  return function (id: string) {
    return userOptions.value.find(item => item.value == id)?.label
  }
})

const roleFilter = computed(function () {
  return function (id: string) {
    return roleOptions.value.find(item => item.id == id)?.label
  }
})


onMounted(() => {
  getRoleOptions()
  getUserOptions()
})
</script>

<style lang="scss" scoped>
.option {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}
</style>
