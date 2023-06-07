<template>
  <div>
    <div class="option">
      <el-button type="primary" icon="Plus" @click="handleAdd">新增Action</el-button>
    </div>
    <el-card>
      <el-table v-loading="loading" stripe :data="dataList">
        <el-table-column label="ID" width="80px" align="center" prop="id"/>
        <el-table-column label="actionName" align="center" prop="actionName" />
        <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
          <template #default="{ row }">
            <!-- <el-button size="small" type="primary" link icon="Edit" @click="handleUpdate(row)">修改
            </el-button> -->
            <el-button size="small" type="primary" link icon="Delete" @click="handleDelete(row)" :loading="row.deleteLoading">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加或修改岗位对话框 -->
    <el-dialog :title="`${open === Status.ADD ? '新增' : '修改'}`" titleIcon="modify" v-model="visible" width="500px" append-to-body
      :before-close="cancel">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="action" prop="actionId">
          <el-select v-model="form.actionId" placeholder="请选择action">
            <el-option
              v-for="item in actionOption"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
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
import { ElForm, ElInput, ElMessage, ElMessageBox } from 'element-plus';

import { getRoleActions, saveRoleAction, updateRoleAction, delRoleAction } from '~/api/client/resource-type/roles/action'
import { getActions } from '~/api/client/resource-type/action'

const route = useRoute()
const { clientId, type, roleName } = route.params

interface Form {
  id: undefined | Number,
  actionId: undefined | String
}

enum Status {
  CLOSE = 0,
  ADD = 1,
  EDIT = 2
}

const actionOption = ref<any>([])

const state = reactive({
  // 遮罩层
  loading: false,
  dataList: [],
  // 是否显示弹出层
  open: Status.CLOSE, // 0:关闭 1:新增 2:修改
  // 表单参数
  form: {
    id: undefined,
    actionId: undefined,
  } as Form,
  // 表单校验
  rules: {
    name: [
      { required: true, message: 'actionId 不能为空', trigger: 'blur' }
    ]
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

const viewDialogVisible = ref(false)

/** 查询列表 */
function getList() {
  state.loading = true
  getRoleActions(clientId, type, roleName).then((res:any) => {
    state.dataList = res
  }).finally(() => {
    state.loading = false
  })
}

// 获取action列表
async function getActionsList() {
  const actionList: any = await getActions(clientId,type)
  actionOption.value = actionList.map((item: any) => ({
    label: item.name,
    value: item.id
  }))
}

// 表单重置
function resetForm() {
  state.form = {
    id: undefined,
    actionId: undefined,
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
/** 修改按钮操作 */
function handleUpdate(row: any) {
  const {id, actionId } = row
  state.open = Status.EDIT
  nextTick(()=>{
    state.form = {
      id,
      actionId,
    }
  })
}

let updateLoading = ref(false);
/** 提交按钮 */
function submitForm() {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      updateLoading.value = true
      let { id, actionId } = state.form

      const params = [{ actionId }]

      if (state.open === Status.EDIT) {
        updateRoleAction(clientId, type, roleName, id as number, params).then(() => {
          ElMessage({
            showClose: true,
            message: '修改成功',
            type: 'success',
          })
          cancel()
          getList()
        }).finally(() => {
          updateLoading.value = false
        })
      } else {
        saveRoleAction(clientId, type, roleName, params).then(() => {
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
    `是否确认删除${row.actionName}"`,
    'Warning',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async function () {
    row.deleteLoading = true
    await delRoleAction(clientId, type, roleName, row.actionName)
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

onMounted(() => {
  getList()
  getActionsList()
})
</script>

<style lang="scss" scoped>
.option {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}
</style>
