interface List {
  id: number;
  order_id: string;
  trade_id: string;
  trade_type: string;
  rate: string;
  amount: string;
  money: number;
  address: string;
  from_address: string;
  status: number;
  name: string;
  api_type: string;
  return_url: string;
  notify_url: string;
  notify_num: number;
  notify_state: number;
  ref_hash: string;
  ref_block_num: number;
  expired_at: string;
  confirmed_at?: string;
  created_at?: string;
  updated_at?: string;
}

interface FormData {
  form: {
    order_id: string;
    trade_id: string;
    trade_type: string;
    address: string;
    status?: number;
    createTime: string[];
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
