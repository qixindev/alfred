import { ElMessage } from "element-plus";
import { getUser } from "~~/api/common";
interface SelectOption {
  name: string;
  id: number;
}
export const userTenant=ref<any>([])

// 更新列表
export const setTenants = () => {
  getUser()
    .then((res: any) => {
      if (!res) {
        ElMessage({
          message: "当前没有租户，请创建租户",
          type: "error",
        });
      }
      userTenant.value=[...res]
    })
    .finally(() => { });
}