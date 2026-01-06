interface List {
  key: string;
  serial: string;
  name: string;
  content: string;
  searchType: string;
  contentSize: number;
  status: number;
  createTime: string;
}

interface FormData {
  form: {
    serial: string;
    name: string;
    content: string;
    contentSize: string;
    status: number | null;
    createTime: string | null;
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

interface Columns {
  title: string;
  dataIndex?: string;
  checked: boolean;
  slotName?: string;
  align?: string;
}

interface DensityType {
  value: string;
  label: string;
}

export type { List, FormData, RowSelection, Pagination, Columns, DensityType };
