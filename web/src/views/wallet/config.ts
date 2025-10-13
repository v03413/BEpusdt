interface List {
  id: number;
  name: string;
  address: string;
  status: number;
  createTime: string;
  trade_type?: string;
  remark?: string;
  other_notify?: number;
}

interface FormData {
  form: {
    name: string;
    address: string;
    trade_type: string;
  };
  search: boolean;
}

interface AddForm {
  name: string;
  address: string;
  trade_type: string;
  remark: string;
  other_notify: number;
}

interface ModForm {
  id: number;
  name: string;
  status: number;
  address: string;
  trade_type: string;
  remark: string;
  other_notify: number;
}

interface Pagination {
  showPageSize: boolean;
  showTotal: boolean;
  current: number;
  pageSize: number;
  total: number;
}

export type { List, FormData, Pagination, AddForm, ModForm };
