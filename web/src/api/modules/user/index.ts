import axios from "@/api";

// 登录
export const loginAPI = (data: any) => {
  return axios({
    url: "/api/auth/login",
    method: "post",
    data
  });
};

// 获取用户信息
export const getUserInfoAPI = (params?: any) => {
  return axios({
    url: "/api/auth/info",
    method: "get",
    params
  });
};

// 安全设置
export const securityAPI = (data: any) => {
  return axios({
    url: "/api/auth/security",
    method: "post",
    data
  });
};

// 设置密码
export const setPasswordAPI = (data: any) => {
  return axios({
    url: "/api/auth/set_password",
    method: "post",
    data
  });
};
