<template>
  <div>
    <div class="option">
      <el-button type="primary" icon="Plus" @click="handleAdd">新增Tenant</el-button>
    </div>
    <el-card>
      <el-table v-loading="loading" stripe :data="dataList">
        <el-table-column label="ID" align="center" prop="id" />
        <el-table-column label="name" align="center" prop="name" />
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
              @click="viewDevices(row)"
              >device管理
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
              icon="Delete"
              @click="handleDelete(row)"
              :loading="row.deleteLoading"
              >删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <!-- <el-pagination
        v-model:current-page="currentPage4"
        v-model:page-size="pageSize4"
        :page-sizes="[100, 200, 300, 400]"
        :small="small"
        :disabled="disabled"
        :background="background"
        layout="total, sizes, prev, pager, next, jumper"
        :total="400"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      /> -->
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
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="name" prop="name">
          <el-input v-model="form.name" placeholder="请输入name" />
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
import { getUserInfo } from "~/api/user";
import { getTenants, saveTenant, updateTenant, delTenant } from "~/api/tenant";
import { setTenants } from "~/composables/getUser";
const currentPage4 = ref(4);
const pageSize4 = ref(100);
const small = ref(false);
const background = ref(false);
const disabled = ref(false);
const tenant = useTenant();
interface Form {
  id: undefined | Number;
  name: undefined | string;
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
  } as Form,
  // 表单校验
  rules: {
    name: [{ required: true, message: "client name 不能为空", trigger: "blur" }],
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
  getTenants()
    .then((res: any) => {
      state.dataList = res;
    })
    .finally(() => {
      state.loading = false;
    });
  setTenants();
}

// 表单重置
function resetForm() {
  state.form = {
    id: undefined,
    name: undefined,
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
  const { id, name } = row;
  state.open = Status.EDIT;
  nextTick(() => {
    state.form = {
      id,
      name,
    };
  });
}

let updateLoading = ref(false);
/** 提交按钮 */
function submitForm() {
  formRef.value.validate((valid: boolean) => {
    if (valid) {
      updateLoading.value = true;
      let { id, name } = state.form;

      const params = { name };

      if (state.open === Status.EDIT) {
        updateTenant(id as number, params)
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
        getUserInfo().then((res: any) => {
          const info = { name, sub: res.sub };
          saveTenant(info)
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
      await delTenant(row.id);
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
// 分页
const handleSizeChange = (val: number) => {
  console.log(`${val} items per page`);
};
const handleCurrentChange = (val: number) => {
  console.log(`current page: ${val}`);
};
function viewDevices(row: any) {
  navigateTo(`/dashboard/${tenant.value}/device/${row.id}/groups`);
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
