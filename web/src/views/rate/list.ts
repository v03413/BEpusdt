interface List {
  id: number;
  fiat: string;
  crypto: string;
  rate: number;
  raw_rate: number;
  created_at: string;
  key?: string;
}

interface FormData {
  form: {
    fiat: string;
    crypto: string;
    datetime: string[];
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
