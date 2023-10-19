<template>
  <div>
    <div class="option">
      <el-button type="primary" icon="Plus" @click="handleAdd('新增', '')">新增Type</el-button>
    </div>
    <el-card>
      <el-table v-loading="loading" stripe :data="dataList">
        <el-table-column label="ID" align="center" prop="id" />
        <el-table-column label="name" align="center" prop="name" />
        <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
          <template #default="{ row }">
            <el-button size="small" type="primary" link icon="Edit" @click="viewActions(row)">action管理
            </el-button>
            <el-button size="small" type="primary" link icon="Edit" @click="viewResourecs(row)">资源管理
            </el-button>
            <el-button size="small" type="primary" link icon="Edit" @click="viewRoles(row)">角色管理
            </el-button>
            <el-button size="small" type="primary" link icon="Edit" @click="handleAdd('分配', row)">角色分配
            </el-button>
            <el-button size="small" type="primary" link icon="Delete" @click="handleDelete(row)"
              :loading="row.deleteLoading">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加或角色分配岗位对话框 -->
    <el-dialog :title="`${state.open === Status.ADD ? '新增' : '角色分配'}`" titleIcon="modify" v-model="visible" width="500px"
      append-to-body :before-close="cancel">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item v-if="state.open == Status.ADD">
          <el-form-item label="name" prop="name">
            <el-input v-model="form.name" placeholder="请输入 name" />
          </el-form-item>
        </el-form-item>

        <el-form-item label="角色" prop="region" v-if="state.open == Status.EDIT">
          <el-select v-model="form.region" placeholder="请选择角色">
            <el-option v-for="item in roleOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户" prop="user" v-if="state.open == Status.EDIT">
          <el-select v-model="form.user" multiple placeholder="请选择用户" @change="changeSelect">
            <el-checkbox v-model="checked1" label="全选" size="large" @change="selectAll" />
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
import { ElForm, ElInput, ElMessage, ElMessageBox, ElSelect } from 'element-plus';
import { getRoles } from '~~/api/client/resource-type/roles'
import { getTypes, saveType, updateType, delType } from '~/api/client/resource-type/types'
import { getClientUsers } from '~~/api/common'
import { addUser } from "~~/api/client/resource-type/roles-permision"
const tenant = useTenant()
const route = useRoute()
// 角色用户名
const roleOptions = ref<SelectOption[]>([])
const userOptions = ref<SelectOption[]>([])
const { clientId } = route.params
// 全选按钮
const checked1 = ref(false)
interface Form {
  id: undefined | Number,
  name: undefined | string
  user: string,
  region: undefined
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
  open: Status.CLOSE, // 0:关闭 1:新增 2:角色分配
  // 表单参数
  form: {
    user: '',
    id: undefined,
    name: undefined,
    region: undefined
  } as Form,
  // 表单校验
  rules: {
    name: [
      { required: true, message: 'client name 不能为空', trigger: 'blur' }
    ],
    region: [{ required: true, message: '请选择角色', trigger: 'change' }],
  }
})

const {
  loading,
  dataList,
  open,
  form,
  rules
} = toRefs(state)

const formRef = ref(ElForm);

const visible = computed(() => {
  return !!state.open
})

const Visible = ref(false)

/** 查询列表 */
function getList() {
  state.loading = true
  getTypes(clientId).then((res: any) => {
    state.dataList = res
  }).finally(() => {
    state.loading = false
  })
}

// 表单重置
function resetForm() {
  state.form = {
    id: undefined,
    name: undefined,
    region: undefined,
    user: undefined
  }
  formRef.value.resetFields()
  checked1.value = false
}
// 取消按钮
function cancel() {
  resetForm()
  state.open = Status.CLOSE
  Visible.value = false
}
/** 按钮操作 */
const typeid = ref('')
// 全选

function selectAll() {

  state.form.user = []
  if (checked1.value) {
    userOptions.value.map((item: any) => { state.form.user.push(item.value) })
  } else {
    state.form.user = []

  }
}

function changeSelect() {
  if (userOptions.value.length == state.form.user.length) {
    checked1.value = true
  } else {
    checked1.value = false
  }
}
function handleAdd(wordname: string, e: any) {
  typeid.value = e.id
  const { id, } = e
  if (wordname == "新增") {
    state.open = Status.ADD
  } else {
    state.open = Status.EDIT
    getRoles(clientId, id).then((res: any) => {
      roleOptions.value = res.map((item: any) => ({
        label: item.name,
        value: item.id,
        id: item.id
      }))
      roleOptions.value.map((item:any)=>{
        if(item.label=="super-admin"){
            state.form.region=item.value
          }else{
            state.form.region=roleOptions.value[0].value
          }
      })
    })
    getClientUsers(clientId).then((res: any) => {
      userOptions.value = res.map((item: any) => ({
        label: item.displayName,
        value: item.id
      }))
    })
  }
}
let updateLoading = ref(false);
/** 提交按钮 */
function submitForm() {

  formRef.value.validate((valid: boolean) => {
    if (valid) {
      updateLoading.value = true
      let { name, region, user } = state.form
      const params = { name }
      if (state.open === Status.EDIT) {
        // 全选按钮
        const obj = { userId: user };
        const params1 = Object.entries(obj)
          .flatMap(([key, values]) => values.map((value: any) => ({ [key]: value })));
        addUser(clientId, typeid.value, region, params1).then(() => {
          ElMessage({
            showClose: true,
            message: '添加成功',
            type: 'success',
          })
          cancel()
          getList()
        }).finally(() => {
          updateLoading.value = false
        })
      } else {
        saveType(clientId, params).then(() => {
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
    }
  })
}
/** 删除按钮操作 */
function handleDelete(row: any) {
  ElMessageBox.confirm(
    `是否确认删除${row.name}"`,
    'Warning',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async function () {
    row.deleteLoading = true
    await delType(clientId, row.id)
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

function viewRoles(row: any) {
  navigateTo(`/dashboard/${tenant.value}/client/${clientId}/resource-types/${row.id}/roles`)
}

function viewResourecs(row: any) {
  navigateTo(`/dashboard/${tenant.value}/client/${clientId}/resource-types/${row.id}/resources`)
}

function viewActions(row: any) {
  navigateTo(`/dashboard/${tenant.value}/client/${clientId}/resource-types/${row.id}/actions`)
}
onMounted(() => {
  getList()
})
</script>

<style lang="scss" scoped>
.option {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}
</style>
