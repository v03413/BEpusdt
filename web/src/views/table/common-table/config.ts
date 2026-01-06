interface List {
  key: string;
  name: string;
  avatar: string;
  phone: string;
  email: string;
  address: string;
  status: number;
  createTime: string;
}

interface FormData {
  form: {
    name: string;
    phone: string;
    email: string;
    address: string;
    status: number | null;
  };
  search: boolean;
}

interface RowSelection {
  type: string;
  showCheckedAll: boolean;
  onlyCurrent: boolean;
}

interface Pagination {
  showPageSize: boolean;
  showTotal: boolean;
  current: number;
  pageSize: number;
  total: number;
}

export type { List, FormData, RowSelection, Pagination };
