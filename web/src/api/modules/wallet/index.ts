import axios from "@/api";

export const getWalletListAPI = (data: any) => {
  return axios({
    url: "/api/wallet/list",
    method: "post",
    data
  });
};

export const delWalletAPI = (data: any) => {
  return axios({
    url: "/api/wallet/del",
    method: "post",
    data
  });
};

export const addWalletAPI = (data: any) => {
  return axios({
    url: "/api/wallet/add",
    method: "post",
    data
  });
};

export const modWalletAPI = (data: any) => {
  return axios({
    url: "/api/wallet/mod",
    method: "post",
    data
  });
};
