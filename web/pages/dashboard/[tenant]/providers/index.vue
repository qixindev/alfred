<template>
  <div>
    <div class="option">
      <el-button type="primary" @click="handleConfig(row)">SMS配置</el-button>
      <el-button type="primary" @click="handleAdd">新增Providers</el-button>
    </div>
    <el-card>
      <el-table v-loading="loading" stripe :data="dataList">
        <el-table-column label="ID"  align="center" prop="id"/>
        <el-table-column label="name" align="center" prop="name" />
        <el-table-column label="type" align="center" prop="type" />
        <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
          <template #default="{ row }">
            <el-button size="small" type="primary" link icon="Edit" @click="handleUpdate(row)">修改
            </el-button>
            <el-button size="small" type="primary" link icon="Delete" @click="handleDelete(row)" :loading="row.deleteLoading">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加或修改对话框 -->
    <el-dialog :title="`${open === Status.ADD ? '新增' : '修改'}`" titleIcon="modify" v-model="visible" width="500px" append-to-body
      :before-close="cancel">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="name" prop="name">
          <el-input v-model="form.name" placeholder="请输入name" />
        </el-form-item>
        <el-form-item label="type" prop="type">
          <el-select v-model="form.type" placeholder="请选择type">
            <el-option
              v-for="item in typeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.type !== 'sms'" label="agentId" prop="agentId">
          <el-input v-model="form.agentId" placeholder="请输入agentId" />
        </el-form-item>
        <el-form-item v-if="form.type !== 'sms'" label="appSecret" prop="appSecret">
          <el-input v-model="form.appSecret" placeholder="请输入appSecret" />
        </el-form-item>

        <!-- 钉钉参数 -->
        <el-form-item label="appKey" prop="appKey" v-if="form.type === 'dingtalk'">
          <el-input v-model="form.appKey" placeholder="请输入appKey" />
        </el-form-item>

        <!-- 企微参数 -->
        <el-form-item label="corpId" prop="corpId" v-if="form.type === 'wecom'">
          <el-input v-model="form.corpId" placeholder="请输入corpId" />
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

import { getProviders, saveProvider, updateProvider, delProvider, getProvider } from '~/api/providers'

interface Form {
  id: undefined | Number,
  name: undefined | string,
  type: undefined | string,
  agentId?: undefined | string,
  appSecret?: undefined | string,
  corpId?: undefined | string,
  appKey?: undefined | string
}

enum Status {
  CLOSE = 0,
  ADD = 1,
  EDIT = 2
}

const tenant =  useTenant()

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
    type: undefined,
    agentId: undefined,
    appSecret: undefined,
    corpId: undefined,
    appKey: undefined
  } as Form,
  // 表单校验
  rules: {
    name: [
      { required: true, message: 'name 不能为空', trigger: 'blur' }
    ],
    type: [
      { required: true, message: 'type 不能为空', trigger: 'change' }
    ],
    agentId: [
      { required: true, message: 'agentId 不能为空', trigger: 'blur' }
    ],
    appSecret: [
      { required: true, message: 'appSecret 不能为空', trigger: 'blur' }
    ],
    appKey: [
      { required: true, message: 'appKey 不能为空', trigger: 'blur' }
    ],
    corpId: [
      { required: true, message: 'corpId 不能为空', trigger: 'blur' }
    ],
  }
})

const typeOptions = ref([
  {
    label: 'dingtalk',
    value: 'dingtalk'
  },
  {
    label: 'wecom',
    value: 'wecom'
  },
  {
    label: 'sms',
    value: 'sms'
  },
])

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
  getProviders().then((res:any) => {
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
    type: undefined,
    appSecret: undefined,
    agentId: undefined,
    appKey: undefined,
    corpId: undefined,
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
  getProvider(row.id).then((res: any) => {
    const {providerId: id, name, type,agentId, appSecret } = res
    switch (row.type) {
      case 'dingtalk':
        const { appKey } = res
        state.form.appKey = appKey
        break;
    
      case 'wecom':
        const { corpId } = res
        state.form.corpId = corpId
        break;
    
      default:
        break;
    }
    state.open = Status.EDIT

    nextTick(()=>{
      state.form.id = id
      state.form.name = name
      state.form.type = type
      state.form.agentId = agentId
      state.form.appSecret = appSecret
    })
  })
}

function handleConfig(row: any) {
  navigateTo(`/dashboard/${tenant.value}/providers/sms`)
}

let updateLoading = ref(false);
/** 提交按钮 */
function submitForm() {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      updateLoading.value = true
      let { id, name, type, agentId, appSecret } = state.form
      let params;
      switch (type) {
        case 'dingtalk':
          const { appKey } = state.form
          params = { name, type, agentId, appSecret, appKey }
          break;
      
        case 'wecom':
          const { corpId } = state.form
          params = { name, type, agentId, appSecret, corpId }
          break;
      
        case 'sms':
          params = { name, type }
          break;
      
        default:
          break;
      }

      if (state.open === Status.EDIT) {
        updateProvider(id as number, params).then(() => {
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
        saveProvider(params).then(() => {
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
    await delProvider(row.id)
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
})
</script>

<style lang="scss" scoped>
.option {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}
</style>
