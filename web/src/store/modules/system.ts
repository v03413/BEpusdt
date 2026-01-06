import { defineStore } from "pinia";
import persistedstateConfig from "@/store/config/index";
import { getDictAPI } from "@/api/modules/system/index";
/**
 * 用户信息
 * @methods setAccount 设置账号信息
 * @methods setToken 设置token
 * @methods logOut 退出登录
 */
const systemStore = () => {
  // 字典数据
  const dict = ref<any>([]);

  // 设置字典数据
  async function setDictData() {
    let dictData = await getDictAPI();
    dict.value = dictData.data || [];
  }

  return { dict, setDictData };
};

export const useSystemStore = defineStore("system", systemStore, {
  persist: persistedstateConfig("system", ["dict"])
});
