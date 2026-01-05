import axios from "@/api";

const setConfAPI = (data: any) => {
  return axios({
    url: "/api/conf/set",
    method: "post",
    data
  });
};

const getConfAPI = (data: any) => {
  return axios({
    url: "/api/conf/get",
    method: "post",
    data
  });
};

const getsConfAPI = (data: any) => {
  return axios({
    url: "/api/conf/gets",
    method: "post",
    data
  });
};

const setsConfAPI = (data: any) => {
  return axios({
    url: "/api/conf/sets",
    method: "post",
    data
  });
};

const notifierAPI = (data: any) => {
  return axios({
    url: "/api/conf/notifier",
    method: "post",
    data
  });
};

const notifierTestAPI = (data: any) => {
  return axios({
    url: "/api/conf/notifier_test",
    method: "post",
    data
  });
};

const resetApiAuthToken = (data: any) => {
  return axios({
    url: "/api/conf/reset_api_auth_token",
    method: "post",
    data
  });
};

export { setConfAPI, getConfAPI, getsConfAPI, setsConfAPI, notifierAPI, notifierTestAPI, resetApiAuthToken };
