import { ref } from "vue";

export interface Detail {
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

export const useOrderDetail = () => {
  const detailVisible = ref(false);
  const detailData = ref<Detail>({
    id: 0,
    order_id: "",
    trade_id: "",
    trade_type: "",
    rate: "",
    amount: "",
    money: 0,
    address: "",
    from_address: "",
    status: 0,
    name: "",
    api_type: "",
    return_url: "",
    notify_url: "",
    notify_num: 0,
    notify_state: 0,
    ref_hash: "",
    ref_block_num: 0,
    expired_at: ""
  });

  // 显示详情
  const showDetail = (record: Detail) => {
    detailData.value = { ...record };
    detailVisible.value = true;
  };

  // 关闭详情
  const closeDetail = () => {
    detailVisible.value = false;
    detailData.value = {
      id: 0,
      order_id: "",
      trade_id: "",
      trade_type: "",
      rate: "",
      amount: "",
      money: 0,
      address: "",
      from_address: "",
      status: 0,
      name: "",
      api_type: "",
      return_url: "",
      notify_url: "",
      notify_num: 0,
      notify_state: 0,
      ref_hash: "",
      ref_block_num: 0,
      expired_at: ""
    };
  };

  return {
    detailVisible,
    detailData,
    showDetail,
    closeDetail
  };
};
