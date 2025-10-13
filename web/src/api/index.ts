import axios from "axios";

import { Message } from "@arco-design/web-vue";
import { useUserInfoStore } from "@/store/modules/user-info";
import pinia from "@/store/index";

declare module "axios" {
  interface AxiosInstance {
    (config: AxiosRequestConfig): Promise<any>;
  }
}

// 创建axios实例
const service = axios.create();

// 请求拦截器
service.interceptors.request.use(
  function (config: any) {
    // 发送请求之前做什么
    // 获取token鉴权
    let userInfo: any = {};
    if (localStorage.getItem("user-info")) {
      userInfo = JSON.parse(localStorage.getItem("user-info") as string);
    }
    if (userInfo?.token) {
      // 有token，在请求头中携带token
      config.headers.Authorization = userInfo.token;
    }
    return config;
  },
  function (error: any) {
    // 请求错误
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  function (response: any) {
    if (response.status != 200) {
      Message.error("服务器异常，请联系管理员");

      return Promise.reject(response.data);
    }

    let res = response.data;
    if (res.code == 400) {
      Message.error(res.msg);

      return Promise.reject(res);
    }

    if (res.code == 403) {
      // 清除 localStorage
      localStorage.removeItem("user-info");

      // 清除 Pinia store 中的 token 和用户信息
      const userStore = useUserInfoStore(pinia);
      userStore.logOut();

      // 跳转到登录页（保留当前路径前缀，如 /admin）
      const basePath = window.location.pathname.split("#")[0];
      window.location.href = `${basePath}#/login`;

      return Promise.reject(res);
    }
    if (res.code == 404) {
      Message.error("请求连接超时");

      return Promise.reject(res);
    }
    if (res.code != 200) {
      Message.error(res.message);

      return Promise.reject(res);
    }

    // 返回数据
    return Promise.resolve(res);
  },
  function (error: any) {
    // 处理 HTTP 错误状态码
    if (error.response) {
      if (error.response.status === 403) {
        localStorage.removeItem("user-info");
        const userStore = useUserInfoStore(pinia);
        userStore.logOut();

        Message.error("登录已过期，请重新登录");

        const basePath = window.location.pathname.split("#")[0];

        return (window.location.href = `${basePath}#/login`);
      }

      return Message.error(error.message || "请求失败");
    }
  }
);

export default service;
