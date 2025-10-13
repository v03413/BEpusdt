interface SelectedKeys {
  type: string;
  showCheckedAll: boolean;
  onlyCurrent: boolean;
}

interface Form {
  name: string;
}

interface Pagination {
  showPageSize: boolean;
  showTotal: boolean;
  current: number;
  pageSize: number;
  total: number;
}

interface List {
  key: string;
  source: string;
  sourceSvg: string;
  sourceType: string;
  flow: string;
  percent: string;
  status: number;
  createTime: string;
}

export type { SelectedKeys, Form, Pagination, List };
