import axios from "@/api";

export const mockTestAPI = () => {
  return axios({
    url: "/mock/test/success",
    method: "get"
  });
};
