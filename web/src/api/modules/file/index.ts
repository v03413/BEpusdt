import axios from "@/api";

// 获取机构树
export const getDocumentLibraryTreeAPI = () => {
  return axios({
    url: "/mock/document-library/tree",
    method: "get"
  });
};

// 获取列表数据
export const getDocumentLibraryTableAPI = () => {
  return axios({
    url: "/mock/document-library/table",
    method: "get"
  });
};
