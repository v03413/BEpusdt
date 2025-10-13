import { defineStore } from "pinia";
import persistedstateConfig from "@/store/config/index";
import { getUserInfoAPI } from "@/api/modules/user/index";

interface Account {
  user: any; // 用户信息
  roles: string[]; // 角色
  permissions: string[]; // 权限
}

/**
 * 用户信息
 * @methods setAccount 设置账号信息
 * @methods setToken 设置token
 * @methods logOut 退出登录
 */
const userInfoStore = () => {
  // 账号信息
  const account = ref<Account>({
    user: {}, // 用户信息
    roles: [], // 角色
    permissions: [] // 权限
  });

  const trade_type = ref<Record<string, string>>({});
  const trade_fiat = ref<string[]>([]);
  const trade_crypto = ref<string[]>([]);
  const admin_username = ref<string>("");

  async function setAccount() {
    let data = await getUserInfoAPI();

    // 确保返回的数据有效
    if (data && data.data) {
      trade_type.value = data.data.trade_type || {};
      trade_fiat.value = data.data.trade_fiat || [];
      trade_crypto.value = data.data.trade_crypto || [];
      admin_username.value = data.data.admin_username || "";
    }
  }

  // 设置token
  const token = ref<string>("");
  async function setToken(data: string) {
    token.value = data;
  }

  // 退出登录
  async function logOut() {
    // 清除账号数据
    account.value = {
      user: {},
      roles: [],
      permissions: []
    };
    token.value = "";
  }

  return { account, token, setAccount, setToken, logOut, trade_type, trade_fiat, trade_crypto, admin_username };
};

export const useUserInfoStore = defineStore("user-info", userInfoStore, {
  persist: persistedstateConfig("user-info", ["token", "account", "trade_type", "trade_fiat", "trade_crypto", "admin_username"])
});
